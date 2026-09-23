package httpapi
import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "elite.local/enterprise/internal/inventorycontrol"
 "elite.local/enterprise/internal/platform/identity"
)
type containmentBulkRepo struct {bulkHTTPRepo; reservations int}
func(r *containmentBulkRepo) ReserveBulk(_ context.Context,_,_ string,v inventorycontrol.BulkReservation)(inventorycontrol.BulkReservation,error){r.reservations++;return v,nil}
func TestWarehouseHostContainsLegacyReservationOnly(t *testing.T){
 repo:=&containmentBulkRepo{};mux:=http.NewServeMux()
 p:=identity.Principal{Subject:"operator",TenantID:"018f4d4a-7b36-7a21-8d10-2f4c54c29d00",Permissions:map[string]struct{}{"inventory:write":{},"inventory:read":{}},Organizations:map[string]struct{}{"a":{}}}
 InventoryControlModule{Bulk:inventorycontrol.NewBulkService(repo,inventoryIDs{}),RequireDurableReservations:true}.Register(mux,royaltyVerifier{principal:p})
 call:=func(method,path,body string,auth bool)*httptest.ResponseRecorder{r:=httptest.NewRequest(method,path,strings.NewReader(body));r.Header.Set("Content-Type","application/json");if auth{r.Header.Set("Authorization","Bearer fixture")};w:=httptest.NewRecorder();mux.ServeHTTP(w,r);return w}
 for i:=0;i<2;i++{w:=call("POST","/v1/inventory/bulk/reservations",`{"organization_id":"a","bin_id":"bin","item_id":"i","demand_kind":"manual","demand_id":"d","demand_line_id":"l","quantity":"1"}`,true);if w.Code!=409||!strings.Contains(w.Body.String(),"DURABLE_RESERVATION_REQUIRED"){t.Fatal(w.Code,w.Body.String())}}
 if repo.reservations!=0{t.Fatal("legacy side effect")}
 if w:=call("POST","/v1/inventory/bulk/reservations",`{"organization_id":"a"}`,false);w.Code!=401{t.Fatal("authentication",w.Code)}
 if w:=call("POST","/v1/inventory/bulk/reservations",`{"organization_id":"other"}`,true);w.Code!=403{t.Fatal("organization",w.Code)}
 if w:=call("GET","/v1/inventory/bulk/availability?organization_id=a&item_id=i","",true);w.Code!=200{t.Fatal("read changed",w.Code,w.Body.String())}
}
