"use client";
import { createContext, useContext, useEffect, useRef, useState, type ReactNode } from "react";
import { usePathname } from "next/navigation";
import Link from "next/link";
import { PublicIcon } from "./public-ui";
import { publicText as t } from "./messages";
import s from "./public-web.module.css";

const PublicLocaleContext=createContext("es");
export const usePublicWebLocale=()=>useContext(PublicLocaleContext);

// Exact allowlist. Private /admin, /customer, /help, /guide NEVER become public
// through a negative app test. Existing private chrome is supplied unchanged.
export function isPublicWebPath(path: string) {
  return ["/", "/models", "/locations", "/connected"].includes(path) || /^\/models\/[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u.test(path);
}
export function PublicRouteBoundary({ children, privateHeader, privateFooter, brand, supportEmail, locale, privateLocale, navigation }: { children:ReactNode; privateHeader:ReactNode; privateFooter:ReactNode; brand:string; supportEmail:string; locale:string; privateLocale:string; navigation: {href:string; label:string}[] }) {
  const pathname = usePathname();
  if (!isPublicWebPath(pathname)) return <>{privateHeader}<main id="main" className="shell" lang={privateLocale}>{children}</main>{privateFooter}</>;
  return <PublicShell brand={brand} supportEmail={supportEmail} locale={locale} navigation={navigation}>{children}</PublicShell>;
}
function PublicShell({ children, brand, supportEmail, locale, navigation }: {children:ReactNode;brand:string;supportEmail:string;locale:string;navigation:{href:string;label:string}[]}) {
  const pathname=usePathname();
  const [open,setOpen]=useState(false);
  const toggle=useRef<HTMLButtonElement>(null);
  const nav=useRef<HTMLElement>(null);
  const links=navigation.filter(item=>["/models","/locations","/customer"].includes(item.href));
  useEffect(()=>{setOpen(false);},[pathname]);
  useEffect(()=>{if(!open)return; const listener=(event:KeyboardEvent)=>{if(event.key==="Escape"){setOpen(false);toggle.current?.focus();}}; document.addEventListener("keydown",listener);return()=>document.removeEventListener("keydown",listener);},[open]);
  useEffect(()=>{const query=window.matchMedia("(min-width: 761px)");const reset=()=>{if(query.matches)setOpen(false);};query.addEventListener("change",reset);return()=>query.removeEventListener("change",reset);},[]);
  return <PublicLocaleContext.Provider value={locale}><div className={s.root} lang={locale} data-public-web="V403-0.1.0">
    <header className={s.header}>
      <div className={s.headerInner}>
        <Link href="/" className={s.brand}>{brand}</Link>
        <button ref={toggle} type="button" className={s.menuButton} aria-label={t(locale,open?"closeMenu":"menu")} aria-controls="public-navigation" aria-expanded={open} onClick={()=>setOpen(value=>!value)}><PublicIcon name={open?"close":"menu"}/></button>
        <nav ref={nav} id="public-navigation" className={`${s.navigation} ${open?s.navigationOpen:""}`} aria-label={t(locale,"navigation")}>
          {links.map(item=><a key={item.href} href={item.href} className={item.href==="/customer"?s.accountLink:undefined} aria-current={pathname===item.href||(item.href==="/models"&&pathname.startsWith("/models/"))?"page":undefined} onClick={()=>setOpen(false)}>{item.href==="/locations"&&<PublicIcon name="pin"/>}{item.href==="/models"&&<PublicIcon name="model"/>}{t(locale,item.href==="/models"?"models":item.href==="/locations"?"locations":"account")}</a>)}
        </nav>
      </div>
    </header>
    <main id="main" className={s.main}>{children}</main>
    <footer className={s.footer}><Link href="/" className={s.brand}>{brand}</Link><nav aria-label={t(locale,"contact")}><a href={`mailto:${supportEmail}`}><PublicIcon name="mail"/>{t(locale,"contact")}</a><Link href="/locations">{t(locale,"locations")}</Link></nav></footer>
  </div></PublicLocaleContext.Provider>;
}
