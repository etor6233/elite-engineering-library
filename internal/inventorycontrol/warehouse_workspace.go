// AUTHORED application read contract over existing durable warehouse owners.
package inventorycontrol

import (
 "context"
 "encoding/base64"
 "encoding/json"
 "errors"
 "strings"
 "unicode/utf8"
)

var ErrWarehouseWorkspaceQuery = errors.New("invalid warehouse workspace query")
var WarehouseDatasets = map[string]bool{"organizations":true,"items":true,"bins":true,"policies":true,"uoms":true,"balances":true,"reservations":true,"activities":true,"activity-lines":true,"transfers":true,"crossdock":true,"entries":true,"orders":true,"order-lines":true,"bindings":true,"shipments":true,"packaging":true,"transfer-shipments":true,"transfer-receipts":true}
type WarehouseWorkspaceScope struct { Tenant,Organization string; Organizations []string; AllOrganizations bool }
type WarehouseWorkspaceQuery struct { Dataset,Organization,Search,ItemID,ParentID,ID,RequestID,Cursor string; Limit int; After string }
type WarehouseWorkspaceRow struct {
 ID string `json:"id"`
 Label string `json:"label"`
 State string `json:"state"`
 Version string `json:"version"`
 Data json.RawMessage `json:"data"`
}
type WarehouseWorkspacePage struct {
 Schema string `json:"schema"`
 Dataset string `json:"dataset"`
 OrganizationID string `json:"organization_id"`
 Rows []WarehouseWorkspaceRow `json:"rows"`
 NextCursor *string `json:"next_cursor"`
}
type WarehouseWorkspaceReader interface { WarehouseWorkspace(context.Context,WarehouseWorkspaceScope,WarehouseWorkspaceQuery)(WarehouseWorkspacePage,error) }
type warehouseCursor struct { Tenant,Organization,Dataset,Search,ItemID,ParentID,ID,RequestID,After string }
func safeWorkspaceText(s string,n int)bool{if !utf8.ValidString(s)||len(s)>n{return false};for _,r:=range s{if r<32||r==127{return false}};return true}
func (q *WarehouseWorkspaceQuery) Validate(scope WarehouseWorkspaceScope) error {
 if scope.Tenant==""||scope.Organization==""||q.Organization!=scope.Organization||!WarehouseDatasets[q.Dataset]||q.Limit<1||q.Limit>50{return ErrWarehouseWorkspaceQuery}
 for _,s:=range []string{q.Organization,q.ItemID,q.ParentID,q.ID,q.RequestID}{if !safeWorkspaceText(s,128){return ErrWarehouseWorkspaceQuery}}
 if !safeWorkspaceText(q.Search,100)||len(q.Cursor)>4096{return ErrWarehouseWorkspaceQuery}
 if q.Dataset=="activity-lines"||q.Dataset=="order-lines" {if q.ParentID==""{return ErrWarehouseWorkspaceQuery}}
 if q.Cursor!="" {
  b,e:=base64.RawURLEncoding.DecodeString(q.Cursor);if e!=nil{return ErrWarehouseWorkspaceQuery};var c warehouseCursor
  d:=json.NewDecoder(strings.NewReader(string(b)));d.DisallowUnknownFields();if d.Decode(&c)!=nil{return ErrWarehouseWorkspaceQuery}
  want:=q.cursor(scope, c.After);if c!=want||!safeWorkspaceText(c.After,384)||c.After==""{return ErrWarehouseWorkspaceQuery};q.After=c.After
 }
 return nil
}
func(q WarehouseWorkspaceQuery)cursor(s WarehouseWorkspaceScope,after string)warehouseCursor{return warehouseCursor{s.Tenant,s.Organization,q.Dataset,q.Search,q.ItemID,q.ParentID,q.ID,q.RequestID,after}}
func(q WarehouseWorkspaceQuery)Next(s WarehouseWorkspaceScope,after string)string{b,_:=json.Marshal(q.cursor(s,after));return base64.RawURLEncoding.EncodeToString(b)}
