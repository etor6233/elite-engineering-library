# Connected WhatsApp application host

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-WHATSAPP-HOST"
pack_version: "0.3.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Optional connected WhatsApp host using existing durable owners, exact profiles, authenticated reporting and renewable verified service identity; local infrastructure, no live credentials or production claim"
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-ELECTROMOBILITY-APPLICATION", "GO-APP-WIRING", "GO-OIDC-SERVICE-TOKEN-BROKER", "PYTHON-META-WHATSAPP-CLOUD-ADAPTER"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["Existing exact source and dependency records of the selected owners"]
verified_at: "2026-09-13"
```

## 2. Applicability

Select with the full franchise Go app, WhatsApp bridge and OIDC service broker. Without this pack the base application remains buildable with WhatsApp disabled.

## 3. Architecture contract

One existing app/runtime, current scoped contact resolver, existing approval and outbound fence, same inbox/jobs/status observer. Runtime output is a proposal; a different authenticated human approves it. A single selected identity source supplies both worker principal and domain bearer. TLS1.3 collector authentication; worker reporting failure stops the host. No business policy or upstream authorship invented.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/whatsapp.go
CREATE cmd/electromobility-api/whatsapp_test.go
CREATE cmd/electromobility-api/whatsapp_tls_test.go
CREATE cmd/electromobility-api/whatsapp_assembly_test.go
CREATE cmd/electromobility-api/whatsapp_identity_broker_test.go
CREATE docs/whatsapp-connected-host.md
CREATE docs/whatsapp-host-profile.template.json
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/whatsapp.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8c6d6e921949b055a1c0c3521037c2e36250115b1b2294eb847b252f044729e5"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `cmd/electromobility-api/whatsapp_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "207403570410f49c965a7b6f0428c2c21ff3d209a9c42821a8b21fdb44006b70"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"elite.local/enterprise/internal/platform/identity"
)

type hostWAIdentityVerifier struct {
	calls int
	p     identity.Principal
	token string
}

