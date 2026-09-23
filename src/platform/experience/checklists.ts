// AUTHORED DTO glue. The Go owner remains authoritative for published checklists.
import { z } from "zod";
import { validLiteralUnitCode } from "./unit-code";
const code = z.string().regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/).max(64);
export const checklistSchema = z.object({ id: code, organization_id: z.string().min(1).max(128), version: z.number().int().positive(), state: z.literal("published"), title: z.string().trim().min(1).max(160), items: z.array(z.object({ id: code, prompt: z.string().trim().min(1).max(500), response_type: z.enum(["confirmation", "text", "serial", "evidence"]), required: z.boolean() })).min(1).max(64) }).passthrough();
export type PublishedChecklist = z.infer<typeof checklistSchema>;
export const checklistSelectionSchema = z.object({ schema_version: z.literal("1"), tenants: z.record(z.string(), z.object({ organizations: z.record(z.string(), z.array(z.object({ id: code, version: z.number().int().positive() })).max(8)) })) }).strict();
export function checklistResponses(checklist: PublishedChecklist, data: FormData) {
    return checklist.items.flatMap(item => {
        const raw = String(data.get(`answer:${item.id}`) ?? "");
        const answer = item.response_type === "serial" ? raw : raw.trim();
        if (!answer && !item.required)
            return [];
        if (!answer || answer.length > 2048)
            throw new Error("Complete required answer");
        if (item.response_type === "serial" && !validLiteralUnitCode(answer))
            throw new Error("Invalid literal serial observation");
        if (item.response_type === "evidence") {
            if (!/^[a-f0-9]{64}$/.test(answer))
                throw new Error("Evidence SHA-256 required");
            return [{ item_id: item.id, response_text: "Evidence verified", evidence_sha256: answer }];
        }
        return [{ item_id: item.id, response_text: answer }];
    });
}
