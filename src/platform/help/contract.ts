import { z } from "zod";
import { ALL_GUIDES } from "./content";
export const helpResponseSchema = z.object({
  schema: z.literal("journey-help/1"),
  articles: z.array(z.object({
    id: z.enum(ALL_GUIDES.map(guide => guide.id)),
    version: z.string().regex(/^\d+\.\d+\.\d+$/).max(32),
    title: z.string().min(1).max(160),
    paragraphs: z.array(z.string().min(1).max(2000)).min(1).max(20)
  }).strict()).max(ALL_GUIDES.length)
}).strict();
export type HelpResponse = z.infer<typeof helpResponseSchema>;
