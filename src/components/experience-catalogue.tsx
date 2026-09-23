"use client";
// AUTHORED reference composition. Function contracts are specifications, never authorization.
import { useState, useEffect } from "react";
import { StatusBadge, StatePanel, type ExperienceState } from "@/components/experience-ui";
import matrix from "@/platform/experience/business-functions.v403.json";
import {functionOwnerBinding,functionOwnerDestination} from "@/platform/experience/function-owner-bindings";
const journeys = [{ name: "Captación y cotización", area: "Comercial", href: "/models", description: "Conocé el producto y elegí el siguiente paso." }, { name: "Pedidos y entregas", area: "Operaciones", href: "/franchise", description: "Continuá una operación autorizada con su estado actualizado." }, { name: "Seguimiento del cliente", area: "Cliente", href: "/customer", description: "Consultá pedidos, entregas y atención en tu cuenta." }, { name: "Fábrica e inventario", area: "Operaciones", href: "/factory", description: "Planificá y verificá las unidades de la red." }, { name: "Capacitación y ayuda", area: "Equipo", href: "/help", description: "Encontrá la guía del recorrido que estás realizando." }];
export function ExperienceCatalogue() {
    const [query, setQuery] = useState(""), [area, setArea] = useState("Todos"), [brand, setBrand] = useState("forest"), [state, setState] = useState<ExperienceState>("empty"), [selected, setSelected] = useState(matrix.functions[0]!.id);
    useEffect(() => {
        const resume = () => {
            const id = window.location.hash.replace("#/functions/", "");
            if (matrix.functions.some(f => f.id === id))
                setSelected(id);
        };
        resume();
        window.addEventListener("hashchange", resume);
        return () => window.removeEventListener("hashchange", resume);
    }, []);
    const current = matrix.functions.find(x => x.id === selected)!;
    const filtered = journeys.filter(x => (area === "Todos" || x.area === area) && `${x.name} ${x.description}`.toLocaleLowerCase().includes(query.toLocaleLowerCase()));
    return <div className="catalogue" data-brand={brand} lang="es-AR">
 <div className="cluster">
    <StatusBadge tone="warning">Referencia de biblioteca</StatusBadge>
    <span className="muted">Datos de ejemplo · revisión V403 · aprobación visual pendiente</span>
    </div>
 <section className="hero" aria-labelledby="experience-title">
    <div>
    <p className="eyebrow">Una experiencia compartida</p>
    <h1 id="experience-title">Todo claro.<br />El próximo paso, a mano.</h1>
    <p className="lede">Una presentación comercial simple y herramientas para completar el trabajo de cada día. La misma marca, los mismos controles y ayuda en contexto.</p>
    <div className="actions">
    <a className="button" href="#workspace">Explorar recorridos</a>
    <a href="/experience/operations">Probar entrega guiada</a>
    </div>
    </div>
    <div className="visualMark" aria-hidden="true">Tu marca.<br />Tu equipo.<br />Tu próximo paso.</div>
    </section>
 <section className="card stack" aria-labelledby="brand-title">
    <h2 id="brand-title">Identidad configurable</h2>
    <label>Paleta de referencia<select value={brand} onChange={e => setBrand(e.target.value)}>
    <option value="forest">Bosque</option>
    <option value="indigo">Índigo</option>
    </select>
    </label>
    <p className="muted">Ambas paletas usan los mismos componentes. Una marca adicional debe verificar sus pares reales antes de publicarse.</p>
    <div className="cluster">
    <StatusBadge tone="success">Confirmado</StatusBadge>
    <StatusBadge tone="warning">Requiere revisión</StatusBadge>
    <StatusBadge tone="danger">No se pudo completar</StatusBadge>
    <StatusBadge>Sin actividad</StatusBadge>
    </div>
    </section>
 <section id="workspace" className="workspace" aria-labelledby="workspace-title">
    <aside>
    <p className="eyebrow">Espacio de trabajo</p>
    <nav aria-label="Secciones de la referencia">
    <a href="#workspace" aria-current="page">Recorridos</a>
    <a href="#states">Estados y recuperación</a>
    <a href="#functions">Funciones del equipo</a>
    <a href="#help">Ayuda</a>
    </nav>
    </aside>
    <div className="stack">
    <h2 id="workspace-title">¿Qué necesitás hacer?</h2>
    <div className="toolbar">
    <label>Buscar recorrido<input type="search" value={query} onChange={e => setQuery(e.target.value)} placeholder="Por ejemplo: entrega"/>
    </label>
    <label>Área<select value={area} onChange={e => setArea(e.target.value)}>{["Todos", "Comercial", "Operaciones", "Cliente", "Equipo"].map(v => <option key={v}>{v}</option>)}</select>
    </label>
    </div>
    <p aria-live="polite">{filtered.length} recorridos disponibles</p>{filtered.length ? <>
        <div className="tableWrap journeyTable">
        <table>
        <caption>Recorridos disponibles</caption>
        <thead>
        <tr>
        <th scope="col">Recorrido</th>
        <th scope="col">Área</th>
        <th scope="col">Próximo paso</th>
        </tr>
        </thead>
        <tbody>{filtered.map(item => <tr key={item.href}>
            <th scope="row">{item.name}<p className="muted">{item.description}</p>
            </th>
            <td className="areaCell">{item.area}</td>
            <td>
            <a href={item.href} aria-label={`Continuar: ${item.name.toLocaleLowerCase()}`}>Continuar</a>
            </td>
            </tr>)}</tbody>
        </table>
        </div>
        <ul className="journeyCards">{filtered.map(item => <li className="card stack" key={item.href}>
            <div>
            <StatusBadge>{item.area}</StatusBadge>
            <h3>{item.name}</h3>
            <p>{item.description}</p>
            </div>
            <a href={item.href} aria-label={`Continuar: ${item.name.toLocaleLowerCase()}`}>Continuar</a>
            </li>)}</ul>
        </> : <StatePanel state="empty"/>}<p className="muted">Cada destino verifica la identidad, la organización y los permisos reales. Este índice no concede acceso.</p>
    </div>
    </section>
 <section id="states" className="stack" aria-labelledby="states-title">
    <h2 id="states-title">Estados que ayudan a continuar</h2>
    <label>Estado a revisar<select value={state} onChange={e => setState(e.target.value as ExperienceState)}>
    <option value="loading">Carga</option>
    <option value="empty">Sin resultados</option>
    <option value="error">Error de consulta</option>
    <option value="forbidden">Permiso insuficiente</option>
    <option value="slow">Conexión lenta</option>
    <option value="uncertain">Resultado incierto</option>
    </select>
    </label>
    <StatePanel state={state}/>
    <p className="muted">Muestrario de presentación. La entrega guiada prueba la recuperación de una operación; esta selección no simula que el negocio haya ocurrido.</p>
    </section>
 <section id="functions" className="card stack" aria-labelledby="functions-title">
    <h2 id="functions-title">Once funciones. Permisos explícitos.</h2>
    <label>Función o tema<select value={selected} onChange={e => { setSelected(e.target.value); window.history.replaceState(null, "", `#/functions/${e.target.value}`); }}>{matrix.functions.map(f => <option key={f.id} value={f.id}>{f.track_name}</option>)}</select>
    </label>
    <h3>{current.track_name}</h3>
    <p>{current.purpose}</p>
    <div className="cluster">
    <StatusBadge>{current.status.specification_status}</StatusBadge>
    <StatusBadge tone="warning">Ejecución de tareas: no probada en este catálogo</StatusBadge>
    </div>
    <ol className="compactList">{current.tasks.map(task => <li key={task.id}>
        <strong>{task.title}</strong>
        <p>{functionOwnerBinding(task.id)?.required_product_permissions.length ? `Permisos del destino: ${functionOwnerBinding(task.id)!.required_product_permissions.join(", ")}` : "Trabajo registrado en los archivos del proyecto; no crea permisos de producto."}</p>
        {functionOwnerDestination(task.id)?<a href={functionOwnerDestination(task.id)!}>Abrir espacio de trabajo</a>:null}
        <details>
        <summary>Contrato de la tarea</summary>
        <pre>{JSON.stringify(task, null, 2)}</pre>
        </details>
        </li>)}</ol>
    <details>
    <summary>Procedencia y límites de esta función</summary>
    <p>El tema no es un rol de autorización. La especificación, la admisión del método y su ejecución se registran por separado. Material de Galaxy pendiente de admisión no bloquea métodos oficiales ya disponibles.</p>
    <pre>{JSON.stringify(current.status.method_admission, null, 2)}</pre>
    </details>
    <details>
    <summary>Autoridades, límites y continuidad</summary>
    <pre>{JSON.stringify({ systems: current.systems, allowed: current.allowed, prohibited: current.prohibited, implementation: current.implementation_binding, sources: current.source_refs, status: current.status, owners: current.owner_refs, resume: current.resume }, null, 2)}</pre>
    </details>
    </section>
 <section id="help" className="stack">
    <h2>Ayuda en la misma revisión</h2>
    <details>
    <summary>¿Qué hago si no recibo una confirmación?</summary>
    <p>Conservá la operación abierta y consultá su estado. El sistema debe verificar si ocurrió antes de permitir un nuevo envío.</p>
    </details>
    <details>
    <summary>¿Por qué no aparece una acción?</summary>
    <p>Puede depender del permiso, la organización o el estado del registro. Revisá la explicación del recorrido y pedí al responsable que compruebe el acceso.</p>
    </details>
    <a href="/help">Abrir centro de ayuda</a>
    </section>
 </div>;
}
