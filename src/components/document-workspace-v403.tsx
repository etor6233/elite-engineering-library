"use client";
// AUTHORED presentation over the existing document owner. Files/fields never enter
// browser recovery storage; only an opaque id, operation and matching hashes do.
import { useEffect, useRef, useState } from "react";
import { commandSchema, documentId, fields, fieldsHash, MAX_DOCUMENT_BYTES, parsePending, reconciles, sha256, viewSchema, type DocumentAccess, type DocumentFields, type DocumentView, type PendingDocument } from "@/platform/documents-v403/contract";
const endpoint = "/api/enterprise/documents-v403";
const empty: DocumentFields = { invoice_number: "", vendor: "", total: "", currency: "" };
const nameOf = (name: string) => { try { return decodeURIComponent(name); } catch { return name; } };

export function DocumentWorkspaceV403({ language = "es", initialDocumentId }: { language?: string; initialDocumentId?: string }) {
  const tx = (es: string, en: string) => language === "en" ? en : es;
  const [access, setAccess] = useState<DocumentAccess | null>(null), [document, setDocument] = useState<DocumentView | null>(null);
  const [file, setFile] = useState<File | null>(null), [values, setValues] = useState<DocumentFields>(empty);
  const [pending, setPending] = useState<PendingDocument | null>(null), [busy, setBusy] = useState(false), [ready, setReady] = useState(false), [message, setMessage] = useState("");
  const [reason, setReason] = useState(""); const inFlight = useRef(false), preparing = useRef(false), edited = useRef(false), currentDocument = useRef<DocumentView | null>(null), accessRef = useRef<DocumentAccess | null>(null);
  const key = (a: DocumentAccess) => `elite-document-v403:${a.scope}`;

  async function acceptView(value: unknown, p: PendingDocument | null, a: DocumentAccess) {
    const view = viewSchema.parse(value); sessionStorage.setItem(`${key(a)}:current`, view.document_id);
    const keepEdits = edited.current && currentDocument.current?.document_id === view.document_id && view.state === "REVIEW_REQUIRED" && !view.proposal;
    currentDocument.current = view; setDocument(view);
    if (!keepEdits) { setValues(view.proposal?.fields ?? view.suggested ?? empty); edited.current = false; }
    if (p && await reconciles(p, view)) { sessionStorage.removeItem(key(a)); setPending(null); }
    else if (p) setMessage(tx("El resultado sigue pendiente. Consultá el estado antes de continuar.", "The result is still pending. Check its status before continuing."));
  }
  async function refresh(id?: string, knownAccess?: DocumentAccess, knownPending?: PendingDocument | null) {
    const a = knownAccess ?? accessRef.current; if (!a) return;
    const p = knownPending === undefined ? pending : knownPending; const target = id ?? p?.id ?? document?.document_id;
    if (!target || !documentId.safeParse(target).success) return;
    setBusy(true);
    try {
      const response = await fetch(`${endpoint}?id=${target}`, { cache: "no-store", signal: AbortSignal.timeout(10000) });
      if (!response.ok) throw new Error("unavailable");
      await acceptView(await response.json(), p, a);
    } catch { setMessage(tx("No pudimos confirmar el estado. La operación no se reenviará.", "We could not confirm the status. The operation will not be resent.")); }
    finally { setBusy(false); }
  }
  useEffect(() => {
    let active = true;
    void (async () => {
      try {
        const response = await fetch(endpoint, { cache: "no-store", signal: AbortSignal.timeout(10000) });
        if (!response.ok) throw new Error("unavailable");
        const a = await response.json() as DocumentAccess;
        if (!/^[a-f0-9]{64}$/u.test(a.scope) || !a.subject || [a.canWrite, a.canProcess, a.canReview].some(v => typeof v !== "boolean") || !["FIXTURE", "PROVIDER"].includes(a.mode)) throw new Error("invalid access");
        if (!active) return; accessRef.current = a; setAccess(a);
        const raw = sessionStorage.getItem(key(a)); const p = raw === null ? null : parsePending(raw); setPending(p); setReady(true);
        const id = p?.id ?? initialDocumentId ?? sessionStorage.getItem(`${key(a)}:current`) ?? undefined;
        if (id) await refresh(id, a, p);
      } catch { if (active) setMessage(tx("Facturas no está disponible para esta sesión. Volvé a abrir el espacio cuando se restablezca el acceso.", "Invoices is unavailable for this session. Reopen this workspace once access is restored.")); }
    })();
    return () => { active = false; };
  }, [initialDocumentId]);

  async function execute(p: PendingDocument, request: RequestInit) {
    const a = accessRef.current; if (!a || !ready || busy || pending || inFlight.current) return;
    const current = currentDocument.current;
    if (p.action === "receive" ? current !== null : !current || current.document_id !== p.id || current.state !== ({ process: "QUARANTINED", review: "REVIEW_REQUIRED", decision: "REVIEW_PENDING" } as const)[p.action]) return;
    inFlight.current = true; setBusy(true);
    try {
      if (sessionStorage.getItem(key(a)) !== null) throw new Error("pending");
      sessionStorage.setItem(key(a), JSON.stringify(p)); setPending(p);
      const response = await fetch(request.method === "PUT" ? `${endpoint}?id=${p.id}` : endpoint, { ...request, signal: AbortSignal.timeout(p.action === "process" ? 135000 : 10000) });
      const value = await response.json();
      if (!response.ok) {
        if (value.effect === "NOT_ATTEMPTED") { sessionStorage.removeItem(key(a)); setPending(null); setMessage(tx("No se envió la operación. Revisá el archivo o los datos y volvé a intentarlo.", "The operation was not sent. Check the file or fields and try again.")); return; }
        throw new Error("unconfirmed");
      }
      setMessage(""); await acceptView(value, p, a);
    } catch { setMessage(tx("Resultado pendiente. Consultá el estado; no vuelvas a cargar ni confirmar.", "Result pending. Check its status; do not upload or confirm again.")); }
    finally { inFlight.current = false; setBusy(false); }
  }

  async function upload() {
    if (!file || !access?.canWrite || pending || busy || inFlight.current || preparing.current || currentDocument.current) return;
    preparing.current = true; setBusy(true);
    try {
      const name = encodeURIComponent(file.name);
      if (!/^[\x21-\x7e]{1,128}$/u.test(name) || !/\.(pdf|jpe?g)$/iu.test(name) || file.size < 1 || file.size > MAX_DOCUMENT_BYTES) { setMessage(tx("Elegí un JPEG o PDF de hasta 2 MiB con un nombre más breve.", "Choose a JPEG or PDF up to 2 MiB with a shorter name.")); return; }
      const bytes = new Uint8Array(await file.arrayBuffer()), digest = await sha256(bytes), id = crypto.randomUUID();
      await execute({ id, action: "receive", expected: digest }, { method: "PUT", headers: { "content-type": "application/octet-stream", "x-document-name": name, "x-document-sha256": digest }, body: new Uint8Array(bytes).buffer });
    } catch { setMessage(tx("No pudimos leer el archivo. Elegilo nuevamente.", "The file could not be read. Select it again.")); }
    finally { preparing.current = false; setBusy(false); }
  }

  async function act(action: "process" | "review" | "decision", approved?: boolean) {
    if (!document || !access || pending || busy || preparing.current || inFlight.current) return;
    preparing.current = true;
    try {
      const command = commandSchema.parse(action === "process" ? { action, id: document.document_id } : action === "review" ? { action, id: document.document_id, evidence_sha256: document.evidence_sha256, fields: fields.parse(values) } : { action, id: document.document_id, payload_sha256: document.payload_sha256, approved, reason });
      const expected = action === "process" ? document.original_sha256 : action === "review" ? await fieldsHash(values) : document.payload_sha256!;
      await execute({ id: document.document_id, action, expected, ...(command.action === "decision" ? { approved: command.approved } : {}) }, { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(command) });
    } catch { setMessage(tx("Completá los campos indicados antes de confirmar.", "Complete the required fields before confirming.")); }
    finally { preparing.current = false; }
  }

  const locked = !ready || busy || !!pending;
  const states: Record<string, string> = { QUARANTINED: tx("Recibida", "Received"), REVIEW_REQUIRED: tx("Revisar campos", "Check fields"), REVIEW_PENDING: tx("Esperando revisión", "Awaiting review"), PERSISTED: tx("Confirmada", "Confirmed"), REJECTED: tx("Rechazada", "Rejected"), QUARANTINE_TERMINAL: tx("Requiere atención", "Needs attention") };
  return <section aria-label={tx("Facturas", "Invoices")} className="documentWorkspace">
    {!document && access?.canWrite ? <div className="surface"><h2>{tx("Nueva factura", "New invoice")}</h2><label>{tx("Archivo", "File")}<input type="file" accept="image/jpeg,application/pdf,.jpg,.jpeg,.pdf" disabled={locked} onChange={event => { setFile(event.target.files?.[0] ?? null); setMessage(""); }}/></label>
      <p className="muted">{file ? `${file.name} · ${Math.ceil(file.size / 1024)} KiB` : tx("JPEG o PDF · una página · hasta 2 MiB", "JPEG or PDF · one page · up to 2 MiB")}</p>
      <button className="primaryAction" type="button" disabled={locked || !file} onClick={() => void upload()}>{tx("Cargar factura", "Upload invoice")}</button></div> : null}
    {document ? <article className="surface"><header><h2>{nameOf(document.name)}</h2><p role="status">{states[document.state]}</p></header>
      <a href={`${endpoint}?id=${document.document_id}&part=original`}>{tx("Descargar original", "Download original")}</a>
      {document.state === "QUARANTINED" && access?.canProcess ? <button type="button" className="primaryAction" disabled={locked} onClick={() => void act("process")}>{tx("Leer factura", "Read invoice")}</button> : null}
      {document.state === "REVIEW_REQUIRED" && access?.canWrite && access.subject === document.uploader ? <form onSubmit={event => { event.preventDefault(); void act("review"); }}>
        {([["invoice_number", tx("Número de factura", "Invoice number")], ["vendor", tx("Proveedor", "Supplier")], ["total", tx("Total", "Total")], ["currency", tx("Moneda", "Currency")]] as const).map(([field, label]) => <label key={field}>{label}<input required value={values[field]} disabled={locked} onChange={event => { edited.current = true; setValues(previous => ({ ...previous, [field]: event.target.value })); }}/></label>)}
        <button className="primaryAction" disabled={locked}>{tx("Enviar a revisión", "Submit for review")}</button></form> : null}
      {document.proposal ? <dl>{Object.entries(document.proposal.fields).map(([name, value]) => <div key={name}><dt>{({ invoice_number: tx("Número", "Number"), vendor: tx("Proveedor", "Supplier"), total: tx("Total", "Total"), currency: tx("Moneda", "Currency") } as Record<string,string>)[name]}</dt><dd>{value}</dd></div>)}</dl> : null}
      {document.evidence_sha256 ? <details><summary>{tx("Evidencia de lectura", "Reading evidence")}</summary>{["security", "provider", "analysis"].map(part => <p key={part}><a href={`${endpoint}?id=${document.document_id}&part=${part}`}>{({ security: tx("Validación del archivo", "File checks"), provider: tx("Resultado de extracción", "Extraction result"), analysis: tx("Registro de lectura", "Reading record") } as Record<string,string>)[part]}</a></p>)}</details> : null}
      {document.state === "REVIEW_PENDING" ? <>
        <button type="button" disabled={busy} onClick={() => { const url = new URL("/experience/documents", window.location.origin); url.searchParams.set("document", document.document_id); void navigator.clipboard.writeText(url.href).then(() => setMessage(tx("Enlace copiado para la persona revisora.", "Review link copied."))).catch(() => setMessage(tx("No se pudo copiar el enlace.", "The link could not be copied."))); }}>{tx("Copiar enlace de revisión", "Copy review link")}</button>
        {access?.canReview && access.subject !== document.uploader ? <div><label>{tx("Motivo de la decisión", "Decision reason")}<textarea required value={reason} disabled={locked} onChange={event => setReason(event.target.value)}/></label><button className="primaryAction" type="button" disabled={locked || !reason.trim()} onClick={() => void act("decision", true)}>{tx("Confirmar factura", "Confirm invoice")}</button><button type="button" disabled={locked || !reason.trim()} onClick={() => void act("decision", false)}>{tx("Rechazar", "Reject")}</button></div> : <p>{tx("Otra persona debe revisar y confirmar esta factura.", "Another person must review and confirm this invoice.")}</p>}
      </> : null}
      {document.reviewer ? <p>{tx("Revisión registrada", "Review recorded")}{document.reason ? ` · ${document.reason}` : ""}</p> : null}
    </article> : null}
    {message ? <p role="status" aria-live="polite">{message}</p> : null}
    {pending || document ? <button type="button" disabled={busy || !ready} onClick={() => void refresh()}>{tx("Consultar estado", "Check status")}</button> : null}
    {document && !pending && ["PERSISTED", "REJECTED", "QUARANTINE_TERMINAL"].includes(document.state) && access?.canWrite ? <button type="button" disabled={busy} onClick={() => { try { sessionStorage.removeItem(`${key(access)}:current`); currentDocument.current = null; edited.current = false; setDocument(null); setFile(null); setValues(empty); setReason(""); setMessage(""); } catch { setMessage(tx("No se pudo abrir una nueva carga.", "A new upload could not be started.")); } }}>{tx("Nueva factura", "New invoice")}</button> : null}
  </section>;
}
