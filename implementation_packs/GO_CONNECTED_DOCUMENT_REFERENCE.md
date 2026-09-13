# Go Connected Document Reference

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-DOCUMENT-REFERENCE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Connected local fixture document receive/extract/review/commit infrastructure with original pinned SDK, existing security-gate wiring, PostgreSQL fencing and separate human approval; explicitly simulated detector/provider results, no OCR accuracy or automatic accounting claim."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-AWS-TEXTRACT-DOCUMENT-RUNTIME 0.2.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "GO-ENTERPRISE-BACKEND 0.4.x", "GO-HUMAN-APPROVAL-CORE"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/aws/aws-sdk-go-v2/tree/a30468cff35d6e385287a0cbff0ac11aa7202529/service/textract", "https://github.com/Azure/azure-sdk-for-python/tree/8555d14532a9688b751d8408d822d1dd5feb47f6"]
verified_at: "2026-09-13"
```

## 2. Applicability

Required supplier-invoice reference in LIBRARY_INFRASTRUCTURE only. Full composition FRANCHISE_COMPLETE_PACK_PLAN selects existing SDK/security/fixture owners plus this glue. Every private class/corpus remains governed by docs/DOCUMENT_REFERENCE_DECISION.md.

## 3. Architecture contract

Atomic original+job; original SDK and existing security gate; immutable extraction; bounded existing PostgreSQL job generation fencing; four-field immutable proposal; separate human approval; atomic commit/approval/outbox and exact replay. FIXTURE marked permanently and restricted to one public digest. PROVIDER never falls back to simulated detectors.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/document.go
CREATE cmd/electromobility-api/document_test.go
CREATE config/documents/reference-profile.json
CREATE db/migrations/0085_document_review.down.sql
CREATE db/migrations/0085_document_review.up.sql
CREATE docs/DOCUMENT_REFERENCE_DECISION.md
CREATE docs/DOCUMENT_REFERENCE_START.md
CREATE internal/approval/document_test.go
CREATE internal/documentbridge/fixture.go
CREATE internal/documentbridge/pipeline.go
CREATE internal/documentbridge/pipeline_test.go
CREATE internal/platform/httpapi/document.go
CREATE internal/platform/postgres/document.go
CREATE internal/platform/postgres/document_connected_integration_test.go
CREATE tools/verify_document_reference.py
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/document.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0563dd66231952cffc7e3b93a6f88ee31d38a3651fb34236f279a99c6aeb8fe7"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED opt-in host composition. FIXTURE uses an in-process no-network HTTP
// transport. PROVIDER requires the original file gate and external AWS credentials.
import (
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	runtime "example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() { documentModuleFactory = selectedDocumentModule }
func selectedDocumentModule(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	if getenv == nil {
		return nil, doc.ErrContract
	}
	enabled := getenv("DOCUMENTS_ENABLED")
	if enabled == "" || enabled == "false" {
		return nil, nil
	}
	if enabled != "true" || pool == nil {
		return nil, doc.ErrContract
	}
	var ready bool
	e := pool.QueryRow(ctx, `select(select count(*)from pg_trigger where not tgisinternal and tgenabled in('O','A')and(tgrelid=to_regclass('document.original')and tgname='original_immutable' or tgrelid=to_regclass('document.extraction')and tgname='extraction_immutable' or tgrelid=to_regclass('document.committed')and tgname in('committed_immutable','document_commit_guard') or tgrelid=to_regclass('approval.request')and tgname='document_review_guard'))=5 and to_regclass('approval.document_review_once')is not null`).Scan(&ready)
	if e != nil || !ready {
		return nil, doc.ErrContract
	}
	mode, profile, hash, root := getenv("DOCUMENTS_MODE"), getenv("DOCUMENTS_PROFILE_FILE"), getenv("DOCUMENTS_PROFILE_SHA256"), getenv("DOCUMENTS_WORK_ROOT")
	var processor *doc.Pipeline
	switch mode {
	case "FIXTURE":
		processor, e = doc.NewFixturePipeline(profile, hash, root)
	case "PROVIDER":
		command := doc.SecurityCommand{Python: getenv("DOCUMENTS_PYTHON"), PythonSHA: getenv("DOCUMENTS_PYTHON_SHA256"), Script: getenv("DOCUMENTS_SECURITY_SCRIPT"), ScriptSHA: getenv("DOCUMENTS_SECURITY_SCRIPT_SHA256"), Policy: getenv("DOCUMENTS_SECURITY_POLICY"), PolicySHA: getenv("DOCUMENTS_SECURITY_POLICY_SHA256"), Magika: getenv("DOCUMENTS_MAGIKA"), ClamScan: getenv("DOCUMENTS_CLAMSCAN"), ClamDatabase: getenv("DOCUMENTS_CLAM_DATABASE"), Yara: getenv("DOCUMENTS_YARA"), Rules: getenv("DOCUMENTS_YARA_RULES")}
		var stage doc.SecurityStage
		stage, e = command.Stage()
		if e != nil {
			return nil, e
		}
		region, err := runtime.ApprovedRegion(profile, doc.Class)
		if err != nil {
			return nil, err
		}
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			return nil, doc.ErrContract
		}
		processor, e = doc.NewPipeline(textract.NewFromConfig(cfg, func(o *textract.Options) { o.RetryMaxAttempts = 1 }), stage, profile, hash, root, mode)
	default:
		return nil, doc.ErrContract
	}
	if e != nil {
		return nil, e
	}
	store, e := postgres.NewDocuments(pool, doc.Scope{TenantID: getenv("DOCUMENTS_TENANT_ID"), OrganizationID: getenv("DOCUMENTS_ORGANIZATION_ID"), ProfileSHA: hash, Mode: mode}, processor)
	if e != nil {
		return nil, e
	}
	return &httpapi.DocumentModule{Store: store}, nil
}
````

### FILE: `cmd/electromobility-api/document_test.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "17b39cc2741931d4e30014ecb61c476f3130b80b42e29998e2dc80c8e055e07f"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentHostActivation(t *testing.T) {
	ctx := context.Background()
	none := func(string) string { return "" }
	if m, e := selectedDocumentModule(ctx, nil, none); e != nil || m != nil {
		t.Fatal("default activation")
	}
	url := os.Getenv("DOCUMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned document fixture database required")
	}
	parsed, parseErr := urlpkg.Parse(url)
	if parseErr != nil || (parsed.Hostname() != "127.0.0.1" && parsed.Hostname() != "localhost") || (!strings.HasPrefix(parsed.Path, "/elite_document_reference_") && !strings.HasPrefix(parsed.Path, "/elite_payment_connected_")) {
		t.Fatal("owned loopback fixture database required")
	}
	pool, e := pgxpool.New(ctx, url)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	profile := filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "config", "documents", "reference-profile.json")
	b, e := os.ReadFile(profile)
	if e != nil {
		t.Fatal(e)
	}
	env := map[string]string{"DOCUMENTS_ENABLED": "true", "DOCUMENTS_MODE": "FIXTURE", "DOCUMENTS_PROFILE_FILE": profile, "DOCUMENTS_PROFILE_SHA256": doc.Hash(b), "DOCUMENTS_WORK_ROOT": t.TempDir(), "DOCUMENTS_TENANT_ID": uuid.NewString(), "DOCUMENTS_ORGANIZATION_ID": "fixture-org"}
	get := func(k string) string { return env[k] }
	if m, e := selectedDocumentModule(ctx, pool, get); e != nil || m == nil {
		t.Fatal("selected fixture module unavailable", e)
	}
	env["DOCUMENTS_MODE"] = "PROVIDER"
	if _, e = selectedDocumentModule(ctx, pool, get); e == nil {
		t.Fatal("provider mode accepted missing security runtime")
	}
	env["DOCUMENTS_MODE"] = "FIXTURE"
	env["DOCUMENTS_PROFILE_SHA256"] = "bad"
	if _, e = selectedDocumentModule(ctx, pool, get); e == nil {
		t.Fatal("unbound profile")
	}
	env["DOCUMENTS_PROFILE_SHA256"] = doc.Hash(b)
	_, e = pool.Exec(ctx, `alter table document.committed disable trigger document_commit_guard`)
	if e != nil {
		t.Fatal(e)
	}
	_, e = selectedDocumentModule(ctx, pool, get)
	_, restore := pool.Exec(ctx, `alter table document.committed enable trigger document_commit_guard`)
	if restore != nil {
		t.Fatal(restore)
	}
	if e == nil {
		t.Fatal("host accepted absent approval/commit guard")
	}
	t.Log("DOCUMENT_HOST_PASS: opt-in fixture module; provider runtime required; profile hash; database guard refusal")
}
````

### FILE: `config/documents/reference-profile.json`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0e11ae43fefb995ff18d01541a6da19a52bee75d3660d2ebf49eb4e1f55e2887"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-aws-textract-document-profile/v2",
  "provider": "Amazon Textract",
  "sdk": "github.com/aws/aws-sdk-go-v2/service/textract@v1.45.0",
  "region": "us-east-1",
  "classes": [
    {
      "id": "supplier-invoice",
      "decision": "REQUIRED",
      "operation": "AnalyzeExpense",
      "features": [],
      "queries": [],
      "response_schema_version": "supplier-invoice/v1",
      "automatic_storage": false
    }
  ]
}
````

### FILE: `db/migrations/0085_document_review.down.sql`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bac0405f1e3502fc5877877e673dfde2375ad142549d9f29354b3452a9cc1e18"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from document.original)or exists(select 1 from approval.request where kind='document_review')
 then raise exception 'document evidence exists: populated rollback forbidden';end if;
end$$;
drop trigger document_review_guard on approval.request;
drop index approval.document_review_once;
drop schema document cascade;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;



commit;
````

### FILE: `db/migrations/0085_document_review.up.sql`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a689fbccc554062d22bd01851ce39926ee2721a61dd0d2c4703274f7e8bea402"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;



