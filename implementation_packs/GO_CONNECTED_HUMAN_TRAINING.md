# Durable human training bridge

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-HUMAN-TRAINING"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Explicit same-release training content through durable participation and distinct human assessment; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J5 slice. Compose existing help, shared approval, audit/outbox and OIDC owners. Four role curricula guide public content; permissions remain independently administered. No automatic score, labor qualification, employee creation or grant.

## 3. Architecture contract

Immutable profile/content/actor/org identity; append-only participation facts in audit.event, assessment in existing approval.request/decision, atomic existing outbox. One assessment per attempt. Physical guards reject generic approval bypass and audit mutation. Same-release snapshot survives content revisions.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/training.go
CREATE cmd/electromobility-api/training_test.go
CREATE db/migrations/0075_training_assessment.down.sql
CREATE db/migrations/0075_training_assessment.up.sql
CREATE deploy/training/reference.profile.json
CREATE docs/TRAINING_REFERENCE.md
CREATE internal/trainingbridge/browser_issuer_test.go
CREATE internal/trainingbridge/browser_test.go
CREATE internal/trainingbridge/connected_test.go
CREATE internal/trainingbridge/guard_test.go
CREATE internal/trainingbridge/http.go
CREATE internal/trainingbridge/pack_profile_test.go
CREATE internal/trainingbridge/profile.go
CREATE internal/trainingbridge/profile_fuzz_test.go
CREATE internal/trainingbridge/profile_test.go
CREATE internal/trainingbridge/store.go
CREATE training_content/help.bundle.json
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/training.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "655b907536d127ac1406af9dec70bdb07a29f0a69bc98c0f67006104bc956d3a"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED opt-in composition. No environment value grants learner permissions.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/trainingbridge"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

func init() { trainingModuleFactory = selectedTrainingModule }
func selectedTrainingModule(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	if getenv == nil {
		return nil, trainingbridge.ErrContract
	}
	enabled := getenv("TRAINING_ENABLED")
	if enabled == "" || enabled == "false" {
		return nil, nil
	}
	if enabled != "true" {
		return nil, trainingbridge.ErrContract
	}
	if pool == nil {
		return nil, trainingbridge.ErrContract
	}
	var ready bool
	e := pool.QueryRow(ctx, `select
 (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('audit.event') and tgname='training_audit_immutable'
 or tgrelid=to_regclass('approval.request') and tgname='training_assessment_guard'))=2
 and to_regclass('audit.audit_training_attempt_idx') is not null
 and to_regclass('approval.approval_training_review_idx') is not null`).Scan(&ready)
	if e == nil && ready {
		e = pool.QueryRow(ctx, `select to_regclass('approval.training_assessment_attempt_idx') is not null`).Scan(&ready)
	}
	if e != nil || !ready {
		return nil, trainingbridge.ErrContract
	}
	revision, e := strconv.Atoi(getenv("TRAINING_PROFILE_REVISION"))
	if e != nil {
		return nil, trainingbridge.ErrContract
	}
	profile, e := trainingbridge.LoadProfile(getenv("TRAINING_PROFILE_FILE"), getenv("TRAINING_CONTENT_FILE"), trainingbridge.Activation{Enabled: true, ID: getenv("TRAINING_PROFILE_ID"), Revision: revision, SHA256: getenv("TRAINING_PROFILE_SHA256"), TenantID: getenv("TRAINING_TENANT_ID"), OrganizationID: getenv("TRAINING_ORGANIZATION_ID")})
	if e != nil {
		return nil, e
	}
	store, e := trainingbridge.NewStore(pool, profile)
	if e != nil {
		return nil, e
	}
	return &trainingbridge.Module{Store: store}, nil
}
````

### FILE: `cmd/electromobility-api/training_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bbd23a9f91ad09917166260d5be6a889b51adb0bf6d854367e9544dd2ab9ee32"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/sha256"
	tr "elite.local/enterprise/internal/trainingbridge"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"testing"
)

