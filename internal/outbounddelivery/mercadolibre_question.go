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
