package postgres

// AUTHORED same-transaction composition of existing job fencing, human approvals
// and immutable original/extraction records. No automatic storage authority.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	runtime "example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var ErrDocumentScope = errors.New("document not found or forbidden")
var ErrDocumentConflict = errors.New("document conflict; consult recorded state")

type Documents struct {
	pool      *pgxpool.Pool
	scope     doc.Scope
	processor doc.Processor
	reviews   *HumanApprovals
}
type DocumentProposal struct {
	Schema         string          `json:"schema"`
	DocumentID     string          `json:"document_id"`
	OrganizationID string          `json:"organization_id"`
	OriginalSHA    string          `json:"original_sha256"`
	ProfileSHA     string          `json:"profile_sha256"`
	EvidenceSHA    string          `json:"evidence_sha256"`
	Mode           string          `json:"mode"`
	Fields         json.RawMessage `json:"fields"`
	ClassID        string          `json:"class_id,omitempty"`
	SchemaVersion  string          `json:"schema_version,omitempty"`
}
type DocumentView struct {
	ID            string            `json:"document_id"`
	Name          string            `json:"name"`
	Uploader      string            `json:"uploader"`
	OriginalSHA   string            `json:"original_sha256"`
	ProfileSHA    string            `json:"profile_sha256"`
	Mode          string            `json:"mode"`
	State         string            `json:"state"`
	EvidenceSHA   string            `json:"evidence_sha256,omitempty"`
	Suggested     json.RawMessage   `json:"suggested,omitempty"`
	Proposal      *DocumentProposal `json:"proposal,omitempty"`
	PayloadSHA    string            `json:"payload_sha256,omitempty"`
	Reviewer      string            `json:"reviewer,omitempty"`
	Reason        string            `json:"reason,omitempty"`
	ClassID       string            `json:"class_id"`
	SchemaVersion string            `json:"schema_version"`
}

