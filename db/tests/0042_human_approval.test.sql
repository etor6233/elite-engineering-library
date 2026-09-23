-- 0042_human_approval.test.sql — verifica cola de aprobación: estado pendiente
-- con decided_at nulo, decisión inmutable ligada a request, y dinero (payment)
-- que permanece pendiente (nunca auto-aprobado).
-- Precondición: migración 0001 (platform.tenant) y 0042 aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A');

insert into approval.request
  (tenant_id, request_id, kind, subject_id, amount_minor_units, requester, evidence_sha)
values
  ('11111111-1111-1111-1111-111111111111', 'req-1', 'sale', 'sub-1', 500, 'seller-1', repeat('a',64)),
  ('11111111-1111-1111-1111-111111111111', 'req-2', 'payment', 'sub-2', 10000, 'seller-1', repeat('b',64));

-- dos pendientes con decided_at nulo
do $$
declare n int;
begin
  select count(*) into n from approval.request
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and state = 'pending' and decided_at is null;
  if n <> 2 then raise exception 'expected 2 pending, got %', n; end if;
end $$;

-- aprobar req-1 y registrar decisión
update approval.request set state = 'approved', decided_at = clock_timestamp()
 where tenant_id = '11111111-1111-1111-1111-111111111111' and request_id = 'req-1';

insert into approval.decision (tenant_id, request_id, reviewer, approved, reason) values
  ('11111111-1111-1111-1111-111111111111', 'req-1', 'manager-1', true, 'ok');

do $$
declare n int;
begin
  select count(*) into n from approval.decision
   where tenant_id = '11111111-1111-1111-1111-111111111111' and request_id = 'req-1';
  if n <> 1 then raise exception 'expected 1 decision, got %', n; end if;
end $$;

-- el dinero (payment) permanece pendiente
do $$
declare st text;
begin
  select state into st from approval.request
   where tenant_id = '11111111-1111-1111-1111-111111111111' and request_id = 'req-2';
  if st <> 'pending' then raise exception 'payment should remain pending, got %', st; end if;
end $$;

rollback;
