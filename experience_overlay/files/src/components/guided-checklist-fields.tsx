"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
// AUTHORED form adapter, preserving the original owner command and recovery fence.
import { useEffect, useState } from "react";
import { checklistSchema, type PublishedChecklist } from "@/platform/experience/checklists";
import { StatePanel } from "@/components/experience-ui";
import { validLiteralUnitCode } from "@/platform/experience/unit-code";
export function GuidedChecklistFields({ organization, handover, onReady }: {
    organization: string;
    handover?: {
        id: string;
        version: number;
    } | undefined;
    onReady: (ready: boolean) => void;
}) {
    const [evidenceOptions, setEvidenceOptions] = useState<{
        value: string;
        label: string;
    }[]>([]);
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const [options, setOptions] = useState<PublishedChecklist[] | null>(null), [selected, setSelected] = useState(0), [state, setState] = useState<"loading" | "error" | "forbidden" | "slow">("loading"), [attempt, setAttempt] = useState(0);
    useEffect(() => {
        let active = true;
        const controller = new AbortController(), timer = setTimeout(() => {
            if (active)
                setState("slow");
        }, 1500), deadline = setTimeout(() => controller.abort(), 10000);
        setState("loading");
        setOptions(null);
        onReady(false);
        void fetch("/api/enterprise/experience-checklists?" + new URLSearchParams({ organizationId: organization }), { cache: "no-store", signal: controller.signal }).then(async (response) => {
            if (!response.ok) {
                if (response.status === 401 || response.status === 403) {
                    if (active)
                        setState("forbidden");
                    return;
                }
                throw new Error("unavailable");
            }
            const raw = await response.json();
            if (!Array.isArray(raw.checklists) || raw.checklists.length > 32)
                throw new Error("invalid options");
            const values = raw.checklists.map((v: unknown) => checklistSchema.parse(v));
            if (values.some((v: PublishedChecklist) => v.organization_id !== organization))
                throw new Error("scope");
            if (!Array.isArray(raw.evidenceOptions) || raw.evidenceOptions.some((v: {
                value: string;
                label: string;
            }) => !/^[a-f0-9]{64}$/.test(v.value) || typeof v.label !== "string"))
                throw new Error("evidence options");
            if (active) {
                setOptions(values);
                setEvidenceOptions(raw.evidenceOptions);
                setSelected(0);
            }
        }).catch(() => {
            if (active)
                setState("error");
        }).finally(() => { clearTimeout(timer); clearTimeout(deadline); });
        return () => { active = false; clearTimeout(timer); clearTimeout(deadline); controller.abort(); };
    }, [organization, attempt]);
    const checklist = options?.[selected];
    useEffect(() => { onReady(Boolean(handover && checklist && !checklist.items.some(item => item.required && item.response_type === "evidence" && evidenceOptions.length === 0))); }, [handover?.id, checklist, evidenceOptions, onReady]);
    if (!handover)
        return <>
        <p role="alert">{tx("Elegí una entrega desde tus pedidos autorizados para continuar.")}</p>
        </>;
    if (!options || !checklist)
        return <>
        <StatePanel state={options?.length === 0 ? "empty" : state} onRetry={() => setAttempt(v => v + 1)}/>
        </>;
    return <>
    <input type="hidden" name="handoverId" value={handover.id}/>
    <input type="hidden" name="handoverVersion" value={handover.version}/>
    <input type="hidden" name="checklistId" value={checklist.id}/>
    <input type="hidden" name="checklistVersion" value={checklist.version}/>
 
 <label>{tx("Lista de verificación")}<select name="checklistSelection" value={selected} onChange={event => setSelected(Number(event.target.value))}>{options.map((item, index) => <option key={`${item.id}:${item.version}`} value={index}>{item.title}</option>)}</select>
    </label>
 <div className="stack" key={`${checklist.id}:${checklist.version}`}>{checklist.items.map(item => item.response_type === "confirmation" ? <label className="checkRow" key={item.id}>
        <input type="checkbox" name={`answer:${item.id}`} value="confirmed" required={item.required}/>
        <span>{item.prompt}{item.required ? <span aria-hidden="true"> *</span> : tx(" (opcional)")}</span>
        </label> : <label key={item.id}>
        <span>{item.prompt}{item.required ? <span aria-hidden="true"> *</span> : tx(" (opcional)")}</span>
        <>{item.response_type === "evidence" ? <select name={`answer:${item.id}`} required={item.required} defaultValue="">
            <option value="">{tx("Elegí un expediente verificado")}</option>{evidenceOptions.map(e => <option key={e.value} value={e.value}>{e.label}</option>)}</select> : <input name={`answer:${item.id}`} required={item.required} maxLength={item.response_type === "serial" ? 128 : 2048} autoComplete="off" autoCapitalize="off" spellCheck={false} onKeyDown={event => { if (item.response_type === "serial" && event.key === "Enter") { event.preventDefault(); event.stopPropagation(); } }} onInput={event => { if (item.response_type === "serial") event.currentTarget.setCustomValidity(event.currentTarget.value && !validLiteralUnitCode(event.currentTarget.value) ? tx("Revisá el número de serie.") : ""); }}/>}</>{item.response_type === "evidence" ? <small>{tx("Elegí la inspección correspondiente a esta entrega.")}</small> : null}</label>)}</div>
 <input type="hidden" name="selectedChecklist" value={JSON.stringify(checklist)}/>
 <details>
    <summary>{tx("Antes de presentar la entrega")}</summary>
    <p>{tx("Si se corta la conexión al confirmar, usá «Consultar estado» antes de volver a enviar.")}</p>
    </details>
    </>;
}
