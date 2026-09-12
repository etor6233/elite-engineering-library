# Go Marketing Core

## 1. Metadata

```yaml
pack_id: "GO-MARKETING-CORE"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Registro AUTHORED de campañas en memoria, consultas scoped por tenant y acuses locales por destinatario; no implementa envío al proveedor."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; dispatch externo no admitido"]
incompatible_with: ["callers sin tenant explícito", "acuse local duplicado", "campaña sin mensaje o sin destinatarios"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Registro AUTHORED en memoria; no envía SMS/email/push. Due retorna snapshots, Send sólo registra un acuse local y no impone ScheduledAt. El caller debe autorizar tenant, consentimiento/supresión y dispatch; no hay reserva, receipt de proveedor, entrega exactamente una vez ni persistencia.

## 3. Architecture contract

- **Ownership**: `internal/marketing` gobierna campañas y acuses locales; no contiene un adapter de dispatch admitido.
- **Invariantes**: (1) canal admitido, mensaje y destinatarios no vacíos. (2) destinatarios deduplicados. (3) acuse local único por tenant/campaña/destinatario; repetir retorna ErrAlreadySent. (4) completado automático cuando todos enviados. (5) tenant-scoped.
- **Data flow**: `Create` → `Due(tenant, now)` → `Send(tenant, id, recipient)` → marca local de completado.
- **Failure modes**: inválido, duplicado, destinatario desconocido, ya enviado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: Send/Remaining O(R); Due O(C + destinatarios copiados), con C campañas y R destinatarios de una campaña.

## 4. Exact file manifest

```text
CREATE internal/marketing/marketing.go
CREATE internal/marketing/marketing_test.go
```

## 5. Materialization blocks

### FILE: `internal/marketing/marketing.go`
```yaml
block_id: "GO-MARKETING-CORE:internal/marketing/marketing.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1a1af2821d5b9bce83b4454842d8af0ad837b5dbb3f23b1705fe802c8a1e811f"
variables: []
secrets_allowed: false
```
````go
// Package marketing provides tenant-scoped marketing campaigns: a channel, a
// message and a deduplicated recipient list, with tenant-scoped local
// acknowledgement tracking. It does not dispatch messages to a provider.
package marketing

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Channel is the dispatch surface.
type Channel string

const (
	ChannelSMS   Channel = "sms"
	ChannelEmail Channel = "email"
	ChannelPush  Channel = "push"
)

var (
	ErrInvalidCampaign = errors.New("marketing: invalid campaign")
	ErrDuplicate       = errors.New("marketing: duplicate campaign")
	ErrNotFound        = errors.New("marketing: not found")
	ErrAlreadySent     = errors.New("marketing: recipient already sent")
)

// Campaign is a scheduled message to a recipient list.
type Campaign struct {
	TenantID    string
	ID          string
	Name        string
	Channel     Channel
	Message     string
	Recipients  []string
	ScheduledAt time.Time
	Sent        bool
}

// Store holds campaigns and per-recipient send tracking.
type Store struct {
	mu        sync.Mutex
	campaigns map[campaignKey]Campaign
	sentTo    map[deliveryKey]bool // (tenant, campaign, recipient)
}

type campaignKey struct{ tenant, id string }
type deliveryKey struct{ tenant, campaign, recipient string }

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{campaigns: make(map[campaignKey]Campaign), sentTo: make(map[deliveryKey]bool)}
}

// Create registers a campaign, deduplicating recipients and validating.
func (s *Store) Create(c Campaign) error {
	if strings.TrimSpace(c.TenantID) == "" || strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Name) == "" {
		return ErrInvalidCampaign
	}
	switch c.Channel {
	case ChannelSMS, ChannelEmail, ChannelPush:
	default:
		return ErrInvalidCampaign
	}
	if strings.TrimSpace(c.Message) == "" || len(c.Recipients) == 0 {
		return ErrInvalidCampaign
	}
	if c.ScheduledAt.IsZero() {
		return ErrInvalidCampaign
	}
	seen := make(map[string]bool)
	dedup := make([]string, 0, len(c.Recipients))
	for _, r := range c.Recipients {
		if strings.TrimSpace(r) == "" {
			return ErrInvalidCampaign
		}
		if !seen[r] {
			seen[r] = true
			dedup = append(dedup, r)
		}
	}
	c.Recipients = dedup

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.campaigns[campaignKey{tenant: c.TenantID, id: c.ID}]; ok {
		return ErrDuplicate
	}
	s.campaigns[campaignKey{tenant: c.TenantID, id: c.ID}] = c
	return nil
}

// Due returns tenant-scoped snapshots whose schedule arrived and Sent is false.
// The result is not a reservation and can repeat before local acknowledgement.
func (s *Store) Due(tenant string, now time.Time) []Campaign {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Campaign
	for _, c := range s.campaigns {
		if c.TenantID == tenant && !c.Sent && !c.ScheduledAt.After(now) {
			c.Recipients = append([]string(nil), c.Recipients...)
			out = append(out, c)
		}
	}
	return out
}

// Send records a local acknowledgement; duplicates return ErrAlreadySent.
// It does not dispatch, reserve delivery or enforce the schedule.
func (s *Store) Send(tenant, campaignID, recipient string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.campaigns[campaignKey{tenant: tenant, id: campaignID}]
	if !ok {
		return ErrNotFound
	}
	k := deliveryKey{tenant: tenant, campaign: campaignID, recipient: recipient}
	if s.sentTo[k] {
		return ErrAlreadySent
	}
	found := false
	for _, r := range c.Recipients {
		if r == recipient {
			found = true
			break
		}
	}
	if !found {
		return ErrNotFound
	}
	s.sentTo[k] = true
	if s.allSent(c) {
		c.Sent = true
		s.campaigns[campaignKey{tenant: tenant, id: campaignID}] = c
	}
	return nil
}

// Remaining returns the count of unsent recipients.
func (s *Store) Remaining(tenant, campaignID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.campaigns[campaignKey{tenant: tenant, id: campaignID}]
	if !ok {
		return 0
	}
	return len(c.Recipients) - s.sentCount(c)
}

func (s *Store) sentCount(c Campaign) int {
	n := 0
	for _, r := range c.Recipients {
		if s.sentTo[deliveryKey{tenant: c.TenantID, campaign: c.ID, recipient: r}] {
			n++
		}
	}
	return n
}

func (s *Store) allSent(c Campaign) bool {
	return s.sentCount(c) == len(c.Recipients)
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("marketing(%d)", len(s.campaigns))
}
````

### FILE: `internal/marketing/marketing_test.go`
```yaml
block_id: "GO-MARKETING-CORE:internal/marketing/marketing_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2e90bd08651a7eed724bb824b744bb0c05e5139db294fb397a99d1e7fab803af"
variables: []
secrets_allowed: false
```
````go
package marketing

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func camp() Campaign {
	return Campaign{
		TenantID: "t", ID: "C1", Name: "Promo", Channel: ChannelEmail,
		Message: "Oferta", Recipients: []string{"a", "b"}, ScheduledAt: time.Now().Add(-time.Hour),
	}
}

func TestCreateDeduplicatesRecipients(t *testing.T) {
	s := NewStore()
	c := camp()
	c.Recipients = []string{"a", "b", "a"}
	if err := s.Create(c); err != nil {
		t.Fatal(err)
	}
	if len(s.campaigns[campaignKey{tenant: "t", id: "C1"}].Recipients) != 2 {
		t.Fatalf("expected 2 unique recipients, got %d", len(s.campaigns[campaignKey{tenant: "t", id: "C1"}].Recipients))
	}
}

func TestSendAndCompletion(t *testing.T) {
	s := NewStore()
	_ = s.Create(camp())
	due := s.Due("t", time.Now())
	if len(due) != 1 {
		t.Fatalf("expected 1 due, got %d", len(due))
	}
	_ = s.Send("t", "C1", "a")
	if s.Remaining("t", "C1") != 1 {
		t.Fatalf("expected 1 remaining, got %d", s.Remaining("t", "C1"))
	}
	_ = s.Send("t", "C1", "b")
	if s.Remaining("t", "C1") != 0 {
		t.Fatalf("expected 0 remaining, got %d", s.Remaining("t", "C1"))
	}
	if len(s.Due("t", time.Now())) != 0 {
		t.Fatal("campaign should not be due after completion")
	}
}

func TestSendDedupPerRecipient(t *testing.T) {
	s := NewStore()
	_ = s.Create(camp())
	_ = s.Send("t", "C1", "a")
	if err := s.Send("t", "C1", "a"); !errors.Is(err, ErrAlreadySent) {
		t.Fatalf("expected ErrAlreadySent, got %v", err)
	}
}

func TestInvalidCampaign(t *testing.T) {
	s := NewStore()
	bad := camp()
	bad.Channel = "carrier"
	if err := s.Create(bad); !errors.Is(err, ErrInvalidCampaign) {
		t.Fatalf("bad channel accepted: %v", err)
	}
	empty := camp()
	empty.Message = ""
	if err := s.Create(empty); !errors.Is(err, ErrInvalidCampaign) {
		t.Fatalf("empty message accepted: %v", err)
	}
	noRecip := camp()
	noRecip.Recipients = nil
	if err := s.Create(noRecip); !errors.Is(err, ErrInvalidCampaign) {
		t.Fatalf("no recipients accepted: %v", err)
	}
}

func TestDuplicateCampaignRejected(t *testing.T) {
	s := NewStore()
	_ = s.Create(camp())
	if err := s.Create(camp()); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate campaign accepted: %v", err)
	}
}

func TestSendUnknownRecipient(t *testing.T) {
	s := NewStore()
	_ = s.Create(camp())
	if err := s.Send("t", "C1", "zzz"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown recipient accepted: %v", err)
	}
}

func TestCampaignIDsIndependentPerTenant(t *testing.T) {
	s := NewStore()
	a, b := camp(), camp()
	b.TenantID = "other"
	if e := s.Create(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Create(b); e != nil {
		t.Fatalf("another tenant blocked: %v", e)
	}
}
func TestDueCanRestrictCallerTenant(t *testing.T) {
	s := NewStore()
	a, b := camp(), camp()
	b.TenantID = "other"
	b.ID = "OTHER"
	for _, c := range []Campaign{a, b} {
		if e := s.Create(c); e != nil {
			t.Fatal(e)
		}
	}
	// The legacy API has no tenant argument with which to enforce this boundary.
	due := s.Due("t", time.Now())
	if len(due) != 1 || due[0].TenantID != "t" {
		t.Fatalf("unscoped due returned %d campaigns", len(due))
	}
}
func TestDueRecipientsDoNotAliasStore(t *testing.T) {
	s := NewStore()
	if e := s.Create(camp()); e != nil {
		t.Fatal(e)
	}
	out := s.Due("t", time.Now())
	out[0].Recipients[0] = "injected"
	if e := s.Send("t", "C1", "injected"); !errors.Is(e, ErrNotFound) {
		t.Fatalf("snapshot changed real recipients: %v", e)
	}
	if e := s.Send("t", "C1", "a"); e != nil {
		t.Fatalf("original recipient removed: %v", e)
	}
}
func TestDeliveryIdentityBoundary(t *testing.T) {
	s := NewStore()
	a, b := camp(), camp()
	a.ID = "a\x00b"
	a.Recipients = []string{"c"}
	b.ID = "a"
	b.Recipients = []string{"b\x00c"}
	for _, c := range []Campaign{a, b} {
		if e := s.Create(c); e != nil {
			t.Fatal(e)
		}
	}
	if e := s.Send("t", "a\x00b", "c"); e != nil {
		t.Fatal(e)
	}
	if n := s.Remaining("t", "a"); n != 1 {
		t.Errorf("other campaign falsely acknowledged: %d", n)
	}
	if e := s.Send("t", "a", "b\x00c"); e != nil {
		t.Errorf("distinct acknowledgement rejected: %v", e)
	}
}

func TestScopedCampaignStateAndScheduleSemantics(t *testing.T) {
	s := NewStore()
	at := time.Unix(1000000, 0)
	c := camp()
	c.ScheduledAt = at
	c.Recipients = []string{"a", "b"}
	if e := s.Create(c); e != nil {
		t.Fatal(e)
	}
	c.Recipients[0] = "changed"
	if len(s.Due("t", at.Add(-time.Nanosecond))) != 0 || len(s.Due("t", at)) != 1 {
		t.Fatal("schedule boundary changed")
	}
	if len(s.Due("other", at)) != 0 || s.Remaining("other", "C1") != 0 {
		t.Fatal("foreign read")
	}
	if e := s.Send("other", "C1", "a"); !errors.Is(e, ErrNotFound) {
		t.Fatal("foreign ack", e)
	}
	if e := s.Send("t", "C1", "changed"); !errors.Is(e, ErrNotFound) {
		t.Fatal("input slice alias", e)
	}
	if e := s.Send("t", "C1", "a"); e != nil {
		t.Fatal(e)
	}
	if e := s.Send("t", "C1", "a"); !errors.Is(e, ErrAlreadySent) {
		t.Fatal("duplicate ack changed", e)
	}
	if s.Remaining("t", "C1") != 1 {
		t.Fatal("duplicate count")
	}
	// Preserve caller-provided Sent and the absence of a schedule gate in Send.
	c = camp()
	c.ID = "IMPORTED"
	c.Sent = true
	c.ScheduledAt = at.Add(time.Hour)
	if e := s.Create(c); e != nil {
		t.Fatal(e)
	}
	if e := s.Send("t", "IMPORTED", "a"); e != nil {
		t.Fatal("local ack newly restricted", e)
	}
	if s.Remaining("t", "IMPORTED") != 1 {
		t.Fatal("Sent flag incorrectly synthesized acknowledgements")
	}
}
func TestConcurrentCampaignAndAcknowledgementScope(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var creates, acks atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c := camp()
			if i%2 == 0 {
				c.TenantID = "other"
			}
			e := s.Create(c)
			if e == nil {
				creates.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
		}(i)
	}
	wg.Wait()
	if creates.Load() != 2 {
		t.Fatal("campaign identity scope")
	}
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tenant := "t"
			if i%2 == 0 {
				tenant = "other"
			}
			e := s.Send(tenant, "C1", "a")
			if e == nil {
				acks.Add(1)
			} else if !errors.Is(e, ErrAlreadySent) {
				t.Error(e)
			}
			_ = s.String()
			s.Due(tenant, time.Now())
		}(i)
	}
	wg.Wait()
	if acks.Load() != 2 {
		t.Fatal("ack identity scope")
	}
	for _, tenant := range []string{"t", "other"} {
		if s.Remaining(tenant, "C1") != 1 {
			t.Fatal("wrong remaining")
		}
	}
}
func FuzzCampaignAcknowledgementModel(f *testing.F) {
	for _, x := range []string{"a", "a\x00b", "", " ", "客户", "b\x00c"} {
		f.Add(x, []byte{0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0})
		f.Add(x, []byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 1, 0, 0, 0})
	}
	f.Fuzz(func(t *testing.T, x string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(ops) > 125 {
			ops = ops[:125]
		}
		tenants := []string{"t", "other", x, " "}
		ids := []string{"a", "a\x00b", x, " "}
		recipients := []string{"c", "b\x00c", x, " "}
		plans := [][]string{{"c", "next", "c"}, {"b\x00c"}, {x, "c"}, nil, {" "}}
		now := time.Unix(1000000, 0)
		type entry struct {
			c   Campaign
			ack []bool
		}
		var model []entry
		s := NewStore()
		lookup := func(tenant, id string) int {
			for i, e := range model {
				if e.c.TenantID == tenant && e.c.ID == id {
					return i
				}
			}
			return -1
		}
		for i := 0; i+4 < len(ops); i += 5 {
			tenant := tenants[int(ops[i+1])%4]
			id := ids[int(ops[i+2])%4]
			index := lookup(tenant, id)
			var got, want error
			switch ops[i] % 3 {
			case 0:
				c := Campaign{TenantID: tenant, ID: id, Name: "synthetic", Channel: ChannelEmail, Message: "test only", Recipients: plans[int(ops[i+3])%len(plans)], ScheduledAt: now, Sent: ops[i+4]&1 != 0}
				switch ops[i+4] % 8 {
				case 1:
					c.ScheduledAt = now.Add(time.Hour)
				case 2:
					c.ScheduledAt = time.Time{}
				case 3:
					c.Channel = "invalid"
				case 4:
					c.Name = " "
				case 5:
					c.Message = " "
				}
				if strings.TrimSpace(tenant) == "" || strings.TrimSpace(id) == "" || strings.TrimSpace(c.Name) == "" || strings.TrimSpace(c.Message) == "" || c.ScheduledAt.IsZero() || (c.Channel != ChannelSMS && c.Channel != ChannelEmail && c.Channel != ChannelPush) || len(c.Recipients) == 0 {
					want = ErrInvalidCampaign
				}
				var unique []string
				for _, r := range c.Recipients {
					if strings.TrimSpace(r) == "" {
						want = ErrInvalidCampaign
					}
					found := false
					for _, old := range unique {
						if r == old {
							found = true
						}
					}
					if !found {
						unique = append(unique, r)
					}
				}
				if want == nil && index >= 0 {
					want = ErrDuplicate
				}
				got = s.Create(c)
				if want == nil {
					c.Recipients = unique
					model = append(model, entry{c: c, ack: make([]bool, len(unique))})
				}
			case 1:
				recipient := recipients[int(ops[i+3])%4]
				if index < 0 {
					want = ErrNotFound
				} else {
					pos := -1
					for j, r := range model[index].c.Recipients {
						if r == recipient {
							pos = j
						}
					}
					if pos < 0 {
						want = ErrNotFound
					} else if model[index].ack[pos] {
						want = ErrAlreadySent
					} else {
						model[index].ack[pos] = true
						all := true
						for _, done := range model[index].ack {
							all = all && done
						}
						if all {
							model[index].c.Sent = true
						}
					}
				}
				got = s.Send(tenant, id, recipient)
			case 2:
				_ = s.String()
			}
			if !errors.Is(got, want) {
				t.Fatalf("step %d got %v want %v", i/5, got, want)
			}
			for _, tenant := range tenants {
				for _, id := range ids {
					index := lookup(tenant, id)
					remaining := 0
					if index >= 0 {
						for _, done := range model[index].ack {
							if !done {
								remaining++
							}
						}
					}
					if s.Remaining(tenant, id) != remaining {
						t.Fatal("wrong scoped remaining")
					}
				}
				expected := []Campaign{}
				for _, e := range model {
					if e.c.TenantID == tenant && !e.c.Sent && !e.c.ScheduledAt.After(now) {
						expected = append(expected, e.c)
					}
				}
				due := s.Due(tenant, now)
				if len(due) != len(expected) {
					t.Fatal("wrong due count")
				}
				seen := make([]bool, len(expected))
				for _, c := range due {
					found := -1
					for j, e := range expected {
						if c.ID == e.ID && c.TenantID == e.TenantID && c.Name == e.Name && c.Message == e.Message && c.Channel == e.Channel && c.Sent == e.Sent && c.ScheduledAt.Equal(e.ScheduledAt) {
							found = j
							break
						}
					}
					if found < 0 || seen[found] {
						t.Fatal("foreign/duplicate due")
					}
					seen[found] = true
					e := expected[found]
					if len(c.Recipients) != len(e.Recipients) {
						t.Fatal("recipient count")
					}
					for j, r := range c.Recipients {
						if r != e.Recipients[j] {
							t.Fatal("recipient order/ownership")
						}
					}
					if len(c.Recipients) > 0 {
						c.Recipients[0] = "snapshot mutation"
					}
				}
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
| Go standard library | go 1.26.7 (toolchain) | campañas | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/marketing/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/marketing/`.

