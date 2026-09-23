package httpapi
import(
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
 cr "elite.local/enterprise/internal/catalogrelease"
)
type catalogCacheProbe struct{CatalogReleaseService}
func(catalogCacheProbe)Public(context.Context)(cr.Publication,error){
 raw,_:=json.Marshal(cr.Snapshot{Profile:cr.Profile{Market:"AR",Currency:"ARS"}})
 return cr.Publication{Generation:1,SHA256:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",Snapshot:raw},nil
}
func TestCatalogPublicCachePolicy(t *testing.T){
 mux:=http.NewServeMux();CatalogReleaseModule{Service:catalogCacheProbe{}}.Register(mux,nil)
 first:=httptest.NewRecorder();mux.ServeHTTP(first,httptest.NewRequest("GET","/v1/public/catalog",nil))
 if first.Code!=200||first.Header().Get("Cache-Control")!="public, max-age=0, must-revalidate"{t.Fatal(first.Code,first.Header())}
 etag:=first.Header().Get("ETag");if etag==""{t.Fatal("no ETag")}
 req:=httptest.NewRequest("GET","/v1/public/catalog",nil);req.Header.Set("If-None-Match",etag)
 next:=httptest.NewRecorder();mux.ServeHTTP(next,req)
 if next.Code!=304||next.Body.Len()!=0||next.Header().Get("Cache-Control")!="public, max-age=0, must-revalidate"{t.Fatal(next.Code,next.Header())}
}
