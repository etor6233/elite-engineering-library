package fiscal

import (
	"context"
	"errors"
	"time"
)

var ErrNoWork = errors.New("no fiscal work available")

type Provider interface {
	LastAuthorized(context.Context, Invoice) (int64, string, error)
	Consult(context.Context, Invoice) (Authorization, error)
	Authorize(context.Context, Invoice) (Authorization, error)
}

type WorkRepository interface {
	ClaimInvoice(context.Context, string, time.Duration, string) (Invoice, error)
	AssignVoucherNumber(context.Context, Invoice, string, int64, string, string) (Invoice, error)
	FinishInvoice(context.Context, Invoice, string, Authorization, string) (Invoice, error)
	DeferInvoice(context.Context, Invoice, string, string, string, bool) error
}

type Processor struct {
	repository WorkRepository
	provider   Provider
	ids        IDGenerator
	workerID   string
	lease      time.Duration
}

func NewProcessor(repository WorkRepository, provider Provider, ids IDGenerator, workerID string, lease time.Duration) (*Processor, error) {
	if repository == nil || provider == nil || ids == nil || workerID == "" || lease < 30*time.Second || lease > 10*time.Minute {
		return nil, ErrInvalid
	}
	return &Processor{repository: repository, provider: provider, ids: ids, workerID: workerID, lease: lease}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Invoice, error) {
	invoice, err := p.repository.ClaimInvoice(ctx, p.workerID, p.lease, p.ids.New())
	if err != nil {
		return invoice, err
	}
	if invoice.VoucherNumber > 0 {
		consulted, consultErr := p.provider.Consult(ctx, invoice)
		if consultErr != nil {
			_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "consult", "", true)
			return invoice, consultErr
		}
		if !validAuthorization(consulted) {
			_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "consult-invalid", consulted.ResponseHash, true)
			return invoice, ErrInvalid
		}
		if consulted.Found {
			return p.repository.FinishInvoice(ctx, invoice, p.workerID, consulted, p.ids.New())
		}
	}
	if invoice.VoucherNumber == 0 {
		last, responseHash, sequenceErr := p.provider.LastAuthorized(ctx, invoice)
		if sequenceErr != nil || last < 0 || !hashPattern.MatchString(responseHash) {
			_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "last-authorized", responseHash, false)
			if sequenceErr != nil {
				return invoice, sequenceErr
			}
			return invoice, ErrInvalid
		}
		invoice, err = p.repository.AssignVoucherNumber(ctx, invoice, p.workerID, last+1, responseHash, p.ids.New())
		if err != nil {
			return invoice, err
		}
	}
	authorized, authorizationErr := p.provider.Authorize(ctx, invoice)
	if authorizationErr != nil {
		_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "authorize", "", true)
		return invoice, authorizationErr
	}
	if !validAuthorization(authorized) || !authorized.Found {
		_ = p.repository.DeferInvoice(ctx, invoice, p.workerID, "authorize-invalid", authorized.ResponseHash, true)
		return invoice, ErrInvalid
	}
	return p.repository.FinishInvoice(ctx, invoice, p.workerID, authorized, p.ids.New())
}
