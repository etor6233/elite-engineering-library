package storedvaluebridge

// AUTHORED typed transaction and receipt contracts. No pricing or points rule
// is implemented here; Result and Calculation identify the source-derived IPC.
import "encoding/json"

type Entry struct {
	AccountID    string `json:"account_id"`
	PointsDelta  string `json:"points_delta"`
	AppliedMinor int64  `json:"applied_minor_units"`
}
type BoundOperation struct {
	Schema        string      `json:"schema"`
	Request       Request     `json:"request"`
	Requester     string      `json:"requester"`
	Order         Order       `json:"order"`
	OrderState    string      `json:"order_state"`
	GrossMinor    int64       `json:"gross_minor_units"`
	GiftMinor     int64       `json:"gift_minor_units"`
	DiscountMinor int64       `json:"discount_minor_units"`
	EngineSHA256  string      `json:"engine_sha256"`
	ExpiresAt     string      `json:"expires_at"`
	Calculation   Calculation `json:"calculation"`
	Result        Result      `json:"result"`
	Entries       []Entry     `json:"entries"`
}
type ApprovalView struct {
	ApprovalID    string          `json:"approval_id"`
	PayloadSHA256 string          `json:"payload_sha256"`
	State         string          `json:"state"`
	Replay        bool            `json:"replay"`
	Payload       BoundOperation  `json:"payload"`
	Receipt       json.RawMessage `json:"receipt"`
}
type Decision struct {
	OrganizationID string `json:"organization_id"`
	ApprovalID     string `json:"approval_id"`
	PayloadSHA256  string `json:"payload_sha256"`
	Approved       bool   `json:"approved"`
	Reason         string `json:"reason"`
}
type Allocation struct {
	OrderID          string `json:"order_id"`
	OrderVersion     int64  `json:"order_version"`
	Currency         string `json:"currency"`
	GrossMinor       int64  `json:"gross_minor_units"`
	GiftMinor        int64  `json:"gift_minor_units"`
	DiscountMinor    int64  `json:"discount_minor_units"`
	ProviderDueMinor int64  `json:"provider_due_minor_units"`
}
