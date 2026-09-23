package main

// AUTHORED host composition of the existing app, verified inbox, approval and
// outbound fence owners. Provider/runtime secrets are external files; none are
// placed in profiles, URLs, logs or process arguments.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"elite.local/enterprise/internal/app"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/whatsappbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errWhatsAppHost = errors.New("WhatsApp host configuration, authenticated identity or reporting unavailable")

type whatsappHostConfig struct {
	CampaignPolicyFile           string                     `json:"campaign_policy_file,omitempty"`
	CampaignPolicySHA256         string                     `json:"campaign_policy_sha256,omitempty"`
	Schema                       string                     `json:"schema"`
	ReferenceFixture             bool                       `json:"reference_fixture"`
	TenantID                     string                     `json:"tenant_id"`
	TenantCode                   string                     `json:"tenant_code"`
	OrganizationID               string                     `json:"organization_id"`
	ConnectionID                 string                     `json:"connection_id"`
	Purpose                      string                     `json:"consent_purpose"`
	PolicyVersion                string                     `json:"policy_version"`
	RetentionApprovalSHA256      string                     `json:"retention_approval_sha256"`
	ProviderProfileFile          string                     `json:"provider_profile_file"`
	ProviderProfileSHA256        string                     `json:"provider_profile_sha256"`
	Process                      whatsappbridge.Process     `json:"process"`
	ReconcilerSHA256             string                     `json:"reconciler_sha256"`
	WorkerID                     string                     `json:"worker_id"`
	PollSeconds                  int                        `json:"poll_seconds"`
	ConversationRetentionSeconds int                        `json:"conversation_retention_seconds"`
	ConversationMaxAttempts      int                        `json:"conversation_max_attempts"`
	LLMModel                     string                     `json:"llm_model"`
	LLMBaseURL                   string                     `json:"llm_base_url"`
	LLMTokenBudget               int64                      `json:"llm_token_budget"`
	DomainBaseURL                string                     `json:"domain_base_url"`
	ServiceKinds                 map[string]string          `json:"service_kinds"`
	ProductVariants              map[string]string          `json:"product_variants"`
	PriceBookID                  string                     `json:"price_book_id"`
	Conversation                 conversationruntime.Config `json:"conversation"`
	ApprovalPolicy               approval.Policy            `json:"approval_policy"`
	Collector                    struct {
		Address    string `json:"address"`
		ServerName string `json:"server_name"`
		CAFile     string `json:"ca_file"`
		CASHA256   string `json:"ca_sha256"`
	} `json:"status_collector"`
}

