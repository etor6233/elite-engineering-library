// Package helpcenter provides a tenant-scoped knowledge base: articles with
// categories, versioned updates (optimistic concurrency) and a state machine
// draft → published → archived. Search is a local substring scan.
package helpcenter

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
)

// State is the publication lifecycle.
type State string

const (
	StateDraft     State = "draft"
	StatePublished State = "published"
	StateArchived  State = "archived"
)

var (
	categoryRe = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

	ErrInvalidArticle   = errors.New("helpcenter: invalid article")
	ErrDuplicate        = errors.New("helpcenter: duplicate article")
	ErrNotFound         = errors.New("helpcenter: not found")
	ErrVersion          = errors.New("helpcenter: version conflict")
	ErrBadTransition    = errors.New("helpcenter: bad state transition")
	ErrVersionExhausted = errors.New("helpcenter: version exhausted")
)

// Article is a single help document.
type Article struct {
	TenantID  string
	ID        string
	Category  string
	Title     string
	Body      string
	State     State
	Version   int64
	UpdatedAt time.Time
}

// Store holds articles.
type Store struct {
	mu       sync.Mutex
	articles map[articleKey]Article
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{articles: make(map[articleKey]Article)}
}

type articleKey struct{ tenant, id string }

func key(tenant, id string) articleKey { return articleKey{tenant, id} }

func validate(a Article) error {
	if strings.TrimSpace(a.TenantID) == "" || strings.TrimSpace(a.ID) == "" {
		return ErrInvalidArticle
	}
	if !categoryRe.MatchString(a.Category) {
		return ErrInvalidArticle
	}
	if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Body) == "" {
		return ErrInvalidArticle
	}
	return nil
}

// Create adds a draft article (version 1).
func (s *Store) Create(a Article) error {
	if err := validate(a); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[key(a.TenantID, a.ID)]; ok {
		return ErrDuplicate
	}
	a = prepareDraft(a, time.Now().UTC())
	s.articles[key(a.TenantID, a.ID)] = a
	return nil
}

// Update edits title/body with optimistic concurrency (expected version).
func (s *Store) Update(tenant, id string, title, body string, expectedVersion int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	updated, err := ReviseArticle(a, title, body, expectedVersion, time.Now().UTC())
	if err != nil {
		return err
	}
	s.articles[key(tenant, id)] = updated
	return nil
}

// Publish marks draft → published (version-gated).
func (s *Store) Publish(tenant, id string, expectedVersion int64) error {
	return s.transition(tenant, id, expectedVersion, StateDraft, StatePublished)
}

// Archive marks published → archived (version-gated).
func (s *Store) Archive(tenant, id string, expectedVersion int64) error {
	return s.transition(tenant, id, expectedVersion, StatePublished, StateArchived)
}

func (s *Store) transition(tenant, id string, expectedVersion int64, from, to State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	updated, err := transitionArticle(a, expectedVersion, from, to, time.Now().UTC())
	if err != nil {
		return err
	}
	s.articles[key(tenant, id)] = updated
	return nil
}

// Published returns tenant-scoped published articles, optionally by category.
func (s *Store) Published(tenant, category string) []Article {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Article
	for _, a := range s.articles {
		if a.TenantID == tenant && a.State == StatePublished {
			if category == "" || a.Category == category {
				out = append(out, a)
			}
		}
	}
	return out
}

// Search returns published articles whose title/body contain the query
// (case-insensitive). No external search adapter is connected.
func (s *Store) Search(tenant, query string) []Article {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	var out []Article
	for _, a := range s.Published(tenant, "") {
		if strings.Contains(strings.ToLower(a.Title), q) || strings.Contains(strings.ToLower(a.Body), q) {
			out = append(out, a)
		}
	}
	return out
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("helpcenter(%d)", len(s.articles))
}

// AUTHORED storage boundary extraction. Existing source validation/version/state
// rules retained; caller loads a valid article and owns authorization/persistence.
func prepareDraft(a Article, now time.Time) Article {
	a.State = StateDraft
	a.Version = 1
	a.UpdatedAt = now
	return a
}
func NewDraft(a Article, now time.Time) (Article, error) {
	if err := validate(a); err != nil {
		return Article{}, err
	}
	return prepareDraft(a, now), nil
}
func ReviseArticle(a Article, title, body string, expectedVersion int64, now time.Time) (Article, error) {
	if a.Version != expectedVersion {
		return Article{}, ErrVersion
	}
	if strings.TrimSpace(title) == "" || strings.TrimSpace(body) == "" {
		return Article{}, ErrInvalidArticle
	}
	if a.Version == math.MaxInt64 {
		return Article{}, ErrVersionExhausted
	}
	a.Title = title
	a.Body = body
	a.Version++
	a.UpdatedAt = now
	return a, nil
}

func PublishArticle(a Article, expectedVersion int64, now time.Time) (Article, error) {
	return transitionArticle(a, expectedVersion, StateDraft, StatePublished, now)
}
func ArchiveArticle(a Article, expectedVersion int64, now time.Time) (Article, error) {
	return transitionArticle(a, expectedVersion, StatePublished, StateArchived, now)
}
func transitionArticle(a Article, expectedVersion int64, from, to State, now time.Time) (Article, error) {
	if a.Version != expectedVersion {
		return Article{}, ErrVersion
	}
	if a.State != from {
		return Article{}, ErrBadTransition
	}
	if a.Version == math.MaxInt64 {
		return Article{}, ErrVersionExhausted
	}
	a.State = to
	a.Version++
	a.UpdatedAt = now
	return a, nil
}