create schema document;
create table document.original(
 tenant_id uuid not null,document_id uuid not null,organization_id text not null,
 uploader text not null,name text not null,original_sha256 text not null,
 profile_sha256 text not null,mode text not null,content bytea not null,
 received_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id),
 foreign key(tenant_id,organization_id)references org.organization(tenant_id,organization_id),
 check(length(uploader)between 1 and 128),check(length(name)between 1 and 128),
 check(original_sha256~'^[0-9a-f]{64}$'and profile_sha256~'^[0-9a-f]{64}$'),
 check(encode(sha256(content),'hex')=original_sha256),
 check(octet_length(content)between 1 and 2097152),check(mode in('FIXTURE','PROVIDER')),
 constraint original_fixture_identity check(mode<>'FIXTURE'or original_sha256='489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb')
);
create table document.extraction(
 tenant_id uuid not null,document_id uuid not null,job_id uuid not null,attempt integer not null,
 security_receipt bytea not null,provider_response bytea not null,analysis_receipt bytea not null,
 evidence_sha256 text not null,suggested jsonb not null,
 extracted_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id),
 foreign key(tenant_id,document_id)references document.original(tenant_id,document_id),
 foreign key(tenant_id,job_id)references platform.job(tenant_id,job_id),
 check(attempt>0),check(evidence_sha256~'^[0-9a-f]{64}$'),
 check(encode(sha256(analysis_receipt),'hex')=evidence_sha256),
 check(octet_length(security_receipt)between 1 and 65536),
 check(octet_length(provider_response)between 1 and 4194304),
 check(octet_length(analysis_receipt)between 1 and 32768),
 check(jsonb_typeof(suggested)='object')
);
create table document.attempt_failure(
 tenant_id uuid not null,document_id uuid not null,attempt integer not null,code text not null,
 failed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id,attempt),
 foreign key(tenant_id,document_id)references document.original(tenant_id,document_id),
 check(code in('SECURITY_REJECTED','EXTRACTION_UNAVAILABLE','INVALID_CONTRACT'))
);
create table document.committed(
 tenant_id uuid not null,document_id uuid not null,request_id text not null,
 payload_sha256 text not null,fields jsonb not null,reviewer text not null,
 committed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id),unique(tenant_id,request_id),
 foreign key(tenant_id,document_id)references document.extraction(tenant_id,document_id),
 foreign key(tenant_id,request_id)references approval.request(tenant_id,request_id),
 check(payload_sha256~'^[0-9a-f]{64}$'),check(jsonb_typeof(fields)='object'),check(length(reviewer)between 1 and 128)
);
create unique index document_review_once on approval.request(tenant_id,organization_id,(payload->>'document_id'))where kind='document_review';
create trigger original_immutable before update or delete on document.original for each row execute function catalog.release_immutable();
create trigger extraction_immutable before update or delete on document.extraction for each row execute function catalog.release_immutable();
create trigger attempt_failure_immutable before update or delete on document.attempt_failure for each row execute function catalog.release_immutable();
create trigger committed_immutable before update or delete on document.committed for each row execute function catalog.release_immutable();
create function document.guard_review()returns trigger language plpgsql as $$
begin
 if new.kind='document_review' and not exists(
 select 1 from document.original o join document.extraction e using(tenant_id,document_id)
 where o.tenant_id=new.tenant_id and o.document_id::text=new.payload->>'document_id'
 and o.organization_id=new.organization_id and o.original_sha256=new.payload->>'original_sha256'
 and o.profile_sha256=new.payload->>'profile_sha256'and o.mode=new.payload->>'mode'
 and e.evidence_sha256=new.payload->>'evidence_sha256'
 and new.payload->>'schema'='document-review/v1'
 and jsonb_typeof(new.payload->'fields')='object'
 and(select count(*)from jsonb_object_keys(new.payload->'fields'))=4
 and jsonb_typeof(new.payload->'fields'->'invoice_number')='string'and octet_length(new.payload->'fields'->>'invoice_number')between 1 and 128
 and jsonb_typeof(new.payload->'fields'->'vendor')='string'and octet_length(new.payload->'fields'->>'vendor')between 1 and 512
 and jsonb_typeof(new.payload->'fields'->'total')='string'and octet_length(new.payload->'fields'->>'total')between 1 and 64
 and jsonb_typeof(new.payload->'fields'->'currency')='string'and(new.payload->'fields'->>'currency')~'^[A-Z]{3}$'
 and not exists(select 1 from jsonb_each_text(new.payload->'fields')f where f.value<>btrim(f.value)or f.value='')
 )then raise exception 'document review evidence not bound';end if;
 return new;
end$$;
create trigger document_review_guard before insert on approval.request for each row execute function document.guard_review();
create function document.guard_commit()returns trigger language plpgsql as $$
begin
 if not exists(select 1 from approval.request r join approval.decision d using(tenant_id,request_id)
 where r.tenant_id=new.tenant_id and r.request_id=new.request_id and r.kind='document_review'
 and r.state='approved'and d.approved and d.reviewer<>r.requester and d.reviewer=new.reviewer
 and r.evidence_sha=new.payload_sha256 and r.payload->>'document_id'=new.document_id::text
 and r.payload->'fields'=new.fields
 and(select count(*)from approval.decision x where x.tenant_id=r.tenant_id and x.request_id=r.request_id)=1)
 then raise exception 'document commit requires exact human approval';end if;
 return new;
end$$;
create constraint trigger document_commit_guard after insert on document.committed deferrable initially deferred for each row execute function document.guard_commit();
commit;
````

### FILE: `docs/DOCUMENT_REFERENCE_DECISION.md`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d24e5cf8ad6dc948f69b93d560c089e6502ae541fe9d976137f00b92e7ffff6c"
variables: []
secrets_allowed: false
```

````markdown
# Decisión documental de la composición de referencia

Derivada de PROJECT_DOCUMENT_INTELLIGENCE_DECISION_TEMPLATE.md. Alcance único:
LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES / REVIEW_ONLY. No habilita documentos
privados, exactitud OCR, contabilización, reglas fiscales ni almacenamiento automático.
La autorización del usuario es completar infraestructura sin cuentas ni secretos.
Owner de este perfil: mantenimiento de biblioteca; operación/revisión/incidentes:
agente durante el ensayo, usuario cuando materialice su propio target.

## Ingesta y conservación

Canal REQUIRED: API autenticada por el verifier OIDC existente, permiso
`documents:write`, tenant y organización. Un original por PUT y UUID; bytes
SHA-256 y profile SHA obligatorios. Mismo ID+actor+nombre+bytes+perfil devuelve
el mismo registro; diferencias rechazan. JPEG/PDF, máximo2MiB, una página para
el processor. FIXTURE sólo admite el JPEG público exacto, no cualquier PDF.
El cuerpo se recibe completo y acotado antes de abrir la transacción; cargas
parciales no crean original/job. Timeout de proceso2min; lease3min;3intentos;
una extracción por job reclamado. El operador posee `documents:process`.

El original se conserva como bytea con hash calculado también por PostgreSQL,
FK tenant/org, timestamp de la base y trigger que impide update/delete. Es
inmutabilidad de aplicación/base bajo roles acotados, no WORM frente al DBA.
Un superusuario sigue pudiendo cambiar la infraestructura. FIXTURE no contiene
cuentas/secretos del usuario. En este ensayo local, acceso del SO restringido al
owner; cifrado/KMS, residencia, plazo legal, legal hold y borrado de documentos
privados son decisiones del futuro target, sin defaults regulatorios inventados.
La retención del fixture dura el ensayo y se conserva en su evidencia; el módulo
no borra evidencia ni ejecuta down con documentos presentes. Temporales por
intento se eliminan al retornar; huérfanos por caída del host los gestiona el
supervisor local/operación T2809. No se promete eliminación segura del disco.

Estados proyectados: QUARANTINED, REVIEW_REQUIRED, REVIEW_PENDING, REJECTED,
PERSISTED y QUARANTINE_TERMINAL. El hash original, perfil, security receipt,
respuesta completa tipada y analysis receipt se enlazan antes de revisión.
GET de original y tres partes de evidencia está autenticado, devuelve bytes
exactos con SHA, no-store, attachment y nosniff. No se registran campos privados
ni salidas de detectores en logs. El outbox sólo lleva IDs/hash/mode.

## Seguridad y ejecución

PROVIDER llama el gate existente SECURE_LOCAL_FILE_INGESTION_GATE: intérprete,
script y policy con SHA fijados; ClamAV1.5.4, Magika1.1.0 y YARA-X1.20.0 conservan
sus propios hashes/contratos. Firmas oficiales frescas, type/MIME concordante,
reglas YARA fijadas, límites y rechazo permanecen en ese owner. Archives,
macros, archivos cifrados, multipágina y clases desconocidas no se habilitan.
No hay override por receipt aportado desde HTTP. El SDK se llama después del
gate, con una sola tentativa interna. Los reintentos de la cola son explícitos;
una respuesta perdida puede repetir una extracción remota de sólo lectura y
su costo, nunca autoriza duplicar el commit. Fallo de scanner/provider/schema
conserva código seguro por intento y cuarentena; agotamiento permanece terminal.

FIXTURE usa transporte HTTP en memoria del SDK original y contrato de detector
simulado, limitado al SHA público. `security_claim`, mode y commit lo declaran.
No prueba eficacia de antivirus, firmas reales ni libxml2. Esa rama nativa
ACCESS_BLOCKED/Daybreak queda para el expediente final, sin investigación aquí.
No se presenta como una credencial faltante ni se oculta código incompleto.
El runtime PROVIDER está conectado al gate; sólo el ensayo fixture está probado.
No hay datos sensibles a redactar antes del proveedor fixture: no existe egress.
Para datos privados, clasificación/redacción/retención y scanners del target
son condiciones de activación, no hechos inferidos del PASS de este ensayo.

## Clases, una decisión por clase

| Clase | Decisión referencia | Variante / motivo | Owner |
|---|---|---|---|
| invoice / supplier invoice | REQUIRED / REVIEW_ONLY | JPEG público Microsoft fijado; respuesta AWS simulada, sin ground truth de exactitud | mantenimiento de biblioteca |
| receipt / ticket | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| proforma invoice | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| commercial invoice | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| packing list | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| purchase order | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| purchase order confirmation | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| bill of lading / airway bill | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| delivery note / remito | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| certificate of origin | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| quality/inspection certificate | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| customs declaration/clearance | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| insurance/freight document | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| price list/quotation | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| product specification/catalog | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| contract/addendum | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| bank/payment statement | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| warranty/service/claim document | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| identity or regulatory document | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| multi-document package | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |
| other named class | NONE_WITH_REASON | Fuera del journey de factura de referencia; requiere selección y expediente propios al ampliar el target | mantenimiento de biblioteca |

La fila commercial invoice no amplía supplier-invoice a comercio exterior; cada
clase excluida conserva su gate propio. Idioma/país/reglas tributarias no se
infieren de nombres simulados. Volumen: ensayos acotados; concurrencia probada:
dos generaciones del mismo job. No SLO ni carga de documentos privados declarados.

## Campos y revisión

Schema `supplier-invoice/v1`, una entidad por campo, sin line items ni enlaces
PO/invoice/packing list. Proyección directa de los tipos oficiales Textract:

| Campo | Campo AWS | Tipo / límite | Política |
|---|---|---|---|
| invoice_number | INVOICE_RECEIPT_ID.ValueDetection.Text | texto128bytes | REQUIRED / REVIEW_ONLY |
| vendor | VENDOR_NAME.ValueDetection.Text | texto512bytes | REQUIRED / REVIEW_ONLY |
| total | TOTAL.ValueDetection.Text | texto64bytes; no se convierte a dinero | REQUIRED / REVIEW_ONLY |
| currency | TOTAL.Currency.Code | tres letras ASCII mayúsculas | REQUIRED / REVIEW_ONLY |

No se reemplaza una ausencia por cero ni confidence por exactitud. Faltantes
quedan visibles para corrección; la propuesta exige los cuatro valores. La
revisión consulta original y respuesta oficial completa (incluidos campos de
confidence/provenance que el SDK conserve). Una propuesta es inmutable y está
ligada a original/perfil/extracción; otra persona con `documents:review` decide
sobre su SHA exacto. No hay autoapprove aunque exista umbral global positivo.
El commit, decisión y outbox se confirman en la misma transacción. Rechazar no
crea registro committed. Para corregir una propuesta rechazada se crea otro
expediente/document ID, conservando el original previo; no se reescribe historia.
Los cuatro campos siguen siendo un registro documental revisado, no un asiento
contable, pago o comprobante fiscal. Esas mutaciones pertenecen a otros owners.

## Corpus y autoridad

Un fixture público adquirido por commit/tamaño/SHA y licencia MIT original.
Ground truth de OCR:0; entrenamiento:0; test de exactitud:0. El output simulado
no está etiquetado como lectura de la imagen. Las pruebas demuestran transporte,
persistencia, aislamiento, review, recovery y contratos; no métricas de extracción.
Se selecciona AWS Go SDK Textract1.45.0 (15módulos fijados) por compatibilidad con
la composición. Jobs PostgreSQL evita un segundo runtime/in-memory workflow;
DurableTask, Stickler y CEL Java no son necesarios para este recorrido manual.
No se atribuye el glue a Microsoft/AWS/Google ni a personas famosas.

PROVIDER requiere cuenta/IAM/credenciales, región/egress/cuotas/costo admitidos y
scanners preparados. El perfil público sólo habilita supplier-invoice y
`automatic_storage=false`. Ampliar a corpus privado requiere ground truth,
política por clase/campo y jurisdicción del usuario, sin solicitarla en esta sesión.
Backup/restore, alertas, retención del supervisor y carga general se consolidan
en T2808/T2809/T2810; este documento no los promueve por inferencia.

## Evidencia y reapertura

`tools/verify_document_reference.py` exige fixture/Go exactos, DB loopback de
ensayo y receipt nuevo; rechaza SKIP. Tests cubren transacción, negativos de
scope/hash/review, replay, falla de outbox, fencing, cuarentena y cambio de perfil.
Fuzzing finito3s sobre límites de original; no sustituye SAST/carga/seguridad live.
La evidencia pública T2806 se enlaza en DOCUMENT_REFERENCE_RELEASE_V402.md/json.
Cualquier cambio de owner/perfil/input/schema/permiso/gate invalida el claim
relacionado. Cambios de corpus/modelo/país reabren su expediente antes de automatizar.
````

### FILE: `docs/DOCUMENT_REFERENCE_START.md`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5a69ff8b1d83dc12d469a3d3d676b540d07273baa9baef79bc149431b5343908"
variables: []
secrets_allowed: false
```

