# Go Mercado Libre Question Outbound

## 1. Metadata

```yaml
pack_id: "GO-MERCADOLIBRE-QUESTION-OUTBOUND"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa respuesta aprobada a preguntas de Mercado Libre sobre el fence outbound existente, con un único POST automático, confirmación GET v4 exacta y reconciliación sin reenvío."
stacks: ["Go 1.26.7", "Mercado Libre Questions API v4", "GO-PG-OUTBOUND-DELIVERY-FENCE 0.2.x"]
compatible_with: ["GO-CHANNELS-CORE 0.4.x", "GO-PG-OUTBOUND-DELIVERY-FENCE 0.2.x", "GO-MERCADOLIBRE-MARKETPLACE-ADAPTER 0.2.x"]
incompatible_with: ["respuesta automática sin aprobación", "retry ciego", "respuesta sin read-after-write", "token persistido", "respuesta mayor a 2.000 caracteres"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers", "https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html", "https://www.postgresql.org/docs/18/mvcc.html"]
verified_at: "2026-09-05"
```

Los tres archivos son `AUTHORED`. El contrato HTTP, payload, UTF-8, límite y recomendación semi-automática proceden de la documentación oficial vigente de Mercado Libre; el fence de intento/ambigüedad conserva las autoridades AWS/PostgreSQL del owner existente. Nada se atribuye falsamente a esas empresas como código copiado.

## 2. Applicability

Use cuando un proyecto deba responder preguntas pre-venta de Mercado Libre desde el runtime conversacional. Requiere que `GO-PG-OUTBOUND-DELIVERY-FENCE` ya materialice `outbounddelivery.Channel` y su store PostgreSQL. No use para crear preguntas, borrar, ocultar, moderar, publicar items ni mensajería postventa.

## 3. Architecture contract

- **Un owner**: añade un sender/reconciler al package `internal/outbounddelivery`; no crea otra cola, tabla, ledger ni retry loop.
- **Aprobación demostrable**: `tenant + DeliveryKey + question_id + SHA-256(text)` deben coincidir con una aprobación no vencida antes del POST.
- **Contrato oficial**: `POST /answers`, JSON con `question_id` entero y `text`, UTF-8, máximo 2.000 runas; Bearer token sólo en memoria.
- **Confirmación fuerte**: después de un 2xx, `GET /questions/{id}?api_version=4` debe devolver question/seller exactos, `ANSWERED`, answer `ACTIVE` y texto idéntico.
- **Ambigüedad inmovilizada**: timeout, respuesta no confirmada o receipt inválido producen `unknown`; el fence existente impide un segundo POST automático.
- **Reconciliación**: sólo GET; exact match cierra `accepted`, divergencia/terminal cierra `failed_terminal`, `UNANSWERED` permanece `unknown`.
- **Privacidad**: no persiste token, texto ni respuesta provider; el ledger conserva hashes/HMAC y evidencia de aprobación.
- **Semántica honesta**: no afirma exactly-once del proveedor ni cuenta productiva; afirma un máximo de una invocación automática por DeliveryKey.

## 4. Exact file manifest

```text
CREATE internal/outbounddelivery/mercadolibre_question.go
CREATE internal/outbounddelivery/mercadolibre_question_test.go
CREATE internal/outbounddelivery/MERCADOLIBRE_QUESTION_OUTBOUND.md
```

## 5. Materialization blocks

### FILE: `internal/outbounddelivery/mercadolibre_question.go`
```yaml
block_id: "GO-MERCADOLIBRE-QUESTION-OUTBOUND:internal/outbounddelivery/mercadolibre_question.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter governed by official Mercado Libre Questions & Answers contract"
license: "LicenseRef-Workspace-Owner"
sha256: "5c30e9937bad08db395d98243198740e17f186a1da63c64ad49065e35587c6f0"
variables: []
secrets_allowed: false
```
````go
package outbounddelivery

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"elite.local/enterprise/internal/channels"
)

const MercadoLibreAPIBaseURL = "https://api.mercadolibre.com"

var (
	ErrApprovalRequired    = errors.New("outbounddelivery: approved Mercado Libre answer is required")
	ErrProviderUnconfirmed = errors.New("outbounddelivery: Mercado Libre answer is not confirmed")
	mlDigitsRE             = regexp.MustCompile(`^[0-9]{3,20}$`)
)

type MercadoLibreQuestionConfig struct {
	ChannelCode                   string
	SellerID                      string
	OfficialDocumentationReviewed bool
	SellerAuthorizationProven     bool
	WriteScopeProven              bool
	ReconciliationApproved        bool
}

