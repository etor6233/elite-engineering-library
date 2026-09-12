# Go Dashboards Core

## 1. Metadata

```yaml
pack_id: "GO-DASHBOARDS-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencia AUTHORED de agregación pura de KPIs suministrados por el caller: valida forma YYYY-MM y mes 01-12, conteos y scores; etiqueta tenant sin demostrar procedencia o autorización de los datos."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; adaptadores de origen y autorización aún no admitidos"]
incompatible_with: ["rating fuera de rango", "nps fuera de 0-10", "periodo inválido"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Calculadora local de snapshot: no consulta dominios ni comprueba que Input pertenezca al tenant. El caller debe probar origen, aislamiento, moneda y período de sus agregados. Rechaza mes fuera de 01-12, ratings y NPS inválidos.

## 3. Architecture contract

- **Ownership**: `internal/dashboards` agrega KPIs; el origen de cada métrica es su dominio (loyalty/reviews/surveys/pos).
- **Invariantes**: (1) tenant y periodo válidos. (2) conteos ≥ 0. (3) reviews 1-5, NPS 0-10. (4) agregación pura y fail-closed.
- **Data flow**: `Build(input)` → `Snapshot`.
- **Failure modes**: input inválido → error.
- **Seguridad/privacidad**: devuelve el tenant indicado; no hay autorización ni comprobación de pertenencia de Input.
- **Performance budget**: O(N) sobre reviews/NPS.

## 4. Exact file manifest

```text
CREATE internal/dashboards/dashboards.go
CREATE internal/dashboards/dashboards_test.go
```

## 5. Materialization blocks

### FILE: `internal/dashboards/dashboards.go`
```yaml
block_id: "GO-DASHBOARDS-CORE:internal/dashboards/dashboards.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c4ec3fb2c2e63340d1d6b0963a3d8ed30f1a57bf9423c0ab54f146db1c02773a"
variables: []
secrets_allowed: false
```
````go
// Package dashboards builds a pure KPI snapshot from caller-supplied values.
// The caller owns source queries, tenant authorization and period attribution;
// this package does not read domain stores or verify data ownership.
package dashboards

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	periodRe = regexp.MustCompile(`^\d{4}-\d{2}$`)

	ErrInvalidInput = errors.New("dashboards: invalid input")
)

// Input is the raw domain data for a period.
type Input struct {
	SalesMinorUnits int64
	Orders          int
	PointsIssued    int64
	Reviews         []int // ratings 1-5
	NPSScores       []int // 0-10
	Appointments    int
}

// Snapshot is the aggregated KPI view.
type Snapshot struct {
	TenantID        string
	Period          string
	SalesMinorUnits int64
	Orders          int
	PointsIssued    int64
	ReviewsAvg      float64
	ReviewsCount    int
	NPS             float64
	Appointments    int
}

// Build validates and aggregates the input into a snapshot.
func Build(tenant, period string, in Input) (Snapshot, error) {
	if strings.TrimSpace(tenant) == "" {
		return Snapshot{}, ErrInvalidInput
	}
	if !periodRe.MatchString(period) || period[5:] < "01" || period[5:] > "12" {
		return Snapshot{}, ErrInvalidInput
	}
	if in.SalesMinorUnits < 0 || in.Orders < 0 || in.PointsIssued < 0 || in.Appointments < 0 {
		return Snapshot{}, ErrInvalidInput
	}
	for _, r := range in.Reviews {
		if r < 1 || r > 5 {
			return Snapshot{}, ErrInvalidInput
		}
	}
	for _, s := range in.NPSScores {
		if s < 0 || s > 10 {
			return Snapshot{}, ErrInvalidInput
		}
	}

	s := Snapshot{
		TenantID: tenant, Period: period,
		SalesMinorUnits: in.SalesMinorUnits, Orders: in.Orders,
		PointsIssued: in.PointsIssued, ReviewsCount: len(in.Reviews),
		Appointments: in.Appointments,
	}
	if len(in.Reviews) > 0 {
		sum := 0
		for _, r := range in.Reviews {
			sum += r
		}
		s.ReviewsAvg = float64(sum) / float64(len(in.Reviews))
	}
	if len(in.NPSScores) > 0 {
		promoters, detractors := 0, 0
		for _, v := range in.NPSScores {
			switch {
			case v >= 9:
				promoters++
			case v <= 6:
				detractors++
			}
		}
		s.NPS = float64(promoters-detractors) / float64(len(in.NPSScores)) * 100
	}
	return s, nil
}

// String aids debugging.
func (s Snapshot) String() string {
	return fmt.Sprintf("dashboard(%s/%s)", s.TenantID, s.Period)
}
````

### FILE: `internal/dashboards/dashboards_test.go`
```yaml
block_id: "GO-DASHBOARDS-CORE:internal/dashboards/dashboards_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "235dafdda91235e9eac8c99a0334d3c0078dd9d8437b75a1f8249e84960efe34"
variables: []
secrets_allowed: false
```
````go
package dashboards

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestBuildSnapshot(t *testing.T) {
	s, err := Build("t", "2026-09", Input{
		SalesMinorUnits: 10000, Orders: 5, PointsIssued: 300,
		Reviews: []int{5, 4, 5}, NPSScores: []int{10, 9, 6}, Appointments: 7,
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.SalesMinorUnits != 10000 || s.Orders != 5 || s.PointsIssued != 300 || s.Appointments != 7 {
		t.Fatalf("unexpected snapshot: %+v", s)
	}
	// reviews avg = 14/3 ≈ 4.67
	if s.ReviewsAvg < 4.6 || s.ReviewsAvg > 4.7 {
		t.Fatalf("expected ~4.67 reviews avg, got %v", s.ReviewsAvg)
	}
	if s.ReviewsCount != 3 {
		t.Fatalf("expected 3 reviews, got %d", s.ReviewsCount)
	}
	// nps = (2 - 1) / 3 * 100 ≈ 33.33
	if s.NPS < 33 || s.NPS > 34 {
		t.Fatalf("expected NPS ~33.33, got %v", s.NPS)
	}
}

func TestEmptyReviewsAndNPS(t *testing.T) {
	s, err := Build("t", "2026-09", Input{SalesMinorUnits: 100, Orders: 1})
	if err != nil {
		t.Fatal(err)
	}
	if s.ReviewsAvg != 0 || s.NPS != 0 {
		t.Fatalf("empty reviews/nps should be 0, got %v/%v", s.ReviewsAvg, s.NPS)
	}
}

func TestInvalidInput(t *testing.T) {
	if _, err := Build("", "2026-09", Input{}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if _, err := Build("t", "bad", Input{}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("bad period accepted: %v", err)
	}
	if _, err := Build("t", "2026-09", Input{Reviews: []int{6}}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("rating 6 accepted: %v", err)
	}
	if _, err := Build("t", "2026-09", Input{NPSScores: []int{11}}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("nps 11 accepted: %v", err)
	}
}

func TestInvalidCalendarMonth(t *testing.T) {
	for _, period := range []string{"2026-00", "2026-13", "2026-99"} {
		s, e := Build("t", period, Input{})
		if !errors.Is(e, ErrInvalidInput) || s != (Snapshot{}) {
			t.Errorf("invalid month admitted: %q %+v %v", period, s, e)
		}
	}
}

func TestValidCalendarMonthsPreserveShape(t *testing.T) {
	for _, year := range []string{"0000", "2026", "9999"} {
		for month := 1; month <= 12; month++ {
			period := fmt.Sprintf("%s-%02d", year, month)
			s, e := Build(" t ", period, Input{})
			if e != nil || s.TenantID != " t " || s.Period != period {
				t.Fatal("valid month or identity rejected", s, e)
			}
		}
	}
	for _, period := range []string{"2026-1", "26-01", "2026-001", "2026-01\n", "２０２６-01", "2026/01", " 2026-01", "2026-01-01"} {
		s, e := Build("t", period, Input{})
		if !errors.Is(e, ErrInvalidInput) || s != (Snapshot{}) {
			t.Fatal("invalid shape admitted", period, s, e)
		}
	}
}
func TestNegativeAggregatesAndInvalidScores(t *testing.T) {
	for _, in := range []Input{{SalesMinorUnits: -1}, {Orders: -1}, {PointsIssued: -1}, {Appointments: -1}, {Reviews: []int{5, 0}}, {Reviews: []int{1, 6}}, {NPSScores: []int{10, -1}}, {NPSScores: []int{0, 11}}} {
		s, e := Build("t", "2026-01", in)
		if !errors.Is(e, ErrInvalidInput) || s != (Snapshot{}) {
			t.Fatal("invalid input admitted", s, e)
		}
	}
	s, e := Build("t", "2026-12", Input{SalesMinorUnits: math.MaxInt64, PointsIssued: math.MaxInt64})
	if e != nil || s.SalesMinorUnits != math.MaxInt64 || s.PointsIssued != math.MaxInt64 {
		t.Fatal("representable aggregates rejected", s, e)
	}
}
func TestSnapshotOwnsAggregatedValues(t *testing.T) {
	in := Input{Reviews: []int{5}, NPSScores: []int{10}}
	s, e := Build("t", "2026-01", in)
	if e != nil {
		t.Fatal(e)
	}
	in.Reviews[0] = 1
	in.NPSScores[0] = 0
	if s.ReviewsAvg != 5 || s.NPS != 100 {
		t.Fatal("snapshot changed after return")
	}
}
func FuzzCalendarAndKPIModel(f *testing.F) {
	for _, period := range []string{"2026-00", "2026-01", "2026-12", "2026-13", "2026-99", "0000-01", "9999-12", "bad", "2026-01\n"} {
		f.Add("t", period, int64(3), []byte{1, 5, 4}, []byte{0, 7, 9, 10})
		f.Add("t", period, int64(-1), []byte{0, 6}, []byte{11})
	}
	f.Fuzz(func(t *testing.T, tenant, period string, sales int64, reviews, nps []byte) {
		if len(tenant) > 128 || len(period) > 128 {
			return
		}
		if len(reviews) > 256 {
			reviews = reviews[:256]
		}
		if len(nps) > 256 {
			nps = nps[:256]
		}
		in := Input{SalesMinorUnits: sales, Orders: len(reviews), PointsIssued: 7, Appointments: len(nps)}
		valid := strings.TrimSpace(tenant) != "" && sales >= 0
		calendar := len(period) == 7
		if calendar {
			for i := 0; i < 4; i++ {
				calendar = calendar && period[i] >= '0' && period[i] <= '9'
			}
			calendar = calendar && period[4] == '-'
			months := map[string]bool{"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true, "11": true, "12": true}
			calendar = calendar && months[period[5:]]
		}
		valid = valid && calendar
		sum := int64(0)
		balance := int64(0)
		for _, v := range reviews {
			in.Reviews = append(in.Reviews, int(v))
			sum += int64(v)
			valid = valid && v >= 1 && v <= 5
		}
		for _, v := range nps {
			in.NPSScores = append(in.NPSScores, int(v))
			valid = valid && v <= 10
			if v >= 9 {
				balance++
			} else if v <= 6 {
				balance--
			}
		}
		s, e := Build(tenant, period, in)
		if !valid {
			if !errors.Is(e, ErrInvalidInput) || s != (Snapshot{}) {
				t.Fatal("invalid model input admitted", s, e)
			}
			return
		}
		if e != nil {
			t.Fatal("valid model input rejected", e)
		}
		avg, score := 0.0, 0.0
		if len(reviews) > 0 {
			avg, _ = new(big.Rat).SetFrac(big.NewInt(sum), big.NewInt(int64(len(reviews)))).Float64()
		}
		if len(nps) > 0 {
			score, _ = new(big.Rat).SetFrac(big.NewInt(balance*100), big.NewInt(int64(len(nps)))).Float64()
		}
		if math.Abs(s.ReviewsAvg-avg) > 1e-12 || math.Abs(s.NPS-score) > 1e-12 || s.TenantID != tenant || s.Period != period || s.SalesMinorUnits != sales || s.Orders != len(reviews) || s.PointsIssued != 7 || s.Appointments != len(nps) || s.ReviewsCount != len(reviews) {
			t.Fatal("KPI model mismatch", s)
		}
	})
}

// V374: independently tabulated Bain score classes, not a representativeness claim.
func TestSourceNPSGoldenClassesV374(t *testing.T) {
	for _, c := range []struct {
		scores []int
		want   float64
	}{
		{[]int{9, 10, 7, 8, 0, 1, 2, 3, 4, 5, 6}, -500.0 / 11},
		{[]int{9, 10}, 100}, {[]int{0, 6}, -100}, {[]int{7, 8}, 0},
	} {
		got, e := Build("reference", "2026-09", Input{NPSScores: c.scores})
		if e != nil {
			t.Fatal(e)
		}
		if math.Abs(got.NPS-c.want) > 1e-12 {
			t.Fatalf("got %v want %v", got.NPS, c.want)
		}
	}
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | KPIs | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/dashboards/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/dashboards/`.

## 9. Verification

- `go test ./internal/dashboards/ -count=1`: 3/3 PASS (snapshot con reviews/NPS, vacíos, inválido).
- `go test ./... -count=1` (29 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_DASHBOARDS_CORE_2026-09-02_V207.md`.

## Revisión V362

0.1.1 rechaza meses00/13/99; conserva YYYY-MM de cuatro dígitos, incluido año0000
antes admitido, sin nueva política de calendario fiscal. KPI puro y no fuente
de verdad: tenant no valida procedencia, moneda ni ventana de agregación.
La media suma int y requiere tamaño representable del corpus; fuzz se acota a
256scores, sin afirmar cálculo de slices ilimitados ni rendimiento productivo.
V207 conserva evidencia histórica; nueva calificación en
reconstruction_evidence/DASHBOARD_HELP_BOUNDARIES_V362.md.

V362: entre ambos cores4fuentes reconstruidas,19tests x3,8tests históricos
intactos,vet/build,34semillas y4266789ejecuciones fuzz PASS.
Sin -race, auth externa, fuente KPI conectada, búsqueda externa o publicación real.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
