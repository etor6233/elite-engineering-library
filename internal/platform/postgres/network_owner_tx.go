package postgres

// AUTHORED narrow adapter: only four existing network writers are exposed.
import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"github.com/jackc/pgx/v5"
)

type networkTxRepository struct{ tx pgx.Tx }

var _ fulfillment.Repository = networkTxRepository{}

func (r networkTxRepository) CreateShipment(a0 context.Context, a1 string, a2 string, a3 fulfillment.Shipment) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) TransitionShipment(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string, a8 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) CreateCustomerTransport(a0 context.Context, a1 string, a2 fulfillment.TransportIDs, a3 fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, error) {
	return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
}
func (r networkTxRepository) RecordProviderReport(a0 context.Context, a1 string, a2 string, a3 string, a4 fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, error) {
	return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
}
func (r networkTxRepository) OpenServiceCase(a0 context.Context, a1 string, a2 string, a3 fulfillment.ServiceCase) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) TransitionServiceCase(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 int64, a7 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) CreateRecall(a0 context.Context, a1 string, a2 string, a3 fulfillment.Recall) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) ActivateRecall(a0 context.Context, a1 string, a2 string, a3 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) AddRecallUnit(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) CreateOrganization(a0 context.Context, a1 string, a2 string, a3 fulfillment.Organization) error {
	return createOrganizationInTx(a0, r.tx, a1, a2, a3)
}
func (r networkTxRepository) TransitionOrganization(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 int64, a6 string) error {
	return transitionOrganizationInTx(a0, r.tx, a1, a2, a3, a4, a5, a6)
}
func (r networkTxRepository) CreateAgreement(a0 context.Context, a1 string, a2 string, a3 fulfillment.Agreement) error {
	return createAgreementInTx(a0, r.tx, a1, a2, a3)
}
func (r networkTxRepository) TransitionAgreement(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 int64, a7 string) error {
	return transitionAgreementInTx(a0, r.tx, a1, a2, a3, a4, a5, a6, a7)
}
func (r networkTxRepository) QueueMessage(a0 context.Context, a1 string, a2 string, a3 fulfillment.Message) error {
	return fulfillment.ErrConflict
}
func (r networkTxRepository) TransitionMessage(a0 context.Context, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string) error {
	return fulfillment.ErrConflict
}
