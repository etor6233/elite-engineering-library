# Go Reminders Core

## 1. Metadata

```yaml
pack_id: "GO-REMINDERS-CORE"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Registro AUTHORED en memoria de recordatorios, con identidad (tenant,ID), Due(tenant) y MarkSent(tenant,ID) explícitos. Devuelve copias y conserva flag local; no reserva despacho ni demuestra envío, entrega exactamente una vez, persistencia o autorización."
stacks: ["Go 1.26.8"]
compatible_with: ["caller con tenant validado explícito; integración de canales no admitida"]
incompatible_with: ["API 0.1.x sin tenant", "consumo global de Due", "claim de envío exactamente una vez", "recordatorio en pasado", "canal no admitido"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia aislada para registrar tiempos/canal y consultar pendientes dentro de un tenant. El caller debe obtener el tenant de contexto autorizado; un parámetro string no autentica al actor. SMS/email/push son etiquetas validadas, no adapters ni efectos habilitados. No usar este registro como fence de envío.

## 3. Architecture contract

- **Ownership**: `internal/reminders` conserva registros locales; scheduling durable y envío son de los owners PostgreSQL/worker/outbound existentes.
- **Invariantes**: (1) At en el futuro. (2) due = at ≤ now y no enviado. (3) MarkSent es idempotente por (tenant,ID), sin efecto externo. (4) canal admitido. (5) reads/acks con tenant explícito. (6) Due no adquiere una reserva; lecturas concurrentes pueden repetir el mismo registro.
- **Data flow**: `Schedule` → `Due(tenant)` → decisión del caller → `MarkSent(tenant,id)`. No hay dispatch, validación de receipt, cancelación ni reconciliación en este core.
- **Failure modes**: pasado, inválido, duplicado, no encontrado → error.
- **Seguridad/privacidad**: aislamiento local por tenant exacto; caller debe autorizar actor/tenant y decidir acceso al contador global String de diagnóstico.
- **Performance budget**: O(1) schedule, O(N) due.

## 4. Exact file manifest

```text
CREATE internal/reminders/reminders.go
CREATE internal/reminders/reminders_test.go
```

## 5. Materialization blocks

### FILE: `internal/reminders/reminders.go`
```yaml
block_id: "GO-REMINDERS-CORE:internal/reminders/reminders.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0771a0a86db811051231a9094726803c4783210dae3e1964e45e9071e44bbd2e"
variables: []
secrets_allowed: false
```
````go
// Package reminders provides a tenant-scoped in-memory reminder registry.
// Due is a read, not a dispatch reservation; MarkSent only records a local flag.
// Callers must supply an authorized tenant and govern actual outbound effects.
package reminders

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidReminder = errors.New("reminders: invalid reminder")
	ErrDuplicate       = errors.New("reminders: duplicate reminder")
	ErrNotFound        = errors.New("reminders: not found")
	ErrPast            = errors.New("reminders: schedule time must be in the future")
)

// Channel is the dispatch surface.
type Channel string

const (
	ChannelSMS   Channel = "sms"
	ChannelEmail Channel = "email"
	ChannelPush  Channel = "push"
)

// Reminder is a single scheduled notification.
type Reminder struct {
	TenantID      string
	ID            string
	AppointmentID string
	At            time.Time
	Channel       Channel
	Sent          bool
}

type reminderID struct{ tenant, id string }

// Scheduler holds reminders.
type Scheduler struct {
	mu    sync.Mutex
	items map[reminderID]Reminder
	clock func() time.Time
}

// NewScheduler returns an empty scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{items: make(map[reminderID]Reminder), clock: time.Now}
}

func (s *Scheduler) now() time.Time {
	if s.clock != nil {
		return s.clock()
	}
	return time.Now()
}

// Schedule registers a reminder, fail-closed on invalid or duplicate.
func (s *Scheduler) Schedule(r Reminder) error {
	if strings.TrimSpace(r.TenantID) == "" || !idRe.MatchString(r.ID) || strings.TrimSpace(r.AppointmentID) == "" {
		return ErrInvalidReminder
	}
	switch r.Channel {
	case ChannelSMS, ChannelEmail, ChannelPush:
	default:
		return ErrInvalidReminder
	}
	if !r.At.After(s.now()) {
		return ErrPast
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[reminderID{r.TenantID, r.ID}]; ok {
		return ErrDuplicate
	}
	s.items[reminderID{r.TenantID, r.ID}] = r
	return nil
}

// Due returns copies of unsent due reminders only for the exact tenant.
// Repeated/concurrent reads can return the same reminder; no claim is acquired.
func (s *Scheduler) Due(tenant string) []Reminder {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	var out []Reminder
	for _, r := range s.items {
		if r.TenantID == tenant && !r.Sent && !r.At.After(now) {
			out = append(out, r)
		}
	}
	return out
}

// MarkSent sets the local sent flag for (tenant,id); already set is OK.
// It neither sends nor verifies a provider receipt.
func (s *Scheduler) MarkSent(tenant, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.items[reminderID{tenant, id}]
	if !ok {
		return ErrNotFound
	}
	r.Sent = true
	s.items[reminderID{tenant, id}] = r
	return nil
}

// String aids debugging.
func (s *Scheduler) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("reminders(%d)", len(s.items))
}
````

### FILE: `internal/reminders/reminders_test.go`
```yaml
block_id: "GO-REMINDERS-CORE:internal/reminders/reminders_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0054fd9b5af5282a36a0a80dfa5d836c69d1b392bfb4858cc291fdb7a1def8ba"
variables: []
secrets_allowed: false
```
````go
package reminders

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDueAndMarkSent(t *testing.T) {
	s := NewScheduler()
	now := time.Now()
	s.clock = func() time.Time { return now }

	// schedule r1 (10 min) and r2 (1 h) in the future
	if err := s.Schedule(Reminder{TenantID: "t", ID: "r1", AppointmentID: "a1", At: now.Add(10 * time.Minute), Channel: ChannelSMS}); err != nil {
		t.Fatal(err)
	}
	if err := s.Schedule(Reminder{TenantID: "t", ID: "r2", AppointmentID: "a2", At: now.Add(time.Hour), Channel: ChannelEmail}); err != nil {
		t.Fatal(err)
	}

	// advance clock past r1 only
	now = now.Add(30 * time.Minute)
	due := s.Due("t")
	if len(due) != 1 || due[0].ID != "r1" {
		t.Fatalf("expected only r1 due, got %+v", due)
	}
	if err := s.MarkSent("t", "r1"); err != nil {
		t.Fatal(err)
	}
	if len(s.Due("t")) != 0 {
		t.Fatal("r1 should not be due after MarkSent")
	}
}

func TestScheduleRejectsPast(t *testing.T) {
	s := NewScheduler()
	err := s.Schedule(Reminder{TenantID: "t", ID: "r1", AppointmentID: "a1", At: time.Now().Add(-time.Minute), Channel: ChannelSMS})
	if !errors.Is(err, ErrPast) {
		t.Fatalf("expected ErrPast, got %v", err)
	}
}

func TestScheduleRejectsInvalid(t *testing.T) {
	s := NewScheduler()
	if err := s.Schedule(Reminder{TenantID: "t", ID: "r1", AppointmentID: "a1", At: time.Now().Add(time.Hour), Channel: Channel("carrier")}); !errors.Is(err, ErrInvalidReminder) {
		t.Fatalf("bad channel accepted: %v", err)
	}
	if err := s.Schedule(Reminder{TenantID: "", ID: "r1", AppointmentID: "a1", At: time.Now().Add(time.Hour), Channel: ChannelSMS}); !errors.Is(err, ErrInvalidReminder) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
}

func TestDuplicateRejected(t *testing.T) {
	s := NewScheduler()
	r := Reminder{TenantID: "t", ID: "r1", AppointmentID: "a1", At: time.Now().Add(time.Hour), Channel: ChannelSMS}
	_ = s.Schedule(r)
	if err := s.Schedule(r); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", err)
	}
}

func TestMarkSentNotFound(t *testing.T) {
	s := NewScheduler()
	if err := s.MarkSent("t", "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func dueFor(s *Scheduler, tenant string) []Reminder { return s.Due(tenant) }
func markFor(s *Scheduler, tenant, id string) error { return s.MarkSent(tenant, id) }
func reminderFixture(t *testing.T) (*Scheduler, *time.Time) {
	t.Helper()
	now := time.Unix(1000000, 0).UTC()
	s := NewScheduler()
	s.clock = func() time.Time { return now }
	return s, &now
}
func TestSameReminderIDIndependentTenants(t *testing.T) {
	s, now := reminderFixture(t)
	for _, tenant := range []string{"a", "b"} {
		if e := s.Schedule(Reminder{TenantID: tenant, ID: "shared", AppointmentID: "appt", At: now.Add(time.Minute), Channel: ChannelEmail}); e != nil {
			t.Fatalf("tenant %s blocked by another: %v", tenant, e)
		}
	}
}
func TestDueNeverReturnsOtherTenant(t *testing.T) {
	s, now := reminderFixture(t)
	for _, tenant := range []string{"a", "b"} {
		if e := s.Schedule(Reminder{TenantID: tenant, ID: tenant, AppointmentID: "appt-" + tenant, At: now.Add(time.Minute), Channel: ChannelEmail}); e != nil {
			t.Fatal(e)
		}
	}
	*now = now.Add(time.Minute)
	for _, tenant := range []string{"a", "b", "unknown", ""} {
		for _, r := range dueFor(s, tenant) {
			if r.TenantID != tenant {
				t.Fatalf("requested %q got tenant %q", tenant, r.TenantID)
			}
		}
	}
}
func TestOtherTenantCannotMarkReminder(t *testing.T) {
	s, now := reminderFixture(t)
	if e := s.Schedule(Reminder{TenantID: "owner", ID: "private", AppointmentID: "appt", At: now.Add(time.Minute), Channel: ChannelPush}); e != nil {
		t.Fatal(e)
	}
	if e := markFor(s, "other", "private"); !errors.Is(e, ErrNotFound) {
		t.Fatalf("foreign mark accepted: %v", e)
	}
	*now = now.Add(time.Minute)
	if got := dueFor(s, "owner"); len(got) != 1 {
		t.Fatal("foreign mark changed owner state")
	}
}

var _ interface {
	Due(string) []Reminder
	MarkSent(string, string) error
} = (*Scheduler)(nil)

func TestScopedAcknowledgementAndDetachedDueReads(t *testing.T) {
	s, now := reminderFixture(t)
	for _, tenant := range []string{"a", "b"} {
		r := Reminder{TenantID: tenant, ID: "same", AppointmentID: "appt-" + tenant, At: now.Add(time.Minute), Channel: ChannelSMS}
		if e := s.Schedule(r); e != nil {
			t.Fatal(e)
		}
		r.AppointmentID = "replacement"
		if e := s.Schedule(r); !errors.Is(e, ErrDuplicate) {
			t.Fatal(e)
		}
	}
	*now = now.Add(time.Minute)
	first, again := s.Due("a"), s.Due("a")
	if len(first) != 1 || len(again) != 1 {
		t.Fatal("Due should be a read, not a claim")
	}
	first[0].TenantID = "b"
	first[0].AppointmentID = "mutated"
	first[0].Sent = true
	if r := s.Due("a"); len(r) != 1 || r[0].AppointmentID != "appt-a" || r[0].Sent {
		t.Fatal("caller can mutate registry")
	}
	if e := s.MarkSent("a", "same"); e != nil {
		t.Fatal(e)
	}
	if e := s.MarkSent("a", "same"); e != nil {
		t.Fatal("ack replay should succeed")
	}
	if len(s.Due("a")) != 0 || len(s.Due("b")) != 1 {
		t.Fatal("ack scope")
	}
	for _, tenant := range []string{"unknown", "", " "} {
		if len(s.Due(tenant)) != 0 {
			t.Fatal("unknown tenant leaked")
		}
		if e := s.MarkSent(tenant, "same"); !errors.Is(e, ErrNotFound) {
			t.Fatal("foreign ack")
		}
	}
}
func TestRejectedSchedulePreservesIdentityAndDueBoundary(t *testing.T) {
	for _, kind := range []string{"past", "equal", "zero", "channel", "appointment"} {
		t.Run(kind, func(t *testing.T) {
			s, now := reminderFixture(t)
			good := Reminder{TenantID: "t", ID: "r", AppointmentID: "a", At: now.Add(time.Second), Channel: ChannelSMS}
			bad := good
			want := ErrInvalidReminder
			switch kind {
			case "past":
				bad.At = now.Add(-time.Second)
				want = ErrPast
			case "equal":
				bad.At = *now
				want = ErrPast
			case "zero":
				bad.At = time.Time{}
				want = ErrPast
			case "channel":
				bad.Channel = "unsupported"
			case "appointment":
				bad.AppointmentID = " "
			}
			if e := s.Schedule(bad); !errors.Is(e, want) {
				t.Fatalf("rejection got=%v want=%v", e, want)
			}
			if e := s.Schedule(good); e != nil {
				t.Fatalf("rejection consumed identity: %v", e)
			}
			if len(s.Due("t")) != 0 {
				t.Fatal("premature due")
			}
			*now = good.At
			if len(s.Due("t")) != 1 {
				t.Fatal("at==now must be due")
			}
		})
	}
}
func TestConcurrentScopedScheduleAndAcknowledgement(t *testing.T) {
	s, now := reminderFixture(t)
	var wg sync.WaitGroup
	var a, b atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tenant := "a"
			count := &a
			if i%2 == 1 {
				tenant = "b"
				count = &b
			}
			e := s.Schedule(Reminder{TenantID: tenant, ID: "same", AppointmentID: "appt", At: now.Add(time.Minute), Channel: ChannelPush})
			if e == nil {
				count.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
		}(i)
	}
	wg.Wait()
	if a.Load() != 1 || b.Load() != 1 {
		t.Fatal("concurrent duplicate/global ID collision")
	}
	*now = now.Add(time.Minute)
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if e := s.MarkSent("a", "same"); e != nil {
				t.Error(e)
			}
			for _, r := range s.Due("b") {
				if r.TenantID != "b" {
					t.Error("foreign due")
				}
			}
			_ = s.String()
		}()
	}
	wg.Wait()
	if len(s.Due("a")) != 0 || len(s.Due("b")) != 1 {
		t.Fatal("cross-tenant concurrent ack")
	}
}

// The oracle stores a flat list and searches exact tenant/ID fields. It checks
// all exposed tenant views after each operation, independently of map keys.
func FuzzTenantReminderLifecycle(f *testing.F) {
	for _, ops := range [][]byte{{0, 0, 1, 0, 0, 1, 1, 0, 2, 0, 0, 1}, {0, 0, 1, 0, 0, 0, 1, 0}, {0, 0, 1, 0, 2, 0, 0, 2, 1, 1, 0, 0}, {0, 2, 1, 0, 0, 2, 1, 0}, {0, 0, 0, 0, 0, 0, 255, 0}, {0, 0, 1, 3, 2, 0, 0, 255}} {
		f.Add(ops)
	}
	f.Fuzz(func(t *testing.T, ops []byte) {
		if len(ops) > 128 {
			ops = ops[:128]
		}
		now := time.Unix(1000000, 0).UTC()
		s := NewScheduler()
		s.clock = func() time.Time { return now }
		tenants := []string{"a", "b", "a\x00b", " ", ""}
		ids := []string{"shared", "second", "bad id!", ""}
		channels := []Channel{ChannelSMS, ChannelEmail, ChannelPush, "bad"}
		model := []Reminder{}
		for i := 0; i+3 < len(ops); i += 4 {
			action := ops[i] % 3
			ti, ii := int(ops[i+1])%len(tenants), int(ops[i+2])%len(ids)
			tenant, id := tenants[ti], ids[ii]
			found := -1
			for j, r := range model {
				if r.TenantID == tenant && r.ID == id {
					found = j
					break
				}
			}
			if action == 0 {
				channel := channels[int(ops[i+3])%len(channels)]
				r := Reminder{TenantID: tenant, ID: id, AppointmentID: "appt", At: now.Add(time.Duration(int8(ops[i+2])) * time.Second), Channel: channel}
				want := error(nil)
				if ti >= 3 || ii >= 2 || channel == "bad" {
					want = ErrInvalidReminder
				} else if !r.At.After(now) {
					want = ErrPast
				} else if found >= 0 {
					want = ErrDuplicate
				}
				if e := s.Schedule(r); !errors.Is(e, want) {
					t.Fatalf("schedule classification got=%v want=%v", e, want)
				}
				if want == nil {
					model = append(model, r)
				}
			} else if action == 1 {
				want := error(nil)
				if found < 0 {
					want = ErrNotFound
				}
				if e := s.MarkSent(tenant, id); !errors.Is(e, want) {
					t.Fatalf("mark classification got=%v want=%v", e, want)
				}
				if want == nil {
					model[found].Sent = true
				}
			} else {
				now = now.Add(time.Duration(ops[i+3]) * time.Second)
			}
			for _, tn := range append(append([]string{}, tenants...), "unknown") {
				expected := map[string]Reminder{}
				for _, r := range model {
					if r.TenantID == tn && !r.Sent && !r.At.After(now) {
						expected[r.ID] = r
					}
				}
				got := s.Due(tn)
				if len(got) != len(expected) {
					t.Fatalf("due size for tenant %q", tn)
				}
				for _, r := range got {
					e, ok := expected[r.ID]
					if !ok || r != e {
						t.Fatal("due leaked/mutated record or repeated ID")
					}
					delete(expected, r.ID)
				}
				if len(expected) != 0 {
					t.Fatal("due missing")
				}
			}
		}
		_ = fmt.Sprint(s)
	})
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | scheduler | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/reminders/`.
3. Migrar cada `Due()` a `Due(tenant)` y cada `MarkSent(id)` a `MarkSent(tenant,id)` con tenant obtenido del contexto autorizado; no usar un tenant constante ni derivarlo de texto no confiable. No hay shim global.
4. Verificar con `go test ./... -count=1` y `go vet ./...`.
5. Rollback conserva0.1.0 sólo para diagnóstico: volver a su API global reabre FAIL624 y no es una opción segura de promoción.

