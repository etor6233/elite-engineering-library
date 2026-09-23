package leadstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type mercadoLibreQuestionParty struct {
	ID        json.Number `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     struct {
		AreaCode string `json:"area_code"`
		Number   string `json:"number"`
	} `json:"phone"`
}

type mercadoLibreQuestion struct {
	ID          json.Number                `json:"id"`
	SellerID    json.Number                `json:"seller_id"`
	BuyerID     json.Number                `json:"buyer_id"`
	ItemID      string                     `json:"item_id"`
	Status      string                     `json:"status"`
	Text        string                     `json:"text"`
	DateCreated string                     `json:"date_created"`
	From        *mercadoLibreQuestionParty `json:"from"`
}

type MercadoLibreQuestionDecoded struct {
	Raw               RawEvent
	Candidate         *LeadCandidate
	NormalizationCode string
}

// DecodeMercadoLibreQuestion normalizes the authoritative API v4 response
// fetched after a Mercado Libre `questions` notification. It does not verify a
// callback by itself and never grants permission to contact or answer a buyer.
func DecodeMercadoLibreQuestion(payload []byte, expectedSellerID, tenantID, organizationID string, receivedAt time.Time) (MercadoLibreQuestionDecoded, error) {
	if len(payload) == 0 || int64(len(payload)) > MaxPayloadBytes {
		return MercadoLibreQuestionDecoded{}, fmt.Errorf("%w: payload size", ErrInvalidPayload)
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var in mercadoLibreQuestion
	if err := dec.Decode(&in); err != nil {
		return MercadoLibreQuestionDecoded{}, fmt.Errorf("%w: json", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return MercadoLibreQuestionDecoded{}, err
	}
	if strings.TrimSpace(expectedSellerID) == "" || in.SellerID.String() != expectedSellerID {
		return MercadoLibreQuestionDecoded{}, ErrInvalidCredential
	}
	questionID := in.ID.String()
	createdAt, parseErr := time.Parse(time.RFC3339Nano, in.DateCreated)
	if questionID == "" || questionID == "0" || strings.TrimSpace(in.ItemID) == "" || strings.TrimSpace(in.Status) == "" {
		return MercadoLibreQuestionDecoded{}, fmt.Errorf("%w: question identity", ErrInvalidPayload)
	}
	if parseErr != nil {
		createdAt = receivedAt.UTC()
	}
	raw, err := NewRedactedRawEvent(tenantID, organizationID, "mercadolibre", questionID, "mercadolibre.question.v4", "https://api.mercadolibre.com/questions/"+questionID+"?api_version=4", "4", createdAt, receivedAt, payload, payload, "mercadolibre:provider-response:v1")
	if err != nil {
		return MercadoLibreQuestionDecoded{}, err
	}
	if parseErr != nil {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "INVALID_DATE_CREATED"}, nil
	}
	status := strings.ToUpper(strings.TrimSpace(in.Status))
	if status == "BANNED" || status == "DELETED" || status == "DISABLED" || status == "UNDER_REVIEW" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "QUESTION_" + status}, nil
	}
	if status != "UNANSWERED" && status != "ANSWERED" && status != "CLOSED_UNANSWERED" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "UNKNOWN_QUESTION_STATUS"}, nil
	}
	text := strings.TrimSpace(in.Text)
	if text == "" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "MISSING_QUESTION_TEXT"}, nil
	}
	buyerID := in.BuyerID.String()
	if in.From != nil && in.From.ID.String() != "" {
		if buyerID != "" && buyerID != "0" && buyerID != in.From.ID.String() {
			return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "AMBIGUOUS_BUYER_ID"}, nil
		}
		buyerID = in.From.ID.String()
	}
	if buyerID == "" || buyerID == "0" {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "MISSING_BUYER_ID"}, nil
	}
	fields := []Field{{ID: "ITEM_ID", Value: strings.TrimSpace(in.ItemID)}, {ID: "BUYER_ID", Value: buyerID}, {ID: "QUESTION_TEXT", Value: text}}
	if in.From != nil {
		appendField := func(id, value string) {
			if value = strings.TrimSpace(value); value != "" {
				fields = append(fields, Field{ID: id, Value: value})
			}
		}
		appendField("FIRST_NAME", in.From.FirstName)
		appendField("LAST_NAME", in.From.LastName)
		appendField("EMAIL", in.From.Email)
		appendField("PHONE_AREA_CODE", in.From.Phone.AreaCode)
		appendField("PHONE_NUMBER", in.From.Phone.Number)
	}
	candidate := LeadCandidate{
		TenantID: tenantID, OrganizationID: organizationID, Provider: "mercadolibre",
		ProviderLeadID: questionID, LeadStage: status, SourceKind: "MARKETPLACE_QUESTION",
		SubmittedAt: createdAt.UTC(), Fields: fields, ContactEligibility: "pending_policy",
	}
	if err := candidate.Validate(); err != nil {
		return MercadoLibreQuestionDecoded{Raw: raw, NormalizationCode: "INVALID_NORMALIZED_LEAD"}, nil
	}
	return MercadoLibreQuestionDecoded{Raw: raw, Candidate: &candidate}, nil
}
