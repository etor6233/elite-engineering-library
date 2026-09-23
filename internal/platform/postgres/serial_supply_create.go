// AUTHORED single-transaction composition; no new pricing or stock rules.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"errors"
	"github.com/jackc/pgx/v5"
)

// Operations requests an aggregate ID first and then an event ID. Only the
// aggregate is a preallocated recovery reference; events retain random IDs.
type supplyCreationIDs struct{ first string }

func (g *supplyCreationIDs) New() string {
	if g.first != "" {
		v := g.first
		g.first = ""
		return v
	}
	return randomid.Generator{}.New()
}
func (s *SerialSupply) Create(ctx context.Context, p identity.Principal, r sc.CreateRequest) (sc.Receipt, error) {
	if s == nil || !r.Valid() || p.TenantID == "" {
		return sc.Receipt{}, sc.ErrInvalid
	}
	if !approvalPrincipal(p, p.TenantID, r.DestinationOrganizationID, "supply:plan") {
		return sc.Receipt{}, sc.ErrNotFound
	}
	_, hash, err := sc.Canonical(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	// Separate SQL parameters and a namespace; hash collisions only serialize
	// unrelated orders, never grant access or merge their durable identities.
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('serial-supply-create:'||$1::text||':'||$2::text,0))`, p.TenantID, r.PurchaseOrderID); err != nil {
		return sc.Receipt{}, err
	}
	old, err := readSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID)
	if err == nil {
		plan, e := readSupplyPlan(ctx, tx, p.TenantID, r.PurchaseOrderID, false)
		if e != nil {
			return sc.Receipt{}, e
		}
		if plan.DestinationOrganizationID != r.DestinationOrganizationID || old.Actor != p.Subject || old.Kind != "planned" || old.RequestSHA256 != hash {
			return sc.Receipt{}, sc.ErrConflict
		}
		old.Replay = true
		return old, tx.Commit(ctx)
	}
	if !errors.Is(err, sc.ErrNotFound) {
		return sc.Receipt{}, err
	}
	svc := operations.NewService(operationsTxRepository{tx: tx}, &supplyCreationIDs{first: r.PurchaseOrderID})
	_, err = svc.CreatePurchaseOrder(ctx, p.TenantID, operations.PurchaseOrder{SupplierID: r.SupplierID, DestinationOrganizationID: r.DestinationOrganizationID, Currency: r.Currency, TotalMinorUnits: r.TotalMinorUnits})
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	receipt, err := bindSupplyPlanInTx(ctx, tx, p, r.PlanRequest, hash)
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, supplyError(tx.Commit(ctx))
}
