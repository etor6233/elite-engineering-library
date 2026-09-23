// AUTHORED public-journey catalog. Runtime Intl supplies locale/date/plural semantics.
const es = {
  "nav.network": "Red de franquicia",
    "nav.help_content": "Contenido de ayuda",
  "site.description": "CatÃ¡logo y portales de una red de movilidad elÃ©ctrica.",
  "site.skip": "Saltar al contenido", "site.navigation": "Principal", "site.contact": "Contacto:",
  "nav.warranty": "GarantÃ­a", "nav.supply": "Suministro", "nav.catalog_editor": "Editar catÃ¡logo", "nav.training": "CapacitaciÃ³n", "nav.help": "Ayuda", "nav.dashboard": "Mi panel",
  "nav.public_catalog": "Modelos", "nav.locations": "DÃ³nde estamos", "nav.customer": "Mi cuenta",
  "nav.admin": "OperaciÃ³n", "nav.franchise": "Franquicia", "nav.factory": "FÃ¡brica", "nav.publishing": "Publicaciones", "nav.finance": "Finanzas",
  "home.eyebrow": "Movilidad elÃ©ctrica configurable",
  "home.description": "ExplorÃ¡ modelos publicados y enviÃ¡ una consulta con tu consentimiento.",
  "home.action": "Ver modelos", "models.title": "Modelos elÃ©ctricos", "models.short": "Modelos",
  "models.eyebrow": "CatÃ¡logo pÃºblico", "models.description": "ElegÃ­ un modelo para solicitar informaciÃ³n.",
  "models.one": "modelo disponible", "models.other": "modelos disponibles",
  "lead.name": "Nombre", "lead.email": "Email", "lead.consent": "Acepto ser contactado sobre este modelo.",
  "lead.sending": "Enviandoâ€¦", "lead.submit": "Solicitar informaciÃ³n",
  "lead.failed": "No pudimos registrar la solicitud.",
  "lead.invalid": "La respuesta no incluyÃ³ la referencia de la solicitud.",
  "lead.received": "Solicitud recibida. Ya podÃ©s reservar una consulta o prueba.",
  "lead.next": "Reservar turno",
  "locations.eyebrow": "Red de atenciÃ³n", "locations.title": "EncontrÃ¡ tu franquicia",
  "locations.book": "Reservar atenciÃ³n",
  "locations.start": "ElegÃ­ un modelo y enviÃ¡ tu consulta para reservar un turno asociado.",
  "appointment.select": "ElegÃ­ un turno disponible.",
  "appointment.failed": "No pudimos reservar el turno.",
  "appointment.invalid": "No pudimos verificar la referencia del turno. ReintentÃ¡ sin cambiar la selecciÃ³n.",
  "appointment.received": "Turno solicitado. Referencia:",
  "appointment.confirmation": "La franquicia confirmarÃ¡ la disponibilidad.",
  "appointment.slot": "Turno disponible", "appointment.option": "SeleccionÃ¡ una opciÃ³n",
  "appointment.empty": "No hay turnos publicados en los prÃ³ximos 30 dÃ­as.",
  "appointment.sending": "Reservandoâ€¦", "appointment.submit": "Solicitar turno",
  "appointment.one": "lugar", "appointment.other": "lugares",
  "kind.consultation": "Consulta", "kind.test-drive": "Prueba de manejo",
  "kind.delivery": "Entrega", "kind.service": "Servicio"
} as const;
type MessageKey = keyof typeof es;
const en: Record<MessageKey, string> = {
  "nav.network": "Franchise network",
    "nav.help_content": "Help content",
  "site.description": "Catalog and portals for an electric mobility network.",
  "site.skip": "Skip to content", "site.navigation": "Main", "site.contact": "Contact:",
  "nav.warranty": "Warranty", "nav.supply": "Supply", "nav.catalog_editor": "Edit catalog", "nav.training": "Training", "nav.help": "Help", "nav.dashboard": "My workspace",
  "nav.public_catalog": "Models", "nav.locations": "Locations", "nav.customer": "My account",
  "nav.admin": "Operations", "nav.franchise": "Franchise", "nav.factory": "Factory", "nav.publishing": "Publishing", "nav.finance": "Finance",
  "home.eyebrow": "Configurable electric mobility",
  "home.description": "Explore published models and send an inquiry with your consent.",
  "home.action": "View models", "models.title": "Electric models", "models.short": "Models",
  "models.eyebrow": "Public catalog", "models.description": "Choose a model to request information.",
  "models.one": "model available", "models.other": "models available",
  "lead.name": "Name", "lead.email": "Email", "lead.consent": "I agree to be contacted about this model.",
  "lead.sending": "Sendingâ€¦", "lead.submit": "Request information",
  "lead.failed": "We could not register your request.",
  "lead.invalid": "The response did not include a request reference.",
  "lead.received": "Request received. You can now book a consultation or test drive.",
  "lead.next": "Book an appointment",
  "locations.eyebrow": "Service network", "locations.title": "Find your franchise",
  "locations.book": "Book an appointment",
  "locations.start": "Choose a model and send your inquiry to book a related appointment.",
  "appointment.select": "Choose an available appointment.",
  "appointment.failed": "We could not request the appointment.",
  "appointment.invalid": "We could not verify the appointment reference. Retry without changing your selection.",
  "appointment.received": "Appointment requested. Reference:",
  "appointment.confirmation": "The franchise will confirm availability.",
  "appointment.slot": "Available appointment", "appointment.option": "Select an option",
  "appointment.empty": "No appointments have been published for the next 30 days.",
  "appointment.sending": "Bookingâ€¦", "appointment.submit": "Request appointment",
  "appointment.one": "place", "appointment.other": "places",
  "kind.consultation": "Consultation", "kind.test-drive": "Test drive",
  "kind.delivery": "Delivery", "kind.service": "Service"
};
export type PublicLanguage = "es" | "en";
export interface PublicLocale { language: PublicLanguage; locale: string; timeZone: string }

