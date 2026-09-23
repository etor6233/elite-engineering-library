package testfixture

// AUTHORED loopback transport fixture. The actual pinned Google SDK builds and
// sends these requests; this receiver never claims Google account approval.
import (
	mb "elite.local/enterprise/internal/merchantbridge"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
)

type Server struct {
	Mu                                           sync.Mutex
	HTTP                                         *httptest.Server
	Profile                                      mb.Profile
	Input                                        map[string]any
	Inserts, Gets                                int
	DropNext, RejectNext, Pending, ForeignSource bool
}

func New(p mb.Profile) *Server {
	s := &Server{Profile: p}
	s.HTTP = httptest.NewServer(http.HandlerFunc(s.serve))
	return s
}
func (s *Server) Close() { s.HTTP.Close() }
func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	bad := func(status int, message string) {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": status, "message": message}})
	}
	if r.Header.Get("Authorization") != "Bearer synthetic-merchant-sdk-fixture" || r.URL.Query().Get("$alt") != "json;enum-encoding=int" {
		bad(400, "fixture auth/query")
		return
	}
	if r.Method == "POST" && r.URL.Path == "/products/v1/accounts/"+s.Profile.AccountID+"/productInputs:insert" {
		s.Inserts++
		if r.URL.Query().Get("dataSource") != s.Profile.DataSource() {
			bad(400, "fixture data source")
			return
		}
		raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
		if e != nil || len(raw) > 32768 {
			bad(413, "fixture body")
			return
		}
		var input map[string]any
		if json.Unmarshal(raw, &input) != nil {
			bad(400, "fixture JSON")
			return
		}
		if input["offerId"] != s.Profile.Bindings[0].OfferID || input["contentLanguage"] != s.Profile.Language || input["feedLabel"] != s.Profile.FeedLabel || input["versionNumber"] != "1" {
			bad(400, "fixture key/version")
			return
		}
		if s.RejectNext {
			s.RejectNext = false
			bad(400, "fixture provider rejection")
			return
		}
		input["name"] = "accounts/" + s.Profile.AccountID + "/productInputs/" + s.Profile.Language + "~" + s.Profile.FeedLabel + "~" + s.Profile.Bindings[0].OfferID
		input["product"] = s.Profile.Resource(s.Profile.Bindings[0].OfferID)
		s.Input = input
		if s.DropNext {
			s.DropNext = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		json.NewEncoder(w).Encode(input)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/products/v1/"+s.Profile.Resource(s.Profile.Bindings[0].OfferID) {
		s.Gets++
		if s.Pending || s.Input == nil {
			bad(404, "fixture processing pending")
			return
		}
		data := map[string]any{}
		for k, v := range s.Input {
			data[k] = v
		}
		delete(data, "product")
		data["name"] = s.Profile.Resource(s.Profile.Bindings[0].OfferID)
		data["dataSource"] = s.Profile.DataSource()
		if s.ForeignSource {
			data["dataSource"] = "accounts/999999/dataSources/999999"
		}
		data["productStatus"] = map[string]any{"destinationStatuses": []any{map[string]any{"reportingContext": "SHOPPING_ADS", "pendingCountries": []string{"AR"}}}, "lastUpdateDate": "2026-09-13T00:00:00Z"}
		json.NewEncoder(w).Encode(data)
		return
	}
	bad(404, "fixture path")
}
