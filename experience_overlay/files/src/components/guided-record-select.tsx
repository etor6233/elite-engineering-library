"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
import { useEffect, useState, useId, useRef } from "react";
// AUTHORED replacement for operator-entered technical identifiers. No authorization is inferred here.
export function GuidedRecordSelect({ organization, group, name, label, disabled, optional = false, onReady, value, onChange }: {
    organization: string;
    group: string;
    name: string;
    label: string;
    disabled: boolean;
    optional?: boolean;
    onReady?: (ready: boolean) => void;
    value?: string;
    onChange?: (value: string) => void;
}) {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const controlId = useId();
    const control = useRef<HTMLSelectElement>(null);
    const [options, setOptions] = useState<{
        value: string;
        label: string;
    }[] | null>(null), [failure, setFailure] = useState(false), [attempt, setAttempt] = useState(0);
    useEffect(() => {
        let active = true;
        const control = new AbortController(), timer = setTimeout(() => control.abort(), 10000);
        setOptions(null);
        setFailure(false);
        onReady?.(false);
        void fetch("/api/enterprise/experience-options?" + new URLSearchParams({ organizationId: organization, group }), { cache: "no-store", signal: control.signal }).then(async (response) => {
            if (!response.ok)
                throw new Error("unavailable");
            const data = await response.json();
            if (!Array.isArray(data.options) || data.options.length > 100 || data.options.some((v: {
                value?: unknown;
                label?: unknown;
            }) => typeof v.value !== "string" || typeof v.label !== "string"))
                throw new Error("invalid");
            if (active) {
                setOptions(data.options);
                onReady?.(optional || data.options.length > 0);
            }
        }).catch(() => {
            if (active)
                setFailure(true);
        });
        return () => { active = false; clearTimeout(timer); control.abort(); };
    }, [organization, group, attempt]);
    useEffect(() => {
        const form = control.current?.form;
        if (!form)
            return;
        const guard = (event: Event) => {
            const selected = control.current?.value ?? "";
            if (!options || (!optional || selected !== "") && !options.some(option => option.value === selected)) {
                event.preventDefault();
                event.stopPropagation();
                setFailure(true);
            }
        };
        form.addEventListener("submit", guard, true);
        return () => form.removeEventListener("submit", guard, true);
    }, [options, optional]);
    return <div className="stack">
    <label htmlFor={controlId}>{tx(label)}</label>
    <select ref={control} id={controlId} name={name} required={!optional} disabled={disabled || !options} {...(value === undefined ? { defaultValue: "" } : { value, onChange: (event) => onChange?.(event.target.value) })}>
    <option value="">{tx(optional ? "Sin selección (cuando aplique)" : "Elegí una opción")}</option>{options?.map(option => <option key={option.value} value={option.value}>{option.label}</option>)}</select>{failure ? <p role="alert">{tx("No pudimos consultar las opciones autorizadas. ")}<button type="button" onClick={() => setAttempt(v => v + 1)}>{tx("Volver a consultar")}</button>
        </p> : !options ? <p role="status">{tx("Consultando opciones autorizadas…")}</p> : options.length === 0 ? <p role="status">{tx("No hay opciones habilitadas para tu cuenta. Revisá la configuración con el responsable.")}</p> : null}</div>;
}
