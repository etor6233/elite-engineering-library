begin;

drop index if exists sales.customer_shipment_order_idx;
drop trigger if exists customer_shipment_allocation_immutable on inventory.customer_shipment_allocation;
drop trigger if exists customer_shipment_line_immutable on sales.customer_shipment_line;
drop trigger if exists customer_shipment_immutable on sales.customer_shipment;
drop function if exists sales.reject_customer_shipment_mutation();
drop table if exists inventory.customer_shipment_allocation;
drop table if exists sales.customer_shipment_line;
drop table if exists sales.customer_shipment;
alter table sales.customer_order_line drop column if exists shipped_quantity;
alter table sales.customer_order drop column if exists fulfillment_state;

commit;