func readWhatsAppFile(path string, limit int64) ([]byte, error) {
	if !filepath.IsAbs(path) || limit < 1 {
		return nil, errWhatsAppHost
	}
	i, e := os.Lstat(path)
	if e != nil || !i.Mode().IsRegular() || i.Size() > limit {
		return nil, errWhatsAppHost
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, errWhatsAppHost
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil || int64(len(b)) > limit {
		return nil, errWhatsAppHost
	}
	return b, nil
}
func readWhatsAppExact(path, hash string, limit int64) ([]byte, error) {
	h, e := hex.DecodeString(hash)
	if e != nil || len(h) != 32 || strings.ToLower(hash) != hash {
		return nil, errWhatsAppHost
	}
	b, e := readWhatsAppFile(path, limit)
	if e != nil {
		return nil, e
	}
	sum := sha256.Sum256(b)
	if !bytes.Equal(h, sum[:]) {
		return nil, errWhatsAppHost
	}
	return b, nil
}

// Rotating a credential file takes effect on the next request/poll. Errors do
// not disclose its contents or path. The deployment owns private file ACLs.
type whatsappHostSecrets struct{ lookup func(string) string }

func (s whatsappHostSecrets) read(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil || s.lookup == nil {
		return "", errWhatsAppHost
	}
	b, e := readWhatsAppFile(s.lookup(key), 16384)
	if e != nil {
		return "", e
	}
	v := strings.TrimSpace(string(b))
	if v == "" || strings.ContainsAny(v, "\r\n\x00") {
		return "", errWhatsAppHost
	}
	return v, nil
}
func (s whatsappHostSecrets) WhatsAppToken(c context.Context) (string, error) {
	return s.read(c, "WHATSAPP_ACCESS_TOKEN_FILE")
}
func (s whatsappHostSecrets) WhatsAppAppSecret(c context.Context) (string, error) {
	return s.read(c, "WHATSAPP_APP_SECRET_FILE")
}
func (s whatsappHostSecrets) WhatsAppVerifyToken(c context.Context) (string, error) {
	return s.read(c, "WHATSAPP_VERIFY_TOKEN_FILE")
}

type whatsappAccessTokenSource interface {
	AccessToken(context.Context) (string, error)
}
type whatsappFileTokenSource struct{ secrets whatsappHostSecrets }

func (s whatsappFileTokenSource) AccessToken(ctx context.Context) (string, error) {
	return s.secrets.read(ctx, "WHATSAPP_SERVICE_TOKEN_FILE")
}

type whatsappHostIdentity struct {
	tokens               whatsappAccessTokenSource
	verifier             identity.Verifier
	tenant, organization string
}

func (s whatsappHostIdentity) verified(ctx context.Context) (string, identity.Principal, error) {
	var empty identity.Principal
	if s.tokens == nil || s.verifier == nil || s.tenant == "" || s.organization == "" {
		return "", empty, errWhatsAppHost
	}
	raw, e := s.tokens.AccessToken(ctx)
	if e != nil {
		return "", empty, e
	}
	p, e := s.verifier.Verify(ctx, raw)
	if e != nil || p.Subject == "" || p.TenantID != s.tenant || !p.AllowedOrganization(s.organization) || !p.Allowed("whatsapp:process") || !p.Allowed("appointment:manage") {
		return "", empty, errWhatsAppHost
	}
	return raw, p, nil
}
func (s whatsappHostIdentity) Resolve(ctx context.Context) (identity.Principal, error) {
	_, p, e := s.verified(ctx)
	return p, e
}
func (s whatsappHostIdentity) AccessToken(ctx context.Context) (string, error) {
	raw, _, e := s.verified(ctx)
	return raw, e
}

// Select exactly one external identity lifecycle owner. A token file is for an
// existing host token agent. The OIDC broker renews client-credentials grants in
// process, verifies them with the pinned SDK and rereads the secret on renewal.
func selectWhatsAppIdentity(ctx context.Context, c whatsappHostConfig, verifier identity.Verifier, lookup func(string) string) (whatsappHostIdentity, error) {
	var empty whatsappHostIdentity
	if lookup == nil || verifier == nil {
		return empty, errWhatsAppHost
	}
	file := lookup("WHATSAPP_SERVICE_TOKEN_FILE")
	profileFile := lookup("WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE")
	profileHash := lookup("WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256")
	secretFile := lookup("WHATSAPP_SERVICE_CLIENT_SECRET_FILE")
	var tokens whatsappAccessTokenSource
	if file != "" {
		if profileFile != "" || profileHash != "" || secretFile != "" {
			return empty, errWhatsAppHost
		}
		tokens = whatsappFileTokenSource{whatsappHostSecrets{lookup}}
	} else {
		if !filepath.IsAbs(profileFile) || !filepath.IsAbs(secretFile) {
			return empty, errWhatsAppHost
		}
		raw, err := readWhatsAppExact(profileFile, profileHash, 16384)
		if err != nil {
			return empty, errWhatsAppHost
		}
		profile, err := identity.LoadServiceTokenProfile(raw, profileHash)
		if err != nil {
			return empty, errWhatsAppHost
		}
		binding := profile.Binding()
		if binding.Transport == "LOOPBACK_FIXTURE" && !c.ReferenceFixture {
			return empty, errWhatsAppHost
		}
		if binding.TenantID != c.TenantID || !slices.Contains(binding.Organizations, c.OrganizationID) || !slices.Contains(binding.Permissions, "whatsapp:process") || !slices.Contains(binding.Permissions, "appointment:manage") {
			return empty, errWhatsAppHost
		}
		broker, err := identity.NewServiceTokenBroker(ctx, profile, identity.ServiceSecretFile(secretFile))
		if err != nil {
			return empty, errWhatsAppHost
		}
		tokens = broker
	}
	return whatsappHostIdentity{tokens: tokens, verifier: verifier, tenant: c.TenantID, organization: c.OrganizationID}, nil
}

type whatsappHost struct {
	extraRun    func(context.Context) error
	modules     []httpapi.EnterpriseModule
	ingress     *whatsappbridge.WebhookReceiver
	worker      *whatsappbridge.StatusWorker
	identity    whatsappHostIdentity
	reporter    *whatsappbridge.JSONStatusReporter
	interval    time.Duration
	application *app.App
}

func (h *whatsappHost) Register(mux *http.ServeMux, v identity.Verifier) {
	for _, m := range h.modules {
		m.Register(mux, v)
	}
	mux.Handle("/v1/providers/whatsapp/webhook", h.ingress)
}
func (h *whatsappHost) run(ctx context.Context) error {
	if h.extraRun == nil {
		return h.worker.Run(ctx, h.identity, h.reporter, h.interval)
	}
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan error, 2)
	go func() { done <- h.worker.Run(child, h.identity, h.reporter, h.interval) }()
	go func() { done <- h.extraRun(child) }()
	first := <-done
	cancel()
	second := <-done
	return errors.Join(first, second)
}
func (h *whatsappHost) close() error { return h.reporter.Close() }