## 9. Verification

- `go test ./internal/marketing/ -count=1`: 6/6 PASS (dedup destinatarios, envío y completado, dedup por destinatario, inválido, duplicado, destinatario desconocido).
- `go test ./... -count=1` (23 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_MARKETING_CORE_2026-09-02_V201.md`.

## Migración obligatoria 0.2.0

Due(now) → Due(tenant, now); Send(id, recipient) → Send(tenant, id, recipient);
Remaining(id) → Remaining(tenant, id). El tenant debe venir del contexto autorizado,
no de texto no confiable. No shim global; callers antiguos fallan compilación.
Due copia Recipients y no reserva ni ordena envíos. Remaining desconocido retorna0,
sin distinguirlo de completado; Create conserva Sent dado por caller. Send puede
registrarse antes de ScheduledAt porque es acuse local, no dispatch.

## Revisión V360

0.2.0: identidad estructurada y String sincronizado; marketing además exige tenant
y copia snapshots. AUTHORED/LicenseRef-Workspace-Owner y CONDITIONED persisten.
Suite V201 histórica no reejecutada; compatibilidad antigua no demuestra
integración actual. Evidencia: reconstruction_evidence/SURVEY_CAMPAIGN_SCOPE_V360.md.

V360 verificado entre ambos cores: 4/4 fuentes reconstruidas, 20 tests x3,
vet/build, 24 semillas/3864824 ejecuciones fuzz PASS.11tests históricos conservan
aserciones; marketing migra sólo argumentos tenant e índice privado del test.
Tres firmas sin tenant fallan compilación. Sin -race, dispatch ni consentimiento
real probado.

## Auditoría por claim V374

La revisión 0.2.1 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.2.0 y sus bytes quedan preservados en el expediente anterior.