func TestTrainingHostProfileAndGuards(t *testing.T) {
	ctx := context.Background()
	calls := 0
	off := func(k string) string {
		calls++
		if k != "TRAINING_ENABLED" {
			t.Fatal("disabled host read", k)
		}
		return "false"
	}
	if m, e := selectedTrainingModule(ctx, nil, off); e != nil || m != nil || calls != 1 {
		t.Fatal("disabled", e)
	}
	db := os.Getenv("ELITE_TRAINING_DATABASE_URL")
	if db == "" {
		t.Skip("owned fixture DB required")
	}
	pool, e := pgxpool.New(ctx, db)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	b, e := os.ReadFile("../../training_content/help.bundle.json")
	if e != nil {
		t.Fatal(e)
	}
	var bundle tr.ContentBundle
	if e = json.Unmarshal(b, &bundle); e != nil {
		t.Fatal(e)
	}
	sum := sha256.Sum256(b)
	d := tr.ProfileDocument{Schema: "elite-training-profile/v1", ID: "host-reference", Revision: 1, TenantID: "50f38793-8a22-4f6b-983f-81dd0fca8208", OrganizationID: "store-1", Method: tr.Method, ContentSHA256: hex.EncodeToString(sum[:]), SourceSHA256: bundle.SourceSHA256, Courses: []tr.Course{{ID: "resource-onboarding", Title: "Recuperar un registro", Role: "employee", Lessons: []string{"resource-create-view"}, Prompts: []tr.Prompt{{ID: "recovery", Text: "Explicá cómo consultarías el registro."}}}}}
	raw, _ := json.Marshal(d)
	sum = sha256.Sum256(raw)
	dir := t.TempDir()
	profile, content := filepath.Join(dir, "profile.json"), filepath.Join(dir, "content.json")
	if e = os.WriteFile(profile, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(content, b, 0600); e != nil {
		t.Fatal(e)
	}
	values := map[string]string{"TRAINING_ENABLED": "true", "TRAINING_PROFILE_FILE": profile, "TRAINING_CONTENT_FILE": content, "TRAINING_PROFILE_ID": d.ID, "TRAINING_PROFILE_REVISION": "1", "TRAINING_PROFILE_SHA256": hex.EncodeToString(sum[:]), "TRAINING_TENANT_ID": d.TenantID, "TRAINING_ORGANIZATION_ID": d.OrganizationID}
	lookup := func(k string) string { return values[k] }
	if m, e := selectedTrainingModule(ctx, pool, lookup); e != nil || m == nil {
		t.Fatal("valid activation", e)
	}
	values["TRAINING_PROFILE_FILE"] = "relative.json"
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("relative profile")
	}
	values["TRAINING_PROFILE_FILE"] = dir
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("directory profile")
	}
	values["TRAINING_PROFILE_FILE"] = profile
	values["TRAINING_PROFILE_SHA256"] = ""
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("unbound profile")
	}
	values["TRAINING_PROFILE_SHA256"] = hex.EncodeToString(sum[:])
	if _, e = pool.Exec(ctx, `alter table approval.request disable trigger training_assessment_guard`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table approval.request enable trigger training_assessment_guard`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedTrainingModule(ctx, pool, lookup); e == nil {
		t.Fatal("disabled guard")
	}
	t.Log("TRAINING_HOST_PASS opt-in no-read; exact profile/content scope; regular absolute paths; shared participation/review guards and indexes required")
}
````

### FILE: `db/migrations/0075_training_assessment.down.sql`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2143eec3e552ac0fce33374e22698947e6892b01dc707a6cfaf1aa403edd1483"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from approval.request where kind='training_assessment') or exists(select 1 from audit.event where resource_type='training-attempt') then
  raise exception 'training history exists; downgrade refused';
 end if;
end$$;
drop trigger training_assessment_guard on approval.request;
drop function approval.require_training_facts();
drop trigger training_audit_immutable on audit.event;
drop function audit.protect_training_fact();
drop index approval.approval_training_review_idx;
drop index approval.training_assessment_attempt_idx;
drop index audit.audit_training_attempt_idx;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

commit;
````

### FILE: `db/migrations/0075_training_assessment.up.sql`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "10cc4d86b93dc111ae48a5e0d895122874a485710bac6d71a9a9fd15179eef0b"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create index audit_training_attempt_idx on audit.event(tenant_id,resource_id,actor_subject,action,audit_sequence) where resource_type='training-attempt';
create index approval_training_review_idx on approval.request(tenant_id,organization_id,created_at desc,request_id) where kind='training_assessment';
create unique index training_assessment_attempt_idx on approval.request(tenant_id,organization_id,(payload->>'attempt_id')) where kind='training_assessment';
create function audit.protect_training_fact() returns trigger language plpgsql as $$
begin
 if old.resource_type='training-attempt' or (tg_op='UPDATE' and new.resource_type='training-attempt') then
  raise exception 'training participation facts are immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end$$;
create trigger training_audit_immutable before update or delete on audit.event for each row execute function audit.protect_training_fact();
create function approval.require_training_facts() returns trigger language plpgsql as $$
declare r approval.request%rowtype; a jsonb;
begin
 if new.kind<>'training_assessment' then return new;end if;
 select * into strict r from approval.request where tenant_id=new.tenant_id and request_id=new.request_id;
 if r.payload->>'schema'<>'training-assessment/v1' or r.payload->>'learner_subject' is distinct from r.requester
  or r.subject_id is distinct from r.requester or r.payload->>'organization_id' is distinct from r.organization_id then
  raise exception 'training assessment identity mismatch';
 end if;
 a:=r.payload;
 if not exists(select 1 from audit.event e where e.tenant_id=r.tenant_id and e.resource_type='training-attempt'
  and e.resource_id=a->>'attempt_id' and e.actor_subject=r.requester and e.action='training.started'
  and e.evidence->>'organization_id'=r.organization_id and e.evidence->'content'=a->'content') then
  raise exception 'training assessment has no matching immutable start';
 end if;
 if not exists(select 1 from audit.event e where e.tenant_id=r.tenant_id and e.resource_type='training-attempt'
  and e.resource_id=a->>'attempt_id' and e.actor_subject=r.requester and e.action='training.submitted'
  and e.evidence->>'organization_id'=r.organization_id and e.evidence->>'payload_sha256'=r.evidence_sha) then
  raise exception 'training assessment submission not recorded';
 end if;
 if exists(select 1 from jsonb_array_elements_text(a#>'{content,course,lessons}') lesson
  where not exists(select 1 from audit.event e where e.tenant_id=r.tenant_id and e.resource_type='training-attempt'
   and e.resource_id=a->>'attempt_id' and e.actor_subject=r.requester and e.action='training.lesson-read'
   and e.evidence->>'lesson_id'=lesson and e.evidence->>'organization_id'=r.organization_id
   and e.evidence->>'profile_sha256'=a#>>'{content,profile_sha256}')) then
  raise exception 'training assessment requires recorded lesson declarations';
 end if;
 if r.state<>'pending' and not exists(select 1 from approval.decision d join audit.event e
  on e.tenant_id=d.tenant_id and e.actor_subject=d.reviewer and e.resource_type='training-attempt'
   and e.resource_id=a->>'attempt_id' and e.action='training.assessed'
  where d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.reviewer<>r.requester
   and e.evidence->>'organization_id'=r.organization_id and e.evidence->>'payload_sha256'=r.evidence_sha
   and e.evidence->>'reason'=d.reason and (e.evidence->>'approved')::boolean=d.approved
   and d.approved=(r.state='approved')) then
  raise exception 'training assessment decision has no matching human review fact';
 end if;
 return new;
end$$;
create constraint trigger training_assessment_guard after insert or update on approval.request
 deferrable initially deferred for each row execute function approval.require_training_facts();
commit;
````

### FILE: `deploy/training/reference.profile.json`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2ce1e88eb44eef172bd4c1d8e09680e2c7c732aa479bee65d855e3251c53f937"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-training-profile/v1",
  "id": "reference-onboarding",
  "revision": 2,
  "tenant_id": "50f38793-8a22-4f6b-983f-81dd0fca8208",
  "organization_id": "store-1",
  "method": "HUMAN_REVIEW_NO_GRANTS_V1",
  "content_sha256": "eb919b9731d56832c38233572d6f1fd71cc62e7da64ca3d53a97638ab6a3ea20",
  "source_sha256": "1c1ebc33ca2e578725f9d0ca4f51588f6f480e81ab9e8cbbdfc05b963395a581",
  "courses": [
    {
      "id": "resource-onboarding",
      "title": "Registrar recursos y recuperar una respuesta",
      "role": "employee",
      "lessons": [
        "resource-create-view"
      ],
      "prompts": [
        {
          "id": "recovery",
          "text": "Describí qué harías si se pierde la respuesta después de registrar un recurso."
        }
      ]
    },
    {
      "id": "admin-onboarding",
      "title": "Revisión de recursos",
      "role": "admin",
      "lessons": [
        "resource-create-view",
        "supply-role-view",
        "warranty-role-view",
        "network-role-view"
      ],
      "prompts": [
        {
          "id": "review",
          "text": "Describí cómo recuperarías una operación incierta, revisarías contenido y separarías una evaluación de los permisos de acceso."
        }
      ]
    },
    {
      "id": "owner-onboarding",
      "title": "Supervisión de la operación",
      "role": "owner",
      "lessons": [
        "operation-sections-view",
        "network-role-view",
        "training-role-view"
      ],
      "prompts": [
        {
          "id": "review",
          "text": "Describí cómo distinguís una sección no disponible de una lista vacía."
        }
      ]
    },
    {
      "id": "customer-onboarding",
      "title": "Revisar una cotización",
      "role": "customer",
      "lessons": [
        "quote-acceptance-view"
      ],
      "prompts": [
        {
          "id": "review",
          "text": "Describí qué confirma aceptar una cotización y qué se debe consultar ante una respuesta incierta."
        }
      ]
    },
    {
      "id": "content-onboarding",
      "title": "Contenido, publicación y revisión humana",
      "role": "admin",
      "lessons": [
        "help-cms-view",
        "catalog-role-view",
        "training-role-view"
      ],
      "prompts": [
        {
          "id": "review",
          "text": "Explicá cómo publicarías contenido revisado, recuperarías una respuesta incierta y preservarías cursos anteriores sin conceder accesos."
        }
      ]
    }
  ]
}
````

### FILE: `docs/TRAINING_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8329b1c033cc11ab530d4bb62522bd8fbbfa7e19a370ffeb3e38c8ab10b2dd48"
variables: []
secrets_allowed: false
```

````markdown
# Capacitación con evaluación humana — referencia local

El alumno inicia un intento, lee la versión fijada, declara su lectura y presenta
respuestas. Un revisor autorizado distinto registra su evaluación. Cada efecto
usa audit.event, approval.request/decision y platform.outbox_event existentes.
Una evaluación no crea personas, certificaciones ni permisos. El rol del curso
orienta el contenido; la autorización depende de los permisos de la sesión.

Materializar el perfil completo y aplicar migraciones hasta 0075. No incorporar
0064 de candidatos históricos. Este pack conserva todos los kinds de aprobación
anteriores. Hay dos guards físicos y tres índices de identidad/consulta. La
baja con historial se rechaza; nunca borrar evidencia para poder degradar.

Activación explícita del host Go, deshabilitada de forma predeterminada:

| Variable | Valor de referencia |
|---|---|
| TRAINING_ENABLED | true |
| TRAINING_PROFILE_FILE | Ruta absoluta a deploy/training/reference.profile.json |
| TRAINING_CONTENT_FILE | Ruta absoluta a training_content/help.bundle.json |
| TRAINING_PROFILE_ID | reference-onboarding |
| TRAINING_PROFILE_REVISION | 1 |
| TRAINING_PROFILE_SHA256 | 9e897f609970a1e46c16944a6e6ec9abfc2a1ba8c2e4a3a4d0385095cb540907 |
| TRAINING_TENANT_ID | 50f38793-8a22-4f6b-983f-81dd0fca8208 |
| TRAINING_ORGANIZATION_ID | store-1 |

Son identidades sintéticas, sin cuentas ni secretos. Al materializar otro negocio,
fijar sus identidades y programa de capacitación, producir una nueva revisión del
perfil y calcular su SHA-256; no editar intentos ni evaluaciones previas. El host
rechaza paths relativos, archivos no regulares y hashes/identidades inconsistentes.
El proceso normal de identidad debe emitir training:learn o training:review
para la organización; estas variables no conceden autorización.

El contenido contiene 15 guías públicas de la misma release. SHA fuente:
17d210d1ad3bb0c80bdb94a8d8623cc2b8ea71a4fc94824882ef9362c189c4fd; SHA bundle: 3e636fde7b3e6cbdc109de837024424d0b7a35315b8b0a54edd874fec4a042fb.
El exportador del pack web ejecuta solamente el módulo de ayuda admitido en el
proyecto. No usarlo con código de terceros sin admisión. Los antiguos intentos
conservan el snapshot y se pueden consultar; una revisión nueva no los autoriza
a seguir escribiendo. No introducir documentos privados en este bundle público.

Verificación local: TestTrainingDurableBindingsAndAtomicity, perfil, guards,
host, downgrade y TestTrainingHumanBrowserPostgres. Las pruebas de PG requieren
ELITE_TRAINING_DATABASE_URL a una base aislada con las migraciones aplicadas.
El fixture de navegador usa ELITE_TRAINING_BROWSER=1, ELITE_WEB_ROOT y
ELITE_NODE_BIN absolutos, Next previamente compilado y el browser gate ya fijado.
No se contactan proveedores live. Ejecutar el fixture solamente sobre datos
sintéticos aislados; no es un runner sobre una base productiva.

Los tests prueban cuatro hechos/cuatro eventos/una solicitud/una decisión y cero
grants/recursos. Se rechazan aprobación genérica sin evaluación, duplicación de
assessment por intento, mutación de historial, autoevaluación y acceso cruzado.
Ver reconstruction_evidence/TRAINING_CONNECTED_RELEASE_V402.md en la biblioteca.
````

### FILE: `internal/trainingbridge/browser_issuer_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bd98525257a477684dd49b85d4b0b3e821ad85197c0f99385ad122573e993bcd"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

// AUTHORED synthetic RS256/JWKS fixture reused from admitted local browser tests.
import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func trainingBrowserIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": time.Now().Unix(), "exp": time.Now().Add(5 * time.Minute).Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}
````

### FILE: `internal/trainingbridge/browser_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c7739890701f76d16c3a5320fd4336a6ad4a4a623425641d79295f4fe8fc3bd3"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

// AUTHORED real-browser fixture over the actual module and durable owners.
import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTrainingHumanBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_TRAINING_BROWSER") != "1" {
		t.Skip("explicit local training browser fixture required")
	}
	f := connectedFixture(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	verifier, token := trainingBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"learner", "reviewer", "other-learner", "foreign-org", "reader"} {
		permissions := []string{"training:learn"}
		orgs := []string{"store-1"}
		if name == "reviewer" {
			permissions = []string{"training:review"}
		}
		if name == "foreign-org" {
			orgs = []string{"other"}
		}
		if name == "reader" {
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": f.learner.TenantID, "permissions": permissions, "organizations": orgs, "accessToken": token(name, f.learner.TenantID, permissions, orgs)}
	}
	mux := http.NewServeMux()
	Module{Store: f.store}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	controlToken := uuid.NewString()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	web := os.Getenv("ELITE_WEB_ROOT")
	node := os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed web/runtime paths required")
	}
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if strings.HasPrefix(name, "ELITE_") || strings.Contains(name, "DATABASE_URL") || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, _ := json.Marshal(identities)
	env = append(env, "ELITE_TRAINING_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_TRAINING_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-browser-only-"+controlToken, "ELITE_TRAINING_FIXTURE=1")
	artifacts, err := os.MkdirTemp(web, "training-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = log, log
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		resp, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				ready = true
				break
			}
		}
		select {
		case <-ctx.Done():
			t.Fatal("startup cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		t.Fatal("Next did not start", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/training-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}

	var facts, events, requests, decisions, resources int
	err = f.pool.QueryRow(ctx, `select (select count(*) from audit.event where tenant_id=$1 and resource_type='training-attempt'),(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='training-evidence'),(select count(*) from approval.request where tenant_id=$1 and kind='training_assessment'),(select count(*) from approval.decision where tenant_id=$1),(select count(*) from crm.service_resource where tenant_id=$1)`, f.learner.TenantID).Scan(&facts, &events, &requests, &decisions, &resources)
	if err != nil || facts != 4 || events != 4 || requests != 1 || decisions != 1 || resources != 0 {
		t.Fatal("unexpected durable outcomes", facts, events, requests, decisions, resources, err)
	}
	mu.Lock()
	encodedCounts, _ := json.Marshal(counts)
	total := 0
	for _, count := range counts {
		total += count
	}
	mu.Unlock()
	if total != 4 {
		t.Fatal("duplicate or unauthorized backend writes", total)
	}
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), encodedCounts, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("TRAINING_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true learner_start_read_response=true human_assessment=true response_loss_read_only_recovery=true actor_org_isolation=true durable_facts=4 outbox=4 backend_posts=4 permission_grants=0 artifacts=%s", artifacts)
}
````

