import { z } from "zod";
export const surveyID = z.string().regex(/^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$/);
const version = z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/);
export const definitionSchema = z.object({ id: surveyID, prompt: z.string().min(1).max(1000), consent_version: version, consent_notice: z.string().min(1).max(4000), accepting: z.boolean() }).strict();
export const answerSchema = z.object({ score: z.number().int().min(0).max(10), consent_version: version, received_at: z.iso.datetime({ offset: true }) }).strict();
export const resultSchema = z.object({ answer: answerSchema, replay: z.boolean() }).strict();
export const submissionSchema = z.object({ score: z.number().int().min(0).max(10), consent: z.literal(true), consent_version: version }).strict();
export const summarySchema = z.object({ responses: z.number().int().min(0).max(Number.MAX_SAFE_INTEGER), available: z.boolean(), nps: z.number().min(-100).max(100).nullable() }).strict().refine(v => v.available ? v.responses > 0 && v.nps !== null : v.nps === null);
export type SurveyDefinition = z.infer<typeof definitionSchema>;
export type SurveyAnswer = z.infer<typeof answerSchema>;
export type SurveySubmission = z.infer<typeof submissionSchema>;

// Flat, bounded request grammar preserves duplicate-key rejection through the BFF.
export function parseSubmission(text: string): SurveySubmission {
 if (text.length > 2048) throw new Error("invalid survey body");
 let at = 0;
 const ws = () => { while (at < text.length && /[ \t\r\n]/.test(text[at]!)) at++; };
 const expect = (c: string) => { ws(); if (text[at++] !== c) throw new Error("invalid survey body"); };
 const string = () => {
  ws(); const token = /^"(?:[^"\\\u0000-\u001f]|\\(?:["\\/bfnrt]|u[0-9a-fA-F]{4}))*"/.exec(text.slice(at));
  if (!token) throw new Error("invalid survey body"); at += token[0].length; return JSON.parse(token[0]) as string;
 };
 const values: Record<string, unknown> = Object.create(null);
 expect("{"); ws();
 while (text[at] !== "}") {
  const key = string(); if (!["score", "consent", "consent_version"].includes(key) || Object.hasOwn(values, key)) throw new Error("invalid survey body");
  expect(":"); ws();
  if (key === "consent_version") values[key] = string();
  else {
   const token = (key === "score" ? /^-?(?:0|[1-9][0-9]*)/ : /^(?:true|false)/).exec(text.slice(at));
   if (!token) throw new Error("invalid survey body"); at += token[0].length; values[key] = JSON.parse(token[0]) as unknown;
  }
  ws(); if (text[at] === "}") break;
  expect(","); ws(); if (text[at] === "}") throw new Error("invalid survey body");
 }
 expect("}"); ws(); if (at !== text.length) throw new Error("invalid survey body");
 return submissionSchema.parse(values);
}