func (v *hostWAIdentityVerifier) Verify(_ context.Context, raw string) (identity.Principal, error) {
	v.calls++
	if raw != v.token {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return v.p, nil
}

func TestWhatsAppHostActivationNoImplicitInputs(t *testing.T) {
	for _, enabled := range []string{"", "false"} {
		calls := 0
		h, e := selectedWhatsAppHost(context.Background(), nil, nil, func(k string) string {
			calls++
			if k != "WHATSAPP_ENABLED" {
				t.Fatal("disabled host read inputs")
			}
			return enabled
		})
		if e != nil || h != nil || calls != 1 {
			t.Fatalf("%v %v calls%d", h, e, calls)
		}
	}
	for _, enabled := range []string{"yes", "TRUE", "true"} {
		if _, e := selectedWhatsAppHost(context.Background(), nil, nil, func(string) string { return enabled }); e == nil {
			t.Fatal("activation without actual dependencies")
		}
	}
}

func TestWhatsAppHostExactFilesAndSecretRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.json")
	raw := []byte("{\"mode\":\"fixture\"}\n")
	if e := os.WriteFile(path, raw, 0600); e != nil {
		t.Fatal(e)
	}
	h := sha256.Sum256(raw)
	if b, e := readWhatsAppExact(path, hex.EncodeToString(h[:]), 64); e != nil || string(b) != string(raw) {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		path, hash string
		limit      int64
	}{{path, hex.EncodeToString(h[:]), 1}, {"relative", hex.EncodeToString(h[:]), 64}, {dir, hex.EncodeToString(h[:]), 64}, {path, "00", 64}} {
		if _, e := readWhatsAppExact(tc.path, tc.hash, tc.limit); e == nil {
			t.Fatal("invalid exact file")
		}
	}
	tokenFile := filepath.Join(dir, "rotating-token")
	if e := os.WriteFile(tokenFile, []byte("first\n"), 0600); e != nil {
		t.Fatal(e)
	}
	s := whatsappHostSecrets{func(string) string { return tokenFile }}
	if v, e := s.WhatsAppToken(context.Background()); e != nil || v != "first" {
		t.Fatal(v, e)
	}
	if e := os.WriteFile(tokenFile, []byte("second\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if v, e := s.WhatsAppToken(context.Background()); e != nil || v != "second" {
		t.Fatal(v, e)
	}
	if e := os.WriteFile(tokenFile, []byte("two\nlines"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := s.WhatsAppToken(context.Background()); e == nil {
		t.Fatal("multiline secret")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := s.WhatsAppToken(ctx); e == nil {
		t.Fatal("cancelled secret read")
	}
}

func TestWhatsAppHostIdentityRevalidatedAndScoped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if e := os.WriteFile(path, []byte("current"), 0600); e != nil {
		t.Fatal(e)
	}
	p := identity.Principal{Subject: "worker", TenantID: "tenant", Organizations: map[string]struct{}{"org": {}}, Permissions: map[string]struct{}{"appointment:manage": {}, "whatsapp:process": {}}}
	v := &hostWAIdentityVerifier{p: p, token: "current"}
	s := whatsappHostIdentity{whatsappFileTokenSource{whatsappHostSecrets{func(string) string { return path }}}, v, "tenant", "org"}
	if _, e := s.Resolve(context.Background()); e != nil {
		t.Fatal(e)
	}
	if raw, e := s.AccessToken(context.Background()); e != nil || raw != "current" {
		t.Fatal(e)
	}
	if v.calls != 2 {
		t.Fatal("identity cached")
	}
	delete(v.p.Permissions, "whatsapp:process")
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("revoked grant")
	}
	v.p.Permissions["whatsapp:process"] = struct{}{}
	v.p.TenantID = "other"
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("wrong tenant")
	}
	v.p.TenantID = "tenant"
	delete(v.p.Organizations, "org")
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("wrong org")
	}
	v.p.Organizations["org"] = struct{}{}
	if e := os.WriteFile(path, []byte("expired"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("token verification ignored")
	}
}
````

### FILE: `cmd/electromobility-api/whatsapp_tls_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed53d408dd4f00e9159fdaee3d0568d4f86c581eb47fcb55e311e25edcb26b23"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"elite.local/enterprise/internal/whatsappbridge"
)

