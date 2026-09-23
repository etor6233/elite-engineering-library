package helpcenter

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func art() Article {
	return Article{TenantID: "t", ID: "a1", Category: "faq", Title: "Cómo reservar", Body: "Elegí un turno disponible."}
}

func TestCreatePublishArchive(t *testing.T) {
	s := NewStore()
	if err := s.Create(art()); err != nil {
		t.Fatal(err)
	}
	if len(s.Published("t", "")) != 0 {
		t.Fatal("draft should not be published")
	}
	if err := s.Publish("t", "a1", 1); err != nil {
		t.Fatal(err)
	}
	if len(s.Published("t", "")) != 1 {
		t.Fatal("published should be visible")
	}
	if err := s.Archive("t", "a1", 2); err != nil {
		t.Fatal(err)
	}
	if len(s.Published("t", "")) != 0 {
		t.Fatal("archived should not be published")
	}
}

func TestOptimisticConcurrency(t *testing.T) {
	s := NewStore()
	_ = s.Create(art())
	if err := s.Update("t", "a1", "Nuevo", "Contenido", 99); !errors.Is(err, ErrVersion) {
		t.Fatalf("stale version accepted: %v", err)
	}
	if err := s.Update("t", "a1", "Nuevo", "Contenido", 1); err != nil {
		t.Fatal(err)
	}
	// now version is 2
	if err := s.Publish("t", "a1", 2); err != nil {
		t.Fatal(err)
	}
}

func TestBadTransition(t *testing.T) {
	s := NewStore()
	_ = s.Create(art())
	_ = s.Publish("t", "a1", 1)
	if err := s.Publish("t", "a1", 2); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("double publish accepted: %v", err)
	}
}

func TestSearchPublishedOnly(t *testing.T) {
	s := NewStore()
	_ = s.Create(art())
	_ = s.Create(Article{TenantID: "t", ID: "a2", Category: "faq", Title: "Devoluciones", Body: "Cómo devolver"})
	// only a1 is published
	_ = s.Publish("t", "a1", 1)
	results := s.Search("t", "reservar")
	if len(results) != 1 || results[0].ID != "a1" {
		t.Fatalf("search should return only published a1, got %+v", results)
	}
	if len(s.Search("t", "devolver")) != 0 {
		t.Fatal("draft a2 should not appear in search")
	}
}

func TestInvalidAndTenantIsolation(t *testing.T) {
	s := NewStore()
	bad := art()
	bad.Category = "BAD CAT"
	if err := s.Create(bad); !errors.Is(err, ErrInvalidArticle) {
		t.Fatalf("bad category accepted: %v", err)
	}
	_ = s.Create(art())
	if len(s.Published("other", "")) != 0 {
		t.Fatal("cross-tenant leak")
	}
}

func TestDistinctArticleTuples(t *testing.T) {
	s := NewStore()
	a, b := art(), art()
	a.TenantID = "a\x00b"
	a.ID = "c"
	b.TenantID = "a"
	b.ID = "b\x00c"
	if e := s.Create(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Create(b); e != nil {
		t.Fatal("distinct article rejected", e)
	}
}
func TestForeignArticleUpdateRejected(t *testing.T) {
	s := NewStore()
	a := art()
	a.TenantID = "a\x00b"
	a.ID = "c"
	if e := s.Create(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Update("a", "b\x00c", "foreign", "foreign", 1); !errors.Is(e, ErrNotFound) {
		t.Fatal("foreign edit admitted", e)
	}
	if e := s.Publish(a.TenantID, a.ID, 1); e != nil {
		t.Fatal("foreign edit consumed version", e)
	}
	if rows := s.Published(a.TenantID, ""); len(rows) != 1 || rows[0].Title != a.Title {
		t.Fatal("owner article changed", rows)
	}
}
func TestForeignArticleTransitionsRejected(t *testing.T) {
	for _, archive := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "archive"}[archive], func(t *testing.T) {
			s := NewStore()
			a := art()
			a.TenantID = "a\x00b"
			a.ID = "c"
			if e := s.Create(a); e != nil {
				t.Fatal(e)
			}
			var e error
			if archive {
				if e = s.Publish(a.TenantID, a.ID, 1); e != nil {
					t.Fatal(e)
				}
				e = s.Archive("a", "b\x00c", 2)
			} else {
				e = s.Publish("a", "b\x00c", 1)
			}
			if !errors.Is(e, ErrNotFound) {
				t.Fatal("foreign transition admitted", e)
			}
			for _, row := range s.articles {
				want := StateDraft
				if archive {
					want = StatePublished
				}
				if row.State != want {
					t.Fatal("owner state changed")
				}
			}
		})
	}
}
func TestVersionExhaustionRejected(t *testing.T) {
	for _, op := range []string{"update", "publish", "archive"} {
		t.Run(op, func(t *testing.T) {
			s := NewStore()
			if e := s.Create(art()); e != nil {
				t.Fatal(e)
			}
			// Synthetic internal boundary: no public API initializes a MaxInt64 version.
			var before Article
			for k, a := range s.articles {
				a.Version = math.MaxInt64
				if op == "archive" {
					a.State = StatePublished
				}
				s.articles[k] = a
				before = a
			}
			var e error
			switch op {
			case "update":
				e = s.Update("t", "a1", "new", "body", math.MaxInt64)
			case "publish":
				e = s.Publish("t", "a1", math.MaxInt64)
			case "archive":
				e = s.Archive("t", "a1", math.MaxInt64)
			}
			if !errors.Is(e, ErrVersionExhausted) {
				t.Error("exhausted version accepted", e)
			}
			for _, a := range s.articles {
				if a != before {
					t.Error("failed operation changed article")
				}
			}
		})
	}
}

