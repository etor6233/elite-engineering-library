"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
import Link from "next/link";
export function ExperienceHome() {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    return <div className="experienceApp">
 <h1>{tx("Operaciones")}</h1>
 <section className="nextTask" aria-labelledby="next-task-title">
    <div>
    <span className="taskEyebrow">{tx("ENTREGAS")}</span>
    <h2 id="next-task-title">{tx("Prepará la próxima entrega.")}</h2>
    <p>{tx("Revisá el activo y confirmá que está listo.")}</p>
    <Link className="primaryAction" href="/experience/operations">{tx("Preparar entrega")}<span aria-hidden="true">→</span>
    </Link>
    </div>
    <div className="taskSteps" aria-label={tx("Pasos de la entrega")}>
    <p>
    <span>01</span>{tx("Revisar")}</p>
    <p>
    <span>02</span>{tx("Presentar")}</p>
    <p>
    <span>03</span>{tx("Confirmar")}</p>
    </div>
    </section>
 <nav className="quickTasks" aria-label={tx("Otros trabajos")}>
    <Link href="/franchise">
    <span>{tx("Cotizaciones")}</span>
    <span aria-hidden="true">↗</span>
    </Link>
    <Link href="/franchise?task=agenda">
    <span>{tx("Agenda")}</span>
    <span aria-hidden="true">↗</span>
    </Link>
    <Link href="/customer">
    <span>{tx("Atención al cliente")}</span>
    <span aria-hidden="true">↗</span>
    </Link>
    </nav>
 </div>;
}
