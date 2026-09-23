package mercadolibremarketplace

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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const APIBaseURL = "https://api.mercadolibre.com"

type Profile struct {
	Schema                                 string   `json:"schema"`
	Provider                               string   `json:"provider"`
	Decision                               string   `json:"decision"`
	OfficialDocumentationReviewed          bool     `json:"official_documentation_reviewed"`
	ArchivedSDKRejectionAcknowledged       bool     `json:"archived_sdk_rejection_acknowledged"`
	ApplicationAndLegalOwnerProven         bool     `json:"application_and_legal_owner_proven"`
	SellerAuthorizationProven              bool     `json:"seller_authorization_proven"`
	OAuthTokenStorageProven                bool     `json:"oauth_token_storage_proven"`
	ReadScopeProven                        bool     `json:"read_scope_proven"`
	WriteScopeProven                       bool     `json:"write_scope_proven"`
	NoPreproductionEnvironmentAcknowledged bool     `json:"no_preproduction_environment_acknowledged"`
	ItemAndCategoryPolicyApproved          bool     `json:"item_and_category_policy_approved"`
	NotificationOriginControlProven        bool     `json:"notification_origin_control_proven"`
	NotificationTopicsApproved             bool     `json:"notification_topics_approved"`
	QuotaAndCostApproved                   bool     `json:"quota_and_cost_approved"`
	ReconciliationApproved                 bool     `json:"reconciliation_approved"`
	SiteID                                 string   `json:"site_id"`
	SellerID                               string   `json:"seller_id"`
	ApplicationID                          string   `json:"application_id"`
	AccessTokenEnvironmentVariable         string   `json:"access_token_environment_variable"`
	ApprovedNotificationTopics             []string `json:"approved_notification_topics"`
	ValidationCallApproved                 bool     `json:"validation_call_approved"`
	AutomaticBusinessWrite                 bool     `json:"automatic_business_write"`
}

type Request struct {
	Operation  string          `json:"operation"`
	ResourceID string          `json:"resource_id"`
	Item       json.RawMessage `json:"item,omitempty"`
}

