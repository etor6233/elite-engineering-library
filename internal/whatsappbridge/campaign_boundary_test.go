package whatsappbridge

// AUTHORED finite command-envelope and source-selection invariants.
import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func campaignBoundaryFixture() (*Campaigns, CampaignRequest) {
	c := &Campaigns{policy: CampaignPolicy{MaxMembers: 20, MaxSteps: 3}, schedule: &ScheduledNotifications{approvals: &ScheduleApprovals{profile: json.RawMessage(`{"approved_templates":[{"name":"order_update","language_code":"es_AR","body_parameter_count":2}]}`)}}}
	at := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	r := CampaignRequest{ID: "campaign-boundary-001", Sources: []string{"fixture"}, States: []string{"new"}, Members: []CampaignRecipient{{LeadID: "lead-1", Recipient: "5491112345678"}}, Steps: []CampaignStep{{TemplateName: "order_update", LanguageCode: "es_AR", BodyParameters: []string{"a", "b"}, NotBefore: at, ExpiresAt: at.Add(time.Hour)}}}
	return c, r
}

type campaignDecodeWriter struct{ *httptest.ResponseRecorder }

func (campaignDecodeWriter) SetReadDeadline(time.Time) error { return nil }
func TestCampaignSelectionBoundary(t *testing.T) {
	c, r := campaignBoundaryFixture()
	if c.validate(r) != nil {
		t.Fatal("valid bounded selection")
	}
	for name, change := range map[string]func(*CampaignRequest){
		"implicit source":     func(r *CampaignRequest) { r.Sources = nil },
		"converted audience":  func(r *CampaignRequest) { r.States = []string{"converted"} },
		"duplicate member":    func(r *CampaignRequest) { r.Members = append(r.Members, r.Members[0]) },
		"unreviewed template": func(r *CampaignRequest) { r.Steps[0].TemplateName = "invented" },
		"wrong arity":         func(r *CampaignRequest) { r.Steps[0].BodyParameters = []string{"only-one"} },
		"overlap":             func(r *CampaignRequest) { r.Steps = append(r.Steps, r.Steps[0]) },
		"unbounded steps":     func(r *CampaignRequest) { r.Steps = make([]CampaignStep, 4) },
		"unbounded audience":  func(r *CampaignRequest) { r.Members = make([]CampaignRecipient, 21) },
	} {
		t.Run(name, func(t *testing.T) {
			_, r := campaignBoundaryFixture()
			change(&r)
			if c.validate(r) == nil {
				t.Fatal("invalid selection admitted")
			}
		})
	}
	for _, raw := range []string{`{"campaign_id":"one","CAMPAIGN_ID":"two"}`, `{"states":["new"],"predicate_sql":"select anything"}`, `{}{}`, strings.Repeat(" ", 32769)} {
		r := httptest.NewRequest("POST", "/", strings.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
		if decodeSchedule(campaignDecodeWriter{httptest.NewRecorder()}, r, new(CampaignRequest)) == nil {
			t.Fatal("ambiguous command")
		}
	}
	base := &ScheduledNotifications{approvals: &ScheduleApprovals{base: &PostgresAppointmentApprovals{purpose: "appointment"}}}
	policy := CampaignPolicy{Schema: "elite-whatsapp-campaign-policy/v1", ConsentPurpose: "appointment", PolicyVersion: "policy-1", MaxMembers: 20, MaxSteps: 3}
	if _, e := NewCampaigns(base, policy); e == nil {
		t.Fatal("appointment consent reused for marketing")
	}
}
func FuzzCampaignSelection(f *testing.F) {
	_, r := campaignBoundaryFixture()
	raw, _ := json.Marshal(r)
	f.Add(raw)
	f.Add([]byte(`{"campaign_id":"one","CAMPAIGN_ID":"two"}`))
	f.Add([]byte(`{"steps":[{"not_before":"invalid"}]}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32769 {
			t.Skip()
		}
		req := httptest.NewRequest("POST", "/", strings.NewReader(string(raw)))
		req.Header.Set("Content-Type", "application/json")
		var value CampaignRequest
		if decodeSchedule(campaignDecodeWriter{httptest.NewRecorder()}, req, &value) != nil {
			return
		}
		canonical, e := json.Marshal(value)
		if e != nil {
			t.Fatal(e)
		}
		var again CampaignRequest
		if json.Unmarshal(canonical, &again) != nil || !reflect.DeepEqual(value, again) {
			t.Fatal("round-trip lost reviewed command")
		}
		c, _ := campaignBoundaryFixture()
		if c.validate(value) == nil {
			if len(value.Members) < 1 || len(value.Members) > 20 || len(value.Steps) < 1 || len(value.Steps) > 3 {
				t.Fatal("unbounded admitted campaign")
			}
			for i, s := range value.Steps {
				if !s.ExpiresAt.After(s.NotBefore) || i > 0 && !s.NotBefore.After(value.Steps[i-1].ExpiresAt) {
					t.Fatal("overlapping admitted steps")
				}
			}
		}
	})
}