func (c MercadoLibreQuestionConfig) Validate() error {
	if c.ChannelCode != "mercadolibre_questions" || !validMercadoLibreID(c.SellerID) || !c.OfficialDocumentationReviewed || !c.SellerAuthorizationProven || !c.WriteScopeProven || !c.ReconciliationApproved {
		return ErrInvalid
	}
	return nil
}

type MercadoLibreAnswerApproval struct {
	TenantID       string
	DeliveryKey    string
	QuestionID     string
	AnswerSHA256   string
	ApprovedBy     string
	EvidenceSHA256 string
	ApprovedAt     time.Time
	ExpiresAt      time.Time
}

type MercadoLibreApprovalResolver interface {
	ResolveMercadoLibreAnswerApproval(context.Context, string, string) (MercadoLibreAnswerApproval, error)
}

type MercadoLibreTokenSource interface {
	MercadoLibreAccessToken(context.Context) (string, error)
}

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type MercadoLibreQuestionSender struct {
	Config    MercadoLibreQuestionConfig
	Approvals MercadoLibreApprovalResolver
	Tokens    MercadoLibreTokenSource
	HTTP      HTTPDoer
	Now       func() time.Time
}

type mercadoLibreQuestion struct {
	ID       int64  `json:"id"`
	SellerID int64  `json:"seller_id"`
	Status   string `json:"status"`
	Answer   *struct {
		Text        string `json:"text"`
		Status      string `json:"status"`
		DateCreated string `json:"date_created"`
	} `json:"answer"`
}

type MercadoLibreReconciliation struct {
	Decision       string
	Receipt        Receipt
	EvidenceSHA256 string
	FailureCode    string
}

type MercadoLibreReconciliationStore interface {
	ReconcileAccepted(context.Context, channels.Message, string, Receipt) error
	ReconcileFailed(context.Context, channels.Message, string, string, string) error
}

func (s *MercadoLibreQuestionSender) SendWithReceipt(ctx context.Context, message channels.Message) (Receipt, error) {
	now, approval, token, err := s.authorize(ctx, message, true)
	if err != nil {
		return Receipt{}, err
	}
	body, err := json.Marshal(struct {
		QuestionID int64  `json:"question_id"`
		Text       string `json:"text"`
	}{QuestionID: parseQuestionID(message.ExternalID), Text: message.Text})
	if err != nil {
		return Receipt{}, err
	}
	postStatus, postBody, err := s.call(ctx, http.MethodPost, "/answers", body, token)
	if err != nil {
		return Receipt{}, err
	}
	if postStatus < 200 || postStatus > 299 {
		if !documentedAnswerRejection(postStatus, postBody) {
			return Receipt{}, ErrProviderUnconfirmed
		}
		evidence := evidenceHash(postStatus, postBody, 0, nil, approval.EvidenceSHA256)
		return Receipt{}, NewTerminalFailure("MERCADOLIBRE_ANSWER_REJECTED", evidence, fmt.Errorf("%w: answer POST status %d", ErrProviderUnconfirmed, postStatus))
	}
	question, getStatus, getBody, err := s.fetchQuestion(ctx, message.ExternalID, token)
	if err != nil {
		return Receipt{}, err
	}
	if getStatus != http.StatusOK || !questionConfirms(question, s.Config.SellerID, message.ExternalID, message.Text) {
		return Receipt{}, ErrProviderUnconfirmed
	}
	evidence := evidenceHash(postStatus, postBody, getStatus, getBody, approval.EvidenceSHA256)
	return Receipt{ProviderMessageID: "mercadolibre-question:" + message.ExternalID, EvidenceSHA256: evidence, AcceptedAt: now}, nil
}

