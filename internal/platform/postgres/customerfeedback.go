// AUTHORED PostgreSQL implementation; authentication belongs to the existing HTTP identity boundary.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerFeedback struct{ pool *pgxpool.Pool }

func NewCustomerFeedback(pool *pgxpool.Pool) *CustomerFeedback { return &CustomerFeedback{pool: pool} }

const feedbackEligible = `exists(select 1 from crm.customer_profile c where c.tenant_id=d.tenant_id and c.customer_principal_id=$4 and c.status='active')`

func (r *CustomerFeedback) Definition(ctx context.Context, s customerfeedback.Scope, id string) (customerfeedback.Definition, error) {
	var d customerfeedback.Definition
	e := r.pool.QueryRow(ctx, `select d.survey_id,d.prompt,d.consent_version,d.consent_notice,
 d.active and clock_timestamp()>=d.opens_at and clock_timestamp()<d.closes_at
 from crm.survey_definition d where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3
 and clock_timestamp()<d.retain_until
 and (d.active or exists(select 1 from crm.survey_response a where
 (a.tenant_id,a.organization_id,a.survey_id,a.customer_principal_id)=(d.tenant_id,d.organization_id,d.survey_id,$4)))
 and `+feedbackEligible, s.Tenant, s.Organization, id, s.Customer).
		Scan(&d.ID, &d.Prompt, &d.ConsentVersion, &d.ConsentNotice, &d.Accepting)
	if errors.Is(e, pgx.ErrNoRows) {
		return d, customerfeedback.ErrUnavailable
	}
	return d, e
}
func (r *CustomerFeedback) OwnAnswer(ctx context.Context, s customerfeedback.Scope, id string) (customerfeedback.Answer, error) {
	var a customerfeedback.Answer
	e := r.pool.QueryRow(ctx, `select a.score,a.consent_version,a.received_at from crm.survey_response a
 join crm.survey_definition d using(tenant_id,organization_id,survey_id)
 where a.tenant_id=$1 and a.organization_id=$2 and a.survey_id=$3 and a.customer_principal_id=$4
 and clock_timestamp()<d.retain_until and `+feedbackEligible, s.Tenant, s.Organization, id, s.Customer).
		Scan(&a.Score, &a.ConsentVersion, &a.ReceivedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		return a, customerfeedback.ErrNotFound
	}
	return a, e
}
func (r *CustomerFeedback) Submit(ctx context.Context, s customerfeedback.Scope, id string, in customerfeedback.Submission) (customerfeedback.Result, error) {
	var result customerfeedback.Result
	if e := customerfeedback.ValidateSubmission(s, id, in); e != nil {
		return result, e
	}
	tx, e := r.pool.Begin(ctx)
	if e != nil {
		return result, e
	}
	defer tx.Rollback(ctx)
	// Acquire both definition and customer locks before the final eligibility check.
	// No accepting boolean computed before a lock wait is trusted for admission.
	var version string
	e = tx.QueryRow(ctx, `select d.consent_version
 from crm.survey_definition d join crm.customer_profile c on c.tenant_id=d.tenant_id
 where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3
 and c.customer_principal_id=$4 and c.status='active' for share of d,c`,
		s.Tenant, s.Organization, id, s.Customer).Scan(&version)
	if errors.Is(e, pgx.ErrNoRows) {
		return result, customerfeedback.ErrUnavailable
	}
	if e != nil {
		return result, e
	}
	readExisting := func() error {
		return tx.QueryRow(ctx, `select a.score,a.consent_version,a.received_at
 from crm.survey_response a join crm.survey_definition d using(tenant_id,organization_id,survey_id)
 where a.tenant_id=$1 and a.organization_id=$2 and a.survey_id=$3 and a.customer_principal_id=$4
 and clock_timestamp()<d.retain_until`, s.Tenant, s.Organization, id, s.Customer).
			Scan(&result.Answer.Score, &result.Answer.ConsentVersion, &result.Answer.ReceivedAt)
	}
	checkReplay := func() (customerfeedback.Result, error) {
		if result.Answer.Score != in.Score || result.Answer.ConsentVersion != in.ConsentVersion {
			return result, customerfeedback.ErrConflict
		}
		result.Replay = true
		return result, tx.Commit(ctx)
	}
	// Recovery is a read of the original receipt, including after deactivation/closing,
	// while the customer is still eligible and the retention deadline has not passed.
	e = readExisting()
	if e == nil {
		return checkReplay()
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return result, e
	}
	e = tx.QueryRow(ctx, `insert into crm.survey_response(tenant_id,organization_id,survey_id,customer_principal_id,score,consent_version)
 select $1,$2,$3,$4,$5,d.consent_version from crm.survey_definition d
 where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3 and d.consent_version=$6
 and d.active and clock_timestamp()>=d.opens_at and clock_timestamp()<d.closes_at
 and clock_timestamp()<d.retain_until
 on conflict do nothing returning score,consent_version,received_at`,
		s.Tenant, s.Organization, id, s.Customer, in.Score, in.ConsentVersion).
		Scan(&result.Answer.Score, &result.Answer.ConsentVersion, &result.Answer.ReceivedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		e = readExisting()
		if errors.Is(e, pgx.ErrNoRows) {
			return result, customerfeedback.ErrUnavailable
		}
		if e != nil {
			return result, e
		}
		return checkReplay()
	}
	if e != nil {
		return result, e
	}
	return result, tx.Commit(ctx)
}
func (r *CustomerFeedback) Summary(ctx context.Context, s customerfeedback.Scope, id string) (customerfeedback.Summary, error) {
	var n, p, d, m int64
	e := r.pool.QueryRow(ctx, `select count(a.customer_principal_id),count(*) filter(where a.score>=9),
 count(*) filter(where a.score<=6),d.minimum_responses
 from crm.survey_definition d left join crm.survey_response a using(tenant_id,organization_id,survey_id)
 where d.tenant_id=$1 and d.organization_id=$2 and d.survey_id=$3 and clock_timestamp()<d.retain_until
 group by d.minimum_responses`, s.Tenant, s.Organization, id).Scan(&n, &p, &d, &m)
	if errors.Is(e, pgx.ErrNoRows) {
		return customerfeedback.Summary{}, customerfeedback.ErrUnavailable
	}
	if e != nil {
		return customerfeedback.Summary{}, e
	}
	return customerfeedback.NPS(n, p, d, m)
}

// PurgeExpired removes only expired responses in one explicit authorized tenant/org.
// Caller must schedule this bounded operation under its approved retention policy.
func (r *CustomerFeedback) PurgeExpired(ctx context.Context, tenant, organization string, limit int) (int64, error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 1000 {
		return 0, customerfeedback.ErrInvalid
	}
	tag, e := r.pool.Exec(ctx, `with doomed as (
 select a.tenant_id,a.organization_id,a.survey_id,a.customer_principal_id
 from crm.survey_response a join crm.survey_definition d using(tenant_id,organization_id,survey_id)
 where a.tenant_id=$1 and a.organization_id=$2 and d.retain_until<=clock_timestamp()
 order by a.survey_id,a.customer_principal_id for update of a skip locked limit $3)
 delete from crm.survey_response a using doomed x where
 (a.tenant_id,a.organization_id,a.survey_id,a.customer_principal_id)=
 (x.tenant_id,x.organization_id,x.survey_id,x.customer_principal_id)`, tenant, organization, limit)
	return tag.RowsAffected(), e
}
