package conversationruntime

import (
	"context"
	"elite.local/enterprise/internal/llmopenai"
	"errors"
	"fmt"
	"testing"
)

func TestU8InvalidArgumentsHandoffAndTransportRetry(t *testing.T) {
	for _, tc := range []struct {
		name    string
		err     error
		handoff bool
	}{{"invalid", llmopenai.ErrInvalidToolArguments, true}, {"wrapped", fmt.Errorf("adapter: %w", llmopenai.ErrInvalidToolArguments), true}, {"transport", errors.New("transport interrupted"), false}, {"same_text_not_type", errors.New(llmopenai.ErrInvalidToolArguments.Error()), false}} {
		t.Run(tc.name, func(t *testing.T) {
			store, domain := &memoryStore{}, &fakeDomain{}
			model := &fakeModel{startErr: tc.err}
			r := runtimeFor(store, model, domain)
			message := inbound("u8-invalid", "Cotizar")
			result, e := r.Handle(context.Background(), message)
			if domain.quotes != 0 || domain.appointments != 0 || domain.statuses != 0 || model.completes != 0 {
				t.Fatal("invalid input reached effects")
			}
			if tc.handoff {
				if e != nil || store.completion == nil || store.completion.State != StateHandedOff || store.completion.FailureCode != "INVALID_TOOL_ARGUMENTS" || store.retryCode != "" {
					t.Fatal("missing terminal handoff", e)
				}
				again, e := r.Handle(context.Background(), message)
				if e != nil || again != result || model.starts != 1 {
					t.Fatal("handoff replay called model")
				}
			} else if e == nil || store.completion != nil || store.retryCode != "MODEL_START_FAILED" {
				t.Fatal("transport was not retryable")
			}
		})
	}
}
