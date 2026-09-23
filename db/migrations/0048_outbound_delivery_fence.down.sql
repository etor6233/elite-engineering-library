begin;
drop table communication.outbound_delivery_event;
drop table communication.outbound_delivery;
drop function communication.reject_outbound_delivery_event_mutation();
commit;
