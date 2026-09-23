begin;
-- AUTHORED: no campaign, invitation, consent policy or activation is inferred.
create table crm.survey_definition (
 tenant_id uuid not null,
 organization_id text not null,
 survey_id text not null check(survey_id ~ '^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$'),
 prompt text not null check(length(prompt) between 1 and 1000),
 consent_version text not null check(consent_version ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
 consent_notice text not null check(length(consent_notice) between 1 and 4000),
 opens_at timestamptz not null,
 closes_at timestamptz not null,
 retain_until timestamptz not null,
 minimum_responses integer not null check(minimum_responses between 1 and 1000000),
 active boolean not null default false,
 primary key(tenant_id,organization_id,survey_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check(opens_at<closes_at and closes_at<retain_until)
);
create table crm.survey_response (
 tenant_id uuid not null,
 organization_id text not null,
 survey_id text not null,
 customer_principal_id text not null,
 score smallint not null check(score between 0 and 10),
 consent_version text not null,
 received_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,organization_id,survey_id,customer_principal_id),
 foreign key(tenant_id,organization_id,survey_id) references crm.survey_definition(tenant_id,organization_id,survey_id),
 foreign key(tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id)
);
-- Definitions are versioned by survey ID. Activation may change; historical terms may not.
create function crm.survey_definition_immutable() returns trigger language plpgsql as $$
begin
 if (to_jsonb(new)-'active') is distinct from (to_jsonb(old)-'active') then
  raise exception using errcode='23514',message='survey definition is immutable; create a new survey id';
 end if;
 return new;
end $$;
create trigger survey_definition_immutable before update on crm.survey_definition
 for each row execute function crm.survey_definition_immutable();
commit;
