"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
import { useApplicationNavigation } from "@/components/application-chrome";
import { usePathname } from "next/navigation";
import {FRANCHISE_TASK_QUERY} from "@/platform/workspace/navigation";
import adminStyles from "@/components/admin-ops/chrome.module.css";
import { Children, isValidElement, useState, useCallback, useEffect, useId, type ReactNode } from "react";
// AUTHORED focus pattern. Visited panes remain mounted, preserving operation fences and inputs.
export function TaskWorkspace({ children, labels, title = "Elegí una tarea", initialLabel, sidebar = false }: {
    children: ReactNode;
    labels?: string[];
    title?: string;
    initialLabel?: string;
    sidebar?: boolean;
}) {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const pathname = usePathname();
    const adminScope = pathname === "/franchise" || pathname === "/experience" || pathname.startsWith("/experience/");
    const panes = Children.toArray(children);
    const names = labels ?? panes.map((pane, index) => isValidElement<{
        "data-task-label"?: string;
    }>(pane) ? pane.props["data-task-label"] ?? String(index + 1) : String(index + 1));
    const initial = Math.max(0, names.indexOf(initialLabel ?? names[0] ?? ""));
    // Stable task identity preserves drafts/fences when conditional panes are inserted.
    const identities=panes.map((pane,index)=>isValidElement<{"data-task-id"?:string}>(pane)?pane.props["data-task-id"]??names[index]!:names[index]!);
    const identitiesKey=JSON.stringify(identities);
    const [activeId,setActiveId]=useState(identities[initial]??"");
    const [visited,setVisited]=useState<string[]>([identities[initial]??""]);
    const active=Math.max(0,identities.indexOf(activeId));
    const namesKey = JSON.stringify(names);
    const choose=useCallback((index:number)=>{const ids=JSON.parse(identitiesKey) as string[],id=ids[index];if(id===undefined)return;setVisited(previous=>previous.includes(id)?previous:[...previous,id]);setActiveId(id)},[identitiesKey]);
    const select=useCallback((index:number)=>{
        choose(index);
        if(sidebar&&pathname==="/franchise"){
            const label=(JSON.parse(namesKey) as string[])[index],query=label?FRANCHISE_TASK_QUERY[label]:undefined;
            if(query){const url=new URL(window.location.href);if(url.searchParams.get("task")!==query){url.searchParams.set("task",query);window.history.pushState(null,"",url.pathname+url.search+url.hash)}}
        }
    },[choose,sidebar,pathname,namesKey]);
    useEffect(()=>{
        if(!sidebar||pathname!=="/franchise")return;
        const sync=()=>{const query=new URL(window.location.href).searchParams.get("task"),labels=JSON.parse(namesKey) as string[];const index=query===null?labels.indexOf(initialLabel??labels[0]??""):labels.findIndex(label=>FRANCHISE_TASK_QUERY[label]===query);if(index>=0)choose(index)};
        sync();window.addEventListener("popstate",sync);return()=>window.removeEventListener("popstate",sync);
    },[sidebar,pathname,namesKey,choose,initialLabel]);
    useEffect(()=>{const ids=JSON.parse(identitiesKey) as string[];setActiveId(previous=>ids.includes(previous)?previous:ids[0]??"");setVisited(previous=>{const next=previous.filter(id=>ids.includes(id));return next.length===previous.length?previous:next})},[identitiesKey]);
    const id = useId();
    const { register, clear } = useApplicationNavigation();
    useEffect(() => { if (sidebar)
        register({ id, names: JSON.parse(namesKey), active, select }); }, [sidebar, id, namesKey, active, select, register]);
    useEffect(() => () => { if (sidebar)
        clear(id); }, [sidebar, id, clear]);
    return <div className={adminScope ? `taskWorkspace ${adminStyles.workspace}` : "taskWorkspace"}>
    {((sidebar && (adminScope || names[active] !== "Gestión diaria")) || title === "Trabajo de gestión") ? <h1 className={adminScope ? `focusTitle ${adminStyles.taskTitle}` : "focusTitle"}>{tx(names[active] ?? title)}</h1> : null}
    <nav hidden={sidebar} className={adminScope ? `taskNavigation ${adminStyles.taskNavigation}` : "taskNavigation"} aria-label={tx(title)}>{names.length > 4 ? <label className={adminScope ? `taskPicker ${adminStyles.taskPicker}` : "taskPicker"}>{tx("Trabajo a realizar")}<select value={active} onChange={event => select(Number(event.target.value))}>{names.map((name, index) => <option key={name} value={index}>{tx(name)}</option>)}</select>
        </label> : names.map((label, index) => <button key={label} type="button" aria-current={active === index ? "page" : undefined} onClick={() => select(index)}>{tx(label)}</button>)}</nav>
    <div className={adminScope ? `taskBody ${adminStyles.taskBody}` : "taskBody"}>{panes.map((pane, index) => (visited.includes(identities[index]!) || active===index) ? <div key={identities[index]} hidden={active !== index} data-task-pane={names[index]}>{pane}</div> : null)}</div>
    </div>;
}