func whatsappTLSFixture(t *testing.T, linesPerConnection ...int) (whatsappHostConfig, map[string]string, <-chan []byte) {
	t.Helper()
	key, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		t.Fatal(e)
	}
	now := time.Now()
	template := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "local-fixture"}, DNSNames: []string{"collector.local"}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}}
	der, e := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if e != nil {
		t.Fatal(e)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, e := x509.MarshalPKCS8PrivateKey(key)
	if e != nil {
		t.Fatal(e)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	pair, e := tls.X509KeyPair(certPEM, keyPEM)
	if e != nil {
		t.Fatal(e)
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(certPEM)
	listener, e := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{MinVersion: tls.VersionTLS13, Certificates: []tls.Certificate{pair}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: roots})
	if e != nil {
		t.Fatal(e)
	}
	stop := make(chan struct{})
	output := make(chan []byte, 4)
	go func() {
		defer close(stop)
		for {
			conn, e := listener.Accept()
			if e != nil {
				return
			}
			func(conn net.Conn) {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
				lines := 1
				if len(linesPerConnection) == 1 {
					lines = linesPerConnection[0]
				}
				reader := bufio.NewReader(conn)
				for i := 0; i < lines; i++ {
					line, e := reader.ReadBytes('\n')
					if e != nil {
						return
					}
					output <- line
				}
			}(conn)
		}
	}()
	t.Cleanup(func() { _ = listener.Close(); <-stop })
	dir := t.TempDir()
	ca := filepath.Join(dir, "ca.pem")
	cert := filepath.Join(dir, "client.pem")
	private := filepath.Join(dir, "client-key.pem")
	for p, b := range map[string][]byte{ca: certPEM, cert: certPEM, private: keyPEM} {
		if e := os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	var c whatsappHostConfig
	c.Collector.Address = listener.Addr().String()
	c.Collector.ServerName = "collector.local"
	c.Collector.CAFile = ca
	h := sha256.Sum256(certPEM)
	c.Collector.CASHA256 = hex.EncodeToString(h[:])
	return c, map[string]string{"WHATSAPP_COLLECTOR_CLIENT_CERT_FILE": cert, "WHATSAPP_COLLECTOR_CLIENT_KEY_FILE": private}, output
}

func TestWhatsAppHostAuthenticatedReporter(t *testing.T) {
	c, values, out := whatsappTLSFixture(t)
	lookup := func(k string) string { return values[k] }
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, e := dialWhatsAppReporter(ctx, c, lookup)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Close()
	if e = r.Report(ctx, whatsappbridge.StatusPollReport{Outcome: "COMPLETED", Claimed: true, InsertedObservations: 1, NextDelay: time.Second}); e != nil {
		t.Fatal(e)
	}
	select {
	case line := <-out:
		var v map[string]any
		if json.Unmarshal(line, &v) != nil || v["outcome"] != "COMPLETED" || len(v) != 6 {
			t.Fatal(string(line))
		}
	case <-ctx.Done():
		t.Fatal("no authenticated report")
	}
	c.Collector.ServerName = "wrong.local"
	if r, e := dialWhatsAppReporter(ctx, c, lookup); e == nil {
		r.Close()
		t.Fatal("untrusted collector identity")
	}
	c.Collector.ServerName = "collector.local"
	values["WHATSAPP_COLLECTOR_CLIENT_KEY_FILE"] = ""
	if r, e := dialWhatsAppReporter(ctx, c, lookup); e == nil {
		r.Close()
		t.Fatal("missing client key")
	}
}
````

### FILE: `cmd/electromobility-api/whatsapp_assembly_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b53f80bb2e1a460914e2fe2877f0441a91adaf822dcf899faf1b1734aac78a12"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/whatsappbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This gate exercises the real constructor and fixed Python profile validator,
// plus real mutual TLS. It deliberately does not claim PostgreSQL statements,
// LLM calls or provider traffic; those have their own connected fixture gate.
func TestWhatsAppHostAssemblesExistingOwners(t *testing.T) {
	python := os.Getenv("ELITE_WHATSAPP_HOST_PYTHON")
	if python == "" {
		t.Skip("explicit pinned Python runtime is required for this constructor gate")
	}
	dir := t.TempDir()
	hashFile := func(p string) string {
		t.Helper()
		b, e := os.ReadFile(p)
		if e != nil {
			t.Fatal(e)
		}
		h := sha256.Sum256(b)
		return hex.EncodeToString(h[:])
	}
	write := func(name string, b []byte) string {
		t.Helper()
		p := filepath.Join(dir, name)
		if e := os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
		return p
	}
	adapter, e := filepath.Abs("../../whatsapp_cloud")
	if e != nil {
		t.Fatal(e)
	}
	provider := map[string]any{"provider": "Meta WhatsApp Cloud API", "source_id": "meta-whatsapp-api-examples", "source_commit": "de70ee908a67026e642aaee3703d20464e2a9466", "decision": "PROVEN", "business_account_id": "987654321", "graph_api_version": "v99.0", "phone_number_id": "123456789", "access_token_environment_variable": "WHATSAPP_ACCESS_TOKEN", "app_secret_environment_variable": "META_APP_SECRET", "verify_token_environment_variable": "WHATSAPP_VERIFY_TOKEN", "approved_templates": []any{map[string]any{"name": "order_update", "language_code": "es_AR", "body_parameter_count": 2}}, "data_retention": "synthetic local fixture only", "automatic_business_write": false}
	for _, k := range []string{"official_source_license_accepted", "platform_terms_accepted", "business_account_proven", "app_registration_proven", "phone_number_id_proven", "message_template_approval_proven", "webhook_subscription_proven", "test_recipient_consent_proven", "quota_and_cost_approved", "reconciliation_approved"} {
		provider[k] = true
	}
	providerRaw, _ := json.Marshal(provider)
	profilePath := write("synthetic-provider.json", providerRaw)
	c, values, _ := whatsappTLSFixture(t)
	c.Schema = "elite-whatsapp-host/v1"
	c.TenantID = "11111111-1111-4111-8111-111111111111"
	c.TenantCode = "fixture"
	c.OrganizationID = "fixture-org"
	c.ConnectionID = "fixture-whatsapp"
	c.Purpose = "fixture-conversation"
	c.PolicyVersion = "fixture-policy-v1"
	c.RetentionApprovalSHA256 = strings.Repeat("1", 64)
	c.ProviderProfileFile = profilePath
	c.ProviderProfileSHA256 = hashFile(profilePath)
	c.Process = whatsappbridge.Process{PythonExecutable: python, PythonSHA256: hashFile(python), AdapterDirectory: adapter, AdapterSHA256: hashFile(filepath.Join(adapter, "whatsapp_cloud.py")), EvidenceDirectory: dir}
	c.ReconcilerSHA256 = hashFile(filepath.Join(adapter, "status_reconciliation.py"))
	c.WorkerID = "fixture-worker"
	c.PollSeconds = 1
	c.ConversationRetentionSeconds = 3600
	c.ConversationMaxAttempts = 2
	c.LLMModel = "fixture-model"
	c.LLMBaseURL = "https://llm.example.invalid/v1"
	c.LLMTokenBudget = 10000
	c.DomainBaseURL = "https://domain.example.invalid"
	c.ServiceKinds = map[string]string{"fixture": "service"}
	c.ProductVariants = map[string]string{"fixture": "variant"}
	c.PriceBookID = "fixture-price-book"
	c.Conversation = conversationruntime.Config{Instructions: "Synthetic fixture only", PromptCacheKey: "fixture", MaxOutputTokens: 32, MaxHistory: 1, MaxInputBytes: 4096, MaxToolBytes: 4096, StoreApproved: true, HandoffText: "Fixture human handoff"}
	c.ApprovalPolicy = approval.Policy{MaxOpenPerSubject: 5}
	for _, k := range []string{"WHATSAPP_ACCESS_TOKEN_FILE", "WHATSAPP_APP_SECRET_FILE", "WHATSAPP_VERIFY_TOKEN_FILE", "WHATSAPP_SERVICE_TOKEN_FILE", "LLM_API_KEY_FILE"} {
		values[k] = write(k, []byte("synthetic-fixture-secret"))
	}
	values["WHATSAPP_CONTACT_HMAC_KEY_HEX_FILE"] = write("hmac", []byte(strings.Repeat("a1", 32)))
	values["WHATSAPP_ENABLED"] = "true"
	configRaw, _ := json.Marshal(c)
	values["WHATSAPP_HOST_PROFILE_FILE"] = write("host.json", configRaw)
	values["WHATSAPP_HOST_PROFILE_SHA256"] = hashFile(values["WHATSAPP_HOST_PROFILE_FILE"])
	p := identity.Principal{Subject: "fixture-service", TenantID: c.TenantID, Permissions: map[string]struct{}{"appointment:manage": {}, "whatsapp:process": {}}, Organizations: map[string]struct{}{c.OrganizationID: {}}}
	v := &hostWAIdentityVerifier{p: p, token: "synthetic-fixture-secret"}
	pool, e := pgxpool.New(context.Background(), "postgres://fixture:fixture@127.0.0.1:1/constructor_not_connected?sslmode=disable")
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	h, e := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] })
	if e != nil {
		t.Fatal(e)
	}
	defer h.close()
	if _, ok := h.(*whatsappHost).application.Conversation.Domain.(*postgres.ConversationDomain); !ok {
		t.Fatal("host did not install contact-bound domain status")
	}
	if h.(*whatsappHost).application == nil || h.(*whatsappHost).application.Conversation == nil || h.(*whatsappHost).worker == nil || h.(*whatsappHost).ingress == nil || len(h.(*whatsappHost).modules) != 2 {
		t.Fatal("missing existing owner")
	}
	mux := http.NewServeMux()
	h.Register(mux, v)
	for _, path := range []string{"/v1/franchise/whatsapp/replies", "/v1/providers/whatsapp/webhook"} {
		r := httptest.NewRequest("DELETE", path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code == 404 {
			t.Fatalf("unmounted %s", path)
		}
	}
	// The normal profile remains blocked; the constructor cannot activate it.
	provider["decision"] = "BLOCKED_ACCESS_TERMS_AND_RECONCILIATION_REQUIRED"
	blocked, _ := json.Marshal(provider)
	c.ProviderProfileFile = write("blocked-provider.json", blocked)
	c.ProviderProfileSHA256 = hashFile(c.ProviderProfileFile)
	configRaw, _ = json.Marshal(c)
	values["WHATSAPP_HOST_PROFILE_FILE"] = write("blocked-host.json", configRaw)
	values["WHATSAPP_HOST_PROFILE_SHA256"] = hashFile(values["WHATSAPP_HOST_PROFILE_FILE"])
	if bad, e := selectedWhatsAppHost(context.Background(), pool, v, func(k string) string { return values[k] }); e == nil {
		bad.close()
		t.Fatal("blocked provider activated")
	}
}
````

