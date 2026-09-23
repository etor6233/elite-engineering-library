"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
// AUTHORED presentation glue; native HTML semantics, existing React dependency pin.
import type { ReactNode, ButtonHTMLAttributes } from "react";
export type Tone = "success" | "warning" | "danger" | "muted";
export function StatusBadge({ tone = "muted", children }: {
    tone?: Tone;
    children: ReactNode;
}) { return <span className={`badge badge-${tone}`}>{children}</span>; }
export function Button({ children, className = "", ...props }: ButtonHTMLAttributes<HTMLButtonElement>) { return <button className={`button ${className}`} {...props}>{children}</button>; }
export type ExperienceState = "loading" | "empty" | "error" | "forbidden" | "slow" | "uncertain";
const states = { loading: ["Cargando información", "Estamos consultando tus registros."], empty: ["Todavía no hay resultados", "Probá otro filtro o iniciá una nueva operación."], error: ["No pudimos cargar la información", "Tus cambios no se descartaron. Volvé a consultar."], forbidden: ["Esta operación requiere otro permiso", "Pedí acceso al responsable de tu organización."], slow: ["La consulta está tardando", "Podés seguir trabajando en otra tarea. Conservamos tu información."], uncertain: ["Necesitamos comprobar el resultado", "La operación puede haberse registrado. Consultá su estado antes de volver a enviarla."] } as const;
export function StatePanel({ state, onRetry }: {
    state: ExperienceState;
    onRetry?: () => void;
}) {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const [title, body] = states[state];
    return <section className="card" aria-live="polite" aria-busy={state === "loading" || state === "slow"} role={state === "error" || state === "forbidden" ? "alert" : "status"}>
    <h3>{tx(title)}</h3>
    <p>{tx(body)}</p>{onRetry && ["error", "uncertain", "slow"].includes(state) ? <button onClick={onRetry} type="button">{tx("Consultar estado")}</button> : null}</section>;
}
