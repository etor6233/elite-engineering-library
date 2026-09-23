// Finite synthetic qualification host; never a production identity provider.
package main

import (
 "context"
 "crypto/rand"
 "encoding/hex"
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "net"
 "net/http"
 "os"
 "time"

 "elite.local/enterprise/internal/order"
 "elite.local/enterprise/internal/platform/httpapi"
 "elite.local/enterprise/internal/platform/identity"
 "elite.local/enterprise/internal/platform/postgres"
 "github.com/jackc/pgx/v5/pgxpool"
)

type ids struct{}
func (ids) New() string { var b [16]byte; if _,err:=rand.Read(b[:]);err!=nil {panic(err)};return hex.EncodeToString(b[:]) }
type fixtureIdentity struct{}
func (fixtureIdentity) Verify(_ context.Context, raw string)(identity.Principal,error){
 if raw!="synthetic-local-reference-only" {return identity.Principal{},identity.ErrUnauthenticated}
 return identity.Principal{Subject:"synthetic-customer",TenantID:"018f4d4a-7b36-7a21-8d10-2f4c54c28b01",Permissions:map[string]struct{}{"order:create":{}},Organizations:map[string]struct{}{"integration-org":{}}},nil
}
func run() error {
 ttl:=flag.Duration("ttl",18*time.Minute,"finite reference lifetime, at most twenty minutes")
 flag.Parse()
 if *ttl<=0 || *ttl>20*time.Minute {return errors.New("invalid reference lifetime")}
 ctx,cancel:=context.WithTimeout(context.Background(),*ttl);defer cancel()
 pool,err:=pgxpool.New(ctx,os.Getenv("REFERENCE_DATABASE_URL"));if err!=nil{return errors.New("reference database configuration")}
 defer pool.Close()
 if err=pool.Ping(ctx);err!=nil{return errors.New("reference database unavailable")}
 app,metrics,shutdown,err:=instrument(httpapi.New(order.NewService(postgres.NewOrders(pool),ids{}),fixtureIdentity{}));if err!=nil{return errors.New("reference metrics configuration")}
 defer func(){c,done:=context.WithTimeout(context.Background(),5*time.Second);defer done();_ = shutdown(c)}()
 listener,err:=net.Listen("tcp","127.0.0.1:0");if err!=nil{return err}
 metricListener,err:=net.Listen("tcp","127.0.0.1:0");if err!=nil{listener.Close();return err}
 server:=&http.Server{Handler:app,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:15*time.Second,IdleTimeout:30*time.Second,MaxHeaderBytes:16384}
 metricServer:=&http.Server{Handler:metrics,ReadHeaderTimeout:3*time.Second,WriteTimeout:5*time.Second,IdleTimeout:10*time.Second,MaxHeaderBytes:4096}
 failures:=make(chan error,2)
 go func(){failures<-server.Serve(listener)}()
 go func(){failures<-metricServer.Serve(metricListener)}()
 if err=json.NewEncoder(os.Stdout).Encode(map[string]string{"api":listener.Addr().String(),"metrics":metricListener.Addr().String()});err!=nil{return err}
 select {case <-ctx.Done():case err=<-failures:if !errors.Is(err,http.ErrServerClosed){cancel()}}
 c,done:=context.WithTimeout(context.Background(),5*time.Second);defer done()
 a:=server.Shutdown(c);b:=metricServer.Shutdown(c)
 if a!=nil||b!=nil{return errors.New("reference shutdown incomplete")}
 return nil
}
func main(){if err:=run();err!=nil{fmt.Fprintln(os.Stderr,"HTTP_METRICS_REFERENCE_FAILED");os.Exit(1)}}

