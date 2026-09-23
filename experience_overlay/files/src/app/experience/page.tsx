import { notFound } from "next/navigation";
import { ExperienceHome } from "@/components/experience-home";
export const dynamic = "force-dynamic";
export default function Page() {
    if (process.env.ELITE_EXPERIENCE_CATALOGUE !== "1")
        notFound();
    return <ExperienceHome />;
}
