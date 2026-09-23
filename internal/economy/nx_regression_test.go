package economy

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNXOverBudgetPromptNotReturned(t *testing.T) {
	p := Compress("must preserve mandatory", "message", nil, 1)
	if p.System != "" || p.User != "" || p.Tokens > 1 {
		t.Fatalf("returned runnable overbudget prompt: %+v", p)
	}
}
func TestNXGovernorConsumesAllowanceAtomically(t *testing.T) {
	l := NewCostLedger()
	l.SetBudget("t", 10)
	g := Governor{Ledger: l}
	var n atomic.Int64
	var w sync.WaitGroup
	for i := 0; i < 20; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			if p, _ := g.Decide("t", nil, 10); p == PathLLM {
				n.Add(1)
			}
		}()
	}
	w.Wait()
	if n.Load() != 1 {
		t.Fatalf("approved %d calls against one allowance", n.Load())
	}
}
func TestNXMissingTenantCannotUseCachedAnswer(t *testing.T) {
	c := &SemanticCache{}
	c.Store([]float32{1, 0}, "sensitive")
	g := Governor{Cache: c, Threshold: .9}
	if p, _ := g.Decide("", []float32{1, 0}, 10); p == PathCache {
		t.Fatal("missing session tenant hit")
	}
}
