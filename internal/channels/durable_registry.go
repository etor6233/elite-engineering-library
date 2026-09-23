package channels

// DurableDeliveryChannel marks an adapter whose Send path is guarded by a
// durable attempt/receipt store. It is checked by app.NewProduction.
type DurableDeliveryChannel interface {
	Channel
	DurableDelivery() bool
}

// AllDurable reports whether every registered channel proves the durable
// marker. An empty registry is never production-ready.
func (r *Registry) AllDurable() bool {
	if r == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.channels) == 0 {
		return false
	}
	for _, channel := range r.channels {
		durable, ok := channel.(DurableDeliveryChannel)
		if !ok || !durable.DurableDelivery() {
			return false
		}
	}
	return true
}
