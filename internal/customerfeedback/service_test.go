package customerfeedback

import (
	"context"
	"testing"
)

type neverRepository struct{}

func (neverRepository) Definition(context.Context, Scope, string) (Definition, error) {
	panic("invalid request reached persistence")
}
func (neverRepository) Submit(context.Context, Scope, string, Submission) (Result, error) {
	panic("invalid request reached persistence")
}
func (neverRepository) OwnAnswer(context.Context, Scope, string) (Answer, error) {
	panic("invalid request reached persistence")
}
func (neverRepository) Summary(context.Context, Scope, string) (Summary, error) {
	panic("invalid request reached persistence")
}
func TestFeedbackNPSGoldenAndThreshold(t *testing.T) {
	for _, c := range []struct {
		n, p, d, m int64
		available  bool
		want       float64
	}{
		{0, 0, 0, 1, false, 0}, {1, 1, 0, 2, false, 0}, {5, 1, 3, 1, true, -40}, {3, 2, 1, 2, true, 100.0 / 3}, {2, 0, 2, 1, true, -100}, {2, 2, 0, 1, true, 100},
	} {
		got, e := NPS(c.n, c.p, c.d, c.m)
		if e != nil || got.Available != c.available || got.Responses != c.n {
			t.Fatalf("%+v %v", got, e)
		}
		if c.available {
			if got.NPS == nil || *got.NPS-c.want > 1e-12 || c.want-*got.NPS > 1e-12 {
				t.Fatal(got)
			}
		} else if got.NPS != nil {
			t.Fatal("score exposed below threshold")
		}
	}
	for _, c := range [][4]int64{{-1, 0, 0, 1}, {1, 2, 0, 1}, {1, 1, 1, 1}, {1, 0, 0, 0}, {1, -1, 0, 1}} {
		if _, e := NPS(c[0], c[1], c[2], c[3]); e != ErrInvalid {
			t.Fatal(c, e)
		}
	}
}
func TestFeedbackRejectsBeforePersistence(t *testing.T) {
	service := NewService(neverRepository{})
	scope := Scope{"tenant", "org", "customer"}
	for _, input := range []Submission{{-1, true, "v1"}, {11, true, "v1"}, {1, false, "v1"}, {1, true, ""}, {1, true, "v1\n"}} {
		if _, e := service.Submit(context.Background(), scope, "survey", input); e != ErrInvalid {
			t.Fatal(e)
		}
	}
	for _, s := range []Scope{{}, {"tenant", "org", ""}, {"tenant", "\x00", "customer"}} {
		if _, e := service.Definition(context.Background(), s, "survey"); e != ErrInvalid {
			t.Fatal(e)
		}
	}
}