type Receipt struct {
	Schema                 string `json:"schema"`
	CreatedAt              string `json:"created_at"`
	Operation              string `json:"operation"`
	HTTPStatus             int    `json:"http_status"`
	Outcome                string `json:"outcome"`
	SellerIDSHA256         string `json:"seller_id_sha256"`
	ApplicationIDSHA256    string `json:"application_id_sha256"`
	ResourceIDSHA256       string `json:"resource_id_sha256"`
	RequestSHA256          string `json:"request_sha256"`
	ResponseSHA256         string `json:"response_sha256"`
	AutomaticBusinessWrite bool   `json:"automatic_business_write"`
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	sitePattern   = regexp.MustCompile(`^ML[A-Z]$`)
	digitsPattern = regexp.MustCompile(`^[0-9]{3,20}$`)
	itemPattern   = regexp.MustCompile(`^ML[A-Z][0-9]{6,20}$`)
	envPattern    = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,80}$`)
)

func LoadJSON(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(io.LimitReader(file, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("JSON must contain exactly one value")
	}
	return nil
}

func Validate(profile Profile, request Request, token string) error {
	if profile.Schema != "elite-mercadolibre-marketplace-profile/v1" || profile.Provider != "Mercado Libre" || profile.Decision != "PROVEN" {
		return errors.New("provider profile is not PROVEN")
	}
	if !profile.OfficialDocumentationReviewed || !profile.ArchivedSDKRejectionAcknowledged || !profile.ApplicationAndLegalOwnerProven || !profile.SellerAuthorizationProven || !profile.OAuthTokenStorageProven || !profile.ReadScopeProven || !profile.QuotaAndCostApproved || !profile.ReconciliationApproved {
		return errors.New("required provider proof is missing")
	}
	if profile.AutomaticBusinessWrite {
		return errors.New("automatic business write is forbidden")
	}
	if !sitePattern.MatchString(profile.SiteID) || !digitsPattern.MatchString(profile.SellerID) || !digitsPattern.MatchString(profile.ApplicationID) || !envPattern.MatchString(profile.AccessTokenEnvironmentVariable) {
		return errors.New("provider identifiers are invalid")
	}
	if len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return errors.New("access token is invalid")
	}
	switch request.Operation {
	case "GET_ITEM":
		if !itemPattern.MatchString(request.ResourceID) || len(request.Item) != 0 {
			return errors.New("GET_ITEM request is invalid")
		}
	case "GET_ORDER":
		if !digitsPattern.MatchString(request.ResourceID) || len(request.Item) != 0 {
			return errors.New("GET_ORDER request is invalid")
		}
	case "GET_QUESTION":
		if !digitsPattern.MatchString(request.ResourceID) || len(request.Item) != 0 {
			return errors.New("GET_QUESTION request is invalid")
		}
	case "SEARCH_UNANSWERED_QUESTIONS":
		if request.ResourceID != "" || len(request.Item) != 0 {
			return errors.New("SEARCH_UNANSWERED_QUESTIONS request is invalid")
		}
	case "VALIDATE_ITEM":
		if request.ResourceID != "" || !profile.WriteScopeProven || !profile.NoPreproductionEnvironmentAcknowledged || !profile.ItemAndCategoryPolicyApproved || !profile.ValidationCallApproved {
			return errors.New("item validation is not approved")
		}
		if len(request.Item) < 2 || len(request.Item) > 1<<20 || request.Item[0] != '{' || !json.Valid(request.Item) {
			return errors.New("item JSON is invalid")
		}
	default:
		return errors.New("operation is not allowlisted")
	}
	return nil
}

func Execute(ctx context.Context, doer Doer, profile Profile, request Request, token, outputDirectory string) (Receipt, error) {
	if doer == nil {
		return Receipt{}, errors.New("HTTP client is required")
	}
	if err := Validate(profile, request, token); err != nil {
		return Receipt{}, err
	}
	method, path, body := http.MethodGet, "", []byte(nil)
	switch request.Operation {
	case "GET_ITEM":
		path = "/items/" + request.ResourceID
	case "GET_ORDER":
		path = "/orders/" + request.ResourceID
	case "GET_QUESTION":
		path = "/questions/" + request.ResourceID + "?api_version=4"
	case "SEARCH_UNANSWERED_QUESTIONS":
		path = "/questions/search?seller_id=" + profile.SellerID + "&status=UNANSWERED&api_version=4"
	case "VALIDATE_ITEM":
		method, path, body = http.MethodPost, "/items/validate", []byte(request.Item)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, method, APIBaseURL+path, bytes.NewReader(body))
	if err != nil {
		return Receipt{}, err
	}
	httpRequest.Header.Set("Authorization", "Bearer "+token)
	httpRequest.Header.Set("Accept", "application/json")
	if method == http.MethodPost {
		httpRequest.Header.Set("Content-Type", "application/json")
	}
	response, err := doer.Do(httpRequest)
	if err != nil {
		return Receipt{}, fmt.Errorf("Mercado Libre request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return Receipt{}, err
	}
	if len(responseBody) > 4<<20 {
		return Receipt{}, errors.New("provider response exceeds 4 MiB")
	}
	outcome := "ACCEPTED"
	if request.Operation == "VALIDATE_ITEM" && response.StatusCode == http.StatusBadRequest {
		outcome = "REJECTED_BY_VALIDATOR"
	} else if request.Operation == "VALIDATE_ITEM" && response.StatusCode != http.StatusNoContent {
		return Receipt{}, fmt.Errorf("unexpected validation status %d", response.StatusCode)
	} else if request.Operation != "VALIDATE_ITEM" && response.StatusCode != http.StatusOK {
		return Receipt{}, fmt.Errorf("unexpected read status %d", response.StatusCode)
	}
	abs, err := filepath.Abs(outputDirectory)
	if err != nil {
		return Receipt{}, err
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		return Receipt{}, errors.New("output directory must not exist")
	}
	requestBytes, _ := json.Marshal(request)
	receipt := Receipt{Schema: "elite-mercadolibre-marketplace-receipt/v1", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano), Operation: request.Operation, HTTPStatus: response.StatusCode, Outcome: outcome, SellerIDSHA256: hash([]byte(profile.SellerID)), ApplicationIDSHA256: hash([]byte(profile.ApplicationID)), ResourceIDSHA256: hash([]byte(request.ResourceID)), RequestSHA256: hash(requestBytes), ResponseSHA256: hash(responseBody), AutomaticBusinessWrite: false}
	stage := abs + ".stage-" + fmt.Sprint(time.Now().UnixNano())
	if err := os.Mkdir(stage, 0o700); err != nil {
		return Receipt{}, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stage)
		}
	}()
	if err := os.WriteFile(filepath.Join(stage, "PROVIDER_RESPONSE.json"), append(responseBody, '\n'), 0o600); err != nil {
		return Receipt{}, err
	}
	receiptBytes, _ := json.MarshalIndent(receipt, "", "  ")
	if err := os.WriteFile(filepath.Join(stage, "MARKETPLACE_RECEIPT.json"), append(receiptBytes, '\n'), 0o600); err != nil {
		return Receipt{}, err
	}
	if err := os.Rename(stage, abs); err != nil {
		return Receipt{}, err
	}
	committed = true
	return receipt, nil
}

func hash(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
