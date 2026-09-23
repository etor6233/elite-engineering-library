package mercadolibremarketplace

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeDoer struct {
	response *http.Response
	err      error
	request  *http.Request
	body     []byte
}

func (f *fakeDoer) Do(request *http.Request) (*http.Response, error) {
	f.request = request
	f.body, _ = io.ReadAll(request.Body)
	return f.response, f.err
}

func provenProfile() Profile {
	return Profile{Schema: "elite-mercadolibre-marketplace-profile/v1", Provider: "Mercado Libre", Decision: "PROVEN", OfficialDocumentationReviewed: true, ArchivedSDKRejectionAcknowledged: true, ApplicationAndLegalOwnerProven: true, SellerAuthorizationProven: true, OAuthTokenStorageProven: true, ReadScopeProven: true, WriteScopeProven: true, NoPreproductionEnvironmentAcknowledged: true, ItemAndCategoryPolicyApproved: true, NotificationOriginControlProven: true, NotificationTopicsApproved: true, QuotaAndCostApproved: true, ReconciliationApproved: true, SiteID: "MLA", SellerID: "123456789", ApplicationID: "987654321", AccessTokenEnvironmentVariable: "MERCADOLIBRE_ACCESS_TOKEN", ApprovedNotificationTopics: []string{"items", "orders_v2", "questions"}, ValidationCallApproved: true, AutomaticBusinessWrite: false}
}

const token = "APP_USR-test-token-value-1234567890"

func TestGetItemUsesExactOfficialContractAndRedactsReceipt(t *testing.T) {
	doer := &fakeDoer{response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"id":"MLA1234567890","title":"x"}`))}}
	output := filepath.Join(t.TempDir(), "evidence")
	receipt, err := Execute(context.Background(), doer, provenProfile(), Request{Operation: "GET_ITEM", ResourceID: "MLA1234567890"}, token, output)
	if err != nil {
		t.Fatal(err)
	}
	if doer.request.Method != "GET" || doer.request.URL.String() != APIBaseURL+"/items/MLA1234567890" || doer.request.Header.Get("Authorization") != "Bearer "+token {
		t.Fatal("official request contract mismatch")
	}
	serialized, _ := json.Marshal(receipt)
	if strings.Contains(string(serialized), "123456789") || strings.Contains(string(serialized), token) || receipt.AutomaticBusinessWrite {
		t.Fatal("receipt leaked identifiers or enabled write")
	}
}

func TestValidateItemUsesValidatorAndPreservesRejection(t *testing.T) {
	doer := &fakeDoer{response: &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(`{"message":"body.invalid_field_types"}`))}}
	receipt, err := Execute(context.Background(), doer, provenProfile(), Request{Operation: "VALIDATE_ITEM", Item: json.RawMessage(`{"title":"x"}`)}, token, filepath.Join(t.TempDir(), "evidence"))
	if err != nil {
		t.Fatal(err)
	}
	if doer.request.Method != "POST" || doer.request.URL.Path != "/items/validate" || receipt.Outcome != "REJECTED_BY_VALIDATOR" {
		t.Fatal("validator contract mismatch")
	}
}

func TestQuestionReadAndReconciliationUseOfficialV4Contract(t *testing.T) {
	tests := []struct {
		name     string
		request  Request
		path     string
		rawQuery string
	}{
		{name: "question by id", request: Request{Operation: "GET_QUESTION", ResourceID: "11751825075"}, path: "/questions/11751825075", rawQuery: "api_version=4"},
		{name: "unanswered seller reconciliation", request: Request{Operation: "SEARCH_UNANSWERED_QUESTIONS"}, path: "/questions/search", rawQuery: "seller_id=123456789&status=UNANSWERED&api_version=4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doer := &fakeDoer{response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"questions":[]}`))}}
			if _, err := Execute(context.Background(), doer, provenProfile(), tt.request, token, filepath.Join(t.TempDir(), "evidence")); err != nil {
				t.Fatal(err)
			}
			if doer.request.Method != http.MethodGet || doer.request.URL.Path != tt.path || doer.request.URL.RawQuery != tt.rawQuery {
				t.Fatalf("method=%s path=%s query=%s", doer.request.Method, doer.request.URL.Path, doer.request.URL.RawQuery)
			}
		})
	}
}

func TestProfileAndOperationsFailClosed(t *testing.T) {
	profile := provenProfile()
	profile.Decision = "BLOCKED"
	if err := Validate(profile, Request{Operation: "GET_ITEM", ResourceID: "MLA1234567890"}, token); err == nil {
		t.Fatal("blocked profile passed")
	}
	profile = provenProfile()
	profile.ValidationCallApproved = false
	if err := Validate(profile, Request{Operation: "VALIDATE_ITEM", Item: json.RawMessage(`{"title":"x"}`)}, token); err == nil {
		t.Fatal("unapproved validation passed")
	}
	if err := Validate(provenProfile(), Request{Operation: "PUBLISH_ITEM", Item: json.RawMessage(`{}`)}, token); err == nil {
		t.Fatal("unimplemented write passed")
	}
}

func TestProviderFailureAndUnexpectedStatusAreAtomic(t *testing.T) {
	parent := t.TempDir()
	output := filepath.Join(parent, "evidence")
	_, err := Execute(context.Background(), &fakeDoer{err: errors.New("provider unavailable")}, provenProfile(), Request{Operation: "GET_ORDER", ResourceID: "123456789"}, token, output)
	if err == nil || !strings.Contains(err.Error(), "provider unavailable") {
		t.Fatal("provider failure not retained")
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatal("partial output exists")
	}
	_, err = Execute(context.Background(), &fakeDoer{response: &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader(`{}`))}}, provenProfile(), Request{Operation: "GET_ORDER", ResourceID: "123456789"}, token, output)
	if err == nil {
		t.Fatal("unexpected status passed")
	}
}

func TestLoadJSONRejectsUnknownAndTrailing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "request.json")
	if err := os.WriteFile(path, []byte(`{"operation":"GET_ITEM","resource_id":"MLA1234567890","unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var request Request
	if LoadJSON(path, &request) == nil {
		t.Fatal("unknown field passed")
	}
	if err := os.WriteFile(path, []byte(`{"operation":"GET_ITEM","resource_id":"MLA1234567890"} {}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if LoadJSON(path, &request) == nil {
		t.Fatal("trailing value passed")
	}
}
