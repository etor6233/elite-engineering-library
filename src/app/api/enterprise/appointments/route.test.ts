import { beforeEach, describe, expect, it, vi } from "vitest";
import { POST } from "./route";

const { requestAppointment } = vi.hoisted(() => ({ requestAppointment: vi.fn() }));
vi.mock("@/platform/backend/public-client", () => ({ requestAppointment }));

beforeEach(() => requestAppointment.mockReset());

describe("appointment BFF", () => {
  it("rejects requests without an idempotency key", async () => {
    const response = await POST(new Request("https://web.example.test/api/enterprise/appointments", { method: "POST", headers: { "content-type": "application/json" }, body: "{}" }));
    expect(response.status).toBe(400);
    expect(requestAppointment).not.toHaveBeenCalled();
  });

  it("validates and forwards an exact appointment", async () => {
    requestAppointment.mockResolvedValue({ value: { id: "appointment", state: "requested" }, status: 202, replayed: false });
    const response = await POST(new Request("https://web.example.test/api/enterprise/appointments", { method: "POST", headers: { "content-type": "application/json", "idempotency-key": "appointment-key-0001" }, body: JSON.stringify({ leadId: "lead", modelId: "model", kind: "test-drive", startsAt: "2026-09-01T15:00:00Z" }) }));
    expect(response.status).toBe(202);
    expect(requestAppointment).toHaveBeenCalledWith({ leadId: "lead", modelId: "model", kind: "test-drive", startsAt: "2026-09-01T15:00:00Z" }, "appointment-key-0001");
  });
});