### FILE: `cmd/electromobility-api/whatsapp_identity_broker_test.go`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c70f2fed569129fb69519365eb633861d82ce8e6fbdc037435addef4c98d3a71"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	jose "github.com/go-jose/go-jose/v4"
)

// Exercise the host selector and domain/worker identity adapter with the real
// pinned client-credentials and RS256/JWKS verifiers, not a token provider stub.
func TestWhatsAppHostUsesOIDCServiceBroker(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	var requests, grants atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": server.URL, "authorization_endpoint": server.URL + "/authorize", "token_endpoint": server.URL + "/token", "jwks_uri": server.URL + "/jwks", "response_types_supported": []string{"code"}, "subject_types_supported": []string{"public"}, "id_token_signing_alg_values_supported": []string{"RS256"}, "grant_types_supported": []string{"client_credentials"}, "token_endpoint_auth_methods_supported": []string{"client_secret_basic"}})
		case "/jwks":
			_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "fixture-key", Use: "sig", Algorithm: "RS256"}}})
		case "/token":
			grants.Add(1)
			_ = r.ParseForm()
			id, secret, ok := r.BasicAuth()
			if r.Method != "POST" || !ok || id != "whatsapp-worker" || secret != "synthetic-secret" || r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "whatsapp:process appointment:manage" {
				t.Error("invalid grant contract")
				w.WriteHeader(400)
				return
			}
			now := time.Now()
			raw, _ := json.Marshal(map[string]any{"iss": server.URL, "sub": "fixture-worker", "aud": "enterprise-api", "iat": now.Unix(), "exp": now.Add(120 * time.Second).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"whatsapp:process", "appointment:manage"}, "scope": "whatsapp:process appointment:manage", "client_id": "whatsapp-worker"})
			signer, e := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, (&jose.SignerOptions{}).WithHeader("kid", "fixture-key"))
			if e != nil {
				t.Error(e)
				w.WriteHeader(500)
				return
			}
			signed, e := signer.Sign(raw)
			if e != nil {
				t.Error(e)
				w.WriteHeader(500)
				return
			}
			token, e := signed.CompactSerialize()
			if e != nil {
				t.Error(e)
				w.WriteHeader(500)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"token_type": "Bearer", "expires_in": 120, "access_token": token})
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	ctx := context.Background()
	verifier, err := identity.NewOIDCVerifier(ctx, server.URL, "enterprise-api")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	profileFile := filepath.Join(dir, "identity.json")
	secretFile := filepath.Join(dir, "secret")
	raw, _ := json.Marshal(map[string]any{"schema": "elite.oidc.service-token.v1", "profile_id": "wa-fixture", "revision": 1, "transport": "LOOPBACK_FIXTURE", "issuer": server.URL, "token_endpoint": server.URL + "/token", "jwks_endpoint": server.URL + "/jwks", "client_id": "whatsapp-worker", "client_authentication": "client_secret_basic", "audience": "enterprise-api", "subject": "fixture-worker", "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"whatsapp:process", "appointment:manage"}, "scopes": []string{"whatsapp:process", "appointment:manage"}, "maximum_lifetime_seconds": 120, "refresh_before_seconds": 5})
	if os.WriteFile(profileFile, raw, 0600) != nil || os.WriteFile(secretFile, []byte("synthetic-secret\n"), 0600) != nil {
		t.Fatal("fixture files")
	}
	hash := sha256.Sum256(raw)
	env := map[string]string{"WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE": profileFile, "WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256": hex.EncodeToString(hash[:]), "WHATSAPP_SERVICE_CLIENT_SECRET_FILE": secretFile}
	lookup := func(k string) string { return env[k] }
	cfg := whatsappHostConfig{TenantID: "tenant", OrganizationID: "org"}
	before := requests.Load()
	if _, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup); e == nil {
		t.Fatal("loopback profile accepted without reference opt-in")
	}
	if requests.Load() != before {
		t.Fatal("rejected fixture contacted issuer")
	}
	cfg.ReferenceFixture = true
	cfg.OrganizationID = "foreign"
	if _, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup); e == nil {
		t.Fatal("foreign organization accepted")
	}
	if requests.Load() != before {
		t.Fatal("scope rejection contacted issuer")
	}
	cfg.OrganizationID = "org"
	env["WHATSAPP_SERVICE_TOKEN_FILE"] = filepath.Join(dir, "unused")
	if _, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup); e == nil {
		t.Fatal("competing token owners accepted")
	}
	if requests.Load() != before {
		t.Fatal("mixed identity contacted issuer")
	}
	delete(env, "WHATSAPP_SERVICE_TOKEN_FILE")
	source, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup)
	if e != nil {
		t.Fatal(e)
	}
	var group sync.WaitGroup
	for i := 0; i < 12; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if token, e := source.AccessToken(ctx); e != nil || token == "" {
				t.Error("domain bearer unavailable", e)
			}
			if principal, e := source.Resolve(ctx); e != nil || principal.Subject != "fixture-worker" {
				t.Error("worker principal unavailable", e)
			}
		}()
	}
	group.Wait()
	if grants.Load() != 1 {
		t.Fatalf("concurrent host calls performed %d grants", grants.Load())
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := source.AccessToken(cancelled); e == nil {
		t.Fatal("cancelled host got cached token")
	}
	if grants.Load() != 1 {
		t.Fatal("cancelled host performed another grant")
	}
}
````