func (s *MercadoLibreQuestionSender) Reconcile(ctx context.Context, message channels.Message) (MercadoLibreReconciliation, error) {
	now, approval, token, err := s.authorize(ctx, message, false)
	if err != nil {
		return MercadoLibreReconciliation{}, err
	}
	question, status, body, err := s.fetchQuestion(ctx, message.ExternalID, token)
	if err != nil {
		return MercadoLibreReconciliation{}, err
	}
	evidence := evidenceHash(0, nil, status, body, approval.EvidenceSHA256)
	if status != http.StatusOK {
		return MercadoLibreReconciliation{Decision: "STILL_UNKNOWN", EvidenceSHA256: evidence, FailureCode: "QUESTION_LOOKUP_UNAVAILABLE"}, nil
	}
	if questionConfirms(question, s.Config.SellerID, message.ExternalID, message.Text) {
		return MercadoLibreReconciliation{Decision: "ACCEPTED", Receipt: Receipt{ProviderMessageID: "mercadolibre-question:" + message.ExternalID, EvidenceSHA256: evidence, AcceptedAt: now}, EvidenceSHA256: evidence}, nil
	}
	if question.ID != parseQuestionID(message.ExternalID) || fmt.Sprint(question.SellerID) != s.Config.SellerID {
		return MercadoLibreReconciliation{Decision: "STILL_UNKNOWN", EvidenceSHA256: evidence, FailureCode: "QUESTION_IDENTITY_UNCONFIRMED"}, nil
	}
	if question.Status == "ANSWERED" && question.Answer != nil && question.Answer.Text != "" && (question.Answer.Status == "ACTIVE" || question.Answer.Status == "BANNED") {
		return MercadoLibreReconciliation{Decision: "FAILED_TERMINAL", EvidenceSHA256: evidence, FailureCode: "QUESTION_ANSWER_DIVERGED_OR_TERMINAL"}, nil
	}
	return MercadoLibreReconciliation{Decision: "STILL_UNKNOWN", EvidenceSHA256: evidence, FailureCode: "QUESTION_STATE_UNCONFIRMED"}, nil
}

func ReconcileMercadoLibreQuestion(ctx context.Context, sender *MercadoLibreQuestionSender, store MercadoLibreReconciliationStore, message channels.Message) error {
	if sender == nil || store == nil {
		return ErrInvalid
	}
	requestHash, err := MessageSHA256(message)
	if err != nil {
		return err
	}
	result, err := sender.Reconcile(ctx, message)
	if err != nil {
		return err
	}
	switch result.Decision {
	case "ACCEPTED":
		return store.ReconcileAccepted(ctx, message, requestHash, result.Receipt)
	case "FAILED_TERMINAL":
		return store.ReconcileFailed(ctx, message, requestHash, result.EvidenceSHA256, result.FailureCode)
	case "STILL_UNKNOWN":
		return ErrUnknown
	default:
		return ErrInvalid
	}
}

func (s *MercadoLibreQuestionSender) authorize(ctx context.Context, message channels.Message, requireCurrent bool) (time.Time, MercadoLibreAnswerApproval, string, error) {
	if s == nil || s.Config.Validate() != nil || s.Approvals == nil || s.Tokens == nil || s.HTTP == nil || s.Now == nil || message.ChannelCode != s.Config.ChannelCode || message.Direction != channels.DirectionOut || message.Validate() != nil || !validMercadoLibreID(message.ExternalID) || !utf8.ValidString(message.Text) {
		return time.Time{}, MercadoLibreAnswerApproval{}, "", ErrInvalid
	}
	trimmed := strings.TrimSpace(message.Text)
	if trimmed != message.Text || utf8.RuneCountInString(message.Text) < 1 || utf8.RuneCountInString(message.Text) > 2000 {
		return time.Time{}, MercadoLibreAnswerApproval{}, "", ErrInvalid
	}
	now := s.Now().UTC()
	approval, err := s.Approvals.ResolveMercadoLibreAnswerApproval(ctx, message.TenantID, message.DeliveryKey)
	if err != nil || approval.TenantID != message.TenantID || approval.DeliveryKey != message.DeliveryKey || approval.QuestionID != message.ExternalID || approval.AnswerSHA256 != textHash(message.Text) || strings.TrimSpace(approval.ApprovedBy) == "" || !hex64RE.MatchString(approval.EvidenceSHA256) || approval.ApprovedAt.IsZero() || approval.ExpiresAt.IsZero() || !approval.ExpiresAt.After(approval.ApprovedAt) || (requireCurrent && (now.Before(approval.ApprovedAt.UTC()) || !now.Before(approval.ExpiresAt.UTC()))) {
		return time.Time{}, MercadoLibreAnswerApproval{}, "", ErrApprovalRequired
	}
	token, err := s.Tokens.MercadoLibreAccessToken(ctx)
	if err != nil || len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return time.Time{}, MercadoLibreAnswerApproval{}, "", ErrInvalid
	}
	return now, approval, token, nil
}

func (s *MercadoLibreQuestionSender) fetchQuestion(ctx context.Context, questionID, token string) (mercadoLibreQuestion, int, []byte, error) {
	status, body, err := s.call(ctx, http.MethodGet, "/questions/"+questionID+"?api_version=4", nil, token)
	if err != nil {
		return mercadoLibreQuestion{}, 0, nil, err
	}
	var question mercadoLibreQuestion
	if status == http.StatusOK {
		if err := json.Unmarshal(body, &question); err != nil {
			return mercadoLibreQuestion{}, status, body, ErrProviderUnconfirmed
		}
	}
	question.Status = strings.ToUpper(question.Status)
	if question.Answer != nil {
		question.Answer.Status = strings.ToUpper(question.Answer.Status)
	}
	return question, status, body, nil
}

