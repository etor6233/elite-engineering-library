package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Fiscal struct{ pool *pgxpool.Pool }

func NewFiscal(pool *pgxpool.Pool) *Fiscal { return &Fiscal{pool: pool} }

func fiscalConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return fiscal.ErrConflict
	}
	return err
}

func (r *Fiscal) ConfigurePointOfSale(ctx context.Context, tenant, eventID string, value fiscal.PointOfSale) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into fiscal.point_of_sale(tenant_id,point_of_sale_id,organization_id,taxpayer_cuit,environment,point_of_sale_number,active,version) select $1,$2,$3,$4,$5,$6,true,1 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active')`, tenant, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Environment, value.Number)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-point-of-sale',$3,1,'fiscal-point-of-sale.configured',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'environment',$5::text,'point_of_sale_number',$6::integer))`, tenant, eventID, value.ID, value.OrganizationID, value.Environment, value.Number)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) RequestInvoice(ctx context.Context, tenant, idempotency, requestHash, eventID string, value fiscal.Invoice) (fiscal.Invoice, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, false, err
	}
	defer tx.Rollback(ctx)
	var existingHash string
	err = tx.QueryRow(ctx, `select request_hash from fiscal.invoice where tenant_id=$1 and idempotency_key=$2`, tenant, idempotency).Scan(&existingHash)
	if err == nil {
		if existingHash != requestHash {
			return value, false, fiscal.ErrConflict
		}
		existing, loadErr := scanFiscalInvoice(tx.QueryRow(ctx, invoiceSelect+` where i.tenant_id=$1 and i.idempotency_key=$2`, tenant, idempotency))
		if loadErr == nil {
			loadErr = loadFiscalLines(ctx, tx, &existing)
		}
		return existing, true, loadErr
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return value, false, err
	}
	originalType, credit := fiscal.AssociatedOriginalVoucherType(value.VoucherType)
	if (!credit && (!fiscal.SupportedInvoiceVoucherType(value.VoucherType) || len(value.AssociatedVouchers) != 0)) || (credit && (len(value.AssociatedVouchers) != 1 || value.AssociatedVouchers[0].InvoiceID == "")) {
		return value, false, fiscal.ErrConflict
	}
	if credit {
		associated := &value.AssociatedVouchers[0]
		err = tx.QueryRow(ctx, `select original_pos.taxpayer_cuit,original.voucher_type,original_pos.point_of_sale_number,original.voucher_number,original.issued_on from fiscal.invoice original join fiscal.point_of_sale original_pos on original_pos.tenant_id=original.tenant_id and original_pos.point_of_sale_id=original.point_of_sale_id join fiscal.point_of_sale credit_pos on credit_pos.tenant_id=original.tenant_id and credit_pos.point_of_sale_id=$5 and credit_pos.organization_id=original.organization_id and credit_pos.taxpayer_cuit=original_pos.taxpayer_cuit and credit_pos.active where original.tenant_id=$1 and original.organization_id=$2 and original.invoice_id=$3 and original.status='authorized' and original.voucher_type=$4 and original.voucher_number is not null for key share of original,original_pos,credit_pos`, tenant, value.OrganizationID, associated.InvoiceID, originalType, value.PointOfSaleID).Scan(&associated.TaxpayerCUIT, &associated.VoucherType, &associated.PointOfSale, &associated.Number, &associated.IssuedOn)
		if errors.Is(err, pgx.ErrNoRows) {
			return value, false, fiscal.ErrConflict
		}
		if err != nil {
			return value, false, err
		}
		if associated.IssuedOn.After(value.IssuedOn) {
			return value, false, fiscal.ErrConflict
		}
	}
	var orderCurrency, paymentCurrency string
	var orderTotal, paymentTotal int64
	err = tx.QueryRow(ctx, `select o.currency,o.total_minor_units,p.currency,p.amount_minor_units from sales.customer_order o join payment.payment_attempt p on p.tenant_id=o.tenant_id and p.order_id=o.order_id where o.tenant_id=$1 and o.order_id=$2 and o.organization_id=$3 and o.state<>'cancelled' and p.payment_attempt_id=$4 and p.state='captured' for update of o,p`, tenant, value.OrderID, value.OrganizationID, value.PaymentAttemptID).Scan(&orderCurrency, &orderTotal, &paymentCurrency, &paymentTotal)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, false, fiscal.ErrConflict
	}
	if err != nil {
		return value, false, err
	}
	total := value.NetMinorUnits + value.VATMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits
	if orderCurrency != paymentCurrency || orderCurrency != "ARS" || orderTotal != paymentTotal || paymentTotal != total || !fiscalDetailsValid(value) {
		return value, false, fiscal.ErrConflict
	}
	result, err := tx.Exec(ctx, `insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,service_from,service_until,payment_due_on,status,request_hash,idempotency_key,version) select $1,$2,$3,$4,$5,p.point_of_sale_id,$6,$7,$8,$9,$10,$11,'PES',$12,$13,$14,$15,$16,$17,($18::timestamptz at time zone 'UTC')::date,($19::timestamptz at time zone 'UTC')::date,($20::timestamptz at time zone 'UTC')::date,($21::timestamptz at time zone 'UTC')::date,'queued',$22,$23,1 from fiscal.point_of_sale p where p.tenant_id=$1 and p.point_of_sale_id=$24 and p.organization_id=$3 and p.active`, tenant, value.ID, value.OrganizationID, value.OrderID, value.PaymentAttemptID, value.VoucherType, value.Concept, value.RecipientDocumentType, value.RecipientDocument, value.RecipientVATConditionID, orderCurrency, total, value.NetMinorUnits, value.VATMinorUnits, value.ExemptMinorUnits, value.NonTaxedMinorUnits, value.OtherTaxMinorUnits, value.IssuedOn, value.ServiceFrom, value.ServiceUntil, value.PaymentDueOn, requestHash, idempotency, value.PointOfSaleID)
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return value, false, fiscal.ErrConflict
	}
	if credit {
		associated := value.AssociatedVouchers[0]
		if _, err = tx.Exec(ctx, `insert into fiscal.invoice_associated_voucher(tenant_id,invoice_id,associated_invoice_id,taxpayer_cuit,voucher_type,point_of_sale_number,voucher_number,issued_on) values($1,$2,$3,$4,$5,$6,$7,($8::timestamptz at time zone 'UTC')::date)`, tenant, value.ID, associated.InvoiceID, associated.TaxpayerCUIT, associated.VoucherType, associated.PointOfSale, associated.Number, associated.IssuedOn); err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	for _, line := range value.VATLines {
		if _, err = tx.Exec(ctx, `insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units) values($1,$2,$3,$4,$5)`, tenant, value.ID, line.ID, line.BaseMinorUnits, line.AmountMinorUnits); err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	for _, line := range value.OtherTaxLines {
		if _, err = tx.Exec(ctx, `insert into fiscal.invoice_other_tax(tenant_id,invoice_id,tax_id,description,base_minor_units,rate_basis_points,amount_minor_units) values($1,$2,$3,$4,$5,$6,$7)`, tenant, value.ID, line.ID, line.Description, line.BaseMinorUnits, line.RateBasisPoints, line.AmountMinorUnits); err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-invoice',$3,1,'fiscal-invoice.requested',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'order_id',$5::text,'payment_attempt_id',$6::text,'point_of_sale_id',$7::text,'voucher_type',$8::integer,'total_minor_units',$9::bigint,'currency',$10::text))`, tenant, eventID, value.ID, value.OrganizationID, value.OrderID, value.PaymentAttemptID, value.PointOfSaleID, value.VoucherType, total, orderCurrency)
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, false, fiscalConflict(err)
	}
	value.TenantID = tenant
	value.TotalMinorUnits = total
	value.Currency = orderCurrency
	value.Status = "queued"
	value.Version = 1
	return value, false, nil
}