func init() {
	if whatsappRuntimeFactory != nil {
		panic("duplicate WhatsApp host owner")
	}
	whatsappRuntimeFactory = func(ctx context.Context, pool *pgxpool.Pool, verifier identity.Verifier, lookup func(string) string) (whatsappRuntime, error) {
		h, err := prepareWhatsAppHost(ctx, pool, verifier, lookup)
		if err != nil {
			return nil, err
		}
		if h == nil {
			return nil, nil
		}
		return h, nil
	}
}

func prepareWhatsAppHost(ctx context.Context, pool *pgxpool.Pool, verifier identity.Verifier, lookup func(string) string) (*whatsappHost, error) {
	if lookup == nil {
		return nil, errWhatsAppHost
	}
	switch lookup("WHATSAPP_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errWhatsAppHost
	}
	if pool == nil || verifier == nil {
		return nil, errWhatsAppHost
	}
	raw, e := readWhatsAppExact(lookup("WHATSAPP_HOST_PROFILE_FILE"), lookup("WHATSAPP_HOST_PROFILE_SHA256"), 65536)
	if e != nil {
		return nil, e
	}
	var c whatsappHostConfig
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&c) != nil || dec.Decode(new(any)) != io.EOF || c.Schema != "elite-whatsapp-host/v1" || c.TenantID == "" || c.OrganizationID == "" || c.ConnectionID == "" || c.PollSeconds < 1 || c.PollSeconds > 60 || c.ConversationRetentionSeconds < 1 || c.ConversationMaxAttempts < 1 {
		return nil, errWhatsAppHost
	}
	profile, e := readWhatsAppExact(c.ProviderProfileFile, c.ProviderProfileSHA256, 32768)
	if e != nil {
		return nil, e
	}
	if _, e = readWhatsAppExact(c.Process.PythonExecutable, c.Process.PythonSHA256, 128<<20); e != nil {
		return nil, e
	}
	if _, e = readWhatsAppExact(filepath.Join(c.Process.AdapterDirectory, "whatsapp_cloud.py"), c.Process.AdapterSHA256, 1<<20); e != nil {
		return nil, e
	}
	if _, e = readWhatsAppExact(filepath.Join(c.Process.AdapterDirectory, "status_reconciliation.py"), c.ReconcilerSHA256, 1<<20); e != nil {
		return nil, e
	}
	if e = validateWhatsAppProviderProfile(ctx, c.Process, profile, c.ProviderProfileSHA256); e != nil {
		return nil, e
	}
	secrets := whatsappHostSecrets{lookup}
	serviceIdentity, e := selectWhatsAppIdentity(ctx, c, verifier, lookup)
	if e != nil {
		return nil, e
	}
	// Startup validates the current service identity; subsequent polls repeat it.
	if _, e = serviceIdentity.Resolve(ctx); e != nil {
		return nil, e
	}
	for _, key := range []string{"WHATSAPP_ACCESS_TOKEN_FILE", "WHATSAPP_APP_SECRET_FILE", "WHATSAPP_VERIFY_TOKEN_FILE"} {
		if _, e = secrets.read(ctx, key); e != nil {
			return nil, e
		}
	}
	keyText, e := secrets.read(ctx, "WHATSAPP_CONTACT_HMAC_KEY_HEX_FILE")
	if e != nil {
		return nil, e
	}
	key, e := hex.DecodeString(keyText)
	if e != nil || len(key) < 32 || len(key) > 128 {
		return nil, errWhatsAppHost
	}
	if (c.CampaignPolicyFile != "" || c.CampaignPolicySHA256 != "") && lookup("WHATSAPP_SCHEDULE_ENABLED") != "true" {
		return nil, errWhatsAppHost
	}
	base, e := whatsappbridge.NewPostgresAppointmentApprovals(pool, key, c.Purpose, c.PolicyVersion)
	if e != nil {
		return nil, errWhatsAppHost
	}
	replies, e := whatsappbridge.NewPostgresReplyApprovals(base, c.TenantID, c.OrganizationID, c.ConnectionID, profile)
	if e != nil {
		return nil, errWhatsAppHost
	}
	sender := &whatsappbridge.Sender{TenantID: c.TenantID, Profile: profile, Approvals: replies, Tokens: secrets, Process: c.Process}
	fence, e := postgres.NewOutboundDeliveryStore(pool, key, 2*time.Minute)
	if e != nil {
		return nil, errWhatsAppHost
	}
	receiver := whatsappbridge.VerifiedInboxChannel{}
	fenced := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiver, Sender: sender, Store: fence}
	registry := channels.NewRegistry()
	if registry.Register(fenced) != nil {
		return nil, errWhatsAppHost
	}
	store, e := postgres.NewConversationStore(pool, 2*time.Minute, time.Duration(c.ConversationRetentionSeconds)*time.Second, c.ConversationMaxAttempts)
	if e != nil {
		return nil, errWhatsAppHost
	}
	contacts, e := postgres.NewContactIdentityStore(pool, key)
	if e != nil {
		return nil, errWhatsAppHost
	}
	llmKey, e := secrets.read(ctx, "LLM_API_KEY_FILE")
	if e != nil {
		return nil, e
	}
	a, e := app.New(app.Config{LLMAPIKey: llmKey, LLMModel: c.LLMModel, LLMBaseURL: c.LLMBaseURL, DomainBaseURL: c.DomainBaseURL, DomainTokenProvider: serviceIdentity, TenantID: c.TenantID, TenantCode: c.TenantCode, ServiceKinds: c.ServiceKinds, ProductVariants: c.ProductVariants, PriceBookID: c.PriceBookID, ConversationStore: store, ContactResolver: whatsappbridge.ScopedContactResolver{TenantID: c.TenantID, OrganizationID: c.OrganizationID, Resolver: contacts}, ChannelRegistry: registry, Conversation: c.Conversation, LLMTokenBudget: c.LLMTokenBudget, ApprovalPolicy: c.ApprovalPolicy})
	if e != nil {
		return nil, errWhatsAppHost
	}
	contactDomain, e := postgres.NewConversationDomain(pool, a.Gateway, c.TenantID, c.OrganizationID)
	if e != nil {
		return nil, errWhatsAppHost
	}
	a.Conversation.Domain = contactDomain
	// Dispatcher is deliberately absent from the host execution path. Runtime
	// completion becomes a proposal; a separate authenticated human approves it.
	observer := &whatsappbridge.StatusObserver{TenantID: c.TenantID, Profile: profile, Process: c.Process, ReconcilerSHA256: c.ReconcilerSHA256, Secrets: secrets, Approvals: base}
	router, e := whatsappbridge.NewStatusRouter(observer, c.ConnectionID, c.RetentionApprovalSHA256, key, 8)
	if e != nil {
		return nil, errWhatsAppHost
	}
	if router.EnableConversation(whatsappbridge.ConversationRoute{OrganizationID: c.OrganizationID, Runtime: a.Conversation, Proposals: replies}) != nil {
		return nil, errWhatsAppHost
	}
	worker, e := whatsappbridge.NewStatusWorker(router, c.WorkerID, 2*time.Minute, 5*time.Second)
	if e != nil {
		return nil, errWhatsAppHost
	}
	module, e := whatsappbridge.NewReplyModule(replies, sender, fence)
	if e != nil {
		return nil, errWhatsAppHost
	}
	appointments, e := whatsappbridge.NewAppointmentNotificationModule(base, sender, fence, receiver)
	if e != nil {
		return nil, errWhatsAppHost
	}
	ingress, e := whatsappbridge.NewWebhookReceiver(whatsappbridge.WebhookReceiverConfig{TenantID: c.TenantID, ConnectionID: c.ConnectionID, RetentionApprovalSHA256: c.RetentionApprovalSHA256, Profile: profile, Process: c.Process, Secrets: secrets, Verification: secrets, Store: postgres.NewProviderIntegration(pool), MaxConcurrent: 4})
	if e != nil {
		return nil, errWhatsAppHost
	}
	reporter, e := dialWhatsAppReporter(ctx, c, lookup)
	if e != nil {
		return nil, e
	}
	h := &whatsappHost{modules: []httpapi.EnterpriseModule{appointments, module}, ingress: ingress, worker: worker, identity: serviceIdentity, reporter: reporter, interval: time.Duration(c.PollSeconds) * time.Second, application: a}
	if lookup("WHATSAPP_SCHEDULE_ENABLED") != "" && lookup("WHATSAPP_SCHEDULE_ENABLED") != "false" {
		if lookup("WHATSAPP_SCHEDULE_ENABLED") != "true" || whatsappScheduleFactory == nil {
			reporter.Close()
			return nil, errWhatsAppHost
		}
		if e = whatsappScheduleFactory(ctx, h, c, pool, base, sender, fence); e != nil {
			reporter.Close()
			return nil, e
		}
	}
	return h, nil
}

