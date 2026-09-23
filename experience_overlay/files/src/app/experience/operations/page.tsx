import { experienceText } from "@/platform/experience/messages";
import { loadPrivateI18n } from "@/platform/i18n/load-private-i18n";
import { organizationLabel } from "@/platform/experience/options";
import { notFound } from "next/navigation";
import { allowed, readSession } from "@/platform/auth/session";
import { ChecklistCompletionPanel } from "@/components/franchise-command-panel";
import { StatePanel } from "@/components/experience-ui";
import { createHash } from "node:crypto";
// Opt-in fixture page only. Real /franchise already obtains handovers from its Go owner.
export const dynamic = "force-dynamic";
export default async function Page() {
    if (process.env.ELITE_EXPERIENCE_CATALOGUE !== "1")
        notFound();
    const { locale } = await loadPrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const session = await readSession();
    if (!session || !session.organizations.includes("store") || !allowed(session, "handover:manage"))
        return <StatePanel state="forbidden"/>;
    const scope = createHash("sha256").update(JSON.stringify([session.tenantId, session.subject, "store"])).digest("hex");
    return <div className="experienceOperation">
    <a href="/experience">{tx("Volver al inicio")}</a>
    <p className="branchContext">{await organizationLabel(session, "store") ?? tx("Sucursal: ") + "store"}</p>
    <h1 className="pageTitle">{tx("Preparar entrega")}</h1>
    <p className="experienceIntro">{tx("Revisá cada punto antes de presentar el activo.")}</p>
    <ChecklistCompletionPanel organization="store" scope={scope} initialHandover={{ id: "handover-v403", version: 1 }}/>
    </div>;
}
