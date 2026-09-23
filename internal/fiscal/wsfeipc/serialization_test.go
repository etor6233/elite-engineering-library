package wsfeipc

import (
 "encoding/json"
 "testing"
)

func TestEmptyFiscalWireCollectionsRemainArrays(t *testing.T) {
 invoice := fixtureInvoice()
 invoice.VATLines = nil
 invoice.OtherTaxLines = nil
 invoice.AssociatedVouchers = nil
 raw,err := json.Marshal(request(invoice));if err!=nil{t.Fatal(err)}
 var payload map[string]json.RawMessage;if err=json.Unmarshal(raw,&payload);err!=nil{t.Fatal(err)}
 for _,field:=range []string{"vat_lines","other_tax_lines","associated_vouchers"}{if string(payload[field])!="[]"{t.Fatalf("%s=%s: .NET collection contract requires an array",field,payload[field])}}
}
