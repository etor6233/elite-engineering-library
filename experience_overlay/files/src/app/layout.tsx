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
  const session = await readSession();
  const navigation = enabledNavigation(config, session?.permissions ?? null);
  const documentsVisible = process.env.DOCUMENTS_ENABLED === "true" && session && session.tenantId === process.env.DOCUMENTS_TENANT_ID && session.organizations.includes(process.env.DOCUMENTS_ORGANIZATION_ID ?? "") && ["documents:write", "documents:process", "documents:review"].some(permission => allowed(session, permission));
  const appNavigation = [...navigation.map(item => ({ href: item.href, label: message(locale.locale, "nav." + item.id) })), ...(documentsVisible ? [{ href: "/experience/documents", label: locale.language === "en" ? "Invoices" : "Facturas" }] : [])];
  return (
    <html lang={locale.locale}>
      <body><PrivateLocaleProvider locale={locale}><ApplicationNavigationProvider>
        <a className="skipLink" href="#main">{message(locale.locale, "site.skip")}</a>
        <ApplicationChrome kind="header" brand={config.business.name} signedIn={Boolean(session)} canPrepareDelivery={Boolean(session && allowed(session,"handover:manage"))} navigation={appNavigation}><header>
          <div className="shell headerRow">
            <Link className="brand" href="/">{config.business.name}</Link>
            <nav aria-label={message(locale.locale, "site.navigation")}>
              {navigation.map((item) => <Link key={item.id} href={item.href as Route}>{message(locale.locale, "nav." + item.id)}</Link>)}
            </nav>{session?<PrivateLanguageSwitcher/>:null}
          </div>
        </header></ApplicationChrome>
        <main id="main" className="shell" lang={locale.locale}>{children}</main>
        <ApplicationChrome kind="footer"><footer><div className="shell">{message(locale.locale, "site.contact")} <a href={`mailto:${config.business.supportEmail}`}>{config.business.supportEmail}</a></div></footer></ApplicationChrome>
      </ApplicationNavigationProvider></PrivateLocaleProvider></body>
    </html>
  );
}
