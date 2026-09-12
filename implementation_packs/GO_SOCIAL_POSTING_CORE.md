# Go Social Posting Core

## 1. Metadata

```yaml
pack_id: "GO-SOCIAL-POSTING-CORE"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Agenda AUTHORED en memoria por tenant/id, consulta Due y registro local scheduled→posted/failed; no conexión, entrega, retry ni publicación real a redes."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; composición y target deben demostrar condiciones"]
incompatible_with: ["posteo en pasado", "canal no admitido", "posteo duplicado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia de agenda y acuses locales. El nombre del canal no prueba un adapter ni cuenta/proveedor; no usar como publicador productivo.

## 3. Architecture contract

- **Ownership**: `internal/social` gobierna el schedule; la publicación real es del adapter de red (CONDITIONED).
- **Invariantes**: (1) canal admitido, texto no vacío. (2) schedule futuro. (3) estados scheduled→posted/failed. (4) tenant-scoped.
- **Data flow**: `Schedule` → `Due` → red → `MarkPosted`/`MarkFailed`.
- **Failure modes**: pasado, inválido, duplicado, estado inválido → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) schedule, O(N) due.

## 4. Exact file manifest

```text
CREATE internal/social/social.go
CREATE internal/social/social_test.go
```

## 5. Materialization blocks

### FILE: `internal/social/social.go`
```yaml
block_id: "GO-SOCIAL-POSTING-CORE:internal/social/social.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "56bb88d4f04ff1e8e30bf6f9376603a02449f2cef7548c6e33721e43e5e56102"
variables: []
secrets_allowed: false
```
````go
// Package social provides tenant-scoped scheduled social-media posts: channel,
// text and schedule, with a state machine scheduled → posted/failed.
package social

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Channel is a social surface.
type Channel string

const (
	ChannelInstagram Channel = "instagram"
	ChannelFacebook  Channel = "facebook"
	ChannelTikTok    Channel = "tiktok"
	ChannelLinkedIn  Channel = "linkedin"
)

// State is the posting lifecycle.
type State string

const (
	StateScheduled State = "scheduled"
	StatePosted    State = "posted"
	StateFailed    State = "failed"
)

var (
	ErrInvalidPost = errors.New("social: invalid post")
	ErrDuplicate   = errors.New("social: duplicate post")
	ErrNotFound    = errors.New("social: not found")
	ErrBadState    = errors.New("social: bad state")
	ErrPast        = errors.New("social: schedule must be in the future")
)

// Post is a scheduled social post.
type Post struct {
	TenantID      string
	ID            string
	Channel       Channel
	Text          string
	ScheduledAt   time.Time
	State         State
	FailureReason string
}

// Store holds posts.
type Store struct {
	mu    sync.Mutex
	posts map[postKey]Post
	clock func() time.Time
}

type postKey struct{ tenant, id string }

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{posts: make(map[postKey]Post), clock: time.Now}
}

func (s *Store) now() time.Time {
	if s.clock != nil {
		return s.clock()
	}
	return time.Now()
}

// Schedule registers a post, fail-closed on invalid or past schedule.
func (s *Store) Schedule(p Post) error {
	if strings.TrimSpace(p.TenantID) == "" || strings.TrimSpace(p.ID) == "" {
		return ErrInvalidPost
	}
	switch p.Channel {
	case ChannelInstagram, ChannelFacebook, ChannelTikTok, ChannelLinkedIn:
	default:
		return ErrInvalidPost
	}
	if strings.TrimSpace(p.Text) == "" {
		return ErrInvalidPost
	}
	if !p.ScheduledAt.After(s.now()) {
		return ErrPast
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.posts[postKey{p.TenantID, p.ID}]; ok {
		return ErrDuplicate
	}
	p.State = StateScheduled
	s.posts[postKey{p.TenantID, p.ID}] = p
	return nil
}

// Due returns the specified tenant's scheduled posts whose time has arrived.
// This read does not reserve work or send content.
func (s *Store) Due(tenant string) []Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	var out []Post
	for _, p := range s.posts {
		if p.TenantID == tenant && p.State == StateScheduled && !p.ScheduledAt.After(now) {
			out = append(out, p)
		}
	}
	return out
}

// MarkPosted records a local acknowledgement, scheduled → posted.
// It does not validate delivery or require the scheduled time to have arrived.
func (s *Store) MarkPosted(tenant, id string) error {
	return s.transition(tenant, id, StateScheduled, StatePosted, "")
}

// MarkFailed marks scheduled → failed with a reason.
func (s *Store) MarkFailed(tenant, id, reason string) error {
	return s.transition(tenant, id, StateScheduled, StateFailed, reason)
}

func (s *Store) transition(tenant, id string, from, to State, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.posts[postKey{tenant, id}]
	if !ok {
		return ErrNotFound
	}
	if p.State != from {
		return ErrBadState
	}
	p.State = to
	p.FailureReason = reason
	s.posts[postKey{tenant, id}] = p
	return nil
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("social(%d)", len(s.posts))
}
````

### FILE: `internal/social/social_test.go`
```yaml
block_id: "GO-SOCIAL-POSTING-CORE:internal/social/social_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "49a04c495857257f9485979e1e28de065d0229a81edebdef007dace20eabfafd"
variables: []
secrets_allowed: false
```
````go
package social

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func post() Post {
	return Post{TenantID: "t", ID: "p1", Channel: ChannelInstagram, Text: "Nuevo producto", ScheduledAt: time.Now().Add(time.Hour)}
}

func TestScheduleDueAndMarkPosted(t *testing.T) {
	s := NewStore()
	now := time.Now()
	s.clock = func() time.Time { return now }

	if err := s.Schedule(post()); err != nil {
		t.Fatal(err)
	}
	if len(s.Due("t")) != 0 {
		t.Fatal("future post should not be due")
	}
	now = now.Add(2 * time.Hour)
	if len(s.Due("t")) != 1 {
		t.Fatal("post should be due after schedule")
	}
	if err := s.MarkPosted("t", "p1"); err != nil {
		t.Fatal(err)
	}
	if len(s.Due("t")) != 0 {
		t.Fatal("posted should not be due")
	}
}

func TestMarkFailed(t *testing.T) {
	s := NewStore()
	now := time.Now()
	s.clock = func() time.Time { return now }
	_ = s.Schedule(post())
	now = now.Add(2 * time.Hour)
	if err := s.MarkFailed("t", "p1", "rate limited"); err != nil {
		t.Fatal(err)
	}
	if len(s.Due("t")) != 0 {
		t.Fatal("failed should not be due")
	}
}

func TestRejectsPastAndInvalid(t *testing.T) {
	s := NewStore()
	past := post()
	past.ScheduledAt = time.Now().Add(-time.Minute)
	if err := s.Schedule(past); !errors.Is(err, ErrPast) {
		t.Fatalf("past accepted: %v", err)
	}
	bad := post()
	bad.Channel = "myspace"
	if err := s.Schedule(bad); !errors.Is(err, ErrInvalidPost) {
		t.Fatalf("bad channel accepted: %v", err)
	}
	empty := post()
	empty.Text = ""
	if err := s.Schedule(empty); !errors.Is(err, ErrInvalidPost) {
		t.Fatalf("empty text accepted: %v", err)
	}
}

func TestDuplicateRejected(t *testing.T) {
	s := NewStore()
	_ = s.Schedule(post())
	if err := s.Schedule(post()); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Schedule(post())
	other := post()
	other.TenantID = "other"
	other.ID = "p2"
	if err := s.Schedule(other); err != nil {
		t.Fatalf("different tenant should be allowed: %v", err)
	}
}

func TestSamePostIDAcrossTenants(t *testing.T) {
	s := NewStore()
	a, b := post(), post()
	b.TenantID = "other"
	if e := s.Schedule(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Schedule(b); e != nil {
		t.Fatal("other tenant cannot use its own ID", e)
	}
}
func TestDueHasNoTenantBoundary(t *testing.T) {
	s := NewStore()
	now := time.Unix(100, 0)
	s.clock = func() time.Time { return now }
	a, b := post(), post()
	a.ScheduledAt = now.Add(time.Second)
	b.ScheduledAt = a.ScheduledAt
	b.TenantID = "other"
	b.ID = "p2"
	if e := s.Schedule(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Schedule(b); e != nil {
		t.Fatal(e)
	}
	now = now.Add(time.Second)
	for _, p := range s.Due("t") {
		if p.TenantID != "t" {
			t.Error("due view cannot select requesting tenant", p.TenantID)
		}
	}
}
func TestMutationAndReadRequireTenantArgument(t *testing.T) {
	typ := reflect.TypeOf(NewStore())
	for name, n := range map[string]int{"Due": 2, "MarkPosted": 3, "MarkFailed": 4} {
		m, ok := typ.MethodByName(name)
		if !ok || m.Type.NumIn() != n {
			t.Error("API cannot express tenant scope", name)
		}
	}
}

func TestForeignPostTransitionsRejected(t *testing.T) {
	s := NewStore()
	if e := s.Schedule(post()); e != nil {
		t.Fatal(e)
	}
	if e := s.MarkPosted("other", "p1"); !errors.Is(e, ErrNotFound) {
		t.Fatal(e)
	}
	if e := s.MarkFailed("other", "p1", "x"); !errors.Is(e, ErrNotFound) {
		t.Fatal(e)
	}
	if e := s.MarkPosted("t", "p1"); e != nil {
		t.Fatal(e)
	}
}
func TestPostDueCopyAndEqualityBoundary(t *testing.T) {
	now := time.Unix(100, 0)
	s := NewStore()
	s.clock = func() time.Time { return now }
	p := post()
	p.ScheduledAt = now
	if e := s.Schedule(p); !errors.Is(e, ErrPast) {
		t.Fatal(e)
	}
	p.ScheduledAt = now.Add(time.Second)
	if e := s.Schedule(p); e != nil {
		t.Fatal(e)
	}
	now = p.ScheduledAt
	rows := s.Due("t")
	if len(rows) != 1 {
		t.Fatal("exact due time")
	}
	rows[0].Text = "changed"
	if s.Due("t")[0].Text != p.Text {
		t.Fatal("returned copy aliases")
	}
}
func TestConcurrentTenantPostIdentity(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for i := 0; i < 48; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p := post()
			p.TenantID = fmt.Sprint(i)
			if e := s.Schedule(p); e != nil {
				t.Error(e)
			}
			if e := s.MarkPosted(p.TenantID, p.ID); e != nil {
				t.Error(e)
			}
			_ = s.String()
		}(i)
	}
	wg.Wait()
	if len(s.posts) != 48 {
		t.Fatal("lost posts")
	}
}
func FuzzScopedScheduleHistory(f *testing.F) {
	for _, prefix := range []string{"a", "a\x00b", "ñ", " "} {
		for _, ops := range [][]byte{{0, 0, 0, 1, 4, 0, 1, 0, 2, 1}, {0, 0, 1, 0, 2, 0, 3, 0}, {0, 0, 0, 0, 4, 0, 3, 0}, {0, 1, 0, 2, 4, 0, 2, 1}} {
			f.Add(prefix, ops)
		}
	}
	f.Fuzz(func(t *testing.T, prefix string, ops []byte) {
		if len(prefix) > 64 {
			return
		}
		if len(ops) > 128 {
			ops = ops[:128]
		}
		now := time.Unix(100, 0)
		s := NewStore()
		s.clock = func() time.Time { return now }
		model := map[string]map[string]Post{}
		tenants := []string{"t" + prefix, "t" + prefix + "\x00b", "other"}
		ids := []string{"b\x00c", "c", "c"}
		for i := 0; i+1 < len(ops); i += 2 {
			op, slot := ops[i]%5, int(ops[i+1]%3)
			tenant, id := tenants[slot], ids[slot]
			if model[tenant] == nil {
				model[tenant] = map[string]Post{}
			}
			a, exists := model[tenant][id]
			var got, want error
			switch op {
			case 0:
				p := Post{TenantID: tenant, ID: id, Channel: ChannelInstagram, Text: "text", ScheduledAt: now.Add(time.Second)}
				if ops[i]&128 != 0 {
					p.ScheduledAt = now
				}
				if !p.ScheduledAt.After(now) {
					want = ErrPast
				} else if exists {
					want = ErrDuplicate
				} else {
					p.State = StateScheduled
					model[tenant][id] = p
				}
				got = s.Schedule(p)
			case 1, 2:
				if !exists {
					want = ErrNotFound
				} else if a.State != StateScheduled {
					want = ErrBadState
				} else {
					if op == 1 {
						a.State = StatePosted
						a.FailureReason = ""
					} else {
						a.State = StateFailed
						a.FailureReason = "synthetic"
					}
					model[tenant][id] = a
				}
				if op == 1 {
					got = s.MarkPosted(tenant, id)
				} else {
					got = s.MarkFailed(tenant, id, "synthetic")
				}
			case 3:
				got = s.MarkPosted("absent", id)
				want = ErrNotFound
			case 4:
				now = now.Add(time.Second)
			}
			if !errors.Is(got, want) {
				t.Fatal("transition model", got, want)
			}
			total := 0
			for mt, rows := range model {
				total += len(rows)
				due := map[string]Post{}
				for mi, m := range rows {
					found := false
					for _, p := range s.posts {
						if p.TenantID == mt && p.ID == mi {
							found = true
							if p != m {
								t.Fatal("post model mismatch")
							}
						}
					}
					if !found {
						t.Fatal("missing post")
					}
					if m.State == StateScheduled && !m.ScheduledAt.After(now) {
						due[mi] = m
					}
				}
				actual := s.Due(mt)
				if len(actual) != len(due) {
					t.Fatal("due count")
				}
				seen := map[string]bool{}
				for _, p := range actual {
					if p != due[p.ID] || p.TenantID != mt || seen[p.ID] {
						t.Fatal("due scope/model")
					}
					seen[p.ID] = true
				}
			}
			if len(s.posts) != total || len(s.Due("absent")) != 0 {
				t.Fatal("post cardinality/scope")
			}
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | posteos | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/social/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/social/`.

