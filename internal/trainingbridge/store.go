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