func (s *MercadoLibreQuestionSender) call(ctx context.Context, method, path string, body []byte, token string) (int, []byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, MercadoLibreAPIBaseURL+path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	response, err := s.HTTP.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("Mercado Libre request failed")
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	if err != nil {
		return 0, nil, err
	}
	if len(responseBody) > 1<<20 {
		return 0, nil, ErrProviderUnconfirmed
	}
	return response.StatusCode, responseBody, nil
}

func questionConfirms(question mercadoLibreQuestion, sellerID, questionID, answer string) bool {
	return question.ID == parseQuestionID(questionID) && fmt.Sprint(question.SellerID) == sellerID && question.Status == "ANSWERED" && question.Answer != nil && question.Answer.Status == "ACTIVE" && question.Answer.Text == answer
}

func parseQuestionID(value string) int64 {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func validMercadoLibreID(value string) bool {
	if !mlDigitsRE.MatchString(value) {
		return false
	}
	id, err := strconv.ParseInt(value, 10, 64)
	return err == nil && id > 0 && strconv.FormatInt(id, 10) == value
}

func documentedAnswerRejection(status int, body []byte) bool {
	if status != http.StatusBadRequest {
		return false
	}
	var problem struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &problem) != nil {
		return false
	}
	return problem.Error == "invalid_question" || problem.Error == "invalid_post_body"
}

func textHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func evidenceHash(postStatus int, postBody []byte, getStatus int, getBody []byte, approvalEvidence string) string {
	canonical, _ := json.Marshal(struct {
		PostStatus       int    `json:"post_status"`
		PostSHA256       string `json:"post_sha256"`
		GetStatus        int    `json:"get_status"`
		GetSHA256        string `json:"get_sha256"`
		ApprovalEvidence string `json:"approval_evidence_sha256"`
	}{postStatus, textHash(string(postBody)), getStatus, textHash(string(getBody)), approvalEvidence})
	return textHash(string(canonical))
}
````

### FILE: `internal/outbounddelivery/mercadolibre_question_test.go`
```yaml
block_id: "GO-MERCADOLIBRE-QUESTION-OUTBOUND:internal/outbounddelivery/mercadolibre_question_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local executable contract and failure regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "3243d098d830c574f11f1495c6c107939ebe3b2ad2dcea0414778a845e05bb28"
variables: []
secrets_allowed: false
```
````go
package outbounddelivery

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
)

type mlApprovalResolver struct{ approval MercadoLibreAnswerApproval }

func (r mlApprovalResolver) ResolveMercadoLibreAnswerApproval(context.Context, string, string) (MercadoLibreAnswerApproval, error) {
	return r.approval, nil
}

type mlTokenSource struct{ token string }

func (s mlTokenSource) MercadoLibreAccessToken(context.Context) (string, error) { return s.token, nil }

type mlHTTP struct {
	requests  []*http.Request
	bodies    []string
	responses []*http.Response
	err       error
}

func (h *mlHTTP) Do(request *http.Request) (*http.Response, error) {
	h.requests = append(h.requests, request)
	if request.Body != nil {
		raw, _ := io.ReadAll(request.Body)
		h.bodies = append(h.bodies, string(raw))
	}
	if h.err != nil {
		return nil, h.err
	}
	response := h.responses[0]
	h.responses = h.responses[1:]
	return response, nil
}

type mlReconcileStore struct{ accepted, failed int }

func (s *mlReconcileStore) ReconcileAccepted(context.Context, channels.Message, string, Receipt) error {
	s.accepted++
	return nil
}
func (s *mlReconcileStore) ReconcileFailed(context.Context, channels.Message, string, string, string) error {
	s.failed++
	return nil
}

func mlMessage() channels.Message {
	return channels.Message{ChannelCode: "mercadolibre_questions", TenantID: "tenant", ExternalID: "3957150025", DeliveryKey: strings.Repeat("d", 64), Direction: channels.DirectionOut, Text: "Sí, tenemos disponibilidad inmediata."}
}

