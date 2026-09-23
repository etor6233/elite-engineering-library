import { afterEach, describe, expect, it } from "vitest";
import { clearBusinessConfigCache, loadBusinessConfig } from "./load";
import { businessConfigSchema } from "./schema";

afterEach(() => clearBusinessConfigCache());

describe("business presentation configuration", () => {
  it("loads the bounded config file from ./config", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    expect(config.business.id).toBe("electric-mobility-network");
    expect(config.workflows.lead?.states).toEqual(["new", "contacted", "qualified", "converted", "lost"]);
    expect(config.roles.find((role) => role.id === "customer")?.permissions).toEqual(["customer:self"]);
  });

  it("keeps the example lead workflow inside the PostgreSQL state contract", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    const invalid = structuredClone(config);
    invalid.workflows.lead!.states.push("assigned");
    expect(invalid.workflows.lead!.states).not.toEqual(config.workflows.lead!.states);
    expect(config.workflows.lead!.states).not.toContain("assigned");
  });

  it("rejects an invalid module dependency", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    const invalid = structuredClone(config);
    invalid.modules.orders = { enabled: false };
    expect(() => businessConfigSchema.parse(invalid)).toThrow(/requires enabled module orders/);
  });

  it("does not allow config paths to escape ./config", async () => {
    await expect(loadBusinessConfig({ path: "../outside.json", bypassCache: true })).rejects.toThrow();
  });
});