const invoiceSelect = `select i.tenant_id,i.invoice_id,i.organization_id,i.order_id,i.payment_attempt_id,i.point_of_sale_id,p.taxpayer_cuit,p.environment,p.point_of_sale_number,i.voucher_type,i.concept,i.recipient_document_type,i.recipient_document,i.recipient_vat_condition_id,i.currency,i.total_minor_units,i.net_minor_units,i.vat_minor_units,i.exempt_minor_units,i.non_taxed_minor_units,i.other_tax_minor_units,i.issued_on,i.service_from,i.service_until,i.payment_due_on,i.status,coalesce(i.voucher_number,0),coalesce(i.cae,''),i.cae_expires_on,i.version from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id`

type fiscalRow interface{ Scan(...any) error }

func scanFiscalInvoice(row fiscalRow) (fiscal.Invoice, error) {
	var v fiscal.Invoice
	err := row.Scan(&v.TenantID, &v.ID, &v.OrganizationID, &v.OrderID, &v.PaymentAttemptID, &v.PointOfSaleID, &v.TaxpayerCUIT, &v.Environment, &v.PointOfSaleNumber, &v.VoucherType, &v.Concept, &v.RecipientDocumentType, &v.RecipientDocument, &v.RecipientVATConditionID, &v.Currency, &v.TotalMinorUnits, &v.NetMinorUnits, &v.VATMinorUnits, &v.ExemptMinorUnits, &v.NonTaxedMinorUnits, &v.OtherTaxMinorUnits, &v.IssuedOn, &v.ServiceFrom, &v.ServiceUntil, &v.PaymentDueOn, &v.Status, &v.VoucherNumber, &v.CAE, &v.CAEExpiresOn, &v.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, fiscal.ErrNotFound
	}
	return v, err
}
func (r *Fiscal) GetInvoice(ctx context.Context, tenant, organization, id string) (fiscal.Invoice, error) {
	v, err := scanFiscalInvoice(r.pool.QueryRow(ctx, invoiceSelect+` where i.tenant_id=$1 and i.organization_id=$2 and i.invoice_id=$3`, tenant, organization, id))
	if err != nil {
		return v, err
	}
	err = loadFiscalLines(ctx, r.pool, &v)
	return v, err
}

