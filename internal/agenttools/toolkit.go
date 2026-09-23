package agenttools

import "elite.local/enterprise/internal/agent"

// Toolkit wires the four business tools into a ToolRegistry (one per intent).
func Toolkit(d Domain) *agent.ToolRegistry {
	r := agent.NewToolRegistry()
	_ = r.Register(appointmentTool{d: d})
	_ = r.Register(quoteTool{d: d})
	_ = r.Register(orderStatusTool{d: d})
	_ = r.Register(returnTool{d: d})
	return r
}