````markdown
# Pipeline documental de referencia

Infraestructura local/fixtures: recepción → SDK Textract → revisión → commit.
AUTHORED glue; SDK AWS original fijado; fixture Microsoft con licencia MIT exacta.
FIXTURE queda explícito en cada original, evidencia y commit. No es OCR live.

1. Aplicar las migraciones seleccionadas a la DB de referencia, incluida0085.
2. Adquirir el JPEG con `azure_document_intelligence_official_invoice/acquire_official_fixture.ps1` y su lock/approval. Puede importarse offline desde caché con `-ImportDirectory`; el acquirer verifica184686bytes/SHA. Retener `upstream/LICENSE.txt`. El template no habilita archivos privados.
3. En el host, habilitar `DOCUMENTS_ENABLED=true`, `DOCUMENTS_MODE=FIXTURE`, `DOCUMENTS_PROFILE_FILE=config/documents/reference-profile.json`, su SHA256 exacto, `DOCUMENTS_WORK_ROOT` directorio de trabajo existente, `DOCUMENTS_TENANT_ID` y `DOCUMENTS_ORGANIZATION_ID`. El IdP fixture existente entrega permisos explícitos; ninguna variable de entorno otorga roles.
4. PUT `/v1/documents/{uuid}/original`, `Content-Type: application/octet-stream`, headers `X-Document-Name`, `X-Document-SHA256`, `X-Document-Profile-SHA256` y Bearer del uploader. Body: bytes del JPEG fijado.
5. POST `/v1/documents/{uuid}/process` con `{}` y `documents:process`. GET del recurso recupera estado. GET `/original` y `/evidence/security|provider|analysis` conserva bytes/hash para revisión.
6. POST `/review` con `evidence_sha256` y `fields:{invoice_number,vendor,total,currency}`. POST `/decision` por otra persona, con `payload_sha256`, `approved` y `reason`. GET/repetición exacta recuperan resultado después de perder una respuesta.

Verificación en DB sintética loopback ya migrada:

```text
DOCUMENT_CONNECTED_DB_URL=postgres://.../elite_document_reference_<id>
python tools/verify_document_reference.py --go <Go1.26.8-admitido> --fixture <sample_invoice.jpg> --receipt <nuevo-receipt.json>
```

El runner usa DB URL sólo por entorno, nunca en el receipt ni argumentos. Fixture
público SHA `489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb`.
Registra cada prueba y falla si aparece SKIP. El proveedor fixture no abre red.

Para activar PROVIDER en tu futuro target: default credential chain AWS; perfil
por clase/region y runtime de seguridad existente mediante `DOCUMENTS_PYTHON`,
`DOCUMENTS_PYTHON_SHA256`, `DOCUMENTS_SECURITY_SCRIPT`,
`DOCUMENTS_SECURITY_SCRIPT_SHA256`, `DOCUMENTS_SECURITY_POLICY`,
`DOCUMENTS_SECURITY_POLICY_SHA256`, `DOCUMENTS_MAGIKA`, `DOCUMENTS_CLAMSCAN`,
`DOCUMENTS_CLAM_DATABASE`, `DOCUMENTS_YARA`, `DOCUMENTS_YARA_RULES`. Todos los paths
deben ser absolutos y los hashes reales. No pongas secretos en estos archivos.
Lee DOCUMENT_REFERENCE_DECISION.md antes de admitir corpus privado. El rechazo
por falta de scanner no habilita una simulación automática como fallback.
````

