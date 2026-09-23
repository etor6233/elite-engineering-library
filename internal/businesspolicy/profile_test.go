package businesspolicy

// AUTHORED contract tests. No vendor attribution.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func profileHash(raw []byte) string { sum := sha256.Sum256(raw); return hex.EncodeToString(sum[:]) }
func changedProfile(t *testing.T, old, new string) *Profile {
	t.Helper()
	raw := bytes.ReplaceAll(ReferenceJSON(), []byte(old), []byte(new))
	p, err := Load(raw, profileHash(raw))
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestPolicyProfileBehaviorAndBindings(t *testing.T) {
	base := Reference()
	now := time.Date(2026, 9, 11, 12, 0, 0, 0, time.UTC)
	if base.SHA256() != ReferenceSHA256 || base.LeadTime() != 30*time.Minute || base.MaximumSlotDuration() != 8*time.Hour || base.MinimumSlotCapacity() != 1 || base.MaximumSlotCapacity() != 100 {
		t.Fatal("reference drift")
	}
	near := now.Add(time.Minute)
	if base.AllowsSlot(near, near.Add(time.Hour), now, 2) {
		t.Fatal("reference lead time bypass")
	}
	short := changedProfile(t, `"lead_time_seconds": 1800`, `"lead_time_seconds": 0`)
	if !short.AllowsSlot(near, near.Add(time.Hour), now, 2) {
		t.Fatal("compatible config did not change admission")
	}
	capTwo := changedProfile(t, `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 2`)
	start := now.Add(time.Hour)
	if !capTwo.AllowsSlot(start, start.Add(time.Hour), now, 2) || capTwo.AllowsSlot(start, start.Add(time.Hour), now, 3) {
		t.Fatal("capacity configuration ignored")
	}
	oneHour := changedProfile(t, `"maximum_slot_seconds": 28800`, `"maximum_slot_seconds": 3600`)
	if !oneHour.AllowsSlot(start, start.Add(time.Hour), now, 1) || oneHour.AllowsSlot(start, start.Add(time.Hour+time.Nanosecond), now, 1) {
		t.Fatal("duration boundary")
	}
	request := strings.Repeat("a", 64)
	historical, err := base.BindRequestHash(request)
	if err != nil || historical != request {
		t.Fatal("historical identity changed")
	}
	first, err := short.BindRequestHash(request)
	again, _ := short.BindRequestHash(request)
	other, _ := capTwo.BindRequestHash(request)
	if err != nil || first != again || first == historical || first == other {
		t.Fatal("policy request binding")
	}
	if _, err = short.BindRequestHash("invalid"); !errors.Is(err, ErrProfile) {
		t.Fatal("unbound request hash")
	}
	raw := ReferenceJSON()
	p, err := Load(raw, ReferenceSHA256)
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = 'x'
	if !p.Valid() || ReferenceJSON()[0] != '{' {
		t.Fatal("caller mutated immutable profile")
	}
	path := filepath.Join(t.TempDir(), "policy.json")
	if err = os.WriteFile(path, ReferenceJSON(), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = LoadFile(path, ReferenceSHA256); err != nil {
		t.Fatal(err)
	}
	if _, err = LoadFile(path, first); !errors.Is(err, ErrProfile) {
		t.Fatal("file hash mismatch")
	}
}
func TestPolicyProfileRejectsInvalidOrIncompatible(t *testing.T) {
	ref := string(ReferenceJSON())
	cases := map[string]string{
		"unknown":             strings.Replace(ref, `"revision": 1`, `"revision": 1, "extra": true`, 1),
		"duplicate":           strings.Replace(ref, `"revision": 1`, `"revision": 1, "revision": 2`, 1),
		"case-alias":          strings.Replace(ref, `"revision": 1`, `"revision": 1, "Revision": 2`, 1),
		"missing":             strings.Replace(ref, `"lead_time_seconds": 1800,`, "", 1),
		"null":                strings.Replace(ref, `"lead_time_seconds": 1800`, `"lead_time_seconds": null`, 1),
		"negative":            strings.Replace(ref, `"lead_time_seconds": 1800`, `"lead_time_seconds": -1`, 1),
		"overflow":            strings.Replace(ref, `"lead_time_seconds": 1800`, `"lead_time_seconds": 9223372037`, 1),
		"duration":            strings.Replace(ref, `"maximum_slot_seconds": 28800`, `"maximum_slot_seconds": 28801`, 1),
		"capacity":            strings.Replace(ref, `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 101`, 1),
		"minmax":              strings.Replace(strings.Replace(ref, `"minimum_slot_capacity": 1`, `"minimum_slot_capacity": 3`, 1), `"maximum_slot_capacity": 100`, `"maximum_slot_capacity": 2`, 1),
		"occupancy":           strings.Replace(ref, `"confirmed"`, `"completed"`, 1),
		"occupancy-duplicate": strings.Replace(ref, `"confirmed"`, `"requested"`, 1),
		"working":             strings.Replace(ref, `"contains_slot"`, `"ignore"`, 1),
		"unavailable":         strings.Replace(ref, `"reject_overlap"`, `"ignore"`, 1),
		"exclusive":           strings.Replace(ref, `"single_active_per_market_currency"`, `"best_price"`, 1),
		"halfopen":            strings.Replace(ref, `"half_open"`, `"inclusive"`, 1),
		"trailing":            ref + "{}", "oversize": strings.Repeat(" ", MaximumProfileBytes+1), "utf8": ref + string([]byte{0xff}),
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			raw := []byte(value)
			if _, err := Load(raw, profileHash(raw)); !errors.Is(err, ErrProfile) {
				t.Fatal("invalid profile accepted", err)
			}
		})
	}
	for _, hash := range []string{"", strings.Repeat("0", 64), strings.ToUpper(ReferenceSHA256)} {
		if _, err := Load(ReferenceJSON(), hash); !errors.Is(err, ErrProfile) {
			t.Fatal("unverified SHA accepted")
		}
	}
	if (*Profile)(nil).Valid() || new(Profile).Valid() {
		t.Fatal("zero profile accepted")
	}
}
func FuzzPolicyProfile(f *testing.F) {
	f.Add(ReferenceJSON())
	f.Add([]byte(`{"schema":"elite-business-policy/v1","revision":1,"revision":2}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, raw []byte) {
		p, err := Load(raw, profileHash(raw))
		if err == nil {
			if !p.Valid() || p.SHA256() != profileHash(raw) {
				t.Fatal("unbound accepted profile")
			}
			replay, e := Load(raw, p.SHA256())
			if e != nil || replay.SHA256() != p.SHA256() {
				t.Fatal("nondeterministic loader")
			}
		}
	})
}
