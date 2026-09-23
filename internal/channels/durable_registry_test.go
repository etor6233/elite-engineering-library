package channels

import "testing"

type durableFake struct{ *fakeChannel }

func (d durableFake) DurableDelivery() bool { return d.fakeChannel != nil }

func TestRegistryRequiresEveryProductionChannelDurable(t *testing.T) {
	r := NewRegistry()
	if r.AllDurable() {
		t.Fatal("empty registry accepted")
	}
	if err := r.Register(&fakeChannel{code: "email"}); err != nil {
		t.Fatal(err)
	}
	if r.AllDurable() {
		t.Fatal("plain channel accepted")
	}
	r = NewRegistry()
	if err := r.Register(durableFake{fakeChannel: &fakeChannel{code: "whatsapp"}}); err != nil {
		t.Fatal(err)
	}
	if !r.AllDurable() {
		t.Fatal("durable channel rejected")
	}
}
