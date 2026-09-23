import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { listModels } from "./public-client";
import { createServer, type RequestListener } from "node:http";
import { gzipSync } from "node:zlib";

const limit = 1_048_576;
const utf8 = new TextEncoder();
beforeEach(() => {
  vi.stubEnv("ENTERPRISE_API_BASE_URL", "https://api.example.test");
  vi.stubEnv("ENTERPRISE_TENANT_CODE", "tenant");
  vi.stubEnv("ENTERPRISE_ORGANIZATION_CODE", "store");
});

async function withHTTP(handler: RequestListener, run: () => Promise<void>) {
  const server = createServer(handler);
  await new Promise<void>((resolve, reject) => { server.once("error", reject); server.listen(0, "127.0.0.1", resolve); });
  try {
    const address = server.address();
    if (!address || typeof address === "string") throw new Error("loopback address missing");
    vi.stubEnv("ENTERPRISE_API_BASE_URL", `http://127.0.0.1:${address.port}`);
    await run();
  } finally {
    server.closeAllConnections();
    await new Promise<void>((resolve, reject) => server.close((error) => error ? reject(error) : resolve()));
  }
}

describe("native fetch against a real loopback HTTP server", () => {
  it("limits decompressed bytes rather than trusting compressed Content-Length", async () => {
    const compressed = gzipSync(JSON.stringify({ models: [], padding: "é".repeat(limit / 2) }));
    expect(compressed.byteLength).toBeLessThan(limit);
    let requests = 0;
    await withHTTP((_request, response) => {
      requests++;
      response.writeHead(200, { "content-type": "application/json", "content-encoding": "gzip", "content-length": compressed.byteLength });
      response.end(compressed);
    }, async () => {
      await expect(listModels()).rejects.toMatchObject({ code: "RESPONSE_TOO_LARGE" });
      expect(requests).toBe(1);
    });
  });
  it("cancels an oversized chunked response before the server finishes it", async () => {
    let closed!: () => void;
    const closure = new Promise<void>((resolve) => { closed = resolve; });
    await withHTTP((_request, response) => {
      response.on("close", closed);
      response.writeHead(200, { "content-type": "application/json" });
      response.write(" ".repeat(limit + 1));
      // No end: the client must reject and cancel rather than await EOF.
    }, async () => {
      await expect(listModels()).rejects.toMatchObject({ code: "RESPONSE_TOO_LARGE" });
      await closure;
    });
  }, 2000);
  it("retains the five-second deadline while waiting for body bytes", async () => {
    let requests = 0;
    await withHTTP((_request, response) => {
      requests++;
      response.writeHead(200, { "content-type": "application/json" });
      response.write('{"models":[');
    }, async () => {
      await expect(listModels()).rejects.toMatchObject({ code: "INVALID_RESPONSE" });
      expect(requests).toBe(1);
    });
  }, 7000);
});
afterEach(() => { vi.unstubAllGlobals(); vi.unstubAllEnvs(); });

function reply(chunks: Uint8Array[], contentType = "application/json", status = 200, cancel = vi.fn()) {
  let count = 0;
  const stream = new ReadableStream<Uint8Array>({
    pull(controller) { const chunk = chunks[count++]; if (chunk) controller.enqueue(chunk); else controller.close(); },
    cancel
  }, { highWaterMark: 0 });
  const fetch = vi.fn(async () => new Response(stream, { status, headers: { "content-type": contentType, "content-length": "1" } }));
  vi.stubGlobal("fetch", fetch);
  return { stream, cancel, fetch, pulls: () => count };
}

describe("bounded backend response transport", () => {
  it("rejects UTF-8 bytes above budget even when the string is shorter and Content-Length lies", async () => {
    reply([utf8.encode(JSON.stringify({ models: [], padding: "é".repeat(limit / 2) }))]);
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "RESPONSE_TOO_LARGE" });
  });
  it("cancels at the first excessive chunk without draining the rest or retrying", async () => {
    const fixture = reply([utf8.encode(" ".repeat(limit)), utf8.encode(" "), utf8.encode('{"models":[]}')]);
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "RESPONSE_TOO_LARGE" });
    expect(fixture.pulls()).toBe(2);
    expect(fixture.cancel).toHaveBeenCalledOnce();
    expect(fixture.fetch).toHaveBeenCalledOnce();
    expect(fixture.stream.locked).toBe(false);
  });
  it("accepts exactly the byte budget", async () => {
    const json = '{"models":[]}';
    reply([utf8.encode(json + " ".repeat(limit - utf8.encode(json).byteLength))]);
    await expect(listModels()).resolves.toEqual([]);
  });
  it("preserves Unicode split across byte chunks", async () => {
    const bytes = utf8.encode('{"models":[{"displayName":"Eléctrica 🚲"}]}');
    reply(Array.from(bytes, (byte) => new Uint8Array([byte])));
    await expect(listModels()).resolves.toEqual([{ displayName: "Eléctrica 🚲" }]);
  });
  it("cancels invalid media types without reading the body", async () => {
    const fixture = reply([utf8.encode("private body")], "text/html");
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "INVALID_CONTENT_TYPE" });
    expect(fixture.pulls()).toBe(0);
    expect(fixture.cancel).toHaveBeenCalledOnce();
  });
  it("does not wait for a stalled cancellation promise", async () => {
    const fixture = reply([utf8.encode(" ".repeat(limit + 1)), utf8.encode("{}")], "application/json", 200, vi.fn(() => new Promise<void>(() => {})));
    await expect(listModels()).rejects.toMatchObject({ code: "RESPONSE_TOO_LARGE" });
    expect(fixture.cancel).toHaveBeenCalledOnce();
    expect(fixture.stream.locked).toBe(false);
  }, 1000);
  for (const bytes of [utf8.encode('{"secret":"private-token",'), new Uint8Array([123, 34, 120, 34, 58, 34, 255, 34, 125])]) {
    it(`rejects malformed JSON/UTF-8 without propagating payload (${bytes.byteLength})`, async () => {
      reply([bytes]);
      await expect(listModels()).rejects.toMatchObject({ status: 502, code: "INVALID_RESPONSE", message: "backend request failed: INVALID_RESPONSE" });
    });
  }
  it("does not manufacture a successful response after a stream failure", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => new Response(new ReadableStream({ start(controller) { controller.error(new Error("private-token")); } }), { headers: { "content-type": "application/json" } })));
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "INVALID_RESPONSE", message: "backend request failed: INVALID_RESPONSE" });
  });
  for (const body of ["null", '{"code":{"secret":"private-token"}}']) {
    it(`preserves error status but rejects unusable upstream codes (${body})`, async () => {
      reply([utf8.encode(body)], "application/problem+json", 403);
      await expect(listModels()).rejects.toMatchObject({ status: 403, code: "UPSTREAM_ERROR" });
    });
  }
});