### FILE: `internal/trainingbridge/connected_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a90742832d03022b83feb2e8ce6a3dc6370bfdd3bb9eca9fcf1b1ac393b949d6"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
)

type trainingFixture struct {
	store                      *Store
	pool                       *pgxpool.Pool
	learner, reviewer, foreign identity.Principal
	profile                    *Profile
}

func connectedFixture(t *testing.T) *trainingFixture {
	t.Helper()
	raw := os.Getenv("ELITE_TRAINING_DATABASE_URL")
	if raw == "" {
		t.Skip("explicit owned training database required")
	}
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_training_connected_") {
		t.Fatal("dedicated loopback training database required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	d, b := profileFixture(t)
	d.TenantID = uuid.NewString()
	bytes, _ := json.Marshal(d)
	p, e := ParseProfile(bytes, b, Activation{true, d.ID, d.Revision, digest(bytes), d.TenantID, d.OrganizationID})
	if e != nil {
		t.Fatal(e)
	}
	for _, q := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'training-'||$1::text,'Fixture','Fixture')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store-1','store-1','Fixture','store'),($1,'other','other','Other','store')`} {
		if _, e = pool.Exec(ctx, q, d.TenantID); e != nil {
			t.Fatal(e)
		}
	}
	s, e := NewStore(pool, p)
	if e != nil {
		t.Fatal(e)
	}
	learner := identity.Principal{TenantID: d.TenantID, Subject: "learner", Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"training:learn": {}, "training:review": {}}}
	reviewer := identity.Principal{TenantID: d.TenantID, Subject: "reviewer", Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"training:review": {}}}
	foreign := identity.Principal{TenantID: d.TenantID, Subject: "foreign", Organizations: map[string]struct{}{"other": {}}, Permissions: map[string]struct{}{"training:learn": {}, "training:review": {}}}
	return &trainingFixture{s, pool, learner, reviewer, foreign, p}
}
func TestTrainingDurableBindingsAndAtomicity(t *testing.T) {
	f := connectedFixture(t)
	ctx := context.Background()
	id := uuid.NewString()
	hash := f.profile.Hash()
	start := func() error { _, e := f.store.Start(ctx, f.learner, id, "resource-onboarding", hash); return e }
	var wg sync.WaitGroup
	errs := make(chan error, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); errs <- start() }()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	if _, e := f.store.Start(ctx, f.learner, id, "admin-onboarding", hash); e == nil {
		t.Fatal("attempt identity changed course")
	}
	if _, e := f.store.Read(ctx, f.foreign, id); e == nil {
		t.Fatal("cross organization attempt visible")
	}
	other := f.learner
	other.Subject = "other-learner"
	other.Permissions = map[string]struct{}{"training:learn": {}}
	if _, e := f.store.Read(ctx, other, id); e == nil {
		t.Fatal("another learner read attempt")
	}
	answers := map[string]string{"recovery": "Consulto la referencia guardada y no creo otro recurso para forzar un reintento."}
	if _, e := f.store.Submit(ctx, f.learner, id, hash, answers); e == nil {
		t.Fatal("response submitted before declared reading")
	}
	if _, e := f.pool.Exec(ctx, `create function audit.reject_training_fixture() returns trigger language plpgsql as $$ begin if new.event_type='training.lesson-read' then raise exception 'fixture rollback';end if;return new;end $$;create trigger reject_training_fixture before insert on platform.outbox_event for each row execute function audit.reject_training_fixture()`); e != nil {
		t.Fatal(e)
	}
	if _, e := f.store.Acknowledge(ctx, f.learner, id, "resource-create-view", hash); e == nil {
		t.Fatal("outbox failure accepted reading")
	}
	v, e := f.store.Read(ctx, f.learner, id)
	if e != nil || len(v.ReadLessons) != 0 {
		t.Fatal("reading fact survived rollback", e)
	}
	if _, e = f.pool.Exec(ctx, `drop trigger reject_training_fixture on platform.outbox_event;drop function audit.reject_training_fixture()`); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 2; i++ {
		if _, e = f.store.Acknowledge(ctx, f.learner, id, "resource-create-view", hash); e != nil {
			t.Fatal(e)
		}
	}
	a, e := f.store.Submit(ctx, f.learner, id, hash, answers)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.store.Submit(ctx, f.learner, id, hash, map[string]string{"recovery": "Changed immutable response"}); e == nil {
		t.Fatal("submission overwritten")
	}
	if _, e = f.store.Assess(ctx, f.learner, a.RequestID, a.PayloadSHA256, true, "Self review"); e == nil {
		t.Fatal("self review accepted")
	}
	if _, e = f.store.Assess(ctx, f.foreign, a.RequestID, a.PayloadSHA256, true, "Other organization"); e == nil {
		t.Fatal("foreign review accepted")
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, strings.Repeat("f", 64), true, "Wrong payload"); e == nil {
		t.Fatal("different payload assessed")
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, true, "Respuesta revisada contra el contenido de esta versión"); e != nil {
		t.Fatal(e)
	}
	restarted, e := NewStore(f.pool, f.profile)
	if e != nil {
		t.Fatal(e)
	}
	v, e = restarted.Read(ctx, f.learner, id)
	if e != nil || v.Assessment == nil || v.Assessment.State != "approved" || v.Assessment.Reviewer != "reviewer" {
		t.Fatal("assessment not durable", e)
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, false, "Changed decision"); e == nil {
		t.Fatal("decision changed")
	}
	var facts, events, requests, decisions, resources int
	e = f.pool.QueryRow(ctx, `select (select count(*) from audit.event where tenant_id=$1 and resource_type='training-attempt'),(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='training-evidence'),(select count(*) from approval.request where tenant_id=$1 and kind='training_assessment'),(select count(*) from approval.decision where tenant_id=$1),(select count(*) from crm.service_resource where tenant_id=$1)`, f.learner.TenantID).Scan(&facts, &events, &requests, &decisions, &resources)
	if e != nil || facts != 4 || events != 4 || requests != 1 || decisions != 1 || resources != 0 {
		t.Fatal("unexpected durable effects", facts, events, requests, decisions, resources, e)
	}
	if !f.learner.Allowed("training:learn") || f.learner.Allowed("resource:manage") {
		t.Fatal("training changed permissions")
	}
	if _, e = f.pool.Exec(ctx, `update approval.request set payload=jsonb_set(payload,'{answers,recovery}','"overwrite"') where tenant_id=$1 and request_id=$2`, f.learner.TenantID, a.RequestID); e == nil {
		t.Fatal("training payload update accepted")
	}
	t.Log("TRAINING_BINDINGS_PASS concurrent_start_one_fact=true reading_outbox_rollback=true own_subject_org_hash_bound=true human_review_durable=true duplicate_submission_immutable=true permission_grants=0 resource_creation=0")
}
````

### FILE: `internal/trainingbridge/guard_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "547ad406de386cfd17ee06d6d915997ff2f2c35dbc698a4f7fdea2b94ef1189d"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
)

func TestTrainingGenericBypassAndAuditImmutability(t *testing.T) {
	f := connectedFixture(t)
	ctx := context.Background()
	id := uuid.NewString()
	hash := f.profile.Hash()
	if _, e := f.store.Start(ctx, f.learner, id, "resource-onboarding", hash); e != nil {
		t.Fatal(e)
	}
	if _, e := f.store.Acknowledge(ctx, f.learner, id, "resource-create-view", hash); e != nil {
		t.Fatal(e)
	}
	a, e := f.store.Submit(ctx, f.learner, id, hash, map[string]string{"recovery": "Consulto la referencia antes de repetir."})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.store.reviews.Decide(ctx, f.reviewer, f.learner.TenantID, a.RequestID, "store-1", a.PayloadSHA256, true, "Generic bypass", "training:review", nil); e == nil {
		t.Fatal("generic decision bypassed training fact")
	}
	current, e := f.store.Assessment(ctx, f.learner, a.RequestID)
	if e != nil || current.State != approval.StatePending {
		t.Fatal("generic decision leaked", current.State, e)
	}
	raw, _ := json.Marshal(a.Payload)
	_, e = f.store.reviews.Submit(ctx, f.learner, postgres.HumanApprovalSpec{Request: approval.Request{TenantID: f.learner.TenantID, ID: "duplicate-training-identity", Kind: approval.KindTrainingAssessment, SubjectID: f.learner.Subject, Requester: f.learner.Subject, EvidenceSHA: a.PayloadSHA256}, OrganizationID: "store-1", Payload: raw}, "training:learn", nil)
	if e == nil {
		t.Fatal("one attempt acquired a second assessment identity")
	}
	for _, q := range []string{`update audit.event set evidence='{}' where tenant_id=$1 and resource_type='training-attempt'`, `delete from audit.event where tenant_id=$1 and resource_type='training-attempt'`} {
		if _, e = f.pool.Exec(ctx, q, f.learner.TenantID); e == nil {
			t.Fatal("training facts altered")
		}
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, true, "Distinct human review recorded"); e != nil {
		t.Fatal(e)
	}
	var requests, decisions, facts int
	e = f.pool.QueryRow(ctx, `select(select count(*) from approval.request where tenant_id=$1),(select count(*) from approval.decision where tenant_id=$1),(select count(*) from audit.event where tenant_id=$1 and resource_type='training-attempt')`, f.learner.TenantID).Scan(&requests, &decisions, &facts)
	if e != nil || requests != 1 || decisions != 1 || facts != 4 {
		t.Fatal("rejected mutation left effects", requests, decisions, facts, e)
	}
	t.Log("TRAINING_PHYSICAL_GUARDS_PASS generic decision without fact rejected; duplicate assessment refused; audit update/delete denied; normal human review preserved")
}
````