func TestVersionBoundaryAndErrorPrecedence(t *testing.T) {
	s := NewStore()
	if e := s.Create(art()); e != nil {
		t.Fatal(e)
	}
	for k, a := range s.articles {
		a.Version = math.MaxInt64 - 1
		s.articles[k] = a
	}
	if e := s.Update("t", "a1", "last", "body", math.MaxInt64-1); e != nil {
		t.Fatal(e)
	}
	if e := s.Update("t", "a1", "x", "body", 1); !errors.Is(e, ErrVersion) {
		t.Fatal("stale precedence", e)
	}
	if e := s.Update("t", "a1", " ", "body", math.MaxInt64); !errors.Is(e, ErrInvalidArticle) {
		t.Fatal("invalid precedence", e)
	}
	if e := s.Archive("t", "a1", math.MaxInt64); !errors.Is(e, ErrBadTransition) {
		t.Fatal("state precedence", e)
	}
	if e := s.Publish("missing", "a1", math.MaxInt64); !errors.Is(e, ErrNotFound) {
		t.Fatal("lookup precedence", e)
	}
	if e := s.Update("t", "a1", "x", "body", math.MaxInt64); !errors.Is(e, ErrVersionExhausted) {
		t.Fatal("overflow accepted", e)
	}
}
func TestConcurrentVersionHasSingleWinner(t *testing.T) {
	s := NewStore()
	if e := s.Create(art()); e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	var winners atomic.Int32
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e := s.Update("t", "a1", fmt.Sprint(i), "body", 1)
			if e == nil {
				winners.Add(1)
			} else if !errors.Is(e, ErrVersion) {
				t.Error(e)
			}
			_ = s.String()
			s.Published("t", "")
			s.Search("t", "body")
		}(i)
	}
	wg.Wait()
	if winners.Load() != 1 {
		t.Fatal("optimistic concurrency lost", winners.Load())
	}
	if e := s.Publish("t", "a1", 2); e != nil {
		t.Fatal(e)
	}
}
func TestPublishedSearchCopyAndArchivedUpdate(t *testing.T) {
	s := NewStore()
	if e := s.Create(art()); e != nil {
		t.Fatal(e)
	}
	if e := s.Publish("t", "a1", 1); e != nil {
		t.Fatal(e)
	}
	rows := s.Search("t", " RESERVAR ")
	if len(rows) != 1 {
		t.Fatal("substring search", rows)
	}
	rows[0].Title = "mutated"
	if s.Published("t", "")[0].Title != art().Title {
		t.Fatal("returned article aliases stored value")
	}
	if len(s.Published("t", "other")) != 0 || len(s.Search("t", " ")) != 0 || len(s.Search("other", "reservar")) != 0 {
		t.Fatal("filter failure")
	}
	if e := s.Archive("t", "a1", 2); e != nil {
		t.Fatal(e)
	}
	if e := s.Update("t", "a1", "archived edit", "body", 3); e != nil {
		t.Fatal("existing archived edit semantics changed", e)
	}
	if len(s.Published("t", "")) != 0 {
		t.Fatal("archived edit republished article")
	}
}
func FuzzArticleHistoryModel(f *testing.F) {
	for _, prefix := range []string{"a", "x\x00y", "ñ", " "} {
		for _, ops := range [][]byte{{0, 0, 0, 1, 2, 0, 3, 0}, {0, 0, 0, 1, 1, 0, 2, 0, 3, 0}, {0, 0, 4, 0, 1, 0, 2, 0}, {0, 0, 2, 0, 4, 0, 3, 0}} {
			f.Add(prefix, ops)
		}
	}
	f.Fuzz(func(t *testing.T, prefix string, ops []byte) {
		if len(prefix) > 64 {
			return
		}
		if len(ops) > 128 {
			ops = ops[:128]
		}
		s := NewStore()
		model := map[string]map[string]Article{}
		tenants := []string{"t" + prefix + "\x00b", "t" + prefix, "other", "t" + prefix}
		ids := []string{"c", "b\x00c", "c", "c"}
		for i := 0; i+1 < len(ops); i += 2 {
			op := ops[i] % 5
			slot := int(ops[i+1] % 4)
			tenant, id := tenants[slot], ids[slot]
			if model[tenant] == nil {
				model[tenant] = map[string]Article{}
			}
			a, exists := model[tenant][id]
			expected := a.Version
			switch (ops[i+1] / 4) % 4 {
			case 1:
				expected--
			case 2:
				expected = math.MaxInt64
			case 3:
				expected = 0
			}
			title, body := "Title", "body"
			if ops[i]&128 != 0 {
				title = " "
			}
			var want, got error
			switch op {
			case 0:
				fresh := Article{TenantID: tenant, ID: id, Category: "faq", Title: title, Body: body}
				if strings.TrimSpace(title) == "" {
					want = ErrInvalidArticle
				} else if exists {
					want = ErrDuplicate
				} else {
					a = fresh
					a.Version = 1
					a.State = StateDraft
					model[tenant][id] = a
				}
				got = s.Create(fresh)
			case 1, 2, 3:
				if !exists {
					want = ErrNotFound
				} else if a.Version != expected {
					want = ErrVersion
				} else if op == 1 && strings.TrimSpace(title) == "" {
					want = ErrInvalidArticle
				} else if op == 2 && a.State != StateDraft || op == 3 && a.State != StatePublished {
					want = ErrBadTransition
				} else if a.Version == math.MaxInt64 {
					want = ErrVersionExhausted
				} else {
					a.Version++
					if op == 1 {
						a.Title = title
						a.Body = body
					} else if op == 2 {
						a.State = StatePublished
					} else {
						a.State = StateArchived
					}
					model[tenant][id] = a
				}
				if op == 1 {
					got = s.Update(tenant, id, title, body, expected)
				} else if op == 2 {
					got = s.Publish(tenant, id, expected)
				} else {
					got = s.Archive(tenant, id, expected)
				}
			case 4:
				// Synthetic private-state boundary only; no restored/persisted source claimed.
				if exists {
					a.Version = math.MaxInt64
					model[tenant][id] = a
					for k, v := range s.articles {
						if v.TenantID == tenant && v.ID == id {
							v.Version = math.MaxInt64
							s.articles[k] = v
						}
					}
				}
			}
			if !errors.Is(got, want) {
				t.Fatal("model error", op, slot, got, want)
			}
			total := 0
			for mt, articles := range model {
				total += len(articles)
				expectedPublished := map[string]Article{}
				for mi, m := range articles {
					found := false
					for _, actual := range s.articles {
						if actual.TenantID == mt && actual.ID == mi {
							found = true
							if actual.Title != m.Title || actual.Body != m.Body || actual.Category != m.Category || actual.Version != m.Version || actual.State != m.State {
								t.Fatal("stored model mismatch", actual, m)
							}
						}
					}
					if !found {
						t.Fatal("model article missing")
					}
					if m.State == StatePublished {
						expectedPublished[mi] = m
					}
				}
				check := func(rows []Article) {
					if len(rows) != len(expectedPublished) {
						t.Fatal("published count", len(rows), len(expectedPublished))
					}
					seen := map[string]bool{}
					for _, row := range rows {
						m, ok := expectedPublished[row.ID]
						if !ok || seen[row.ID] || row.TenantID != mt || row.Version != m.Version || row.Title != m.Title {
							t.Fatal("published scope/model mismatch")
						}
						seen[row.ID] = true
					}
				}
				check(s.Published(mt, "faq"))
				check(s.Search(mt, " BODY "))
				if len(s.Published(mt, "missing")) != 0 {
					t.Fatal("category ignored")
				}
			}
			if len(s.articles) != total {
				t.Fatal("article cardinality mismatch")
			}
			_ = s.String()
		}
	})
}
