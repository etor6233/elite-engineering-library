# Go i18n Core

## 1. Metadata

```yaml
pack_id: "GO-I18N-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Catálogo AUTHORED con fallback exacto→base→default→key y regla binaria n==1 one/else other; rechaza formas vacías. No implementa CLDR, validación completa de tags ni interpolación."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; composición y target deben demostrar condiciones"]
incompatible_with: ["key vacía/inválida", "locale inválido", "texto vacío"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://cldr.unicode.org"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia local de catálogo. La regla binaria debe ser compatible con los idiomas/números del proyecto; no representa un motor multilingüe CLDR.

## 3. Architecture contract

- **Ownership**: `internal/i18n` gobierna el catálogo; el locale del usuario lo resuelve la sesión/preferencia.
- **Invariantes**: (1) key segura. (2) etiqueta de locale validada por gramática local, sin probar existencia ISO/BCP47. (3) fallback exacto→base→default→key. (4) plural one/other.
- **Data flow**: `Set`/`SetPlural` → `T`/`Tn`.
- **Failure modes**: key/locale inválido → error; missing → key (no inventa texto).
- **Seguridad/privacidad**: sin datos.
- **Performance budget**: O(1) por lookup.

## 4. Exact file manifest

```text
CREATE internal/i18n/i18n.go
CREATE internal/i18n/i18n_test.go
```

## 5. Materialization blocks

### FILE: `internal/i18n/i18n.go`
```yaml
block_id: "GO-I18N-CORE:internal/i18n/i18n.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "818df4f6310c82a45520c2fb858745702c590a9c4d37db14decdd3ecc44b1bac"
variables: []
secrets_allowed: false
```
````go
// Package i18n provides a locale catalog: message keys, per-locale strings,
// a fallback chain and a local binary rule (n == 1: one; otherwise: other).
// This is not a CLDR engine or a complete language-tag validator.
package i18n

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Locale is a syntactically bounded catalog label (e.g. "es", "es-AR", "en").
type Locale string

var (
	keyRe = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,127}$`)

	ErrInvalidKey    = errors.New("i18n: invalid key")
	ErrInvalidLocale = errors.New("i18n: invalid locale")
)

// PluralCategory is the local binary plural set; not CLDR equivalence.
type PluralCategory string

const (
	PluralOne   PluralCategory = "one"
	PluralOther PluralCategory = "other"
)

// Catalog holds translations.
type Catalog struct {
	mu            sync.RWMutex
	messages      map[string]map[Locale]string
	plurals       map[string]map[Locale]map[PluralCategory]string
	defaultLocale Locale
}

// NewCatalog returns an empty catalog with the given default locale.
func NewCatalog(defaultLocale Locale) *Catalog {
	return &Catalog{
		messages:      make(map[string]map[Locale]string),
		plurals:       make(map[string]map[Locale]map[PluralCategory]string),
		defaultLocale: defaultLocale,
	}
}

func validLocale(l Locale) bool {
	s := string(l)
	return len(s) >= 2 && len(s) <= 16 && regexp.MustCompile(`^[a-z]{2}(-[A-Za-z0-9]{2,8})*$`).MatchString(s)
}

// Set registers a translation for a key and locale.
func (c *Catalog) Set(key string, locale Locale, text string) error {
	if !keyRe.MatchString(key) {
		return ErrInvalidKey
	}
	if !validLocale(locale) || strings.TrimSpace(text) == "" {
		return ErrInvalidLocale
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.messages[key] == nil {
		c.messages[key] = make(map[Locale]string)
	}
	c.messages[key][locale] = text
	return nil
}

// SetPlural registers one/other plural forms.
func (c *Catalog) SetPlural(key string, locale Locale, one, other string) error {
	if !keyRe.MatchString(key) || !validLocale(locale) {
		return ErrInvalidKey
	}
	if strings.TrimSpace(one) == "" || strings.TrimSpace(other) == "" {
		return ErrInvalidLocale
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.plurals[key] == nil {
		c.plurals[key] = make(map[Locale]map[PluralCategory]string)
	}
	if c.plurals[key][locale] == nil {
		c.plurals[key][locale] = make(map[PluralCategory]string)
	}
	c.plurals[key][locale][PluralOne] = one
	c.plurals[key][locale][PluralOther] = other
	return nil
}

// T resolves a message with fallback: exact locale → base → default → key.
func (c *Catalog) T(key string, locale Locale) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if m, ok := c.messages[key]; ok {
		if s, ok := m[locale]; ok {
			return s
		}
		base := baseLocale(locale)
		if s, ok := m[base]; ok {
			return s
		}
		if s, ok := m[c.defaultLocale]; ok {
			return s
		}
	}
	return key
}

// Tn resolves a plural form (n == 1 → one, else other) with the same fallback.
func (c *Catalog) Tn(key string, locale Locale, n int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cat := PluralOther
	if n == 1 {
		cat = PluralOne
	}
	if m, ok := c.plurals[key]; ok {
		for _, loc := range []Locale{locale, baseLocale(locale), c.defaultLocale} {
			if forms, ok := m[loc]; ok {
				if s, ok := forms[cat]; ok {
					return s
				}
			}
		}
	}
	return key
}

func baseLocale(l Locale) Locale {
	s := string(l)
	if i := strings.IndexByte(s, '-'); i > 0 {
		return Locale(s[:i])
	}
	return l
}

// String aids debugging.
func (c *Catalog) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return fmt.Sprintf("i18n(%d)", len(c.messages))
}
````

