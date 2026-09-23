"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import{useState,useRef}from"react";
import{visibleMetricFamilies,parseMetric,type Metric}from"@/platform/metrics/contract";
export function RoleMetrics({organizations,permissions,locale="es"}:{organizations:string[];permissions:string[];locale?:"es"|"en"}){
 const families=visibleMetricFamilies(permissions),[org,setOrg]=useState(organizations[0]??""),[kind,setKind]=useState<string>(families[0]?.id??""),[id,setID]=useState(""),[value,setValue]=useState<Metric|null>(null),[status,setStatus]=useState<"idle"|"loading"|"unavailable">("idle"),generation=useRef(0),en=locale==="en";
 const reset=()=>{generation.current++;setValue(null);setStatus("idle")};
 async function read(){const current=++generation.current;setValue(null);setStatus("loading");try{const query=new URLSearchParams({kind,organization_id:org,...(kind==="survey"||kind==="stored-value"?{id}:{})}),r=await fetch("/api/enterprise/metrics?"+query,{cache:"no-store",signal:AbortSignal.timeout(6500)});if(!r.ok)throw new Error();const v=parseMetric(kind,await r.json(),org,id);if(current===generation.current){setValue(v);setStatus("idle")}}catch{if(current===generation.current)setStatus("unavailable")}}
 if(!families.length)return <p>{en?"No metric permission is assigned.":"No tenés permisos de indicadores asignados."}</p>;
 return <section aria-label={en?"Operational metrics":"Indicadores operativos"}>
 <h2>{en?"Operational metrics":"Indicadores operativos"}</h2>
 <p>{en?"Current registered records. Each consultation has its own observation time. Order amounts are not collected revenue; currencies and program points remain separate.":"Registros actuales. Cada consulta tiene su propia fecha de observación. Los importes de pedidos no son ingresos cobrados; monedas y puntos de programas permanecen separados."}</p>
 <form onSubmit={e=>{e.preventDefault();void read()}} style={{display:"grid",gap:"1rem",maxWidth:"50rem"}}>
 <label>{en?"Organization":"Organización"}<select value={org} onChange={e=>{reset();setOrg(e.target.value)}}>{organizations.map(o=><option key={o} value={o}>{o}</option>)}</select></label>
 <label>{en?"Indicator":"Indicador"}<select value={kind} onChange={e=>{reset();setKind(e.target.value);setID("")}}>{families.map(f=><option key={f.id} value={f.id}>{f[locale]}</option>)}</select></label>
 {(kind==="stored-value"||kind==="survey")&&<GuidedRecordSelect organization={org} group="metricReferences" name="reference:id" label="Programa o encuesta" disabled={status==="loading"} value={id} onChange={nextValue=>{reset();setID(nextValue)}}/>}
 <button type="submit" disabled={status==="loading"||!org}>{status==="loading"?(en?"Reading…":"Consultando…"):(en?"Consult indicator":"Consultar indicador")}</button></form>
 <div aria-live="polite">{status==="unavailable"&&<p role="alert">{en?"Consultation unavailable. No zero or previous result is shown. Check access and service availability, then retry.":"Consulta no disponible. No se muestra cero ni un resultado anterior. Revisá acceso y disponibilidad del servicio y volvé a consultar."}</p>}
 {value&&<article><h3>{families.find(f=>f.id===kind)?.[locale]}</h3><p>{en?"Organization":"Organización"}: {value.organization_id} · {en?"Observed":"Observado"}: <time dateTime={value.observed_at}>{value.observed_at}</time></p>
 {"scope"in value&&<p>{en?"Scope":"Alcance"}: {value.scope}</p>}
 {"responses"in value?<><p>{en?"Responses":"Respuestas"}: {value.responses}</p><p>{value.available?"NPS: "+value.nps:(en?"NPS unavailable: configured minimum not met.":"NPS no disponible: no se alcanzó el mínimo configurado.")}</p></>:
 <>{value.rows.length===0?<p>{en?"No registered records for this scope.":"Sin registros para este alcance."}</p>:<div style={{overflowX:"auto"}}><table><caption>{en?"Current values by status or operation":"Valores actuales por estado u operación"}</caption><thead><tr><th>{en?"State / operation":"Estado / operación"}</th><th>{en?"Count":"Cantidad"}</th><th>{en?"Currency":"Moneda"}</th><th>{en?"Amount in minor units":"Importe en unidades menores"}</th>{"program_id"in value&&<th>{en?"Point delta":"Variación de puntos"}</th>}</tr></thead><tbody>{"program_id"in value?value.rows.map(r=><tr key={r.operation}><td>{r.operation}</td><td>{r.operations}</td><td>{value.currency}</td><td>{r.applied_minor_units}</td><td>{r.points_delta}</td></tr>):value.rows.map(r=><tr key={r.state+":"+r.currency}><td>{r.state}</td><td>{r.count}</td><td>{r.currency??"—"}</td><td>{r.total_minor_units??"—"}</td></tr>)}</tbody></table></div>}</>}
 {"program_id"in value&&<p>{en?"Program":"Programa"}: {value.program_id} · {value.program_kind} · {value.currency}</p>}
 <details><summary>{en?"Definition and data source":"Definición y fuente de datos"}</summary><p>{value.basis}</p><p>{value.source}</p>{"profile_sha256"in value&&<p style={{overflowWrap:"anywhere"}}>SHA-256: {value.profile_sha256}</p>}</details></article>}</div></section>
}
