import { z } from "zod";
const id=z.string().regex(/^[a-z][a-z0-9-]{0,63}$/),hash=z.string().regex(/^[0-9a-f]{64}$/),attemptId=z.uuid();
const article=z.object({id,version:z.string().regex(/^\d{1,6}\.\d{1,6}\.\d{1,6}$/),title:z.string().min(1).max(160),paragraphs:z.array(z.string().min(1).max(4000)).min(1).max(20)}).strict();
const course=z.object({id,title:z.string().min(1).max(160),role:z.enum(["owner","admin","employee","customer"]),lessons:z.array(id).min(1).max(4),prompts:z.array(z.object({id,text:z.string().min(1).max(1000)}).strict()).min(1).max(8)}).strict();
export const courseView=z.object({course,articles:z.array(article).min(1).max(4),profile_id:id,profile_revision:z.number().int().positive(),profile_sha256:hash,content_sha256:hash,method:z.literal("HUMAN_REVIEW_NO_GRANTS_V1")}).strict();
const answers=z.record(id,z.string().min(1).max(2048)).refine(value=>Object.keys(value).length>=1&&Object.keys(value).length<=8);
export const assessmentSchema=z.object({request_id:z.string().regex(/^training:[0-9a-f]{64}$/),payload_sha256:hash,state:z.enum(["pending","approved","rejected"]),payload:z.object({schema:z.literal("training-assessment/v1"),attempt_id:attemptId,learner_subject:z.string().min(1).max(128),organization_id:z.string().min(1).max(128),content:courseView,answers}).strict(),reviewer:z.string().max(128).optional(),reason:z.string().max(2048).optional(),approved:z.boolean().optional()}).strict();
export const attemptSchema=z.object({attempt_id:attemptId,learner_subject:z.string().min(1).max(128),organization_id:z.string().min(1).max(128),content:courseView,read_lessons:z.array(id).max(4),assessment:assessmentSchema.optional(),current_profile:z.boolean()}).strict();
export const trainingCommand=z.discriminatedUnion("action",[
 z.object({action:z.literal("start"),attempt_id:attemptId,course_id:id,profile_sha256:hash}).strict(),
 z.object({action:z.literal("acknowledge"),attempt_id:attemptId,lesson_id:id,profile_sha256:hash}).strict(),
 z.object({action:z.literal("submit"),attempt_id:attemptId,profile_sha256:hash,answers}).strict(),
 z.object({action:z.literal("assess"),request_id:z.string().regex(/^training:[0-9a-f]{64}$/),payload_sha256:hash,approved:z.boolean(),reason:z.string().trim().min(1).max(2048)}).strict(),
]);
export const trainingReference=z.object({attempt_id:attemptId,course_id:id,profile_sha256:hash}).strict();
export type CourseView=z.infer<typeof courseView>;
export type Assessment=z.infer<typeof assessmentSchema>;
export type Attempt=z.infer<typeof attemptSchema>;
export type TrainingReference=z.infer<typeof trainingReference>;