### FILE: `internal/trainingbridge/http.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e539425cb7332126e3591ee6d75e97fb5084b387d6e7edf13713006bd1a63dd3"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Module struct{ Store *Store }

func trainingJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func trainingError(w http.ResponseWriter, e error) {
	code, status := "UNAVAILABLE", 503
	switch {
	case errors.Is(e, ErrScope):
		code, status = "NOT_FOUND", 404
	case errors.Is(e, ErrContract):
		code, status = "INVALID_CONTRACT", 400
	case errors.Is(e, ErrConflict), errors.Is(e, approval.ErrDuplicate), errors.Is(e, approval.ErrNotPending), errors.Is(e, approval.ErrSeparation):
		code, status = "CONSULT_RECORDED_STATE", 409
	}
	trainingJSON(w, status, map[string]string{"code": code})
}
func (m Module) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier, permission string) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		trainingJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return identity.Principal{}, false
	}
	p, e := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if e != nil {
		trainingJSON(w, 401, map[string]string{"code": "UNAUTHENTICATED"})
		return p, false
	}
	allowed := m.Store != nil && (m.Store.allowed(p, permission) || (permission == "training:view" && (m.Store.allowed(p, "training:learn") || m.Store.allowed(p, "training:review"))))
	if !allowed {
		trainingJSON(w, 403, map[string]string{"code": "FORBIDDEN"})
		return p, false
	}
	if r.URL.RawQuery != "" {
		trainingJSON(w, 400, map[string]string{"code": "INVALID_QUERY"})
		return p, false
	}
	return p, true
}
func trainingBody(w http.ResponseWriter, r *http.Request, out any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrContract
	}
	b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		return ErrContract
	}
	return strict(b, out)
}
func (m Module) Register(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("GET /v1/training/courses", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Courses(p)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("GET /v1/training/assessments", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Assessments(r.Context(), p)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("GET /v1/training/assessments/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Assessment(r.Context(), p, r.PathValue("id"))
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("GET /v1/training/attempts/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:view")
		if !ok {
			return
		}
		value, e := m.Store.Read(r.Context(), p, r.PathValue("id"))
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("POST /v1/training/attempts", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:learn")
		if !ok {
			return
		}
		var body struct {
			ID         string `json:"attempt_id"`
			Course     string `json:"course_id"`
			ProfileSHA string `json:"profile_sha256"`
		}
		if e := trainingBody(w, r, &body); e != nil {
			trainingError(w, e)
			return
		}
		value, e := m.Store.Start(r.Context(), p, body.ID, body.Course, body.ProfileSHA)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 201, value)
	})
	mux.HandleFunc("POST /v1/training/attempts/{id}/acknowledgements", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:learn")
		if !ok {
			return
		}
		var body struct {
			Lesson     string `json:"lesson_id"`
			ProfileSHA string `json:"profile_sha256"`
		}
		if e := trainingBody(w, r, &body); e != nil {
			trainingError(w, e)
			return
		}
		value, e := m.Store.Acknowledge(r.Context(), p, r.PathValue("id"), body.Lesson, body.ProfileSHA)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("POST /v1/training/attempts/{id}/submissions", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:learn")
		if !ok {
			return
		}
		var body struct {
			Answers    map[string]string `json:"answers"`
			ProfileSHA string            `json:"profile_sha256"`
		}
		if e := trainingBody(w, r, &body); e != nil {
			trainingError(w, e)
			return
		}
		value, e := m.Store.Submit(r.Context(), p, r.PathValue("id"), body.ProfileSHA, body.Answers)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
	mux.HandleFunc("POST /v1/training/assessments/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "training:review")
		if !ok {
			return
		}
		var body struct {
			SHA      string `json:"payload_sha256"`
			Approved *bool  `json:"approved"`
			Reason   string `json:"reason"`
		}
		if e := trainingBody(w, r, &body); e != nil || body.Approved == nil {
			trainingError(w, ErrContract)
			return
		}
		value, e := m.Store.Assess(r.Context(), p, r.PathValue("id"), body.SHA, *body.Approved, body.Reason)
		if e != nil {
			trainingError(w, e)
			return
		}
		trainingJSON(w, 200, value)
	})
}
````

### FILE: `internal/trainingbridge/pack_profile_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5475fab73fd2d7afc5abead19aaed3fbc0e2a6db5b548cf17424a38baf40ce98"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTrainingMaterializedProfile(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/training/reference.profile.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc ProfileDocument
	if e = json.Unmarshal(raw, &doc); e != nil {
		t.Fatal(e)
	}
	expected, b := profileFixture(t)
	if !reflect.DeepEqual(doc, expected) {
		t.Fatal("materialized profile differs from exercised curriculum")
	}
	path, _ := filepath.Abs("../../deploy/training/reference.profile.json")
	content, _ := filepath.Abs("../../training_content/help.bundle.json")
	p, e := LoadProfile(path, content, Activation{true, doc.ID, doc.Revision, digest(raw), doc.TenantID, doc.OrganizationID})
	if e != nil || len(p.Courses()) != 5 || digest(b) != doc.ContentSHA256 {
		t.Fatalf("materialized profile: %v", e)
	}
}
````

### FILE: `internal/trainingbridge/profile.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c205ee46c56e5133758e0c0d05c4044051d67e74af27bca520790eeff2e016a1"
variables: []
secrets_allowed: false
```

````go
// Package trainingbridge composes versioned help, participation evidence and
// the existing human approval owner. It never scores or grants permissions.
package trainingbridge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"elite.local/enterprise/internal/approval"
)

var ErrContract = errors.New("training: invalid or inactive contract")
var ErrScope = errors.New("training: not authorized or not found")
var ErrConflict = errors.New("training: recorded identity or content differs")
var codeRE = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
var hexRE = regexp.MustCompile(`^[0-9a-f]{64}$`)

const Method = "HUMAN_REVIEW_NO_GRANTS_V1"

type Article struct {
	ID         string   `json:"id"`
	Version    string   `json:"version"`
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}
type ContentBundle struct {
	Schema       string    `json:"schema"`
	SourceSHA256 string    `json:"source_sha256"`
	Articles     []Article `json:"articles"`
}
type Prompt struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type Course struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Role    string   `json:"role"`
	Lessons []string `json:"lessons"`
	Prompts []Prompt `json:"prompts"`
}
type ProfileDocument struct {
	Schema         string   `json:"schema"`
	ID             string   `json:"id"`
	Revision       int      `json:"revision"`
	TenantID       string   `json:"tenant_id"`
	OrganizationID string   `json:"organization_id"`
	Method         string   `json:"method"`
	ContentSHA256  string   `json:"content_sha256"`
	SourceSHA256   string   `json:"source_sha256"`
	Courses        []Course `json:"courses"`
}
type Activation struct {
	Enabled        bool
	ID             string
	Revision       int
	SHA256         string
	TenantID       string
	OrganizationID string
}
type Profile struct {
	document ProfileDocument
	hash     string
	articles map[string]Article
}
type CourseView struct {
	Course          Course    `json:"course"`
	Articles        []Article `json:"articles"`
	ProfileID       string    `json:"profile_id"`
	ProfileRevision int       `json:"profile_revision"`
	ProfileSHA256   string    `json:"profile_sha256"`
	ContentSHA256   string    `json:"content_sha256"`
	Method          string    `json:"method"`
}

func digest(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func strict(raw []byte, out any) error {
	if len(raw) == 0 || len(raw) > 65536 || !utf8.Valid(raw) {
		return ErrContract
	}
	if _, _, e := approval.CanonicalPayload(raw); e != nil {
		return ErrContract
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(out) != nil {
		return ErrContract
	}
	var trailing any
	if dec.Decode(&trailing) != io.EOF {
		return ErrContract
	}
	return nil
}
func text(s string, max int) bool {
	if len(s) < 1 || len(s) > max || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if (r < 32 && r != '\n' && r != '\t') || r == 127 {
			return false
		}
	}
	return true
}
func readBounded(path string) ([]byte, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrContract
	}
	before, e := os.Lstat(path)
	if e != nil || !before.Mode().IsRegular() || before.Size() > 65536 {
		return nil, ErrContract
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() > 65536 || !os.SameFile(before, st) {
		return nil, ErrContract
	}
	b, e := io.ReadAll(io.LimitReader(f, 65537))
	if len(b) > 65536 {
		return nil, ErrContract
	}
	return b, e
}
func LoadProfile(path, contentPath string, a Activation) (*Profile, error) {
	raw, e := readBounded(path)
	if e != nil {
		return nil, e
	}
	content, e := readBounded(contentPath)
	if e != nil {
		return nil, e
	}
	return ParseProfile(raw, content, a)
}
func ParseProfile(raw, content []byte, a Activation) (*Profile, error) {
	var d ProfileDocument
	var b ContentBundle
	if !a.Enabled || !hexRE.MatchString(a.SHA256) || digest(raw) != a.SHA256 || strict(raw, &d) != nil || strict(content, &b) != nil {
		return nil, ErrContract
	}
	if d.Schema != "elite-training-profile/v1" || d.Method != Method || !codeRE.MatchString(d.ID) || d.ID != a.ID || d.Revision < 1 || d.Revision != a.Revision || d.TenantID != a.TenantID || d.OrganizationID != a.OrganizationID || !text(d.TenantID, 128) || !text(d.OrganizationID, 128) || !hexRE.MatchString(d.SourceSHA256) || digest(content) != d.ContentSHA256 || b.Schema != "elite-training-content/v1" || b.SourceSHA256 != d.SourceSHA256 || len(b.Articles) < 1 || len(b.Articles) > 64 || len(d.Courses) < 1 || len(d.Courses) > 8 {
		return nil, ErrContract
	}
	p := &Profile{document: d, hash: a.SHA256, articles: map[string]Article{}}
	for _, v := range b.Articles {
		if !codeRE.MatchString(v.ID) || !regexp.MustCompile(`^\d{1,6}\.\d{1,6}\.\d{1,6}$`).MatchString(v.Version) || !text(v.Title, 160) || len(v.Paragraphs) < 1 || len(v.Paragraphs) > 20 {
			return nil, ErrContract
		}
		if _, ok := p.articles[v.ID]; ok {
			return nil, ErrContract
		}
		for _, s := range v.Paragraphs {
			if !text(s, 4000) {
				return nil, ErrContract
			}
		}
		p.articles[v.ID] = v
	}
	seen := map[string]bool{}
	for _, c := range d.Courses {
		if !codeRE.MatchString(c.ID) || seen[c.ID] || !text(c.Title, 160) || len(c.Lessons) < 1 || len(c.Lessons) > 4 || len(c.Prompts) < 1 || len(c.Prompts) > 8 {
			return nil, ErrContract
		}
		seen[c.ID] = true
		switch c.Role {
		case "owner", "admin", "employee", "customer":
		default:
			return nil, ErrContract
		}
		lessons := map[string]bool{}
		for _, id := range c.Lessons {
			if _, ok := p.articles[id]; !ok || lessons[id] {
				return nil, ErrContract
			}
			lessons[id] = true
		}
		prompts := map[string]bool{}
		for _, q := range c.Prompts {
			if !codeRE.MatchString(q.ID) || !text(q.Text, 1000) || prompts[q.ID] {
				return nil, ErrContract
			}
			prompts[q.ID] = true
		}
		view, _ := p.Course(c.ID)
		encoded, _ := json.Marshal(view)
		if len(encoded) > 32768 {
			return nil, ErrContract
		}
	}
	return p, nil
}
func (p *Profile) Scope() (string, string) { return p.document.TenantID, p.document.OrganizationID }
func (p *Profile) Hash() string            { return p.hash }
func (p *Profile) Courses() []CourseView {
	v := []CourseView{}
	for _, c := range p.document.Courses {
		x, _ := p.Course(c.ID)
		v = append(v, x)
	}
	return v
}
func (p *Profile) Course(id string) (CourseView, error) {
	for _, c := range p.document.Courses {
		if c.ID == id {
			c.Lessons = append([]string(nil), c.Lessons...)
			c.Prompts = append([]Prompt(nil), c.Prompts...)
			v := CourseView{Course: c, ProfileID: p.document.ID, ProfileRevision: p.document.Revision, ProfileSHA256: p.hash, ContentSHA256: p.document.ContentSHA256, Method: Method}
			for _, id := range c.Lessons {
				article := p.articles[id]
				article.Paragraphs = append([]string(nil), article.Paragraphs...)
				v.Articles = append(v.Articles, article)
			}
			return v, nil
		}
	}
	return CourseView{}, ErrContract
}
````

### FILE: `internal/trainingbridge/profile_fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "293240da673d4eb235dd69eb3ba73938de731735a7452b29f3bb3b3a9a2e534a"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"encoding/json"
	"os"
	"testing"
)

