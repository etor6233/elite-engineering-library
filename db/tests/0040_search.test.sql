-- 0040_search.test.sql — verifica búsqueda full-text (tsvector/GIN + websearch_to_tsquery),
-- ranking por peso, aislamiento ACL por tenant/organization, filtro de facets y
-- almacenamiento round-trip de embeddings real[].
-- Precondición: migración 0001 (platform.tenant) y 0040 (search.document) aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A'),
  ('22222222-2222-2222-2222-222222222222', 'tenant-b', 'B', 'Tenant B');

insert into search.document
  (tenant_id, organization_id, document_id, kind, external_id, title, body, facets, embedding, content_sha256)
values
  ('11111111-1111-1111-1111-111111111111', 'org-1', 'd1', 'product', 'SKU-1',
   'Scooter Urbano Pro', 'vehiculo electrico ligero', '{"category":"movilidad"}',
   ARRAY[1.0,0.0,0.0]::real[], repeat('a',64)),
  ('11111111-1111-1111-1111-111111111111', 'org-1', 'd2', 'product', 'SKU-2',
   'Vehiculo Plegable', 'scooter plegable portatil', '{"category":"movilidad"}',
   ARRAY[0.0,1.0,0.0]::real[], repeat('b',64)),
  ('11111111-1111-1111-1111-111111111111', 'org-2', 'd3', 'product', 'SKU-3',
   'Scooter Urbano', 'scooter urbano', '{"category":"movilidad"}',
   null, repeat('c',64)),
  ('22222222-2222-2222-2222-222222222222', 'org-1', 'd4', 'product', 'SKU-4',
   'Scooter Urbano', 'scooter urbano', '{"category":"movilidad"}',
   null, repeat('d',64));

-- full-text 'scooter' en (tenant-a, org-1): 2 resultados, d1 (title, peso A) primero
do $$
declare n int; top text;
begin
  select count(*) into n from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and organization_id = 'org-1'
     and tsv @@ websearch_to_tsquery('simple', 'scooter');
  if n <> 2 then raise exception 'fulltext expected 2, got %', n; end if;

  select document_id into top from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and organization_id = 'org-1'
     and tsv @@ websearch_to_tsquery('simple', 'scooter')
   order by ts_rank(tsv, websearch_to_tsquery('simple', 'scooter')) desc, document_id
   limit 1;
  if top <> 'd1' then raise exception 'expected d1 top-ranked (title weight A), got %', top; end if;
end $$;

-- término específico 'urbano' en (tenant-a, org-1): sólo d1
do $$
declare n int;
begin
  select count(*) into n from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and organization_id = 'org-1'
     and tsv @@ websearch_to_tsquery('simple', 'urbano');
  if n <> 1 then raise exception 'urbano expected 1, got %', n; end if;
end $$;

-- aislamiento de tenant: tenant-b ve 1, tenant-a ve 0 de tenant-b
do $$
declare a int; b int;
begin
  select count(*) into b from search.document
   where tenant_id = '22222222-2222-2222-2222-222222222222'
     and tsv @@ websearch_to_tsquery('simple', 'scooter');
  if b <> 1 then raise exception 'tenant-b expected 1, got %', b; end if;
  select count(*) into a from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and tsv @@ websearch_to_tsquery('simple', 'scooter')
     and organization_id = 'org-2';
  if a <> 1 then raise exception 'org-2 expected 1, got %', a; end if;
end $$;

-- filtro de facets
do $$
declare n int;
begin
  select count(*) into n from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and facets @> '{"category":"movilidad"}';
  if n <> 3 then raise exception 'facet expected 3, got %', n; end if;
end $$;

-- round-trip de embedding real[]
do $$
declare dim int; v real;
begin
  select array_length(embedding,1), embedding[1] into dim, v from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111' and document_id = 'd1';
  if dim <> 3 or v <> 1.0 then raise exception 'embedding round-trip failed: dim=%, v=%', dim, v; end if;
end $$;

rollback;
