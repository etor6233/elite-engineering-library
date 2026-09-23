"use client";
import { useEffect, useRef, type ReactNode } from "react";
import { adminText } from "./messages";
import styles from "./operations.module.css";

// AUTHORED detail frame. Hiding a visited frame does not reset its owner component.
export function OperationDetail({ active, title, reference, state, language, onBack, children, testId = "operation-detail" }: {
  active: boolean; title: string; reference: string; state?: string | undefined; language: string; onBack: () => void; children: ReactNode; testId?: string;
}) {
  const heading = useRef<HTMLHeadingElement>(null);
  useEffect(() => { if (active) heading.current?.focus({ preventScroll: true }); }, [active]);
  return <div hidden={!active} className={styles.detail} data-testid={testId} data-record-id={reference}>
    <button type="button" className={styles.back} onClick={onBack} data-testid={`${testId}-back`}><svg aria-hidden="true" viewBox="0 0 24 24"><path d="m14 6-6 6 6 6M8 12h12"/></svg>{adminText("back", language)}</button>
    <div className={styles.detailHeading}><h3 ref={heading} tabIndex={-1}>{title}</h3>{state ? <span className={styles.state}>{state}</span> : null}</div>
    <details className={styles.reference}><summary>{adminText("reference", language)}</summary><code>{reference}</code></details>
    <div className={styles.detailBody}>{children}</div>
  </div>;
}
