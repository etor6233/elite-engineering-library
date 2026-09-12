# Go Return Fiscal Credit Note Worker

## 1. Metadata

```yaml
pack_id: "GO-RETURN-FISCAL-CREDIT-NOTE-WORKER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Convierte un request fiscal de devolución full-refund en una única nota de crédito ARCA asociada, espera autorización real y reconcilia caídas sin duplicar el comprobante."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "ARCA WSFEv1 4.6"]
compatible_with: ["GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.x", "GO-RETURN-EFFECT-EXECUTION-WORKER 0.1.x", "GO-OFFICIAL-RETURN-REFUND-WORKER 0.1.x", "GO-RETURN-ACCOUNTING-REVERSAL-WORKER 0.1.x", "GO-ARCA-FISCAL-ISSUANCE-API 0.6.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND ARCA public specification"
upstream_sources: ["https://www.arca.gob.ar/facturacion/comprobantes/", "https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf", "https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/sales-returns", "https://www.postgresql.org/docs/18/"]
verified_at: "2026-08-31"
```

Los ocho archivos son `AUTHORED`: no se presentan como código copiado de ARCA, Microsoft o PostgreSQL. ARCA gobierna el comprobante asociado; Microsoft documenta la separación del proceso de devolución; PostgreSQL gobierna la transacción, locks e idempotencia.

## 2. Applicability

Use únicamente para un refund total ya recibido y dispuesto, con inventario, devolución de pago y reversión contable exitosos, y exactamente una factura ARCA local autorizada. Rechace refund parcial, exchange, múltiples facturas originales, moneda/importe divergentes o reglas fiscales no aprobadas.

## 3. Architecture contract

El journey crea el request `fiscal`; este owner no inventa uno nuevo. Reclama por `SKIP LOCKED`, exige los tres efectos previos exitosos y compara refund, posting contable y factura original. Copia campos/IVA/tributos de la factura, selecciona sólo `3→1`, `8→6`, `13→11` y entrega únicamente el ID original al owner fiscal. IDs, idempotency y request hash son deterministas. Crear la solicitud no es éxito: retorna retry mientras esté queued/authorizing/reconcile, bloquea rechazo y cierra sólo tras CAE autorizado. El vínculo es inmutable y una caída entre invoice/link/finalización se reconcilia sin segunda nota.

## 4. Exact file manifest

```text
CREATE internal/returnfiscal/processor.go
CREATE internal/returnfiscal/processor_test.go
CREATE internal/platform/postgres/returnfiscal.go
CREATE internal/platform/postgres/returnfiscal_integration_test.go
CREATE cmd/return-fiscal-worker/main.go
CREATE db/migrations/0024_return_fiscal_credit_note.up.sql
CREATE db/migrations/0024_return_fiscal_credit_note.down.sql
CREATE db/tests/0024_return_fiscal_credit_note.test.sql
```

## 5. Materialization blocks

### FILE: `internal/returnfiscal/processor.go`
```yaml
block_id: "GO-RETURN-FISCAL:processor:v1"
operation: CREATE
provenance: AUTHORED
source: "local durable return fiscal orchestration"
license: "LicenseRef-Workspace-Owner"
sha256: "2fbbfa69135259beec948f822a387f02a67ebfdc5e6f28b33e0f4cf83f097f19"
variables: []
secrets_allowed: false
```
````go
package returnfiscal

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

var (
	ErrNoWork   = errors.New("no return fiscal work")
	ErrPending  = errors.New("credit note authorization pending")
	ErrRejected = errors.New("credit note rejected")
	ErrConflict = errors.New("return fiscal source conflict")
)

var workerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Result struct {
	RequestID, OriginalInvoiceID, CreditInvoiceID, Status, ResultSHA256 string
}

type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error)
	RequestOrObserve(context.Context, returneffects.Work, string) (Result, error)
	Finish(context.Context, returneffects.Work, string, returneffects.Completion) error
}

type IDGenerator interface{ New() string }

