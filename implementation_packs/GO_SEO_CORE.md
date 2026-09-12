# Go SEO Core

## 1. Metadata

```yaml
pack_id: "GO-SEO-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencia AUTHORED: lista URL por tenant exacto, robots sin CR/LF/NUL y JSON-LD con URL HTTP(S)/hostname; política local de meta60/160 code points. No sitemap XML, dominio autorizado ni certificación Google."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; composición y target deben demostrar condiciones"]
incompatible_with: ["URL no absoluta", "meta título/descripción fuera de rango"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://developers.google.com/search/docs"]
verified_at: "2026-09-10"
```

## 2. Applicability

Use para generar sitemap/robots/meta de la web pública. Rechace para URL no absoluta o meta fuera de rango.

## 3. Architecture contract

- **Ownership**: `internal/seo` emite el contrato SEO; el render de la página es del frontend.
- **Invariantes**: (1) URLs absolutas http(s). (2) sitemap sin duplicados. (3) title ≤60, description ≤160 como política local AUTHORED, no límite de Google. (4) JSON-LD válido.
- **Data flow**: `Add` (sitemap) / `Allow/Disallow` (robots) / `JSONLD` / `ValidateMeta`.
- **Failure modes**: URL/meta/JSON-LD inválido → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(N) por render.

## 4. Exact file manifest

```text
CREATE internal/seo/seo.go
CREATE internal/seo/seo_test.go
```

## 5. Materialization blocks

### FILE: `internal/seo/seo.go`
```yaml
block_id: "GO-SEO-CORE:internal/seo/seo.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d5aecce729bce9f8ab3717d81795cdebdc2aeb65f5a8b74169f1408c07a5ebff"
variables: []
secrets_allowed: false
```
````go
// Package seo provides a tenant-scoped SEO contract: a deduplicated sitemap,
// robots rules and LocalBusiness JSON-LD with meta length validation.
// AUTHORED local policy; meta length bounds are not Google requirements.
package seo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

var (
	ErrInvalidURL  = errors.New("seo: invalid URL")
	ErrInvalidMeta = errors.New("seo: meta length out of range")
	ErrInvalidBiz  = errors.New("seo: invalid local business")
)

// Sitemap holds deduplicated, tenant-scoped URLs in insertion order.
type Sitemap struct {
	mu    sync.Mutex
	urls  map[sitemapKey]bool
	order []sitemapKey
}

type sitemapKey struct{ tenant, url string }

// NewSitemap returns an empty sitemap.
func NewSitemap() *Sitemap {
	return &Sitemap{urls: make(map[sitemapKey]bool)}
}

// Add registers a URL (deduplicated). Absolute http(s) URLs only.
func (s *Sitemap) Add(tenant, rawURL string) error {
	if strings.TrimSpace(tenant) == "" {
		return ErrInvalidURL
	}
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return ErrInvalidURL
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := sitemapKey{tenant, u.String()}
	if s.urls[k] {
		return nil
	}
	s.urls[k] = true
	s.order = append(s.order, k)
	return nil
}

// URLs returns the tenant-scoped URLs in insertion order.
func (s *Sitemap) URLs(tenant string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for _, k := range s.order {
		if k.tenant == tenant {
			out = append(out, k.url)
		}
	}
	return out
}

// Robots holds per-tenant allow/disallow rules.
type Robots struct {
	mu         sync.Mutex
	allowed    map[string][]string
	disallowed map[string][]string
}

// NewRobots returns an empty robots registry.
func NewRobots() *Robots {
	return &Robots{allowed: make(map[string][]string), disallowed: make(map[string][]string)}
}

