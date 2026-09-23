import type { BusinessConfig } from "@/platform/config/schema";
import { enabledNavigation } from "@/platform/config/registry";
export type Section = { id: string; label: string; href: string };

// Call only after readSession succeeds. Common links are explicit, never wildcard grants.
export function visibleSections(permissions: readonly string[], config: BusinessConfig): Section[] {
  const ids: Record<string,string> = { dashboard: "overview", franchise: "sales", public_catalog: "catalog" };
  const sections = enabledNavigation(config, permissions).map(({id,label,href})=>({id:ids[id]??id,label,href}));
  if (config.features.role_workspace === true) sections.push({id:"guide",label:"Guía de uso",href:"/guide"});
  return sections;
}
export const ROLE_GUIDES = [
  {id:"owner",label:"Dueño",blurb:"Accesos de supervisión disponibles en tu sesión.",sections:["overview","admin","sales","factory","catalog","help","training","catalog_editor","supply","warranty","network"]},
  {id:"admin",label:"Administrador",blurb:"Accesos disponibles para la operación diaria.",sections:["overview","admin","sales","catalog","help","training","catalog_editor","supply","warranty","network"]},
  {id:"employee",label:"Empleado",blurb:"Accesos disponibles para atención y operación.",sections:["overview","sales","catalog","help","training","catalog_editor","supply","warranty"]},
  {id:"customer",label:"Cliente",blurb:"Accesos disponibles para tus consultas y entregas.",sections:["overview","customer","catalog","locations","help","training","warranty"]}
] as const;
export function roleGuide(id: string) { return ROLE_GUIDES.find(guide=>guide.id===id); }
export function sectionsForRole(id: string, permissions: readonly string[], config: BusinessConfig): Section[] {
  const guide=roleGuide(id);if(!guide)return [];
  return visibleSections(permissions,config).filter(section=>(guide.sections as readonly string[]).includes(section.id));
}