## 9. Verification

- `go test ./internal/social/ -count=1`: 5/5 PASS (schedule/due/posted, failed, pasado/inválido, duplicado, aislamiento tenant).
- `go test ./... -count=1` (27 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_SOCIAL_POSTING_CORE_2026-09-02_V205.md`.

## Migración obligatoria V363 / 0.2.0

Due(tenant), MarkPosted(tenant,id), MarkFailed(tenant,id,reason) requieren ámbito
explícito. No shim global. Due no reserva trabajo; transiciones no prueban acuse
del proveedor, no exigen tiempo vencido y no envían mensajes. Esas reglas previas
se conservan; caller debe demostrar auth, delivery durable, retry/idempotencia,
reconciliación y política editorial antes de componer. Misma ID entre tenants
es válida; duplicado del mismo tenant sigue ErrDuplicate. String sincronizado.
No orden de Due garantizado; sin persistencia y sin API pública de reloj inyectable.
Evidencia reconstruction_evidence/PUBLIC_CORE_SCOPE_V363.md.

V363: entre los tres cores6fuentes reconstruidas,32tests x3,15tests originales
conservados (social migra argumentos de tenant),vet/build,62semillas y
4518771ejecuciones fuzz PASS. Tres firmas sociales legacy rechazadas.
Sin -race, publicación externa ni admisión del proyecto.

## Auditoría por claim V374

La revisión 0.2.1 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.2.0 y sus bytes quedan preservados en el expediente anterior.
