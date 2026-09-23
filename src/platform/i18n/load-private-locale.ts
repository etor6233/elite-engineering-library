import"server-only";
import{cache}from"react";
import{createHash}from"node:crypto";
import{cookies,headers}from"next/headers";
import{loadBusinessConfig}from"@/platform/config/load";
import{readSession,type PortalSession}from"@/platform/auth/session";
import{resolvePrivateLocale}from"./private-locale";
// Non-secret display preference is namespaced to the authenticated user/tenant.
// It never contributes to authorization or to the signed session.
export function privateLocaleCookieName(s:Pick<PortalSession,"tenantId"|"subject">){
 return(process.env.NODE_ENV==="production"?"__Host-":"")+"elite_locale_"+createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex");
}
export const loadPrivateLocale=cache(async()=>{
 const [config,session]=await Promise.all([loadBusinessConfig(),readSession()]);
 const market=config.markets.find(x=>x.code===config.business.defaultMarket);if(!market)throw new Error("configured market missing");
 // Anonymous public pages keep their configured public language.
 if(!session)return resolvePrivateLocale(config.business.defaultLocale,market.timeZone,null,null);
 const [jar,requestHeaders]=await Promise.all([cookies(),headers()]);
 return resolvePrivateLocale(config.business.defaultLocale,market.timeZone,jar.get(privateLocaleCookieName(session))?.value,requestHeaders.get("accept-language"));
});