func FuzzTrainingProfileIsolation(f *testing.F) {
	content, e := os.ReadFile("../../training_content/help.bundle.json")
	if e != nil {
		f.Fatal(e)
	}
	var bundle ContentBundle
	if e = json.Unmarshal(content, &bundle); e != nil {
		f.Fatal(e)
	}
	d := ProfileDocument{Schema: "elite-training-profile/v1", ID: "reference-onboarding", Revision: 1, TenantID: "50f38793-8a22-4f6b-983f-81dd0fca8208", OrganizationID: "store-1", Method: Method, ContentSHA256: digest(content), SourceSHA256: bundle.SourceSHA256, Courses: []Course{{ID: "resource-onboarding", Title: "Revisar referencia", Role: "employee", Lessons: []string{"resource-create-view"}, Prompts: []Prompt{{ID: "recovery", Text: "Explicá cómo consultarías el registro."}}}}}
	raw, _ := json.Marshal(d)
	f.Add(raw)
	f.Add([]byte("{}"))
	f.Add([]byte(`{"revision":1,"revision":2}`))
	f.Fuzz(func(t *testing.T, candidate []byte) {
		if len(candidate) > 65536 {
			return
		}
		p, e := ParseProfile(candidate, content, Activation{true, d.ID, d.Revision, digest(candidate), d.TenantID, d.OrganizationID})
		if e != nil {
			return
		}
		views := p.Courses()
		tenant, org := p.Scope()
		if p.Hash() != digest(candidate) || tenant != d.TenantID || org != d.OrganizationID || len(views) < 1 || len(views) > 8 {
			t.Fatal("accepted profile escaped activation")
		}
		for _, v := range views {
			if v.Method != Method || v.ProfileSHA256 != digest(candidate) || v.ContentSHA256 != digest(content) || len(v.Articles) < 1 || len(v.Articles) > 4 {
				t.Fatal("view binding changed")
			}
			before, _ := json.Marshal(v)
			v.Articles[0].Paragraphs[0] = "caller mutation"
			again, e := p.Course(v.Course.ID)
			if e != nil {
				t.Fatal(e)
			}
			actual, _ := json.Marshal(again)
			if string(before) != string(actual) {
				t.Fatal("consumer mutated frozen content")
			}
		}
	})
}
````

### FILE: `internal/trainingbridge/profile_test.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "57c207c18d3afdfb536c03ad0cc10da671060f7dd761b79cc5ceaae2d491e400"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func profileFixture(t *testing.T) (ProfileDocument, []byte) {
	t.Helper()
	b, e := os.ReadFile("../../training_content/help.bundle.json")
	if e != nil {
		t.Fatal(e)
	}
	var content ContentBundle
	if e = json.Unmarshal(b, &content); e != nil {
		t.Fatal(e)
	}
	return ProfileDocument{Schema: "elite-training-profile/v1", ID: "reference-onboarding", Revision: 2, TenantID: "50f38793-8a22-4f6b-983f-81dd0fca8208", OrganizationID: "store-1", Method: Method, ContentSHA256: digest(b), SourceSHA256: content.SourceSHA256, Courses: []Course{{ID: "resource-onboarding", Title: "Registrar recursos y recuperar una respuesta", Role: "employee", Lessons: []string{"resource-create-view"}, Prompts: []Prompt{{ID: "recovery", Text: "Describí qué harías si se pierde la respuesta después de registrar un recurso."}}}, {ID: "admin-onboarding", Title: "Revisión de recursos", Role: "admin", Lessons: []string{"resource-create-view", "supply-role-view", "warranty-role-view", "network-role-view"}, Prompts: []Prompt{{ID: "review", Text: "Describí cómo recuperarías una operación incierta, revisarías contenido y separarías una evaluación de los permisos de acceso."}}}, {ID: "owner-onboarding", Title: "Supervisión de la operación", Role: "owner", Lessons: []string{"operation-sections-view", "network-role-view", "training-role-view"}, Prompts: []Prompt{{ID: "review", Text: "Describí cómo distinguís una sección no disponible de una lista vacía."}}}, {ID: "customer-onboarding", Title: "Revisar una cotización", Role: "customer", Lessons: []string{"quote-acceptance-view"}, Prompts: []Prompt{{ID: "review", Text: "Describí qué confirma aceptar una cotización y qué se debe consultar ante una respuesta incierta."}}}, {ID: "content-onboarding", Title: "Contenido, publicación y revisión humana", Role: "admin", Lessons: []string{"help-cms-view", "catalog-role-view", "training-role-view"}, Prompts: []Prompt{{ID: "review", Text: "Explicá cómo publicarías contenido revisado, recuperarías una respuesta incierta y preservarías cursos anteriores sin conceder accesos."}}}}}, b
}
func activatedFixture(t *testing.T) (*Profile, []byte, []byte, Activation) {
	t.Helper()
	d, b := profileFixture(t)
	raw, _ := json.Marshal(d)
	a := Activation{true, d.ID, d.Revision, digest(raw), d.TenantID, d.OrganizationID}
	p, e := ParseProfile(raw, b, a)
	if e != nil {
		t.Fatal(e)
	}
	return p, raw, b, a
}
func TestTrainingProfileFixedContentAndNoGrants(t *testing.T) {
	p, raw, b, a := activatedFixture(t)
	if len(p.Courses()) != 5 {
		t.Fatal("role views")
	}
	for name, mutate := range map[string]func([]byte, []byte, Activation) ([]byte, []byte, Activation){
		"disabled":     func(x, y []byte, a Activation) ([]byte, []byte, Activation) { a.Enabled = false; return x, y, a },
		"wrong-tenant": func(x, y []byte, a Activation) ([]byte, []byte, Activation) { a.TenantID = "other"; return x, y, a },
		"wrong-org": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			a.OrganizationID = "other"
			return x, y, a
		},
		"wrong-revision":  func(x, y []byte, a Activation) ([]byte, []byte, Activation) { a.Revision++; return x, y, a },
		"content-altered": func(x, y []byte, a Activation) ([]byte, []byte, Activation) { return x, append(y, ' '), a },
		"policy-incompatible": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			x = []byte(strings.Replace(string(x), Method, "AUTO_SCORE_AND_GRANT", 1))
			a.SHA256 = digest(x)
			return x, y, a
		},
		"duplicate-key": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			x = []byte(strings.Replace(string(x), `"revision":2`, `"revision":2,"revision":2`, 1))
			a.SHA256 = digest(x)
			return x, y, a
		},
		"unknown-field": func(x, y []byte, a Activation) ([]byte, []byte, Activation) {
			x = append([]byte(`{"auto_approve":true,`), x[1:]...)
			a.SHA256 = digest(x)
			return x, y, a
		},
	} {
		t.Run(name, func(t *testing.T) {
			x, y, c := mutate(append([]byte(nil), raw...), append([]byte(nil), b...), a)
			if _, e := ParseProfile(x, y, c); e == nil {
				t.Fatal("unsafe activation accepted")
			}
		})
	}
}
````

