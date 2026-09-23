package rolemetrics

import (
	"elite.local/enterprise/internal/platform/identity"
	"net/url"
	"strings"
	"testing"
)

func FuzzMetricBoundary(f *testing.F) {
	for _, v := range []struct{ k, o string }{{"orders", "store"}, {"own-orders", "a:b"}, {"../orders", "store"}, {"orders", "x&tenant=other"}, {"", ""}, {"factory-owned", "factory"}} {
		f.Add(v.k, v.o)
	}
	f.Fuzz(func(t *testing.T, kind, org string) {
		root := identity.Principal{TenantID: "tenant", Subject: "reader", Permissions: map[string]struct{}{"*": {}}}
		if Authorized(root, kind, org) && (Permission(kind) == "" || !ID(org)) {
			t.Fatal("unknown kind or unsafe scope authorized")
		}
		if ID(org) {
			if len(org) > 128 || url.QueryEscape(org) == "" || strings.ContainsAny(org, "/&?=#% \n\r\x00") {
				t.Fatal("identifier changes query shape")
			}
			v, e := url.ParseQuery("organization_id=" + url.QueryEscape(org))
			if e != nil || len(v) != 1 || v.Get("organization_id") != org {
				t.Fatal("scope does not round trip")
			}
		}
		for _, p := range []identity.Principal{{TenantID: "tenant", Subject: "reader"}, {TenantID: "tenant", Permissions: root.Permissions}, {Subject: "reader", Permissions: root.Permissions}} {
			if Authorized(p, kind, org) {
				t.Fatal("missing identity or permission authorized")
			}
		}
	})
}
