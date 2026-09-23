"use client";
import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState, type ReactNode } from "react";
import { usePathname } from "next/navigation";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
import { PrivateLanguageSwitcher } from "@/components/private-language-switcher";
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
export function ApplicationChrome({ kind, children, brand, signedIn = false, canPrepareDelivery = false, navigation = [] }: {
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
    if (!app)
        return children;
    if (kind === "footer")
        return null;
    const primary = workspace && pathname === "/franchise" ? workspace : null;
    const spaces = navigation.filter(item => item.href !== "/help" && !(primary && item.href === "/franchise"));
    const help = navigation.find(item => item.href === "/help");
    function body(mobile: boolean) {
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
    <aside className="applicationSidebar">{body(false)}</aside>
    <header className="applicationMobileHeader">
    <button ref={opener} type="button" className="iconButton" aria-label={tx("Abrir menú")} aria-haspopup="dialog" aria-expanded={open} aria-controls="application-mobile-menu" onClick={() => { dialog.current?.showModal(); setOpen(true); closeButton.current?.focus(); }}>
    <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
    <path d="M4 6h16 M4 12h16 M4 18h16"/>
    </svg>
    </button>
    <span>{brand}</span>
    </header>
    <dialog id="application-mobile-menu" className="applicationDrawer" ref={dialog} aria-label={tx("Menú")} onKeyDown={event => {
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