// Allow adds an allow rule for a tenant.
func (r *Robots) Allow(tenant, path string) error {
	if strings.TrimSpace(tenant) == "" || path == "" || strings.ContainsAny(path, "\r\n\x00") {
		return ErrInvalidURL
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.allowed[tenant] = append(r.allowed[tenant], path)
	return nil
}

// Disallow adds a disallow rule for a tenant.
func (r *Robots) Disallow(tenant, path string) error {
	if strings.TrimSpace(tenant) == "" || path == "" || strings.ContainsAny(path, "\r\n\x00") {
		return ErrInvalidURL
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.disallowed[tenant] = append(r.disallowed[tenant], path)
	return nil
}

// Render returns the robots.txt content for a tenant.
func (r *Robots) Render(tenant string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var b strings.Builder
	b.WriteString("User-agent: *\n")
	for _, p := range r.allowed[tenant] {
		b.WriteString("Allow: " + p + "\n")
	}
	for _, p := range r.disallowed[tenant] {
		b.WriteString("Disallow: " + p + "\n")
	}
	return b.String()
}

// LocalBusiness is a minimal JSON-LD LocalBusiness entity.
type LocalBusiness struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Telephone string `json:"telephone"`
	Address   string `json:"address,omitempty"`
}

// Validate enforces required fields.
func (b LocalBusiness) Validate() error {
	if strings.TrimSpace(b.Name) == "" || strings.TrimSpace(b.URL) == "" {
		return ErrInvalidBiz
	}
	u, err := url.Parse(b.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return ErrInvalidBiz
	}
	return nil
}

// JSONLD returns the JSON-LD script body (without the <script> tag).
func (b LocalBusiness) JSONLD() (string, error) {
	if err := b.Validate(); err != nil {
		return "", err
	}
	obj := map[string]any{
		"@context":  "https://schema.org",
		"@type":     "LocalBusiness",
		"name":      b.Name,
		"url":       b.URL,
		"telephone": b.Telephone,
	}
	if b.Address != "" {
		obj["address"] = b.Address
	}
	raw, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// ValidateMeta enforces this library's local policy: title 1..60 and
// description 1..160 Unicode code points after TrimSpace, not Google limits.
func ValidateMeta(title, description string) error {
	t := len([]rune(strings.TrimSpace(title)))
	d := len([]rune(strings.TrimSpace(description)))
	if t < 1 || t > 60 || d < 1 || d > 160 {
		return ErrInvalidMeta
	}
	return nil
}

// String aids debugging.
func (s *Sitemap) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("seo.sitemap(%d)", len(s.urls))
}
````

### FILE: `internal/seo/seo_test.go`
```yaml
block_id: "GO-SEO-CORE:internal/seo/seo_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b934a766ae605d2d7c7ea725b79e88fadc646f3e4406f19751beab9619ee3e44"
variables: []
secrets_allowed: false
```
````go
package seo

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestSitemapDedupAndOrder(t *testing.T) {
	s := NewSitemap()
	if err := s.Add("t", "https://a.com/x"); err != nil {
		t.Fatal(err)
	}
	if err := s.Add("t", "https://a.com/x"); err != nil {
		t.Fatal(err)
	}
	if err := s.Add("t", "https://a.com/y"); err != nil {
		t.Fatal(err)
	}
	got := s.URLs("t")
	if len(got) != 2 || got[0] != "https://a.com/x" || got[1] != "https://a.com/y" {
		t.Fatalf("unexpected urls: %v", got)
	}
}

func TestSitemapInvalidURL(t *testing.T) {
	s := NewSitemap()
	if err := s.Add("t", "not-a-url"); !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("bad url accepted: %v", err)
	}
	if err := s.Add("", "https://a.com/x"); !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
}

func TestRobotsRender(t *testing.T) {
	r := NewRobots()
	if err := r.Allow("t", "/public/"); err != nil {
		t.Fatal(err)
	}
	if err := r.Disallow("t", "/admin/"); err != nil {
		t.Fatal(err)
	}
	out := r.Render("t")
	if !strings.Contains(out, "Allow: /public/") || !strings.Contains(out, "Disallow: /admin/") {
		t.Fatalf("unexpected robots: %q", out)
	}
}

