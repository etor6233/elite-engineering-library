"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { adminChromeText } from "@/components/admin-ops/chrome-messages";
import styles from "@/components/admin-ops/chrome.module.css";
import Link from "next/link";
export function ExperienceHome() {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => adminChromeText(text, locale.language);
    return <div className={styles.home} data-admin-home="v403">
        <h1 className={styles.homeTitle}>{tx("Operaciones")}</h1>
        <p className={styles.homeLead}>{tx("Un lugar para cada tarea.")}</p>
        <section className={styles.nextTask} aria-labelledby="next-task-title">
            <div className={styles.nextContent}>
                <span className={styles.taskMark} aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M3 6h11v12H3z M14 10h4l3 4v4h-7 M7 18v3 M18 18v3"/></svg></span>
                <div><h2 id="next-task-title">{tx("Preparar entrega")}</h2><p>{tx("Revisá, presentá y confirmá.")}</p></div>
            </div>
            <Link className={styles.primary} href="/experience/operations">{tx("Continuar entrega")}<svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M4 12h16 M14 6l6 6-6 6"/></svg></Link>
        </section>
        <h2 className={styles.quickTitle}>{tx("Accesos rápidos")}</h2>
        <nav className={styles.quickTasks} aria-label={tx("Accesos rápidos")}>
            <Link className={styles.quickTask} href="/franchise?task=orders"><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M6 3h12v18l-3-2-3 2-3-2-3 2z M9 8h6 M9 12h6"/></svg><span><strong>{tx("Pedidos")}</strong><small>{tx("Consultar y dar seguimiento")}</small></span></Link>
            <Link className={styles.quickTask} href="/franchise?task=agenda"><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M5 5h14v15H5z M8 3v4 M16 3v4 M5 10h14 M9 14h2"/></svg><span><strong>{tx("Agenda")}</strong><small>{tx("Organizar el día")}</small></span></Link>
            <Link className={styles.quickTask} href="/customer"><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M4 4h16v13H9l-5 4z M8 8h8 M8 12h5"/></svg><span><strong>{tx("Atención al cliente")}</strong><small>{tx("Ayuda y consultas")}</small></span></Link>
        </nav>
    </div>;
}
