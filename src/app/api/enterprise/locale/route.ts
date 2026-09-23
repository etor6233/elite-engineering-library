import{NextResponse,type NextRequest}from"next/server";
import{readSession}from"@/platform/auth/session";
import{applicationBaseUrl}from"@/platform/auth/oidc-client";
import{boundedCommandBody}from"@/platform/backend/bounded-command";
import{privateLocaleCookieName}from"@/platform/i18n/load-private-locale";
const reply=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"private, no-store",vary:"Cookie"}});
export async function POST(r:NextRequest){
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return reply("CROSS_ORIGIN_REJECTED",403)}catch{return reply("UNAVAILABLE",503)}
 const session=await readSession();if(!session)return reply("UNAUTHENTICATED",401);
 if(r.nextUrl.searchParams.size)return reply("INVALID_QUERY",400);
 if(r.headers.get("content-type")!=="application/json")return reply("UNSUPPORTED_MEDIA_TYPE",415);
 let text:string;try{text=await boundedCommandBody(r,64)}catch{return reply("INVALID_PREFERENCE",400)}
 // Exact one-field request rejects duplicate keys as well as identity/policy fields.
 const parsed=/^\{\s*"language"\s*:\s*"(es|en)"\s*\}$/.exec(text);if(!parsed)return reply("INVALID_PREFERENCE",400);
 const response=NextResponse.json({language:parsed[1]},{headers:{"cache-control":"private, no-store",vary:"Cookie"}});
 response.cookies.set(privateLocaleCookieName(session),parsed[1]!,{httpOnly:true,secure:process.env.NODE_ENV==="production",sameSite:"lax",path:"/",maxAge:31536000});
 return response;
}
