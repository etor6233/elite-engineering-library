import { notFound } from "next/navigation";
import { ExperienceCatalogue } from "@/components/experience-catalogue";
export const dynamic = "force-dynamic";
export default function Page() {
    if (process.env.ELITE_EXPERIENCE_CATALOGUE !== "1")
        notFound();
    return <ExperienceCatalogue />;
}
