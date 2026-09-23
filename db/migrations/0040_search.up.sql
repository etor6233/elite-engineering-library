begin;

create schema if not exists search;

create table search.document (
  tenant_id uuid not null,
  organization_id text,
  document_id text not null,
  kind text not null,
  external_id text not null,
  title text not null,
  body text not null default '',
  facets jsonb not null default '{}'::jsonb,
  embedding real[],
  content_sha256 text not null,
  tsv tsvector generated always as (
    setweight(to_tsvector('simple', title), 'A') ||
    setweight(to_tsvector('simple', body), 'B')
  ) stored,
  indexed_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, document_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(document_id) between 1 and 128),
  check (length(external_id) between 1 and 200),
  check (title <> ''),
  check (kind ~ '^[a-z][a-z0-9._-]{0,63}$'),
  check (organization_id is null or length(organization_id) between 1 and 128),
  check (jsonb_typeof(facets) = 'object'),
  check (embedding is null or cardinality(embedding) > 0),
  check (content_sha256 ~ '^[0-9a-f]{64}$'),
  check (updated_at >= indexed_at)
);

create index search_document_tsv_idx on search.document using gin (tsv);
create index search_document_scope_idx on search.document (tenant_id, kind, external_id);
create index search_document_facets_idx on search.document using gin (facets);

commit;