type Processor struct {
	store        Store
	ids          IDGenerator
	workerID     string
	lease, retry time.Duration
}

func NewProcessor(store Store, ids IDGenerator, workerID string, lease, retry time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retry < time.Second || retry > time.Hour {
		return nil, fmt.Errorf("invalid return fiscal processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retry: retry}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	work, err := p.store.Claim(ctx, "fiscal", p.workerID, p.ids.New(), p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, observeErr := p.store.RequestOrObserve(ctx, *work, p.workerID)
	completion := returneffects.Completion{Outcome: "succeeded", ProviderReference: result.CreditInvoiceID, ResultSHA256: result.ResultSHA256}
	if observeErr != nil {
		completion = returneffects.Completion{Outcome: "retry", ErrorCode: "FISCAL_TRANSIENT_FAILURE", RetryAfter: p.retry}
		if errors.Is(observeErr, ErrPending) {
			completion.ErrorCode = "FISCAL_AUTHORIZATION_PENDING"
		}
		if errors.Is(observeErr, ErrRejected) {
			completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "FISCAL_AUTHORIZATION_REJECTED", ProviderReference: result.CreditInvoiceID}
		}
		if errors.Is(observeErr, ErrConflict) {
			completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "FISCAL_SOURCE_CONFLICT"}
		}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(observeErr, finishErr)
	}
	if observeErr != nil {
		return Result{}, observeErr
	}
	return result, nil
}
````

### FILE: `internal/returnfiscal/processor_test.go`
```yaml
block_id: "GO-RETURN-FISCAL:processor-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local processor regression"
license: "LicenseRef-Workspace-Owner"
sha256: "752ab293891558cfd422b313cc1e0aa4d42d6a7c1d1fae1268c36cd62853cb42"
variables: []
secrets_allowed: false
```
````go
package returnfiscal

import (
	"context"
	"errors"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

type ids struct{}

func (ids) New() string { return "claim" }

type fakeStore struct {
	work       *returneffects.Work
	result     Result
	err        error
	completion returneffects.Completion
}

func (f *fakeStore) Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error) {
	return f.work, nil
}
func (f *fakeStore) RequestOrObserve(context.Context, returneffects.Work, string) (Result, error) {
	return f.result, f.err
}
func (f *fakeStore) Finish(_ context.Context, _ returneffects.Work, _ string, c returneffects.Completion) error {
	f.completion = c
	return nil
}

