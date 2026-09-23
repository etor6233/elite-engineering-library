package postgres

// AUTHORED bounded, scoped catalog/inbox composition. Cursor grants no access.
import (
	"bytes"
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/base64"
	"encoding/json"
	"time"
)

func (d *Documents) validateFields(v DocumentView, raw []byte) error {
	if v.Mode == "TYPED_FIXTURE" {
		return doc.ValidateTypedFields(v.ClassID, v.SchemaVersion, raw)
	}
	var f doc.Fields
	de := json.NewDecoder(bytes.NewReader(raw))
	de.DisallowUnknownFields()
	if de.Decode(&f) != nil {
		return doc.ErrContract
	}
	return f.Validate()
}
func (d *Documents) Catalog(p identity.Principal) (doc.ClassCatalog, error) {
	if !d.canRead(p) {
		return doc.ClassCatalog{}, ErrDocumentScope
	}
	return doc.Catalog(d.scope.Mode, d.scope.ProfileSHA), nil
}

type DocumentInbox struct {
	Schema     string         `json:"schema"`
	Items      []DocumentView `json:"items"`
	NextCursor *string        `json:"next_cursor"`
}
type documentCursor struct {
	At    time.Time `json:"at"`
	ID    string    `json:"id"`
	Scope string    `json:"scope"`
}

func (d *Documents) List(ctx context.Context, p identity.Principal, cursor string, limit int) (DocumentInbox, error) {
	out := DocumentInbox{Schema: "document-inbox/v1", Items: []DocumentView{}}
	if !d.canRead(p) {
		return out, ErrDocumentScope
	}
	if limit < 1 || limit > 50 || len(cursor) > 1024 {
		return out, doc.ErrContract
	}
	privileged := d.Allowed(p, "documents:review") || d.Allowed(p, "documents:process")
	scopeBytes, _ := json.Marshal([]any{p.TenantID, d.scope.OrganizationID, p.Subject, privileged})
	scope := doc.Hash(scopeBytes)
	c := documentCursor{At: time.Unix(0, 0).UTC(), ID: "00000000-0000-0000-0000-000000000000", Scope: scope}
	if cursor != "" {
		b, e := base64.RawURLEncoding.DecodeString(cursor)
		if e != nil || json.Unmarshal(b, &c) != nil || !documentID(c.ID) || c.Scope != scope || c.At.IsZero() {
			return out, doc.ErrContract
		}
	}
	rows, e := d.pool.Query(ctx, `select document_id::text,received_at from document.original where tenant_id=$1 and organization_id=$2 and ($3 or uploader=$4) and (received_at,document_id)>($5,$6::uuid) order by received_at,document_id limit $7`, p.TenantID, d.scope.OrganizationID, privileged, p.Subject, c.At, c.ID, limit+1)
	if e != nil {
		return out, e
	}
	type ref struct {
		id string
		at time.Time
	}
	refs := []ref{}
	for rows.Next() {
		var x ref
		if e = rows.Scan(&x.id, &x.at); e != nil {
			rows.Close()
			return out, e
		}
		refs = append(refs, x)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return out, e
	}
	more := len(refs) > limit
	if more {
		refs = refs[:limit]
	}
	for _, x := range refs {
		v, e := d.Read(ctx, p, x.id)
		if e != nil {
			return out, e
		}
		out.Items = append(out.Items, v)
	}
	if more {
		last := refs[len(refs)-1]
		b, _ := json.Marshal(documentCursor{At: last.at, ID: last.id, Scope: scope})
		next := base64.RawURLEncoding.EncodeToString(b)
		out.NextCursor = &next
	}
	return out, nil
}