func mlSender(httpClient *mlHTTP) *MercadoLibreQuestionSender {
	message := mlMessage()
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	return &MercadoLibreQuestionSender{
		Config:    MercadoLibreQuestionConfig{ChannelCode: "mercadolibre_questions", SellerID: "179571326", OfficialDocumentationReviewed: true, SellerAuthorizationProven: true, WriteScopeProven: true, ReconciliationApproved: true},
		Approvals: mlApprovalResolver{MercadoLibreAnswerApproval{TenantID: message.TenantID, DeliveryKey: message.DeliveryKey, QuestionID: message.ExternalID, AnswerSHA256: textHash(message.Text), ApprovedBy: "operator-7", EvidenceSHA256: strings.Repeat("a", 64), ApprovedAt: now.Add(-time.Minute), ExpiresAt: now.Add(time.Hour)}},
		Tokens:    mlTokenSource{strings.Repeat("t", 32)}, HTTP: httpClient, Now: func() time.Time { return now },
	}
}

func mlResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}
}

func confirmedQuestion(answer string) string {
	return `{"id":3957150025,"seller_id":179571326,"status":"ANSWERED","answer":{"text":"` + answer + `","status":"ACTIVE","date_created":"2026-09-05T12:00:00Z"}}`
}

func TestMercadoLibreAnswerUsesExactOfficialContractAndReadAfterWrite(t *testing.T) {
	message := mlMessage()
	httpClient := &mlHTTP{responses: []*http.Response{mlResponse(200, `{"status":"ok"}`), mlResponse(200, confirmedQuestion(message.Text))}}
	receipt, err := mlSender(httpClient).SendWithReceipt(context.Background(), message)
	if err != nil {
		t.Fatal(err)
	}
	if len(httpClient.requests) != 2 || httpClient.requests[0].Method != http.MethodPost || httpClient.requests[0].URL.String() != MercadoLibreAPIBaseURL+"/answers" || httpClient.requests[1].Method != http.MethodGet || httpClient.requests[1].URL.String() != MercadoLibreAPIBaseURL+"/questions/3957150025?api_version=4" {
		t.Fatalf("unexpected requests: %+v", httpClient.requests)
	}
	if httpClient.bodies[0] != `{"question_id":3957150025,"text":"Sí, tenemos disponibilidad inmediata."}` || httpClient.requests[0].Header.Get("Content-Type") != "application/json; charset=utf-8" || httpClient.requests[0].Header.Get("Authorization") != "Bearer "+strings.Repeat("t", 32) {
		t.Fatalf("contract mismatch body=%q headers=%v", httpClient.bodies[0], httpClient.requests[0].Header)
	}
	if receipt.ProviderMessageID != "mercadolibre-question:3957150025" || len(receipt.EvidenceSHA256) != 64 || receipt.AcceptedAt.IsZero() {
		t.Fatalf("invalid receipt: %+v", receipt)
	}
}

func TestMercadoLibreAnswerFailsClosedBeforeProvider(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*MercadoLibreQuestionSender, *channels.Message)
	}{
		{"missing approval", func(sender *MercadoLibreQuestionSender, _ *channels.Message) { sender.Approvals = mlApprovalResolver{} }},
		{"expired approval", func(sender *MercadoLibreQuestionSender, _ *channels.Message) {
			a := sender.Approvals.(mlApprovalResolver).approval
			a.ExpiresAt = a.ApprovedAt
			sender.Approvals = mlApprovalResolver{a}
		}},
		{"answer changed", func(_ *MercadoLibreQuestionSender, message *channels.Message) { message.Text = "texto distinto" }},
		{"more than 2000 characters", func(_ *MercadoLibreQuestionSender, message *channels.Message) {
			message.Text = strings.Repeat("á", 2001)
		}},
		{"write scope absent", func(sender *MercadoLibreQuestionSender, _ *channels.Message) { sender.Config.WriteScopeProven = false }},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			httpClient := &mlHTTP{}
			sender, message := mlSender(httpClient), mlMessage()
			test.mutate(sender, &message)
			if _, err := sender.SendWithReceipt(context.Background(), message); err == nil {
				t.Fatal("invalid send accepted")
			}
			if len(httpClient.requests) != 0 {
				t.Fatal("provider called before validation")
			}
		})
	}
}

func TestMercadoLibreAmbiguousCallIsNeverRetriedByFence(t *testing.T) {
	httpClient := &mlHTTP{err: errors.New("timeout after write")}
	sender := mlSender(httpClient)
	store := &memoryStore{}
	channel := &Channel{CodeValue: "mercadolibre_questions", Receiver: testReceiver{code: "mercadolibre_questions"}, Sender: sender, Store: store}
	message := mlMessage()
	if err := channel.Send(context.Background(), message); !errors.Is(err, ErrUnknown) {
		t.Fatalf("first=%v", err)
	}
	if err := channel.Send(context.Background(), message); !errors.Is(err, ErrUnknown) {
		t.Fatalf("replay=%v", err)
	}
	if len(httpClient.requests) != 1 || store.state != "unknown" {
		t.Fatalf("requests=%d state=%s", len(httpClient.requests), store.state)
	}
}

