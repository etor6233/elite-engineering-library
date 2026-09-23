package returnaccounting

import (
	"context"
	"elite.local/enterprise/internal/returneffects"
	"errors"
	"testing"
	"time"
)

type ids struct{ v string }

func (i ids) New() string { return i.v }

type store struct {
	work       *returneffects.Work
	err        error
	completion returneffects.Completion
}

func (s *store) Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error) {
	return s.work, nil
}
func (s *store) PostOrReconcile(context.Context, returneffects.Work, string) (Result, error) {
	return Result{RequestID: "a", ReversalJournalID: "r"}, s.err
}
func (s *store) Finish(_ context.Context, _ returneffects.Work, _ string, c returneffects.Completion) error {
	s.completion = c
	return nil
}
func TestProcessor(t *testing.T) {
	s := &store{work: &returneffects.Work{RequestID: "a"}}
	p, e := NewProcessor(s, ids{"c"}, "worker", time.Minute, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	r, e := p.ProcessOne(context.Background())
	if e != nil || r.ReversalJournalID != "r" {
		t.Fatalf("%+v %v", r, e)
	}
	s.err = ErrConflict
	if _, e = p.ProcessOne(context.Background()); !errors.Is(e, ErrConflict) || s.completion.Outcome != "blocked" {
		t.Fatalf("%v %+v", e, s.completion)
	}
}
func TestNoWorkAndConfig(t *testing.T) {
	if _, e := NewProcessor(&store{}, ids{"c"}, "bad worker", time.Minute, time.Second); e == nil {
		t.Fatal("invalid accepted")
	}
	p, e := NewProcessor(&store{}, ids{"c"}, "worker", time.Minute, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = p.ProcessOne(context.Background()); !errors.Is(e, ErrNoWork) {
		t.Fatal(e)
	}
}
