import { describe, expect, it } from "vitest";
import { publicAppointmentTime, publicCount, publicMessage, resolvePublicLocale } from "./public-catalog";

describe("public journey localization", () => {
  it.each(["en-US", "en-GB", "EN-au"])("uses the English catalog for %s", (locale) => {
    expect(publicMessage(locale, "lead.name")).toBe("Name");
    expect(resolvePublicLocale(locale, "UTC").language).toBe("en");
  });
  it.each(["es-AR", "es-ES", "ES-mx"])("uses the Spanish catalog for %s", (locale) => {
    expect(publicMessage(locale, "lead.name")).toBe("Nombre");
    expect(resolvePublicLocale(locale, "UTC").language).toBe("es");
  });
  it.each(["fr-FR", "xx", "", "not_a_locale", "__proto__", "<script>"])("falls back honestly for %s", (locale) => {
    expect(resolvePublicLocale(locale, "UTC")).toEqual({ language: "es", locale: "es", timeZone: "UTC" });
    expect(publicMessage(locale, "lead.name")).toBe("Nombre");
  });
  it("returns unknown keys literally without reading prototypes", () => {
    expect(publicMessage("en-US", "missing")).toBe("missing");
    expect(publicMessage("en-US", "__proto__")).toBe("__proto__");
    expect(publicMessage("en-US", "constructor")).toBe("constructor");
  });
  it.each([
    ["en-US", 0, "0 models available"], ["en-US", 1, "1 model available"],
    ["en-US", 2, "2 models available"], ["es-AR", 0, "0 modelos disponibles"],
    ["es-AR", 1, "1 modelo disponible"], ["es-AR", 2, "2 modelos disponibles"]
  ])("formats %s count %s", (locale, count, expected) => {
    expect(publicCount(String(locale), "models", Number(count))).toBe(expected);
  });
  it.each([-1, 1.5, NaN, Infinity, Number.MAX_SAFE_INTEGER + 1])("rejects invalid count %s", (count) => {
    expect(() => publicCount("en", "models", count)).toThrow(RangeError);
  });
  it("keeps explicit configured time zones and crosses daylight saving", () => {
    const config = resolvePublicLocale("en-US", "America/New_York");
    const winter = publicAppointmentTime(config, "2026-01-15T15:00:00Z");
    const summer = publicAppointmentTime(config, "2026-07-15T15:00:00Z");
    expect(winter).toContain("10:00");
    expect(summer).toContain("11:00");
    expect(publicAppointmentTime(resolvePublicLocale("es-AR", "America/Argentina/Buenos_Aires"), "2026-01-15T15:00:00Z")).toContain("12:00");
  });
  it("rejects invalid dates and time zones without a silent fallback", () => {
    expect(() => resolvePublicLocale("en", "Mars/Olympus")).toThrow(RangeError);
    expect(() => publicAppointmentTime(resolvePublicLocale("en", "UTC"), "bad")).toThrow(RangeError);
  });
  it("does not share mutable locale state across concurrent renders", async () => {
    const values = await Promise.all(Array.from({ length: 100 }, (_, i) =>
      Promise.resolve(publicMessage(i % 2 ? "en-GB" : "es-AR", "lead.name"))));
    expect(values.filter((v) => v === "Name")).toHaveLength(50);
    expect(values.filter((v) => v === "Nombre")).toHaveLength(50);
  });
});