func (r *Fiscal) ClaimInvoice(ctx context.Context, worker string, lease time.Duration, attemptID string) (fiscal.Invoice, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fiscal.Invoice{}, err
	}
	defer tx.Rollback(ctx)
	query := invoiceSelect + ` where i.status in ('queued','reconcile_required') and p.active and (p.active_invoice_id is null or p.lease_until < clock_timestamp()) order by i.created_at,i.tenant_id,i.invoice_id for update of i,p skip locked limit 1`
	v, err := scanFiscalInvoice(tx.QueryRow(ctx, query))
	if errors.Is(err, fiscal.ErrNotFound) {
		return v, fiscal.ErrNoWork
	}
	if err != nil {
		return v, err
	}
	if err = loadFiscalLines(ctx, tx, &v); err != nil {
		return v, err
	}
	v.ClaimedFrom = v.Status
	result, err := tx.Exec(ctx, `update fiscal.point_of_sale set active_invoice_id=$3,lease_owner=$4,lease_until=clock_timestamp()+$5::interval where tenant_id=$1 and point_of_sale_id=$2 and (active_invoice_id is null or lease_until<clock_timestamp())`, v.TenantID, v.PointOfSaleID, v.ID, worker, lease.String())
	if err != nil {
		return v, err
	}
	if result.RowsAffected() != 1 {
		return v, fiscal.ErrConflict
	}
	err = tx.QueryRow(ctx, `update fiscal.invoice set attempt_count=attempt_count+1,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2 and point_of_sale_id=$3 returning version`, v.TenantID, v.ID, v.PointOfSaleID).Scan(&v.Version)
	if err != nil {
		return v, err
	}
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,outcome) select tenant_id,$2,invoice_id,'claim',voucher_number,request_hash,'started' from fiscal.invoice where tenant_id=$1 and invoice_id=$3`, v.TenantID, attemptID, v.ID)
	if err != nil {
		return v, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return v, fiscalConflict(err)
	}
	return v, nil
}

func (r *Fiscal) AssignVoucherNumber(ctx context.Context, value fiscal.Invoice, worker string, number int64, responseHash, eventID string) (fiscal.Invoice, error) {
	if number < 1 {
		return value, fiscal.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var tenant, requestHash string
	err = tx.QueryRow(ctx, `select i.tenant_id,i.request_hash from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id where i.tenant_id=$5 and i.invoice_id=$1 and i.point_of_sale_id=$2 and i.status='queued' and i.voucher_number is null and i.version=$3 and p.active_invoice_id=i.invoice_id and p.lease_owner=$4 and p.lease_until>=clock_timestamp() for update of i,p`, value.ID, value.PointOfSaleID, value.Version, worker, value.TenantID).Scan(&tenant, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+value.PointOfSaleID+":"+string(rune(value.VoucherType))); err != nil {
		return value, err
	}
	err = tx.QueryRow(ctx, `update fiscal.invoice set voucher_number=$3,status='authorizing',last_response_hash=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2 returning version,status`, tenant, value.ID, number, responseHash).Scan(&value.Version, &value.Status)
	if err != nil {
		return value, fiscalConflict(err)
	}
	value.VoucherNumber = number
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,response_hash,outcome) values($1,$2,$3,'assign-number',$4,$5,$6,'observed')`, tenant, eventID, value.ID, number, requestHash, responseHash)
	if err != nil {
		return value, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, fiscalConflict(err)
	}
	return value, nil
}

