"use client";
import { CaptureAssist } from "@/components/capture-next/capture-assist";
import { useId, useState, type KeyboardEvent } from "react";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { matchAuthorizedUnitCode, type AuthorizedCodeUnit, type CodeMatch } from "@/platform/experience/unit-code";

// AUTHORED context-bound selection; optional capture resolves on the authenticated server first.
export function AuthorizedUnitCodePicker({ units, complete, selectedIds, disabled, onChoose }: {
  units: AuthorizedCodeUnit[]; complete: boolean; selectedIds: string[]; disabled: boolean; onChoose: (id: string) => void;
}) {
  const { locale } = usePrivateI18n(); const en = locale.language === "en";
  const id = useId(); const [code, setCode] = useState(""); const [result, setResult] = useState<CodeMatch | null>(null);
  const messages = en ? {
    INVALID: "Check the code: it is too long or contains unsupported characters.",
    PARTIAL: "This order view is partial. Load the full order before selecting by code.",
    NONE: "No matching unit in this order. Check the code or the selected order.",
    AMBIGUOUS: "More than one unit matches. Review the order before choosing.",
    UNAVAILABLE: "This unit is not available for the current operation.",
  } : {
    INVALID: "Revisá el código: es demasiado largo o contiene caracteres no admitidos.",
    PARTIAL: "La vista de este pedido es parcial. Cargá el pedido completo antes de elegir por código.",
    NONE: "No hay una unidad coincidente en este pedido. Revisá el código o el pedido elegido.",
    AMBIGUOUS: "Coincide más de una unidad. Revisá el pedido antes de elegir.",
    UNAVAILABLE: "Esta unidad no está disponible para la operación actual.",
  };
  function lookup() { if (!disabled) setResult(matchAuthorizedUnitCode(code, units, complete)); }
  function keyboard(event: KeyboardEvent<HTMLInputElement>) { if (event.key === "Enter") { event.preventDefault(); event.stopPropagation(); lookup(); } }
  const selected = result?.state === "MATCH" && selectedIds.includes(result.unit.id);
  return <section className="unitCodePicker" aria-label={en ? "Find a unit" : "Buscar una unidad"}>
    <label htmlFor={id}>{en ? "Unit code" : "Código de unidad"}</label>
    <div className="unitCodeEntry"><input id={id} type="text" value={code} disabled={disabled} maxLength={128} autoComplete="off" spellCheck={false} aria-describedby={id + "-hint " + id + "-result"} aria-invalid={result?.state === "INVALID" || undefined} onChange={event => { setCode(event.target.value); setResult(null); }} onKeyDown={keyboard}/><button type="button" disabled={disabled || !code} onClick={lookup}>{en ? "Find" : "Buscar"}</button></div>
    <p id={id + "-hint"} className="muted">{en ? "Scan with a keyboard reader, type or paste the serial number." : "Leé con el lector, escribí o pegá el número de serie."}</p>
    <CaptureAssist disabled={disabled||!complete} language={locale.language} acceptRecord={record=>record.kind==="factory_unit"&&units.some(unit=>unit.id===record.id&&unit.serial_number===record.serial&&unit.selectable&&!selectedIds.includes(unit.id))} onResolved={record=>{
      // Decoder/resolver does not authorize mutation. Recheck the complete current order owner.
      const current=matchAuthorizedUnitCode(record.serial,units,complete);
      if(!disabled&&current.state==="MATCH"&&current.unit.id===record.id&&!selectedIds.includes(record.id)){setCode(record.serial);setResult(current);onChoose(record.id)}else setResult(current);
    }}/>
    <div id={id + "-result"} role="status" aria-live="polite">{result?.state === "MATCH" ? <><p><strong>{result.unit.serial_number}</strong> · {result.unit.product} · {result.unit.organization}</p><button type="button" disabled={disabled || selected} onClick={() => {
      // Revalidate against current authorized data immediately before changing local selection.
      const current = matchAuthorizedUnitCode(code, units, complete);
      if (current.state === "MATCH" && current.unit.id === result.unit.id && !selectedIds.includes(current.unit.id)) onChoose(current.unit.id);
      else setResult(current);
    }}>{selected ? (en ? "Unit selected" : "Unidad elegida") : (en ? "Choose this unit" : "Elegir esta unidad")}</button></> : result ? <p>{messages[result.state]}</p> : null}</div>
  </section>;
}
