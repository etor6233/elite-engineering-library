"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
import { useApplicationNavigation } from "@/components/application-chrome";
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
    const panes = Children.toArray(children);
    const names = labels ?? panes.map((pane, index) => isValidElement<{
        "data-task-label"?: string;
    }>(pane) ? pane.props["data-task-label"] ?? String(index + 1) : String(index + 1));
    const initial = Math.max(0, names.indexOf(initialLabel ?? names[0] ?? ""));
    const [active, setActive] = useState(initial);
    const [visited, setVisited] = useState<number[]>([initial]);
    const select = useCallback((index: number) => { setVisited(previous => previous.includes(index) ? previous : [...previous, index]); setActive(index); }, []);
    const id = useId();
    const { register, clear } = useApplicationNavigation();
    const namesKey = JSON.stringify(names);
    useEffect(() => { if (sidebar)
        register({ id, names: JSON.parse(namesKey), active, select }); }, [sidebar, id, namesKey, active, select, register]);
    useEffect(() => () => { if (sidebar)
        clear(id); }, [sidebar, id, clear]);
    return <div className="taskWorkspace">
    {((sidebar && names[active] !== "Gestión diaria") || title === "Trabajo de gestión") ? <h1 className="focusTitle">{tx(names[active] ?? title)}</h1> : null}
    <nav hidden={sidebar} className="taskNavigation" aria-label={tx(title)}>{names.length > 4 ? <label className="taskPicker">{tx("Trabajo a realizar")}<select value={active} onChange={event => select(Number(event.target.value))}>{names.map((name, index) => <option key={name} value={index}>{tx(name)}</option>)}</select>
        </label> : names.map((label, index) => <button key={label} type="button" aria-current={active === index ? "page" : undefined} onClick={() => select(index)}>{tx(label)}</button>)}</nav>
    <div className="taskBody">{panes.map((pane, index) => visited.includes(index) ? <div key={names[index]} hidden={active !== index} data-task-pane={names[index]}>{pane}</div> : null)}</div>
    </div>;
}