func dialWhatsAppReporter(ctx context.Context, c whatsappHostConfig, lookup func(string) string) (*whatsappbridge.JSONStatusReporter, error) {
	if c.Collector.Address == "" || c.Collector.ServerName == "" {
		return nil, errWhatsAppHost
	}
	ca, e := readWhatsAppExact(c.Collector.CAFile, c.Collector.CASHA256, 1<<20)
	if e != nil {
		return nil, e
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(ca) {
		return nil, errWhatsAppHost
	}
	cert, e := readWhatsAppFile(lookup("WHATSAPP_COLLECTOR_CLIENT_CERT_FILE"), 1<<20)
	if e != nil {
		return nil, e
	}
	key, e := readWhatsAppFile(lookup("WHATSAPP_COLLECTOR_CLIENT_KEY_FILE"), 1<<20)
	if e != nil {
		return nil, e
	}
	pair, e := tls.X509KeyPair(cert, key)
	if e != nil {
		return nil, errWhatsAppHost
	}
	dialer := tls.Dialer{NetDialer: &net.Dialer{Timeout: 5 * time.Second}, Config: &tls.Config{MinVersion: tls.VersionTLS13, ServerName: c.Collector.ServerName, RootCAs: roots, Certificates: []tls.Certificate{pair}}}
	conn, e := dialer.DialContext(ctx, "tcp", c.Collector.Address)
	if e != nil {
		return nil, errWhatsAppHost
	}
	reporter, e := whatsappbridge.NewJSONStatusReporter(conn)
	if e != nil {
		_ = conn.Close()
		return nil, errWhatsAppHost
	}
	return reporter, nil
}

type whatsappProfileOutput struct{ bytes.Buffer }

func (b *whatsappProfileOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 1024 {
		return 0, errWhatsAppHost
	}
	return b.Buffer.Write(p)
}
func validateWhatsAppProviderProfile(ctx context.Context, p whatsappbridge.Process, profile []byte, expected string) error {
	input, e := json.Marshal(struct {
		Schema  string `json:"schema"`
		Profile []byte `json:"profile"`
	}{"elite-whatsapp-validate-profile/v1", profile})
	if e != nil {
		return errWhatsAppHost
	}
	bounded, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	command := exec.CommandContext(bounded, p.PythonExecutable, "-I", "-B", filepath.Join(p.AdapterDirectory, "whatsapp_cloud.py"), "--validate-profile-bridge")
	command.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		command.Env = append(command.Env, "SystemRoot="+root)
	}
	command.Stdin = bytes.NewReader(input)
	command.Stderr = io.Discard
	command.WaitDelay = time.Second
	var output whatsappProfileOutput
	command.Stdout = &output
	if command.Run() != nil {
		return errWhatsAppHost
	}
	var result struct {
		Schema string `json:"schema"`
		SHA256 string `json:"profile_sha256"`
	}
	d := json.NewDecoder(bytes.NewReader(output.Bytes()))
	d.DisallowUnknownFields()
	if d.Decode(&result) != nil || d.Decode(new(any)) != io.EOF || result.Schema != "elite-whatsapp-profile-validation/v1" || result.SHA256 != expected {
		return errWhatsAppHost
	}
	return nil
}

var whatsappScheduleFactory func(context.Context, *whatsappHost, whatsappHostConfig, *pgxpool.Pool, *whatsappbridge.PostgresAppointmentApprovals, *whatsappbridge.Sender, *postgres.OutboundDeliveryStore) error
