package postgres

// AUTHORED provider evidence binding. Raw SDK bytes survive PostgreSQL JSONB
// whitespace changes; observed processing is separate from input acknowledgement.
import (
	"context"
	mb "elite.local/enterprise/internal/merchantbridge"
)

func (s *MerchantPublication) recordMerchantObservation(ctx context.Context, r mb.Intent, result mb.Result, confirmed bool) error {
	if result.Operation != "INSERT" && result.Operation != "GET" || result.State != "OBSERVED" && result.State != "NOT_FOUND" && result.State != "REJECTED" && result.State != "UNKNOWN" || confirmed && result.State != "OBSERVED" || result.Operation == "INSERT" && !confirmed {
		return mb.ErrBinding
	}
	raw := []byte("{}")
	if result.State == "OBSERVED" {
		if result.Resource != r.Resource || len(result.Response) == 0 || mb.Hash(result.Response) != result.ResponseSHA256 {
			return mb.ErrBinding
		}
		raw = result.Response
	}
	command, e := s.catalog.pool.Exec(ctx, `insert into catalog.merchant_observation(tenant_id,approval_id,operation,result_state,confirmed,response,response_sha256)
 select i.tenant_id,i.delivery_key,$3,$4,$5,$6,$7 from catalog.merchant_effect i
 join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.delivery_key=$2 and i.profile_sha256=$8 and i.source_sha256=$9
 and d.state in ('sending','unknown','accepted')
 on conflict(tenant_id,approval_id,operation,result_state,response_sha256) do nothing`, r.TenantID, r.ApprovalID, result.Operation, result.State, confirmed, raw, mb.Hash(raw), r.ProfileSHA256, r.SourceSHA256)
	if e != nil {
		return e
	}
	if command.RowsAffected() == 1 {
		return nil
	}
	var exists bool
	e = s.catalog.pool.QueryRow(ctx, `select exists(select 1 from catalog.merchant_observation o join catalog.merchant_effect i
 on i.tenant_id=o.tenant_id and i.channel_code=o.channel_code and i.delivery_key=o.approval_id
 where o.tenant_id=$1 and o.approval_id=$2 and o.operation=$3 and o.result_state=$4 and o.confirmed=$5
 and o.response_sha256=$6 and i.profile_sha256=$7 and i.source_sha256=$8)`, r.TenantID, r.ApprovalID, result.Operation, result.State, confirmed, mb.Hash(raw), r.ProfileSHA256, r.SourceSHA256).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		return mb.ErrBinding
	}
	return nil
}
