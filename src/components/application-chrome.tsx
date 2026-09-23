"use client";
import {UnifiedWorkspaceChrome} from "@/components/unified-workspace/workspace-chrome";
import {isUnifiedWorkspaceRoute,type WorkspaceGroup} from "@/platform/workspace/navigation";
import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState, type ReactNode } from "react";
import { usePathname } from "next/navigation";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
import { PrivateLanguageSwitcher } from "@/components/private-language-switcher";
import { adminChromeText } from "@/components/admin-ops/chrome-messages";
import adminStyles from "@/components/admin-ops/chrome.module.css";
type Workspace = {
    id: string;
    names: string[];
    active: number;
    select: (index: number) => void;
};
const NavigationContext = createContext<{
    workspace: Workspace | null;
    register: (workspace: Workspace) => void;
    clear: (id: string) => void;
}>({ workspace: null, register: () => { }, clear: () => { } });
// AUTHORED shared navigation glue. The page supplies already permission-filtered areas.
export function ApplicationNavigationProvider({ children }: {
    children: ReactNode;
}) {
    const [workspace, setWorkspace] = useState<Workspace | null>(null);
    const register = useCallback((next: Workspace) => setWorkspace(next), []);
    const clear = useCallback((id: string) => setWorkspace(current => current?.id === id ? null : current), []);
    const value = useMemo(() => ({ workspace, register, clear }), [workspace, register, clear]);
    return <NavigationContext.Provider value={value}>{children}</NavigationContext.Provider>;
}
export const useApplicationNavigation = () => useContext(NavigationContext);
// Original, geometric inline paths; no external icon or brand assets are incorporated.
function NavigationIcon({ kind }: {
    kind: string;
}) {
    const paths: Record<string, string> = {
        "Pedidos": "M4 4h16v16H4z M4 9h16 M9 9v11", "Entregas": "M3 7h11v11H3z M14 11h4l3 4v3h-7 M7 18v2 M18 18v2",
        "Agenda": "M5 5h14v15H5z M8 3v4 M16 3v4 M5 10h14", "Gestión diaria": "M4 5h6v6H4z M14 5h6v6h-6z M4 15h6v5H4z M14 15h6v5h-6z",
        "Ayuda": "M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18 M9 9a3 3 0 0 1 6 0c0 2-3 2-3 4 M12 16v1",
    };
    return <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
    <path d={paths[kind] ?? "M4 5h16v14H4z M4 10h16 M9 10v9"}/>
    </svg>;
}
// AUTHORED geometric pictograms. Reference observations transfer hierarchy, never assets.
function AdminNavigationIcon({ kind }: { kind: string }) {
    const paths: Record<string, string> = {
        "Gestión diaria": "M4 4h6v6H4z M14 4h6v6h-6z M4 14h6v6H4z M14 14h6v6h-6z",
        "Operaciones": "M4 4h6v6H4z M14 4h6v6h-6z M4 14h6v6H4z M14 14h6v6h-6z",
        "Pedidos": "M6 3h12v18l-3-2-3 2-3-2-3 2z M9 8h6 M9 12h6",
        "Entregas": "M3 6h11v12H3z M14 10h4l3 4v4h-7 M7 18a2 2 0 1 0 0 4 2 2 0 0 0 0-4 M18 18a2 2 0 1 0 0 4 2 2 0 0 0 0-4",
        "Agenda": "M5 5h14v15H5z M8 3v4 M16 3v4 M5 10h14 M9 14h2 M14 14h1",
        "Mensajes": "M4 4h16v13H9l-5 4z M8 8h8 M8 12h5",
        "Saldos": "M4 6h16v14H4z M4 6V4h13 M14 11h6v5h-6z M17 13.5h.01",
        "Reintentar": "M4 9a8 8 0 1 1 0 6 M4 4v5h5",
        "Ayuda": "M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18 M9 9a3 3 0 0 1 6 0c0 2-3 2-3 4 M12 16v1",
        "/experience/documents": "M6 3h8l4 4v14H6z M14 3v5h4 M9 12h6 M9 16h5",
        "/factory": "M3 21V10l6 3V8l6 3V3h4l2 18z M7 17h1 M12 17h1 M17 17h1",
        "/supply": "M3 7l9-4 9 4v10l-9 4-9-4z M3 7l9 5 9-5 M12 12v9",
        "/warranty": "M12 3l8 4v5c0 5-8 9-8 9s-8-4-8-9V7z M8 12l3 3 5-6",
        "/network": "M4 4h5v5H4z M15 15h5v5h-5z M15 4h5v5h-5z M9 6h6 M17 9v6 M6 9v8h9",
    };
    return <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d={paths[kind] ?? "M4 5h16v14H4z M4 10h16 M9 10v9"}/></svg>;
}
function WorkspaceMark() {
    return <span className={adminStyles.brandMark} aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M4 6h7v12H4z M15 6h5v5h-5z M15 15h5v3h-5z"/></svg></span>;
}
export function ApplicationChrome({ kind, children, brand, signedIn = false, canPrepareDelivery = false, navigation = [], workspaceNavigation = [] }: {
    workspaceNavigation?: WorkspaceGroup[];
    kind: "header" | "footer";
    children: ReactNode;
    brand?: string;
    signedIn?: boolean;
    canPrepareDelivery?: boolean;
    navigation?: {
        href: string;
        label: string;
    }[];
}) {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const pathname = usePathname();
    const adminScope = pathname === "/franchise" || pathname === "/experience" || pathname.startsWith("/experience/");
    const at = (text: string) => adminChromeText(text, locale.language);
    const { workspace } = useApplicationNavigation();
    const dialog = useRef<HTMLDialogElement>(null);
    const opener = useRef<HTMLButtonElement>(null);
    const closeButton = useRef<HTMLButtonElement>(null);
    const [open, setOpen] = useState(false);
    const app = ["/experience", "/franchise", "/factory", "/supply", "/warranty", "/network"].some(path => pathname === path || pathname.startsWith(path + "/"));
    const closeMenu = useCallback(() => { dialog.current?.close(); setOpen(false); opener.current?.focus(); }, []);
    useEffect(() => {
        const query = window.matchMedia("(min-width: 901px)");
        const resize = () => { if (query.matches && dialog.current?.open)
            closeMenu(); };
        query.addEventListener("change", resize);
        return () => query.removeEventListener("change", resize);
    }, [closeMenu]);
    if (signedIn && isUnifiedWorkspaceRoute(pathname)) return kind === "footer" ? null : <UnifiedWorkspaceChrome brand={brand ?? ""} groups={workspaceNavigation} workspace={workspace}/>;
    if (!app)
        return children;
    if (kind === "footer")
        return null;
    const primary = workspace && pathname === "/franchise" ? workspace : null;
    const spaces = navigation.filter(item => item.href !== "/help" && !(primary && item.href === "/franchise"));
    const help = navigation.find(item => item.href === "/help");
    function adminBody(mobile: boolean) {
        const taskOrder = ["Gestión diaria", "Pedidos", "Entregas", "Agenda", "Reintentar"];
        const operations = primary?.names.map((name, index) => ({ name, index })).filter(item => !["Saldos", "Mensajes"].includes(item.name)).sort((a, b) => {
            const left = taskOrder.indexOf(a.name), right = taskOrder.indexOf(b.name);
            return (left < 0 ? taskOrder.length : left) - (right < 0 ? taskOrder.length : right);
        }) ?? [];
        const management = primary?.names.map((name, index) => ({ name, index })).filter(item => ["Saldos", "Mensajes"].includes(item.name)) ?? [];
        const workspaceLinks = spaces.filter(item => !(pathname.startsWith("/experience") && item.href === "/experience"));
        const renderTasks = (items: { name: string; index: number }[]) => items.map(({name, index}) => <button className={adminStyles.navItem} type="button" key={name} aria-current={primary?.active === index ? "page" : undefined} onClick={() => { primary?.select(index); if (mobile) closeMenu(); }}><AdminNavigationIcon kind={name}/><span>{tx(name)}</span></button>);
        return <>
            <div className={adminStyles.brand}><a className={adminStyles.brandLink} href={pathname.startsWith("/experience") ? "/experience" : "/franchise"}><WorkspaceMark/><span>{brand}</span></a>{mobile ? <button ref={closeButton} type="button" className={adminStyles.close} aria-label={tx("Cerrar menú")} onClick={closeMenu}><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M6 6l12 12 M18 6 6 18"/></svg></button> : null}</div>
            <nav className={adminStyles.navigation} aria-label={tx("Navegación de la aplicación")}>
                {(operations.length || pathname.startsWith("/experience")) ? <div className={adminStyles.navGroup}><span className={adminStyles.groupTitle}>{at("Operación")}</span>{primary ? renderTasks(operations) : <>
                    {pathname.startsWith("/experience") ? <><a className={adminStyles.navItem} href="/experience" aria-current={pathname === "/experience" ? "page" : undefined}><AdminNavigationIcon kind="Operaciones"/><span>{tx("Operaciones")}</span></a>{canPrepareDelivery ? <a className={adminStyles.navItem} href="/experience/operations" aria-current={pathname === "/experience/operations" ? "page" : undefined}><AdminNavigationIcon kind="Entregas"/><span>{tx("Preparar entrega")}</span></a> : null}</> : null}
                </>}</div> : null}
                {management.length ? <div className={adminStyles.navGroup}><span className={adminStyles.groupTitle}>{at("Gestión")}</span>{renderTasks(management)}</div> : null}
                {workspaceLinks.length ? <div className={adminStyles.navGroup}><span className={adminStyles.groupTitle}>{at("Espacios")}</span>{workspaceLinks.map(item => <a className={adminStyles.navItem} key={item.href} href={item.href} aria-current={pathname === item.href ? "page" : undefined}><AdminNavigationIcon kind={item.href}/><span>{item.label}</span></a>)}</div> : null}
            </nav>
            <div className={adminStyles.bottom}>{help ? <a className={adminStyles.navItem} href={help.href}><AdminNavigationIcon kind="Ayuda"/><span>{tx("Ayuda")}</span></a> : null}{signedIn ? <PrivateLanguageSwitcher compact/> : null}</div>
        </>;
    }
    function body(mobile: boolean) {
        if (adminScope) return adminBody(mobile);
        return <>
      <div className="applicationBrand">
        <a href={pathname.startsWith("/experience") ? "/experience" : "/franchise"}>{brand}</a>{mobile ? <button ref={closeButton} type="button" className="iconButton" aria-label={tx("Cerrar menú")} onClick={closeMenu}>×</button> : null}</div>
      <nav className="applicationNavigation" aria-label={tx("Navegación de la aplicación")}>
        {primary ? primary.names.map((name, index) => <button type="button" key={name} aria-current={primary.active === index ? "page" : undefined} onClick={() => { primary.select(index); if (mobile)
                closeMenu(); }}>
            <NavigationIcon kind={name}/>
            <span>{tx(name)}</span>
            </button>) : <>
          {pathname.startsWith("/experience") ? <>
                <a href="/experience" aria-current={pathname === "/experience" ? "page" : undefined}>
                <NavigationIcon kind="Gestión diaria"/>
                <span>{tx("Operaciones")}</span>
                </a>
                {canPrepareDelivery ? <a href="/experience/operations" aria-current={pathname === "/experience/operations" ? "page" : undefined}>
                <NavigationIcon kind="Entregas"/>
                <span>{tx("Preparar entrega")}</span>
                </a> : null}
                </> : null}
          {spaces.filter(item => ["/franchise", "/supply", "/factory", "/warranty", "/experience/documents"].includes(item.href)).map(item => <a key={item.href} href={item.href} aria-current={pathname === item.href ? "page" : undefined}>
                <NavigationIcon kind={item.label}/>
                <span>{item.label}</span>
                </a>)}
        </>}
      </nav>
      {spaces.length ? <details className="applicationSpaces">
            <summary>{tx("Más espacios")}</summary>
            <nav aria-label={tx("Otros espacios")}>{spaces.map(item => <a key={item.href} href={item.href}>{item.label}</a>)}</nav>
            </details> : null}
      <div className="applicationBottom">{help ? <a href={help.href}>
            <NavigationIcon kind="Ayuda"/>
            <span>{tx("Ayuda")}</span>
            </a> : null}{signedIn ? <PrivateLanguageSwitcher compact/> : null}</div>
    </>;
    }
    return <>
    <aside className={adminScope ? `applicationSidebar ${adminStyles.sidebar}` : "applicationSidebar"} data-admin-chrome={adminScope ? "v403" : undefined}>{body(false)}</aside>
    <header className={adminScope ? `applicationMobileHeader ${adminStyles.mobileHeader}` : "applicationMobileHeader"}>
    {adminScope ? <a className={adminStyles.brandLink} href={pathname.startsWith("/experience") ? "/experience" : "/franchise"}><WorkspaceMark/><span>{brand}</span></a> : null}
    <button ref={opener} type="button" className={adminScope ? `iconButton ${adminStyles.menuButton}` : "iconButton"} aria-label={tx("Abrir menú")} aria-haspopup="dialog" aria-expanded={open} aria-controls="application-mobile-menu" onClick={() => { dialog.current?.showModal(); setOpen(true); closeButton.current?.focus(); }}>
    <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
    <path d="M4 6h16 M4 12h16 M4 18h16"/>
    </svg>
    </button>
    {adminScope ? null : <span>{brand}</span>}
    </header>
    <dialog id="application-mobile-menu" className={adminScope ? `applicationDrawer ${adminStyles.drawer}` : "applicationDrawer"} ref={dialog} aria-label={tx("Menú")} onKeyDown={event => {
        if (event.key !== "Tab") return;
        const container = event.currentTarget;
        const controls = Array.from(container.querySelectorAll<HTMLElement>('a[href],button:not([disabled]),input:not([disabled]):not([type="hidden"]),select:not([disabled]),textarea:not([disabled]),summary,[tabindex]:not([tabindex="-1"])')).filter(control => {
            if (!control.getClientRects().length || getComputedStyle(control).visibility === "hidden") return false;
            // Chromium retains layout boxes for children of closed details. They are not tab stops.
            for (let parent = control.parentElement; parent && parent !== container; parent = parent.parentElement) {
                if (parent instanceof HTMLDetailsElement && !parent.open && !parent.querySelector(":scope > summary")?.contains(control)) return false;
            }
            return true;
        });
        const first = controls[0], last = controls[controls.length - 1];
        if (!first || !last) { event.preventDefault(); return; }
        if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
        else if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
        else if (!controls.includes(document.activeElement as HTMLElement)) { event.preventDefault(); first.focus(); }
    }} onCancel={event => { event.preventDefault(); closeMenu(); }} onClose={() => setOpen(false)}>{body(true)}</dialog>
  </>;
}
