package serialsupply

// AUTHORED finite-fuzz invariants for physical serial identity, receipt hashes
// and immutable command versions; not a deployment security certification.
import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func FuzzSerialSupplyCommandIdentity(f *testing.F) {
	for _, s := range []string{"FRAME-0001", "BATERÍA-ñ-2", " serial ", "a\x00b", "\xff", "電池-7", strings.Repeat("A", 129)} {
		f.Add(s, int64(9007199254740993))
	}
	f.Fuzz(func(t *testing.T, serial string, version int64) {
		r := Command{PurchaseOrderID: "po", CommandID: "register", Kind: "register", LineID: "line", SerialNumber: serial, ExpectedVersion: version, EvidenceSHA256: strings.Repeat("a", 64)}
		if !r.Valid() {
			return
		}
		raw, hash, err := Canonical(r)
		if err != nil {
			t.Fatal(err)
		}
		var decoded Command
		if err = json.Unmarshal(raw, &decoded); err != nil || decoded.SerialNumber != serial || decoded.ExpectedVersion != version || !decoded.Valid() {
			t.Fatal("serial/version identity lost", err)
		}
		again, other, err := Canonical(decoded)
		if err != nil || hash != other || !bytes.Equal(raw, again) {
			t.Fatal("unstable durable command")
		}
		decoded.CommandID += "-changed"
		_, different, err := Canonical(decoded)
		if err != nil || different == hash {
			t.Fatal("command identity not hash-bound")
		}
		shipment := Command{PurchaseOrderID: "po", CommandID: "receive", Kind: "receive", ShipmentID: "asn", Units: []string{"unit", "unit"}, ExpectedVersion: version, EvidenceSHA256: r.EvidenceSHA256}
		if shipment.Valid() {
			t.Fatal("duplicate physical unit admitted")
		}
	})
}
