# Public journey localization

AUTHORED integration under the existing BFF and journey-portal owners. No new dependency.

The configured business.defaultLocale selects Spanish or English, including regional Intl formatting.
Unsupported or malformed tags fall back to Spanish and are labeled as Spanish. Unknown message keys remain
literal keys. The fixed catalogs contain complete Spanish and English messages; they are not an arbitrary
translation CMS. Intl.PluralRules and NumberFormat handle integer counts, with one/other messages for these
two languages. No claim of general CLDR translation coverage or interpolation.

The default market's explicit timeZone controls displayed appointment time on both server and client.
Invalid time zones fail configuration; no guessed business time zone. The backend's starts_at and kind,
lead and model IDs, consent boolean, endpoints and idempotency keys remain unchanged.
The displayed localized time is never parsed back into a request.

The localized reference journey covers home, models, connected catalog, lead request, locations and appointment
request, including failures, empty state, pending state and durable receipt. Product names and vehicle-class
content come from the backend unchanged. Navigation, footer, metadata and JSON-LD use the resolved language.
Private portals retain their existing Spanish content under main lang=es; public sections override that
language explicitly. Customer/operator portal translation, user language negotiation, URL language variants,
hreflang and multilingual content storage require their own integration and target review.
Consent wording is a translation of the existing technical fixture, not a legal approval.

References: TC39 ECMA-402 Intl specification, Node internationalization documentation, W3C declaring
language in HTML. Exact runtime behavior and connected tests are recorded by V378. Existing native/runtime
admission conditions are retained; using Intl does not close their security or distribution gates.

