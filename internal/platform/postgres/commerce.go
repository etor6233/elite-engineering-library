package postgres

import (
	"context"
	"elite.local/enterprise/internal/bcamounts"
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/commerce"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Commerce struct {
	pool   *pgxpool.Pool
	policy *businesspolicy.Profile
}

func NewCommerce(pool *pgxpool.Pool) *Commerce {
	return &Commerce{pool: pool, policy: businesspolicy.Reference()}
}

// The v1 profile admits the existing single-active-source contract only.
// A mode that would need price ranking is rejected by businesspolicy.Load.
func NewCommerceWithProfile(pool *pgxpool.Pool, policy *businesspolicy.Profile) (*Commerce, error) {
	if pool == nil || !policy.Valid() {
		return nil, businesspolicy.ErrProfile
	}
	return &Commerce{pool: pool, policy: policy}, nil
}
func (r *Commerce) BusinessPolicySHA256() string {
	if r == nil {
		return ""
	}
	return r.policy.SHA256()
}

func (r *Commerce) Operations(ctx context.Context, tenant, organization string) (commerce.OperationSnapshot, error) {
	value := commerce.OperationSnapshot{Orders: []commerce.OperationOrder{}, Stock: []commerce.OperationStock{}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `select order_id,state,version from sales.customer_order where tenant_id=$1 and organization_id=$2 order by created_at desc,order_id limit 26`, tenant, organization)
	if err != nil {
		return value, err
	}
	for rows.Next() {
		var o commerce.OperationOrder
		if err = rows.Scan(&o.ID, &o.State, &o.Version); err != nil {
			rows.Close()
			return value, err
		}
		o.Lines = []commerce.OperationLine{}
		o.Payments = []commerce.OperationStatus{}
		o.Handovers = []commerce.OperationStatus{}
		value.Orders = append(value.Orders, o)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return value, err
	}
	if len(value.Orders) > 25 {
		value.Truncated = true
		value.Orders = value.Orders[:25]
	}
	for i := range value.Orders {
		o := &value.Orders[i]
		rows, err = tx.Query(ctx, `select l.line_id,l.variant_id,v.display_name,l.quantity,coalesce(l.allocated_stock_unit_id,'') from sales.customer_order_line l join catalog.vehicle_variant v on v.tenant_id=l.tenant_id and v.variant_id=l.variant_id where l.tenant_id=$1 and l.order_id=$2 order by l.line_id limit 101`, tenant, o.ID)
		if err != nil {
			return value, err
		}
		for rows.Next() {
			var l commerce.OperationLine
			if err = rows.Scan(&l.ID, &l.VariantID, &l.Name, &l.Quantity, &l.StockID); err != nil {
				rows.Close()
				return value, err
			}
			o.Lines = append(o.Lines, l)
		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return value, err
		}
		if len(o.Lines) > 100 {
			value.Truncated = true
			o.Lines = o.Lines[:100]
		}
		for _, query := range []struct {
			sql    string
			target *[]commerce.OperationStatus
		}{
			{`select payment_attempt_id,state from payment.payment_attempt where tenant_id=$1 and order_id=$2 order by payment_attempt_id limit 101`, &o.Payments},
			{`select h.handover_id,h.state,exists(select 1 from sales.customer_order co join inventory.stock_unit su on su.tenant_id=co.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where co.tenant_id=h.tenant_id and co.order_id=h.order_id and co.organization_id=h.organization_id and co.customer_principal_id=h.customer_principal_id) from sales.delivery_handover h where h.tenant_id=$1 and h.order_id=$2 and h.organization_id=$3 order by h.handover_id limit 101`, &o.Handovers},
		} {
			args := []any{tenant, o.ID}
			if query.target == &o.Handovers {
				args = append(args, organization)
			}
			rows, err = tx.Query(ctx, query.sql, args...)
			if err != nil {
				return value, err
			}
			for rows.Next() {
				var s commerce.OperationStatus
				coherent := true
				columns := []any{&s.ID, &s.State}
				if query.target == &o.Handovers {
					columns = append(columns, &coherent)
				}
				if err = rows.Scan(columns...); err != nil {
					rows.Close()
					return value, err
				}
				if !coherent {
					rows.Close()
					return commerce.OperationSnapshot{}, errors.New("operation handover scope inconsistent")
				}
				*query.target = append(*query.target, s)
			}
			rows.Close()
			if err = rows.Err(); err != nil {
				return value, err
			}
			if len(*query.target) > 100 {
				value.Truncated = true
				*query.target = (*query.target)[:100]
			}
		}
	}
	rows, err = tx.Query(ctx, `select stock_unit_id,variant_id,serial_number,version from inventory.stock_unit where tenant_id=$1 and organization_id=$2 and state='available' order by stock_unit_id limit 201`, tenant, organization)
	if err != nil {
		return value, err
	}
	for rows.Next() {
		var s commerce.OperationStock
		if err = rows.Scan(&s.ID, &s.VariantID, &s.Serial, &s.Version); err != nil {
			rows.Close()
			return value, err
		}
		value.Stock = append(value.Stock, s)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return value, err
	}
	if len(value.Stock) > 200 {
		value.Truncated = true
		value.Stock = value.Stock[:200]
	}
	return value, tx.Commit(ctx)
}
func (r *Commerce) CreatePriceBook(ctx context.Context, tenant, eventID string, value commerce.PriceBook) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = r.createPriceBookTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) ActivatePriceBook(ctx context.Context, tenant, id, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = r.activatePriceBookTx(ctx, tx, tenant, id, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) PublicPrice(ctx context.Context, tenantCode, market, variant string) (commerce.PriceEntry, error) {
	rows, err := r.pool.Query(ctx, `select e.variant_id,e.amount_minor_units,e.tax_mode from pricing.price_book b join pricing.price_book_entry e on e.tenant_id=b.tenant_id and e.price_book_id=b.price_book_id join platform.tenant t on t.tenant_id=b.tenant_id `+bcPriceTimeContextSQL+` where t.tenant_code=$1 and t.status='active' and b.market=$2 and `+bcPriceEligibilitySQL+` and e.variant_id=$3 limit 2`, tenantCode, market, variant)
	if err != nil {
		return commerce.PriceEntry{}, err
	}
	defer rows.Close()
	values := []commerce.PriceEntry{}
	for rows.Next() {
		var value commerce.PriceEntry
		if err := rows.Scan(&value.VariantID, &value.AmountMinorUnits, &value.TaxMode); err != nil {
			return value, err
		}
		values = append(values, value)
	}
	if len(values) != 1 {
		return commerce.PriceEntry{}, commerce.ErrConflict
	}
	return values[0], rows.Err()
}
func (r *Commerce) AddOrderLine(ctx context.Context, tenant, eventID string, line commerce.OrderLine, expectedVersion int64) (commerce.OrderLine, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return line, err
	}
	defer tx.Rollback(ctx)
	var currency string
	var total int64
	err = tx.QueryRow(ctx, `select currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id=$2 and organization_id=$3 and state='draft' and version=$4 for update`, tenant, line.OrderID, line.OrganizationID, expectedVersion).Scan(&currency, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return line, commerce.ErrConflict
	}
	if err != nil {
		return line, err
	}
	var amount int64
	err = tx.QueryRow(ctx, `select e.amount_minor_units from pricing.price_book b join pricing.price_book_entry e on e.tenant_id=b.tenant_id and e.price_book_id=b.price_book_id `+bcPriceTimeContextSQL+` where b.tenant_id=$1 and b.price_book_id=$2 and `+bcPriceCurrencyBindingSQL+` and `+bcPriceEligibilitySQL+` and e.variant_id=$4`, tenant, line.PriceBookID, currency, line.VariantID).Scan(&amount)
	if errors.Is(err, pgx.ErrNoRows) {
		return line, commerce.ErrConflict
	}
	if err != nil {
		return line, err
	}
	var current int64
	if err = tx.QueryRow(ctx, `select coalesce(sum(quantity*unit_price_minor_units),0) from sales.customer_order_line where tenant_id=$1 and order_id=$2`, tenant, line.OrderID).Scan(&current); err != nil {
		return line, err
	}
	lineAmount, amountErr := bcamounts.LineAmount(int64(line.Quantity), amount, 0)
	if amountErr != nil || lineAmount <= 0 || current < 0 || total < 0 || lineAmount > total || current > total-lineAmount {
		return line, commerce.ErrConflict
	}
	line.UnitPriceMinorUnits = amount
	line.OrderVersion = expectedVersion + 1
	_, err = tx.Exec(ctx, `insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units)values($1,$2,$3,$4,$5,$6)`, tenant, line.OrderID, line.LineID, line.VariantID, line.Quantity, amount)
	if err != nil {
		return line, err
	}
	result, err := tx.Exec(ctx, `update sales.customer_order set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3`, tenant, line.OrderID, expectedVersion)
	if err != nil {
		return line, err
	}
	if result.RowsAffected() != 1 {
		return line, commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'order',$3,$4,'order.line-added',1,clock_timestamp(),jsonb_build_object('line_id',$5::text,'variant_id',$6::text))`, tenant, eventID, line.OrderID, line.OrderVersion, line.LineID, line.VariantID)
	if err != nil {
		return line, err
	}
	return line, tx.Commit(ctx)
}
func (r *Commerce) PlaceOrder(ctx context.Context, tenant, organization, orderID string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var total int64
	err = tx.QueryRow(ctx, `select total_minor_units from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 and state='draft' and version=$4 for update`, tenant, organization, orderID, version).Scan(&total)
	if errors.Is(err, pgx.ErrNoRows) {
		return commerce.ErrConflict
	}
	if err != nil {
		return err
	}
	var lines int64
	if err = tx.QueryRow(ctx, `select coalesce(sum(quantity*unit_price_minor_units),0) from sales.customer_order_line where tenant_id=$1 and order_id=$2`, tenant, orderID).Scan(&lines); err != nil {
		return err
	}
	if total != lines || lines <= 0 {
		return commerce.ErrConflict
	}
	result, err := tx.Exec(ctx, `update sales.customer_order set state='placed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3`, tenant, orderID, version)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'order',$3,$4,'order.placed',1,clock_timestamp(),'{}')`, tenant, eventID, orderID, version+1)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) AllocateStock(ctx context.Context, tenant, authorizedOrganization, orderID, lineID, stockID string, orderVersion, stockVersion int64, eventID string) error {
	return r.allocateStock(ctx, tenant, authorizedOrganization, orderID, lineID, stockID, orderVersion, stockVersion, eventID, "")
}
func (r *Commerce) AllocateStockAs(ctx context.Context, tenant, organization, orderID, lineID, stockID string, orderVersion, stockVersion int64, eventID, actor string) error {
	if actor == "" {
		return commerce.ErrConflict
	}
	return r.allocateStock(ctx, tenant, organization, orderID, lineID, stockID, orderVersion, stockVersion, eventID, actor)
}
func (r *Commerce) allocateStock(ctx context.Context, tenant, authorizedOrganization, orderID, lineID, stockID string, orderVersion, stockVersion int64, eventID, actor string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var variant, organization string
	err = tx.QueryRow(ctx, `select l.variant_id,o.organization_id from sales.customer_order_line l join sales.customer_order o on o.tenant_id=l.tenant_id and o.order_id=l.order_id where l.tenant_id=$1 and o.organization_id=$2 and l.order_id=$3 and l.line_id=$4 and l.allocated_stock_unit_id is null and l.quantity=1 and o.state in ('placed','confirmed') and o.version=$5 for update of o,l`, tenant, authorizedOrganization, orderID, lineID, orderVersion).Scan(&variant, &organization)
	if errors.Is(err, pgx.ErrNoRows) {
		return commerce.ErrConflict
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and stock_unit_id=$2 and organization_id=$3 and variant_id=$4 and state='available' and version=$5`, tenant, stockID, organization, variant, stockVersion)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.serial_reservation(tenant_id,reservation_id,organization_id,stock_unit_id,variant_id,demand_kind,demand_id,demand_line_id,status,cancellation_disallowed,version) values($1,$2,$3,$4,$5,'customer-order',$6,$7,'reservation',true,1)`, tenant, eventID, organization, stockID, variant, orderID, lineID)
	if err != nil {
		return err
	}
	result, err = tx.Exec(ctx, `update sales.customer_order_line set allocated_stock_unit_id=$4 where tenant_id=$1 and order_id=$2 and line_id=$3 and allocated_stock_unit_id is null`, tenant, orderID, lineID, stockID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `update sales.customer_order set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3`, tenant, orderID, orderVersion)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'order',$3,$4,'order.stock-allocated',1,clock_timestamp(),jsonb_build_object('line_id',$5::text,'stock_unit_id',$6::text) || case when $7::text<>'' then jsonb_build_object('actor_subject',$7::text) else '{}'::jsonb end)`, tenant, eventID, orderID, orderVersion+1, lineID, stockID, actor)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Commerce) CreatePaymentAttempt(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt) error {
	_, err := r.recordPaymentIntent(ctx, tenant, eventID, idempotency, value, "", false, false)
	return err
}
func (r *Commerce) RecordPaymentIntent(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt, actor string) (commerce.PaymentAttempt, error) {
	return r.recordPaymentIntent(ctx, tenant, eventID, idempotency, value, actor, true, false)
}
func (r *Commerce) RecordOrderPayment(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt, actor string) (commerce.PaymentAttempt, error) {
	return r.recordPaymentIntent(ctx, tenant, eventID, idempotency, value, actor, true, true)
}
func (r *Commerce) recordPaymentIntent(ctx context.Context, tenant, eventID, idempotency string, value commerce.PaymentAttempt, actor string, allowReplay, deriveOrder bool) (commerce.PaymentAttempt, error) {
	empty := commerce.PaymentAttempt{}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return empty, err
	}
	defer tx.Rollback(ctx)
	var currency, orderState string
	var total int64
	// Lock the authorized order before reading its amount/state and writing the intent.
	// Existing receipts remain recoverable after a later order transition.
	err = tx.QueryRow(ctx, `select currency,total_minor_units,state from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 for update`, tenant, value.OrganizationID, value.OrderID).Scan(&currency, &total, &orderState)
	if errors.Is(err, pgx.ErrNoRows) {
		return empty, commerce.ErrConflict
	}
	if err != nil {
		return empty, err
	}
	total, err = orderProviderDue(ctx, tx, tenant, value.OrganizationID, value.OrderID, currency, total)
	if err != nil {
		return empty, err
	}
	if deriveOrder {
		value.Currency = currency
		value.AmountMinorUnits = total
	}
	readReceipt := func() (commerce.PaymentAttempt, error) {
		var saved commerce.PaymentAttempt
		err := tx.QueryRow(ctx, `select p.payment_attempt_id,p.order_id,o.organization_id,p.provider_code,p.state,p.currency,p.amount_minor_units,p.version from payment.payment_attempt p join sales.customer_order o on o.tenant_id=p.tenant_id and o.order_id=p.order_id where p.tenant_id=$1 and p.provider_code=$2 and p.idempotency_key=$3 and o.organization_id=$4`, tenant, value.ProviderCode, idempotency, value.OrganizationID).Scan(&saved.ID, &saved.OrderID, &saved.OrganizationID, &saved.ProviderCode, &saved.State, &saved.Currency, &saved.AmountMinorUnits, &saved.Version)
		if err != nil {
			return empty, err
		}
		if !allowReplay || saved.OrderID != value.OrderID || saved.Currency != value.Currency || saved.AmountMinorUnits != value.AmountMinorUnits {
			return empty, commerce.ErrConflict
		}
		return saved, nil
	}
	if saved, err := readReceipt(); err == nil {
		return saved, tx.Commit(ctx)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return empty, err
	}
	// Every public creation path requests the remaining provider due. Enforce initial
	// intent exclusivity here, not just in the new frontend route.
	var exists bool
	if err := tx.QueryRow(ctx, `select exists(select 1 from payment.payment_attempt where tenant_id=$1 and order_id=$2)`, tenant, value.OrderID).Scan(&exists); err != nil {
		return empty, err
	}
	if exists {
		return empty, commerce.ErrConflict
	}
	if total <= 0 || currency != value.Currency || total != value.AmountMinorUnits || (orderState != "placed" && orderState != "confirmed") {
		return empty, commerce.ErrConflict
	}
	result, err := tx.Exec(ctx, `insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,idempotency_key,state,currency,amount_minor_units,version)values($1,$2,$3,$4,$5,'created',$6,$7,1) on conflict(tenant_id,provider_code,idempotency_key) do nothing`, tenant, value.ID, value.OrderID, value.ProviderCode, idempotency, value.Currency, value.AmountMinorUnits)
	if err != nil {
		return empty, err
	}
	if result.RowsAffected() != 1 {
		// Another order can contend for this provider/key. Read Committed gives
		// this subsequent statement visibility of the committed winning intent.
		saved, err := readReceipt()
		if errors.Is(err, pgx.ErrNoRows) {
			err = commerce.ErrConflict
		}
		if err != nil {
			return empty, err
		}
		return saved, tx.Commit(ctx)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'payment',$3,1,'payment.requested',1,clock_timestamp(),jsonb_build_object('provider_code',$4::text,'order_id',$5::text) || case when $6::text<>'' then jsonb_build_object('actor_subject',$6::text) else '{}'::jsonb end)`, tenant, eventID, value.ID, value.ProviderCode, value.OrderID, actor)
	if err != nil {
		return empty, err
	}
	value.State = "created"
	value.Version = 1
	return value, tx.Commit(ctx)
}
func (r *Commerce) TransitionPayment(ctx context.Context, tenant, organization, id, current, target string, version int64, providerReference, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update payment.payment_attempt p set state=$6,version=p.version+1,provider_reference=coalesce(nullif($7,''),p.provider_reference),updated_at=clock_timestamp() from sales.customer_order o where p.tenant_id=$1 and p.payment_attempt_id=$3 and p.state=$4 and p.version=$5 and o.tenant_id=p.tenant_id and o.order_id=p.order_id and o.organization_id=$2`, tenant, organization, id, current, version, target, providerReference)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'payment',$3,$4,$5,1,clock_timestamp(),'{}')`, tenant, eventID, id, version+1, "payment."+target)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original activation SQL retained.
func (r *Commerce) activatePriceBookTx(ctx context.Context, tx pgx.Tx, tenant, id, eventID string) error {
	var market, currency string
	var from time.Time
	var until *time.Time
	err := tx.QueryRow(ctx, `select market,currency,valid_from,valid_until from pricing.price_book where tenant_id=$1 and price_book_id=$2 and status='draft' for update`, tenant, id).Scan(&market, &currency, &from, &until)
	if errors.Is(err, pgx.ErrNoRows) {
		return commerce.ErrConflict
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+market+":"+currency)
	if err != nil {
		return err
	}
	var overlap bool
	err = tx.QueryRow(ctx, `select exists(select 1 from pricing.price_book where tenant_id=$1 and market=$2 and currency=$3 and status='active' and valid_from<coalesce($5::timestamptz,'infinity'::timestamptz) and coalesce(valid_until,'infinity'::timestamptz)>$4)`, tenant, market, currency, from, until).Scan(&overlap)
	if err != nil {
		return err
	}
	if overlap {
		return commerce.ErrConflict
	}
	_, err = tx.Exec(ctx, `update pricing.price_book set status='active' where tenant_id=$1 and price_book_id=$2 and status='draft'`, tenant, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'price-book',$3,2,'price-book.activated',1,clock_timestamp(),jsonb_build_object('policy_sha256',$4::text))`, tenant, eventID, id, r.policy.SHA256())
	if err != nil {
		return err
	}
	return nil
}

// AUTHORED transaction extraction; original creation SQL/event retained.
func (r *Commerce) createPriceBookTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value commerce.PriceBook) error {
	_, err := tx.Exec(ctx, `insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,valid_until,status)values($1,$2,$3,$4,$5,$6,'draft')`, tenant, value.ID, value.Market, value.Currency, value.ValidFrom, value.ValidUntil)
	if err != nil {
		return err
	}
	for _, entry := range value.Entries {
		if _, err = tx.Exec(ctx, `insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,$2,$3,$4,$5)`, tenant, value.ID, entry.VariantID, entry.AmountMinorUnits, entry.TaxMode); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'price-book',$3,1,'price-book.created',1,clock_timestamp(),jsonb_build_object('market',$4::text,'currency',$5::text))`, tenant, eventID, value.ID, value.Market, value.Currency)
	if err != nil {
		return err
	}
	return nil
}
