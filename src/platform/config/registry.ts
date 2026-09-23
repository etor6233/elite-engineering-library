import type { BusinessConfig } from "./schema";

export interface NavigationItem {
  id: string;
  label: string;
  href: string;
  module?: string;
  feature?: string;
  permission?: string;
  anyPermissions?: readonly string[];
  guestEntry?: boolean;
}

const navigation: NavigationItem[] = [
  {id:"help_content",label:"Contenido de ayuda",href:"/help/library",feature:"help_cms",anyPermissions:["help:read","help:write","help:publish"]},
  {id:"network",label:"Red de franquicia",href:"/network",feature:"network_portal",anyPermissions:["network:admin","franchise:write"]},
  {id:"warranty",label:"GarantÃ­a",href:"/warranty",feature:"warranty_portal",anyPermissions:["warranty:read","warranty:self","warranty:factory-read"]},
  {id:"supply",label:"Suministro",href:"/supply",feature:"supply_portal",anyPermissions:["supply:read","supply:factory-read"]},
  { id: "dashboard", label: "Mi panel", href: "/dashboard", feature: "role_workspace" },
  { id: "help", label: "Ayuda", href: "/help" },
    { id:"catalog_editor",label:"Editar catálogo",href:"/admin/catalog",feature:"catalog_editor",permission:"catalog:read" },
  { id:"publishing",label:"Publicaciones",href:"/admin/publishing",feature:"social_publishing",permission:"social:read" },
  { id:"finance",label:"Finanzas",href:"/admin/finance",feature:"finance_workspace",anyPermissions:["accounting:read","royalty:read"] },

  { id: "training", label: "CapacitaciÃ³n", href: "/guide/training", feature: "training_portal", anyPermissions: ["training:learn","training:review"] },
  { id: "public_catalog", label: "Modelos", href: "/models", module: "catalog", feature: "public_catalog" },
  { id: "locations", label: "DÃ³nde estamos", href: "/locations", module: "crm" },
  { id: "customer", label: "Mi cuenta", href: "/customer", feature: "customer_portal", permission: "customer:self", guestEntry: true },
  { id: "admin", label: "OperaciÃ³n", href: "/admin", permission: "admin:read" },
  { id: "franchise", label: "Franquicia", href: "/franchise", module: "crm", anyPermissions: ["inventory:allocate", "payment:create", "handover:manage", "admin:read", "lead:read", "resource:manage", "availability:read", "availability:manage", "appointment:manage"] },
  { id: "factory", label: "FÃ¡brica", href: "/factory", module: "procurement", feature: "factory_portal", permission: "factory:read" }
];

export function enabledNavigation(config: BusinessConfig, permissions: readonly string[] | null = null): NavigationItem[] {
  return navigation.filter((item) => {
    if (item.module && !config.modules[item.module]?.enabled) return false;
    if (item.feature && !config.features[item.feature]) return false;
    const can = (permission: string) => permissions !== null && (permissions.includes("*") || permissions.includes(permission));
    if (item.permission && !(permissions === null && item.guestEntry) && !can(item.permission)) return false;
    if (item.anyPermissions && !item.anyPermissions.some(can)) return false;
    return true;
  });
}

export function enabledModules(config: BusinessConfig): string[] {
  return Object.entries(config.modules).filter(([, value]) => value.enabled).map(([name]) => name);
}
