package aifoundation

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeProvider struct {
	resp CompletionResponse
	err  error
}

func (f fakeProvider) Generate(context.Context, CompletionRequest) (CompletionResponse, error) {
	return f.resp, f.err
}

func TestEvalSuiteGateDecision(t *testing.T) {
	p := fakeProvider{}
	suite := EvalSuite{
		ID:          "v1",
		MinPassRate: 1.0,
		Cases: []EvalCase{
			{ID: "pass", Assert: func(CompletionResponse) bool { return true }},
			{ID: "fail", Assert: func(CompletionResponse) bool { return false }},
		},
	}
	res, err := suite.Run(p)
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 2 || res.Passed != 1 {
		t.Fatalf("unexpected result: %+v", res)
	}
	if res.GatePassed {
		t.Fatal("gate passed at 50%% with MinPassRate 1.0")
	}

	suite.MinPassRate = 0.5
	res, err = suite.Run(p)
	if err != nil {
		t.Fatal(err)
	}
	if !res.GatePassed {
		t.Fatal("gate should pass at 50%% with MinPassRate 0.5")
	}
}

func TestEvalSuiteRejectsEmpty(t *testing.T) {
	suite := EvalSuite{ID: "empty", MinPassRate: 1.0}
	if _, err := suite.Run(fakeProvider{}); err == nil {
		t.Fatal("empty suite accepted")
	}
}

func TestEvalSuitePassesThroughErrors(t *testing.T) {
	p := fakeProvider{err: context.Canceled}
	suite := EvalSuite{
		ID:          "err",
		MinPassRate: 1.0,
		Cases:       []EvalCase{{ID: "c", Assert: func(CompletionResponse) bool { return true }}},
	}
	if _, err := suite.Run(p); err == nil {
		t.Fatal("provider error swallowed")
	}
}

func TestEvalSuiteJSONAssert(t *testing.T) {
	p := fakeProvider{resp: CompletionResponse{Content: json.RawMessage(`{"answer":"ok"}`)}}
	suite := EvalSuite{
		ID:          "json",
		MinPassRate: 1.0,
		Cases: []EvalCase{{
			ID: "schema",
			Assert: func(r CompletionResponse) bool {
				return ValidateStructuredOutput([]byte(closedSchemaJSON), r.Content) == nil
			},
		}},
	}
	res, err := suite.Run(p)
	if err != nil {
		t.Fatal(err)
	}
	if !res.GatePassed {
		t.Fatal("schema assertion failed")
	}
}