### FILE: `internal/trainingbridge/store.go`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a969dac0e661a28d0af5d05c79cb095f6725a69b4171cc9198b4d23964c61e52"
variables: []
secrets_allowed: false
```

````go
package trainingbridge

// AUTHORED durable adapter. Participation is an append-only audit projection;
// assessment and separation of duties belong to the shared approval owner.
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool    *pgxpool.Pool
	profile *Profile
	reviews *postgres.HumanApprovals
}
type Attempt struct {
	ID             string      `json:"attempt_id"`
	Learner        string      `json:"learner_subject"`
	OrganizationID string      `json:"organization_id"`
	View           CourseView  `json:"content"`
	ReadLessons    []string    `json:"read_lessons"`
	Assessment     *Assessment `json:"assessment,omitempty"`
	CurrentProfile bool        `json:"current_profile"`
}
type AssessmentPayload struct {
	Schema         string            `json:"schema"`
	AttemptID      string            `json:"attempt_id"`
	Learner        string            `json:"learner_subject"`
	OrganizationID string            `json:"organization_id"`
	Content        CourseView        `json:"content"`
	Answers        map[string]string `json:"answers"`
}
type Assessment struct {
	RequestID     string            `json:"request_id"`
	PayloadSHA256 string            `json:"payload_sha256"`
	State         approval.State    `json:"state"`
	Payload       AssessmentPayload `json:"payload"`
	Reviewer      string            `json:"reviewer,omitempty"`
	Reason        string            `json:"reason,omitempty"`
	Approved      *bool             `json:"approved,omitempty"`
}
type startFact struct {
	Schema         string     `json:"schema"`
	AttemptID      string     `json:"attempt_id"`
	OrganizationID string     `json:"organization_id"`
	Content        CourseView `json:"content"`
}
type recordFact struct {
	Schema         string `json:"schema"`
	OrganizationID string `json:"organization_id"`
	ProfileSHA256  string `json:"profile_sha256"`
	LessonID       string `json:"lesson_id,omitempty"`
	PayloadSHA256  string `json:"payload_sha256,omitempty"`
	Approved       *bool  `json:"approved,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

func NewStore(pool *pgxpool.Pool, profile *Profile) (*Store, error) {
	if pool == nil || profile == nil {
		return nil, ErrContract
	}
	return &Store{pool, profile, postgres.NewHumanApprovals(pool)}, nil
}
func (s *Store) allowed(p identity.Principal, permission string) bool {
	tenant, org := s.profile.Scope()
	return p.TenantID == tenant && text(p.Subject, 128) && p.Allowed(permission) && p.AllowedOrganization(org)
}
func safeAttempt(id string) bool {
	parsed, e := uuid.Parse(id)
	return e == nil && parsed.String() == id && parsed != uuid.Nil
}
func tuple(parts ...string) string { b, _ := json.Marshal(parts); return string(b) }
func eventID(parts ...string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(tuple(parts...))).String()
}
func assessmentID(tenant, org, attempt string) string {
	return "training:" + digest([]byte(tuple(tenant, org, attempt)))
}
func (s *Store) Courses(p identity.Principal) ([]CourseView, error) {
	if !s.allowed(p, "training:learn") && !s.allowed(p, "training:review") {
		return nil, ErrScope
	}
	return s.profile.Courses(), nil
}

