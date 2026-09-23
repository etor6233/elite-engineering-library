package approval

import (
	"sync"
	"time"
)

type record struct {
	req       Request
	state     State
	approvers map[string]bool
}

// Registry is the approval state machine. Money effects (payment/refund) never
// auto-approve; high-value requests require dual control; velocity is bounded.
type Registry struct {
	mu        sync.Mutex
	policy    Policy
	records   map[string]*record
	openCount map[string]int
	decisions []Decision
	clock     func() time.Time
}

// NewRegistry returns a registry with the given policy.
func NewRegistry(p Policy) *Registry {
	return &Registry{
		policy:    p,
		records:   make(map[string]*record),
		openCount: make(map[string]int),
		clock:     time.Now,
	}
}

func (r *Registry) now() time.Time {
	if r.clock != nil {
		return r.clock()
	}
	return time.Now().UTC()
}

func key(tenant, id string) string          { return tenant + "\x00" + id }
func subjKey(tenant, subject string) string { return tenant + "\x00" + subject }

// Submit validates and registers a request. It auto-approves only low-value,
// non-money effects; everything else stays pending. Velocity and duplicates
// fail closed.
func (r *Registry) Submit(req Request) (State, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.records[key(req.TenantID, req.ID)]; ok {
		return "", ErrDuplicate
	}
	sk := subjKey(req.TenantID, req.SubjectID)
	if r.policy.MaxOpenPerSubject > 0 && r.openCount[sk] >= r.policy.MaxOpenPerSubject {
		return "", ErrSubjectFlagged
	}
	rec := &record{req: req, state: StatePending, approvers: map[string]bool{}}
	r.records[key(req.TenantID, req.ID)] = rec
	r.openCount[sk]++
	if r.autoApprove(req) {
		rec.state = StateApproved
		r.openCount[sk]--
		r.decisions = append(r.decisions, Decision{
			RequestID: req.ID, Reviewer: "system", Approved: true, Reason: "auto", At: r.now(),
		})
	}
	return rec.state, nil
}

func (r *Registry) autoApprove(req Request) bool {
	if r.policy.AutoApproveMinorUnits <= 0 {
		return false
	}
	if req.Kind == KindPayment || req.Kind == KindRefund || req.Kind == KindWhatsAppReply || req.Kind == KindSocialPublish || req.Kind == KindSocialRevoke || req.Kind == KindStoredValueOperation || req.Kind == KindWarrantyRepair || req.Kind == KindSerialQuality || req.Kind == KindCatalogReview || req.Kind == KindTrainingAssessment || req.Kind == KindMarketplaceMutation || req.Kind == KindWhatsAppSchedule || req.Kind == KindDocumentReview {
		return false // money movements always require a human
	}
	return req.AmountMinorUnits <= r.policy.AutoApproveMinorUnits
}

// Approve records a distinct reviewer. Dual-control requests stay pending until
// two distinct reviewers have approved.
func (r *Registry) Approve(tenant, id, reviewer, reason string) (State, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.records[key(tenant, id)]
	if !ok {
		return "", ErrNotFound
	}
	if rec.state != StatePending {
		return "", ErrNotPending
	}
	if reviewer == "" || reviewer == rec.req.Requester {
		return "", ErrSeparation
	}
	if rec.approvers[reviewer] {
		return "", ErrDuplicateApprover
	}
	rec.approvers[reviewer] = true
	needsDual := r.policy.DualControlMinorUnits > 0 && rec.req.AmountMinorUnits >= r.policy.DualControlMinorUnits
	if needsDual && len(rec.approvers) < 2 {
		return StatePending, nil // awaiting a second reviewer
	}
	rec.state = StateApproved
	r.openCount[subjKey(tenant, rec.req.SubjectID)]--
	r.decisions = append(r.decisions, Decision{
		RequestID: id, Reviewer: reviewer, Approved: true, Reason: reason, At: r.now(),
	})
	return StateApproved, nil
}

// Reject records a rejection (separation of duties enforced).
func (r *Registry) Reject(tenant, id, reviewer, reason string) (State, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.records[key(tenant, id)]
	if !ok {
		return "", ErrNotFound
	}
	if rec.state != StatePending {
		return "", ErrNotPending
	}
	if reviewer == "" || reviewer == rec.req.Requester {
		return "", ErrSeparation
	}
	rec.state = StateRejected
	r.openCount[subjKey(tenant, rec.req.SubjectID)]--
	r.decisions = append(r.decisions, Decision{
		RequestID: id, Reviewer: reviewer, Approved: false, Reason: reason, At: r.now(),
	})
	return StateRejected, nil
}

// State returns the current state of a request.
func (r *Registry) State(tenant, id string) (State, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.records[key(tenant, id)]
	if !ok {
		return "", false
	}
	return rec.state, true
}

// Audit returns a copy of the decision log.
func (r *Registry) Audit() []Decision {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Decision, len(r.decisions))
	copy(out, r.decisions)
	return out
}