// Only trusted business configuration selects the locale; never headers/query/cookies.
export function resolvePublicLocale(requested: string, timeZone: string): PublicLocale {
  let locale = "es";
  try {
    const canonical = Intl.getCanonicalLocales(requested)[0];
    if (canonical && ["es", "en"].includes(new Intl.Locale(canonical).language)) locale = canonical;
  } catch { /* Unsupported or malformed configuration falls back to the actual Spanish catalog. */ }
  const language = new Intl.Locale(locale).language as PublicLanguage;
  // Invalid time zone fails configuration explicitly; do not silently book in another zone.
  new Intl.DateTimeFormat(locale, { timeZone }).format(0);
  return { language, locale, timeZone };
}

export function publicMessage(locale: string, key: string): string {
  const language = resolvePublicLocale(locale, "UTC").language;
  const catalog: Readonly<Record<string, string>> = language === "en" ? en : es;
  return Object.hasOwn(catalog, key) ? catalog[key]! : key;
}

export function publicCount(locale: string, kind: "models" | "appointment", count: number): string {
  if (!Number.isSafeInteger(count) || count < 0) throw new RangeError("Invalid public count");
  const resolved = resolvePublicLocale(locale, "UTC");
  const category = new Intl.PluralRules(resolved.locale).select(count);
  const suffix = publicMessage(resolved.locale, kind + (category === "one" ? ".one" : ".other"));
  return new Intl.NumberFormat(resolved.locale).format(count) + " " + suffix;
}

export function publicAppointmentTime(config: PublicLocale, startsAt: string): string {
  const date = new Date(startsAt);
  if (!Number.isFinite(date.getTime())) throw new RangeError("Invalid appointment date");
  return new Intl.DateTimeFormat(config.locale, {
    dateStyle: "medium", timeStyle: "short", timeZone: config.timeZone
  }).format(date);
}

