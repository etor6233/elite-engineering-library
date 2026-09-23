package postgres

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"errors"
	"github.com/jackc/pgx/v5"
)

// PublicMedia only resolves an immutable normalized asset in the currently
// eligible publication. Knowing another draft's hash does not publish its media.
func (s *CatalogRelease) PublicMedia(ctx context.Context, hash string) ([]byte, error) {
	if !cr.ValidSHA(hash) {
		return nil, cr.ErrNotFound
	}
	var raw []byte
	err := s.pool.QueryRow(ctx, `select a.png from catalog.release_media a
 join catalog.release_public_current c using(tenant_id)
 where a.tenant_id=$1 and a.sha256=$2 and exists(select 1 from jsonb_array_elements(c.snapshot->'models') m
 where m->'media'->>'id'=a.media_id and m->'media'->>'sha256'=a.sha256) limit 1`, s.profile.TenantID, hash).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, cr.ErrNotFound
	}
	return raw, err
}
func (s *CatalogRelease) PublicSearch(ctx context.Context, q string) ([]cr.SearchHit, error) {
	if !cr.ValidText(q, 200) {
		return nil, cr.ErrInvalid
	}
	rows, err := s.pool.Query(ctx, `select x.external_id,x.title,x.body,c.generation,c.snapshot_sha256
 from search.document x join catalog.release_public_current c on c.tenant_id=x.tenant_id and c.snapshot_sha256=x.content_sha256
 cross join websearch_to_tsquery('simple',$2) query
 where x.tenant_id=$1 and x.kind='catalog-model' and x.tsv@@query
 order by ts_rank(x.tsv,query) desc,x.external_id limit 20`, s.profile.TenantID, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []cr.SearchHit{}
	for rows.Next() {
		var hit cr.SearchHit
		if err = rows.Scan(&hit.ModelID, &hit.Title, &hit.URL, &hit.Generation, &hit.SourceSHA256); err != nil {
			return nil, err
		}
		out = append(out, hit)
	}
	return out, rows.Err()
}
