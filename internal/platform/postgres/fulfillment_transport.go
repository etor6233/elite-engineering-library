package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func scanCustomerTransport(row pgx.Row) (fulfillment.CustomerTransport, error) {
	var value fulfillment.CustomerTransport
	err := row.Scan(&value.ID, &value.RequestID, &value.CustomerShipmentID, &value.OrderID, &value.OriginOrganizationID, &value.DestinationCustomerSubject, &value.ProviderCode, &value.ProviderServiceCode, &value.ProviderReference, &value.State, &value.Version, &value.OrderFulfillmentState)
	return value, err
}

func (r *Fulfillment) readCustomerTransportReplay(ctx context.Context, tx pgx.Tx, tenant string, command fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, bool, error) {
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `select s.shipment_id,coalesce(s.request_id,''),coalesce(s.customer_shipment_id,''),coalesce(s.order_id,''),s.origin_organization_id,coalesce(s.destination_customer_principal_id,''),s.provider_code,coalesce(s.provider_service_code,''),coalesce(s.provider_reference,''),s.state,s.version,o.fulfillment_state from logistics.shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.request_id=$2 for share of s,o`, tenant, command.RequestID))
	if errors.Is(err, pgx.ErrNoRows) {
		return fulfillment.CustomerTransport{}, false, nil
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, false, err
	}
	if value.CustomerShipmentID != command.CustomerShipmentID || value.OriginOrganizationID != command.OriginOrganizationID || value.ProviderCode != command.ProviderCode || value.ProviderServiceCode != command.ProviderServiceCode || value.ProviderReference != command.ProviderReference {
		return fulfillment.CustomerTransport{}, false, fulfillment.ErrConflict
	}
	return value, true, nil
}

