package postgres

// AUTHORED reporting projection. Existing profile/auth/immutable ledgers govern data.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5"
	"time"
)

type StoredValueMetricRow struct {
	Operation    string `json:"operation"`
	Operations   string `json:"operations"`
	PointsDelta  string `json:"points_delta"`
	AppliedMinor string `json:"applied_minor_units"`
}
type StoredValueMetrics struct {
	OrganizationID string                 `json:"organization_id"`
	ProgramID      string                 `json:"program_id"`
	ProgramKind    string                 `json:"program_kind"`
	Currency       string                 `json:"currency"`
	ProfileSHA256  string                 `json:"profile_sha256"`
	Source         string                 `json:"source"`
	Basis          string                 `json:"basis"`
	ObservedAt     time.Time              `json:"observed_at"`
	Rows           []StoredValueMetricRow `json:"rows"`
}

func (s *StoredValue) Metrics(ctx context.Context, p identity.Principal, org, programID string) (StoredValueMetrics, error) {
	var empty StoredValueMetrics
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(programID) {
		return empty, sv.ErrBinding
	}
	tenant, bound := s.profile.Scope()
	if org != bound {
		return empty, sv.ErrBinding
	}
	program, e := s.profile.Program(programID)
	if e != nil {
		return empty, e
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "set local statement_timeout='3s'"); e != nil {
		return empty, e
	}
	out := StoredValueMetrics{OrganizationID: org, ProgramID: programID, ProgramKind: program.Calculation.Kind, Currency: program.Calculation.Currency, ProfileSHA256: s.profile.SHA256(), Source: "stored_value.operation + entry + account", Basis: "committed_ledger_by_operation", Rows: []StoredValueMetricRow{}}
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&out.ObservedAt); e != nil {
		return empty, e
	}
	rows, e := tx.Query(ctx, `select op.operation,count(distinct op.operation_id)::text,sum(e.points_delta)::text,sum(e.applied_minor_units::numeric)::text
 from stored_value.operation op join stored_value.entry e on e.tenant_id=op.tenant_id and e.operation_id=op.operation_id
 join stored_value.account a on a.tenant_id=e.tenant_id and a.account_id=e.account_id
 where op.tenant_id=$1 and op.organization_id=$2 and op.program_id=$3 and a.organization_id=$2 and a.program_id=$3 and a.currency=$4 and a.kind=$5
 group by op.operation order by op.operation`, tenant, org, programID, out.Currency, out.ProgramKind)
	if e != nil {
		return empty, e
	}
	for rows.Next() {
		var v StoredValueMetricRow
		if e = rows.Scan(&v.Operation, &v.Operations, &v.PointsDelta, &v.AppliedMinor); e != nil {
			rows.Close()
			return empty, e
		}
		out.Rows = append(out.Rows, v)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return empty, e
	}
	return out, tx.Commit(ctx)
}
