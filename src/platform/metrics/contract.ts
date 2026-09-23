// AUTHORED presentation/binding glue over existing domain read models.
import {z} from "zod";
export const metricID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/);
const count=z.string().regex(/^(0|[1-9][0-9]{0,37})$/),signed=z.string().regex(/^-?(0|[1-9][0-9]{0,37})$/);
const point=z.string().regex(/^-?(0|[1-9][0-9]{0,37})(\.[0-9]{1,6})?$/),stamp=z.iso.datetime({offset:true});
export const metricFamilies=[
 {id:"orders",permission:"admin:read",es:"Pedidos de la organización",en:"Organization orders"},
 {id:"own-orders",permission:"customer:self",es:"Mis pedidos",en:"My orders"},
 {id:"leads",permission:"lead:read",es:"Leads por estado",en:"Leads by status"},
 {id:"stock",permission:"admin:read",es:"Stock por estado",en:"Stock by status"},
 {id:"cases",permission:"admin:read",es:"Casos de servicio",en:"Service cases"},
 {id:"own-cases",permission:"customer:self",es:"Mis casos de servicio",en:"My service cases"},
 {id:"shipments",permission:"admin:read",es:"Envíos con origen o destino local",en:"Local origin or destination shipments"},
 {id:"appointments",permission:"appointment:manage",es:"Turnos de la organización",en:"Organization appointments"},
 {id:"own-appointments",permission:"customer:self",es:"Mis turnos",en:"My appointments"},
 {id:"factory-destination",permission:"factory:read",es:"Producción para este destino",en:"Production for this destination"},
 {id:"factory-owned",permission:"supply:factory-read",es:"Producción de esta fábrica",en:"This factory's production"},
 {id:"supply",permission:"supply:read",es:"Compras seriadas del destino",en:"Destination serial purchases"},
 {id:"stored-value",permission:"stored_value:read",es:"Puntos y valor por programa",en:"Points and value by program"},
 {id:"survey",permission:"surveys:read",es:"NPS de una encuesta",en:"Survey NPS"}
] as const;
export function visibleMetricFamilies(permissions:readonly string[]){return metricFamilies.filter(f=>permissions.includes("*")||permissions.includes(f.permission))}
export const operationsMetric=z.object({kind:metricID,organization_id:metricID,scope:z.enum(["organization","customer","origin_or_destination","destination","factory"]),source:z.string().min(1).max(200),basis:z.literal("current_registered_records_by_state"),observed_at:stamp,rows:z.array(z.object({state:metricID,currency:z.string().regex(/^[A-Z]{3}$/).optional(),count,total_minor_units:signed.optional()}).strict()).max(100)}).strict();
export const storedMetric=z.object({organization_id:metricID,program_id:metricID,program_kind:z.enum(["gift_card","loyalty"]),currency:z.string().regex(/^[A-Z]{3}$/),profile_sha256:z.string().regex(/^[a-f0-9]{64}$/),source:z.literal("stored_value.operation + entry + account"),basis:z.literal("committed_ledger_by_operation"),observed_at:stamp,rows:z.array(z.object({operation:z.enum(["issue","accrue","redeem","reverse"]),operations:count,points_delta:point,applied_minor_units:signed}).strict()).max(4)}).strict();
export const surveyMetric=z.object({organization_id:metricID,survey_id:metricID,source:z.literal("crm.survey_definition + crm.survey_response"),basis:z.literal("retained_survey_population_with_configured_minimum"),responses:count,available:z.boolean(),nps:z.number().finite().min(-100).max(100).nullable(),observed_at:stamp}).strict().refine(v=>v.available===(v.nps!==null));
export type Metric= z.infer<typeof operationsMetric>|z.infer<typeof storedMetric>|z.infer<typeof surveyMetric>;
export function parseMetric(kind:string,value:unknown,org:string,id:string){
 const v=kind==="stored-value"?storedMetric.parse(value):kind==="survey"?surveyMetric.parse(value):operationsMetric.parse(value);
 if(v.organization_id!==org||("kind"in v&&v.kind!==kind)||("program_id"in v&&v.program_id!==id)||("survey_id"in v&&v.survey_id!==id))throw new Error("metric scope mismatch");
 return v;
}
