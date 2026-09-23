import { loadBusinessConfig } from "@/platform/config/load";
import { resolvePublicLocale } from "./public-catalog";

export async function loadPublicLocale() {
  const config = await loadBusinessConfig();
  const market = config.markets.find((item) => item.code === config.business.defaultMarket);
  if (!market) throw new Error("Configured market is missing");
  return resolvePublicLocale(config.business.defaultLocale, market.timeZone);
}

