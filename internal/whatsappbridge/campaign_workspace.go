package whatsappbridge

// AUTHORED views over the existing campaign/CRM owner. No new send path.
import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

type CampaignSummary struct {
	ID        string    `json:"campaign_id"`
	Template  string    `json:"template_name"`
	Language  string    `json:"language_code"`
	Members   int       `json:"members"`
	Steps     int       `json:"steps"`
	CreatedAt time.Time `json:"created_at"`
	Stopped   bool      `json:"stopped"`
}
type CampaignAudienceOption struct {
	LeadID      string `json:"lead_id"`
	Name        string `json:"name"`
	ContactHint string `json:"contact_hint"`
	Source      string `json:"source"`
	State       string `json:"state"`
}
type CampaignWorkspacePage struct {
	Campaigns  []CampaignSummary        `json:"campaigns,omitempty"`
	Audience   []CampaignAudienceOption `json:"audience,omitempty"`
	NextCursor string                   `json:"next_cursor,omitempty"`
	Kind       string                   `json:"kind"`
}
type CampaignSelection struct {
	ID      string         `json:"campaign_id"`
	LeadIDs []string       `json:"lead_ids"`
	Steps   []CampaignStep `json:"steps"`
}

func campaignWorkspaceQuery(raw string) (kind, after string, err error) {
	q, e := url.ParseQuery(raw)
	if e != nil {
		return "", "", ErrApproval
	}
	for k, v := range q {
		if (k != "kind" && k != "after") || len(v) != 1 {
			return "", "", ErrApproval
		}
	}
	kind = q.Get("kind")
	if kind == "" {
		kind = "campaigns"
	}
	if kind != "campaigns" && kind != "audience" {
		return "", "", ErrApproval
	}
	after = q.Get("after")
	if after != "" && !cr.ValidID(after) {
		return "", "", ErrApproval
	}
	return kind, after, nil
}

// Only this owner resolves CRM phone data. The UI selects a lead, not a phone/subject/HMAC.
func (c *Campaigns) selectionMember(ctx context.Context, tx pgx.Tx, id string) (CampaignMember, string, error) {
	if !cr.ValidID(id) {
		return CampaignMember{}, "", ErrApproval
	}
	s := c.schedule.approvals
	var name, phone string
	e := tx.QueryRow(ctx, `select coalesce(nullif(l.contact_payload->>'name',''),nullif(cp.display_name,''),'Contacto'),coalesce(l.contact_payload->>'phone','')
 from crm.lead l left join crm.customer_profile cp on cp.tenant_id=l.tenant_id and cp.customer_principal_id=l.customer_principal_id
 where l.tenant_id=$1 and l.organization_id=$2 and l.lead_id=$3 for share of l`, s.tenant, s.org, id).Scan(&name, &phone)
	if e != nil {
		return CampaignMember{}, "", e
	}
	// Normalization is deliberately narrow; ambiguous phone values are ineligible.
	phone = strings.TrimPrefix(phone, "+")
	if len(phone) < 8 || len(phone) > 15 || phone[0] == '0' {
		return CampaignMember{}, "", ErrApproval
	}
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return CampaignMember{}, "", ErrApproval
		}
	}
	m, e := c.member(ctx, tx, CampaignRecipient{LeadID: id, Recipient: phone})
	if e != nil {
		return m, "", e
	}
	if m.Lifecycle != "new" && m.Lifecycle != "contacted" && m.Lifecycle != "qualified" {
		return m, "", ErrApproval
	}
	if !cr.ValidText(name, 256) {
		name = "Contacto"
	}
	return m, name, nil
}