func TestMercadoLibreExplicitProviderRejectionIsTerminal(t *testing.T) {
	httpClient := &mlHTTP{responses: []*http.Response{mlResponse(400, `{"error":"invalid_question"}`)}}
	sender := mlSender(httpClient)
	store := &memoryStore{}
	channel := &Channel{CodeValue: "mercadolibre_questions", Receiver: testReceiver{code: "mercadolibre_questions"}, Sender: sender, Store: store}
	message := mlMessage()
	if err := channel.Send(context.Background(), message); !errors.Is(err, ErrTerminal) {
		t.Fatalf("first=%v", err)
	}
	if err := channel.Send(context.Background(), message); !errors.Is(err, ErrTerminal) {
		t.Fatalf("replay=%v", err)
	}
	if len(httpClient.requests) != 1 || store.state != "failed_terminal" {
		t.Fatalf("requests=%d state=%s", len(httpClient.requests), store.state)
	}
}

func TestMercadoLibreReadAfterWriteMustConfirmExactAnswerAndSeller(t *testing.T) {
	message := mlMessage()
	for _, body := range []string{
		`{"id":3957150025,"seller_id":999999999,"status":"ANSWERED","answer":{"text":"Sí, tenemos disponibilidad inmediata.","status":"ACTIVE"}}`,
		`{"id":3957150025,"seller_id":179571326,"status":"ANSWERED","answer":{"text":"otra respuesta","status":"ACTIVE"}}`,
		`{"id":3957150025,"seller_id":179571326,"status":"ANSWERED","answer":{"text":"","status":"BANNED"}}`,
	} {
		httpClient := &mlHTTP{responses: []*http.Response{mlResponse(200, `{}`), mlResponse(200, body)}}
		if _, err := mlSender(httpClient).SendWithReceipt(context.Background(), message); !errors.Is(err, ErrProviderUnconfirmed) {
			t.Fatalf("body=%s err=%v", body, err)
		}
	}
}

func TestMercadoLibreReconciliationClosesOnlyProvenStates(t *testing.T) {
	message := mlMessage()
	t.Run("accepted", func(t *testing.T) {
		httpClient := &mlHTTP{responses: []*http.Response{mlResponse(200, confirmedQuestion(message.Text))}}
		store := &mlReconcileStore{}
		if err := ReconcileMercadoLibreQuestion(context.Background(), mlSender(httpClient), store, message); err != nil || store.accepted != 1 || store.failed != 0 {
			t.Fatalf("err=%v store=%+v", err, store)
		}
	})
	t.Run("still unknown", func(t *testing.T) {
		httpClient := &mlHTTP{responses: []*http.Response{mlResponse(200, `{"id":3957150025,"seller_id":179571326,"status":"UNANSWERED","answer":null}`)}}
		store := &mlReconcileStore{}
		if err := ReconcileMercadoLibreQuestion(context.Background(), mlSender(httpClient), store, message); !errors.Is(err, ErrUnknown) || store.accepted != 0 || store.failed != 0 {
			t.Fatalf("err=%v store=%+v", err, store)
		}
	})
	t.Run("terminal divergence", func(t *testing.T) {
		httpClient := &mlHTTP{responses: []*http.Response{mlResponse(200, `{"id":3957150025,"seller_id":179571326,"status":"ANSWERED","answer":{"text":"otra","status":"ACTIVE"}}`)}}
		store := &mlReconcileStore{}
		if err := ReconcileMercadoLibreQuestion(context.Background(), mlSender(httpClient), store, message); err != nil || store.accepted != 0 || store.failed != 1 {
			t.Fatalf("err=%v store=%+v", err, store)
		}
	})
}

func TestMercadoLibreNumericIDsRejectLossBeforeHTTP(t *testing.T) {
	for _, id := range []string{"9223372036854775808", "18446744073709551616", "99999999999999999999", "000", "03957150025"} {
		t.Run(id, func(t *testing.T) {
			h := &mlHTTP{err: errors.New("must never reach provider")}
			s := mlSender(h)
			m := mlMessage()
			m.ExternalID = id
			a := s.Approvals.(mlApprovalResolver).approval
			a.QuestionID = id
			s.Approvals = mlApprovalResolver{a}
			_, err := s.SendWithReceipt(context.Background(), m)
			if !errors.Is(err, ErrInvalid) || len(h.requests) != 0 {
				t.Fatalf("invalid ID reached HTTP: error=%v calls=%d", err, len(h.requests))
			}
		})
	}
}

