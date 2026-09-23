"use client";
// AUTHORED native-list presentation, sourced only from the authorized owner snapshot.
import { useRef, useState, type ReactNode } from "react";
import { adminText } from "./messages";
import { OperationDetail } from "./operation-detail";
import styles from "./operations.module.css";

export type OperationListItem = { id: string; title: string; subtitle?: string; state?: string };
export function OperationsList({ items, selected, onSelect, language, label, testId = "operations-list" }: {
  items: OperationListItem[]; selected: string | null; onSelect: (id: string) => void; language: string; label: string; testId?: string;
}) {
  const [query, setQuery] = useState(""), [status, setStatus] = useState("");
  const normalized = query.trim().toLocaleLowerCase();
  const visible = items.filter(item => (!status || item.state === status) && (!normalized || [item.title, item.id, item.subtitle, item.state].filter(Boolean).join(" ").toLocaleLowerCase().includes(normalized)));
  const states = Array.from(new Set(items.map(item => item.state).filter((state): state is string => Boolean(state))));
  return <div className={styles.list} data-testid={testId}>
    <div className={styles.filters}>
      <label className={styles.search}><span className={styles.srOnly}>{adminText("search", language)} {label}</span>
        <svg aria-hidden="true" viewBox="0 0 24 24"><circle cx="10.5" cy="10.5" r="6.5"/><path d="m16 16 4 4"/></svg>
        <input type="search" value={query} placeholder={adminText("search", language)} onChange={event => setQuery(event.target.value)} />
      </label>
      {states.length > 1 ? <label className={styles.statusFilter}><span className={styles.srOnly}>{adminText("status", language)}</span><select value={status} onChange={event => setStatus(event.target.value)}><option value="">{adminText("allStates", language)}</option>{states.map(value => <option key={value} value={value}>{value}</option>)}</select></label> : null}
    </div>
    <p className={styles.count} aria-live="polite">{visible.length} {adminText(visible.length === 1 ? "record" : "recordCount", language)}</p>
    <ul className={styles.rows} aria-label={label}>{visible.map(item => <li key={item.id}>
      <button type="button" className={styles.row} aria-current={selected === item.id ? "true" : undefined} data-testid={`${testId}-row`} data-record-id={item.id} onClick={() => onSelect(item.id)}>
        <span className={styles.recordIcon} aria-hidden="true"><svg viewBox="0 0 24 24"><path d="M7 4h10l3 5v11H4V9l3-5Z"/><path d="M4 9h16M9 13h6M12 4v5"/></svg></span>
        <span className={styles.rowMain}><strong>{item.title}</strong>{item.subtitle ? <span>{item.subtitle}</span> : null}</span>
        <span className={styles.rowEnd}>{item.state ? <span className={styles.state}>{item.state}</span> : null}<svg aria-hidden="true" viewBox="0 0 24 24"><path d="m9 6 6 6-6 6"/></svg></span>
      </button>
    </li>)}</ul>
    {!visible.length ? <div className={styles.empty} role="status"><p>{adminText(items.length ? "noMatches" : "empty", language)}</p>{items.length ? <button type="button" className={styles.secondary} onClick={() => { setQuery(""); setStatus(""); }}>{adminText("clearFilters", language)}</button> : null}</div> : null}
  </div>;
}

// One persistent owner instance per visited record; selection never replays a command.
export function OperationsCollection({items, language, label, testId, renderDetail}: {
  items: OperationListItem[]; language: string; label: string; testId: string; renderDetail: (id: string) => ReactNode;
}) {
  const [selected, setSelected] = useState<string|null>(null), [visited, setVisited] = useState<string[]>([]);
  const list = useRef<HTMLDivElement>(null);
  function choose(id: string) { if (!items.some(item => item.id === id)) return; setSelected(id); setVisited(previous => previous.includes(id) ? previous : [...previous, id]); }
  function back() { const previous = selected; setSelected(null); requestAnimationFrame(() => Array.from(list.current?.querySelectorAll<HTMLButtonElement>("button[data-record-id]") ?? []).find(button => button.dataset.recordId === previous)?.focus({ preventScroll: true })); }
  return <div className={styles.split} data-detail-open={Boolean(selected)}>
    <div className={styles.listColumn} ref={list}><OperationsList items={items} language={language} label={label} selected={selected} onSelect={choose} testId={`${testId}-list`} /></div>
    <div>{!selected ? <div className={styles.placeholder}><h3>{adminText("choose", language)}</h3><p>{adminText("oneSelection", language)}</p></div> : null}
      {items.filter(item => visited.includes(item.id)).map(item => <OperationDetail key={item.id} active={selected === item.id} title={item.title} reference={item.id} state={item.state} language={language} onBack={back} testId={`${testId}-detail`}>{renderDetail(item.id)}</OperationDetail>)}
    </div>
  </div>;
}