func TestProcessorWaitsForAuthorizationAndClosesOnlyAuthorized(t *testing.T) {
	work := &returneffects.Work{RequestID: "fiscal", EffectKind: "fiscal", OwnerContext: "fiscal"}
	pending := &fakeStore{work: work, result: Result{CreditInvoiceID: "credit", Status: "queued"}, err: ErrPending}
	p, _ := NewProcessor(pending, ids{}, "worker", time.Minute, time.Second)
	if _, err := p.ProcessOne(context.Background()); !errors.Is(err, ErrPending) || pending.completion.Outcome != "retry" || pending.completion.ErrorCode != "FISCAL_AUTHORIZATION_PENDING" {
		t.Fatalf("pending=%+v err=%v", pending.completion, err)
	}
	success := &fakeStore{work: work, result: Result{CreditInvoiceID: "credit", Status: "authorized", ResultSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	p, _ = NewProcessor(success, ids{}, "worker", time.Minute, time.Second)
	if result, err := p.ProcessOne(context.Background()); err != nil || result.Status != "authorized" || success.completion.Outcome != "succeeded" {
		t.Fatalf("result=%+v completion=%+v err=%v", result, success.completion, err)
	}
}
````

### FILE: `internal/platform/postgres/returnfiscal.go`
```yaml
block_id: "GO-RETURN-FISCAL:postgres:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL adapter governed by ARCA and PostgreSQL contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "3104845953fdc124595be6055bfdf0978879a623a0e47878643b00419b36eb97"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnfiscal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnFiscal struct {
	pool    *pgxpool.Pool
	effects *ReturnEffects
	fiscal  *Fiscal
}

func NewReturnFiscal(pool *pgxpool.Pool) *ReturnFiscal {
	return &ReturnFiscal{pool: pool, effects: NewReturnEffects(pool), fiscal: NewFiscal(pool)}
}
func (r *ReturnFiscal) Claim(ctx context.Context, owner, worker, token string, lease time.Duration) (*returneffects.Work, error) {
	return r.effects.Claim(ctx, owner, worker, token, lease)
}
func (r *ReturnFiscal) Finish(ctx context.Context, work returneffects.Work, worker string, c returneffects.Completion) error {
	return r.effects.Finish(ctx, work, worker, c)
}

func stableReturnFiscalUUID(scope, value string) string {
	sum := sha256.Sum256([]byte(scope + "\x00" + value))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
func stableReturnFiscalHash(scope, value string) string {
	sum := sha256.Sum256([]byte(scope + "\x00" + value))
	return hex.EncodeToString(sum[:])
}

func (r *ReturnFiscal) RequestOrObserve(ctx context.Context, work returneffects.Work, worker string) (returnfiscal.Result, error) {
	if work.EffectKind != "fiscal" || work.OwnerContext != "fiscal" {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	var originalID string
	var decidedAt time.Time
	err := r.pool.QueryRow(ctx, `
select original.invoice_id,d.decided_at
from sales.return_effect_execution x
join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id and d.customer_remedy='refund'
join sales.return_receipt receipt on receipt.tenant_id=d.tenant_id and receipt.receipt_id=d.receipt_id
join sales.return_effect_request inventory on inventory.tenant_id=e.tenant_id and inventory.disposition_id=e.disposition_id and inventory.effect_kind='inventory'
join sales.return_effect_execution inventory_x on inventory_x.tenant_id=inventory.tenant_id and inventory_x.request_id=inventory.request_id and inventory_x.status='succeeded'
join sales.return_effect_request remedy on remedy.tenant_id=e.tenant_id and remedy.disposition_id=e.disposition_id and remedy.effect_kind='refund'
join sales.return_effect_execution remedy_x on remedy_x.tenant_id=remedy.tenant_id and remedy_x.request_id=remedy.request_id and remedy_x.status='succeeded'
join payment.return_refund refund on refund.tenant_id=remedy.tenant_id and refund.request_id=remedy.request_id and refund.state='succeeded'
join sales.return_effect_request accounting on accounting.tenant_id=e.tenant_id and accounting.disposition_id=e.disposition_id and accounting.effect_kind='accounting'
join sales.return_effect_execution accounting_x on accounting_x.tenant_id=accounting.tenant_id and accounting_x.request_id=accounting.request_id and accounting_x.status='succeeded'
join accounting.return_effect_posting posting on posting.tenant_id=accounting.tenant_id and posting.request_id=accounting.request_id and posting.status='posted'
join fiscal.invoice original on original.tenant_id=receipt.tenant_id and original.organization_id=receipt.organization_id and original.order_id=receipt.order_id and original.status='authorized' and original.voucher_type in(1,6,11)
where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
and e.effect_kind='fiscal' and e.owner_context='fiscal' and refund.amount_minor_units=original.total_minor_units and refund.currency=original.currency and posting.reversed_minor_units=original.total_minor_units
and 1=(select count(*) from fiscal.invoice candidate where candidate.tenant_id=original.tenant_id and candidate.organization_id=original.organization_id and candidate.order_id=original.order_id and candidate.status='authorized' and candidate.voucher_type in(1,6,11))`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&originalID, &decidedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	if err != nil {
		return returnfiscal.Result{}, err
	}
	original, err := r.fiscal.GetInvoice(ctx, work.TenantID, work.OrganizationID, originalID)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	creditType := map[int]int{1: 3, 6: 8, 11: 13}[original.VoucherType]
	if creditType == 0 {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	creditID := stableReturnFiscalUUID("credit-invoice", work.TenantID+"/"+work.RequestID)
	credit := original
	credit.TenantID = ""
	credit.ID = creditID
	credit.VoucherType = creditType
	credit.IssuedOn = decidedAt.UTC()
	credit.Status = "queued"
	credit.VoucherNumber = 0
	credit.CAE = ""
	credit.CAEExpiresOn = nil
	credit.Version = 1
	credit.AssociatedVouchers = []fiscal.AssociatedVoucher{{InvoiceID: original.ID}}
	payload, err := json.Marshal(struct {
		RequestID, OriginalInvoiceID, CreditInvoiceID string
		VoucherType                                   int
		IssuedOn                                      string
		Total                                         int64
	}{work.RequestID, original.ID, creditID, creditType, credit.IssuedOn.Format("2006-01-02"), credit.TotalMinorUnits})
	if err != nil {
		return returnfiscal.Result{}, err
	}
	digest := sha256.Sum256(payload)
	requestHash := hex.EncodeToString(digest[:])
	idempotency := stableReturnFiscalHash("credit-idempotency", work.TenantID+"/"+work.RequestID)
	created, _, err := r.fiscal.RequestInvoice(ctx, work.TenantID, idempotency, requestHash, stableReturnFiscalUUID("credit-event", work.TenantID+"/"+work.RequestID), credit)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	inserted, err := r.pool.Exec(ctx, `insert into fiscal.return_credit_note_link(tenant_id,request_id,disposition_id,original_invoice_id,credit_invoice_id,currency,total_minor_units,request_sha256_hex)values($1,$2,$3,$4,$5,$6,$7,$8) on conflict do nothing`, work.TenantID, work.RequestID, work.DispositionID, original.ID, created.ID, created.Currency, created.TotalMinorUnits, requestHash)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	_ = inserted
	var exact int
	err = r.pool.QueryRow(ctx, `select count(*) from fiscal.return_credit_note_link where tenant_id=$1 and request_id=$2 and disposition_id=$3 and original_invoice_id=$4 and credit_invoice_id=$5 and currency=$6 and total_minor_units=$7 and request_sha256_hex=$8`, work.TenantID, work.RequestID, work.DispositionID, original.ID, created.ID, created.Currency, created.TotalMinorUnits, requestHash).Scan(&exact)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	if exact != 1 {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	result := returnfiscal.Result{RequestID: work.RequestID, OriginalInvoiceID: original.ID, CreditInvoiceID: created.ID, Status: created.Status}
	if created.Status == "rejected" {
		return result, returnfiscal.ErrRejected
	}
	if created.Status != "authorized" {
		return result, returnfiscal.ErrPending
	}
	resultBytes, _ := json.Marshal(struct {
		RequestID, CreditInvoiceID, CAE string
		VoucherNumber                   int64
	}{work.RequestID, created.ID, created.CAE, created.VoucherNumber})
	resultDigest := sha256.Sum256(resultBytes)
	result.ResultSHA256 = hex.EncodeToString(resultDigest[:])
	return result, nil
}
````

### FILE: `internal/platform/postgres/returnfiscal_integration_test.go`
```yaml
block_id: "GO-RETURN-FISCAL:postgres-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL integration regression"
license: "LicenseRef-Workspace-Owner"
sha256: "ffb229094889308726915e2542510f548e6753abde1179fea27782697e43ec23"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnfiscal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReturnFiscalRequestsOnceAndWaitsForAuthorization(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e146"
	cleanup := func() {
		tx, e := pool.Begin(ctx)
		if e != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"fiscal.return_credit_note_link", "fiscal.invoice_associated_voucher", "fiscal.issuance_attempt", "fiscal.invoice_vat", "fiscal.invoice_other_tax", "fiscal.invoice", "fiscal.point_of_sale", "accounting.return_effect_posting", "payment.return_refund_observation", "payment.return_refund", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "payment.payment_attempt", "sales.customer_order", "platform.outbox_event", "org.organization", "platform.tenant"} {
			_, _ = tx.Exec(ctx, `delete from `+table+` where tenant_id=$1`, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	defer cleanup()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `set local session_replication_role=replica`); err != nil {
		t.Fatal(err)
	}
	queries := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'return-fiscal','Return Fiscal','Return Fiscal')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','return-fiscal-store','Store','store')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',12100,4)`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','stripe','pay-ref','return-fiscal-payment','captured','ARS',12100,2)`,
		`insert into fiscal.point_of_sale(tenant_id,point_of_sale_id,organization_id,taxpayer_cuit,environment,point_of_sale_number,active,version)values($1,'pos','store','30715117564','homologation',1,true,1)`,
		`insert into fiscal.invoice(tenant_id,invoice_id,organization_id,order_id,payment_attempt_id,point_of_sale_id,voucher_type,concept,recipient_document_type,recipient_document,recipient_vat_condition_id,currency,provider_currency,total_minor_units,net_minor_units,vat_minor_units,exempt_minor_units,non_taxed_minor_units,other_tax_minor_units,issued_on,status,voucher_number,cae,cae_expires_on,authorized_at,request_hash,idempotency_key,version)values($1,'original','store','order','payment','pos',6,1,99,'0',5,'ARS','PES',12100,10000,2100,0,0,0,'2026-08-15','authorized',1,'12345678901234','2026-08-31',clock_timestamp(),repeat('a',64),'original-fiscal-0001',2)`,
		`insert into fiscal.invoice_vat(tenant_id,invoice_id,vat_id,base_minor_units,amount_minor_units)values($1,'original',5,10000,2100)`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values($1,'receipt','auth','store','order','stock','customer','SERIAL','damaged','Received',repeat('b',64),'operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject,decided_at)values($1,'disposition','receipt','quarantine','refund','Refund','operator','2026-08-31T12:00:00Z')`,
		`insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,'inventory','disposition','inventory','inventory','requested','inventory-effect-0001'),($1,'refund','disposition','refund','payment','requested','refund-effect-000001'),($1,'accounting','disposition','accounting','accounting','requested','accounting-effect-01'),($1,'fiscal','disposition','fiscal','fiscal','requested','fiscal-effect-000001')`,
		`insert into sales.return_effect_execution(tenant_id,request_id,status,result_sha256_hex)values($1,'inventory','succeeded',repeat('c',64)),($1,'refund','succeeded',repeat('d',64)),($1,'accounting','succeeded',repeat('e',64)),($1,'fiscal','requested',null)`,
		`insert into payment.return_refund(tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,provider_refund_reference,provider_status,state,response_sha256_hex)values($1,'refund','payment','order','line','stock','stripe','pay-ref','refund-provider-0001','ARS',12100,'refund-ref','succeeded','succeeded',repeat('f',64))`,
		`insert into accounting.return_effect_posting(tenant_id,request_id,disposition_id,source_order_id,remedy,original_journal_id,reversal_journal_id,currency,reversed_minor_units,status,result_sha256_hex)values($1,'accounting','disposition','order','refund','sale-journal','reversal-journal','ARS',12100,'posted',repeat('1',64))`,
	}
	for _, query := range queries {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	store := NewReturnFiscal(pool)
	work, err := store.Claim(ctx, "fiscal", "worker", "claim-1", time.Minute)
	if err != nil || work == nil {
		t.Fatalf("work=%+v err=%v", work, err)
	}
	result, err := store.RequestOrObserve(ctx, *work, "worker")
	if !errors.Is(err, returnfiscal.ErrPending) || result.CreditInvoiceID == "" || result.Status != "queued" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if err = store.Finish(ctx, *work, "worker", returneffects.Completion{Outcome: "retry", ErrorCode: "FISCAL_AUTHORIZATION_PENDING", ProviderReference: result.CreditInvoiceID, RetryAfter: time.Second}); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update fiscal.invoice set status='authorized',voucher_number=2,cae='22345678901234',cae_expires_on='2026-09-10',authorized_at=clock_timestamp(),version=version+1 where tenant_id=$1 and invoice_id=$2`, tenant, result.CreditInvoiceID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update sales.return_effect_execution set available_at=clock_timestamp() where tenant_id=$1 and request_id='fiscal'`, tenant); err != nil {
		t.Fatal(err)
	}
	work, err = store.Claim(ctx, "fiscal", "worker", "claim-2", time.Minute)
	if err != nil || work == nil {
		t.Fatalf("work2=%+v err=%v", work, err)
	}
	result2, err := store.RequestOrObserve(ctx, *work, "worker")
	if err != nil || result2.Status != "authorized" || result2.CreditInvoiceID != result.CreditInvoiceID {
		t.Fatalf("result2=%+v err=%v", result2, err)
	}
	if err = store.Finish(ctx, *work, "worker", returneffects.Completion{Outcome: "succeeded", ProviderReference: result2.CreditInvoiceID, ResultSHA256: result2.ResultSHA256}); err != nil {
		t.Fatal(err)
	}
	var links, credits, attempts int
	var state string
	if err = pool.QueryRow(ctx, `select (select count(*) from fiscal.return_credit_note_link where tenant_id=$1),(select count(*) from fiscal.invoice where tenant_id=$1 and voucher_type=8),(select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id='fiscal'),(select status from sales.return_effect_execution where tenant_id=$1 and request_id='fiscal')`, tenant).Scan(&links, &credits, &attempts, &state); err != nil {
		t.Fatal(err)
	}
	if links != 1 || credits != 1 || attempts != 2 || state != "succeeded" {
		t.Fatalf("links=%d credits=%d attempts=%d state=%s", links, credits, attempts, state)
	}
	if _, err = pool.Exec(ctx, `update fiscal.return_credit_note_link set total_minor_units=1 where tenant_id=$1`, tenant); err == nil {
		t.Fatal("immutable return credit link update accepted")
	}
}
````

### FILE: `cmd/return-fiscal-worker/main.go`
```yaml
block_id: "GO-RETURN-FISCAL:command:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded worker entrypoint"
license: "LicenseRef-Workspace-Owner"
sha256: "41c63626c9c94e23162cb20bb7da651d9cf1462b2943b662fc154c60bac9cc39"
variables: []
secrets_allowed: true
```
````go
package main

import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/returnfiscal"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	url, id := os.Getenv("DATABASE_URL"), os.Getenv("RETURN_FISCAL_WORKER_ID")
	if url == "" || id == "" {
		os.Exit(2)
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil {
		os.Exit(2)
	}
	cfg.MaxConns = 4
	p, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		os.Exit(1)
	}
	defer p.Close()
	if e = p.Ping(ctx); e != nil {
		os.Exit(1)
	}
	processor, e := returnfiscal.NewProcessor(postgres.NewReturnFiscal(p), randomid.Generator{}, id, 2*time.Minute, 30*time.Second)
	if e != nil {
		os.Exit(2)
	}
	if e = run(ctx, processor); e != nil && !errors.Is(e, context.Canceled) {
		os.Exit(1)
	}
}
func run(ctx context.Context, p *returnfiscal.Processor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		_, e := p.ProcessOne(ctx)
		if e == nil {
			continue
		}
		if errors.Is(e, context.Canceled) {
			return e
		}
		wait := 2 * time.Second
		if errors.Is(e, returnfiscal.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else if !errors.Is(e, returnfiscal.ErrPending) {
			slog.Warn("return fiscal deferred")
		}
		timer.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
````

### FILE: `db/migrations/0024_return_fiscal_credit_note.up.sql`
```yaml
block_id: "GO-RETURN-FISCAL:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL immutable link"
license: "LicenseRef-Workspace-Owner"
sha256: "7f3641342abd4afe345251c40449e6b1db0d8a67aeee5056c530c9c63821b7a7"
variables: []
secrets_allowed: false
```
````sql
begin;
create table fiscal.return_credit_note_link (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  original_invoice_id text not null,
  credit_invoice_id text not null,
  currency text not null check(currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check(total_minor_units>0),
  request_sha256_hex text not null check(request_sha256_hex ~ '^[0-9a-f]{64}$'),
  created_at timestamptz not null default clock_timestamp(),
  primary key(tenant_id,request_id),
  unique(tenant_id,disposition_id),
  unique(tenant_id,credit_invoice_id),
  foreign key(tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key(tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  foreign key(tenant_id,original_invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  foreign key(tenant_id,credit_invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  check(original_invoice_id<>credit_invoice_id)
);
create or replace function fiscal.prevent_return_credit_link_mutation() returns trigger language plpgsql as $$ begin raise exception using errcode='23514',message='return credit note link is immutable'; end $$;
create trigger return_credit_note_link_immutable before update or delete on fiscal.return_credit_note_link for each row execute function fiscal.prevent_return_credit_link_mutation();
commit;
````

### FILE: `db/migrations/0024_return_fiscal_credit_note.down.sql`
```yaml
block_id: "GO-RETURN-FISCAL:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "ca4bb6c941714d3faba162c0a489aa6b91dc6c1bed1a545b21c1965a20b598c7"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table if exists fiscal.return_credit_note_link;
drop function if exists fiscal.prevent_return_credit_link_mutation();
commit;
````

### FILE: `db/tests/0024_return_fiscal_credit_note.test.sql`
```yaml
block_id: "GO-RETURN-FISCAL:migration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL contract regression"
license: "LicenseRef-Workspace-Owner"
sha256: "1c59d76b08675ac3f9919c7e5f773653025c48c087d1225a02e80a8dfe9112a2"
variables: []
secrets_allowed: false
```
````sql
begin;
do $$ begin
  if to_regclass('fiscal.return_credit_note_link') is null then raise exception 'return credit note link missing'; end if;
  if (select count(*) from information_schema.table_constraints where table_schema='fiscal' and table_name='return_credit_note_link' and constraint_type='FOREIGN KEY')<>4 then raise exception 'return credit note foreign keys missing'; end if;
  if not exists(select 1 from pg_trigger where tgrelid='fiscal.return_credit_note_link'::regclass and tgname='return_credit_note_link_immutable' and not tgisinternal) then raise exception 'return credit note immutability missing'; end if;
end $$;
rollback;
````

## 6. Configuration surface

| Variable | Type/default | Secret | Validation/effect |
|---|---|---|---|
| `DATABASE_URL` | PostgreSQL URL / none | yes | TLS/secret injection belongs to target |
| `RETURN_FISCAL_WORKER_ID` | bounded worker ID / none | no | required, 1–64 safe characters |
| lease/retry | 2m/30s in entrypoint | no | bounded by processor contract |

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | 1.26.7 | processor/worker | BSD-3-Clause | build/runtime | go.dev |
| pgx | 5.10.0 | PostgreSQL transactions | MIT | build/runtime | github.com/jackc/pgx |
| PostgreSQL | 18.6 | durable execution/link | PostgreSQL | runtime/test | postgresql.org |
| ARCA WSFEv1 | manual 4.6 | credit-note association authority | public specification | contract | arca.gob.ar |

## 8. Apply order

Materialize after fiscal 0.6.0 and return owners; apply 0024 after 0023; run full Go/SQL/PostgreSQL gates; start fiscal issuance and this observer in homologation only. Roll back by draining both workers and removing 0024 when retained history permits.

## 9. Verification

Require 8/8 exact blocks, gofmt, full Go test/vet/build, PostgreSQL 18.6 migrations 0001–0024, SQL gate and integration proving request→pending→authorized→succeeded, one credit invoice, two attempts, immutable link and no duplicate on replay. Live ARCA homologation, partial notes, debit notes, rendering and accountant/business acceptance remain project gates.

## 10. Reconstruction evidence

V146 evidence is recorded in `reconstruction_evidence/RETURN_FISCAL_CREDIT_NOTE_EXECUTION_INVENTORY_2026-08-31_V146.md`. It does not claim production authorization or homologation.
