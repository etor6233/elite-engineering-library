package main

// AUTHORED opt-in composition. No environment value grants learner permissions.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/trainingbridge"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

func init() { trainingModuleFactory = selectedTrainingModule }
func selectedTrainingModule(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	if getenv == nil {
		return nil, trainingbridge.ErrContract
	}
	enabled := getenv("TRAINING_ENABLED")
	if enabled == "" || enabled == "false" {
		return nil, nil
	}
	if enabled != "true" {
		return nil, trainingbridge.ErrContract
	}
	if pool == nil {
		return nil, trainingbridge.ErrContract
	}
	var ready bool
	e := pool.QueryRow(ctx, `select
 (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('audit.event') and tgname='training_audit_immutable'
 or tgrelid=to_regclass('approval.request') and tgname='training_assessment_guard'))=2
 and to_regclass('audit.audit_training_attempt_idx') is not null
 and to_regclass('approval.approval_training_review_idx') is not null`).Scan(&ready)
	if e == nil && ready {
		e = pool.QueryRow(ctx, `select to_regclass('approval.training_assessment_attempt_idx') is not null`).Scan(&ready)
	}
	if e != nil || !ready {
		return nil, trainingbridge.ErrContract
	}
	revision, e := strconv.Atoi(getenv("TRAINING_PROFILE_REVISION"))
	if e != nil {
		return nil, trainingbridge.ErrContract
	}
	profile, e := trainingbridge.LoadProfile(getenv("TRAINING_PROFILE_FILE"), getenv("TRAINING_CONTENT_FILE"), trainingbridge.Activation{Enabled: true, ID: getenv("TRAINING_PROFILE_ID"), Revision: revision, SHA256: getenv("TRAINING_PROFILE_SHA256"), TenantID: getenv("TRAINING_TENANT_ID"), OrganizationID: getenv("TRAINING_ORGANIZATION_ID")})
	if e != nil {
		return nil, e
	}
	store, e := trainingbridge.NewStore(pool, profile)
	if e != nil {
		return nil, e
	}
	return &trainingbridge.Module{Store: store}, nil
}
