package agent

import "testing"

func TestKeywordRouterIntents(t *testing.T) {
	var r KeywordRouter
	cases := map[string]Intent{
		"hola, buenos dias":                    IntentGreeting,
		"quiero hablar con una persona":        IntentHandoff,
		"quiero agendar un turno":              IntentAppointment,
		"necesito reservar una cita":           IntentAppointment,
		"cuanto cuesta el scooter?":            IntentQuote,
		"quiero cotizar un vehiculo":           IntentQuote,
		"donde esta mi pedido?":                IntentOrderStatus,
		"quiero devolver mi compra":            IntentReturnRequest,
		"que es el plan de mantenimiento?":     IntentFAQ,
		"asdfghjkl":                            IntentFallback,
	}
	for text, want := range cases {
		if got := r.Route(text); got != want {
			t.Errorf("Route(%q) = %q, want %q", text, got, want)
		}
	}
}

func TestKeywordRouterPrecedence(t *testing.T) {
	var r KeywordRouter
	// "cita" (appointment) but also "con un humano" (handoff) → handoff wins.
	if got := r.Route("quiero una cita pero prefiero hablar con un humano"); got != IntentHandoff {
		t.Fatalf("handoff should win precedence, got %q", got)
	}
}
