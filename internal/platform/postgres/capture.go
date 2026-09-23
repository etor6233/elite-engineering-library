package postgres

// AUTHORED read-only resolver. Scanned input never selects tenant or grants a
// business action. Every returned row is scoped by the actor's existing roles.
import (
	"context"
	cap "elite.local/enterprise/internal/capture"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/qr"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"unicode"
)

type CaptureStore struct{ pool *pgxpool.Pool }

func NewCaptureStore(p *pgxpool.Pool) *CaptureStore { return &CaptureStore{p} }
func CaptureAllowed(p identity.Principal) bool {
	return documentID(p.TenantID) && p.Subject != "" && (len(p.Organizations) > 0 || p.Allowed("*")) && (captureCommerceAllowed(p) || p.Allowed("inventory:read") || p.Allowed("factory:read") || p.Allowed("supply:read") || p.Allowed("supply:factory-read"))
}
func captureCommerceAllowed(p identity.Principal) bool {
	return p.Allowed("inventory:allocate") || p.Allowed("payment:create") || p.Allowed("handover:manage") || p.Allowed("admin:read")
}

type CaptureRecord struct {
	Kind           string `json:"kind"`
	ID             string `json:"id"`
	Label          string `json:"label"`
	Serial         string `json:"serial"`
	OrganizationID string `json:"organization_id"`
	State          string `json:"state"`
	Version        string `json:"version"`
	VariantID      string `json:"variant_id"`
}
type CaptureResolution struct {
	Schema         string          `json:"schema"`
	Status         string          `json:"status"`
	Records        []CaptureRecord `json:"records"`
	Authorized     bool            `json:"authorized"`
	BusinessEffect bool            `json:"business_effect"`
	Operation      string          `json:"operation"`
}
type captureResolver func(context.Context, string, qr.Kind, string) (bool, error)

func (f captureResolver) Resolve(c context.Context, t string, k qr.Kind, id string) (bool, error) {
	return f(c, t, k, id)
}

var ErrCaptureNotFound = errors.New("capture target not found or unauthorized")

func (s *CaptureStore) ResolveCapture(ctx context.Context, p identity.Principal, text, encoding, operation string) (CaptureResolution, error) {
	out := CaptureResolution{Schema: "capture-resolution/v1", Records: []CaptureRecord{}, Operation: "lookup"}
	if s == nil || s.pool == nil || !CaptureAllowed(p) {
		return out, ErrCaptureNotFound
	}
	if operation != "lookup" || len(text) == 0 || len(text) > 1024 || strings.TrimSpace(text) != text {
		return out, cap.ErrContract
	}
	for _, r := range text {
		if unicode.IsControl(r) {
			return out, cap.ErrContract
		}
	}
	if encoding == "QR_REFERENCE" {
		var lookupError error
		_, err := qr.Verify(ctx, text, p.TenantID, captureResolver(func(c context.Context, tenant string, k qr.Kind, key string) (bool, error) {
			if tenant != p.TenantID {
				return false, nil
			}
			kind := ""
			switch k {
			case qr.KindProduct:
				kind = "PRODUCT"
			case qr.KindOrder:
				kind = "ORDER"
			default:
				return false, nil
			}
			var err error
			out, err = s.captureRows(c, p, kind, key)
			lookupError = err
			return err == nil && (out.Authorized || out.Status == "REFINE_SELECTION"), err
		}))
		if err != nil {
			if lookupError != nil && !errors.Is(lookupError, ErrCaptureNotFound) {
				return CaptureResolution{}, lookupError
			}
			return CaptureResolution{}, ErrCaptureNotFound
		}
		return out, nil
	}
	if encoding != "SERIAL" {
		return out, cap.ErrContract
	}
	return s.captureRows(ctx, p, "SERIAL", text)
}
func (s *CaptureStore) captureRows(ctx context.Context, p identity.Principal, kind, id string) (CaptureResolution, error) {
	out := CaptureResolution{Schema: "capture-resolution/v1", Records: []CaptureRecord{}, Operation: "lookup"}
	if len(id) > 128 || strings.Contains(id, "://") {
		return out, cap.ErrContract
	}
	orgs := []string{}
	for org := range p.Organizations {
		orgs = append(orgs, org)
	}
	all := p.Allowed("*")
	stockAllowed := captureCommerceAllowed(p) || p.Allowed("inventory:read")
	// A QR product is a catalog variant reference, never silently a stock-unit ID.
	// Matching units still require explicit selection. Order QR support returns
	// only the operational record; customer and appointment QR are not enabled.
	query := `select 'stock_unit',i.stock_unit_id,v.display_name,i.serial_number,i.organization_id,i.state,i.version::text,i.variant_id
 from inventory.stock_unit i join catalog.vehicle_variant v using(tenant_id,variant_id)
 where i.tenant_id=$1 and ($2 or i.organization_id=any($3)) and (($4='SERIAL' and i.serial_number=$5)or($4='PRODUCT' and i.variant_id=$5))
 and ($6 or ($7 and exists(select 1 from procurement.serial_supply_manifest m join procurement.serial_supply_plan sp using(tenant_id,purchase_order_id)where m.tenant_id=i.tenant_id and m.stock_unit_id=i.stock_unit_id and sp.destination_organization_id=i.organization_id)))
 union all
 select 'factory_unit',f.production_unit_id,v.display_name,f.serial_number,case when $9 and sp.factory_organization_id is not null and($2 or sp.factory_organization_id=any($3))then sp.factory_organization_id else po.destination_organization_id end,f.state,''::text,f.variant_id
 from factory.production_unit f join procurement.purchase_order po using(tenant_id,purchase_order_id)join catalog.vehicle_variant v using(tenant_id,variant_id)left join procurement.serial_supply_plan sp using(tenant_id,purchase_order_id)
 where f.tenant_id=$1 and (($4='SERIAL'and f.serial_number=$5)or($4='PRODUCT'and f.variant_id=$5))
 and (($8 and ($2 or po.destination_organization_id=any($3)))or($9 and sp.factory_organization_id is not null and($2 or sp.factory_organization_id=any($3)))or($7 and sp.purchase_order_id is not null and($2 or po.destination_organization_id=any($3))))
 union all
 select 'order',o.order_id,o.order_id,'',o.organization_id,o.state,o.version::text,'' from sales.customer_order o where o.tenant_id=$1 and $4='ORDER'and o.order_id=$5 and $10 and($2 or o.organization_id=any($3))
 order by 1,2 limit 9`
	rows, e := s.pool.Query(ctx, query, p.TenantID, all, orgs, kind, id, stockAllowed, p.Allowed("supply:read"), p.Allowed("factory:read"), p.Allowed("supply:factory-read"), captureCommerceAllowed(p))
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var r CaptureRecord
		if e = rows.Scan(&r.Kind, &r.ID, &r.Label, &r.Serial, &r.OrganizationID, &r.State, &r.Version, &r.VariantID); e != nil {
			return out, e
		}
		out.Records = append(out.Records, r)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.Records) == 0 {
		return out, ErrCaptureNotFound
	}
	if len(out.Records) > 8 {
		return CaptureResolution{Schema: out.Schema, Status: "REFINE_SELECTION", Records: []CaptureRecord{}, Operation: "lookup"}, nil
	}
	out.Status = "RESOLVED"
	if len(out.Records) > 1 {
		out.Status = "MULTIPLE"
	}
	out.Authorized = true
	return out, nil
}
