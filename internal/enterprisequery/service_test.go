package enterprisequery

import (
	"context"
	"testing"
)

type fakeRepository struct{ calls int }

func (f *fakeRepository) Overview(context.Context, string, string) (Overview, error) {
	f.calls++
	return Overview{}, nil
}
func (f *fakeRepository) Orders(context.Context, string, string, string, int, string) (Page[Order], error) {
	f.calls++
	return Page[Order]{}, nil
}
func (f *fakeRepository) FactoryUnits(context.Context, string, string, int, string) (Page[FactoryUnit], error) {
	f.calls++
	return Page[FactoryUnit]{}, nil
}
func (f *fakeRepository) ServiceCases(context.Context, string, string, string, int, string) (Page[ServiceCase], error) {
	f.calls++
	return Page[ServiceCase]{}, nil
}

func TestServiceRejectsUnboundedAndMissingScope(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)
	if _, err := service.Orders(context.Background(), "tenant", "org", "", 101, ""); err == nil {
		t.Fatal("unbounded page accepted")
	}
	if _, err := service.Overview(context.Background(), "tenant", ""); err == nil {
		t.Fatal("missing organization accepted")
	}
	if repository.calls != 0 {
		t.Fatal("invalid query reached repository")
	}
}
