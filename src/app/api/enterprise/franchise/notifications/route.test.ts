import { beforeEach, describe, expect, it, vi } from "vitest";
import { GET } from "./route";
import { BackendProblem } from "@/platform/backend/public-client";
const mocks = vi.hoisted(() => ({ readSession: vi.fn(), protectedGet: vi.fn(), loadBusinessConfig: vi.fn() }));
vi.mock("@/platform/config/load", () => ({ loadBusinessConfig: mocks.loadBusinessConfig }));
vi.mock("@/platform/auth/session", () => ({ readSession: mocks.readSession, allowed: (s: { permissions: string[] }, p: string) => s.permissions.includes(p) }));
vi.mock("@/platform/backend/protected-client", () => ({ protectedGet: mocks.protectedGet }));
const session = { subject: "operator", tenantId: "tenant", permissions: ["appointment:manage"], organizations: ["store"], accessToken: "server-secret" };
const empty = { organization_id: "store", appointment_id: "appointment", items: [] };
const request = (query = "organizationId=store&appointmentId=appointment", site = "same-origin") => new Request(`https://portal.example/api/enterprise/franchise/notifications?${query}`, { headers: { "sec-fetch-site": site } });
beforeEach(() => { vi.clearAllMocks(); mocks.readSession.mockResolvedValue(session); mocks.protectedGet.mockResolvedValue(empty); mocks.loadBusinessConfig.mockResolvedValue({ features: { whatsapp_status_history: true } }); });
describe("notification history BFF", () => {
  it.each([{}, { whatsapp_status_history: false }])("does not activate an unconfigured capability", async features => {
    mocks.loadBusinessConfig.mockResolvedValue({ features });
    const response = await GET(request()); expect(response.status).toBe(404);
    expect(await response.json()).toEqual({ code: "CAPABILITY_NOT_ENABLED" });
    expect(mocks.readSession).not.toHaveBeenCalled(); expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("fails closed when configuration is unavailable", async () => {
    mocks.loadBusinessConfig.mockRejectedValue(new Error("private config path"));
    const response = await GET(request()); expect(response.status).toBe(503);
    expect(await response.text()).not.toContain("private"); expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("forwards only authenticated scope and returns no token", async () => {
    const response = await GET(request()); expect(response.status).toBe(200); expect(await response.json()).toEqual(empty); expect(response.headers.get("cache-control")).toBe("no-store");
    expect(mocks.protectedGet).toHaveBeenCalledWith(session, "/v1/franchise/appointments/appointment/whatsapp-confirmations", { organization_id: "store" });
  });
  it("rejects missing session, permissions and foreign organization", async () => {
    mocks.readSession.mockResolvedValue(null); expect((await GET(request())).status).toBe(401);
    mocks.readSession.mockResolvedValue({ ...session, permissions: [] }); expect((await GET(request())).status).toBe(403);
    mocks.readSession.mockResolvedValue(session); expect((await GET(request("organizationId=other&appointmentId=appointment"))).status).toBe(403);
    expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("rejects cross-site before resolving credentials", async () => {
    expect((await GET(request(undefined, "cross-site"))).status).toBe(403); expect(mocks.readSession).not.toHaveBeenCalled();
  });
  it.each(["organizationId=store&appointmentId=a&tenantId=other", "organizationId=store&appointmentId=a&appointmentId=b", "organizationId=store", "organizationId=store&appointmentId=..", "organizationId=store&appointmentId=a%2Fb", "organizationId=store&appointmentId=" + "a".repeat(129)])("rejects ambiguous/unbounded query %s", async query => { expect((await GET(request(query))).status).toBe(400); expect(mocks.protectedGet).not.toHaveBeenCalled(); });
  it.each([{ ...empty, organization_id: "other" }, { ...empty, appointment_id: "other" }, { ...empty, recipient: "secret" }, { ...empty, items: [{}] }])("rejects foreign or malformed backend projection", async raw => { mocks.protectedGet.mockResolvedValue(raw); const response = await GET(request()); expect(response.status).toBe(502); expect(await response.text()).not.toContain("secret"); });
  it("preserves approved backend errors without arbitrary details", async () => {
    mocks.protectedGet.mockRejectedValue(new BackendProblem(409, "NOTIFICATION_HISTORY_LIMIT")); const response = await GET(request()); expect(response.status).toBe(409); expect(await response.json()).toEqual({ code: "NOTIFICATION_HISTORY_LIMIT" });
    mocks.protectedGet.mockRejectedValue(new Error("private payload")); expect(await (await GET(request())).text()).not.toContain("private");
  });
});
