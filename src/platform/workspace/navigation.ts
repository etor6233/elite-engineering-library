// AUTHORED server-side selection contract; navigation never grants an action.
import type {BusinessConfig} from "@/platform/config/schema";
import {WORKSPACE_TASKS,WORKSPACE_GROUPS,INSTALLED_WORKSPACE_ROUTES} from "./task-registry.generated";
export type WorkspaceEntry={id:string;href:string;label:string;icon:string;taskLabel?:string};
export type WorkspaceGroup={id:string;label:string;items:WorkspaceEntry[]};
export type WorkspaceIdentity={tenantId:string;organizations:readonly string[];permissions:readonly string[]};
export function createWorkspaceNavigation(config:Pick<BusinessConfig,"features"|"modules">,identity:WorkspaceIdentity|null,options:{installedRoutes?:readonly string[];documentsVisible?:boolean;campaignsVisible?:boolean;aiActivityVisible?:boolean}={}):WorkspaceGroup[]{
 if(!identity?.tenantId||!identity.organizations.length)return [];
 const installed=new Set(options.installedRoutes??INSTALLED_WORKSPACE_ROUTES);
 const can=(permission:string)=>identity.permissions.includes("*")||identity.permissions.includes(permission);
 const tasks=WORKSPACE_TASKS.filter(task=>installed.has(task.href.split("?")[0]!)&&task.features.every(flag=>config.features[flag]===true)&&task.modules.every(module=>config.modules[module]?.enabled===true)&&(!task.anyPermissions.length||task.anyPermissions.some(can))&&task.allPermissions.every(can));
 const groups:WorkspaceGroup[]=WORKSPACE_GROUPS.map(group=>({...group,items:tasks.filter(task=>task.group===group.id).map(({id,href,label,icon,taskLabel})=>({id,href,label,icon,...(taskLabel?{taskLabel}:{})}))}));
 // The document owner computes visibility from its tenant/org/actor contract on the server.
 // Installation remains independent. An environment flag alone cannot install this route.
 if(installed.has("/documents")&&options.documentsVisible===true&&["documents:read","documents:write","documents:review","documents:process"].some(can))groups.find(g=>g.id==="documents")!.items.push({id:"documents:inbox",href:"/documents",label:"Bandeja de documentos",icon:"/documents"});
 if(installed.has("/franchise/campaigns")&&options.campaignsVisible===true&&can("marketing:read")&&can("notification:read"))groups.find(g=>g.id==="commercial")!.items.push({id:"campaigns:workspace",href:"/franchise/campaigns",label:"Campañas",icon:"Mensajes"});
 if(installed.has("/intelligence/activity")&&options.aiActivityVisible===true&&config.features.ai_activity===true&&can("whatsapp:approve"))groups.find(g=>g.id==="analysis")!.items.push({id:"intelligence:activity",href:"/intelligence/activity",label:"Actividad de IA",icon:"Análisis"});
 return groups.filter(group=>group.items.length);
}
export function isUnifiedWorkspaceRoute(pathname:string){return ["/franchise","/factory","/supply","/warranty","/network","/admin","/customer","/dashboard","/guide","/documents","/intelligence"].some(root=>pathname===root||pathname.startsWith(root+"/"))||pathname==="/help"||pathname.startsWith("/help/");}
export const FRANCHISE_TASK_QUERY:Readonly<Record<string,string>>={"Pedidos":"orders","Entregas":"deliveries","Agenda":"agenda","Gestión diaria":"daily"};
export function workspaceEntryActive(entry:WorkspaceEntry,pathname:string,activeTask?:string){
 const path=entry.href.split("?")[0]!;
 if(entry.taskLabel)return pathname==="/franchise"&&entry.taskLabel===activeTask;
 // Exact matching avoids two current entries, e.g. /admin and /admin/catalog.
 return pathname===path||(pathname.startsWith(path+"/")&&!WORKSPACE_TASKS.some(other=>other.href.split("?")[0]===pathname));
}