### FILE: `docs/whatsapp-connected-host.md`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bb9ffab7fb8698608063fe2a8eb46b6f8ea758152bb1fb917d701f127c0dea9b"
variables: []
secrets_allowed: false
```

````markdown
# Connected WhatsApp host

The electromobility API mounts the existing WhatsApp inbox, appointment notification and human reply modules when `WHATSAPP_ENABLED=true`. The same host starts the existing scoped job worker and constructs `app.New` once. Incoming text reaches `Runtime.Handle`, becomes a durable proposal and requires a different authenticated human reviewer before sending through the PostgreSQL outbound fence. The dispatcher is not run by this host.

Keep activation disabled while accounts are absent. Copy `docs/whatsapp-host-profile.template.json` to a private deployment configuration directory and bind its exact SHA-256 through `WHATSAPP_HOST_PROFILE_SHA256`. `WHATSAPP_HOST_PROFILE_FILE` must be absolute. The template is deliberately incomplete; the reference fixture is synthetic and is never evidence of real account access, terms or consent.

The provider profile is the existing `whatsapp_cloud/provider-profile.template.json` contract, with its own absolute path and exact hash. The fixed Python owner validates its original bytes at startup through a bounded, network-free command. A blocked provider profile prevents activation. Python executable, `whatsapp_cloud.py` and `status_reconciliation.py` also require their exact hashes. `Process` JSON property names follow the existing exported Go contract: `PythonExecutable`, `PythonSHA256`, `AdapterDirectory`, `AdapterSHA256`, `EvidenceDirectory`.

Tenant, organization, provider connection, contact consent purpose/policy and raw-payload retention approval must refer to the same materialized deployment records. The host never derives these from webhook input. Use one durable HMAC key across contact identity, outbound delivery and WhatsApp approvals. Create the evidence directory with private deployment permissions. The ordinary provider webhook registry and this receiver must not activate competing consumers for the same connection.

The domain gateway gets tenant and product/service/price-book configuration from this profile. Its legacy organization/lead defaults are both empty: `ScopedContactResolver` supplies the authorized contact organization and lead for each message. The configured model, instructions, finite token budget, retention and handoff text are passed to the existing app/runtime owners. This host does not make the current process-local economy/approval objects durable; the separate AI assurance gate retains that scope until its owner is completed.

Account secrets remain in externally protected files. The following environment values are absolute file paths, never secret contents:

```text
WHATSAPP_ACCESS_TOKEN_FILE=
WHATSAPP_APP_SECRET_FILE=
WHATSAPP_VERIFY_TOKEN_FILE=
WHATSAPP_SERVICE_TOKEN_FILE=
WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE=
WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256=
WHATSAPP_SERVICE_CLIENT_SECRET_FILE=
WHATSAPP_CONTACT_HMAC_KEY_HEX_FILE=
LLM_API_KEY_FILE=
WHATSAPP_COLLECTOR_CLIENT_CERT_FILE=
WHATSAPP_COLLECTOR_CLIENT_KEY_FILE=
```

Select exactly one service identity owner: `WHATSAPP_SERVICE_TOKEN_FILE` supplied by an external token agent, or all three `WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE`, `WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256`, `WHATSAPP_SERVICE_CLIENT_SECRET_FILE` values. Mixed or incomplete configuration fails closed. The OIDC profile uses the existing `elite.oidc.service-token.v1` contract. Its tenant, organization and required permissions must match the WhatsApp host before discovery or a grant is requested. The pinned OAuth/OIDC SDKs obtain and renew client-credentials tokens, verify RS256 and exact configured claims, and serialize concurrent renewal. The secret file is reread for each new grant. A failed renewal does not return a cached expired token.

WhatsApp token/app-secret/verify-token and externally supplied service tokens are reread when used. The actual OIDC verifier checks the selected service token on every poll and domain request; the selected tenant/organization and `whatsapp:process` plus `appointment:manage` grants are mandatory. Expired, malformed or wrong-scope tokens block work. JWT verification does not promise immediate remote revocation of an otherwise valid unexpired token; use the provider's admitted short token lifetime and disable the durable connection for local emergency denial. HMAC and LLM configuration are bound at startup and require a controlled restart to change. HMAC rotation also requires the existing identity/receipt migration procedure; changing the file alone cannot migrate stored identities.

The report collector uses TLS1.3 with exact configured CA bytes, server-name verification and a client certificate/key. The existing bounded JSON reporter writes only its six fixed operational fields. If reporting fails, the host cancels its worker and stops serving; reporting loss is not silently ignored. This validates authenticated transport and local writes, not collector retention/acknowledgement. Retention, supervision and restart policy have separate operations gates.

Public callback: `/v1/providers/whatsapp/webhook`. Human endpoints are the existing `/v1/franchise/whatsapp/replies` and appointment-notification routes, with their OIDC scope checks. Do not log callback subscription query tokens. The deployment supplies TLS/edge policy for the API itself.

Verification: three focused activation/file/identity cases, real mutual-TLS report with wrong-host/missing-key negatives, and full constructor with actual fixed Python validation and blocked-profile rejection. The constructor gate uses a lazy pool and performs no database statements or provider/model calls. Connected PostgreSQL/SDK/conversation behavior has its separate exact-source receipt. New host code is AUTHORED integration glue, with no upstream company authorship claim.
````

### FILE: `docs/whatsapp-host-profile.template.json`

```yaml
block_id: "GO-CONNECTED-WHATSAPP-HOST:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f86f34bdccd5ed089b3b709a31f6fbe4ea08069264838bb586db8ad0dbbbcd0b"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-whatsapp-host/v1",
  "tenant_id": "",
  "tenant_code": "",
  "organization_id": "",
  "connection_id": "",
  "consent_purpose": "",
  "policy_version": "",
  "retention_approval_sha256": "",
  "provider_profile_file": "",
  "provider_profile_sha256": "",
  "process": {
    "PythonExecutable": "",
    "PythonSHA256": "",
    "AdapterDirectory": "",
    "AdapterSHA256": "",
    "EvidenceDirectory": ""
  },
  "reconciler_sha256": "",
  "worker_id": "",
  "poll_seconds": 1,
  "conversation_retention_seconds": 0,
  "conversation_max_attempts": 2,
  "llm_model": "",
  "llm_base_url": "",
  "llm_token_budget": 0,
  "domain_base_url": "",
  "service_kinds": {},
  "product_variants": {},
  "price_book_id": "",
  "conversation": {
    "Instructions": "",
    "PromptCacheKey": "",
    "MaxOutputTokens": 0,
    "MaxHistory": 0,
    "MaxInputBytes": 0,
    "MaxToolBytes": 0,
    "StoreApproved": false,
    "HandoffText": ""
  },
  "approval_policy": {
    "AutoApproveMinorUnits": 0,
    "DualControlMinorUnits": 0,
    "MaxOpenPerSubject": 0
  },
  "status_collector": {
    "address": "",
    "server_name": "",
    "ca_file": "",
    "ca_sha256": ""
  },
  "reference_fixture": false
}
````

