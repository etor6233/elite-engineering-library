package refundworker

import (
 "context"
 "errors"
 "io"
 "net/http"
 "net/http/httptest"
 "testing"

 stripe "github.com/stripe/stripe-go/v86"
)

// AUTHORED contract regression. The response to GET re_expected must not
// silently replace that resource with another refund on the same payment.
func TestRefundRetrieveRejectsDifferentRefundResource(t *testing.T) {
 server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method != http.MethodGet || r.URL.Path != "/v1/refunds/re_expected" { t.Errorf("wrong lookup %s %s",r.Method,r.URL.Path) }
  w.Header().Set("Content-Type","application/json")
  _,_ = io.WriteString(w,`{"id":"re_other","object":"refund","amount":1000,"currency":"ars","payment_intent":"pi_1","status":"succeeded"}`)
 }))
 defer server.Close()
 backend:=stripe.GetBackendWithConfig(stripe.APIBackend,&stripe.BackendConfig{URL:stripe.String(server.URL),HTTPClient:server.Client(),MaxNetworkRetries:stripe.Int64(0)})
 p,err:=newStripeProvider("sk_test_local",backend)
 if err!=nil { t.Fatal(err) }
 value:=baseRefund();value.ProviderRefundReference="re_expected"
 if result,err:=p.Retrieve(context.Background(),value); !errors.Is(err,ErrResponseMismatch) { t.Fatalf("different refund accepted: %+v err=%v",result,err) }
}

func TestRefundCommonResultPinsRetrievedResource(t *testing.T) {
 value:=baseRefund();value.ProviderRefundReference="re_expected"
 result:=ProviderResult{ProviderPaymentReference:value.ProviderPaymentReference,ProviderRefundReference:"re_other",ProviderStatus:"succeeded",Currency:value.Currency,AmountMinorUnits:value.AmountMinorUnits}
 if err:=validateResult(value,result);!errors.Is(err,ErrResponseMismatch){t.Fatalf("different refund accepted: %v",err)}
 result.ProviderRefundReference="re_expected"
 if err:=validateResult(value,result);err!=nil{t.Fatalf("matching refund rejected: %v",err)}
 value.ProviderRefundReference=""
 if err:=validateResult(value,result);err!=nil{t.Fatalf("newly created refund rejected: %v",err)}
}