func (s *Store) fact(ctx context.Context, tx pgx.Tx, p identity.Principal, attempt, action, suffix string, payload any) error {
	raw, e := json.Marshal(payload)
	if e != nil {
		return e
	}
	canonical, _, e := approval.CanonicalPayload(raw)
	if e != nil {
		return e
	}
	tenant, org := s.profile.Scope()
	id := eventID(tenant, org, attempt, action, suffix)
	tag, e := tx.Exec(ctx, `insert into audit.event(tenant_id,event_id,actor_subject,action,resource_type,resource_id,decision,evidence) values($1,$2,$3,$4,'training-attempt',$5,'allowed',$6) on conflict(tenant_id,event_id) do nothing`, tenant, id, p.Subject, action, attempt, canonical)
	if e != nil {
		return e
	}
	if tag.RowsAffected() == 0 {
		var same bool
		e = tx.QueryRow(ctx, `select actor_subject=$3 and action=$4 and resource_type='training-attempt' and resource_id=$5 and decision='allowed' and evidence=$6::jsonb from audit.event where tenant_id=$1 and event_id=$2`, tenant, id, p.Subject, action, attempt, canonical).Scan(&same)
		if e != nil || !same {
			return ErrConflict
		}
		return nil
	}
	eventPayload, _ := json.Marshal(map[string]any{"attempt_id": attempt, "organization_id": org, "audit_event_id": id, "action": action})
	_, e = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'training-evidence',$3,1,$4,1,clock_timestamp(),$5)`, tenant, id, attempt+suffix, action, eventPayload)
	return e
}

type queryer interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func (s *Store) readAttempt(ctx context.Context, q queryer, p identity.Principal, id string, learnerOnly bool) (Attempt, error) {
	if !safeAttempt(id) || (!s.allowed(p, "training:learn") && !s.allowed(p, "training:review")) {
		return Attempt{}, ErrScope
	}
	tenant, org := s.profile.Scope()
	var actor string
	var raw []byte
	e := q.QueryRow(ctx, `select actor_subject,evidence from audit.event where tenant_id=$1 and event_id=$2 and action='training.started' and resource_type='training-attempt' and resource_id=$3`, tenant, eventID(tenant, org, id, "training.started", ""), id).Scan(&actor, &raw)
	if e != nil {
		return Attempt{}, ErrScope
	}
	if actor != p.Subject && (learnerOnly || !s.allowed(p, "training:review")) {
		return Attempt{}, ErrScope
	}
	var start startFact
	if strict(raw, &start) != nil || start.Schema != "training-start/v1" || start.OrganizationID != org || start.AttemptID != id || start.Content.Method != Method || !hexRE.MatchString(start.Content.ProfileSHA256) {
		return Attempt{}, ErrContract
	}
	v := Attempt{ID: id, Learner: actor, OrganizationID: org, View: start.Content, ReadLessons: []string{}, CurrentProfile: start.Content.ProfileSHA256 == s.profile.Hash()}
	rows, e := q.Query(ctx, `select evidence from audit.event where tenant_id=$1 and actor_subject=$2 and resource_type='training-attempt' and resource_id=$3 and action='training.lesson-read' order by audit_sequence`, tenant, actor, id)
	if e != nil {
		return Attempt{}, e
	}
	defer rows.Close()
	valid := map[string]bool{}
	for _, lesson := range v.View.Course.Lessons {
		valid[lesson] = true
	}
	for rows.Next() {
		var data []byte
		if rows.Scan(&data) != nil {
			return Attempt{}, ErrContract
		}
		var f recordFact
		if strict(data, &f) != nil || f.Schema != "training-read-declaration/v1" || f.OrganizationID != org || f.ProfileSHA256 != v.View.ProfileSHA256 || !valid[f.LessonID] {
			return Attempt{}, ErrContract
		}
		v.ReadLessons = append(v.ReadLessons, f.LessonID)
		delete(valid, f.LessonID)
	}
	return v, rows.Err()
}
func (s *Store) Start(ctx context.Context, p identity.Principal, id, course, profileHash string) (Attempt, error) {
	if !s.allowed(p, "training:learn") || !safeAttempt(id) {
		return Attempt{}, ErrScope
	}
	if profileHash != s.profile.Hash() {
		return Attempt{}, ErrContract
	}
	view, e := s.profile.Course(course)
	if e != nil {
		return Attempt{}, e
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return Attempt{}, e
	}
	defer tx.Rollback(ctx)
	_, org := s.profile.Scope()
	e = s.fact(ctx, tx, p, id, "training.started", "", startFact{"training-start/v1", id, org, view})
	if e != nil {
		return Attempt{}, e
	}
	if e = tx.Commit(ctx); e != nil {
		return Attempt{}, e
	}
	return s.Read(ctx, p, id)
}
func (s *Store) Read(ctx context.Context, p identity.Principal, id string) (Attempt, error) {
	v, e := s.readAttempt(ctx, s.pool, p, id, false)
	if e != nil {
		return v, e
	}
	tenant, org := s.profile.Scope()
	a, e := s.Assessment(ctx, p, assessmentID(tenant, org, id))
	if e == nil {
		v.Assessment = &a
	} else if !errors.Is(e, ErrScope) {
		return v, e
	}
	return v, nil
}
func (s *Store) Acknowledge(ctx context.Context, p identity.Principal, id, lesson, profileHash string) (Attempt, error) {
	if !s.allowed(p, "training:learn") {
		return Attempt{}, ErrScope
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if e != nil {
		return Attempt{}, e
	}
	defer tx.Rollback(ctx)
	v, e := s.readAttempt(ctx, tx, p, id, true)
	if e != nil {
		return Attempt{}, e
	}
	if !v.CurrentProfile || profileHash != s.profile.Hash() {
		return Attempt{}, ErrContract
	}
	exists := false
	for _, x := range v.View.Course.Lessons {
		if x == lesson {
			exists = true
		}
	}
	if !exists {
		return Attempt{}, ErrContract
	}
	e = s.fact(ctx, tx, p, id, "training.lesson-read", ":"+lesson, recordFact{Schema: "training-read-declaration/v1", OrganizationID: v.OrganizationID, ProfileSHA256: profileHash, LessonID: lesson})
	if e != nil {
		return Attempt{}, e
	}
	if e = tx.Commit(ctx); e != nil {
		return Attempt{}, e
	}
	return s.Read(ctx, p, id)
}
func (s *Store) Submit(ctx context.Context, p identity.Principal, id, profileHash string, answers map[string]string) (Assessment, error) {
	if !s.allowed(p, "training:learn") {
		return Assessment{}, ErrScope
	}
	v, e := s.readAttempt(ctx, s.pool, p, id, true)
	if e != nil {
		return Assessment{}, e
	}
	if !v.CurrentProfile || profileHash != s.profile.Hash() || len(answers) != len(v.View.Course.Prompts) {
		return Assessment{}, ErrContract
	}
	for _, q := range v.View.Course.Prompts {
		if !text(answers[q.ID], 2048) || strings.TrimSpace(answers[q.ID]) == "" {
			return Assessment{}, ErrContract
		}
	}
	payload := AssessmentPayload{"training-assessment/v1", id, p.Subject, v.OrganizationID, v.View, answers}
	raw, e := json.Marshal(payload)
	if e != nil {
		return Assessment{}, e
	}
	canonical, hash, e := approval.CanonicalPayload(raw)
	if e != nil {
		return Assessment{}, ErrContract
	}
	tenant, org := s.profile.Scope()
	request := assessmentID(tenant, org, id)
	_, e = s.reviews.Submit(ctx, p, postgres.HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: request, Kind: approval.KindTrainingAssessment, SubjectID: p.Subject, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: canonical}, "training:learn", func(ctx context.Context, tx pgx.Tx) error {
		current, e := s.readAttempt(ctx, tx, p, id, true)
		if e != nil {
			return e
		}
		if !current.CurrentProfile || len(current.ReadLessons) != len(current.View.Course.Lessons) {
			return ErrConflict
		}
		return s.fact(ctx, tx, p, id, "training.submitted", "", recordFact{Schema: "training-submission/v1", OrganizationID: org, ProfileSHA256: profileHash, PayloadSHA256: hash})
	})
	if e != nil {
		return Assessment{}, e
	}
	return s.Assessment(ctx, p, request)
}
func (s *Store) Assessment(ctx context.Context, p identity.Principal, id string) (Assessment, error) {
	if !s.allowed(p, "training:learn") && !s.allowed(p, "training:review") {
		return Assessment{}, ErrScope
	}
	tenant, org := s.profile.Scope()
	var v Assessment
	var raw []byte
	var learner string
	e := s.pool.QueryRow(ctx, `select request_id,evidence_sha,state,payload,requester from approval.request where tenant_id=$1 and organization_id=$2 and request_id=$3 and kind='training_assessment'`, tenant, org, id).Scan(&v.RequestID, &v.PayloadSHA256, &v.State, &raw, &learner)
	if e != nil {
		return v, ErrScope
	}
	if learner != p.Subject && !s.allowed(p, "training:review") {
		return Assessment{}, ErrScope
	}
	_, hash, e := approval.CanonicalPayload(raw)
	if e != nil || hash != v.PayloadSHA256 || strict(raw, &v.Payload) != nil || v.Payload.Schema != "training-assessment/v1" || v.Payload.OrganizationID != org || v.Payload.Learner != learner || v.Payload.Content.Method != Method {
		return Assessment{}, ErrContract
	}
	if v.State != approval.StatePending {
		var approved bool
		e = s.pool.QueryRow(ctx, `select reviewer,approved,reason from approval.decision where tenant_id=$1 and request_id=$2 and (select count(*) from approval.decision where tenant_id=$1 and request_id=$2)=1`, tenant, id).Scan(&v.Reviewer, &approved, &v.Reason)
		if e != nil || v.Reviewer == learner || approved != (v.State == approval.StateApproved) {
			return Assessment{}, ErrContract
		}
		v.Approved = &approved
	}
	return v, nil
}
func (s *Store) Assess(ctx context.Context, p identity.Principal, id, expectedHash string, approved bool, reason string) (Assessment, error) {
	if !s.allowed(p, "training:review") || !text(reason, 2048) || strings.TrimSpace(reason) == "" {
		return Assessment{}, ErrScope
	}
	v, e := s.Assessment(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.PayloadSHA256 != expectedHash {
		return v, ErrConflict
	}
	tenant, org := s.profile.Scope()
	_, e = s.reviews.Decide(ctx, p, tenant, id, org, expectedHash, approved, reason, "training:review", func(ctx context.Context, tx pgx.Tx) error {
		return s.fact(ctx, tx, p, v.Payload.AttemptID, "training.assessed", "", recordFact{Schema: "training-human-assessment/v1", OrganizationID: org, ProfileSHA256: v.Payload.Content.ProfileSHA256, PayloadSHA256: expectedHash, Approved: &approved, Reason: reason})
	})
	if e != nil {
		return Assessment{}, e
	}
	return s.Assessment(ctx, p, id)
}
func (s *Store) Assessments(ctx context.Context, p identity.Principal) ([]Assessment, error) {
	if !s.allowed(p, "training:learn") && !s.allowed(p, "training:review") {
		return nil, ErrScope
	}
	tenant, org := s.profile.Scope()
	rows, e := s.pool.Query(ctx, `select request_id from approval.request where tenant_id=$1 and organization_id=$2 and kind='training_assessment' and ($3::boolean or requester=$4) order by created_at desc,request_id limit 50`, tenant, org, s.allowed(p, "training:review"), p.Subject)
	if e != nil {
		return nil, e
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if e = rows.Scan(&id); e != nil {
			rows.Close()
			return nil, e
		}
		ids = append(ids, id)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return nil, e
	}
	out := []Assessment{}
	for _, id := range ids {
		v, e := s.Assessment(ctx, p, id)
		if e != nil {
			return nil, fmt.Errorf("assessment projection: %w", e)
		}
		out = append(out, v)
	}
	return out, nil
}
````

### FILE: `training_content/help.bundle.json`

```yaml
block_id: "GO-CONNECTED-HUMAN-TRAINING:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "eb919b9731d56832c38233572d6f1fd71cc62e7da64ca3d53a97638ab6a3ea20"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-training-content/v1",
  "source_sha256": "1c1ebc33ca2e578725f9d0ca4f51588f6f480e81ab9e8cbbdfc05b963395a581",
  "articles": [
    {
      "id": "quote-acceptance-view",
      "version": "1.0.0",
      "title": "Aceptar una cotización",
      "paragraphs": [
        "Revisá el importe, la moneda y la vigencia. Aceptar solicita crear un pedido con esos datos del servidor.",
        "Aceptar una cotización no confirma el pago, el stock ni la entrega.",
        "Si se corta la conexión, usá Actualizar estado o volvé a abrir Mis cotizaciones. No repitas la aceptación hasta comprobar el estado.",
        "Un pedido se confirma aquí cuando la lectura muestra su referencia. Si el problema continúa, contactá al soporte de la empresa con la referencia de cotización y organización, sólo por un canal autorizado; no compartas contraseñas ni tokens.",
        "Práctica de esta versión en el entorno de capacitación: reconocer un resultado incierto, recuperar la lectura y ubicar el pedido sin reenviar. No practiques creando pedidos en producción."
      ]
    },
    {
      "id": "whatsapp-status-view",
      "version": "1.0.0",
      "title": "Consultar el estado de WhatsApp",
      "paragraphs": [
        "Compartí esta referencia sólo con soporte autorizado de tu organización. No adjuntes teléfonos, tokens ni el contenido del mensaje.",
        "Si el resultado es incierto o contradictorio: preservá el intento, consultá su evidencia y no reenvíes ni borres el historial. La entrega informada no confirma una venta ni la aceptación del turno."
      ]
    },
    {
      "id": "availability-cancel-view",
      "version": "1.0.0",
      "title": "Cancelar un intervalo",
      "paragraphs": [
        "Cancelar este registro no cancela ni reprograma citas de clientes. El servidor comprueba los permisos, la versión y las restricciones existentes. Revisá recurso, fechas y tipo antes de actuar.",
        "Si se pierde la respuesta, Consultar intervalo sólo lee el servidor. Una cancelación consultada no demuestra quién la solicitó; la auditoría autorizada conserva al actor. La referencia de esta pestaña guarda sólo la versión, no motivos, horarios ni datos personales.",
        "Si la consulta falla, el intervalo no aparece en la ventana consultada o permanece activo, conservamos el bloqueo: solicitá al soporte autorizado revisar la referencia, hora y estado. No borres la referencia ni abras otra pestaña para forzar un reintento.",
        "Práctica en pruebas: cancelar un intervalo sin citas, perder la respuesta y consultar Cancelado sin reenviar. No practicar con la agenda real. Reactivar o crear otro intervalo es una operación distinta."
      ]
    },
    {
      "id": "availability-create-view",
      "version": "1.0.0",
      "title": "Registrar un intervalo",
      "paragraphs": [
        "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el registro aunque quede fuera del rango de la agenda. Preparar otro intervalo no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro intervalo."
      ]
    },
    {
      "id": "checklist-completion-view",
      "version": "1.0.0",
      "title": "Presentar una entrega",
      "paragraphs": [
        "La presentación conserva respuestas inmutables. Guardamos sólo identidades, versiones y una huella; no guardamos las respuestas ni evidencias en la referencia del navegador. Una respuesta perdida puede ocultar un cambio correcto: consultá sin reenviarlo. El estado actual puede avanzar a aceptado o rechazado; consultar o preparar otro formulario no lo cambia. Si el contenido no coincide, pedí revisión. Practicá con datos sintéticos antes de operar."
      ]
    },
    {
      "id": "checklist-publication-view",
      "version": "1.0.0",
      "title": "Publicar una checklist",
      "paragraphs": [
        "Una versión publicada es inmutable. Conservamos sólo su ID, versión y huella de contenido; no guardamos el título ni los textos en la referencia del navegador. Una respuesta perdida puede ocultar una publicación correcta: consultá sin volver a enviarla. Si no coincide el contenido o no se puede recuperar, pedí revisión sin borrar la referencia. Preparar otra versión no cambia ni vuelve a publicar la anterior. Practicá este recorrido con datos sintéticos antes de usarlo en operación."
      ]
    },
    {
      "id": "delivery-resolution-view",
      "version": "1.0.0",
      "title": "Resolver una discrepancia",
      "paragraphs": [
        "Una respuesta perdida no demuestra rechazo. Conservamos en esta pestaña sólo la versión y la decisión pendiente, sin notas ni credenciales. Consultar vuelve a leer el servidor sin reenviar la acción.",
        "Si la consulta falla, el registro no aparece o hay otra decisión, pedí revisión al responsable autorizado. No borres la referencia ni abras otra pestaña para forzar un reintento. Informá la referencia de la entrega y la hora, sin compartir notas privadas ni tokens.",
        "Práctica: en el entorno de prueba, recuperá una resolución tras perder la respuesta y comprobá su estado antes del siguiente paso. Preparada o autorizada no significa entrega física, reembolso ni cambio finalizados."
      ]
    },
    {
      "id": "handover-read-view",
      "version": "1.0.0",
      "title": "Consultar entregas",
      "paragraphs": [
        "Si la información no está disponible, volvé a consultar antes de aceptar o rechazar. No crees otra solicitud para recuperar una lectura.",
        "Si el problema continúa, contactá al soporte autorizado e indicá la acción y la hora. No envíes contraseñas, tokens ni información de otras personas.",
        "Práctica en pruebas: interrumpí la consulta, verificá que no haya acciones de entrega disponibles y restablecé la conexión. La próxima consulta debe recuperar el estado sin enviar una aceptación automática."
      ]
    },
    {
      "id": "lead-command-view",
      "version": "1.0.0",
      "title": "Actualizar un lead",
      "paragraphs": [
        "Una respuesta perdida no demuestra que el cambio haya sido rechazado. Consultar lee el servidor sin reenviar comandos. Una versión posterior permite revisar el estado real; no prueba que haya sido modificado exclusivamente por tu solicitud.",
        "Conservamos sólo versión y tipo de acción en esta pestaña, no responsables ni datos de contacto. Si falta el lead, falla su consulta o la versión no avanzó, pedí revisión al soporte autorizado; no borres la referencia ni abras otra pestaña para forzar un reintento.",
        "Práctica: recuperá un cambio en el entorno de prueba, comprobá responsable y estado y continuá con la versión consultada. No asumas una venta o un cliente convertido por una asignación."
      ]
    },
    {
      "id": "operation-sections-view",
      "version": "1.0.0",
      "title": "Recuperar una sección",
      "paragraphs": [
        "Una consulta fallida no significa que no existan registros. Volvé a consultar antes de decidir sobre esa sección; las secciones independientes conservan sus propios permisos y estados.",
        "La consulta no reenvía acciones. Si el problema persiste, informá al soporte autorizado la sección, la hora y la operación esperada; no compartas tokens, contraseñas ni datos personales en capturas.",
        "Práctica: ante una consulta fallida, identificá la sección, recuperala mediante consulta y verificá su estado antes de actuar."
      ]
    },
    {
      "id": "order-operations-view",
      "version": "1.1.0",
      "title": "Continuar un pedido",
      "paragraphs": [
        "Elegí una unidad por su serie y reservá una sola vez. El servidor vuelve a comprobar organización, variante, disponibilidad y versiones.",
        "Si se pierde la respuesta o alguien actuó primero, usá Actualizar pedidos. No crees otro pedido ni cambies identificadores para forzar la operación.",
        "Pago: registrar una solicitud guarda una intención del total; no confirma dinero recibido ni libera entrega. Una solicitud existente se consulta, no se duplica. Si falla, no borres la clave del navegador: consultá el estado y escalá al responsable antes de intentar otro pago.",
        "Práctica segura: en pruebas, dos operadores solicitan el pago del mismo pedido; sólo debe existir una intención inicial. Simulá respuesta perdida y recuperá el registro con Actualizar pedidos. No practiques con pedidos reales.",
        "Para soporte, comunicá la referencia del pedido o solicitud, acción y hora al responsable autorizado. No adjuntes tokens ni datos de clientes. Cobro, conciliación y entrega mantienen sus gates separados."
      ]
    },
    {
      "id": "quote-create-view",
      "version": "1.0.0",
      "title": "Emitir una cotización",
      "paragraphs": [
        "El servidor determina precio y moneda. La referencia conserva sólo un identificador de operación. Una respuesta perdida no permite emitir otra: consultá el resultado. Si no aparece, conservá la referencia y pedí revisión al responsable autorizado. Practicá emisión, respuesta perdida y consulta con datos sintéticos; no borres la referencia para repetir una cotización."
      ]
    },
    {
      "id": "resource-create-view",
      "version": "1.0.0",
      "title": "Registrar un recurso",
      "paragraphs": [
        "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta identifica el registro por la referencia de esta operación, sin buscar por nombre. Preparar otro recurso no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro recurso."
      ]
    },
    {
      "id": "return-operations-view",
      "version": "1.0.0",
      "title": "Recibir y decidir devoluciones",
      "paragraphs": [
        "Recepción, decisión y solicitudes conservan evidencia inmutable. La referencia del navegador guarda sólo operación, IDs y huella; no guarda serie, notas ni tokens. Una respuesta perdida puede ocultar un cambio correcto. Consultá el caso exacto aunque la lista falle o no lo muestre. Si la evidencia no coincide, pedí revisión. Continuar operando actualiza el caso consultado y habilita una acción nueva; no reenvía la anterior. Una devolución con reembolso crea cuatro solicitudes; un cambio crea tres y no solicita emisión fiscal. Practicá con datos sintéticos antes de operar."
      ]
    },
    {
      "id": "slot-create-view",
      "version": "1.0.0",
      "title": "Registrar un turno",
      "paragraphs": [
        "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el turno por su referencia, aunque esté cerrado, completo o fuera del listado público. Publicar no confirma una reserva; siguen vigentes la jornada y las restricciones de agenda. Preparar otro turno no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro turno."
      ]
    },
    {
      "id": "supply-role-view",
      "version": "1.0.0",
      "title": "Guía de suministro",
      "paragraphs": [
        "Compras registra proveedor, destino, fábrica, importe y cantidades. La fábrica confirma la orden, registra cada serie y pide revisión de calidad. Otra persona autorizada decide su liberación o rechazo con evidencia.",
        "El despacho identifica un envío y sus series. Recepción registra sólo lo recibido; queda en cuarentena hasta la decisión de otra persona autorizada. Las cantidades y versiones se comprueban en cada operación.",
        "Si se pierde una respuesta, consultá el resultado pendiente. Esa consulta no vuelve a enviar. No borres referencias para forzar un reintento. En capacitación, practicá con datos sintéticos el recorrido y la recuperación."
      ]
    },
    {
      "id": "warranty-role-view",
      "version": "1.0.0",
      "title": "Guía de garantía",
      "paragraphs": [
        "Leé los términos vinculados a la cotización antes de acusar su recepción. La garantía activa conserva esos términos vendidos. Un reclamo necesita entrega y atención registradas.",
        "Diagnóstico, plan y aprobación preceden al trabajo. La persona revisora es diferente de quien propone; calidad es diferente de quien ejecuta. El cliente acepta y la fábrica concilia el servicio. La conciliación es un acuse interno, no un pago.",
        "Si se pierde una respuesta, consultá el resultado pendiente: no vuelve a enviar. Conservá su referencia y pedí revisión autorizada si sigue incierto. Practicá con datos sintéticos en capacitación."
      ]
    },
    {
      "id": "network-role-view",
      "version": "1.0.0",
      "title": "Guía de red y acuerdos",
      "paragraphs": [
        "Registrá la organización bajo un padre autorizado y usá la referencia devuelta para consultar su estado. Crear una organización no concede accesos a personas. Los nuevos accesos se asignan mediante el proveedor de identidad autorizado.",
        "La estructura activa permite preparar acuerdos y sucursales. No equivale a autorización de despliegue o producción. Territorio, versión de términos y fechas deben proceder del acuerdo revisado; no se generan condiciones comerciales aquí.",
        "Antes de operar, completá la capacitación de tu rol. La evaluación no concede permisos. Si una respuesta se pierde, recuperá su resultado y luego consultá el estado actual. No repitas una escritura incierta.",
        "Para cerrar una organización, revisá primero sus sucursales y acuerdos activos. No se eliminan los registros ni se revocan sesiones desde este formulario."
      ]
    },
    {
      "id": "help-cms-view",
      "version": "1.0.0",
      "title": "Guía de contenidos",
      "paragraphs": [
        "El contenido pertenece a una organización y un idioma. Guardá un borrador, revisalo y publicalo con el permiso correspondiente. Archivar retira el artículo de la consulta de lectores y conserva su historial.",
        "Las versiones publicadas no se editan. Para reemplazarlas, prepará un artículo nuevo y archivá el anterior después de revisar la sustitución. El texto se muestra literalmente; no ejecuta HTML.",
        "Publicar aquí no modifica cursos aprobados ni concede permisos. La incorporación a capacitación requiere revisión y una nueva versión del curso."
      ]
    },
    {
      "id": "catalog-role-view",
      "version": "1.0.0",
      "title": "Guía de catálogo y publicación",
      "paragraphs": [
        "Creá modelo, variantes y precios con importes y vigencias explícitos. Conservá las referencias devueltas, incorporá la imagen PNG y prepará una versión del catálogo.",
        "La versión requiere las revisiones legal, técnica, de imagen y de publicación. Cada decisión refiere al contenido exacto y su evidencia. Publicar comprueba la generación vigente; una revisión desactualizada no autoriza contenido nuevo.",
        "Consultar un resultado pendiente recupera la operación sin repetirla. Antes de publicar otra versión o volver a una anterior, consultá la publicación vigente. La publicación del catálogo no confirma stock, cobro ni habilitación comercial. Practicá únicamente con datos sintéticos."
      ]
    },
    {
      "id": "training-role-view",
      "version": "1.0.0",
      "title": "Guía de capacitación y evaluación",
      "paragraphs": [
        "Elegí la guía del rol e iniciá una práctica. Registrá la lectura de cada lección y presentá respuestas para revisión humana. Elegir una guía no cambia los permisos de tu sesión.",
        "Una persona autorizada diferente revisa las respuestas contra la versión presentada. Una evaluación favorable conserva evidencia, pero no asigna accesos ni reemplaza la aceptación del responsable de la operación.",
        "Ante una respuesta incierta, consultá el intento guardado. Un curso nuevo no cambia la evidencia de intentos anteriores; prepará otra práctica para la versión activa. Los artículos privados del CMS requieren incorporación explícita a un curso revisado."
      ]
    }
  ]
}
````

