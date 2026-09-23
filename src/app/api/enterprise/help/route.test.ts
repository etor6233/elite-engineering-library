import { beforeEach, expect, test, vi } from "vitest";
const mocks = vi.hoisted(() => ({ readSession: vi.fn(), loadBusinessConfig: vi.fn() }));
vi.mock("@/platform/auth/session", () => ({ readSession: mocks.readSession, allowed: (s: {permissions: string[]}, p: string) => s.permissions.includes(p) || s.permissions.includes("*") }));
vi.mock("@/platform/config/load", () => ({ loadBusinessConfig: mocks.loadBusinessConfig }));
import { GET } from "./route";
import { helpResponseSchema } from "@/platform/help/contract";
import { QUOTE_GUIDE, NOTIFICATION_GUIDE, OPERATIONAL_GUIDES, ALL_GUIDES } from "@/platform/help/content";
import { readFileSync } from "node:fs";
const session = (permissions: string[]) => ({ subject: "private-subject", tenantId: "private-tenant", permissions, organizations: ["private-org"], accessToken: "never-return-token" });
const get = (query = "", headers = {}) => GET(new Request("https://example.test/api/enterprise/help" + query, { headers }));
beforeEach(() => { vi.resetAllMocks(); mocks.readSession.mockResolvedValue(session(["customer:self"])); mocks.loadBusinessConfig.mockResolvedValue({ features: { whatsapp_status_history: true } }); });
test("customer receives only applicable versioned public help without identity", async () => {
 const r = await get(); expect(r.status).toBe(200); expect(r.headers.get("cache-control")).toBe("no-store"); expect(r.headers.get("x-content-type-options")).toBe("nosniff");
 const data = helpResponseSchema.parse(await r.json()); expect(data.articles.map(g=>g.id)).toEqual(["quote-acceptance-view","handover-read-view"]); expect(JSON.stringify(data)).not.toMatch(/private-|never-return-token/);
});
test("operator receives only enabled notification guidance", async () => { mocks.readSession.mockResolvedValue(session(["appointment:manage"])); expect((await (await get()).json()).articles.map((g:{id:string})=>g.id).sort()).toEqual(["operation-sections-view","slot-create-view","whatsapp-status-view"]); mocks.loadBusinessConfig.mockResolvedValue({ features: {} }); expect((await (await get()).json()).articles.map((g:{id:string})=>g.id).sort()).toEqual(["operation-sections-view","slot-create-view"]); });
test("wildcard still obeys feature selection", async () => { mocks.readSession.mockResolvedValue(session(["*"])); expect((await (await get()).json()).articles).toHaveLength(21); mocks.loadBusinessConfig.mockResolvedValue({ features: {} }); expect((await (await get()).json()).articles.map((g:{id:string})=>g.id).sort()).toEqual(ALL_GUIDES.filter(g=>g.id!=="whatsapp-status-view").map(g=>g.id).sort()); });
test("no permission is denied and rechecked on subsequent reads", async () => { expect((await get()).status).toBe(200); mocks.readSession.mockResolvedValue(session(["factory:read"])); const r=await get(); expect(r.status).toBe(403); expect(await r.json()).toEqual({ code: "FORBIDDEN" }); });
test("missing session returns401 without loading configuration", async () => { mocks.readSession.mockResolvedValue(null); expect((await get()).status).toBe(401); expect(mocks.loadBusinessConfig).not.toHaveBeenCalled(); });
test.each(["cross-site", "none"])("rejects %s before credentials", async value => { expect((await get("", { "sec-fetch-site": value })).status).toBe(403); expect(mocks.readSession).not.toHaveBeenCalled(); });
test("exact version succeeds and unavailable version never silently upgrades", async () => { expect((await get("?article=quote-acceptance-view&version=1.0.0")).status).toBe(200); for(const q of ["?article=quote-acceptance-view&version=9.0.0", "?article=whatsapp-status-view&version=1.0.0", "?article=unknown&version=1.0.0"]) { const r=await get(q); expect(r.status).toBe(404); expect(await r.json()).toEqual({code:"GUIDE_NOT_AVAILABLE"}); } });
test("literal search is case/NFC insensitive and does not interpret regex or markup", async () => { for(const q of ["COTIZACIÓN", "cotizacio\u0301n"]) expect((await (await get("?q="+encodeURIComponent(q))).json()).articles).toHaveLength(1); for(const q of [".*", "<script>alert(1)</script>", "unmatched"]) expect((await (await get("?q="+encodeURIComponent(q))).json()).articles).toHaveLength(0); });
test.each(["?q=a&q=b", "?tenant=other", "?role=*", "?article=quote-acceptance-view", "?version=1.0.0", "?article=../escape&version=1.0.0", "?article=quote-acceptance-view&version=latest", "?article=quote-acceptance-view&version=1.0.0&q=", "?q=%00", "?q="+"x".repeat(161), "?article="+"x".repeat(2050)])("invalid query rejects before credentials: %s", async query => { expect((await get(query)).status).toBe(400); expect(mocks.readSession).not.toHaveBeenCalled(); });
test("unexpected session/config failures are sanitized and not cached", async () => { mocks.readSession.mockRejectedValue(new Error("private-secret")); let r=await get(); expect(r.status).toBe(503); expect(await r.text()).not.toContain("private-secret"); mocks.readSession.mockResolvedValue(session(["*"])); mocks.loadBusinessConfig.mockRejectedValue(new Error("config-secret")); r=await get(); expect(r.status).toBe(503); expect(r.headers.get("cache-control")).toBe("no-store"); expect(await r.json()).toEqual({code:"HELP_UNAVAILABLE"}); });
test("inline guides share the source and retain their safety instructions", () => { const quote=readFileSync("src/components/customer-quote-actions.tsx","utf8"), notification=readFileSync("src/components/appointment-notification-status.tsx","utf8"); expect(quote).toContain("QUOTE_GUIDE.paragraphs.map"); expect(notification).toContain("NOTIFICATION_GUIDE.paragraphs[1]"); expect(QUOTE_GUIDE.paragraphs.join(" ")).toContain("no confirma el pago, el stock ni la entrega"); expect(NOTIFICATION_GUIDE.paragraphs.join(" ")).toContain("no reenvíes ni borres el historial"); });

