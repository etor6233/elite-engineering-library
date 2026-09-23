// AUTHORED HTTP binding to the existing authenticated role transport and
// connected catalog publication owner. No provider credentials are read.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type CatalogReleaseService interface {
	UploadPNG(context.Context, identity.Principal, string, []byte) (cr.Receipt, error)
	CreateDraft(context.Context, identity.Principal, cr.DraftRequest) (cr.Receipt, error)
	Draft(context.Context, identity.Principal, string) (cr.Draft, error)
	Review(context.Context, identity.Principal, cr.ReviewRequest) (cr.Receipt, error)
	Publish(context.Context, identity.Principal, cr.PublishRequest) (cr.Receipt, error)
	CommandReceipt(context.Context, identity.Principal, string) (cr.Receipt, error)
	Public(context.Context) (cr.Publication, error)
	PublicMedia(context.Context, string) ([]byte, error)
	PublicSearch(context.Context, string) ([]cr.SearchHit, error)
}
type CatalogCurrentService interface {
	Current(context.Context, identity.Principal) (cr.Publication, error)
}
type CatalogSourceService interface {
	CreateSource(context.Context, identity.Principal, cr.SourceRequest) (cr.Receipt, error)
}
type CatalogFeedService interface {
	Send(context.Context, identity.Principal, int64) (cr.FeedStatus, error)
	Status(context.Context, identity.Principal, int64) (cr.FeedStatus, error)
	Reconcile(context.Context, identity.Principal, int64) (cr.FeedStatus, error)
}
type CatalogReleaseModule struct {
	Service          CatalogReleaseService
	Feed             CatalogFeedService
	DeferPublicMedia bool
}

func (m CatalogReleaseModule) DeferredPublicMedia() http.Handler {
	if !m.DeferPublicMedia {
		return nil
	}
	return m.publicMedia()
}

func (m CatalogReleaseModule) publicMedia() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sha := r.PathValue("sha")
		if sha == "" {
			const prefix = "/v1/public/catalog/media/"
			rest := strings.TrimPrefix(r.URL.Path, prefix)
			if rest != r.URL.Path && rest != "" && !strings.Contains(rest, "/") {
				sha = rest
			}
		}
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		if sha == "" || len(r.URL.Query()) != 0 {
			catalogError(w, cr.ErrInvalid)
			return
		}
		raw, e := m.Service.PublicMedia(r.Context(), sha)
		if catalogError(w, e) {
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		etag := `"` + sha + `"`
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			return
		}
		_, _ = w.Write(raw)
	})
}

