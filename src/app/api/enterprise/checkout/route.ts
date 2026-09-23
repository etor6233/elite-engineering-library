import { NextResponse } from "next/server";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { checkoutID, checkoutSchema } from "@/platform/payments/checkout";

const reply=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"Cache-Control":"no-store","Referrer-Policy":"no-referrer","X-Content-Type-Options":"nosniff"}});
export async function GET(request: Request) {
 if (["cross-site","none"].includes(request.headers.get("sec-fetch-site")??"")) return reply({code:"CROSS_SITE_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 if(!allowed(session,"customer:self"))return reply({code:"FORBIDDEN"},403);
 const query=new URL(request.url).searchParams;
 const order=query.get("orderId"),organization=query.get("organizationId");
 if(request.url.length>2048||[...query.keys()].length!==2||query.getAll("orderId").length!==1||query.getAll("organizationId").length!==1||!checkoutID.safeParse(order).success||!checkoutID.safeParse(organization).success)return reply({code:"INVALID_CHECKOUT_QUERY"},400);
 if(!session.organizations.includes(organization!)&&!session.permissions.includes("*"))return reply({code:"ORGANIZATION_FORBIDDEN"},403);
 try {
  const raw=await protectedGet<unknown>(session,"/v1/customer/orders/"+encodeURIComponent(order!)+"/checkout",{organization_id:organization!});
  const parsed=checkoutSchema.safeParse(raw);
  if(!parsed.success||parsed.data.order_id!==order)return reply({code:"CHECKOUT_UNAVAILABLE"},502);
  return reply(parsed.data);
 } catch(error) {
  if(error instanceof BackendProblem&&[401,403,404,503].includes(error.status))return reply({code:error.status===404?"CHECKOUT_NOT_AVAILABLE":"CHECKOUT_UNAVAILABLE"},error.status);
  return reply({code:"CHECKOUT_UNAVAILABLE"},503);
 }
}
