package main
import("context";"os";"strings";"testing";nr "elite.local/enterprise/internal/networkrole";"elite.local/enterprise/internal/platform/identity";"elite.local/enterprise/internal/platform/postgres";"github.com/google/uuid";"github.com/jackc/pgx/v5/pgxpool")
func TestNetworkRoleHostGuards(t *testing.T){
 ctx:=context.Background();calls:=0
 if m,e:=selectedNetworkRoleModule(ctx,nil,func(k string)string{calls++;if k!="NETWORK_ROLE_ENABLED"{t.Fatal("unexpected lookup")};return"false"});m!=nil||e!=nil||calls!=1{t.Fatal(m,e,calls)}
 if _,e:=selectedNetworkRoleModule(ctx,nil,func(string)string{return"TRUE"});e==nil{t.Fatal("invalid flag")}
 raw:=os.Getenv("PAYMENT_CONNECTED_DB_URL");if raw==""{t.Skip("owned database required")};pool,e:=pgxpool.New(ctx,raw);if e!=nil{t.Fatal(e)};defer pool.Close()
 lookup:=func(string)string{return"true"}
 enabled:=func(want bool){t.Helper();m,e:=selectedNetworkRoleModule(ctx,pool,lookup);if (e==nil&&m!=nil)!=want{t.Fatal("activation mismatch",want,e)}}
 enabled(true)
 pairs:=[][2]string{
 {`alter table franchise.network_command_receipt disable trigger network_receipt_immutable`,`alter table franchise.network_command_receipt enable trigger network_receipt_immutable`},
 {`alter table org.organization disable trigger organization_hierarchy_guard`,`alter table org.organization enable trigger organization_hierarchy_guard`},
 {`alter table franchise.agreement disable trigger franchise_territory_non_overlap`,`alter table franchise.agreement enable trigger franchise_territory_non_overlap`},
 {`alter index franchise.network_receipt_actor_idx rename to network_receipt_actor_saved`,`alter index franchise.network_receipt_actor_saved rename to network_receipt_actor_idx`},
 {`alter table franchise.network_command_receipt drop constraint network_command_receipt_pkey`,`alter table franchise.network_command_receipt add constraint network_command_receipt_pkey primary key(tenant_id,command_id)`},
 }
 for _,pair:=range pairs{func(){if _,e:=pool.Exec(ctx,pair[0]);e!=nil{t.Fatal(e)};defer func(){if _,e:=pool.Exec(context.Background(),pair[1]);e!=nil{t.Error(e)}}();enabled(false)}();enabled(true)}
 down,e:=os.ReadFile("../../db/migrations/0077_network_role_commands.down.sql");if e!=nil{t.Fatal(e)}
 up,e:=os.ReadFile("../../db/migrations/0077_network_role_commands.up.sql");if e!=nil{t.Fatal(e)}
 if _,e=pool.Exec(ctx,string(down));e!=nil{t.Fatal("empty downgrade",e)};enabled(false)
 if _,e=pool.Exec(ctx,string(up));e!=nil{t.Fatal("reapply",e)};enabled(true)
 tenant:=uuid.NewString();if _,e=pool.Exec(ctx,`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'network-host','Synthetic','Synthetic')`,tenant);e!=nil{t.Fatal(e)}
 store,_:=postgres.NewNetworkRole(pool);p:=identity.Principal{TenantID:tenant,Subject:"fixture",Permissions:map[string]struct{}{"*":{}}}
 if _,e=store.Execute(ctx,p,nr.Command{CommandID:"create",Action:"create-organization",EntityID:"root",Code:"root",DisplayName:"Synthetic root",Type:"franchisor"});e!=nil{t.Fatal(e)}
 tx,e:=pool.Begin(ctx);if e!=nil{t.Fatal(e)};defer tx.Rollback(ctx)
 body:=strings.TrimSpace(string(down));body=strings.TrimPrefix(body,"begin;");body=strings.TrimSuffix(body,"commit;")
 if _,e=tx.Exec(ctx,body);e==nil{t.Fatal("populated evidence discarded")};tx.Rollback(ctx);enabled(true)
 t.Log("NETWORK_ROLE_HOST_PASS five_missing_guards_block=true empty_downgrade_reapply=true populated_downgrade_refused=true no_credentials=true")
}
