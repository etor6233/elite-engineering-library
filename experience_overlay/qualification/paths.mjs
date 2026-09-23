import {resolve} from "node:path";
export const target=resolve(process.env.UI_TARGET||".");
export const evidence=resolve(process.env.UI_EVIDENCE||resolve(import.meta.dirname,"evidence"));
