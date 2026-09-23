// Package aifoundation defines a provider-agnostic contract for serving
// generative models: exact model pinning, a single-active version registry,
// closed-schema structured output, fail-closed safety gates and eval-gated
// promotion. It is AUTHORED glue over the authoritative AI corpus; it does
// not vendor any model, SDK or provider behavior.
package aifoundation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ModelRef is an exact, immutable identity for a model revision.
type ModelRef struct {
	Provider string
	Model    string
	Version  string
	Digest   string // hex sha256 of the model/policy artifact pin
}

var (
	identRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	hex64Re = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

// ErrInvalidModelRef reports an identity that violates the exact-pin contract.
var ErrInvalidModelRef = errors.New("aifoundation: invalid model ref")

// Validate enforces exact pins: no empty fields, no mutable "latest"/"auto"
// versions and a 64-hex digest. Mutable versions are rejected fail-closed.
func (m ModelRef) Validate() error {
	if !identRe.MatchString(m.Provider) {
		return fmt.Errorf("%w: bad provider", ErrInvalidModelRef)
	}
	if !identRe.MatchString(m.Model) {
		return fmt.Errorf("%w: bad model", ErrInvalidModelRef)
	}
	if !identRe.MatchString(m.Version) {
		return fmt.Errorf("%w: bad version", ErrInvalidModelRef)
	}
	if strings.EqualFold(m.Version, "latest") || strings.EqualFold(m.Version, "auto") {
		return fmt.Errorf("%w: mutable version %q forbidden", ErrInvalidModelRef, m.Version)
	}
	if !hex64Re.MatchString(strings.ToLower(m.Digest)) {
		return fmt.Errorf("%w: digest must be 64 hex chars", ErrInvalidModelRef)
	}
	return nil
}

// RegistryState is the promotion state of a model version.
type RegistryState string

const (
	StateCandidate  RegistryState = "candidate"
	StateShadow     RegistryState = "shadow"
	StateActive     RegistryState = "active"
	StateRolledBack RegistryState = "rolledback"
)

// AuditEvent is an immutable registry transition record.
type AuditEvent struct {
	At     time.Time
	Ref    ModelRef
	From   RegistryState
	To     RegistryState
	Reason string
}

type version struct {
	ref   ModelRef
	state RegistryState
}

// VersionRegistry maintains the promotion lifecycle and enforces that at most
// one version is active at any time. The zero value is ready to use.
type VersionRegistry struct {
	mu       sync.Mutex
	versions map[string]*version // key: lowercase digest
	order    []string
	audit    []AuditEvent
	clock    func() time.Time
}

func (r *VersionRegistry) now() time.Time {
	if r.clock != nil {
		return r.clock()
	}
	return time.Now().UTC()
}

func (r *VersionRegistry) ensure() {
	if r.versions == nil {
		r.versions = make(map[string]*version)
	}
}

// Register adds a version as a candidate. Duplicate digests are rejected.
func (r *VersionRegistry) Register(ref ModelRef) error {
	if err := ref.Validate(); err != nil {
		return err
	}
	key := strings.ToLower(ref.Digest)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensure()
	if _, ok := r.versions[key]; ok {
		return fmt.Errorf("%w: duplicate digest", ErrInvalidModelRef)
	}
	r.versions[key] = &version{ref: ref, state: StateCandidate}
	r.order = append(r.order, key)
	r.audit = append(r.audit, AuditEvent{At: r.now(), Ref: ref, From: "", To: StateCandidate, Reason: "register"})
	return nil
}

// Promote activates a candidate only when the release gate passed. Any prior
// active version is rolled back and recorded. Exactly one version is active.
func (r *VersionRegistry) Promote(ref ModelRef, gatePassed bool) error {
	if err := ref.Validate(); err != nil {
		return err
	}
	key := strings.ToLower(ref.Digest)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensure()
	v, ok := r.versions[key]
	if !ok {
		return errors.New("aifoundation: promote unknown version")
	}
	if v.state != StateCandidate && v.state != StateShadow {
		return fmt.Errorf("aifoundation: cannot promote from %s", v.state)
	}
	if !gatePassed {
		return errors.New("aifoundation: promote rejected: release gate not passed")
	}
	from := v.state
	// Demote any active version first.
	for _, k := range r.order {
		if k != key && r.versions[k].state == StateActive {
			r.versions[k].state = StateRolledBack
			r.audit = append(r.audit, AuditEvent{At: r.now(), Ref: r.versions[k].ref, From: StateActive, To: StateRolledBack, Reason: "superseded"})
		}
	}
	v.state = StateActive
	r.audit = append(r.audit, AuditEvent{At: r.now(), Ref: ref, From: from, To: StateActive, Reason: "promote"})
	return nil
}

// Rollback deactivates the active version. It does not auto-reactivate a prior
// version; that decision is a new, explicit Promote.
func (r *VersionRegistry) Rollback(ref ModelRef, reason string) error {
	key := strings.ToLower(ref.Digest)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensure()
	v, ok := r.versions[key]
	if !ok {
		return errors.New("aifoundation: rollback unknown version")
	}
	if v.state != StateActive {
		return fmt.Errorf("aifoundation: cannot rollback from %s", v.state)
	}
	v.state = StateRolledBack
	r.audit = append(r.audit, AuditEvent{At: r.now(), Ref: ref, From: StateActive, To: StateRolledBack, Reason: reason})
	return nil
}

// Active returns the single active version, if any.
func (r *VersionRegistry) Active() (ModelRef, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ensure()
	for _, k := range r.order {
		if r.versions[k].state == StateActive {
			return r.versions[k].ref, true
		}
	}
	return ModelRef{}, false
}

// Audit returns a copy of the transition log.
func (r *VersionRegistry) Audit() []AuditEvent {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]AuditEvent, len(r.audit))
	copy(out, r.audit)
	return out
}
