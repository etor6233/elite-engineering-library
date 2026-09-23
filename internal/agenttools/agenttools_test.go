package agenttools

import (
	"context"
	"errors"
	"testing"

	"elite.local/enterprise/internal/agent"
)

type fakeDomain struct {
	appt  AppointmentInput
	quote QuoteInput
	order string
	ret   ReturnInput
}

func (f *fakeDomain) BookAppointment(_ context.Context, _ string, in AppointmentInput) (string, error) {
	f.appt = in
	return "turno confirmado", nil
}
func (f *fakeDomain) CreateQuote(_ context.Context, _ string, in QuoteInput) (string, error) {
	f.quote = in
	return "cotización lista", nil
}
func (f *fakeDomain) OrderStatus(_ context.Context, _ string, orderID string) (string, error) {
	f.order = orderID
	return "en camino", nil
}
func (f *fakeDomain) RequestReturn(_ context.Context, _ string, in ReturnInput) (string, error) {
	f.ret = in
	return "devolución solicitada", nil
}

func TestToolkitWiresFourIntents(t *testing.T) {
	r := Toolkit(&fakeDomain{})
	for _, intent := range []agent.Intent{
		agent.IntentAppointment, agent.IntentQuote, agent.IntentOrderStatus, agent.IntentReturnRequest,
	} {
		if _, ok := r.For(intent); !ok {
			t.Fatalf("missing tool for intent %q", intent)
		}
	}
}

func TestAppointmentToolExtractsAndCalls(t *testing.T) {
	d := &fakeDomain{}
	tool, _ := Toolkit(d).For(agent.IntentAppointment)
	out, err := tool.Run(context.Background(), "t", "servicio: corte, cuando: mañana")
	if err != nil || out != "turno confirmado" {
		t.Fatalf("unexpected: %q/%v", out, err)
	}
	if d.appt.Service != "corte" || d.appt.When != "mañana" {
		t.Fatalf("wrong args: %+v", d.appt)
	}
}

func TestAppointmentToolNeedsInfo(t *testing.T) {
	tool, _ := Toolkit(&fakeDomain{}).For(agent.IntentAppointment)
	if _, err := tool.Run(context.Background(), "t", "quiero un turno"); !errors.Is(err, agent.ErrNeedsInfo) {
		t.Fatalf("expected ErrNeedsInfo, got %v", err)
	}
}

func TestQuoteToolQuantity(t *testing.T) {
	d := &fakeDomain{}
	tool, _ := Toolkit(d).For(agent.IntentQuote)
	if _, err := tool.Run(context.Background(), "t", "producto: scooter, cantidad: 2"); err != nil {
		t.Fatal(err)
	}
	if d.quote.Product != "scooter" || d.quote.Quantity != 2 {
		t.Fatalf("wrong quote args: %+v", d.quote)
	}
}

func TestQuoteToolInvalidQuantity(t *testing.T) {
	tool, _ := Toolkit(&fakeDomain{}).For(agent.IntentQuote)
	if _, err := tool.Run(context.Background(), "t", "producto: x, cantidad: cero"); !errors.Is(err, agent.ErrNeedsInfo) {
		t.Fatalf("expected ErrNeedsInfo for bad quantity, got %v", err)
	}
}

func TestOrderStatusAndReturn(t *testing.T) {
	d := &fakeDomain{}
	tool, _ := Toolkit(d).For(agent.IntentOrderStatus)
	if _, err := tool.Run(context.Background(), "t", "pedido: ORD-9"); err != nil {
		t.Fatal(err)
	}
	if d.order != "ORD-9" {
		t.Fatalf("wrong order id: %q", d.order)
	}

	ret, _ := Toolkit(d).For(agent.IntentReturnRequest)
	if _, err := ret.Run(context.Background(), "t", "pedido: ORD-9, motivo: no sirve"); err != nil {
		t.Fatal(err)
	}
	if d.ret.OrderID != "ORD-9" || d.ret.Reason != "no sirve" {
		t.Fatalf("wrong return args: %+v", d.ret)
	}
}
