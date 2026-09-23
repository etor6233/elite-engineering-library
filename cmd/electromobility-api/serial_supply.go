package main

// AUTHORED opt-in activation of the explicit serialized reference policy.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errSerialSupplyConfiguration = errors.New("serial supply activation is invalid")

func init() { serialSupplyModuleFactory = selectedSerialSupplyModule }
func selectedSerialSupplyModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errSerialSupplyConfiguration
	}
	switch lookup("SERIAL_SUPPLY_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errSerialSupplyConfiguration
	}
	hash := sha256.Sum256([]byte(sc.PolicyJSON))
	if pool == nil || lookup("SERIAL_SUPPLY_POLICY_SHA256") != hex.EncodeToString(hash[:]) {
		return nil, errSerialSupplyConfiguration
	}
	var ready bool
	err := pool.QueryRow(ctx, `select
 (select count(*) from pg_trigger where not tgisinternal and tgenabled in ('O','A') and (
  tgrelid=to_regclass('procurement.purchase_order') and tgname='serial_supply_purchase_guard'
  or tgrelid=to_regclass('factory.production_unit') and tgname='serial_supply_factory_guard'
  or tgrelid=to_regclass('inventory.stock_unit') and tgname='serial_supply_stock_guard'
  or tgrelid=to_regclass('approval.request') and tgname='serial_supply_approval_guard'
  or tgrelid=to_regclass('procurement.serial_supply_plan') and tgname in ('serial_supply_plan_guard','serial_supply_plan_immutable')
  or tgrelid=to_regclass('procurement.serial_supply_line') and tgname in ('serial_supply_line_guard','serial_supply_line_immutable')
  or tgrelid=to_regclass('procurement.serial_supply_step') and tgname='serial_supply_step_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_unit') and tgname='serial_supply_unit_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_shipment') and tgname='serial_supply_shipment_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_manifest') and tgname='serial_supply_manifest_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_receipt') and tgname='serial_supply_receipt_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_quality') and tgname='serial_supply_quality_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_effect') and tgname='serial_supply_effect_immutable'
 ))=15 and exists(select 1 from pg_constraint where conrelid=to_regclass('approval.request') and conname='request_kind_check'
 and pg_get_constraintdef(oid) like '%serial_quality%')`).Scan(&ready)
	if err != nil || !ready {
		return nil, errSerialSupplyConfiguration
	}
	store, err := postgres.NewSerialSupply(pool)
	if err != nil {
		return nil, errSerialSupplyConfiguration
	}
	return httpapi.SerialSupplyModule{Service: store}, nil
}
