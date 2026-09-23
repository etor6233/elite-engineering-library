import { beforeEach, describe, expect, it, test } from "vitest";
import { readFileSync } from "node:fs";
import { businessConfigSchema, type BusinessConfig } from "@/platform/config/schema";
import { enabledNavigation } from "@/platform/config/registry";
import { ROLE_GUIDES, roleGuide, sectionsForRole, visibleSections } from "./role-visibility";
let config:BusinessConfig;
beforeEach(()=>{config=businessConfigSchema.parse(JSON.parse(readFileSync("config/business.example.json","utf8")));Object.assign(config.features,{role_workspace:true,customer_portal:true,factory_portal:true,public_catalog:true});for(const module of ["catalog","crm","procurement"])config.modules[module]!.enabled=true;});
const common=["dashboard","help","locations","public_catalog"];
const matrix=[
  {
    "name": "none",
    "permissions": [],
    "privateIds": []
  },
  {
    "name": "customer",
    "permissions": [
      "customer:self"
    ],
    "privateIds": [
      "customer"
    ]
  },
  {
    "name": "factory",
    "permissions": [
      "factory:read"
    ],
    "privateIds": [
      "factory"
    ]
  },
  {
    "name": "admin",
    "permissions": [
      "admin:read"
    ],
    "privateIds": [
      "admin",
      "franchise"
    ]
  },
  {
    "name": "lead",
    "permissions": [
      "lead:read"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "inventory",
    "permissions": [
      "inventory:allocate"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "payment",
    "permissions": [
      "payment:create"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "handover",
    "permissions": [
      "handover:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "resource",
    "permissions": [
      "resource:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "availability-read",
    "permissions": [
      "availability:read"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "availability-manage",
    "permissions": [
      "availability:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "appointment",
    "permissions": [
      "appointment:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "quote-only",
    "permissions": [
      "quote:write"
    ],
    "privateIds": []
  },
  {
    "name": "wrong-case",
    "permissions": [
      "Admin:Read",
      "owner",
      "lead:reader",
      "* "
    ],
    "privateIds": []
  },
  {
    "name": "wildcard",
    "permissions": [
      "*"
    ],
    "privateIds": [
      "admin",
      "customer",
      "factory",
      "franchise"
    ]
  }
] as const;
test.each(matrix)("header and workspace use actual permissions: $name",row=>{
 const nav=enabledNavigation(config,row.permissions);expect(nav.map(s=>s.id).sort()).toEqual([...common,...row.privateIds].sort());
 const map:Record<string,string>={dashboard:"overview",franchise:"sales",public_catalog:"catalog"};
 expect(visibleSections(row.permissions,config).map(s=>s.id).sort()).toEqual([...common,...row.privateIds].map(id=>map[id]??id).concat("guide").sort());
});
describe("common authenticated access and exact grants",()=>{
 it("retains common help and guide for non-wildcard sessions",()=>{for(const permissions of [[],["customer:self"],["lead:read"]])expect(visibleSections(permissions,config).map(s=>s.id)).toEqual(expect.arrayContaining(["overview","help","guide"]));});
 it("preserves the customer login entry only for guests",()=>{expect(enabledNavigation(config,null).map(s=>s.id).sort()).toEqual([...common,"customer"].sort());expect(enabledNavigation(config,[]).some(s=>s.id==="customer")).toBe(false);});
 it("feature flags win over wildcard and the workspace stays opt-in",()=>{Object.assign(config.features,{role_workspace:false,customer_portal:false,factory_portal:false,public_catalog:false});const ids=visibleSections(["*"],config).map(s=>s.id);for(const id of ["overview","guide","customer","factory","catalog"])expect(ids).not.toContain(id);expect(ids).toContain("help");});
 it("disabled modules remove their destinations even for wildcard",()=>{for(const module of ["catalog","crm","procurement"])config.modules[module]!.enabled=false;expect(enabledNavigation(config,["*"]).map(s=>s.id)).not.toEqual(expect.arrayContaining(["factory","franchise","locations","public_catalog"]));for(const id of ["factory","franchise","locations","public_catalog"])expect(enabledNavigation(config,["*"]).map(s=>s.id)).not.toContain(id);});
 it("keeps the four labels without treating any label as a permission",()=>{expect(ROLE_GUIDES.map(g=>g.id)).toEqual(["owner","admin","employee","customer"]);for(const role of ROLE_GUIDES){expect(sectionsForRole(role.id,[role.id],config).map(s=>s.id)).not.toEqual(expect.arrayContaining(["sales","factory","admin","customer"]));}});
 it("every role view is a subset of the same effective-permission workspace",()=>{for(const row of matrix)for(const role of ROLE_GUIDES){const visible=new Set(visibleSections(row.permissions,config).map(s=>s.id));for(const section of sectionsForRole(role.id,row.permissions,config))expect(visible.has(section.id)).toBe(true);}});
 it("unknown and prototype role names never select a guide",()=>{for(const role of ["unknown","__proto__","constructor","OWNER","../owner"]){expect(roleGuide(role)).toBeUndefined();expect(sectionsForRole(role,["*"],config)).toEqual([]);}});
});
