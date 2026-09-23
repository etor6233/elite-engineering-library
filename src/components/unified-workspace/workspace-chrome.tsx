"use client";
import {useCallback,useEffect,useRef,useState} from "react";
import Link from "next/link";
import type {Route} from "next";
import {usePathname} from "next/navigation";
import {PrivateLanguageSwitcher} from "@/components/private-language-switcher";
import {usePrivateI18n} from "@/platform/i18n/private-provider";
import {workspaceEntryActive,type WorkspaceGroup,type WorkspaceEntry} from "@/platform/workspace/navigation";
import {workspaceText} from "./messages";
import {WorkspaceIcon} from "./workspace-icon";
import s from "./workspace.module.css";
type ActiveWorkspace={names:string[];active:number;select:(index:number)=>void};
// The server provides only permitted installed destinations, never access tokens or grants.
export function UnifiedWorkspaceChrome({brand,groups,workspace}:{brand:string;groups:WorkspaceGroup[];workspace:ActiveWorkspace|null}){
 const {locale}=usePrivateI18n(),pathname=usePathname(),t=(text:string)=>workspaceText(text,locale.language);
 const [query,setQuery]=useState(""),[open,setOpen]=useState(false),[expanded,setExpanded]=useState<Record<string,boolean>>({});
 const dialog=useRef<HTMLDialogElement>(null),opener=useRef<HTMLButtonElement>(null),closer=useRef<HTMLButtonElement>(null);
 const close=useCallback(()=>{dialog.current?.close();setOpen(false);opener.current?.focus()},[]);
 useEffect(()=>{const media=window.matchMedia("(min-width: 901px)");const resized=()=>{if(media.matches&&dialog.current?.open)close()};media.addEventListener("change",resized);return()=>media.removeEventListener("change",resized)},[close]);
 useEffect(()=>{if(dialog.current?.open)close();setQuery("")},[pathname,close]);
 const active=pathname==="/franchise"?workspace:null;
 const activeTask=active?.names[active.active];
 useEffect(()=>{const current=groups.find(group=>group.items.some(entry=>workspaceEntryActive(entry,pathname,activeTask)));if(current)setExpanded(previous=>previous[current.id]===undefined?{...previous,[current.id]:true}:previous)},[groups,pathname,activeTask]);
 const normalize=(value:string)=>value.normalize("NFD").replace(/[\u0300-\u036f]/g,"").toLocaleLowerCase(locale.language).trim();
 const needle=normalize(query);
 const help=groups.flatMap(group=>group.items).find(item=>item.href==="/help");
 const home=groups.flatMap(group=>group.items).find(item=>item.href==="/dashboard")??groups.flatMap(group=>group.items)[0];
 const displayed=groups.map(group=>({...group,items:group.items.filter(item=>{
   if(!needle&&item.href==="/help")return false;
   if(active&&item.taskLabel&&!active.names.includes(item.taskLabel))return false;
   return !needle||normalize(t(item.label)).includes(needle)||normalize(t(group.label)).includes(needle);
 })})).filter(group=>group.items.length);
 // Retry is an actual owner pane, visible only while the page reports a failed read.
 if(active?.names.includes("Reintentar")&&(!needle||normalize(t("Reintentar")).includes(needle))){
   let group=displayed.find(group=>group.id==="my-work");
   if(!group){group={id:"my-work",label:"Mi trabajo",items:[]};displayed.unshift(group)}
   group.items.push({id:"franchise:retry",href:"/franchise",label:"Reintentar",icon:"Reintentar",taskLabel:"Reintentar"});
 }
 function item(entry:WorkspaceEntry,mobile:boolean){
   const current=workspaceEntryActive(entry,pathname,activeTask),index=entry.taskLabel&&active?active.names.indexOf(entry.taskLabel):-1;
   const content=<><WorkspaceIcon kind={entry.icon}/><span>{t(entry.label)}</span></>;
   if(active&&index>=0)return <button key={entry.id} type="button" className={s.item} data-task-id={entry.id} aria-current={current?"page":undefined} onClick={()=>{active.select(index);if(mobile)close()}}>{content}</button>;
   return <Link key={entry.id} prefetch={false} href={entry.href as Route} className={s.item} data-task-id={entry.id} aria-current={current?"page":undefined} onClick={()=>{if(mobile)close()}}>{content}</Link>;
 }
 function body(mobile:boolean){return <>
   <div className={s.brandRow}><Link prefetch={false} href={(home?.href??"/") as Route} className={s.brand}><span className={s.brandMark} aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6"><path d="M4 5h6v14H4z M14 5h6v6h-6z M14 15h6v4h-6z"/></svg></span><span>{brand}</span></Link>{mobile?<button ref={closer} type="button" className={s.iconButton} aria-label={t("Cerrar menú")} onClick={close}><svg aria-hidden="true" viewBox="0 0 24 24"><path d="m6 6 12 12M18 6 6 18"/></svg></button>:null}</div>
   <div className={s.search}><svg aria-hidden="true" viewBox="0 0 24 24"><circle cx="10.5" cy="10.5" r="6.5"/><path d="m15.5 15.5 5 5"/></svg><input type="search" aria-label={t("Buscar tareas")} placeholder={t("Buscar tareas")} value={query} maxLength={100} onChange={event=>setQuery(event.target.value)} autoComplete="off"/>{query?<button type="button" aria-label={t("Limpiar búsqueda")} onClick={()=>setQuery("")}><span aria-hidden="true">×</span></button>:null}</div>
   <nav className={s.navigation} aria-label={t("Navegación de la aplicación")} data-workspace-navigation="U2">{displayed.map(group=>{
     const current=group.items.some(entry=>workspaceEntryActive(entry,pathname,activeTask));
     const visible=Boolean(needle)||(expanded[group.id]??(current||group.id==="my-work"));
     const groupId=`workspace-${mobile?"mobile":"desktop"}-${group.id}`;
     return <section className={s.group} key={group.id} data-navigation-group={group.id} data-current-group={current||undefined}>
       <h2><button className={s.groupToggle} type="button" aria-expanded={visible} aria-controls={groupId} onClick={()=>setExpanded(previous=>({...previous,[group.id]:!visible}))}>{t(group.label)}<svg aria-hidden="true" viewBox="0 0 24 24"><path d={visible?"m7 14 5-5 5 5":"m9 7 5 5-5 5"}/></svg></button></h2>
       <div id={groupId} className={s.groupItems} hidden={!visible}>{group.items.map(entry=>item(entry,mobile))}</div>
     </section>
   })}{displayed.length===0?<p className={s.empty} role="status">{t(needle?"Sin coincidencias":"Sin tareas disponibles")}</p>:null}</nav>
   <div className={s.bottom}>{!needle&&help?item(help,mobile):null}<PrivateLanguageSwitcher compact/></div>
 </>}
 return <>
   <aside className={`applicationSidebar ${s.sidebar}`} data-unified-chrome="U1-U2">{body(false)}</aside>
   <header className={`applicationMobileHeader ${s.mobileHeader}`}><Link className={s.mobileBrand} prefetch={false} href={(home?.href??"/") as Route}>{brand}</Link><button ref={opener} className={s.iconButton} type="button" aria-label={t("Abrir menú")} aria-haspopup="dialog" aria-expanded={open} aria-controls="unified-workspace-menu" onClick={()=>{dialog.current?.showModal();setOpen(true);closer.current?.focus()}}><svg aria-hidden="true" viewBox="0 0 24 24"><path d="M4 6h16M4 12h16M4 18h16"/></svg></button></header>
   <dialog id="unified-workspace-menu" ref={dialog} className={s.drawer} aria-label={t("Menú")} onCancel={event=>{event.preventDefault();close()}} onClose={()=>setOpen(false)} onKeyDown={event=>{
     if(event.key!=="Tab")return;
     const controls=Array.from(event.currentTarget.querySelectorAll<HTMLElement>('a[href],button:not([disabled]),input:not([disabled]),select:not([disabled]),summary,[tabindex="0"]')).filter(control=>{if(!control.getClientRects().length||getComputedStyle(control).visibility==="hidden")return false;for(let parent=control.parentElement;parent&&parent!==event.currentTarget;parent=parent.parentElement){if(parent instanceof HTMLDetailsElement&&!parent.open&&!parent.querySelector(":scope > summary")?.contains(control))return false}return true});
     const first=controls[0],last=controls[controls.length-1];if(!first||!last){event.preventDefault();return}
     if(event.shiftKey&&document.activeElement===first){event.preventDefault();last.focus()}else if(!event.shiftKey&&document.activeElement===last){event.preventDefault();first.focus()}
   }}>{body(true)}</dialog>
 </>;
}
