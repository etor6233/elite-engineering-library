# Go Surveys Core

## 1. Metadata

```yaml
pack_id: "GO-SURVEYS-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa encuestas/NPS: respuestas con score 0-10, una por cliente por encuesta, tenant-scoped y NPS estándar (promotores 9-10, detractores 0-6, pasivos 7-8)."
stacks: ["Go 1.26.8"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-REVIEWS-CORE 0.1.x"]
incompatible_with: ["score fuera de 0-10", "respuesta duplicada"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria para respuestas y fórmula NPS existente. El caller autoriza tenant/encuesta/cliente; no valida compra, consentimiento, invitación, representatividad ni retención durable.

## 3. Architecture contract

- **Ownership**: `internal/surveys` gobierna respuestas y NPS.
- **Invariantes**: (1) score 0-10. (2) una respuesta por cliente por encuesta. (3) NPS = (promotores - detractores) / total * 100. (4) tenant-scoped.
- **Data flow**: `Submit` → `NPS`/`Responses`.
- **Failure modes**: score inválido, duplicado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) submit, O(N) NPS.

## 4. Exact file manifest

```text
CREATE internal/surveys/surveys.go
CREATE internal/surveys/surveys_test.go
```

## 5. Materialization blocks

### FILE: `internal/surveys/surveys.go`
```yaml
block_id: "GO-SURVEYS-CORE:internal/surveys/surveys.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cc7cc6be6f8ee6ff330911311060500228e090def644ead9979a433fa380665d"
variables: []
secrets_allowed: false
```
````go
// Package surveys provides tenant-scoped survey responses and NPS: a score
// 0-10, one response per customer per survey, and the standard NPS formula
// (promoters 9-10, detractors 0-6, passives 7-8).
package surveys

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidScore = errors.New("surveys: score must be 0..10")
	ErrDuplicate    = errors.New("surveys: customer already responded")
)

// Response is a single survey answer.
type Response struct {
	TenantID   string
	SurveyID   string
	CustomerID string
	Score      int
	At         time.Time
}

// Store holds responses.
type Store struct {
	mu        sync.Mutex
	responses map[identity]Response
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{responses: make(map[identity]Response)}
}

type identity struct{ tenant, survey, customer string }

func key(tenant, survey, customer string) identity {
	return identity{tenant: tenant, survey: survey, customer: customer}
}

// Submit records a response, one per customer per survey, score 0-10.
func (s *Store) Submit(r Response) error {
	if r.Score < 0 || r.Score > 10 {
		return ErrInvalidScore
	}
	if strings.TrimSpace(r.TenantID) == "" || strings.TrimSpace(r.SurveyID) == "" || strings.TrimSpace(r.CustomerID) == "" {
		return ErrInvalidScore
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.responses[key(r.TenantID, r.SurveyID, r.CustomerID)]; ok {
		return ErrDuplicate
	}
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	s.responses[key(r.TenantID, r.SurveyID, r.CustomerID)] = r
	return nil
}

// Responses returns the tenant-scoped responses for a survey.
func (s *Store) Responses(tenant, survey string) []Response {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Response
	for _, r := range s.responses {
		if r.TenantID == tenant && r.SurveyID == survey {
			out = append(out, r)
		}
	}
	return out
}

// NPS returns the Net Promoter Score (-100..100) over all responses.
func (s *Store) NPS(tenant, survey string) float64 {
	rs := s.Responses(tenant, survey)
	if len(rs) == 0 {
		return 0
	}
	promoters, detractors := 0, 0
	for _, r := range rs {
		switch {
		case r.Score >= 9:
			promoters++
		case r.Score <= 6:
			detractors++
		}
	}
	return float64(promoters-detractors) / float64(len(rs)) * 100
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("surveys(%d)", len(s.responses))
}
````

### FILE: `internal/surveys/surveys_test.go`
```yaml
block_id: "GO-SURVEYS-CORE:internal/surveys/surveys_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c2b2efd795a156f22874e6d3dc37054a300517227ab76ceb2c9cb9a155d165e2"
variables: []
secrets_allowed: false
```
````go
package surveys

import (
	"errors"
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func resp(score int) Response {
	return Response{TenantID: "t", SurveyID: "S1", CustomerID: "c", Score: score}
}

func TestNPSCalculation(t *testing.T) {
	s := NewStore()
	// 3 promoters (10,9,9), 2 detractors (0,6), 2 passives (7,8)
	for _, c := range []struct {
		id    string
		score int
	}{{"p1", 10}, {"p2", 9}, {"p3", 9}, {"d1", 0}, {"d2", 6}, {"n1", 7}, {"n2", 8}} {
		_ = s.Submit(Response{TenantID: "t", SurveyID: "S1", CustomerID: c.id, Score: c.score})
	}
	// (3 - 2) / 7 * 100 ≈ 14.29
	nps := s.NPS("t", "S1")
	if nps < 14 || nps > 15 {
		t.Fatalf("expected NPS ~14.29, got %v", nps)
	}
}

func TestNPSEmptyIsZero(t *testing.T) {
	s := NewStore()
	if s.NPS("t", "S1") != 0 {
		t.Fatalf("empty should be 0, got %v", s.NPS("t", "S1"))
	}
}

func TestInvalidScore(t *testing.T) {
	s := NewStore()
	if err := s.Submit(resp(-1)); !errors.Is(err, ErrInvalidScore) {
		t.Fatalf("score -1 accepted: %v", err)
	}
	if err := s.Submit(resp(11)); !errors.Is(err, ErrInvalidScore) {
		t.Fatalf("score 11 accepted: %v", err)
	}
}

func TestDuplicateCustomerRejected(t *testing.T) {
	s := NewStore()
	_ = s.Submit(resp(8))
	if err := s.Submit(resp(9)); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate customer accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Submit(resp(10))
	if len(s.Responses("other", "S1")) != 0 {
		t.Fatal("cross-tenant leak")
	}
}

func TestDistinctSurveyTuplesDoNotCollide(t *testing.T) {
	for _, pair := range [][2]Response{
		{{TenantID: "a\x00b", SurveyID: "c", CustomerID: "d", Score: 10}, {TenantID: "a", SurveyID: "b\x00c", CustomerID: "d", Score: 0}},
		{{TenantID: "a", SurveyID: "b\x00c", CustomerID: "d", Score: 10}, {TenantID: "a", SurveyID: "b", CustomerID: "c\x00d", Score: 0}},
	} {
		s := NewStore()
		if e := s.Submit(pair[0]); e != nil {
			t.Fatal(e)
		}
		if e := s.Submit(pair[1]); e != nil {
			t.Errorf("distinct response rejected: %v", e)
		}
		if s.NPS(pair[0].TenantID, pair[0].SurveyID) != 100 || s.NPS(pair[1].TenantID, pair[1].SurveyID) != -100 {
			t.Error("independent survey aggregate missing")
		}
	}
}

func TestSurveySnapshotAndScoreBoundaries(t *testing.T) {
	for score := 0; score <= 10; score++ {
		s := NewStore()
		if e := s.Submit(resp(score)); e != nil {
			t.Fatal(e)
		}
		want := float64(0)
		if score <= 6 {
			want = -100
		}
		if score >= 9 {
			want = 100
		}
		if s.NPS("t", "S1") != want {
			t.Fatal("score classification", score)
		}
		out := s.Responses("t", "S1")
		out[0].Score = 1000
		out[0].TenantID = "other"
		if s.Responses("t", "S1")[0].Score != score {
			t.Fatal("snapshot aliases store")
		}
	}
}
func TestConcurrentSurveyUniquenessAndDebugRead(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var created atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := resp(10)
			if i%2 == 0 {
				r.TenantID = "other"
			}
			e := s.Submit(r)
			if e == nil {
				created.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
			s.NPS(r.TenantID, "S1")
		}(i)
	}
	wg.Wait()
	if created.Load() != 2 {
		t.Fatal("tenant uniqueness", created.Load())
	}
}
func FuzzSurveyResponseModel(f *testing.F) {
	for _, x := range []string{"a", "a\x00b", "", " ", "客户", "b\x00c"} {
		f.Add(x, []byte{1, 0, 0, 11, 0, 1, 0, 1, 0, 0, 1, 7})
		f.Add(x, []byte{0, 0, 0, 1, 0, 0, 1, 8, 0, 0, 2, 10, 0, 0, 0, 11})
	}
	f.Fuzz(func(t *testing.T, x string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(ops) > 96 {
			ops = ops[:96]
		}
		tenants := []string{"a", "a\x00b", x, " "}
		surveys := []string{"c", "b\x00c", x, " "}
		customers := []string{"d", "c\x00d", x, " "}
		var model []Response
		s := NewStore()
		for i := 0; i+3 < len(ops); i += 4 {
			r := Response{TenantID: tenants[int(ops[i])%4], SurveyID: surveys[int(ops[i+1])%4], CustomerID: customers[int(ops[i+2])%4], Score: int(ops[i+3])%13 - 1}
			var want error
			if r.Score < 0 || r.Score > 10 || strings.TrimSpace(r.TenantID) == "" || strings.TrimSpace(r.SurveyID) == "" || strings.TrimSpace(r.CustomerID) == "" {
				want = ErrInvalidScore
			} else {
				for _, old := range model {
					if old.TenantID == r.TenantID && old.SurveyID == r.SurveyID && old.CustomerID == r.CustomerID {
						want = ErrDuplicate
						break
					}
				}
			}
			if e := s.Submit(r); !errors.Is(e, want) {
				t.Fatalf("got %v want %v", e, want)
			}
			if want == nil {
				model = append(model, r)
			}
			for _, tenant := range tenants {
				for _, survey := range surveys {
					var expected []Response
					net := int64(0)
					for _, r := range model {
						if r.TenantID == tenant && r.SurveyID == survey {
							expected = append(expected, r)
							if r.Score <= 6 {
								net--
							}
							if r.Score >= 9 {
								net++
							}
						}
					}
					actual := s.Responses(tenant, survey)
					if len(actual) != len(expected) {
						t.Fatal("response count mismatch")
					}
					seen := make([]bool, len(expected))
					for _, r := range actual {
						found := -1
						for j, e := range expected {
							if r.TenantID == e.TenantID && r.SurveyID == e.SurveyID && r.CustomerID == e.CustomerID && r.Score == e.Score {
								found = j
								break
							}
						}
						if found < 0 || seen[found] {
							t.Fatal("foreign/duplicate response")
						}
						seen[found] = true
					}
					value := float64(0)
					if len(expected) > 0 {
						value, _ = new(big.Rat).SetFrac(big.NewInt(net*100), big.NewInt(int64(len(expected)))).Float64()
					}
					nps := s.NPS(tenant, survey)
					if math.IsNaN(nps) || math.Abs(nps-value) > 1e-12 || nps < -100 || nps > 100 {
						t.Fatal("rational NPS mismatch", nps, value)
					}
				}
			}
		}
	})
}

// V374: source-class golden population; this does not validate survey collection.
func TestSourceNPSGoldenPopulationV374(t *testing.T) {
	s := NewStore()
	for score := 0; score <= 10; score++ {
		e := s.Submit(Response{TenantID: "reference", SurveyID: "all-classes", CustomerID: strings.Repeat("x", score+1), Score: score})
		if e != nil {
			t.Fatal(e)
		}
	}
	if got := s.NPS("reference", "all-classes"); math.Abs(got-(-500.0/11)) > 1e-12 {
		t.Fatal(got)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | encuestas | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/surveys/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/surveys/`.

## 9. Verification

- `go test ./internal/surveys/ -count=1`: 5/5 PASS (NPS cálculo, vacío 0, score inválido, duplicado, aislamiento tenant).
- `go test ./... -count=1` (22 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_SURVEYS_CORE_2026-09-02_V200.md`.

## Revisión V360

0.1.1: identidad estructurada y String sincronizado; marketing además exige tenant
y copia snapshots. AUTHORED/LicenseRef-Workspace-Owner y CONDITIONED persisten.
Suite V200 histórica no reejecutada; compatibilidad antigua no demuestra
integración actual. Evidencia: reconstruction_evidence/SURVEY_CAMPAIGN_SCOPE_V360.md.

V360 verificado entre ambos cores: 4/4 fuentes reconstruidas, 20 tests x3,
vet/build, 24 semillas/3864824 ejecuciones fuzz PASS.11tests históricos conservan
aserciones; marketing migra sólo argumentos tenant e índice privado del test.
Tres firmas sin tenant fallan compilación. Sin -race, dispatch ni consentimiento
real probado.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