func TestMercadoLibreUncertainHTTPPreservesUnknown(t *testing.T) {
	for _, status := range []int{301, 400, 408, 409, 429, 500, 502, 503, 504} {
		h := &mlHTTP{responses: []*http.Response{mlResponse(status, "{}")}}
		store := &memoryStore{}
		ch := &Channel{CodeValue: "mercadolibre_questions", Receiver: testReceiver{code: "mercadolibre_questions"}, Sender: mlSender(h), Store: store}
		for i := 0; i < 2; i++ {
			if err := ch.Send(context.Background(), mlMessage()); !errors.Is(err, ErrUnknown) {
				t.Errorf("HTTP %d dispatch %d: expected unknown, got %v", status, i, err)
			}
		}
		if len(h.requests) != 1 || store.state != "unknown" {
			t.Errorf("HTTP %d calls=%d state=%s", status, len(h.requests), store.state)
		}
	}
}

func TestMercadoLibreUntrustedLookupDoesNotCloseStore(t *testing.T) {
	for _, body := range []string{
		"{}",
		`{"id":3957150025,"seller_id":999999999,"status":"ANSWERED","answer":{"text":"otra","status":"ACTIVE"}}`,
		`{"id":3957150026,"seller_id":179571326,"status":"ANSWERED","answer":{"text":"otra","status":"ACTIVE"}}`,
		`{"id":3957150025,"seller_id":179571326,"status":"NEW_UNKNOWN_STATE"}`,
		`{"id":3957150025,"seller_id":179571326,"status":"ANSWERED","answer":null}`,
	} {
		h := &mlHTTP{responses: []*http.Response{mlResponse(200, body)}}
		store := &mlReconcileStore{}
		err := ReconcileMercadoLibreQuestion(context.Background(), mlSender(h), store, mlMessage())
		if !errors.Is(err, ErrUnknown) || store.accepted != 0 || store.failed != 0 {
			t.Errorf("untrusted lookup closed store: error=%v accepted=%d failed=%d", err, store.accepted, store.failed)
		}
	}
}

func TestMercadoLibreIDBoundaryIsLossless(t *testing.T) {
	for _, id := range []string{"100", "3957150025", "9223372036854775807"} {
		if !validMercadoLibreID(id) || parseQuestionID(id) <= 0 {
			t.Fatalf("valid boundary rejected: %s", id)
		}
	}
	for _, id := range []string{"000", "0123", "9223372036854775808", "-100", "+100", "1e3"} {
		s := mlSender(&mlHTTP{})
		s.Config.SellerID = id
		if s.Config.Validate() == nil {
			t.Fatalf("invalid seller accepted: %s", id)
		}
	}
}
````

### FILE: `internal/outbounddelivery/MERCADOLIBRE_QUESTION_OUTBOUND.md`
```yaml
block_id: "GO-MERCADOLIBRE-QUESTION-OUTBOUND:internal/outbounddelivery/MERCADOLIBRE_QUESTION_OUTBOUND.md:v1"
operation: CREATE
provenance: AUTHORED
source: "operator boundary and non-claims"
license: "LicenseRef-Workspace-Owner"
sha256: "b9148eac5ffe182703c3393d73587c30ee45bc9c8377196814e90df3bd24b9ea"
variables: []
secrets_allowed: false
```
````markdown
# Mercado Libre question outbound

`MercadoLibreQuestionSender` is an authored adapter governed by Mercado Libre's official Questions & Answers contract (last update 2026-01-15). It uses `POST /answers` with integer `question_id` and UTF-8 `text` limited to 2,000 characters, then confirms the exact seller/question/text/status through `GET /questions/{id}?api_version=4`.

The sender plugs into the existing `outbounddelivery.Channel`; it does not create another ledger. A durable, answer-hash-bound operator approval is required before the provider call. Network ambiguity becomes `unknown`, never an automatic retry. `ReconcileMercadoLibreQuestion` performs a GET-only reconciliation and transitions the existing PostgreSQL store only when the result is proved accepted or terminal.

No access token, answer body or provider response is written by this package. Target-account authorization, terms, quota, moderation behavior, privacy, callback origin, load, security and operational acceptance remain project gates.

