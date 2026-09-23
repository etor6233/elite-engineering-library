// Package agenttools binds the conversational agent's intents to a
// business-agnostic Domain contract (appointments, quotes, order status,
// returns). The chatbot never touches money/inventory directly: it delegates
// to the domain, which any franchise implements.
package agenttools

import "context"

// Domain is the business-agnostic capability contract the chatbot delegates
// to. Any franchise implements these four operations.
type Domain interface {
	BookAppointment(ctx context.Context, tenantID string, in AppointmentInput) (string, error)
	CreateQuote(ctx context.Context, tenantID string, in QuoteInput) (string, error)
	OrderStatus(ctx context.Context, tenantID, orderID string) (string, error)
	RequestReturn(ctx context.Context, tenantID string, in ReturnInput) (string, error)
}

// AppointmentInput carries a booking request.
type AppointmentInput struct {
	Service string
	When    string
}

// QuoteInput carries a quote request.
type QuoteInput struct {
	Product  string
	Quantity int
}

// ReturnInput carries a return request.
type ReturnInput struct {
	OrderID string
	Reason  string
}
