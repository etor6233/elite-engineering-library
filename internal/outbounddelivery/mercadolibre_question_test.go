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