## 9. Verification

V355:5tests originales conservados, adaptando sólo argumentos tenant de las dos
firmas.3regresiones de aislamiento fallan contra0.1.0 y pasan contra0.2.0.
11tests ordinarios x3,vet/build en candidato PASS; copia devuelta, frontera de
tiempo, rechazo sin consumir ID y64goroutines verificados. Caller0.1.x original
rechazado por compilación de ambas firmas. Reconstrucción2/2fuentes exactas;
GO_NATIVE_FUZZ_GATE10s/4workers/6semillas: 2112907ejecuciones PASS.
No se declara -race ni envío exactamente una vez. V199/21paquetes es histórico.
La cabecera antigua decía4tests, pero se ejecutaron los5presentes en la fuente.

## 10. Reconstruction evidence

V355: reconstruction_evidence/REMINDER_TENANT_API_V355.md.
Antes: reconstruction_evidence/GO_REMINDERS_CORE_2026-09-02_V199.md.

AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE y CONDITIONED persisten.
Fuera del perfil integral. Schedule conserva el Sent que aporte el caller;
ni ese flag ni MarkSent prueban aceptación del proveedor. No hay storage durable,
claim de job, delivery receipt, retry/fence, cancelación/reagenda ni autorización
de usuario. Reutilizar los owners reales de V293 y cerrar sus gates de proyecto.

## Auditoría por claim V374

La revisión 0.2.1 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.2.0 y sus bytes quedan preservados en el expediente anterior.