### FILE: `internal/approval/document_test.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3a835a5862161a22622777cd0eaf1c37e6c7c81a6255f4b472154ab9a8602b58"
variables: []
secrets_allowed: false
```

````go
package approval

import (
	"strings"
	"testing"
)

func TestDocumentReviewNeverAutoApproves(t *testing.T) {
	r := NewRegistry(Policy{AutoApproveMinorUnits: 10000})
	request := Request{TenantID: "fixture-tenant", ID: "document:fixture", Kind: KindDocumentReview, SubjectID: "fixture-document", Requester: "maker", EvidenceSHA: strings.Repeat("a", 64)}
	state, e := r.Submit(request)
	if e != nil || state != StatePending {
		t.Fatal("document automatic approval", state, e)
	}
	if _, e = r.Approve(request.TenantID, request.ID, "maker", "self"); e != ErrSeparation {
		t.Fatal("self approval", e)
	}
	state, e = r.Approve(request.TenantID, request.ID, "reviewer", "reviewed fields")
	if e != nil || state != StateApproved {
		t.Fatal(state, e)
	}
}
````

### FILE: `internal/documentbridge/fixture.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fbb8cee6f4c3517e64b57486353f15c42ae26e054bbbaea6e168d50bf32ced44"
variables: []
secrets_allowed: false
```

````go
package documentbridge

// AUTHORED simulator. No detector binaries or live OCR run here. FIXTURE mode is
// permanent on original/evidence/commit records, restricted to one public hash.
import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const FixtureResponse = `{"DocumentMetadata":{"Pages":1},"ExpenseDocuments":[{"ExpenseIndex":1,"SummaryFields":[{"Type":{"Text":"INVOICE_RECEIPT_ID"},"ValueDetection":{"Text":"FIXTURE-INVOICE-1"}},{"Type":{"Text":"VENDOR_NAME"},"ValueDetection":{"Text":"Fixture vendor"}},{"Type":{"Text":"TOTAL"},"ValueDetection":{"Text":"123.45"},"Currency":{"Code":"USD"}}],"LineItemGroups":[]}]}`

type fixtureHTTP struct{}

func (fixtureHTTP) Do(r *http.Request) (*http.Response, error) {
	if r.Method != "POST" || r.Header.Get("X-Amz-Target") != "Textract.AnalyzeExpense" || !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 Credential=AKIDFIXTURE/") {
		return nil, ErrContract
	}
	var wire struct{ Document struct{ Bytes []byte } }
	b, e := io.ReadAll(io.LimitReader(r.Body, 3*1024*1024))
	if e != nil || json.Unmarshal(b, &wire) != nil || Hash(wire.Document.Bytes) != PublicFixtureSHA {
		return nil, ErrContract
	}
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/x-amz-json-1.1"}, "X-Amzn-Requestid": []string{"simulated-public-document-only"}}, Body: io.NopCloser(strings.NewReader(FixtureResponse)), Request: r}, nil
}
func FixtureSecurity(ctx context.Context, input, output string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	b, e := regular(input, MaxBytes)
	if e != nil || Hash(b) != PublicFixtureSHA {
		return ErrSecurity
	}
	synthetic := Hash([]byte("SIMULATED DETECTOR CONTRACT; NO SCANNER EXECUTED"))
	now := time.Now().UTC().Format(time.RFC3339Nano)
	receipt := map[string]any{
		"schema": "elite-secure-local-file-receipt/v1", "decision": "ADMITTED", "reason": "ALL_SELECTED_GATES_PASSED", "started_at": now, "completed_at": now,
		"approval_id": "user-authorized-public-fixture-only", "policy_sha256": synthetic, "input": map[string]any{"sha256": PublicFixtureSHA, "bytes": len(b)}, "business_storage_authorized": false,
		"security_claim": "SIMULATED detector contract for exact public fixture; no scanner executed; no malware or private-file admission proof",
		"clamav":         map[string]any{"exit_code": 0, "output_sha256": synthetic}, "content_type": map[string]any{"label": "jpeg", "mime_type": "image/jpeg", "score": 1.0},
		"yara_x": map[string]any{"compiled_rules_sha256": synthetic, "matching_rules": []string{}},
		"tools":  map[string]any{"clamscan": map[string]any{"sha256": synthetic, "version": "ClamAV 1.5.4/SIMULATED"}, "magika": map[string]any{"sha256": synthetic, "version": "magika 1.1.0"}, "yara_x": map[string]any{"sha256": synthetic, "version": "1.20.0"}},
	}
	raw, e := json.Marshal(receipt)
	if e != nil {
		return e
	}
	if e = os.Mkdir(output, 0700); e != nil {
		return e
	}
	return os.WriteFile(filepath.Join(output, "security-receipt.json"), raw, 0600)
}
func NewFixturePipeline(profile, sha, root string) (*Pipeline, error) {
	client := textract.NewFromConfig(aws.Config{Region: "us-east-1", Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
		return aws.Credentials{AccessKeyID: "AKIDFIXTURE", SecretAccessKey: "fixture-only-no-live-authority"}, nil
	}), HTTPClient: fixtureHTTP{}}, func(o *textract.Options) { o.RetryMaxAttempts = 1 })
	return NewPipeline(client, FixtureSecurity, profile, sha, root, "FIXTURE")
}
````

### FILE: `internal/documentbridge/pipeline.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ee982ae84a9c5f14fb7d47ba5b109a0c0220e645895079bf6122d7cbdb29bf8b"
variables: []
secrets_allowed: false
```

````go
// Package documentbridge is AUTHORED bounded orchestration glue around the
// original secure-file gate and pinned official AWS Textract SDK adapter.
package documentbridge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

	runtime "example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

const MaxBytes = 2 * 1024 * 1024
const PublicFixtureSHA = "489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb"
const Class = "supplier-invoice"

var ErrContract = errors.New("document contract rejected")
var ErrSecurity = errors.New("document security gate rejected")
var ErrExtraction = errors.New("document extraction unavailable")

type Scope struct {
	TenantID       string `json:"tenant_id"`
	OrganizationID string `json:"organization_id"`
	ProfileSHA     string `json:"profile_sha256"`
	Mode           string `json:"mode"`
}
type Original struct {
	Name   string
	SHA256 string
	Bytes  []byte
}
type Fields struct {
	InvoiceNumber string `json:"invoice_number"`
	Vendor        string `json:"vendor"`
	Total         string `json:"total"`
	Currency      string `json:"currency"`
}
type Evidence struct {
	Security, Provider, Receipt []byte
	Suggested                   Fields
	Mode                        string
}
type Processor interface {
	Run(context.Context, Original) (Evidence, error)
}
type SecurityStage func(context.Context, string, string) error

func Hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func Text(v string, n int) bool {
	return v != "" && len(v) <= n && utf8.ValidString(v) && !strings.ContainsRune(v, 0) && strings.TrimSpace(v) == v
}
func Hex(v string) bool {
	b, e := hex.DecodeString(v)
	return e == nil && len(b) == 32 && strings.ToLower(v) == v
}
func (v Fields) Validate() error {
	// Values remain reviewed text, never amounts to post, tax calculations or OCR truth.
	if !Text(v.InvoiceNumber, 128) || !Text(v.Vendor, 512) || !Text(v.Total, 64) || len(v.Currency) != 3 {
		return ErrContract
	}
	for _, r := range v.Currency {
		if r < 'A' || r > 'Z' {
			return ErrContract
		}
	}
	return nil
}
func (v Original) Validate(mode string) error {
	if !Text(v.Name, 128) || strings.ContainsAny(v.Name, "/\\:") || v.Name == "." || v.Name == ".." || len(v.Bytes) == 0 || len(v.Bytes) > MaxBytes || Hash(v.Bytes) != v.SHA256 {
		return ErrContract
	}
	ext := strings.ToLower(filepath.Ext(v.Name))
	mime := http.DetectContentType(v.Bytes)
	if !(ext == ".pdf" && mime == "application/pdf" || (ext == ".jpg" || ext == ".jpeg") && mime == "image/jpeg") {
		return ErrContract
	}
	if mode != "PROVIDER" && mode != "FIXTURE" || mode == "FIXTURE" && v.SHA256 != PublicFixtureSHA {
		return ErrContract
	}
	return nil
}

type Pipeline struct {
	client                                  runtime.Client
	security                                SecurityStage
	profile, profileSHA, root, region, mode string
}

func NewPipeline(client runtime.Client, security SecurityStage, profile, expectedSHA, root, mode string) (*Pipeline, error) {
	if client == nil || security == nil || !Hex(expectedSHA) || (mode != "PROVIDER" && mode != "FIXTURE") {
		return nil, ErrContract
	}
	b, e := regular(profile, 32768)
	if e != nil || Hash(b) != expectedSHA {
		return nil, ErrContract
	}
	region, e := runtime.ApprovedRegion(profile, Class)
	if e != nil {
		return nil, ErrContract
	}
	info, e := os.Lstat(root)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, ErrContract
	}
	return &Pipeline{client, security, profile, expectedSHA, root, region, mode}, nil
}
func regular(path string, max int64) ([]byte, error) {
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Size() > max {
		return nil, ErrContract
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, max+1))
	if e != nil || int64(len(b)) > max {
		return nil, ErrContract
	}
	return b, nil
}
func (p *Pipeline) Run(ctx context.Context, v Original) (Evidence, error) {
	if p == nil || v.Validate(p.mode) != nil {
		return Evidence{}, ErrContract
	}
	profile, e := regular(p.profile, 32768)
	if e != nil || Hash(profile) != p.profileSHA {
		return Evidence{}, ErrContract
	}
	dir, e := os.MkdirTemp(p.root, "document-")
	if e != nil {
		return Evidence{}, e
	}
	// The caller persists all successful evidence in PostgreSQL. Temporary originals
	// are removed on return; interrupted attempts require host orphan cleanup policy.
	defer os.RemoveAll(dir)
	input := filepath.Join(dir, "original"+strings.ToLower(filepath.Ext(v.Name)))
	if e = os.WriteFile(input, v.Bytes, 0600); e != nil {
		return Evidence{}, e
	}
	securityDir := filepath.Join(dir, "security")
	if e = p.security(ctx, input, securityDir); e != nil {
		return Evidence{}, ErrSecurity
	}
	security, e := regular(filepath.Join(securityDir, "security-receipt.json"), 65536)
	if e != nil {
		return Evidence{}, ErrSecurity
	}
	receipt, e := runtime.AnalyzeToEvidence(ctx, p.client, runtime.Request{Region: p.region, DocumentClass: Class, InputPath: input, OutputDirectory: filepath.Join(dir, "extraction"), ProfilePath: p.profile, SecurityReceiptPath: filepath.Join(securityDir, "security-receipt.json"), MaxBytes: MaxBytes})
	if e != nil {
		return Evidence{}, ErrExtraction
	}
	raw, e := regular(filepath.Join(dir, "extraction", "provider-response.json"), 4*1024*1024)
	if e != nil {
		return Evidence{}, ErrContract
	}
	original, e := regular(input, MaxBytes)
	if e != nil || !bytes.Equal(original, v.Bytes) || receipt.InputSHA256 != v.SHA256 || receipt.DocumentProfileSHA256 != p.profileSHA || receipt.SecurityReceiptSHA256 != Hash(security) || receipt.ProviderResponseSHA256 != Hash(raw) || receipt.AutomaticStorageAuthorized {
		return Evidence{}, ErrContract
	}
	var output textract.AnalyzeExpenseOutput
	if json.Unmarshal(raw, &output) != nil || len(output.ExpenseDocuments) != 1 {
		return Evidence{}, ErrContract
	}
	fields := Fields{}
	seen := map[string]bool{}
	for _, f := range output.ExpenseDocuments[0].SummaryFields {
		if f.Type == nil || f.ValueDetection == nil {
			continue
		}
		key := aws.ToString(f.Type.Text)
		if key != "INVOICE_RECEIPT_ID" && key != "VENDOR_NAME" && key != "TOTAL" {
			continue
		}
		if seen[key] {
			return Evidence{}, ErrContract
		}
		seen[key] = true
		switch key {
		case "INVOICE_RECEIPT_ID":
			fields.InvoiceNumber = aws.ToString(f.ValueDetection.Text)
		case "VENDOR_NAME":
			fields.Vendor = aws.ToString(f.ValueDetection.Text)
		case "TOTAL":
			fields.Total = aws.ToString(f.ValueDetection.Text)
			if f.Currency != nil {
				fields.Currency = aws.ToString(f.Currency.Code)
			}
		}
	}
	// Missing extraction values are preserved for manual correction; validation
	// of all four required fields belongs to the immutable review proposal.
	if len(fields.InvoiceNumber) > 128 || len(fields.Vendor) > 512 || len(fields.Total) > 64 || len(fields.Currency) > 3 {
		return Evidence{}, ErrContract
	}
	receiptBytes, e := json.Marshal(receipt)
	if e != nil {
		return Evidence{}, e
	}
	return Evidence{security, raw, receiptBytes, fields, p.mode}, nil
}

type SecurityCommand struct {
	Python, PythonSHA, Script, ScriptSHA, Policy, PolicySHA string
	Magika, ClamScan, ClamDatabase, Yara, Rules             string
}

func (c SecurityCommand) Stage() (SecurityStage, error) {
	for _, v := range []struct{ path, hash string }{{c.Python, c.PythonSHA}, {c.Script, c.ScriptSHA}, {c.Policy, c.PolicySHA}} {
		if !filepath.IsAbs(v.path) || !Hex(v.hash) {
			return nil, ErrContract
		}
		b, e := regular(v.path, 64*1024*1024)
		if e != nil || Hash(b) != v.hash {
			return nil, ErrContract
		}
	}
	for _, p := range []string{c.Magika, c.ClamScan, c.ClamDatabase, c.Yara, c.Rules} {
		if !filepath.IsAbs(p) {
			return nil, ErrContract
		}
	}
	return func(ctx context.Context, input, output string) error {
		// Recheck executable/script/policy hashes per attempt. Native binary hashes,
		// official database age, YARA rule hash and detector limits remain gate-owned.
		for _, v := range []struct{ path, hash string }{{c.Python, c.PythonSHA}, {c.Script, c.ScriptSHA}, {c.Policy, c.PolicySHA}} {
			b, e := regular(v.path, 64*1024*1024)
			if e != nil || Hash(b) != v.hash {
				return ErrSecurity
			}
		}
		cmd := exec.CommandContext(ctx, c.Python, "-I", "-X", "utf8", "-B", c.Script, "--input", input, "--output", output, "--policy", c.Policy, "--magika-exe", c.Magika, "--clamscan-exe", c.ClamScan, "--clamav-database", c.ClamDatabase, "--yara-exe", c.Yara, "--compiled-rules", c.Rules)
		cmd.Env = []string{"PYTHONIOENCODING=utf-8"}
		for _, k := range []string{"SystemRoot", "TEMP", "TMP"} {
			if v := os.Getenv(k); v != "" {
				cmd.Env = append(cmd.Env, k+"="+v)
			}
		}
		// No raw detector output/private paths are copied to application logs.
		if e := cmd.Run(); e != nil {
			return fmt.Errorf("%w: detector process", ErrSecurity)
		}
		return nil
	}, nil
}
````

### FILE: `internal/documentbridge/pipeline_test.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b67d05933a1d882971ba1eeff382f137b289cd7e9db673acf7db3c5b9e405152"
variables: []
secrets_allowed: false
```

````go
package documentbridge

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentIntakeBoundsAndFixtureIsolation(t *testing.T) {
	b := []byte("%PDF-1.7 fixture")
	o := Original{Name: "invoice.pdf", SHA256: Hash(b), Bytes: b}
	if o.Validate("PROVIDER") != nil {
		t.Fatal("reference structural PDF rejected")
	}
	if o.Validate("FIXTURE") == nil {
		t.Fatal("arbitrary bytes reached public fixture lane")
	}
	for _, name := range []string{"../invoice.pdf", "C:invoice.pdf", "folder\\invoice.pdf", "invoice.jpg", "invoice.pdf ", strings.Repeat("a", 129) + ".pdf"} {
		v := o
		v.Name = name
		if v.Validate("PROVIDER") == nil {
			t.Fatal("unbounded filename or mismatched type", name)
		}
	}
	o.SHA256 = strings.Repeat("0", 64)
	if o.Validate("PROVIDER") == nil {
		t.Fatal("unbound bytes")
	}
	if FixtureSecurity(context.Background(), "absent", filepath.Join(t.TempDir(), "receipt")) == nil {
		t.Fatal("fixture security accepted missing source")
	}
	if _, e := (SecurityCommand{Python: os.Args[0], PythonSHA: strings.Repeat("0", 64)}).Stage(); e == nil {
		t.Fatal("unbound security process selected")
	}
}
func FuzzDocumentOriginalBoundary(f *testing.F) {
	f.Add("invoice.pdf", []byte("%PDF-1.7 reference"))
	f.Add("../invoice.jpg", []byte{255, 216, 255, 224})
	f.Add("file.zip", []byte("PK"))
	f.Fuzz(func(t *testing.T, name string, b []byte) {
		if len(b) > MaxBytes+1 {
			return
		}
		o := Original{name, Hash(b), b}
		e := o.Validate("PROVIDER")
		if e == nil {
			if len(b) == 0 || len(b) > MaxBytes || strings.ContainsAny(name, "/\\:") {
				t.Fatal("accepted path or byte-budget violation")
			}
			mutated := append([]byte(nil), b...)
			mutated[0] ^= 1
			o.Bytes = mutated
			if o.Validate("PROVIDER") == nil {
				t.Fatal("content mutation preserved authority")
			}
		}
	})
}
````

### FILE: `internal/platform/httpapi/document.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "292a9c94de0567361310b5ff57826b9f4e5477ff027e42759308a5b2d86e995e"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED role/scoping and bounded HTTP glue. No user-supplied evidence receipt.
import (
	"bytes"
	"elite.local/enterprise/internal/approval"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type DocumentModule struct{ Store *postgres.Documents }

func documentJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func documentError(w http.ResponseWriter, e error) {
	status, code := 503, "UNAVAILABLE"
	switch {
	case errors.Is(e, postgres.ErrDocumentScope):
		status, code = 404, "NOT_FOUND"
	case errors.Is(e, doc.ErrContract):
		status, code = 400, "INVALID_CONTRACT"
	case errors.Is(e, doc.ErrSecurity):
		status, code = 422, "QUARANTINED"
	case errors.Is(e, postgres.ErrDocumentConflict), errors.Is(e, approval.ErrDuplicate), errors.Is(e, approval.ErrNotPending), errors.Is(e, approval.ErrSeparation):
		status, code = 409, "CONSULT_RECORDED_STATE"
	}
	documentJSON(w, status, map[string]string{"code": code})
}
func (m DocumentModule) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier, permission string) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		documentJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return identity.Principal{}, false
	}
	p, e := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if e != nil {
		documentJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return p, false
	}
	allowed := m.Store != nil && (m.Store.Allowed(p, permission) || permission == "documents:read" && (m.Store.Allowed(p, "documents:write") || m.Store.Allowed(p, "documents:review") || m.Store.Allowed(p, "documents:process")))
	if !allowed {
		documentJSON(w, 403, map[string]string{"code": "FORBIDDEN"})
		return p, false
	}
	if r.URL.RawQuery != "" {
		documentJSON(w, 400, map[string]string{"code": "INVALID_QUERY"})
		return p, false
	}
	return p, true
}
func documentBody(w http.ResponseWriter, r *http.Request, out any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return doc.ErrContract
	}
	b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		return doc.ErrContract
	}
	canonical, _, e := approval.CanonicalPayload(b)
	if e != nil {
		return doc.ErrContract
	}
	d := json.NewDecoder(bytes.NewReader(canonical))
	d.DisallowUnknownFields()
	if d.Decode(out) != nil {
		return doc.ErrContract
	}
	return nil
}
func (m DocumentModule) Register(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("GET /v1/documents/{id}/evidence/{part}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		b, e := m.Store.EvidencePart(r.Context(), p, r.PathValue("id"), r.PathValue("part"))
		if e != nil {
			documentError(w, e)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\"evidence.json\"")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Document-Evidence-SHA256", doc.Hash(b))
		_, _ = w.Write(b)
	})

	mux.HandleFunc("PUT /v1/documents/{id}/original", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:write")
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			documentError(w, doc.ErrContract)
			return
		}
		for _, h := range []string{"X-Document-Name", "X-Document-SHA256", "X-Document-Profile-SHA256"} {
			if len(r.Header.Values(h)) != 1 {
				documentError(w, doc.ErrContract)
				return
			}
		}
		b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, doc.MaxBytes))
		if e != nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Receive(r.Context(), p, r.PathValue("id"), r.Header.Get("X-Document-Profile-SHA256"), doc.Original{Name: r.Header.Get("X-Document-Name"), SHA256: r.Header.Get("X-Document-SHA256"), Bytes: b})
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 201, result)
	})
	mux.HandleFunc("GET /v1/documents/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		result, e := m.Store.Read(r.Context(), p, r.PathValue("id"))
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("GET /v1/documents/{id}/original", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:read")
		if !ok {
			return
		}
		result, e := m.Store.Original(r.Context(), p, r.PathValue("id"))
		if e != nil {
			documentError(w, e)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"original.bin\"")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Length", strconv.Itoa(len(result.Bytes)))
		w.Header().Set("X-Document-SHA256", result.SHA256)
		_, _ = w.Write(result.Bytes)
	})
	mux.HandleFunc("POST /v1/documents/{id}/process", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:process")
		if !ok {
			return
		}
		var b struct{}
		if documentBody(w, r, &b) != nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Process(r.Context(), p, r.PathValue("id"))
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("POST /v1/documents/{id}/review", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:write")
		if !ok {
			return
		}
		var b struct {
			EvidenceSHA string     `json:"evidence_sha256"`
			Fields      doc.Fields `json:"fields"`
		}
		if documentBody(w, r, &b) != nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Submit(r.Context(), p, r.PathValue("id"), b.EvidenceSHA, b.Fields)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
	mux.HandleFunc("POST /v1/documents/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "documents:review")
		if !ok {
			return
		}
		var b struct {
			SHA      string `json:"payload_sha256"`
			Approved *bool  `json:"approved"`
			Reason   string `json:"reason"`
		}
		if documentBody(w, r, &b) != nil || b.Approved == nil {
			documentError(w, doc.ErrContract)
			return
		}
		result, e := m.Store.Decide(r.Context(), p, r.PathValue("id"), b.SHA, *b.Approved, b.Reason)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, 200, result)
	})
}
````