## 6. Configuration surface

WHATSAPP_ENABLED defaults false. Exact host/provider/process/CA profiles and absolute external secret-file references. Choose token file or service-broker profile/hash/client-secret file. LOOPBACK_FIXTURE requires explicit reference_fixture=true in the hash-bound host profile; default is false.

## 7. Dependency bill

No new versions. Go OIDC/OAuth, PG, existing app/runtime/WhatsApp and crypto/TLS owners are selected separately. Python adapter and executable retain exact hashes.

## 8. Apply order

Materialize full franchise and broker, provide reference fixtures or later private account configuration, run migrations and activate this optional host. No implicit account acquisition.

## 9. Verification

Three changed-identity/assembly/broker tests passed on the composed source, including twelve concurrent domain/worker requests with exactly one real fixture grant. Earlier unchanged host file/activation and mutual-TLS tests are reused with hashes. Browser/PG/SDK tests are separately bound; global source/security/operations gates remain separate.

## 10. Reconstruction evidence

See reconstruction_evidence/COMMUNICATIONS_RUNTIME_V402.md and exact composition receipt. Current in-process LLM economy/approval persistence remains a separate T2807 task; this host does not close it by wiring.


V402 composed delta: Scheduled WhatsApp exact source/approval/job/host glue and template receipt recovery; SCHEDULED_COMMUNICATIONS_RELEASE_V402.md/json.

V402 composed delta: Optional campaign source/typed job scope and fixed policy host hook; appointment compatibility proven; CAMPAIGN_CONNECTED_RELEASE_V402.md/json.

V402 composed delta: T2807 connected reference315: real domain quotation and contact-bound order status, explicit single-vehicle limit, revalidated contact, structured history roles, per-case required eval gates and canonical host mounting. Local fixtures only. AI_CONNECTED_REFERENCE_RELEASE_V402.md/json. No new upstream dependency or live model quality claim.
