import{it,expect}from"vitest";
import{publicationDetail,publicationQuery,publishingCommand,publicState}from"./contract";
export const request={intent:{tenant_id:"tenant",page_id:"123",profile_sha256:"a".repeat(64),approval_id:"publication-00001",delivery_key:"social_"+"b".repeat(64),operation:"publish",message:"Una nueva colección llega esta semana.",content_sha256:"c".repeat(64)},scheduled_at:"2026-09-14T12:00:00Z",expires_at:"2026-09-15T12:00:00Z",original_approval_id:"",provider_reference:""};
export const detail={organization_id:"store",requested_by_self:false,status:{approval_id:request.intent.approval_id,approval_state:"approved",request_sha256:"d".repeat(64),payload:request,delivery_state:"",terminal_code:""}};
export const receipt={...request.intent,provider_reference:"123_200",state:"published",evidence_sha256:"e".repeat(64)};
it("approved request is scheduled, accepted requires correctly bound receipt",()=>{
 expect(publicState(publicationDetail.parse(detail))).toBe("scheduled");expect(publicationDetail.safeParse({...detail,status:{...detail.status,delivery_state:"accepted"}}).success).toBe(false);
 const{message:_m,operation:_o,...r}=receipt;const accepted={...detail,status:{...detail.status,receipt:r,delivery_state:"accepted"}};
 expect(publicState(publicationDetail.parse(accepted))).toBe("published");for(const key of["tenant_id","page_id","approval_id","profile_sha256","content_sha256","delivery_key","provider_reference"]){expect(publicationDetail.safeParse({...accepted,status:{...accepted.status,receipt:{...r,[key]:"other"}}}).success).toBe(false)}
});
it("query cannot override tenant, organization or duplicate cursor shape",()=>{for(const q of[{id:"x",after:"y"},{tenant_id:"other"},{id:"../../secret"}])expect(publicationQuery.safeParse(q).success).toBe(false)});
it("draft validates text bytes, exact action, date order and existing-reference-only removal",()=>{
 const draft={approval_id:request.intent.approval_id,operation:"publish",message:request.intent.message,original_approval_id:"",scheduled_at:request.scheduled_at,expires_at:request.expires_at};expect(publishingCommand.safeParse({action:"prepare",draft}).success).toBe(true);
 for(const delta of[{approval_id:"short"},{message:"🙂".repeat(5000)},{expires_at:request.scheduled_at},{page_id:"999"},{operation:"revoke",message:"",original_approval_id:""}])expect(publishingCommand.safeParse({action:"prepare",draft:{...draft,...delta}}).success).toBe(false);
 expect(publishingCommand.safeParse({action:"decide",approval_id:"a",request_sha256:"d".repeat(64),reason:""}).success).toBe(false);
});
