package whatsappbridge

import "testing"

func TestCampaignWorkspaceQuery(t *testing.T) {
	for _, q := range []string{"kind=audience&kind=campaigns", "after=ok&after=other", "unknown=x", "after=bad%2Fcursor", "kind=wrong", "kind=campaigns&after=%zz"} {
		if _, _, e := campaignWorkspaceQuery(q); e == nil {
			t.Fatalf("accepted %q", q)
		}
	}
	for _, q := range []string{"", "kind=campaigns", "kind=audience&after=lead-25"} {
		if _, _, e := campaignWorkspaceQuery(q); e != nil {
			t.Fatalf("rejected %q", q)
		}
	}
}