func catalogError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, cr.ErrInvalid):
		writeProblem(w, 400, "CATALOG_INVALID", "invalid bounded catalog request")
	case errors.Is(err, cr.ErrNotFound):
		writeProblem(w, 404, "CATALOG_NOT_FOUND", "catalog resource not available in this scope")
	case errors.Is(err, cr.ErrConflict), errors.Is(err, approval.ErrNotPending), errors.Is(err, approval.ErrSeparation):
		writeProblem(w, 409, "CATALOG_CONFLICT", "consult the immutable command and current publication")
	case errors.As(err, &pg) && (pg.Code == "23505" || pg.Code == "23514" || pg.Code == "P0001" || pg.Code == "40001" || pg.Code == "40P01"):
		writeProblem(w, 409, "CATALOG_CONFLICT", "consult the immutable command and current publication")
	default:
		writeProblem(w, 503, "CATALOG_UNCONFIRMED", "result unconfirmed; consult the immutable command receipt")
	}
	return true
}
func catalogDecode(w http.ResponseWriter, r *http.Request, v any) bool {
	raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		writeProblem(w, 413, "CATALOG_TOO_LARGE", "bounded command required")
		return false
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		writeProblem(w, 400, "CATALOG_JSON", "unambiguous JSON required")
		return false
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(v) != nil || dec.Decode(new(any)) != io.EOF {
		writeProblem(w, 400, "CATALOG_JSON", "typed command required without extra fields")
		return false
	}
	return true
}
func (m CatalogReleaseModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protect := func(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
		w.Header().Set("Cache-Control", "no-store")
		q := r.URL.Query()
		org := q.Get("organization_id")
		if len(q) != 1 || len(q["organization_id"]) != 1 || !cr.ValidID(org) {
			writeProblem(w, 400, "CATALOG_SCOPE", "one explicit organization required")
			return identity.Principal{}, false
		}
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		p.Permissions = map[string]struct{}{permission: {}}
		p.Organizations = map[string]struct{}{org: {}}
		return p, true
	}
	reply := func(w http.ResponseWriter, v cr.Receipt, e error) {
		if catalogError(w, e) {
			return
		}
		status := 201
		if v.Replay {
			status = 200
		}
		writeJSON(w, status, v)
	}
	mux.HandleFunc("POST /v1/admin/catalog/media/{command}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:draft")
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "image/png" {
			writeProblem(w, 415, "CATALOG_PNG", "PNG required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if e != nil {
			writeProblem(w, 413, "CATALOG_TOO_LARGE", "PNG exceeds limit")
			return
		}
		v, e := m.Service.UploadPNG(r.Context(), p, r.PathValue("command"), raw)
		reply(w, v, e)
	})
	if current, ok := m.Service.(CatalogCurrentService); ok {
		mux.HandleFunc("GET /v1/admin/catalog/current", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "catalog:read")
			if !ok {
				return
			}
			v, e := current.Current(r.Context(), p)
			if catalogError(w, e) {
				return
			}
			value, e := cr.ProjectPublic(v)
			if !catalogError(w, e) {
				writeJSON(w, 200, value)
			}
		})
	}
	if source, ok := m.Service.(CatalogSourceService); ok {
		mux.HandleFunc("POST /v1/admin/catalog/sources", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "catalog:draft")
			if !ok {
				return
			}
			var in cr.SourceRequest
			if !catalogDecode(w, r, &in) {
				return
			}
			v, e := source.CreateSource(r.Context(), p, in)
			reply(w, v, e)
		})
	}
	mux.HandleFunc("POST /v1/admin/catalog/drafts", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:draft")
		if !ok {
			return
		}
		var in cr.DraftRequest
		if !catalogDecode(w, r, &in) {
			return
		}
		v, e := m.Service.CreateDraft(r.Context(), p, in)
		reply(w, v, e)
	})
	mux.HandleFunc("GET /v1/admin/catalog/drafts/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:read")
		if !ok {
			return
		}
		v, e := m.Service.Draft(r.Context(), p, r.PathValue("id"))
		if !catalogError(w, e) {
			writeJSON(w, 200, v)
		}
	})
	mux.HandleFunc("GET /v1/admin/catalog/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:read")
		if !ok {
			return
		}
		v, e := m.Service.CommandReceipt(r.Context(), p, r.PathValue("id"))
		if !catalogError(w, e) {
			writeJSON(w, 200, v)
		}
	})
	for _, stage := range []string{"legal", "technical", "media", "publication"} {
		mux.HandleFunc("POST /v1/admin/catalog/drafts/{id}/review/"+stage, func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "catalog:review:"+stage)
			if !ok {
				return
			}
			var in cr.ReviewRequest
			if !catalogDecode(w, r, &in) {
				return
			}
			if in.DraftID != r.PathValue("id") || in.Stage != stage {
				catalogError(w, cr.ErrInvalid)
				return
			}
			v, e := m.Service.Review(r.Context(), p, in)
			reply(w, v, e)
		})
	}
	mux.HandleFunc("POST /v1/admin/catalog/drafts/{id}/publish", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:publish")
		if !ok {
			return
		}
		var in cr.PublishRequest
		if !catalogDecode(w, r, &in) {
			return
		}
		if in.DraftID != r.PathValue("id") {
			catalogError(w, cr.ErrInvalid)
			return
		}
		v, e := m.Service.Publish(r.Context(), p, in)
		reply(w, v, e)
	})
	public := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		if len(r.URL.Query()) != 0 {
			catalogError(w, cr.ErrInvalid)
			return
		}
		v, e := m.Service.Public(r.Context())
		if catalogError(w, e) {
			return
		}
		etag := fmt.Sprintf(`"catalog-%d-%s"`, v.Generation, v.SHA256)
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			return
		}
		value, e := cr.ProjectPublic(v)
		if !catalogError(w, e) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(value)
		}
	}
	mux.HandleFunc("GET /v1/public/catalog", public)
	mux.HandleFunc("GET /v1/public/catalog/feed", public)
	if !m.DeferPublicMedia {
		mux.Handle("GET /v1/public/catalog/media/{sha}", m.publicMedia())
	}
	mux.HandleFunc("GET /v1/public/catalog/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		q := r.URL.Query()
		if len(q) != 1 || len(q["q"]) != 1 {
			catalogError(w, cr.ErrInvalid)
			return
		}
		rows, e := m.Service.PublicSearch(r.Context(), q.Get("q"))
		if !catalogError(w, e) {
			writeJSON(w, 200, rows)
		}
	})
	if m.Feed != nil {
		for _, route := range []struct{ method, suffix, permission string }{{"POST", "", "catalog:feed"}, {"GET", "", "catalog:read"}, {"POST", "/reconcile", "catalog:feed"}} {
			mux.HandleFunc(route.method+" /v1/admin/catalog/feed/{generation}"+route.suffix, func(w http.ResponseWriter, r *http.Request) {
				p, ok := protect(w, r, route.permission)
				if !ok {
					return
				}
				raw := r.PathValue("generation")
				generation, e := strconv.ParseInt(raw, 10, 64)
				if e != nil || generation < 1 || strconv.FormatInt(generation, 10) != raw {
					catalogError(w, cr.ErrInvalid)
					return
				}
				if r.Body != nil {
					value, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 1))
					if e != nil || len(value) != 0 {
						catalogError(w, cr.ErrInvalid)
						return
					}
				}
				var status cr.FeedStatus
				switch {
				case route.method == "GET":
					status, e = m.Feed.Status(r.Context(), p, generation)
				case route.suffix == "/reconcile":
					status, e = m.Feed.Reconcile(r.Context(), p, generation)
				default:
					status, e = m.Feed.Send(r.Context(), p, generation)
				}
				if !catalogError(w, e) {
					writeJSON(w, 200, status)
				}
			})
		}
	}
}
