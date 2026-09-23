// AUTHORED transaction adapter. Existing operations.Service performs all preparation
// and transition validation; existing PostgreSQL owner helpers perform the writes.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"github.com/jackc/pgx/v5"
)

type operationsTxRepository struct{ tx pgx.Tx }

var _ operations.Repository = operationsTxRepository{}

func (r operationsTxRepository) CreatePurchaseOrder(ctx context.Context, tenant, eventID string, value operations.PurchaseOrder) error {
	return createPurchaseOrderInTx(ctx, r.tx, tenant, eventID, value)
}
func (r operationsTxRepository) TransitionPurchaseOrder(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	return transitionPurchaseOrderInTx(ctx, r.tx, tenant, organization, id, current, version, target, eventID)
}
func (r operationsTxRepository) CreateProductionUnit(ctx context.Context, tenant, eventID string, value operations.ProductionUnit) error {
	return createProductionUnitInTx(ctx, r.tx, tenant, eventID, value)
}
func (r operationsTxRepository) TransitionProductionUnit(ctx context.Context, tenant, organization, id, current, target, eventID string) error {
	return transitionProductionUnitInTx(ctx, r.tx, tenant, organization, id, current, target, eventID)
}
func (r operationsTxRepository) CreateStockUnit(ctx context.Context, tenant, eventID string, value operations.StockUnit) error {
	return createStockUnitInTx(ctx, r.tx, tenant, eventID, value)
}
func (r operationsTxRepository) TransitionStockUnit(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	return transitionStockUnitInTx(ctx, r.tx, tenant, organization, id, current, version, target, eventID)
}
