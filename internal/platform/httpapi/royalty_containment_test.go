package httpapi

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "elite.local/enterprise/internal/platform/identity"
 "elite.local/enterprise/internal/royalty"
)

type royaltyContainmentRepo struct { royaltyRepo; opens int }
func (r *royaltyContainmentRepo) OpenSettlement(context.Context,string,string,royalty.Settlement) error { r.opens++;return nil }

func TestRoyaltyDurableHostContainsOnlyReceiptlessCreation(t *testing.T) {
 repo:=&royaltyContainmentRepo{}
 principal:=identity.Principal{Subject:"controller",TenantID:"018f4d4a-7b36-7a21-8d10-2f4c54c29d00",Permissions:map[string]struct{}{"royalty:policy":{},"royalty:settle":{},"royalty:post":{}},Organizations:map[string]struct{}{"franchise":{}}}
 mux:=http.NewServeMux()
 RoyaltyModule{Service:royalty.NewService(repo,&royaltyIDs{}),RequireDurableCommands:true}.Register(mux,royaltyVerifier{principal:principal})
 request:=func(path,body string,authorized bool)*httptest.ResponseRecorder {
  r:=httptest.NewRequest("POST",path,strings.NewReader(body));r.Header.Set("Content-Type","application/json");r.Header.Set("Idempotency-Key","event-1")
  if authorized {r.Header.Set("Authorization","Bearer fixture")};w:=httptest.NewRecorder();mux.ServeHTTP(w,r);return w
 }
 for _,path:=range []string{"/v1/royalties/policies","/v1/royalties/settlements"} {
  for i:=0;i<2;i++ { w:=request(path,`{"organization_id":"franchise"}`,true);if w.Code!=409 || !strings.Contains(w.Body.String(),"DURABLE_COMMAND_REQUIRED") || !strings.Contains(w.Body.String(),"/v1/finance/commands") { t.Fatalf("uncontained %s: %d %s",path,w.Code,w.Body.String()) } }
  if w:=request(path,`{"organization_id":"franchise"}`,false);w.Code!=401 {t.Fatalf("auth %s: %d",path,w.Code)}
  if w:=request(path,`{"organization_id":"other"}`,true);w.Code!=403 {t.Fatalf("scope %s: %d",path,w.Code)}
 }
 if repo.policies!=0 || repo.opens!=0 {t.Fatalf("legacy side effects: policies=%d opens=%d",repo.policies,repo.opens)}
 w:=request("/v1/royalties/accruals/from-payment",`{"organization_id":"franchise","payment_attempt_id":"payment","payment_expected_version":3,"state":"captured","occurred_at":"2026-01-01T00:00:00Z"}`,true)
 if w.Code!=201 || repo.accruals!=1 {t.Fatalf("unrelated operation changed: %d %s",w.Code,w.Body.String())}
}
