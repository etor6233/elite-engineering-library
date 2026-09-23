import{parseHelpQuery}from"./query";
export{parseHelpQuery}from"./query";
import "server-only";
import { allowed, type PortalSession } from "@/platform/auth/session";
import { ALL_GUIDES, type GuideID } from "./content";

// Match the existing page/section permission predicates. Reading is not command authorization.
const operationPermissions = ["inventory:allocate", "payment:create", "handover:manage", "admin:read"];
const sectionPermissions = [...operationPermissions, "lead:read", "resource:manage", "availability:read", "availability:manage", "appointment:manage"];
export function availableGuides(session: PortalSession, notificationsEnabled: boolean) {
  const can = (permission: string) => allowed(session, permission);
  const policy: Record<GuideID, boolean> = {
    "catalog-role-view": can("catalog:read"),
    "supply-role-view": can("supply:read") || can("supply:factory-read"),
    "warranty-role-view": can("warranty:read") || can("warranty:self") || can("warranty:factory-read"),
    "network-role-view": can("network:admin") || can("franchise:write"),
    "help-cms-view": can("help:read") || can("help:write") || can("help:publish"),
    "training-role-view": can("training:learn") || can("training:review"),

    "quote-acceptance-view": can("customer:self"),
    "handover-read-view": can("customer:self"),
    "whatsapp-status-view": notificationsEnabled && can("appointment:manage"),
    "order-operations-view": operationPermissions.some(can),
    "operation-sections-view": sectionPermissions.some(can),
    "lead-command-view": can("lead:read"),
    "quote-create-view": can("lead:read") && can("quote:write"),
    "availability-cancel-view": can("availability:read"),
    "availability-create-view": can("availability:manage"),
    "resource-create-view": can("resource:manage"),
    "slot-create-view": can("appointment:manage"),
    "delivery-resolution-view": can("handover:manage"),
    "checklist-publication-view": can("handover:manage"),
    "checklist-completion-view": can("handover:manage"),
    "return-operations-view": can("handover:manage")
  };
  // Project only the response contract, excluding inline-rendering metadata.
  return ALL_GUIDES.filter(g => policy[g.id]).map(({ id, version, title, paragraphs }) => ({ id, version, title, paragraphs }));
}
export function selectGuides(guides: ReturnType<typeof availableGuides>, query: NonNullable<ReturnType<typeof parseHelpQuery>>) {
  if (query.article !== null) return guides.filter(g => g.id === query.article && g.version === query.version);
  return guides.filter(g => [g.title, ...g.paragraphs].join("\n").normalize("NFC").toLocaleLowerCase("es").includes(query.q));
}
