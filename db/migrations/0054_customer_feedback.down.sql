begin;
drop table crm.survey_response;
drop trigger survey_definition_immutable on crm.survey_definition;
drop function crm.survey_definition_immutable();
drop table crm.survey_definition;

commit;
