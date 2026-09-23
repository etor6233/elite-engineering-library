import { describe, expect, it } from "vitest";
import { z } from "zod";
import { BackendProblem } from "@/platform/backend/public-client";
import { errorResponse } from "./problem";

describe("BFF problem responses", () => {
  it("does not expose validation input", async () => {
    let failure: unknown;
    try { z.object({ email: z.email() }).parse({ email: "secret-invalid" }); } catch (error) { failure = error; }
    const response = errorResponse(failure);
    expect(response.status).toBe(400);
    expect(await response.text()).not.toContain("secret-invalid");
  });

  it("maps upstream server failures to a bounded gateway response", async () => {
    const response = errorResponse(new BackendProblem(503, "UPSTREAM_OVERLOAD"));
    expect(response.status).toBe(502);
    expect(await response.json()).toMatchObject({ code: "UPSTREAM_OVERLOAD" });
  });
});
