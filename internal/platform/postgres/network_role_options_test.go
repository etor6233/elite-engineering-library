package postgres_test

import (
 "context"
 "testing"
 nr "elite.local/enterprise/internal/networkrole"
 "elite.local/enterprise/internal/platform/identity"
 db "elite.local/enterprise/internal/platform/postgres"
 "github.com/google/uuid"
)
func TestNetworkOptionsScopeAndNewOwnerRecords(t *testing.T){
 pool:=roleMetricPool(t);ctx:=context.Background();tenant:=uuid.NewString();otherTenant:=uuid.NewString()
 for _,id:=range []string{tenant,otherTenant}{
  _,e:=pool.Exec(ctx,`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'network-options-'||$1::text,'Fixture','Fixture')`,id);if e!=nil{t.Fatal(e)}
  _,e=pool.Exec(ctx,`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type,status,version)values($1,'allowed','allowed','Allowed','franchisor','active',1),($1,'foreign','foreign','Foreign','franchisee','active',1)`,id);if e!=nil{t.Fatal(e)}
 }
 owner,e:=db.NewNetworkRole(pool);if e!=nil{t.Fatal(e)}
 p:=identity.Principal{TenantID:tenant,Subject:"operator",Permissions:map[string]struct{}{"network:admin":{}},Organizations:map[string]struct{}{"allowed":{}}}
 rows,e:=owner.Options(ctx,p,"organizations");if e!=nil||len(rows)!=1||rows[0].Value!="allowed"{t.Fatal(rows,e)}
 denied:=p;denied.Permissions=map[string]struct{}{};if _,e=owner.Options(ctx,denied,"organizations");e==nil{t.Fatal("permission missing admitted")}
 admin:=p;admin.Permissions=map[string]struct{}{"*":{}}
 id:=uuid.NewString();_,e=owner.Execute(ctx,admin,nr.Command{CommandID:uuid.NewString(),Action:"create-organization",ScopeOrganizationID:"allowed",EntityID:id,Code:"new-branch",DisplayName:"New branch",Type:"store"});if e!=nil{t.Fatal(e)}
 rows,e=owner.Options(ctx,admin,"organizations");if e!=nil||len(rows)!=3{t.Fatal(rows,e)}
 found:=false;for _,v:=range rows{if v.Value==id{found=true}};if !found{t.Fatal("new owner record absent")}
 rows,e=owner.Options(ctx,p,"entities");if e!=nil||len(rows)!=1{t.Fatal("scope changed implicitly",rows,e)}
 if _,e=owner.Options(ctx,admin,"unbounded");e==nil{t.Fatal("invalid kind admitted")}
}
