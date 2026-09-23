package postgres

// AUTHORED read adapter over the existing network owner; no new grants or writes.
import (
 "context"
 nr "elite.local/enterprise/internal/networkrole"
 "elite.local/enterprise/internal/platform/identity"
)
func (s *NetworkRole) Options(ctx context.Context,p identity.Principal,kind string)([]nr.Option,error){
 if p.TenantID==""||p.Subject==""||(!p.Allowed("network:admin")&&!p.Allowed("franchise:write"))||(kind!="organizations"&&kind!="entities"){return nil,nr.ErrInvalid}
 orgs:=[]string{};for org:=range p.Organizations {orgs=append(orgs,org)}
 rows,e:=s.pool.Query(ctx,`select value,label from (
 select organization_id as value,display_name || ' · ' || organization_code as label from org.organization
 where tenant_id=$1 and ($2 or organization_id=any($3::text[]))
 union all
 select agreement_id as value,territory_code || ' · ' || terms_version as label from franchise.agreement
 where tenant_id=$1 and $4='entities' and ($2 or franchise_organization_id=any($3::text[]))
 ) x order by label,value limit 101`,p.TenantID,p.Allowed("*"),orgs,kind)
 if e!=nil{return nil,e};defer rows.Close();out:=[]nr.Option{}
 for rows.Next(){var v nr.Option;if e=rows.Scan(&v.Value,&v.Label);e!=nil{return nil,e};out=append(out,v)}
 if e=rows.Err();e!=nil{return nil,e};if len(out)>100{return nil,nr.ErrConflict};return out,nil
}