func (c *Campaigns) Workspace(ctx context.Context, p identity.Principal, kind, after string) (CampaignWorkspacePage, error) {
	out := CampaignWorkspacePage{Kind: kind}
	if !c.authorized(p, "marketing:read") || !c.authorized(p, "notification:read") || after != "" && !cr.ValidID(after) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	if kind == "campaigns" {
		rows, e := s.base.pool.Query(ctx, `select c.campaign_id,coalesce(c.payload#>>'{request,steps,0,template_name}',''),coalesce(c.payload#>>'{request,steps,0,language_code}',''),jsonb_array_length(c.payload->'members'),jsonb_array_length(c.payload#>'{request,steps}'),c.created_at,exists(select 1 from communication.whatsapp_campaign_stop x where x.tenant_id=c.tenant_id and x.campaign_id=c.campaign_id)
   from communication.whatsapp_campaign c where c.tenant_id=$1 and c.organization_id=$2 and c.campaign_id>$3 order by c.campaign_id limit 26`, s.tenant, s.org, after)
		if e != nil {
			return out, e
		}
		defer rows.Close()
		out.Campaigns = []CampaignSummary{}
		for rows.Next() {
			var v CampaignSummary
			if e = rows.Scan(&v.ID, &v.Template, &v.Language, &v.Members, &v.Steps, &v.CreatedAt, &v.Stopped); e != nil {
				return out, e
			}
			out.Campaigns = append(out.Campaigns, v)
		}
		if e = rows.Err(); e != nil {
			return out, e
		}
		if len(out.Campaigns) > 25 {
			out.NextCursor = out.Campaigns[24].ID
			out.Campaigns = out.Campaigns[:25]
		}
		return out, nil
	}
	if kind != "audience" || !c.authorized(p, "marketing:request") || !c.authorized(p, "lead:read") {
		return out, ErrApproval
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	rows, e := tx.Query(ctx, `select lead_id from crm.lead where tenant_id=$1 and organization_id=$2 and lead_id>$3 and lifecycle_state in('new','contacted','qualified') order by lead_id limit 26`, s.tenant, s.org, after)
	if e != nil {
		return out, e
	}
	var ids []string
	for rows.Next() {
		var id string
		if e = rows.Scan(&id); e != nil {
			rows.Close()
			return out, e
		}
		ids = append(ids, id)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return out, e
	}
	if len(ids) > 25 {
		out.NextCursor = ids[24]
		ids = ids[:25]
	}
	out.Audience = []CampaignAudienceOption{}
	for _, id := range ids {
		m, name, e := c.selectionMember(ctx, tx, id)
		if e == ErrApproval || e == pgx.ErrNoRows {
			continue
		}
		if e != nil {
			return out, e
		}
		out.Audience = append(out.Audience, CampaignAudienceOption{LeadID: id, Name: name, ContactHint: "•••• " + m.Recipient[len(m.Recipient)-4:], Source: m.Source, State: m.Lifecycle})
	}
	return out, tx.Commit(ctx)
}

func (c *Campaigns) PrepareSelection(ctx context.Context, p identity.Principal, in CampaignSelection) (CampaignSnapshot, error) {
	if !c.authorized(p, "marketing:request") || !c.authorized(p, "lead:read") || !cr.ValidID(in.ID) || len(in.ID) < 16 || len(in.LeadIDs) == 0 || len(in.LeadIDs) > c.policy.MaxMembers {
		return CampaignSnapshot{}, ErrApproval
	}
	tx, e := c.schedule.approvals.base.pool.Begin(ctx)
	if e != nil {
		return CampaignSnapshot{}, e
	}
	defer tx.Rollback(ctx)
	r := CampaignRequest{ID: in.ID, Steps: in.Steps}
	sources, states, seen := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, id := range in.LeadIDs {
		if seen[id] {
			return CampaignSnapshot{}, ErrApproval
		}
		seen[id] = true
		m, _, e := c.selectionMember(ctx, tx, id)
		if e != nil {
			return CampaignSnapshot{}, e
		}
		r.Members = append(r.Members, CampaignRecipient{LeadID: id, Recipient: m.Recipient})
		sources[m.Source] = true
		states[m.Lifecycle] = true
	}
	for x := range sources {
		r.Sources = append(r.Sources, x)
	}
	for x := range states {
		r.States = append(r.States, x)
	}
	sort.Strings(r.Sources)
	sort.Strings(r.States)
	out, e := c.prepare(ctx, tx, p, r)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}

func (c *Campaigns) registerWorkspace(route func(string, string, string, func(http.ResponseWriter, *http.Request, identity.Principal))) {
	route("GET", "/workspace/result/{id}", "marketing:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		if !c.authorized(p, "notification:read") {
			notificationProblem(w, 403, "FORBIDDEN")
			return
		}
		s := c.schedule.approvals
		var exists bool
		if e := s.base.pool.QueryRow(r.Context(), "select exists(select 1 from communication.whatsapp_campaign where tenant_id=$1 and organization_id=$2 and campaign_id=$3)", s.tenant, s.org, r.PathValue("id")).Scan(&exists); e != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		if !exists {
			notificationProblem(w, 404, "NOT_FOUND")
			return
		}
		out, e := c.Status(r.Context(), p, r.PathValue("id"))
		if e != nil {
			notificationProblem(w, 409, "RESULT_UNCONFIRMED")
			return
		}
		replyJSON(w, out)
	})
	route("GET", "/workspace", "marketing:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		kind, after, e := campaignWorkspaceQuery(r.URL.RawQuery)
		if e != nil {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		out, e := c.Workspace(r.Context(), p, kind, after)
		if e != nil {
			notificationProblem(w, 409, "WORKSPACE_UNAVAILABLE")
			return
		}
		replyJSON(w, out)
	})
	route("GET", "/workspace/templates", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		// Explicit projection: never return connection credentials or the provider profile.
		var profile struct {
			Templates []struct {
				Name     string `json:"name"`
				Language string `json:"language_code"`
				Count    int    `json:"body_parameter_count"`
			} `json:"approved_templates"`
		}
		if json.Unmarshal(c.schedule.approvals.profile, &profile) != nil {
			notificationProblem(w, 409, "TEMPLATES_UNAVAILABLE")
			return
		}
		replyJSON(w, map[string]any{"templates": profile.Templates, "max_members": c.policy.MaxMembers, "max_steps": c.policy.MaxSteps})
	})
	route("POST", "/workspace/prepare", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in CampaignSelection
		if decodeSchedule(w, r, &in) != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		out, e := c.PrepareSelection(r.Context(), p, in)
		if e != nil {
			notificationProblem(w, 409, "SOURCE_NOT_CURRENT")
			return
		}
		_, h, e := campaignHash(out)
		if e != nil {
			notificationProblem(w, 409, "SOURCE_NOT_CURRENT")
			return
		}
		replyJSON(w, map[string]any{"snapshot": out, "campaign_sha256": h})
	})
	route("POST", "/workspace/create", "marketing:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Selection CampaignSelection `json:"selection"`
			Hash      string            `json:"campaign_sha256"`
		}
		if decodeSchedule(w, r, &in) != nil || !validDigest(in.Hash) {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !c.authorized(p, "notification:request") || !c.authorized(p, "lead:read") {
			notificationProblem(w, 403, "FORBIDDEN")
			return
		}
		s := c.schedule.approvals
		var exists bool
		if e := s.base.pool.QueryRow(r.Context(), "select exists(select 1 from communication.whatsapp_campaign where tenant_id=$1 and campaign_id=$2)", s.tenant, in.Selection.ID).Scan(&exists); e != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		if exists {
			stored, h, e := c.load(r.Context(), s.base.pool, in.Selection.ID)
			if e != nil || h != in.Hash || stored.Requester != p.Subject {
				notificationProblem(w, 409, "REQUEST_DIVERGENT")
				return
			}
			original := CampaignSelection{ID: stored.Request.ID, Steps: stored.Request.Steps}
			for _, m := range stored.Request.Members {
				original.LeadIDs = append(original.LeadIDs, m.LeadID)
			}
			a, _ := json.Marshal(original)
			b, _ := json.Marshal(in.Selection)
			if !bytes.Equal(a, b) {
				notificationProblem(w, 409, "REQUEST_DIVERGENT")
				return
			}
			out, e := c.Resume(r.Context(), p, in.Selection.ID, h)
			if e != nil {
				notificationProblem(w, 409, "RESULT_UNCONFIRMED")
				return
			}
			replyJSON(w, out)
			return
		}
		snapshot, e := c.PrepareSelection(r.Context(), p, in.Selection)
		if e != nil {
			notificationProblem(w, 409, "SOURCE_NOT_CURRENT")
			return
		}
		_, h, e := campaignHash(snapshot)
		if e != nil || h != in.Hash {
			notificationProblem(w, 409, "SOURCE_NOT_CURRENT")
			return
		}
		out, e := c.Create(r.Context(), p, snapshot)
		if e != nil {
			notificationProblem(w, 409, "RESULT_UNCONFIRMED")
			return
		}
		replyJSON(w, out)
	})
}