### FILE: `internal/platform/postgres/document.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "94ca7f9d284a379368fdbdd22a5be88fc8e7a7b95003471efc6362cafa54684d"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED same-transaction composition of existing job fencing, human approvals
// and immutable original/extraction records. No automatic storage authority.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	runtime "example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var ErrDocumentScope = errors.New("document not found or forbidden")
var ErrDocumentConflict = errors.New("document conflict; consult recorded state")

type Documents struct {
	pool      *pgxpool.Pool
	scope     doc.Scope
	processor doc.Processor
	reviews   *HumanApprovals
}
type DocumentProposal struct {
	Schema         string     `json:"schema"`
	DocumentID     string     `json:"document_id"`
	OrganizationID string     `json:"organization_id"`
	OriginalSHA    string     `json:"original_sha256"`
	ProfileSHA     string     `json:"profile_sha256"`
	EvidenceSHA    string     `json:"evidence_sha256"`
	Mode           string     `json:"mode"`
	Fields         doc.Fields `json:"fields"`
}
type DocumentView struct {
	ID          string            `json:"document_id"`
	Name        string            `json:"name"`
	Uploader    string            `json:"uploader"`
	OriginalSHA string            `json:"original_sha256"`
	ProfileSHA  string            `json:"profile_sha256"`
	Mode        string            `json:"mode"`
	State       string            `json:"state"`
	EvidenceSHA string            `json:"evidence_sha256,omitempty"`
	Suggested   *doc.Fields       `json:"suggested,omitempty"`
	Proposal    *DocumentProposal `json:"proposal,omitempty"`
	PayloadSHA  string            `json:"payload_sha256,omitempty"`
	Reviewer    string            `json:"reviewer,omitempty"`
	Reason      string            `json:"reason,omitempty"`
}

