drop table if exists fiscal.parameter_refresh_schedule;
drop trigger if exists fiscal_parameter_decision_reason on fiscal.parameter_decision;
drop function if exists fiscal.require_parameter_decision_reason();
alter table fiscal.parameter_decision drop column if exists decision_reason;