Version 0.1.1 validates canonical positive int64 identifiers with Go's `strconv.ParseInt` and a lossless decimal round-trip. No overflowing or aliased identifier reaches HTTP. Only HTTP 400 carrying the documented `invalid_question` or `invalid_post_body` error is a proven rejection. Other non-2xx responses remain uncertain, with GET-only reconciliation and no automatic second POST. Reconciliation requires the exact seller and question; incomplete responses, unknown states or another identity cannot close the ledger.

Method authorities: https://pkg.go.dev/strconv#ParseInt defines numeric range errors; https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/ explains uncertain side effects and reconciliation. These corrections are AUTHORED and are not copied provider implementation. Use an HTTP transport with bounded timeouts and no automatic redirects or POST retries; the sender's single invocation does not itself prove transport-level exactly-once delivery.
````

## 6. Configuration surface

| Campo | Default seguro | Gate | Secreto | Efecto |
|---|---|---|---|---|
| `ChannelCode` | `mercadolibre_questions` exacto | obligatorio | no | routing único |
| `SellerID` | vacío | seller autorizado/probado | identificador sensible | binding de respuesta |
| proofs docs/seller/write/reconciliation | `false` | todos requeridos | no | fail closed |
| `MercadoLibreApprovalResolver` | ausente | aprobación durable exacta | evidencia sensible | autoriza una respuesta |
| `MercadoLibreTokenSource` | ausente | secret manager target | sí, sólo valor runtime | Bearer auth |
| `HTTPDoer` | ausente | cliente con timeout/TLS/observabilidad | no | provider boundary |
| `Now` | ausente | reloj UTC inyectado | no | vigencia/receipt |

## 7. Dependency bill

| Package/API | Pin/contrato | Uso | Licencia/términos | Fuente oficial |
|---|---|---|---|---|
| Go stdlib | Go 1.26.7 | HTTP, JSON, UTF-8, SHA-256, tests | BSD-3-Clause | https://go.dev/dl/ |
| Mercado Libre Questions API | docs actualizadas 2026-01-15 | POST answer + GET v4 | términos de plataforma; no redistribuidos | https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers |
| owner outbound PostgreSQL | pack 0.2.x | claim, receipt, unknown, reconciliation | workspace + MIT pgx indirecto | pack local gobernado por AWS/PostgreSQL |

## 8. Apply order

1. Completar readiness: aplicación/seller, write scope, token store/rotation, términos, privacidad, moderación, quota/costo, callback origin y reconciliación.
2. Materializar después de `GO-CHANNELS-CORE` y `GO-PG-OUTBOUND-DELIVERY-FENCE`; registrar el sender como `mercadolibre_questions` dentro de `outbounddelivery.Channel`.
3. Persistir aprobaciones autorizadas y su evidencia; no derivarlas del modelo ni de texto similar.
4. Ejecutar tests/vet/build y el journey inbound question → candidate → runtime → aprobación → POST → GET → receipt.
5. Ante `unknown`, ejecutar únicamente `ReconcileMercadoLibreQuestion`; nunca repetir el mismo POST/DeliveryKey.
6. Rollback: cortar routing, revocar/rotar token, preservar ledger/approval/evidence y reauditar el contrato oficial.

## 9. Verification

V259 / version 0.1.1: red-before/green-after tests cover overflowing/noncanonical IDs, uncertain HTTP statuses and mismatched/incomplete reconciliation identities. Go `strconv.ParseInt` supplies the checked conversion. Only documented HTTP 400 validation errors close terminal; the remaining responses require GET reconciliation. Four new tests bring the composed outbound package to 14 top-level tests. Full composed Go test/vet/build, 49 PostgreSQL migrations and PostgreSQL/app/outbound suites pass. See `reconstruction_evidence/MERCADOLIBRE_OUTBOUND_HARDENING_V259.md`; provider calls are simulated and no live delivery is claimed.

- 3/3 archivos y hashes byte-idénticos.
- Go 1.26.7: tests focales y full graph, `go vet` y build.
- Contrato exacto: POST `/answers`, integer `question_id`, UTF-8, GET v4, Bearer y texto exactos.
- Falta de approval/proofs/scope/token, approval vencida/divergente, >2.000 caracteres, seller/text/status divergentes y provider ambiguity fallan cerrados.
- Replay tras timeout demuestra una sola llamada provider; reconciliación accepted/unknown/terminal no reenvía.
- Cuenta/token/callback/quotas/moderación/privacidad/carga/seguridad/aceptación reales permanecen gates del proyecto.

## 10. Reconstruction evidence

`reconstruction_evidence/MERCADOLIBRE_QUESTION_OUTBOUND_2026-09-05_V258.md` conserva fuentes, toolchain, hashes, pruebas, composición y no-claims.
