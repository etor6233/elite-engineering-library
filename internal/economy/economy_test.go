package economy

import (
	"errors"
	"testing"
	"time"
)

func TestEstimateTokens(t *testing.T) {
	if EstimateTokens("") != 0 {
		t.Fatal("empty text should be 0 tokens")
	}
	if EstimateTokens("hola") < 1 {
		t.Fatal("non-empty text should be >= 1 token")
	}
}

func TestSemanticCacheHitAndMiss(t *testing.T) {
	c := &SemanticCache{}
	c.StoreScoped(CacheScope{"t", "u", "p", "k"}, []float32{1, 0, 0}, "respuesta A", time.Minute)
	c.StoreScoped(CacheScope{"t", "u", "p", "k"}, []float32{0, 1, 0}, "respuesta B", time.Minute)

	if ans, ok := c.LookupScoped(CacheScope{"t", "u", "p", "k"}, []float32{1, 0, 0}, 0.9); !ok || ans != "respuesta A" {
		t.Fatalf("exact hit failed: %q ok=%v", ans, ok)
	}
	if ans, ok := c.LookupScoped(CacheScope{"t", "u", "p", "k"}, []float32{0.9, 0.1, 0}, 0.9); !ok || ans != "respuesta A" {
		t.Fatalf("near hit failed: %q ok=%v", ans, ok)
	}
	if _, ok := c.LookupScoped(CacheScope{"t", "u", "p", "k"}, []float32{0, 0, 1}, 0.9); ok {
		t.Fatal("orthogonal vector should miss")
	}
	if _, ok := c.LookupScoped(CacheScope{"t", "u", "p", "k"}, []float32{1, 0}, 0.9); ok {
		t.Fatal("dim mismatch should miss, not error")
	}
}

func TestCompressRespectsBudget(t *testing.T) {
	// budget just enough for system + user, no context
	sys := "system prompt"
	usr := "pregunta del usuario"
	budget := EstimateTokens(sys) + EstimateTokens(usr)
	p := Compress(sys, usr, []string{"contexto largo que no entra", "otro"}, budget)
	if len(p.Context) != 0 {
		t.Fatalf("expected no context within tight budget, got %v", p.Context)
	}
	if p.Tokens != budget {
		t.Fatalf("expected tokens=%d, got %d", budget, p.Tokens)
	}

	// larger budget admits context up to the limit
	big := budget + 1000
	p2 := Compress(sys, usr, []string{"a", "b"}, big)
	if len(p2.Context) != 2 {
		t.Fatalf("expected both context items, got %v", p2.Context)
	}
	if p2.Tokens > big {
		t.Fatalf("compressed prompt exceeded budget: %d > %d", p2.Tokens, big)
	}
}

func TestLedgerFailClosed(t *testing.T) {
	l := NewCostLedger()
	l.SetBudget("t", 100)
	if !l.CanSpend("t", 100) {
		t.Fatal("exact budget should be spendable")
	}
	if l.CanSpend("t", 101) {
		t.Fatal("over-budget should not be spendable")
	}
	if err := l.Spend("t", 40); err != nil {
		t.Fatal(err)
	}
	if l.Remaining("t") != 60 {
		t.Fatalf("remaining should be 60, got %d", l.Remaining("t"))
	}
	if err := l.Spend("t", 61); !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("expected ErrBudgetExceeded, got %v", err)
	}
}

func TestGovernorDecision(t *testing.T) {
	g := &Governor{
		Cache:     &SemanticCache{},
		Ledger:    NewCostLedger(),
		Threshold: 0.9,
	}
	g.Cache.StoreScoped(CacheScope{"t", "u", "p", "k"}, []float32{1, 0}, "cacheada", time.Minute)
	g.Ledger.SetBudget("t", 100)

	// cache hit → PathCache
	if p, ans := g.DecideScoped(CacheScope{"t", "u", "p", "k"}, []float32{1, 0}, 1000); p != PathCache || ans != "cacheada" {
		t.Fatalf("expected cache hit, got path=%d ans=%q", p, ans)
	}
	// miss but within budget → PathLLM
	if p, _ := g.DecideScoped(CacheScope{"t", "u", "p", "k"}, []float32{0, 1}, 50); p != PathLLM {
		t.Fatalf("expected PathLLM, got %d", p)
	}
	// miss and over budget → PathBlocked
	if p, _ := g.DecideScoped(CacheScope{"t", "u", "p", "k"}, []float32{0, 1}, 200); p != PathBlocked {
		t.Fatalf("expected PathBlocked, got %d", p)
	}
}
