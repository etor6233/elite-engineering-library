// AUTHORED durable/authenticated adapter around the existing helpcenter article kernel.
package postgres

import (
	"context"
	hc "elite.local/enterprise/internal/helpcenter"
	cms "elite.local/enterprise/internal/helpcms"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"time"
)

type HelpCMS struct{ pool *pgxpool.Pool }

func NewHelpCMS(pool *pgxpool.Pool) (*HelpCMS, error) {
	if pool == nil {
		return nil, cms.ErrInvalid
	}
	return &HelpCMS{pool}, nil
}

type cmsReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
type cmsHead struct {
	org, locale, category string
	version               int64
}

func cmsError(e error) error {
	if errors.Is(e, pgx.ErrNoRows) {
		return cms.ErrNotFound
	}
	if errors.Is(e, hc.ErrInvalidArticle) {
		return cms.ErrInvalid
	}
	for _, v := range []error{hc.ErrDuplicate, hc.ErrNotFound, hc.ErrVersion, hc.ErrBadTransition, hc.ErrVersionExhausted} {
		if errors.Is(e, v) {
			return cms.ErrConflict
		}
	}
	return e
}
func readCMSHead(ctx context.Context, q cmsReader, tenant, id string, locked bool) (cmsHead, error) {
	var h cmsHead
	sql := `select organization_id,locale,category,current_version from help.article where tenant_id=$1 and article_id=$2`
	if locked {
		sql += " for update"
	}
	e := q.QueryRow(ctx, sql, tenant, id).Scan(&h.org, &h.locale, &h.category, &h.version)
	return h, cmsError(e)
}
func readCMSReceipt(ctx context.Context, q cmsReader, tenant, id string, version int64, command string) (cms.Receipt, error) {
	var v cms.Receipt
	var raw []byte
	sql := `select command_id,action,actor_subject,request_sha256,article_sha256,snapshot from help.revision where tenant_id=$1 and article_id=$2 and version=$3`
	args := []any{tenant, id, version}
	if command != "" {
		sql = `select command_id,action,actor_subject,request_sha256,article_sha256,snapshot from help.revision where tenant_id=$1 and command_id=$2`
		args = []any{tenant, command}
	}
	e := q.QueryRow(ctx, sql, args...).Scan(&v.CommandID, &v.Action, &v.Actor, &v.RequestSHA256, &v.ArticleSHA256, &raw)
	if e != nil {
		return v, cmsError(e)
	}
	if json.Unmarshal(raw, &v.Article) != nil {
		return v, cms.ErrConflict
	}
	_, hash, e := cms.Canonical(v.Article)
	if e != nil || hash != v.ArticleSHA256 {
		return v, cms.ErrConflict
	}
	return v, nil
}
func cmsScope(h cmsHead, a cms.Article, id, org string) bool {
	return h.org == org && a.OrganizationID == org && a.ID == id && a.Locale == h.locale && a.Category == h.category
}
func (s *HelpCMS) Result(ctx context.Context, p identity.Principal, org, command string) (cms.Receipt, error) {
	if !cms.ID(command) || !cms.Editor(p, org) {
		return cms.Receipt{}, cms.ErrNotFound
	}
	v, e := readCMSReceipt(ctx, s.pool, p.TenantID, "", 0, command)
	if e != nil {
		return v, e
	}
	if v.Article.OrganizationID != org || v.Actor != p.Subject || !cms.Authorized(p, cms.Permission(v.Action), org) {
		return cms.Receipt{}, cms.ErrNotFound
	}
	v.Replay = true
	return v, nil
}
func (s *HelpCMS) Article(ctx context.Context, p identity.Principal, org, id string, version int64) (cms.Article, error) {
	var empty cms.Article
	if !cms.ID(id) || version < 0 || !cms.Reader(p, org) {
		return empty, cms.ErrNotFound
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	h, e := readCMSHead(ctx, tx, p.TenantID, id, false)
	if e != nil {
		return empty, e
	}
	if h.org != org {
		return empty, cms.ErrNotFound
	}
	current, e := readCMSReceipt(ctx, tx, p.TenantID, id, h.version, "")
	if e != nil {
		return empty, e
	}
	if !cmsScope(h, current.Article, id, org) {
		return empty, cms.ErrConflict
	}
	editor := cms.Editor(p, org)
	if !editor && current.Article.State != "published" {
		return empty, cms.ErrNotFound
	}
	selected := current
	if version != 0 && version != h.version {
		selected, e = readCMSReceipt(ctx, tx, p.TenantID, id, version, "")
		if e != nil {
			return empty, e
		}
	}
	if !cmsScope(h, selected.Article, id, org) {
		return empty, cms.ErrConflict
	}
	if !editor && selected.Article.State != "published" {
		return empty, cms.ErrNotFound
	}
	return selected.Article, tx.Commit(ctx)
}
func (s *HelpCMS) Execute(ctx context.Context, p identity.Principal, c cms.Command) (cms.Receipt, error) {
	var empty cms.Receipt
	if !c.Valid() {
		return empty, cms.ErrInvalid
	}
	if !cms.Authorized(p, cms.Permission(c.Action), c.OrganizationID) {
		return empty, cms.ErrNotFound
	}
	_, hash, e := cms.Canonical(c)
	if e != nil {
		return empty, e
	}
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('help-command:'||$1::text||':'||$2::text,0))`, p.TenantID, c.CommandID); e != nil {
		return empty, e
	}
	old, e := readCMSReceipt(ctx, tx, p.TenantID, "", 0, c.CommandID)
	if e == nil {
		if old.Actor != p.Subject || old.Action != c.Action || old.Article.OrganizationID != c.OrganizationID || old.Article.ID != c.ArticleID || old.RequestSHA256 != hash {
			return empty, cms.ErrConflict
		}
		old.Replay = true
		return old, tx.Commit(ctx)
	}
	if !errors.Is(e, cms.ErrNotFound) {
		return empty, e
	}
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended('help-article:'||$1::text||':'||$2::text,0))`, p.TenantID, c.ArticleID); e != nil {
		return empty, e
	}
	h, e := readCMSHead(ctx, tx, p.TenantID, c.ArticleID, true)
	if e != nil && !errors.Is(e, cms.ErrNotFound) {
		return empty, e
	}
	exists := e == nil
	var article hc.Article
	if c.Action == "create" {
		if exists {
			return empty, cms.ErrConflict
		}
		h = cmsHead{org: c.OrganizationID, locale: c.Locale, category: c.Category, version: 1}
		article = hc.Article{TenantID: p.TenantID, ID: c.ArticleID, Category: c.Category, Title: c.Title, Body: c.Body}
	} else {
		if !exists || h.org != c.OrganizationID {
			return empty, cms.ErrNotFound
		}
		saved, e := readCMSReceipt(ctx, tx, p.TenantID, c.ArticleID, h.version, "")
		if e != nil {
			return empty, e
		}
		if !cmsScope(h, saved.Article, c.ArticleID, c.OrganizationID) {
			return empty, cms.ErrConflict
		}
		v := saved.Article
		article = hc.Article{TenantID: p.TenantID, ID: c.ArticleID, Category: v.Category, Title: v.Title, Body: v.Body, State: hc.State(v.State), Version: h.version, UpdatedAt: v.UpdatedAt}
	}
	var now time.Time
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&now); e != nil {
		return empty, e
	}
	expected, _ := strconv.ParseInt(c.Version, 10, 64)
	var next hc.Article
	switch c.Action {
	case "create":
		next, e = hc.NewDraft(article, now)
	case "update":
		// Narrow CMS profile: published/archived snapshots are never edited.
		if article.State != hc.StateDraft {
			return empty, cms.ErrConflict
		}
		next, e = hc.ReviseArticle(article, c.Title, c.Body, expected, now)
	case "publish":
		next, e = hc.PublishArticle(article, expected, now)
	case "archive":
		next, e = hc.ArchiveArticle(article, expected, now)
	}
	if e != nil {
		return empty, cmsError(e)
	}
	view := cms.Article{ID: next.ID, OrganizationID: h.org, Locale: h.locale, Category: next.Category, Title: next.Title, Body: next.Body, State: string(next.State), Version: strconv.FormatInt(next.Version, 10), UpdatedAt: next.UpdatedAt}
	payload, payloadHash, e := cms.Canonical(view)
	if e != nil {
		return empty, e
	}
	if !exists {
		if _, e = tx.Exec(ctx, `insert into help.article(tenant_id,article_id,organization_id,locale,category,current_version)values($1,$2,$3,$4,$5,1)`, p.TenantID, c.ArticleID, h.org, h.locale, h.category); e != nil {
			return empty, e
		}
	}
	_, e = tx.Exec(ctx, `insert into help.revision(tenant_id,article_id,version,command_id,action,actor_subject,request_sha256,article_sha256,snapshot,title,body,state)values($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10,$11,$12)`, p.TenantID, c.ArticleID, next.Version, c.CommandID, c.Action, p.Subject, hash, payloadHash, string(payload), view.Title, view.Body, view.State)
	if e != nil {
		return empty, e
	}
	if exists {
		tag, e := tx.Exec(ctx, `update help.article set current_version=$3 where tenant_id=$1 and article_id=$2 and current_version=$4`, p.TenantID, c.ArticleID, next.Version, h.version)
		if e != nil {
			return empty, e
		}
		if tag.RowsAffected() != 1 {
			return empty, cms.ErrConflict
		}
	}
	if e = outbox(ctx, tx, p.TenantID, randomid.Generator{}.New(), "help-article", c.ArticleID, next.Version, "help.article-"+c.Action, `jsonb_build_object('organization_id',$7::text,'article_sha256',$8::text)`, h.org, payloadHash); e != nil {
		return empty, e
	}
	receipt := cms.Receipt{CommandID: c.CommandID, Action: c.Action, Actor: p.Subject, RequestSHA256: hash, ArticleSHA256: payloadHash, Article: view}
	return receipt, tx.Commit(ctx)
}
func (s *HelpCMS) List(ctx context.Context, p identity.Principal, org, locale, query, after string) (cms.Page, error) {
	out := cms.Page{Items: []cms.Summary{}}
	if !cms.Reader(p, org) || (locale != "es" && locale != "en") || (after != "" && !cms.ID(after)) || len(query) > 160 || (query != "" && !cms.Text(query, 160, false)) {
		return out, cms.ErrNotFound
	}
	editor := cms.Editor(p, org)
	rows, e := s.pool.Query(ctx, `select h.article_id,h.locale,h.category,r.title,r.state,r.version::text
 from help.article h join help.revision r on r.tenant_id=h.tenant_id and r.article_id=h.article_id and r.version=h.current_version
 where h.tenant_id=$1 and h.organization_id=$2 and h.locale=$3 and h.article_id>$4
 and($5::boolean or r.state='published')and($6=''or strpos(lower((r.title||E'\n'||r.body)collate pg_catalog.pg_unicode_fast),lower($6::text collate pg_catalog.pg_unicode_fast))>0)
 order by h.article_id limit 51`, p.TenantID, org, locale, after, editor, strings.TrimSpace(query))
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var v cms.Summary
		if e = rows.Scan(&v.ID, &v.Locale, &v.Category, &v.Title, &v.State, &v.Version); e != nil {
			return out, e
		}
		out.Items = append(out.Items, v)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.Items) > 50 {
		out.Items = out.Items[:50]
		out.Next = out.Items[49].ID
	}
	return out, nil
}
func (s *HelpCMS) History(ctx context.Context, p identity.Principal, org, id string, before int64) (cms.Page, error) {
	out := cms.Page{Items: []cms.Summary{}}
	if !cms.ID(id) || before < 0 || !cms.Editor(p, org) {
		return out, cms.ErrNotFound
	}
	rows, e := s.pool.Query(ctx, `select h.article_id,h.locale,h.category,r.title,r.state,r.version::text
 from help.article h join help.revision r on r.tenant_id=h.tenant_id and r.article_id=h.article_id
 where h.tenant_id=$1 and h.organization_id=$2 and h.article_id=$3 and($4::bigint=0 or r.version<$4)order by r.version desc limit 51`, p.TenantID, org, id, before)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var v cms.Summary
		if e = rows.Scan(&v.ID, &v.Locale, &v.Category, &v.Title, &v.State, &v.Version); e != nil {
			return out, e
		}
		out.Items = append(out.Items, v)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.Items) > 50 {
		out.Items = out.Items[:50]
		out.Next = out.Items[49].Version
	}
	return out, nil
}
