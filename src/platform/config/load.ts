import { readFile } from "node:fs/promises";
import { basename, join } from "node:path";
import { businessConfigSchema, type BusinessConfig } from "./schema";

let cached: BusinessConfig | undefined;

export async function loadBusinessConfig(options?: { path?: string; bypassCache?: boolean }): Promise<BusinessConfig> {
  if (cached && !options?.bypassCache) return cached;

  const configFile = basename(options?.path ?? process.env.BUSINESS_CONFIG_FILE ?? "business.example.json");
  const configPath = join(process.cwd(), "config", configFile);
  const text = await readFile(configPath, "utf8");
  const raw = JSON.parse(text.charCodeAt(0) === 0xfeff ? text.slice(1) : text) as unknown;
  const parsed = businessConfigSchema.parse(raw);
  if (!options?.bypassCache) cached = parsed;
  return parsed;
}

export function clearBusinessConfigCache(): void {
  cached = undefined;
}