func NewDocuments(pool *pgxpool.Pool, s doc.Scope, p doc.Processor) (*Documents, error) {
	if pool == nil || p == nil || !documentID(s.TenantID) || !doc.Text(s.OrganizationID, 128) || !doc.Hex(s.ProfileSHA) || (s.Mode != "FIXTURE" && s.Mode != "PROVIDER") {
		return nil, doc.ErrContract
	}
	return &Documents{pool, s, p, NewHumanApprovals(pool)}, nil
}
func documentID(v string) bool {
	x, e := uuid.Parse(v)
	return e == nil && x != uuid.Nil && x.String() == v
}
func (d *Documents) Allowed(p identity.Principal, permission string) bool {
	return d != nil && p.TenantID == d.scope.TenantID && doc.Text(p.Subject, 128) && p.AllowedOrganization(d.scope.OrganizationID) && p.Allowed(permission)
}
func (d *Documents) canRead(p identity.Principal) bool {
	return d.Allowed(p, "documents:write") || d.Allowed(p, "documents:review") || d.Allowed(p, "documents:process")
}
func (d *Documents) requestID(id string) string { return "document:" + id }
func (d *Documents) jobPayload(id string) []byte {
	b, _ := json.Marshal(map[string]string{"document_id": id, "organization_id": d.scope.OrganizationID, "profile_sha256": d.scope.ProfileSHA})
	return b
}
func documentJobID(tenant, org, profile, id string) string {
	payload, _ := json.Marshal(map[string]string{"document_id": id, "organization_id": org, "profile_sha256": profile})
	return uuid.NewSHA1(uuid.NameSpaceOID, append([]byte("document-extraction:"+tenant+":"), payload...)).String()
}
func (d *Documents) jobID(id string) string {
	return documentJobID(d.scope.TenantID, d.scope.OrganizationID, d.scope.ProfileSHA, id)
}
func (d *Documents) Receive(ctx context.Context, p identity.Principal, id, profileSHA string, v doc.Original) (DocumentView, error) {
	if !d.Allowed(p, "documents:write") || !documentID(id) {
		return DocumentView{}, ErrDocumentScope
	}
	if profileSHA != d.scope.ProfileSHA || v.Validate(d.scope.Mode) != nil {
		return DocumentView{}, doc.ErrContract
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return DocumentView{}, e
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, `insert into document.original(tenant_id,document_id,organization_id,uploader,name,original_sha256,profile_sha256,mode,content)values($1,$2,$3,$4,$5,$6,$7,$8,$9)on conflict do nothing`, p.TenantID, id, d.scope.OrganizationID, p.Subject, v.Name, v.SHA256, profileSHA, d.scope.Mode, v.Bytes)
	if e != nil {
		return DocumentView{}, e
	}
	var same bool
	e = tx.QueryRow(ctx, `select organization_id=$3 and uploader=$4 and name=$5 and original_sha256=$6 and profile_sha256=$7 and mode=$8 and content=$9 from document.original where tenant_id=$1 and document_id=$2`, p.TenantID, id, d.scope.OrganizationID, p.Subject, v.Name, v.SHA256, profileSHA, d.scope.Mode, v.Bytes).Scan(&same)
	if e != nil || !same {
		return DocumentView{}, ErrDocumentConflict
	}
	_, e = tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts)values($1,$2,'document-extraction','document-extraction',1,$3,3)on conflict do nothing`, p.TenantID, d.jobID(id), d.jobPayload(id))
	if e != nil {
		return DocumentView{}, e
	}
	if e = tx.Commit(ctx); e != nil {
		return DocumentView{}, e
	}
	return d.Read(ctx, p, id)
}
func (d *Documents) Read(ctx context.Context, p identity.Principal, id string) (DocumentView, error) {
	v := DocumentView{}
	if !d.canRead(p) || !documentID(id) {
		return v, ErrDocumentScope
	}
	e := d.pool.QueryRow(ctx, `select document_id::text,name,uploader,original_sha256,profile_sha256,mode from document.original where tenant_id=$1 and document_id=$2 and organization_id=$3`, p.TenantID, id, d.scope.OrganizationID).Scan(&v.ID, &v.Name, &v.Uploader, &v.OriginalSHA, &v.ProfileSHA, &v.Mode)
	if errors.Is(e, pgx.ErrNoRows) {
		return v, ErrDocumentScope
	}
	if e != nil {
		return v, e
	}
	if v.Uploader != p.Subject && !d.Allowed(p, "documents:review") && !d.Allowed(p, "documents:process") {
		return DocumentView{}, ErrDocumentScope
	}
	v.State = "QUARANTINED"
	var raw []byte
	e = d.pool.QueryRow(ctx, `select evidence_sha256,suggested from document.extraction where tenant_id=$1 and document_id=$2`, p.TenantID, id).Scan(&v.EvidenceSHA, &raw)
	if e == nil {
		var f doc.Fields
		if json.Unmarshal(raw, &f) != nil {
			return v, doc.ErrContract
		}
		v.Suggested = &f
		v.State = "REVIEW_REQUIRED"
	} else if !errors.Is(e, pgx.ErrNoRows) {
		return v, e
	} else {
		var terminal *string
		e = d.pool.QueryRow(ctx, `select terminal_error_code from platform.job where tenant_id=$1 and job_id=$2`, p.TenantID, documentJobID(p.TenantID, d.scope.OrganizationID, v.ProfileSHA, id)).Scan(&terminal)
		if e != nil {
			return v, e
		}
		if terminal != nil {
			v.State = "QUARANTINE_TERMINAL"
		}
		return v, nil
	}
	var state approval.State
	e = d.pool.QueryRow(ctx, `select evidence_sha,payload,state from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='document_review'`, p.TenantID, d.requestID(id), d.scope.OrganizationID).Scan(&v.PayloadSHA, &raw, &state)
	if errors.Is(e, pgx.ErrNoRows) {
		return v, nil
	}
	if e != nil {
		return v, e
	}
	var proposal DocumentProposal
	_, h, e := approval.CanonicalPayload(raw)
	if e != nil || h != v.PayloadSHA || json.Unmarshal(raw, &proposal) != nil || proposal.DocumentID != id || proposal.EvidenceSHA != v.EvidenceSHA || proposal.OriginalSHA != v.OriginalSHA || proposal.ProfileSHA != v.ProfileSHA {
		return v, doc.ErrContract
	}
	v.Proposal = &proposal
	v.State = "REVIEW_PENDING"
	if state != approval.StatePending {
		var approved bool
		e = d.pool.QueryRow(ctx, `select reviewer,approved,reason from approval.decision where tenant_id=$1 and request_id=$2 and(select count(*)from approval.decision where tenant_id=$1 and request_id=$2)=1`, p.TenantID, d.requestID(id)).Scan(&v.Reviewer, &approved, &v.Reason)
		if e != nil {
			return v, e
		}
		if approved != (state == approval.StateApproved) {
			return v, doc.ErrContract
		}
		v.State = "REJECTED"
		if approved {
			var bound bool
			e = d.pool.QueryRow(ctx, `select request_id=$3 and payload_sha256=$4 and reviewer=$5 from document.committed where tenant_id=$1 and document_id=$2`, p.TenantID, id, d.requestID(id), v.PayloadSHA, v.Reviewer).Scan(&bound)
			if e != nil || !bound {
				return v, doc.ErrContract
			}
			v.State = "PERSISTED"
		}
	}
	return v, nil
}
func (d *Documents) Original(ctx context.Context, p identity.Principal, id string) (doc.Original, error) {
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return doc.Original{}, e
	}
	original := doc.Original{Name: v.Name, SHA256: v.OriginalSHA}
	e = d.pool.QueryRow(ctx, `select content from document.original where tenant_id=$1 and document_id=$2 and organization_id=$3`, p.TenantID, id, d.scope.OrganizationID).Scan(&original.Bytes)
	if e != nil {
		return original, e
	}
	if original.Validate(v.Mode) != nil {
		return doc.Original{}, doc.ErrContract
	}
	return original, nil
}
func (d *Documents) Process(ctx context.Context, p identity.Principal, id string) (DocumentView, error) {
	if !d.Allowed(p, "documents:process") {
		return DocumentView{}, ErrDocumentScope
	}
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.ProfileSHA != d.scope.ProfileSHA || v.Mode != d.scope.Mode {
		return v, ErrDocumentConflict
	}
	if v.EvidenceSHA != "" {
		return v, nil
	}
	// Existing claim generation/DB clock owns concurrency and stale-worker fencing.
	jobs := NewJobs(d.pool)
	worker := uuid.NewString()
	claimed, e := jobs.ClaimScoped(ctx, "document-extraction", worker, 3*time.Minute, 1, JobScope{d.scope.TenantID, "document-extraction", 1, d.jobPayload(id)})
	if e != nil {
		return v, e
	}
	if len(claimed) != 1 {
		// Preserve a crashed final attempt as terminal evidence, using the
		// existing owner fence. No counter reset or new provider call.
		var exhausted Job
		err := d.pool.QueryRow(ctx, `select tenant_id::text,job_id::text,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job where tenant_id=$1 and job_id=$2 and completed_at is null and terminal_error_code is null and attempts>=max_attempts and claimed_until<clock_timestamp()`, p.TenantID, d.jobID(id)).Scan(&exhausted.TenantID, &exhausted.JobID, &exhausted.Queue, &exhausted.JobType, &exhausted.SchemaVersion, &exhausted.Payload, &exhausted.Attempts, &exhausted.MaxAttempts)
		if err == nil {
			tx, err := d.pool.Begin(ctx)
			if err != nil {
				return v, err
			}
			defer tx.Rollback(ctx)
			if err = ExhaustJobInTx(ctx, tx, exhausted); err != nil {
				return v, err
			}
			if err = tx.Commit(ctx); err != nil {
				return v, err
			}
			return d.Read(ctx, p, id)
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return v, err
		}
		return v, ErrDocumentConflict
	}
	job := claimed[0]
	original, e := d.Original(ctx, p, id)
	if e != nil {
		return v, e
	}
	callCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	evidence, e := d.processor.Run(callCtx, original)
	cancel()
	// If the caller disconnected the finite lease recovers on the next explicit
	// worker invocation. No network error can fabricate an extraction receipt.
	if e != nil {
		return v, d.recordDocumentFailure(ctx, job, worker, id, e)
	}
	if e = validateDocumentEvidence(evidence, original, v.ProfileSHA, v.Mode); e != nil {
		return v, d.recordDocumentFailure(ctx, job, worker, id, e)
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return v, e
	}
	defer tx.Rollback(ctx)
	if e = CompleteJobInTx(ctx, tx, job, worker); e != nil {
		return v, e
	}
	fields, _ := json.Marshal(evidence.Suggested)
	_, e = tx.Exec(ctx, `insert into document.extraction(tenant_id,document_id,job_id,attempt,security_receipt,provider_response,analysis_receipt,evidence_sha256,suggested)values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, p.TenantID, id, job.JobID, job.Attempts, evidence.Security, evidence.Provider, evidence.Receipt, doc.Hash(evidence.Receipt), fields)
	if e != nil {
		return v, e
	}
	if e = tx.Commit(ctx); e != nil {
		return v, e
	}
	return d.Read(ctx, p, id)
}
func validateDocumentEvidence(e doc.Evidence, o doc.Original, profile, mode string) error {
	var r runtime.Receipt
	if e.Mode != mode || len(e.Security) == 0 || len(e.Security) > 65536 || len(e.Provider) == 0 || len(e.Provider) > 4194304 || len(e.Receipt) > 32768 || json.Unmarshal(e.Receipt, &r) != nil || r.AutomaticStorageAuthorized || r.InputSHA256 != o.SHA256 || r.InputBytes != len(o.Bytes) || r.DocumentProfileSHA256 != profile || r.SecurityReceiptSHA256 != doc.Hash(e.Security) || r.ProviderResponseSHA256 != doc.Hash(e.Provider) || r.PageCount != 1 || r.DocumentClass != doc.Class || r.SDKVersion != "v1.45.0" {
		return doc.ErrContract
	}
	return nil
}
func (d *Documents) recordDocumentFailure(ctx context.Context, job Job, worker, id string, cause error) error {
	code := "EXTRACTION_UNAVAILABLE"
	if errors.Is(cause, doc.ErrSecurity) {
		code = "SECURITY_REJECTED"
	}
	if errors.Is(cause, doc.ErrContract) {
		code = "INVALID_CONTRACT"
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	_, e = FailJobInTx(ctx, tx, job, worker, code, time.Second)
	if e != nil {
		return e
	}
	_, e = tx.Exec(ctx, `insert into document.attempt_failure(tenant_id,document_id,attempt,code)values($1,$2,$3,$4)`, job.TenantID, id, job.Attempts, code)
	if e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	return cause
}
func (d *Documents) Submit(ctx context.Context, p identity.Principal, id, evidenceSHA string, fields doc.Fields) (DocumentView, error) {
	if !d.Allowed(p, "documents:write") {
		return DocumentView{}, ErrDocumentScope
	}
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.Uploader != p.Subject || v.ProfileSHA != d.scope.ProfileSHA || v.EvidenceSHA == "" || v.EvidenceSHA != evidenceSHA || fields.Validate() != nil {
		return v, ErrDocumentConflict
	}
	proposal := DocumentProposal{"document-review/v1", id, d.scope.OrganizationID, v.OriginalSHA, v.ProfileSHA, v.EvidenceSHA, v.Mode, fields}
	raw, _ := json.Marshal(proposal)
	canonical, h, e := approval.CanonicalPayload(raw)
	if e != nil {
		return v, e
	}
	_, e = d.reviews.Submit(ctx, p, HumanApprovalSpec{approval.Request{TenantID: p.TenantID, ID: d.requestID(id), Kind: approval.KindDocumentReview, SubjectID: id, Requester: p.Subject, EvidenceSHA: h}, d.scope.OrganizationID, canonical}, "documents:write", nil)
	if e != nil {
		return v, e
	}
	return d.Read(ctx, p, id)
}
func (d *Documents) Decide(ctx context.Context, p identity.Principal, id, expectedSHA string, approved bool, reason string) (DocumentView, error) {
	if !d.Allowed(p, "documents:review") || !doc.Text(reason, 2048) {
		return DocumentView{}, ErrDocumentScope
	}
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.Proposal == nil || v.ProfileSHA != d.scope.ProfileSHA || v.PayloadSHA != expectedSHA {
		return v, ErrDocumentConflict
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return v, e
	}
	defer tx.Rollback(ctx)
	_, e = d.reviews.decideTx(ctx, tx, p, p.TenantID, d.requestID(id), d.scope.OrganizationID, expectedSHA, approved, reason, "documents:review", func(ctx context.Context, tx pgx.Tx) error {
		if !approved {
			return nil
		}
		f, _ := json.Marshal(v.Proposal.Fields)
		_, err := tx.Exec(ctx, `insert into document.committed(tenant_id,document_id,request_id,payload_sha256,fields,reviewer)values($1,$2,$3,$4,$5,$6)`, p.TenantID, id, d.requestID(id), expectedSHA, f, p.Subject)
		if err != nil {
			return err
		}
		eventID := uuid.NewSHA1(uuid.NameSpaceOID, []byte("document-committed:"+p.TenantID+":"+id)).String()
		payload, _ := json.Marshal(map[string]string{"document_id": id, "organization_id": d.scope.OrganizationID, "payload_sha256": expectedSHA, "mode": v.Mode})
		_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'document',$3,1,'document.committed',1,clock_timestamp(),$4)`, p.TenantID, eventID, id, payload)
		return err
	})
	if errors.Is(e, approval.ErrNotPending) {
		_ = tx.Rollback(ctx)
		current, readErr := d.Read(ctx, p, id)
		if readErr != nil {
			return current, readErr
		}
		if current.PayloadSHA != expectedSHA || current.Reviewer != p.Subject || current.Reason != reason || (approved && current.State != "PERSISTED") || (!approved && current.State != "REJECTED") {
			return current, ErrDocumentConflict
		}
		return current, nil
	}
	if e != nil {
		return v, e
	}
	if e = tx.Commit(ctx); e != nil {
		return v, e
	}
	return d.Read(ctx, p, id)
}

// EvidencePart returns the original serialized bytes for authorized review;
// fixed part names never become paths or SQL fragments.
func (d *Documents) EvidencePart(ctx context.Context, p identity.Principal, id, part string) ([]byte, error) {
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return nil, e
	}
	if v.EvidenceSHA == "" {
		return nil, ErrDocumentConflict
	}
	var security, provider, receipt []byte
	e = d.pool.QueryRow(ctx, `select security_receipt,provider_response,analysis_receipt from document.extraction where tenant_id=$1 and document_id=$2`, p.TenantID, id).Scan(&security, &provider, &receipt)
	if e != nil {
		return nil, e
	}
	if doc.Hash(receipt) != v.EvidenceSHA {
		return nil, doc.ErrContract
	}
	switch part {
	case "security":
		return security, nil
	case "provider":
		return provider, nil
	case "analysis":
		return receipt, nil
	default:
		return nil, doc.ErrContract
	}
}
````

### FILE: `internal/platform/postgres/document_connected_integration_test.go`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3aa49785fa237c77c1a6028709b76ab4881a445e7d7aaaaa76fe80a76822cf7d"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED real PostgreSQL/HTTP and original SDK fixture composition. Simulated
// detector/provider output cannot establish private document or OCR accuracy.
import (
	"bytes"
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type documentVerifier map[string]identity.Principal

func (v documentVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	p, ok := v[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}

type documentProcessorFunc func(context.Context, doc.Original) (doc.Evidence, error)

func (f documentProcessorFunc) Run(ctx context.Context, o doc.Original) (doc.Evidence, error) {
	return f(ctx, o)
}

type documentFixture struct {
	pool            *pgxpool.Pool
	scope           doc.Scope
	maker, reviewer identity.Principal
	original        doc.Original
	pipeline        *doc.Pipeline
}

func newDocumentFixture(t *testing.T) documentFixture {
	t.Helper()
	url := os.Getenv("DOCUMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned local document DB required")
	}
	parsed, parseErr := urlpkg.Parse(url)
	if parseErr != nil || (parsed.Hostname() != "127.0.0.1" && parsed.Hostname() != "localhost") || (!strings.HasPrefix(parsed.Path, "/elite_document_reference_") && !strings.HasPrefix(parsed.Path, "/elite_payment_connected_")) {
		t.Fatal("owned loopback fixture database required")
	}
	pool, e := pgxpool.New(context.Background(), url)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	tenant := uuid.NewString()
	org := "document-store"
	_, e = pool.Exec(context.Background(), `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Fixture','Fixture')`, tenant, "document-"+tenant)
	if e != nil {
		t.Fatal(e)
	}
	_, e = pool.Exec(context.Background(), `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,$2,$2,'Fixture','store')`, tenant, org)
	if e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile(os.Getenv("ELITE_DOCUMENT_REFERENCE_FIXTURE"))
	if e != nil || doc.Hash(raw) != doc.PublicFixtureSHA {
		t.Fatal("exact public fixture missing")
	}
	profile := filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "config", "documents", "reference-profile.json")
	b, e := os.ReadFile(profile)
	if e != nil {
		t.Fatal(e)
	}
	hash := doc.Hash(b)
	processor, e := doc.NewFixturePipeline(profile, hash, t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	principal := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Organizations: map[string]struct{}{org: {}}, Permissions: map[string]struct{}{"documents:write": {}, "documents:review": {}, "documents:process": {}}}
	}
	return documentFixture{pool, doc.Scope{TenantID: tenant, OrganizationID: org, ProfileSHA: hash, Mode: "FIXTURE"}, principal("uploader"), principal("reviewer"), doc.Original{Name: "invoice.jpg", SHA256: doc.PublicFixtureSHA, Bytes: raw}, processor}
}
func TestDocumentConnectedReference(t *testing.T) {
	f := newDocumentFixture(t)
	ctx := context.Background()
	var calls atomic.Int32
	service, e := db.NewDocuments(f.pool, f.scope, documentProcessorFunc(func(c context.Context, o doc.Original) (doc.Evidence, error) {
		calls.Add(1)
		return f.pipeline.Run(c, o)
	}))
	if e != nil {
		t.Fatal(e)
	}
	wrong := f.maker
	wrong.Organizations = map[string]struct{}{"other": {}}
	alien := f.maker
	alien.TenantID = uuid.NewString()
	unprivileged := f.maker
	unprivileged.Permissions = map[string]struct{}{}
	mux := http.NewServeMux()
	httpapi.DocumentModule{Store: service}.Register(mux, documentVerifier{"maker": f.maker, "reviewer": f.reviewer, "wrong": wrong, "alien": alien, "none": unprivileged})
	var drop atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, r)
		if drop.Load() && r.Method == "POST" && rec.Code == 200 {
			drop.Store(false)
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
			} else {
				_ = conn.Close()
			}
			return
		}
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		_, _ = w.Write(rec.Body.Bytes())
	}))
	defer server.Close()
	call := func(method, path, token string, body []byte) (int, []byte, error) {
		r, _ := http.NewRequest(method, server.URL+"/v1/documents/"+path, bytes.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", "application/json")
		if method == "PUT" {
			r.Header.Set("Content-Type", "application/octet-stream")
			r.Header.Set("X-Document-Name", f.original.Name)
			r.Header.Set("X-Document-SHA256", f.original.SHA256)
			r.Header.Set("X-Document-Profile-SHA256", f.scope.ProfileSHA)
		}
		response, err := server.Client().Do(r)
		if err != nil {
			return 0, nil, err
		}
		defer response.Body.Close()
		b, err := io.ReadAll(response.Body)
		return response.StatusCode, b, err
	}
	expect := func(method, path, token string, body []byte, status int) db.DocumentView {
		t.Helper()
		code, b, err := call(method, path, token, body)
		if err != nil || code != status {
			t.Fatalf("%s %s: status=%d want=%d err=%v body=%s", method, path, code, status, err, b)
		}
		var v db.DocumentView
		if status < 300 && json.Unmarshal(b, &v) != nil {
			t.Fatal("bad view")
		}
		return v
	}
	raw := func(v any) []byte {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	id := uuid.NewString()
	v := expect("PUT", id+"/original", "maker", f.original.Bytes, 201)
	if v.State != "QUARANTINED" || v.Mode != "FIXTURE" {
		t.Fatal(v)
	}
	expect("PUT", id+"/original", "maker", f.original.Bytes, 201)
	expect("PUT", id+"/original", "reviewer", f.original.Bytes, 409)
	expect("PUT", uuid.NewString()+"/original", "maker", []byte("changed bytes"), 400)
	for _, token := range []string{"wrong", "alien", "none"} {
		expect("GET", id, token, nil, 403)
	}
	expect("GET", id, "unknown", nil, 401)
	code, original, e := call("GET", id+"/original", "maker", nil)
	if e != nil || code != 200 || !bytes.Equal(original, f.original.Bytes) {
		t.Fatal("original retention")
	}
	drop.Store(true)
	if _, _, e = call("POST", id+"/process", "maker", []byte(`{}`)); e == nil {
		t.Fatal("response loss not injected")
	}
	v = expect("POST", id+"/process", "maker", []byte(`{}`), 200)
	if v.State != "REVIEW_REQUIRED" || calls.Load() != 1 || v.Suggested == nil {
		t.Fatal("durable extraction replay", v, calls.Load())
	}
	var storedSecurity []byte
	if e = f.pool.QueryRow(ctx, `select security_receipt from document.extraction where tenant_id=$1 and document_id=$2`, f.scope.TenantID, id).Scan(&storedSecurity); e != nil || !bytes.Contains(storedSecurity, []byte("SIMULATED")) {
		t.Fatal("lost fixture provenance")
	}
	fields := *v.Suggested
	fields.Vendor = "Explicitly corrected fixture vendor"
	proposal := map[string]any{"evidence_sha256": v.EvidenceSHA, "fields": fields}
	expect("POST", id+"/review", "maker", []byte(`{"evidence_sha256":"a","EVIDENCE_SHA256":"b"}`), 400)
	v = expect("POST", id+"/review", "maker", raw(proposal), 200)
	if v.Proposal == nil || v.Proposal.Fields.Vendor != fields.Vendor {
		t.Fatal("correction not retained")
	}
	decision := map[string]any{"payload_sha256": v.PayloadSHA, "approved": true, "reason": "Compared four fields with public original; fixture only"}
	expect("POST", id+"/decision", "maker", raw(decision), 409)
	stale := map[string]any{"payload_sha256": strings.Repeat("0", 64), "approved": true, "reason": "stale"}
	expect("POST", id+"/decision", "reviewer", raw(stale), 409)
	// Force an outbox failure and verify decision+commit roll back together.
	_, e = f.pool.Exec(ctx, `create function document.fixture_fail_outbox()returns trigger language plpgsql as $$begin if new.event_type='document.committed'then raise exception 'fixture outbox unavailable';end if;return new;end$$;create trigger fixture_outbox_failure before insert on platform.outbox_event for each row execute function document.fixture_fail_outbox()`)
	if e != nil {
		t.Fatal(e)
	}
	expect("POST", id+"/decision", "reviewer", raw(decision), 503)
	var state string
	var n int
	e = f.pool.QueryRow(ctx, `select state,(select count(*)from document.committed where tenant_id=$1 and document_id=$2)from approval.request where tenant_id=$1 and request_id=$3`, f.scope.TenantID, id, "document:"+id).Scan(&state, &n)
	if e != nil || state != "pending" || n != 0 {
		t.Fatal("partial approval/commit", e, state, n)
	}
	_, e = f.pool.Exec(ctx, `drop trigger fixture_outbox_failure on platform.outbox_event;drop function document.fixture_fail_outbox()`)
	if e != nil {
		t.Fatal(e)
	}
	drop.Store(true)
	if _, _, e = call("POST", id+"/decision", "reviewer", raw(decision)); e == nil {
		t.Fatal("approval response loss not injected")
	}
	v = expect("POST", id+"/decision", "reviewer", raw(decision), 200)
	if v.State != "PERSISTED" || v.Reviewer != f.reviewer.Subject {
		t.Fatal(v)
	}
	decision["reason"] = "changed replay"
	expect("POST", id+"/decision", "reviewer", raw(decision), 409)
	for _, q := range []string{
		`update document.original set name='changed'where tenant_id=$1 and document_id=$2`,
		`delete from document.extraction where tenant_id=$1 and document_id=$2`,
		`update document.committed set reviewer='changed'where tenant_id=$1 and document_id=$2`,
	} {
		if _, e = f.pool.Exec(ctx, q, f.scope.TenantID, id); e == nil {
			t.Fatal("immutable document changed")
		}
	}
	if _, e = f.pool.Exec(ctx, `update approval.request set payload=payload||'{"unbound":true}'::jsonb where tenant_id=$1 and request_id=$2`, f.scope.TenantID, "document:"+id); e == nil {
		t.Fatal("immutable proposal changed")
	}
	e = f.pool.QueryRow(ctx, `select count(*)from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='document.committed'`, f.scope.TenantID, id).Scan(&n)
	if e != nil || n != 1 {
		t.Fatal("duplicate commit outbox", e, n)
	}
	// A rejection is terminal and has no persisted business document.
	rejected := uuid.NewString()
	if _, e = service.Receive(ctx, f.maker, rejected, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	rv, e := service.Process(ctx, f.maker, rejected)
	if e != nil {
		t.Fatal(e)
	}
	rv, e = service.Submit(ctx, f.maker, rejected, rv.EvidenceSHA, fields)
	if e != nil {
		t.Fatal(e)
	}
	rv, e = service.Decide(ctx, f.reviewer, rejected, rv.PayloadSHA, false, "fixture rejected")
	if e != nil || rv.State != "REJECTED" {
		t.Fatal(e, rv)
	}
	if _, e = service.Decide(ctx, f.reviewer, rejected, rv.PayloadSHA, true, "changed decision"); e == nil {
		t.Fatal("rejection reversed")
	}
	down, e := os.ReadFile(filepath.Join(os.Getenv("DOCUMENT_TEST_ROOT"), "db", "migrations", "0085_document_review.down.sql"))
	if e != nil {
		t.Fatal(e)
	}
	conn, e := f.pool.Acquire(ctx)
	if e != nil {
		t.Fatal(e)
	}
	_, e = conn.Exec(ctx, string(down))
	_, _ = conn.Exec(ctx, "rollback")
	conn.Release()
	if e == nil {
		t.Fatal("populated down migration erased evidence")
	}
	t.Log("DOCUMENT_CONNECTED_PASS: original+job atomic; scoped HTTP; original SDK; explicit simulated detectors; corrected immutable review; independent reviewer; commit+approval+outbox atomic; lost response replay; rejection; immutable evidence; populated rollback refused")
}
func TestDocumentRecoveryAndFencing(t *testing.T) {
	f := newDocumentFixture(t)
	ctx := context.Background()
	var calls atomic.Int32
	wrapper := documentProcessorFunc(func(c context.Context, o doc.Original) (doc.Evidence, error) {
		if calls.Add(1) == 1 {
			return doc.Evidence{}, doc.ErrExtraction
		}
		return f.pipeline.Run(c, o)
	})
	store, e := db.NewDocuments(f.pool, f.scope, wrapper)
	if e != nil {
		t.Fatal(e)
	}
	id := uuid.NewString()
	if _, e = store.Receive(ctx, f.maker, id, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	if _, e = store.Process(ctx, f.maker, id); !errors.Is(e, doc.ErrExtraction) {
		t.Fatal(e)
	}
	ready := func(id string) {
		t.Helper()
		_, err := f.pool.Exec(ctx, `update platform.job set available_at=clock_timestamp()-interval '1 second' where tenant_id=$1 and payload->>'document_id'=$2`, f.scope.TenantID, id)
		if err != nil {
			t.Fatal(err)
		}
	}
	ready(id)
	v, e := store.Process(ctx, f.maker, id)
	if e != nil || v.State != "REVIEW_REQUIRED" || calls.Load() != 2 {
		t.Fatal(v, e, calls.Load())
	}
	var count int
	e = f.pool.QueryRow(ctx, `select count(*)from document.attempt_failure where tenant_id=$1 and document_id=$2`, f.scope.TenantID, id).Scan(&count)
	if e != nil || count != 1 {
		t.Fatal("failure receipt lost", e, count)
	}
	// A delayed old generation may finish its extraction, but cannot commit it.
	started, release := make(chan struct{}), make(chan struct{})
	var attempts atomic.Int32
	fenced, e := db.NewDocuments(f.pool, f.scope, documentProcessorFunc(func(c context.Context, o doc.Original) (doc.Evidence, error) {
		if attempts.Add(1) == 1 {
			close(started)
			select {
			case <-release:
			case <-c.Done():
				return doc.Evidence{}, c.Err()
			}
		}
		return f.pipeline.Run(c, o)
	}))
	if e != nil {
		t.Fatal(e)
	}
	stale := uuid.NewString()
	if _, e = fenced.Receive(ctx, f.maker, stale, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() { _, err := fenced.Process(ctx, f.maker, stale); done <- err }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("worker not started")
	}
	_, e = f.pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and payload->>'document_id'=$2`, f.scope.TenantID, stale)
	if e != nil {
		t.Fatal(e)
	}
	v, e = fenced.Process(ctx, f.maker, stale)
	close(release)
	if e != nil || v.State != "REVIEW_REQUIRED" {
		t.Fatal(e, v)
	}
	select {
	case err := <-done:
		if !errors.Is(err, db.ErrJobClaimLost) {
			t.Fatal("stale worker committed", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("stale worker hung")
	}
	e = f.pool.QueryRow(ctx, `select count(*)from document.extraction where tenant_id=$1 and document_id=$2 and attempt=2`, f.scope.TenantID, stale).Scan(&count)
	if e != nil || count != 1 {
		t.Fatal("fence lost", e, count)
	}
	// Security refusal never reaches the official provider; finite attempts stop.
	blocked, e := db.NewDocuments(f.pool, f.scope, documentProcessorFunc(func(context.Context, doc.Original) (doc.Evidence, error) { return doc.Evidence{}, doc.ErrSecurity }))
	if e != nil {
		t.Fatal(e)
	}
	bad := uuid.NewString()
	if _, e = blocked.Receive(ctx, f.maker, bad, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 3; i++ {
		ready(bad)
		if _, e = blocked.Process(ctx, f.maker, bad); !errors.Is(e, doc.ErrSecurity) {
			t.Fatal(e)
		}
	}
	v, e = blocked.Read(ctx, f.maker, bad)
	if e != nil || v.State != "QUARANTINE_TERMINAL" {
		t.Fatal(e, v)
	}
	// A crashed last attempt is exhausted by the original job owner, not reset.
	exhausted := uuid.NewString()
	if _, e = blocked.Receive(ctx, f.maker, exhausted, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	_, e = f.pool.Exec(ctx, `update platform.job set attempts=max_attempts,claimed_by='crashed-fixture',claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and payload->>'document_id'=$2`, f.scope.TenantID, exhausted)
	if e != nil {
		t.Fatal(e)
	}
	v, e = blocked.Process(ctx, f.maker, exhausted)
	if e != nil || v.State != "QUARANTINE_TERMINAL" {
		t.Fatal(e, v)
	}
	t.Log("DOCUMENT_RECOVERY_PASS: durable transient failure; bounded retry; generation fencing; security terminal quarantine; exhausted crashed attempt retained")
}

func TestDocumentProfileHistoryAndDatabaseReviewContract(t *testing.T) {
	f := newDocumentFixture(t)
	ctx := context.Background()
	store, e := db.NewDocuments(f.pool, f.scope, f.pipeline)
	if e != nil {
		t.Fatal(e)
	}
	id := uuid.NewString()
	if _, e = store.Receive(ctx, f.maker, id, f.scope.ProfileSHA, f.original); e != nil {
		t.Fatal(e)
	}
	changed := f.scope
	changed.ProfileSHA = strings.Repeat("b", 64)
	newStore, e := db.NewDocuments(f.pool, changed, f.pipeline)
	if e != nil {
		t.Fatal(e)
	}
	v, e := newStore.Read(ctx, f.maker, id)
	if e != nil || v.State != "QUARANTINED" || v.ProfileSHA != f.scope.ProfileSHA {
		t.Fatal("stored profile history unreadable", e, v)
	}
	if _, e = newStore.Process(ctx, f.maker, id); !errors.Is(e, db.ErrDocumentConflict) {
		t.Fatal("stale profile processed", e)
	}
	v, e = store.Process(ctx, f.maker, id)
	if e != nil {
		t.Fatal(e)
	}
	payload := map[string]any{"schema": "document-review/v1", "document_id": id, "organization_id": f.scope.OrganizationID, "original_sha256": v.OriginalSHA, "profile_sha256": v.ProfileSHA, "evidence_sha256": v.EvidenceSHA, "mode": "FIXTURE", "fields": map[string]string{"invoice_number": "only one field"}}
	raw, _ := json.Marshal(payload)
	if _, e = f.pool.Exec(ctx, `insert into approval.request(tenant_id,request_id,kind,subject_id,amount_minor_units,requester,evidence_sha,organization_id,payload)values($1,$2,'document_review',$3,0,'fixture-internal-caller',$4,$5,$6)`, f.scope.TenantID, "document:"+id, id, doc.Hash(raw), f.scope.OrganizationID, raw); e == nil {
		t.Fatal("database admitted incomplete fields")
	}
	// Restart evidence: earlier connected run rows and corrected commit survive.
	var n int
	e = f.pool.QueryRow(ctx, `select count(*)from document.committed c join document.original o using(tenant_id,document_id)where o.mode='FIXTURE'and encode(sha256(o.content),'hex')=o.original_sha256 and c.fields->>'vendor'='Explicitly corrected fixture vendor'`).Scan(&n)
	if e != nil || n < 1 {
		t.Fatal("prior committed reference did not survive host restart", e, n)
	}

	mux := http.NewServeMux()
	httpapi.DocumentModule{Store: store}.Register(mux, documentVerifier{"maker": f.maker})
	for _, part := range []string{"security", "provider", "analysis"} {
		req := httptest.NewRequest("GET", "/v1/documents/"+id+"/evidence/"+part, nil)
		req.Header.Set("Authorization", "Bearer maker")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		exact, err := store.EvidencePart(ctx, f.maker, id, part)
		if err != nil || rec.Code != 200 || !bytes.Equal(exact, rec.Body.Bytes()) || rec.Header().Get("X-Document-Evidence-SHA256") != doc.Hash(exact) || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("exact evidence review failed", part, err, rec.Code)
		}
	}
	if _, e = store.EvidencePart(ctx, f.maker, id, "../original"); !errors.Is(e, doc.ErrContract) {
		t.Fatal("unbounded evidence part", e)
	}
	t.Log("DOCUMENT_REVIEW_DELTA_PASS: stored-profile history; stale execution refusal; database field contract; original and commit survive process restart")
}
````

### FILE: `tools/verify_document_reference.py`

```yaml
block_id: "GO-CONNECTED-DOCUMENT-REFERENCE:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c9f33d2db369d40d99ed485d3098f799fa1593620ccfc5622cf6ffa1cb1da156"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED finite local fixture runner; does not provision or address live services."""
from pathlib import Path
import argparse,hashlib,json,os,subprocess,time
from urllib.parse import urlsplit
p=argparse.ArgumentParser();p.add_argument('--go',type=Path,required=True);p.add_argument('--fixture',type=Path,required=True);p.add_argument('--receipt',type=Path,required=True);args=p.parse_args()
root=Path(__file__).resolve().parents[1]
sha=lambda path:hashlib.sha256(path.read_bytes()).hexdigest()
if sha(args.go)!='21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9':raise SystemExit('exact admitted Go1.26.8 executable required')
if sha(args.fixture)!='489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb':raise SystemExit('exact public fixture required')
url=urlsplit(os.environ.get('DOCUMENT_CONNECTED_DB_URL',''))
if url.scheme not in ('postgres','postgresql') or url.hostname not in ('127.0.0.1','localhost') or not url.path.startswith(('/elite_document_reference_','/elite_payment_connected_')):raise SystemExit('owned loopback synthetic document database required; apply selected migrations first')
if args.receipt.exists() or not args.receipt.parent.is_dir():raise SystemExit('receipt must be new with existing parent')
env=os.environ.copy();env.pop('GOROOT',None);env.update(GOTOOLCHAIN='local',GOPROXY='off',GOSUMDB='off',GOWORK='off',GOFLAGS='-mod=readonly',GOMAXPROCS='4',DOCUMENT_TEST_ROOT=str(root),ELITE_DOCUMENT_REFERENCE_FIXTURE=str(args.fixture.resolve()))
steps=[('connected',['test','./internal/platform/postgres','-run','^TestDocument','-count=1','-v','-timeout','120s']),('host',['test','./cmd/electromobility-api','-run','^TestDocumentHostActivation$','-count=1','-v','-timeout','60s']),('bounds',['test','./internal/documentbridge','./internal/approval','-run','^(TestDocument|FuzzDocument)','-count=1','-v']),('fuzz',['test','./internal/documentbridge','-run','^$','-fuzz','^FuzzDocumentOriginalBoundary$','-fuzztime','3s','-parallel','2']),('vet',['vet','./internal/documentbridge','./internal/approval','./internal/platform/postgres','./internal/platform/httpapi','./cmd/electromobility-api'])]
result={'schema':'elite-document-reference-verification/v1','state':'RUNNING','scope':'LOCAL_FIXTURES','live_effects':False,'detectors':'SIMULATED; no scanner effectiveness proof','go_sha256':sha(args.go),'fixture_sha256':sha(args.fixture),'steps':[]};start=time.monotonic()
try:
 for label,argv in steps:
  log=args.receipt.with_name(args.receipt.stem+'-'+label+'.log')
  with log.open('xb')as f:q=subprocess.run([str(args.go),*argv],cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT,timeout=150,creationflags=getattr(subprocess,'CREATE_NO_WINDOW',0))
  raw=log.read_bytes();result['steps'].append({'name':label,'exit_code':q.returncode,'log_sha256':sha(log),'log':str(log),'no_skip':b'--- SKIP:' not in raw})
  if q.returncode or b'--- SKIP:'in raw:raise RuntimeError(label+' failed; consult preserved log')
 result['state']='PASS'
except BaseException as e:
 result.update(state='FAIL',failure_type=type(e).__name__);raise
finally:
 result['seconds']=round(time.monotonic()-start,3);args.receipt.write_text(json.dumps(result,indent=2)+'\n',encoding='utf-8',newline='\n')
````

## 6. Configuration surface

docs/DOCUMENT_REFERENCE_START.md and docs/DOCUMENT_REFERENCE_DECISION.md. Feature opt-in, fixed profile SHA, explicit tenant/org/permissions, work directory; future credentials only via official SDK credential chain.

## 7. Dependency bill

Existing official AWS GoSDK Textract1.45.0 and15exact modules, original PG/Go/OIDC/approval owners; no new module identities beyond verified existing SDK profile. Native detector qualification remains its original separate gate; no scanner-effectiveness claim from fixtures.

## 8. Apply order

Select full franchise composition and apply selected migrations through0085. SDK owner also added to every backend plan because its common go.mod requires the local module. Reference-only document files/fixture/scanner glue are selected only in complete franchise and HTTP-metrics superset.

## 9. Verification

tools/verify_document_reference.py: original+job, HTTP scope, official SDK serialization, independent review, rollback atomicity, response-loss recovery, stale worker fence, terminal quarantine, profile history, exact evidence download, immutable DB, no populated downgrade; finite fuzz3s;11original secure-gate contract tests. SDK tests use simulation explicitly, no live OCR.

## 10. Reconstruction evidence

reconstruction_evidence/DOCUMENT_REFERENCE_RELEASE_V402.md/json; connected-admission G0-G8, SDK zero-finding union, exact-source delta lint24reviewed, reconstructed full profiles. No globalREADY or production authorization until remaining controls complete.

