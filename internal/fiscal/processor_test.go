package fiscal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type workRepo struct {
	invoice  Invoice
	assigned int64
	finished int
	deferred string
}

func (r *workRepo) ClaimInvoice(context.Context, string, time.Duration, string) (Invoice, error) {
	return r.invoice, nil
}
func (r *workRepo) AssignVoucherNumber(_ context.Context, value Invoice, _ string, number int64, _, _ string) (Invoice, error) {
	r.assigned = number
	value.VoucherNumber = number
	value.Status = "authorizing"
	return value, nil
}
func (r *workRepo) FinishInvoice(_ context.Context, value Invoice, _ string, result Authorization, _ string) (Invoice, error) {
	r.finished++
	if result.Authorized {
		value.Status = "authorized"
		value.CAE = result.CAE
	} else {
		value.Status = "rejected"
	}
	return value, nil
}
func (r *workRepo) DeferInvoice(_ context.Context, _ Invoice, _ string, phase, _ string, _ bool) error {
	r.deferred = phase
	return nil
}

type providerFake struct {
	last         int64
	consult      Authorization
	authorize    Authorization
	consultErr   error
	authorizeErr error
}

func (p providerFake) LastAuthorized(context.Context, Invoice) (int64, string, error) {
	return p.last, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil
}
func (p providerFake) Consult(context.Context, Invoice) (Authorization, error) {
	return p.consult, p.consultErr
}
func (p providerFake) Authorize(context.Context, Invoice) (Authorization, error) {
	return p.authorize, p.authorizeErr
}

func TestProcessorSequencesAuthorizesAndReconcilesBeforeRetry(t *testing.T) {
	expires := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	ok := Authorization{Found: true, Authorized: true, CAE: "12345678901234", CAEExpiresOn: &expires, ResponseHash: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
	repo := &workRepo{invoice: Invoice{ID: "invoice", Status: "queued"}}
	processor, err := NewProcessor(repo, providerFake{last: 40, authorize: ok}, &testIDs{}, "worker", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	value, err := processor.ProcessOne(context.Background())
	if err != nil || repo.assigned != 41 || repo.finished != 1 || value.Status != "authorized" {
		t.Fatalf("value=%+v assigned=%d finished=%d err=%v", value, repo.assigned, repo.finished, err)
	}
	repo = &workRepo{invoice: Invoice{ID: "ambiguous", Status: "reconcile_required", VoucherNumber: 41}}
	processor, _ = NewProcessor(repo, providerFake{consult: ok, authorizeErr: errors.New("must not authorize")}, &testIDs{}, "worker", time.Minute)
	value, err = processor.ProcessOne(context.Background())
	if err != nil || repo.finished != 1 || value.CAE != ok.CAE {
		t.Fatalf("reconciliation failed: value=%+v err=%v", value, err)
	}
	repo = &workRepo{invoice: Invoice{ID: "timeout", Status: "queued"}}
	processor, _ = NewProcessor(repo, providerFake{last: 9, authorizeErr: errors.New("timeout")}, &testIDs{}, "worker", time.Minute)
	if _, err = processor.ProcessOne(context.Background()); err == nil || repo.deferred != "authorize" || repo.assigned != 10 {
		t.Fatalf("ambiguous authorization was not deferred: assigned=%d deferred=%s err=%v", repo.assigned, repo.deferred, err)
	}
}