### FILE: `internal/i18n/i18n_test.go`
```yaml
block_id: "GO-I18N-CORE:internal/i18n/i18n_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "53fe04950e4f95d65315983964a90375c83160005101d57947db1924a06e1092"
variables: []
secrets_allowed: false
```
````go
package i18n

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestTWithFallback(t *testing.T) {
	c := NewCatalog("es")
	_ = c.Set("hello", "es", "Hola")
	_ = c.Set("hello", "en", "Hello")

	if c.T("hello", "es") != "Hola" {
		t.Fatalf("es lookup failed")
	}
	if c.T("hello", "en") != "Hello" {
		t.Fatalf("en lookup failed")
	}
	// es-AR falls back to base es
	if c.T("hello", "es-AR") != "Hola" {
		t.Fatalf("es-AR should fall back to es")
	}
	// missing locale falls back to default es
	if c.T("hello", "pt") != "Hola" {
		t.Fatalf("pt should fall back to default es")
	}
	// missing key returns the key
	if c.T("missing", "es") != "missing" {
		t.Fatalf("missing key should return key")
	}
}

func TestPlural(t *testing.T) {
	c := NewCatalog("es")
	_ = c.SetPlural("item", "es", "1 artículo", "%d artículos")
	_ = c.SetPlural("item", "en", "1 item", "%d items")

	if c.Tn("item", "es", 1) != "1 artículo" {
		t.Fatalf("one form failed")
	}
	if c.Tn("item", "es", 3) != "%d artículos" {
		t.Fatalf("other form failed")
	}
	if c.Tn("item", "en", 1) != "1 item" {
		t.Fatalf("en one failed")
	}
	if c.Tn("item", "en", 5) != "%d items" {
		t.Fatalf("en other failed")
	}
}

func TestInvalidKeyAndLocale(t *testing.T) {
	c := NewCatalog("es")
	if err := c.Set("BAD KEY!", "es", "x"); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("bad key accepted: %v", err)
	}
	if err := c.Set("k", "123", "x"); !errors.Is(err, ErrInvalidLocale) {
		t.Fatalf("bad locale accepted: %v", err)
	}
}

func TestBaseLocale(t *testing.T) {
	if baseLocale("es-AR") != "es" {
		t.Fatalf("base of es-AR should be es")
	}
	if baseLocale("en") != "en" {
		t.Fatalf("base of en should be en")
	}
}

func TestPluralRejectsBlankWithoutReplacing(t *testing.T) {
	for _, forms := range [][2]string{{"", "many"}, {"one", ""}, {" \t", "many"}, {"one", "\n"}} {
		c := NewCatalog("es")
		if e := c.SetPlural("items", "es", "one", "many"); e != nil {
			t.Fatal(e)
		}
		if e := c.SetPlural("items", "es", forms[0], forms[1]); e == nil {
			t.Error("blank form admitted")
		}
		if c.Tn("items", "es", 1) != "one" || c.Tn("items", "es", 2) != "many" {
			t.Error("rejected plural changed catalog")
		}
	}
}

