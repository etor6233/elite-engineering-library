package postgres

// AUTHORED scoped, bounded operator reads; monetary calculations remain in the
// source-derived process and allocations remain in the existing SQL projection.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
)

func (s *StoredValue) ProfileDocument(ctx context.Context, p identity.Principal, org string) (json.RawMessage, error) {
	if !s.allowed(p, "stored_value:read") {
		return nil, sv.ErrBinding
	}
	_, bound := s.profile.Scope()
	if org != bound {
		return nil, sv.ErrBinding
	}
	return s.profile.Document(), nil
}

type StoredValueApprovalPage struct {
	IDs  []string `json:"ids"`
	Next string   `json:"next_cursor,omitempty"`
}

func (s *StoredValue) Approvals(ctx context.Context, p identity.Principal, order, after string) (StoredValueApprovalPage, error) {
	out := StoredValueApprovalPage{IDs: []string{}}
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(order) || after != "" && !sv.ValidID(after) {
		return out, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	rows, e := s.pool.Query(ctx, `select request_id from approval.request where tenant_id=$1 and organization_id=$2 and kind='stored_value_operation' and subject_id=$3 and request_id>$4 order by request_id limit 51`, tenant, org, order, after)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if e = rows.Scan(&id); e != nil {
			return out, e
		}
		out.IDs = append(out.IDs, id)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.IDs) > 50 {
		out.IDs = out.IDs[:50]
		out.Next = out.IDs[49]
	}
	return out, nil
}

func (s *StoredValue) Scope() (string, string) {
	if s == nil {
		return "", ""
	}
	return s.profile.Scope()
}