func (r *Fiscal) FinishInvoice(ctx context.Context, value fiscal.Invoice, worker string, authorization fiscal.Authorization, eventID string) (fiscal.Invoice, error) {
	codes, err := json.Marshal(append([]string{}, authorization.ProviderCodes...))
	if err != nil {
		return value, fiscal.ErrInvalid
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	status := "rejected"
	outcome := "rejected"
	if authorization.Authorized {
		status = "authorized"
		outcome = "authorized"
	}
	var tenant, requestHash string
	err = tx.QueryRow(ctx, `select i.tenant_id,i.request_hash from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id where i.tenant_id=$5 and i.invoice_id=$1 and i.point_of_sale_id=$2 and i.status in ('authorizing','reconcile_required') and i.voucher_number=$3 and p.active_invoice_id=i.invoice_id and p.lease_owner=$4 and p.lease_until>=clock_timestamp() for update of i,p`, value.ID, value.PointOfSaleID, value.VoucherNumber, worker, value.TenantID).Scan(&tenant, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrConflict
	}
	if err != nil {
		return value, err
	}
	err = tx.QueryRow(ctx, `update fiscal.invoice set status=$3,cae=$4,cae_expires_on=($5::timestamptz at time zone 'UTC')::date,last_response_hash=$6,provider_codes=$7,authorized_at=case when $3='authorized' then clock_timestamp() else null end,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2 returning version,status,coalesce(cae,''),cae_expires_on`, tenant, value.ID, status, nullIfEmpty(authorization.CAE), authorization.CAEExpiresOn, authorization.ResponseHash, codes).Scan(&value.Version, &value.Status, &value.CAE, &value.CAEExpiresOn)
	if err != nil {
		return value, fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `update fiscal.point_of_sale set active_invoice_id=null,lease_owner=null,lease_until=null where tenant_id=$1 and point_of_sale_id=$2 and active_invoice_id=$3 and lease_owner=$4`, tenant, value.PointOfSaleID, value.ID, worker)
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,response_hash,outcome,provider_codes) values($1,$2,$3,'finish',$4,$5,$6,$7,$8)`, tenant, eventID, value.ID, value.VoucherNumber, requestHash, authorization.ResponseHash, outcome, codes)
	if err != nil {
		return value, fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'fiscal-invoice',$2,$3,$4,1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'order_id',$6::text,'voucher_number',$7::bigint,'voucher_type',$8::integer,'cae',$9::text))`, tenant, value.ID, value.Version, "fiscal-invoice."+status, value.OrganizationID, value.OrderID, value.VoucherNumber, value.VoucherType, value.CAE)
	if err != nil {
		return value, fiscalConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, fiscalConflict(err)
	}
	return value, nil
}