const roleMatrix = [
  {
    "role": "none",
    "permissions": [],
    "expected": []
  },
  {
    "role": "factory",
    "permissions": [
      "factory:read"
    ],
    "expected": []
  },
  {
    "role": "customer",
    "permissions": [
      "customer:self"
    ],
    "expected": [
      "handover-read-view",
      "quote-acceptance-view"
    ]
  },
  {
    "role": "lead",
    "permissions": [
      "lead:read"
    ],
    "expected": [
      "lead-command-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "quote-without-lead-read",
    "permissions": [
      "quote:write"
    ],
    "expected": []
  },
  {
    "role": "quote",
    "permissions": [
      "lead:read",
      "quote:write"
    ],
    "expected": [
      "lead-command-view",
      "operation-sections-view",
      "quote-create-view"
    ]
  },
  {
    "role": "availability-read",
    "permissions": [
      "availability:read"
    ],
    "expected": [
      "availability-cancel-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "availability-manage",
    "permissions": [
      "availability:manage"
    ],
    "expected": [
      "availability-create-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "resource",
    "permissions": [
      "resource:manage"
    ],
    "expected": [
      "operation-sections-view",
      "resource-create-view"
    ]
  },
  {
    "role": "appointment",
    "permissions": [
      "appointment:manage"
    ],
    "expected": [
      "operation-sections-view",
      "slot-create-view",
      "whatsapp-status-view"
    ]
  },
  {
    "role": "inventory",
    "permissions": [
      "inventory:allocate"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "payment",
    "permissions": [
      "payment:create"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "handover",
    "permissions": [
      "handover:manage"
    ],
    "expected": [
      "checklist-completion-view",
      "checklist-publication-view",
      "delivery-resolution-view",
      "operation-sections-view",
      "order-operations-view",
      "return-operations-view"
    ]
  },
  {
    "role": "admin",
    "permissions": [
      "admin:read"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "wildcard",
    "permissions": [
      "*"
    ],
    "expected": [
      "availability-cancel-view",
      "availability-create-view",
      "catalog-role-view",
      "checklist-completion-view",
      "checklist-publication-view",
      "delivery-resolution-view",
      "handover-read-view",
      "help-cms-view",
      "lead-command-view",
      "network-role-view",
      "operation-sections-view",
      "order-operations-view",
      "quote-acceptance-view",
      "quote-create-view",
      "resource-create-view",
      "return-operations-view",
      "slot-create-view",
      "supply-role-view",
      "training-role-view",
      "warranty-role-view",
      "whatsapp-status-view"
]
  }
] as const;
test.each(roleMatrix)("complete guide access for $role", async row => {
 mocks.readSession.mockResolvedValue(session([...row.permissions]));
 const r=await get();expect(r.status).toBe(row.expected.length?200:403);
 if(row.expected.length){const data=helpResponseSchema.parse(await r.json());expect(data.articles.map(g=>g.id).sort()).toEqual([...row.expected]);}
 for(const guide of ALL_GUIDES){const exact=await get(`?article=${guide.id}&version=${guide.version}`);const included=(row.expected as readonly string[]).includes(guide.id);expect(exact.status).toBe(included?200:row.expected.length?404:403);}
});