func NewDocuments(pool *pgxpool.Pool, s doc.Scope, p doc.Processor) (*Documents, error) {
	if pool == nil || p == nil || !documentID(s.TenantID) || !doc.Text(s.OrganizationID, 128) || !doc.Hex(s.ProfileSHA) || (s.Mode != "FIXTURE" && s.Mode != "PROVIDER" && s.Mode != "TYPED_FIXTURE") {
		return nil, doc.ErrContract
	}
	if s.Mode == "TYPED_FIXTURE" && s.ProfileSHA != doc.TypedProfileSHA() {
		return nil, doc.ErrContract
	}
	return &Documents{pool, s, p, NewHumanApprovals(pool)}, nil
}
func documentID(v string) bool {
	x, e := uuid.Parse(v)
	return e == nil && x != uuid.Nil && x.String() == v
}
func (d *Documents) Allowed(p identity.Principal, permission string) bool {
	return d != nil && p.TenantID == d.scope.TenantID && doc.Text(p.Subject, 128) && p.AllowedOrganization(d.scope.OrganizationID) && p.Allowed(permission)
}
func (d *Documents) canRead(p identity.Principal) bool {
	return d.Allowed(p, "documents:read") || d.Allowed(p, "documents:write") || d.Allowed(p, "documents:review") || d.Allowed(p, "documents:process")
}
func (d *Documents) requestID(id string) string { return "document:" + id }
func (d *Documents) jobPayload(id string) []byte {
	b, _ := json.Marshal(map[string]string{"document_id": id, "organization_id": d.scope.OrganizationID, "profile_sha256": d.scope.ProfileSHA})
	return b
}
func documentJobID(tenant, org, profile, id string) string {
	payload, _ := json.Marshal(map[string]string{"document_id": id, "organization_id": org, "profile_sha256": profile})
	return uuid.NewSHA1(uuid.NameSpaceOID, append([]byte("document-extraction:"+tenant+":"), payload...)).String()
}
func (d *Documents) jobID(id string) string {
	return documentJobID(d.scope.TenantID, d.scope.OrganizationID, d.scope.ProfileSHA, id)
}
func (d *Documents) Receive(ctx context.Context, p identity.Principal, id, profileSHA string, v doc.Original) (DocumentView, error) {
	if !d.Allowed(p, "documents:write") || !documentID(id) {
		return DocumentView{}, ErrDocumentScope
	}
	if profileSHA != d.scope.ProfileSHA || v.Validate(d.scope.Mode) != nil {
		return DocumentView{}, doc.ErrContract
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return DocumentView{}, e
	}
	defer tx.Rollback(ctx)
	if v.ClassID == "" {
		v.ClassID = doc.Class
	}
	if v.SchemaVersion == "" {
		v.SchemaVersion = "1"
	}
	_, e = tx.Exec(ctx, `insert into document.original(tenant_id,document_id,organization_id,uploader,name,original_sha256,profile_sha256,mode,content,class_id,schema_version)values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)on conflict do nothing`, p.TenantID, id, d.scope.OrganizationID, p.Subject, v.Name, v.SHA256, profileSHA, d.scope.Mode, v.Bytes, v.ClassID, v.SchemaVersion)
	if e != nil {
		return DocumentView{}, e
	}
	var same bool
	e = tx.QueryRow(ctx, `select organization_id=$3 and uploader=$4 and name=$5 and original_sha256=$6 and profile_sha256=$7 and mode=$8 and content=$9 and class_id=$10 and schema_version=$11 from document.original where tenant_id=$1 and document_id=$2`, p.TenantID, id, d.scope.OrganizationID, p.Subject, v.Name, v.SHA256, profileSHA, d.scope.Mode, v.Bytes, v.ClassID, v.SchemaVersion).Scan(&same)
	if e != nil || !same {
		return DocumentView{}, ErrDocumentConflict
	}
	_, e = tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts)values($1,$2,'document-extraction','document-extraction',1,$3,3)on conflict do nothing`, p.TenantID, d.jobID(id), d.jobPayload(id))
	if e != nil {
		return DocumentView{}, e
	}
	if e = tx.Commit(ctx); e != nil {
		return DocumentView{}, e
	}
	return d.Read(ctx, p, id)
}
func (d *Documents) Read(ctx context.Context, p identity.Principal, id string) (DocumentView, error) {
	v := DocumentView{}
	if !d.canRead(p) || !documentID(id) {
		return v, ErrDocumentScope
	}
	e := d.pool.QueryRow(ctx, `select document_id::text,name,uploader,original_sha256,profile_sha256,mode,class_id,schema_version from document.original where tenant_id=$1 and document_id=$2 and organization_id=$3`, p.TenantID, id, d.scope.OrganizationID).Scan(&v.ID, &v.Name, &v.Uploader, &v.OriginalSHA, &v.ProfileSHA, &v.Mode, &v.ClassID, &v.SchemaVersion)
	if errors.Is(e, pgx.ErrNoRows) {
		return v, ErrDocumentScope
	}
	if e != nil {
		return v, e
	}
	if v.Uploader != p.Subject && !d.Allowed(p, "documents:review") && !d.Allowed(p, "documents:process") {
		return DocumentView{}, ErrDocumentScope
	}
	v.State = "QUARANTINED"
	var raw []byte
	e = d.pool.QueryRow(ctx, `select evidence_sha256,suggested from document.extraction where tenant_id=$1 and document_id=$2`, p.TenantID, id).Scan(&v.EvidenceSHA, &raw)
	if e == nil {
		if d.validateFields(v, raw) != nil {
			return v, doc.ErrContract
		}
		v.Suggested = append(json.RawMessage(nil), raw...)
		v.State = "REVIEW_REQUIRED"
	} else if !errors.Is(e, pgx.ErrNoRows) {
		return v, e
	} else {
		var terminal *string
		e = d.pool.QueryRow(ctx, `select terminal_error_code from platform.job where tenant_id=$1 and job_id=$2`, p.TenantID, documentJobID(p.TenantID, d.scope.OrganizationID, v.ProfileSHA, id)).Scan(&terminal)
		if e != nil {
			return v, e
		}
		if terminal != nil {
			v.State = "QUARANTINE_TERMINAL"
		}
		return v, nil
	}
	var state approval.State
	e = d.pool.QueryRow(ctx, `select evidence_sha,payload,state from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='document_review'`, p.TenantID, d.requestID(id), d.scope.OrganizationID).Scan(&v.PayloadSHA, &raw, &state)
	if errors.Is(e, pgx.ErrNoRows) {
		return v, nil
	}
	if e != nil {
		return v, e
	}
	var proposal DocumentProposal
	_, h, e := approval.CanonicalPayload(raw)
	if e != nil || h != v.PayloadSHA || json.Unmarshal(raw, &proposal) != nil || proposal.DocumentID != id || proposal.EvidenceSHA != v.EvidenceSHA || proposal.OriginalSHA != v.OriginalSHA || proposal.ProfileSHA != v.ProfileSHA {
		return v, doc.ErrContract
	}
	v.Proposal = &proposal
	if v.Mode == "TYPED_FIXTURE" && (proposal.Schema != "document-review/v2" || proposal.ClassID != v.ClassID || proposal.SchemaVersion != v.SchemaVersion || d.validateFields(v, proposal.Fields) != nil) {
		return v, doc.ErrContract
	}
	v.State = "REVIEW_PENDING"
	if state != approval.StatePending {
		var approved bool
		e = d.pool.QueryRow(ctx, `select reviewer,approved,reason from approval.decision where tenant_id=$1 and request_id=$2 and(select count(*)from approval.decision where tenant_id=$1 and request_id=$2)=1`, p.TenantID, d.requestID(id)).Scan(&v.Reviewer, &approved, &v.Reason)
		if e != nil {
			return v, e
		}
		if approved != (state == approval.StateApproved) {
			return v, doc.ErrContract
		}
		v.State = "REJECTED"
		if approved {
			var bound bool
			e = d.pool.QueryRow(ctx, `select request_id=$3 and payload_sha256=$4 and reviewer=$5 from document.committed where tenant_id=$1 and document_id=$2`, p.TenantID, id, d.requestID(id), v.PayloadSHA, v.Reviewer).Scan(&bound)
			if e != nil || !bound {
				return v, doc.ErrContract
			}
			v.State = "PERSISTED"
		}
	}
	return v, nil
}
func (d *Documents) Original(ctx context.Context, p identity.Principal, id string) (doc.Original, error) {
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return doc.Original{}, e
	}
	original := doc.Original{Name: v.Name, SHA256: v.OriginalSHA, ClassID: v.ClassID, SchemaVersion: v.SchemaVersion}
	e = d.pool.QueryRow(ctx, `select content from document.original where tenant_id=$1 and document_id=$2 and organization_id=$3`, p.TenantID, id, d.scope.OrganizationID).Scan(&original.Bytes)
	if e != nil {
		return original, e
	}
	if original.Validate(v.Mode) != nil {
		return doc.Original{}, doc.ErrContract
	}
	return original, nil
}
func (d *Documents) Process(ctx context.Context, p identity.Principal, id string) (DocumentView, error) {
	if !d.Allowed(p, "documents:process") {
		return DocumentView{}, ErrDocumentScope
	}
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.ProfileSHA != d.scope.ProfileSHA || v.Mode != d.scope.Mode {
		return v, ErrDocumentConflict
	}
	if v.EvidenceSHA != "" {
		return v, nil
	}
	// Existing claim generation/DB clock owns concurrency and stale-worker fencing.
	jobs := NewJobs(d.pool)
	worker := uuid.NewString()
	claimed, e := jobs.ClaimScoped(ctx, "document-extraction", worker, 3*time.Minute, 1, JobScope{d.scope.TenantID, "document-extraction", 1, d.jobPayload(id)})
	if e != nil {
		return v, e
	}
	if len(claimed) != 1 {
		// Preserve a crashed final attempt as terminal evidence, using the
		// existing owner fence. No counter reset or new provider call.
		var exhausted Job
		err := d.pool.QueryRow(ctx, `select tenant_id::text,job_id::text,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job where tenant_id=$1 and job_id=$2 and completed_at is null and terminal_error_code is null and attempts>=max_attempts and claimed_until<clock_timestamp()`, p.TenantID, d.jobID(id)).Scan(&exhausted.TenantID, &exhausted.JobID, &exhausted.Queue, &exhausted.JobType, &exhausted.SchemaVersion, &exhausted.Payload, &exhausted.Attempts, &exhausted.MaxAttempts)
		if err == nil {
			tx, err := d.pool.Begin(ctx)
			if err != nil {
				return v, err
			}
			defer tx.Rollback(ctx)
			if err = ExhaustJobInTx(ctx, tx, exhausted); err != nil {
				return v, err
			}
			if err = tx.Commit(ctx); err != nil {
				return v, err
			}
			return d.Read(ctx, p, id)
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return v, err
		}
		return v, ErrDocumentConflict
	}
	job := claimed[0]
	original, e := d.Original(ctx, p, id)
	if e != nil {
		return v, e
	}
	callCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	evidence, e := d.processor.Run(callCtx, original)
	cancel()
	// If the caller disconnected the finite lease recovers on the next explicit
	// worker invocation. No network error can fabricate an extraction receipt.
	if e != nil {
		return v, d.recordDocumentFailure(ctx, job, worker, id, e)
	}
	if e = validateDocumentEvidence(evidence, original, v.ProfileSHA, v.Mode); e != nil {
		return v, d.recordDocumentFailure(ctx, job, worker, id, e)
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return v, e
	}
	defer tx.Rollback(ctx)
	if e = CompleteJobInTx(ctx, tx, job, worker); e != nil {
		return v, e
	}
	fields, _ := json.Marshal(evidence.Suggested)
	if v.Mode == "TYPED_FIXTURE" {
		fields = evidence.TypedFields
	}
	_, e = tx.Exec(ctx, `insert into document.extraction(tenant_id,document_id,job_id,attempt,security_receipt,provider_response,analysis_receipt,evidence_sha256,suggested)values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, p.TenantID, id, job.JobID, job.Attempts, evidence.Security, evidence.Provider, evidence.Receipt, doc.Hash(evidence.Receipt), fields)
	if e != nil {
		return v, e
	}
	if e = tx.Commit(ctx); e != nil {
		return v, e
	}
	return d.Read(ctx, p, id)
}
func validateDocumentEvidence(e doc.Evidence, o doc.Original, profile, mode string) error {
	if mode == "TYPED_FIXTURE" {
		return doc.ValidateTypedEvidence(e, o, profile)
	}
	var r runtime.Receipt
	if e.Mode != mode || len(e.Security) == 0 || len(e.Security) > 65536 || len(e.Provider) == 0 || len(e.Provider) > 4194304 || len(e.Receipt) > 32768 || json.Unmarshal(e.Receipt, &r) != nil || r.AutomaticStorageAuthorized || r.InputSHA256 != o.SHA256 || r.InputBytes != len(o.Bytes) || r.DocumentProfileSHA256 != profile || r.SecurityReceiptSHA256 != doc.Hash(e.Security) || r.ProviderResponseSHA256 != doc.Hash(e.Provider) || r.PageCount != 1 || r.DocumentClass != doc.Class || r.SDKVersion != "v1.45.0" {
		return doc.ErrContract
	}
	return nil
}
func (d *Documents) recordDocumentFailure(ctx context.Context, job Job, worker, id string, cause error) error {
	code := "EXTRACTION_UNAVAILABLE"
	if errors.Is(cause, doc.ErrSecurity) {
		code = "SECURITY_REJECTED"
	}
	if errors.Is(cause, doc.ErrContract) {
		code = "INVALID_CONTRACT"
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	_, e = FailJobInTx(ctx, tx, job, worker, code, time.Second)
	if e != nil {
		return e
	}
	_, e = tx.Exec(ctx, `insert into document.attempt_failure(tenant_id,document_id,attempt,code)values($1,$2,$3,$4)`, job.TenantID, id, job.Attempts, code)
	if e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	return cause
}
func (d *Documents) Submit(ctx context.Context, p identity.Principal, id, evidenceSHA string, fields doc.Fields) (DocumentView, error) {
	raw, _ := json.Marshal(fields)
	return d.SubmitFields(ctx, p, id, evidenceSHA, raw)
}
func (d *Documents) SubmitFields(ctx context.Context, p identity.Principal, id, evidenceSHA string, fields json.RawMessage) (DocumentView, error) {
	if !d.Allowed(p, "documents:write") {
		return DocumentView{}, ErrDocumentScope
	}
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.Uploader != p.Subject || v.ProfileSHA != d.scope.ProfileSHA || v.EvidenceSHA == "" || v.EvidenceSHA != evidenceSHA || d.validateFields(v, fields) != nil {
		return v, ErrDocumentConflict
	}
	proposal := DocumentProposal{Schema: "document-review/v1", DocumentID: id, OrganizationID: d.scope.OrganizationID, OriginalSHA: v.OriginalSHA, ProfileSHA: v.ProfileSHA, EvidenceSHA: v.EvidenceSHA, Mode: v.Mode, Fields: fields}
	if v.Mode == "TYPED_FIXTURE" {
		proposal.Schema = "document-review/v2"
		proposal.ClassID = v.ClassID
		proposal.SchemaVersion = v.SchemaVersion
	}
	raw, _ := json.Marshal(proposal)
	canonical, h, e := approval.CanonicalPayload(raw)
	if e != nil {
		return v, e
	}
	_, e = d.reviews.Submit(ctx, p, HumanApprovalSpec{approval.Request{TenantID: p.TenantID, ID: d.requestID(id), Kind: approval.KindDocumentReview, SubjectID: id, Requester: p.Subject, EvidenceSHA: h}, d.scope.OrganizationID, canonical}, "documents:write", nil)
	if e != nil {
		return v, e
	}
	return d.Read(ctx, p, id)
}
func (d *Documents) Decide(ctx context.Context, p identity.Principal, id, expectedSHA string, approved bool, reason string) (DocumentView, error) {
	if !d.Allowed(p, "documents:review") || !doc.Text(reason, 2048) {
		return DocumentView{}, ErrDocumentScope
	}
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return v, e
	}
	if v.Proposal == nil || v.ProfileSHA != d.scope.ProfileSHA || v.PayloadSHA != expectedSHA {
		return v, ErrDocumentConflict
	}
	tx, e := d.pool.Begin(ctx)
	if e != nil {
		return v, e
	}
	defer tx.Rollback(ctx)
	_, e = d.reviews.decideTx(ctx, tx, p, p.TenantID, d.requestID(id), d.scope.OrganizationID, expectedSHA, approved, reason, "documents:review", func(ctx context.Context, tx pgx.Tx) error {
		if !approved {
			return nil
		}
		f, _ := json.Marshal(v.Proposal.Fields)
		_, err := tx.Exec(ctx, `insert into document.committed(tenant_id,document_id,request_id,payload_sha256,fields,reviewer)values($1,$2,$3,$4,$5,$6)`, p.TenantID, id, d.requestID(id), expectedSHA, f, p.Subject)
		if err != nil {
			return err
		}
		eventID := uuid.NewSHA1(uuid.NameSpaceOID, []byte("document-committed:"+p.TenantID+":"+id)).String()
		payload, _ := json.Marshal(map[string]string{"document_id": id, "organization_id": d.scope.OrganizationID, "payload_sha256": expectedSHA, "mode": v.Mode})
		_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'document',$3,1,'document.committed',1,clock_timestamp(),$4)`, p.TenantID, eventID, id, payload)
		return err
	})
	if errors.Is(e, approval.ErrNotPending) {
		_ = tx.Rollback(ctx)
		current, readErr := d.Read(ctx, p, id)
		if readErr != nil {
			return current, readErr
		}
		if current.PayloadSHA != expectedSHA || current.Reviewer != p.Subject || current.Reason != reason || (approved && current.State != "PERSISTED") || (!approved && current.State != "REJECTED") {
			return current, ErrDocumentConflict
		}
		return current, nil
	}
	if e != nil {
		return v, e
	}
	if e = tx.Commit(ctx); e != nil {
		return v, e
	}
	return d.Read(ctx, p, id)
}

// EvidencePart returns the original serialized bytes for authorized review;
// fixed part names never become paths or SQL fragments.
func (d *Documents) EvidencePart(ctx context.Context, p identity.Principal, id, part string) ([]byte, error) {
	v, e := d.Read(ctx, p, id)
	if e != nil {
		return nil, e
	}
	if v.EvidenceSHA == "" {
		return nil, ErrDocumentConflict
	}
	var security, provider, receipt []byte
	e = d.pool.QueryRow(ctx, `select security_receipt,provider_response,analysis_receipt from document.extraction where tenant_id=$1 and document_id=$2`, p.TenantID, id).Scan(&security, &provider, &receipt)
	if e != nil {
		return nil, e
	}
	if doc.Hash(receipt) != v.EvidenceSHA {
		return nil, doc.ErrContract
	}
	switch part {
	case "security":
		return security, nil
	case "provider":
		return provider, nil
	case "analysis":
		return receipt, nil
	default:
		return nil, doc.ErrContract
	}
}