## 6. Configuration surface

docs/TRAINING_REFERENCE.md and deploy/training/reference.profile.json bind optional host activation by exact profile/content SHA and absolute regular files; 2guards/3indexes required. Disabled means no profile read. Synthetic identity only, no credentials.

## 7. Dependency bill

AUTHORED unavoidable typed configuration, persistence, authorization, UI and orchestration/test glue; no new upstream/dependency or corporate attribution. Existing pgx/UUID, Go/PG, Next/React/Zod/OIDC/Playwright pins and notices retained.

## 8. Apply order

Materialize complete owner closure into absent target via MARKDOWN-COMPOSITOR0.3.0. Migration0075 follows current catalog73/74; all prior kinds preserved. Populated down refuses history; empty down/up verified. Historical0064 candidate must not be selected.

## 9. Verification

PG72migrations, concurrent start produces one fact; immutable answers/content and actor/org; self-review denied; generic approval/duplicate assessment and audit mutation blocked. Actual Next/BFF/Chromium learner+reviewer recover lost POST responses by GET:4facts/4outbox/1request/1decision/0resources/0grants. Host/path/hash/index/down proof, profile2s fuzz101053execs, static materialized curriculum matches exercised profile. Focused frontend3tests/full webpack build.

## 10. Reconstruction evidence

TRAINING_CONNECTED_RELEASE_V402.md/json binds source, before/after, RED locator and final PASS receipts, independent profile reconstruction and provenance. This closes training subclaim; other T2804 writes/CMS/private i18n remain open. No global readiness or production claim.


V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.
