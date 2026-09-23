// AUTHORED original V403 pictograms reused; no external assets.
export function WorkspaceIcon({ kind }: { kind: string }) {
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
        "/documents": "M6 3h8l4 4v14H6z M14 3v5h4 M9 12h6 M9 16h5",
        "/factory": "M3 21V10l6 3V8l6 3V3h4l2 18z M7 17h1 M12 17h1 M17 17h1",
        "/supply": "M3 7l9-4 9 4v10l-9 4-9-4z M3 7l9 5 9-5 M12 12v9",
        "/warranty": "M12 3l8 4v5c0 5-8 9-8 9s-8-4-8-9V7z M8 12l3 3 5-6",
        "/network": "M4 4h5v5H4z M15 15h5v5h-5z M15 4h5v5h-5z M9 6h6 M17 9v6 M6 9v8h9",
    };
    return <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d={paths[kind] ?? "M4 5h16v14H4z M4 10h16 M9 10v9"}/></svg>;
}
