package agenttools

import (
	"context"
	"strconv"

	"elite.local/enterprise/internal/agent"
)

type appointmentTool struct{ d Domain }

func (appointmentTool) Name() string         { return "book-appointment" }
func (appointmentTool) Intent() agent.Intent { return agent.IntentAppointment }
func (t appointmentTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	svc, hasSvc := kv["servicio"]
	when, hasWhen := kv["cuando"]
	if !hasSvc || !hasWhen {
		return "Necesito servicio y cuándo (ej. servicio: corte, cuando: mañana).", agent.ErrNeedsInfo
	}
	return t.d.BookAppointment(ctx, tenantID, AppointmentInput{Service: svc, When: when})
}

type quoteTool struct{ d Domain }

func (quoteTool) Name() string         { return "create-quote" }
func (quoteTool) Intent() agent.Intent { return agent.IntentQuote }
func (t quoteTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	product, has := kv["producto"]
	if !has {
		return "Necesito el producto (ej. producto: scooter).", agent.ErrNeedsInfo
	}
	qty := 1
	if q, ok := kv["cantidad"]; ok {
		n, err := strconv.Atoi(q)
		if err != nil || n <= 0 {
			return "Cantidad inválida.", agent.ErrNeedsInfo
		}
		qty = n
	}
	return t.d.CreateQuote(ctx, tenantID, QuoteInput{Product: product, Quantity: qty})
}

type orderStatusTool struct{ d Domain }

func (orderStatusTool) Name() string         { return "order-status" }
func (orderStatusTool) Intent() agent.Intent { return agent.IntentOrderStatus }
func (t orderStatusTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	order, has := kv["pedido"]
	if !has {
		return "Necesito el número de pedido (ej. pedido: ORD-123).", agent.ErrNeedsInfo
	}
	return t.d.OrderStatus(ctx, tenantID, order)
}

type returnTool struct{ d Domain }

func (returnTool) Name() string         { return "request-return" }
func (returnTool) Intent() agent.Intent { return agent.IntentReturnRequest }
func (t returnTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	order, hasOrder := kv["pedido"]
	reason, hasReason := kv["motivo"]
	if !hasOrder || !hasReason {
		return "Necesito pedido y motivo (ej. pedido: ORD-1, motivo: no sirve).", agent.ErrNeedsInfo
	}
	return t.d.RequestReturn(ctx, tenantID, ReturnInput{OrderID: order, Reason: reason})
}
