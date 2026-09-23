import { describe, expect, it, vi } from "vitest";
import {
  htmlEscape,
  setAnchorHref,
  unwrapHtml,
} from "./google-safevalues";

describe("Google SafeValues DOM-XSS primitives", () => {
  it("escapes attacker-controlled HTML instead of creating markup", () => {
    const escaped = unwrapHtml(htmlEscape('<img src=x onerror="alert(1)">'));
    expect(String(escaped)).toBe("&lt;img src=x onerror=&quot;alert(1)&quot;&gt;");
  });

  it("refuses a javascript URL without overwriting the current anchor", () => {
    const error = vi.spyOn(console, "error").mockImplementation(() => undefined);
    const anchor = { href: "https://example.invalid/safe" } as HTMLAnchorElement;
    setAnchorHref(anchor, "javascript:alert(document.domain)");
    expect(anchor.href).toBe("https://example.invalid/safe");
    expect(error).toHaveBeenCalledOnce();
    error.mockRestore();
  });

  it("preserves an ordinary HTTPS URL", () => {
    const anchor = { href: "" } as HTMLAnchorElement;
    setAnchorHref(anchor, "https://example.invalid/catalog?model=e-bike");
    expect(anchor.href).toBe("https://example.invalid/catalog?model=e-bike");
  });
});
