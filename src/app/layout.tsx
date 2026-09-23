import {CONNECTED_WORKSPACE_ROUTES} from "@/platform/workspace/installed-route-bindings";
import {documentConfig} from "@/platform/documents-vnext/server";
import {createWorkspaceNavigation} from "@/platform/workspace/navigation";
import {PublicRouteBoundary} from "@/components/public-web/public-shell";
import {ApplicationChrome,ApplicationNavigationProvider} from "@/components/application-chrome";
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import {PrivateLocaleProvider} from "@/platform/i18n/private-provider";
import {PrivateLanguageSwitcher} from "@/components/private-language-switcher";
import {loadPrivateLocale} from "@/platform/i18n/load-private-locale";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import type { Metadata, Route } from "next";
import Link from "next/link";
import { loadBusinessConfig } from "@/platform/config/load";
import { allowed, readSession } from "@/platform/auth/session";
import { enabledNavigation } from "@/platform/config/registry";
import { connection } from "next/server";
import "./globals.css";

export async function generateMetadata(): Promise<Metadata> {
  const config = await loadBusinessConfig();
  const locale = await loadPublicLocale();
  return {
    title: { default: config.business.name, template: `%s | ${config.business.name}` },
    description: message(locale.locale, "site.description"),
    robots: { index: false, follow: false }
  };
}

export default async function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  await connection();
  const config = await loadBusinessConfig();
  const locale = await loadPrivateLocale();
  const publicLocale = await loadPublicLocale();
  const session = await readSession();
  const navigation = enabledNavigation(config, session?.permissions ?? null);
  let documentsVisible=false;
  if(session){try{const dc=documentConfig();documentsVisible=session.tenantId===dc.tenant&&session.organizations.includes(dc.organization)&&["documents:read","documents:write","documents:review","documents:process"].some(permission=>allowed(session,permission));}catch{/* No navigation to an unconfigured document owner. */}}
  const campaignsVisible=process.env.CAMPAIGN_WORKSPACE_ENABLED==="true";
  const aiActivityVisible=config.features.ai_activity===true;
  const appNavigation=[...navigation.map(item=>({href:item.href,label:message(locale.locale,"nav."+item.id)})),...(documentsVisible?[{href:"/documents",label:locale.language==="en"?"Documents":"Documentos"}]:[])];
  const workspaceNavigation=createWorkspaceNavigation(config,session,{installedRoutes:CONNECTED_WORKSPACE_ROUTES,documentsVisible,campaignsVisible,aiActivityVisible});
  return (
    <html lang={locale.locale}>
      <body><PrivateLocaleProvider locale={locale}><ApplicationNavigationProvider>
        <a className="skipLink" href="#main">{message(locale.locale, "site.skip")}</a>
        <PublicRouteBoundary brand={config.business.name} supportEmail={config.business.supportEmail} locale={publicLocale.locale} privateLocale={locale.locale} navigation={navigation.map(item=>({href:item.href,label:item.label}))} privateHeader={<ApplicationChrome kind="header" brand={config.business.name} signedIn={Boolean(session)} canPrepareDelivery={Boolean(session && allowed(session,"handover:manage"))} navigation={appNavigation} workspaceNavigation={workspaceNavigation}><header>
          <div className="shell headerRow">
            <Link className="brand" href="/">{config.business.name}</Link>
            <nav aria-label={message(locale.locale, "site.navigation")}>
              {navigation.map((item) => <Link key={item.id} href={item.href as Route}>{message(locale.locale, "nav." + item.id)}</Link>)}
            </nav>{session?<PrivateLanguageSwitcher/>:null}
          </div>
        </header></ApplicationChrome>} privateFooter={<ApplicationChrome kind="footer" signedIn={Boolean(session)}><footer><div className="shell">{message(locale.locale, "site.contact")} <a href={`mailto:${config.business.supportEmail}`}>{config.business.supportEmail}</a></div></footer></ApplicationChrome>}>{children}</PublicRouteBoundary>
      </ApplicationNavigationProvider></PrivateLocaleProvider></body>
    </html>
  );
}