func (r *Fulfillment) CreateCustomerTransport(ctx context.Context, tenant string, ids fulfillment.TransportIDs, command fulfillment.CustomerTransportCommand) (fulfillment.CustomerTransport, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	defer tx.Rollback(ctx)
	if replay, ok, err := r.readCustomerTransportReplay(ctx, tx, tenant, command); err != nil || ok {
		return replay, err
	}
	state := "planned"
	if command.ProviderReference != "" {
		state = "booked"
	}
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_service_code,provider_reference,origin_organization_id,destination_organization_id,destination_customer_principal_id,state,version,request_id,customer_shipment_id,order_id) select s.tenant_id,$3,$4,nullif($5,''),nullif($6,''),s.organization_id,null,o.customer_principal_id,$7,1,$8,s.shipment_id,s.order_id from sales.customer_shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.shipment_id=$2 and s.organization_id=$9 returning shipment_id,request_id,customer_shipment_id,order_id,origin_organization_id,destination_customer_principal_id,provider_code,coalesce(provider_service_code,''),coalesce(provider_reference,''),state,version,(select fulfillment_state from sales.customer_order where tenant_id=$1 and order_id=logistics.shipment.order_id)`, tenant, command.CustomerShipmentID, ids.ShipmentID, command.ProviderCode, command.ProviderServiceCode, command.ProviderReference, state, command.RequestID, command.OriginOrganizationID))
	if postgresConflict(err) {
		if replay, ok, replayErr := r.readCustomerTransportReplay(ctx, tx, tenant, command); replayErr != nil || ok {
			return replay, replayErr
		}
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	if err = outbox(ctx, tx, tenant, ids.EventID, "shipment", value.ID, value.Version, "shipment."+state, `jsonb_build_object('customer_shipment_id',$7::text,'order_id',$8::text,'automatic_business_delivery_acceptance',false)`, value.CustomerShipmentID, value.OrderID); err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	return value, tx.Commit(ctx)
}

func providerTransitionAllowed(current, target string) bool {
	return (current == "booked" && (target == "picked-up" || target == "exception")) ||
		(current == "picked-up" && (target == "in-transit" || target == "exception")) ||
		(current == "in-transit" && (target == "in-transit" || target == "delivered" || target == "exception"))
}

func (r *Fulfillment) readProviderReportReplay(ctx context.Context, tx pgx.Tx, tenant, shipmentID string, command fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, bool, error) {
	var storedStatus, storedSchema, storedEvidence string
	var storedOccurred time.Time
	var storedAutomatic bool
	err := tx.QueryRow(ctx, `select provider_status,report_schema,evidence_sha256_hex,occurred_at,automatic_business_delivery_acceptance from logistics.shipment_provider_report where tenant_id=$1 and shipment_id=$2 and provider_event_id=$3`, tenant, shipmentID, command.ProviderEventID).Scan(&storedStatus, &storedSchema, &storedEvidence, &storedOccurred, &storedAutomatic)
	if errors.Is(err, pgx.ErrNoRows) {
		return fulfillment.CustomerTransport{}, false, nil
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, false, err
	}
	if storedStatus != command.ProviderStatus || storedSchema != command.ReportSchema || storedEvidence != command.EvidenceSHA256 || !storedOccurred.Equal(command.OccurredAt) || storedAutomatic != command.AutomaticBusinessDeliveryAcceptance {
		return fulfillment.CustomerTransport{}, false, fulfillment.ErrConflict
	}
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `select s.shipment_id,s.request_id,s.customer_shipment_id,s.order_id,s.origin_organization_id,s.destination_customer_principal_id,s.provider_code,coalesce(s.provider_service_code,''),coalesce(s.provider_reference,''),s.state,s.version,o.fulfillment_state from logistics.shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.shipment_id=$2`, tenant, shipmentID))
	return value, true, err
}

func (r *Fulfillment) RecordProviderReport(ctx context.Context, tenant, shipmentID, eventID string, command fulfillment.ProviderReportCommand) (fulfillment.CustomerTransport, error) {
	target, err := fulfillment.NormalizeProviderReport(command)
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	defer tx.Rollback(ctx)
	if replay, ok, err := r.readProviderReportReplay(ctx, tx, tenant, shipmentID, command); err != nil || ok {
		return replay, err
	}
	var current, orderID string
	var version int64
	err = tx.QueryRow(ctx, `select state,version,order_id from logistics.shipment where tenant_id=$1 and shipment_id=$2 and origin_organization_id=$3 and customer_shipment_id is not null for update`, tenant, shipmentID, command.OriginOrganizationID).Scan(&current, &version, &orderID)
	if errors.Is(err, pgx.ErrNoRows) || version != command.Version || !providerTransitionAllowed(current, target) {
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	_, err = tx.Exec(ctx, `insert into logistics.shipment_provider_report(tenant_id,shipment_id,provider_event_id,provider_status,normalized_state,report_schema,evidence_sha256_hex,occurred_at,automatic_business_delivery_acceptance,shipment_version) values($1,$2,$3,$4,$5,$6,$7,$8,false,$9)`, tenant, shipmentID, command.ProviderEventID, command.ProviderStatus, target, command.ReportSchema, command.EvidenceSHA256, command.OccurredAt, version+1)
	if postgresConflict(err) {
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	result, err := tx.Exec(ctx, `update logistics.shipment set state=$4,version=version+1,last_event_at=$5 where tenant_id=$1 and shipment_id=$2 and version=$3`, tenant, shipmentID, version, target, command.OccurredAt)
	if err != nil || result.RowsAffected() != 1 {
		if err != nil {
			return fulfillment.CustomerTransport{}, err
		}
		return fulfillment.CustomerTransport{}, fulfillment.ErrConflict
	}
	if target == "delivered" {
		_, err = tx.Exec(ctx, `update sales.customer_order o set fulfillment_state='delivered',version=version+1,updated_at=clock_timestamp() where o.tenant_id=$1 and o.order_id=$2 and o.fulfillment_state='shipped' and not exists (select 1 from sales.customer_shipment cs left join logistics.shipment ls on ls.tenant_id=cs.tenant_id and ls.customer_shipment_id=cs.shipment_id where cs.tenant_id=o.tenant_id and cs.order_id=o.order_id and (ls.shipment_id is null or ls.state<>'delivered'))`, tenant, orderID)
		if err != nil {
			return fulfillment.CustomerTransport{}, err
		}
	}
	if err = outbox(ctx, tx, tenant, eventID, "shipment", shipmentID, version+1, "shipment.provider-report-recorded", `jsonb_build_object('provider_status',$7::text,'normalized_state',$8::text,'evidence_sha256',$9::text,'automatic_business_delivery_acceptance',false)`, command.ProviderStatus, target, command.EvidenceSHA256); err != nil {
		return fulfillment.CustomerTransport{}, err
	}
	value, err := scanCustomerTransport(tx.QueryRow(ctx, `select s.shipment_id,s.request_id,s.customer_shipment_id,s.order_id,s.origin_organization_id,s.destination_customer_principal_id,s.provider_code,coalesce(s.provider_service_code,''),coalesce(s.provider_reference,''),s.state,s.version,o.fulfillment_state from logistics.shipment s join sales.customer_order o on o.tenant_id=s.tenant_id and o.order_id=s.order_id where s.tenant_id=$1 and s.shipment_id=$2`, tenant, shipmentID))
	if err != nil {
		return fulfillment.CustomerTransport{}, fmt.Errorf("read connected transport: %w", err)
	}
	return value, tx.Commit(ctx)
}