func TestRejectedPluralDoesNotCreateKey(t *testing.T) {
	c := NewCatalog("es")
	if e := c.SetPlural("item", "es", "", "other"); !errors.Is(e, ErrInvalidLocale) {
		t.Fatal(e)
	}
	if c.Tn("item", "es", 1) != "item" || c.Tn("item", "es", 2) != "item" {
		t.Fatal("rejected plural reserved key")
	}
}
func TestPluralErrorPrecedence(t *testing.T) {
	c := NewCatalog("es")
	if e := c.SetPlural("BAD!", "123", "", ""); !errors.Is(e, ErrInvalidKey) {
		t.Fatal(e)
	}
	if e := c.SetPlural("item", "123", "one", "other"); !errors.Is(e, ErrInvalidKey) {
		t.Fatal("existing locale error changed", e)
	}
}
func TestConcurrentCatalogReadWrite(t *testing.T) {
	c := NewCatalog("en")
	var wg sync.WaitGroup
	for i := 0; i < 48; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("key%d", i)
			if e := c.Set(k, "en", "value"); e != nil {
				t.Error(e)
			}
			if e := c.SetPlural(k, "en", "one", "other"); e != nil {
				t.Error(e)
			}
			if c.T(k, "en") != "value" || c.Tn(k, "en", 2) != "other" {
				t.Error("lookup mismatch")
			}
			_ = c.String()
		}(i)
	}
	wg.Wait()
	if len(c.messages) != 48 || len(c.plurals) != 48 {
		t.Fatal("lost keys")
	}
}
func FuzzBinaryPluralFallback(f *testing.F) {
	for _, loc := range []string{"es-AR", "es", "en", "pt", "bad label"} {
		for _, one := range []string{"one", "", " \t"} {
			f.Add(loc, one, "other", int64(1))
			f.Add(loc, one, "other", int64(-2))
		}
	}
	f.Fuzz(func(t *testing.T, requested, one, other string, n int64) {
		if len(requested) > 64 || len(one) > 128 || len(other) > 128 {
			return
		}
		c := NewCatalog("en")
		if e := c.SetPlural("item", "en", "default-one", "default-other"); e != nil {
			t.Fatal(e)
		}
		if e := c.SetPlural("item", "es", "base-one", "base-other"); e != nil {
			t.Fatal(e)
		}
		e := c.SetPlural("item", "es-AR", one, other)
		valid := strings.TrimSpace(one) != "" && strings.TrimSpace(other) != ""
		if valid && e != nil || !valid && !errors.Is(e, ErrInvalidLocale) {
			t.Fatal("validation mismatch", e)
		}
		singular := n == 1
		want := "default-other"
		if singular {
			want = "default-one"
		}
		base := strings.SplitN(requested, "-", 2)[0]
		if requested == "es" || base == "es" {
			want = "base-other"
			if singular {
				want = "base-one"
			}
		}
		if requested == "es-AR" && valid {
			want = other
			if singular {
				want = one
			}
		}
		if got := c.Tn("item", Locale(requested), int(n)); got != want {
			t.Fatal("fallback model", got, want)
		}
		if c.Tn("missing", Locale(requested), int(n)) != "missing" {
			t.Fatal("missing key")
		}
	})
}

// V374: CLDR French cardinal zero is one; the local binary rule is deliberately
// not equivalent. Keep the refusal of a general CLDR claim executable.
func TestSourceCLDRNonEquivalenceV374(t *testing.T) {
	c := NewCatalog("fr")
	if e := c.SetPlural("items", "fr", "un", "autres"); e != nil {
		t.Fatal(e)
	}
	if c.Tn("items", "fr", 0) != "autres" || c.Tn("items", "fr", 1) != "un" {
		t.Fatal("binary contract changed")
	}
	if c.Tn("items", "fr", 0) == "un" {
		t.Fatal("audit requires explicit CLDR mismatch")
	}
}
````


## 6. Configuration surface

Sin variables ni secretos. El `defaultLocale` se fija en construcción.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | catálogo | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/i18n/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/i18n/`.

## 9. Verification

- `go test ./internal/i18n/ -count=1`: 4/4 PASS (fallback, plural, inválido, base locale).
- `go test ./... -count=1` (26 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_I18N_CORE_2026-09-02_V204.md`.

## Revisión V363

0.1.1 conserva firmas/fallback y n==1, rechaza one/other vacío o whitespace antes
de mutar y sincroniza String. SetPlural conserva ErrInvalidKey para key/locale
inválidos y usa ErrInvalidLocale para forma vacía, como Set para texto vacío.
NewCatalog no valida defaultLocale: el caller debe elegir uno válido; no se
admite fallback/base del default, CLDR, regionalización completa, formato numérico
ni interpolación automática. La regla binaria/local no puede atribuirse a Unicode.
Autoridad contrastada2026-09-09: https://cldr.unicode.org/index/cldr-spec/plural-rules
V204 histórica conservada. Evidencia reconstruction_evidence/PUBLIC_CORE_SCOPE_V363.md.

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
