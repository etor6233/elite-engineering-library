"use client";
import type {ReactNode} from "react";
import type {DocumentClass,DocumentFields} from "@/platform/documents-vnext/contract";
import s from "./documents.module.css";
export function DocumentIcon(){return <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M6 3h8l4 4v14H6zM14 3v5h4M9 12h6M9 16h5"/></svg>}
export function DocumentNotice({children,error=false}:{children:ReactNode;error?:boolean}){return <p className={error?s.error:s.notice} role={error?"alert":"status"} aria-live="polite">{children}</p>}
export function DocumentFieldsForm({kind,values,onChange,disabled=false}:{kind:DocumentClass;values:DocumentFields;onChange:(key:string,value:string)=>void;disabled?:boolean}){return <div className={s.fields}>{kind.fields.map(field=><label key={field.key}>{field.label_es}<input name={field.key} required disabled={disabled} value={values[field.key]??""} maxLength={field.max_length} inputMode={field.type==="decimal"?"decimal":field.type==="integer"?"numeric":"text"} type={field.type==="date"?"date":"text"} autoComplete="off" onChange={event=>onChange(field.key,event.target.value)}/></label>)}</div>}
export const documentState=(state:string,en:boolean)=>({QUARANTINED:en?"Received":"Recibido",REVIEW_REQUIRED:en?"Check fields":"Revisar campos",REVIEW_PENDING:en?"Awaiting review":"Esperando revisión",PERSISTED:en?"Confirmed":"Confirmado",REJECTED:en?"Rejected":"Rechazado",QUARANTINE_TERMINAL:en?"Needs attention":"Requiere atención"}[state]??state);
export function documentName(name:string){try{return decodeURIComponent(name)}catch{return name}}
