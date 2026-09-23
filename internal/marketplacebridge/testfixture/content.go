package testfixture

import (
	mb "elite.local/enterprise/internal/marketplacebridge"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

// ContentReceiver adds only the new bounded association and three-field update
// contracts. Existing upload/initial/mutation fixture handlers remain unchanged.
type ContentReceiver struct {
	*InitialReceiver
	Multiple, MissingSales bool
	Sales                  int64
	ContentWrites          int
}

func NewContent(p mb.Profile) *ContentReceiver {
	base := NewInitial(p)
	base.Server.Close()
	s := &ContentReceiver{InitialReceiver: base}
	s.Exists = true
	s.Item.FamilyName = "Old supplier name"
	s.Item.Pictures = []mb.Picture{{ID: "654321-MLA123456789_092026"}}
	base.Server = httptest.NewServer(http.HandlerFunc(s.serveContent))
	return s
}
func (s *ContentReceiver) serveContent(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	handled := true
	defer func() {
		s.Mu.Unlock()
		if !handled {
			s.InitialReceiver.serveInitial(w, r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer "+Token {
		w.WriteHeader(401)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/users/"+s.Profile.SellerID+"/items/search" && r.URL.Query().Get("user_product_id") != "" {
		s.Reads++
		if r.URL.Query().Get("user_product_id") != s.Item.UserProductID || r.URL.Query().Get("limit") != "2" {
			w.WriteHeader(400)
			return
		}
		ids := []string{s.Item.ID}
		if s.Multiple {
			ids = append(ids, "MLA987654321")
		}
		json.NewEncoder(w).Encode(map[string]any{"seller_id": s.Profile.SellerID, "results": ids, "paging": map[string]any{"total": len(ids), "offset": 0}})
		return
	}
	if r.Method == "PUT" && r.URL.Path == "/items/"+s.Item.ID {
		raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
		if e != nil || len(raw) > 32768 {
			w.WriteHeader(413)
			return
		}
		var input map[string]json.RawMessage
		var wanted mb.Item
		if json.Unmarshal(raw, &input) != nil || len(input) != 3 || json.Unmarshal(raw, &wanted) != nil || wanted.FamilyName == "" || len(wanted.Pictures) != 1 || wanted.Pictures[0].ID != PictureID {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"invalid_content"}`)
			return
		}
		s.ContentWrites++
		s.Writes = append(s.Writes, Write{Path: r.URL.Path, Body: string(raw)})
		if s.RejectNext {
			s.RejectNext = false
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"content_validation"}`)
			return
		}
		s.Item.FamilyName = wanted.FamilyName
		s.Item.Attributes = wanted.Attributes
		s.Item.Pictures = wanted.Pictures
		s.Effects++
		if s.DropNext {
			s.DropNext = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		json.NewEncoder(w).Encode(s.Item)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/items/"+s.Item.ID {
		s.Item.SoldQuantity = nil
		if !s.MissingSales {
			n := s.Sales
			s.Item.SoldQuantity = &n
		}
	}
	handled = false
}