func TestJSONLD(t *testing.T) {
	b := LocalBusiness{Name: "Cafe X", URL: "https://cafex.com", Telephone: "+54 11 0000 0000", Address: "Calle 1"}
	out, err := b.JSONLD()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"@type":"LocalBusiness"`) || !strings.Contains(out, `"name":"Cafe X"`) {
		t.Fatalf("unexpected JSON-LD: %s", out)
	}
	if _, err := (LocalBusiness{Name: "", URL: "https://cafex.com"}).JSONLD(); !errors.Is(err, ErrInvalidBiz) {
		t.Fatalf("invalid biz accepted: %v", err)
	}
}

func TestValidateMeta(t *testing.T) {
	if err := ValidateMeta("OK", "desc"); err != nil {
		t.Fatalf("valid meta rejected: %v", err)
	}
	if err := ValidateMeta(strings.Repeat("x", 61), "desc"); !errors.Is(err, ErrInvalidMeta) {
		t.Fatalf("long title accepted: %v", err)
	}
	if err := ValidateMeta("OK", strings.Repeat("x", 161)); !errors.Is(err, ErrInvalidMeta) {
		t.Fatalf("long description accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewSitemap()
	if err := s.Add("t1", "https://a.com/x"); err != nil {
		t.Fatal(err)
	}
	if got := s.URLs("t2"); len(got) != 0 {
		t.Fatalf("sitemap leaked across tenants: %v", got)
	}
}

func TestSitemapExactTenant(t *testing.T) {
	s := NewSitemap()
	if e := s.Add("a\x00b", "https://example.test/path"); e != nil {
		t.Fatal(e)
	}
	if got := s.URLs("a"); len(got) != 0 {
		t.Fatal("prefix tenant leaked URLs", got)
	}
}
func TestRobotsRejectLineInjection(t *testing.T) {
	for _, path := range []string{"/ok\nDisallow: /", "/ok\rUser-agent: secret", "/ok\x00hidden"} {
		for _, allow := range []bool{true, false} {
			r := NewRobots()
			before := r.Render("t")
			var e error
			if allow {
				e = r.Allow("t", path)
			} else {
				e = r.Disallow("t", path)
			}
			if !errors.Is(e, ErrInvalidURL) || r.Render("t") != before {
				t.Error("injected rule accepted", path, e)
			}
		}
	}
}
func TestLocalBusinessAbsoluteURL(t *testing.T) {
	for _, raw := range []string{"relative/path", "javascript:alert(1)", "https://:80/path"} {
		s, e := (LocalBusiness{Name: "A", URL: raw}).JSONLD()
		if !errors.Is(e, ErrInvalidBiz) || s != "" {
			t.Error("invalid business URL accepted", raw, s, e)
		}
	}
}
func TestSitemapRejectsHostlessURL(t *testing.T) {
	s := NewSitemap()
	if e := s.Add("t", "https://:80/path"); !errors.Is(e, ErrInvalidURL) {
		t.Fatal("missing hostname accepted", e)
	}
}

func TestSitemapReturnedSliceCannotMutateStore(t *testing.T) {
	s := NewSitemap()
	if e := s.Add("t", "https://example.test/a"); e != nil {
		t.Fatal(e)
	}
	out := s.URLs("t")
	out[0] = "changed"
	if s.URLs("t")[0] != "https://example.test/a" {
		t.Fatal("caller changed store")
	}
}
func TestJSONLDEscapesScriptText(t *testing.T) {
	b := LocalBusiness{Name: "</script><script>alert(1)</script>", URL: "https://example.test"}
	out, e := b.JSONLD()
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out, "</script>") {
		t.Fatal("unsafe script body")
	}
	var decoded map[string]any
	if e = json.Unmarshal([]byte(out), &decoded); e != nil || decoded["name"] != b.Name {
		t.Fatal("JSON roundtrip", e)
	}
}
func TestConcurrentSitemapDedup(t *testing.T) {
	s := NewSitemap()
	var wg sync.WaitGroup
	for i := 0; i < 48; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if e := s.Add("t", fmt.Sprintf("https://example.test/%d", i%8)); e != nil {
				t.Error(e)
			}
			s.URLs("t")
			_ = s.String()
		}(i)
	}
	wg.Wait()
	if len(s.URLs("t")) != 8 {
		t.Fatal("dedup mismatch")
	}
}
func FuzzTenantURLsAndRobotLines(f *testing.F) {
	for _, tenant := range []string{"t", "t\x00x", " ", "ñ"} {
		for _, path := range []string{"/ok", "/x\nDisallow: /", "/x\r", "/x\x00"} {
			f.Add(tenant, path, []byte{0, 1, 2, 3, 0, 2})
		}
	}
	f.Fuzz(func(t *testing.T, tenant, path string, ops []byte) {
		if len(tenant) > 64 || len(path) > 128 {
			return
		}
		if len(ops) > 64 {
			ops = ops[:64]
		}
		s := NewSitemap()
		model := map[string][]string{}
		tenants := []string{tenant, tenant + "\x00x", "other"}
		urls := []string{"https://example.test/a", "http://example.test/b", "https://example.test/a?q=1"}
		for _, op := range ops {
			who := tenants[int(op)%3]
			u := urls[int(op/3)%3]
			e := s.Add(who, u)
			if strings.TrimSpace(who) == "" {
				if !errors.Is(e, ErrInvalidURL) {
					t.Fatal(e)
				}
				continue
			}
			if e != nil {
				t.Fatal(e)
			}
			found := false
			for _, v := range model[who] {
				found = found || v == u
			}
			if !found {
				model[who] = append(model[who], u)
			}
			for _, q := range append(tenants, tenant+"x") {
				if !reflect.DeepEqual(s.URLs(q), model[q]) {
					t.Fatal("tenant/order model", q, s.URLs(q), model[q])
				}
			}
		}
		r := NewRobots()
		before := r.Render(tenant)
		valid := strings.TrimSpace(tenant) != "" && path != ""
		for _, c := range path {
			if c == '\r' || c == '\n' || c == 0 {
				valid = false
			}
		}
		e := r.Allow(tenant, path)
		d := r.Disallow(tenant, path)
		if !valid {
			if !errors.Is(e, ErrInvalidURL) || !errors.Is(d, ErrInvalidURL) || r.Render(tenant) != before {
				t.Fatal("injected rule mutated registry")
			}
		} else if e != nil || d != nil || r.Render(tenant) != "User-agent: *\nAllow: "+path+"\nDisallow: "+path+"\n" {
			t.Fatal("rule roundtrip")
		}
	})
}

// V374: the local title length constraint must never be attributed to Google.
func TestSourceMetaPolicyIsLocalV374(t *testing.T) {
	if e := ValidateMeta(strings.Repeat("x", 61), "description"); !errors.Is(e, ErrInvalidMeta) {
		t.Fatal(e)
	}
	if e := ValidateMeta(strings.Repeat("x", 60), "description"); e != nil {
		t.Fatal(e)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | SEO | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/seo/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/seo/`.

## 9. Verification

- `go test ./internal/seo/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (37 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_SEO_CORE_2026-09-02_V215.md`.

## Revisión V363

0.1.1 elimina matching por prefijo de tenant; usa tupla tenant/URL y orden
de inserción. URLs devuelve copia; no verifica propiedad/autorización del dominio.
Sitemap sólo lista URLs, no emite XML. URL/JSONLD exigen HTTP(S) y hostname no vacío;
no DNS, fetch, allowlist de red, deduplicación semántica ni certificación de schema.
Robots rechaza CR/LF/NUL sin mutar; no es parser RFC completo ni access control.
60/160 siguen como política local sobre code points tras TrimSpace, no graphemes
ni garantía Google. No se cambia esa política de producto por atribución incorrecta.
Fuentes verificadas2026-09-09: Google no fija esas longitudes máximas; la visualización
puede truncarse según ancho del dispositivo:
https://developers.google.com/search/docs/appearance/title-link
https://developers.google.com/search/docs/appearance/snippet
Historial previo conservado; evidencia reconstruction_evidence/PUBLIC_CORE_SCOPE_V363.md.

V363: entre los tres cores6fuentes reconstruidas,32tests x3,15tests originales
conservados (social migra argumentos de tenant),vet/build,62semillas y
4518771ejecuciones fuzz PASS. Tres firmas sociales legacy rechazadas.
Sin -race, publicación externa ni admisión del proyecto.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