func (r *Fiscal) DeferInvoice(ctx context.Context, value fiscal.Invoice, worker, phase, responseHash string, ambiguous bool) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	status := "queued"
	if ambiguous || value.VoucherNumber > 0 {
		status = "reconcile_required"
	}
	var tenant, requestHash string
	err = tx.QueryRow(ctx, `select i.tenant_id,i.request_hash from fiscal.invoice i join fiscal.point_of_sale p on p.tenant_id=i.tenant_id and p.point_of_sale_id=i.point_of_sale_id where i.tenant_id=$3 and i.invoice_id=$1 and p.active_invoice_id=i.invoice_id and p.lease_owner=$2 for update of i,p`, value.ID, worker, value.TenantID).Scan(&tenant, &requestHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return fiscal.ErrConflict
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `update fiscal.invoice set status=$3,last_response_hash=case when $4='' then last_response_hash else $4 end,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and invoice_id=$2`, tenant, value.ID, status, responseHash)
	if err != nil {
		return fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `update fiscal.point_of_sale set active_invoice_id=null,lease_owner=null,lease_until=null where tenant_id=$1 and point_of_sale_id=$2 and active_invoice_id=$3 and lease_owner=$4`, tenant, value.PointOfSaleID, value.ID, worker)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into fiscal.issuance_attempt(tenant_id,attempt_id,invoice_id,phase,voucher_number,request_hash,response_hash,outcome,error_class) values($1,gen_random_uuid(),$2,'defer',$3,$4,nullif($5,''),'deferred',$6)`, tenant, value.ID, nullIfZero(value.VoucherNumber), requestHash, responseHash, phase)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
func nullIfZero(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

type fiscalQueryer interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func loadFiscalLines(ctx context.Context, q fiscalQueryer, value *fiscal.Invoice) error {
	rows, err := q.Query(ctx, `select vat_id,base_minor_units,amount_minor_units from fiscal.invoice_vat where tenant_id=$1 and invoice_id=$2 order by vat_id`, value.TenantID, value.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	value.VATLines = []fiscal.VATLine{}
	for rows.Next() {
		var line fiscal.VATLine
		if err = rows.Scan(&line.ID, &line.BaseMinorUnits, &line.AmountMinorUnits); err != nil {
			return err
		}
		value.VATLines = append(value.VATLines, line)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	rows.Close()
	rows, err = q.Query(ctx, `select tax_id,description,base_minor_units,rate_basis_points,amount_minor_units from fiscal.invoice_other_tax where tenant_id=$1 and invoice_id=$2 order by tax_id`, value.TenantID, value.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	value.OtherTaxLines = []fiscal.OtherTaxLine{}
	for rows.Next() {
		var line fiscal.OtherTaxLine
		if err = rows.Scan(&line.ID, &line.Description, &line.BaseMinorUnits, &line.RateBasisPoints, &line.AmountMinorUnits); err != nil {
			return err
		}
		value.OtherTaxLines = append(value.OtherTaxLines, line)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	rows.Close()
	rows, err = q.Query(ctx, `select associated_invoice_id,taxpayer_cuit,voucher_type,point_of_sale_number,voucher_number,issued_on from fiscal.invoice_associated_voucher where tenant_id=$1 and invoice_id=$2`, value.TenantID, value.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	value.AssociatedVouchers = []fiscal.AssociatedVoucher{}
	for rows.Next() {
		var associated fiscal.AssociatedVoucher
		if err = rows.Scan(&associated.InvoiceID, &associated.TaxpayerCUIT, &associated.VoucherType, &associated.PointOfSale, &associated.Number, &associated.IssuedOn); err != nil {
			return err
		}
		value.AssociatedVouchers = append(value.AssociatedVouchers, associated)
	}
	return rows.Err()
}
func fiscalDetailsValid(value fiscal.Invoice) bool {
	var base, vat, other int64
	seenVAT := map[int]bool{}
	for _, line := range value.VATLines {
		if line.ID < 1 || line.ID > 999 || line.BaseMinorUnits <= 0 || line.AmountMinorUnits < 0 || seenVAT[line.ID] {
			return false
		}
		seenVAT[line.ID] = true
		base += line.BaseMinorUnits
		vat += line.AmountMinorUnits
	}
	seenTax := map[int]bool{}
	for _, line := range value.OtherTaxLines {
		if line.ID < 1 || line.ID > 999 || len(line.Description) < 1 || len(line.Description) > 80 || line.BaseMinorUnits < 0 || line.RateBasisPoints < 0 || line.RateBasisPoints > 100000 || line.AmountMinorUnits < 0 || seenTax[line.ID] {
			return false
		}
		seenTax[line.ID] = true
		other += line.AmountMinorUnits
	}
	return vat == value.VATMinorUnits && other == value.OtherTaxMinorUnits && ((vat == 0 && len(value.VATLines) == 0) || (vat > 0 && base == value.NetMinorUnits)) && ((other == 0 && len(value.OtherTaxLines) == 0) || (other > 0 && len(value.OtherTaxLines) > 0))
}
