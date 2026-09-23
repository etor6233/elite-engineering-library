begin;

drop trigger if exists lead_candidate_immutable on integration.lead_candidate;
drop trigger if exists lead_ingress_raw_immutable on integration.lead_ingress_raw;
drop function if exists integration.reject_lead_ingress_mutation();
drop table if exists integration.lead_candidate;
drop table if exists integration.lead_ingress_raw;

commit;
