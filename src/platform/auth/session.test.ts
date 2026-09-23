import { afterEach, describe, expect, it } from "vitest";
import { allowed, openFlow, openSession, sealFlow, sealSession } from "./session";

const original = process.env.AUTH_SESSION_SECRET;
afterEach(() => { process.env.AUTH_SESSION_SECRET = original; });

describe("encrypted portal session", () => {
  it("round trips bounded flow and session values", async () => {
    process.env.AUTH_SESSION_SECRET = "a-secure-test-secret-with-more-than-thirty-two-characters";
    const flow = { codeVerifier: "verifier", state: "state", nonce: "nonce", returnTo: "/admin" };
    await expect(openFlow(await sealFlow(flow))).resolves.toEqual(flow);
    const session = { subject: "user", tenantId: "tenant", permissions: ["admin:read"], organizations: ["org"], accessToken: "access" };
    const opened = await openSession(await sealSession(session, 60));
    expect(opened).toEqual(session);
    expect(allowed(opened, "admin:read")).toBe(true);
    expect(allowed(opened, "factory:read")).toBe(false);
  });

  it("rejects tampering and weak secrets", async () => {
    process.env.AUTH_SESSION_SECRET = "short";
    await expect(sealFlow({ codeVerifier: "v", state: "s", nonce: "n", returnTo: "/" })).rejects.toThrow(/32/);
    process.env.AUTH_SESSION_SECRET = "a-secure-test-secret-with-more-than-thirty-two-characters";
    const token = await sealFlow({ codeVerifier: "v", state: "s", nonce: "n", returnTo: "/" });
    const position = Math.floor(token.length / 2);
    const tampered = token.slice(0, position) + (token[position] === "a" ? "b" : "a") + token.slice(position + 1);
    await expect(openFlow(tampered)).rejects.toThrow();
  });

  it("rejects sessions that cannot fit safely in one cookie", async () => {
    process.env.AUTH_SESSION_SECRET = "a-secure-test-secret-with-more-than-thirty-two-characters";
    const session = { subject: "user", tenantId: "tenant", permissions: ["admin:read"], organizations: ["org"], accessToken: "x".repeat(2_001) };
    await expect(sealSession(session, 60)).rejects.toThrow(/access token size/);
  });
});
