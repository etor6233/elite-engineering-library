package agent

import "strings"

// Intent is a bounded user goal. The set is closed so the agent can authorize
// exactly one tool per intent and hand off anything unknown.
type Intent string

const (
	IntentGreeting      Intent = "greeting"
	IntentHandoff       Intent = "handoff"
	IntentAppointment   Intent = "appointment"
	IntentQuote         Intent = "quote"
	IntentOrderStatus   Intent = "order_status"
	IntentReturnRequest Intent = "return_request"
	IntentFAQ           Intent = "faq"
	IntentFallback      Intent = "fallback"
)

// Router classifies a user message into a bounded intent.
type Router interface {
	Route(text string) Intent
}

// KeywordRouter is a deterministic baseline router. It is not a language
// model; an LLM-based router can replace it through the same interface when a
// provider is admitted.
type KeywordRouter struct{}

// Route applies precedence: handoff > appointment > quote > return_request >
// order_status > greeting > faq > fallback.
func (KeywordRouter) Route(text string) Intent {
	t := strings.ToLower(strings.TrimSpace(text))

	containsAny := func(words ...string) bool {
		for _, w := range words {
			if strings.Contains(t, w) {
				return true
			}
		}
		return false
	}

	switch {
	case containsAny("hablar con una persona", "agente humano", "asesor", "representante", "con un humano"):
		return IntentHandoff
	case containsAny("turno", "cita", "reserva", "agendar", "reservar", "disponibilidad"):
		return IntentAppointment
	case containsAny("cotizar", "cotizacion", "presupuesto", "precio", "cuanto cuesta", "cuanto sale", "comprar", "venta"):
		return IntentQuote
	case containsAny("devolver", "devolucion", "reembolso", "cancelar", "cambio"):
		return IntentReturnRequest
	case containsAny("estado de", "seguimiento", "tracking", "donde esta mi", "mi pedido"):
		return IntentOrderStatus
	case containsAny("hola", "buenos dias", "buenas tardes", "buenas noches", "buen dia"):
		return IntentGreeting
	case strings.Contains(t, "?"), containsAny("que es", "como ", "cual ", "cuando ", "donde ", "por que", "que necesito"):
		return IntentFAQ
	default:
		return IntentFallback
	}
}
