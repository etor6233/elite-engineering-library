# Customer surveys 1.0.0

The connected reference accepts one versioned 0–10 response from an authenticated
active CRM customer, stores it in PostgreSQL, recovers the original receipt and
exposes a separately authorized aggregate. This is AUTHORED code. It replaces the
in-memory store only for this narrow journey; it does not change or certify the
GO-SURVEYS-CORE source.

## Enable and configure

Apply migrations 0001–0054 in numeric order with the existing database migration
procedure. 0054 creates two CRM tables and a definition immutability trigger.
The down migration deletes these tables and their responses: use only on an
empty disposable database, or after the project's approved recovery procedure.

Set CUSTOMER_SURVEYS_ENABLED=1 on the API and customer_surveys=true in the
business configuration features map. Both are disabled when absent.
The BFF must use its existing configured API origin, OIDC session and HTTPS
APP_BASE_URL. A flag does not grant any permission.

An authorized database operator creates a crm.survey_definition with all of:
tenant_id, organization_id, survey_id, prompt, consent_version, consent_notice,
opens_at, closes_at, retain_until, minimum_responses. The default active=false
allows review before setting active=true. Use ASCII identifiers (letters,
digits, underscore or hyphen; initial letter/digit); versions additionally
permit dot and colon. Questions and notices are plain text.
closes_at must follow opens_at, and retain_until must follow closes_at.
Use a new survey_id when terms, question, dates or reporting threshold change.
Only active can be changed on an existing definition.

The project supplies the real question, notice, dates and reporting threshold.
Use an appropriate recommendation question before describing its result as NPS.
The metric does not prove representativeness, consent sufficiency or legal
compliance. No real policy is supplied by the synthetic test fixture.

Customer invitation URL:
/customer/surveys?organizationId=ORG&surveyId=SURVEY

Operator results URL:
/admin/surveys?organizationId=ORG&surveyId=SURVEY

The project distributes these links through its authorized existing workflow.
No outbound messaging, campaign, invitation list or purchase verification is
implemented here. The customer needs customer:self and the exact organization;
the result reader needs surveys:read and the exact organization. API identity
and active CRM membership are checked again on the server.

## Response and recovery

The response key is tenant/organization/survey/customer. Equal resubmission
returns the original receipt; a different score/version conflicts. Consent must
be explicit and match the definition. The database checks the admission window
after acquiring customer/definition locks. Closing or deactivating a survey
prevents new responses but preserves recovery until retention expires.

If the POST outcome is unknown, the UI disables submission and offers a GET-only
recovery button. Reload first reads the existing response. No automatic POST
retry or replacement of a stored answer occurs. The UI shows the actual stored
score, including zero. A new attempt with changed input cannot overwrite it.

All supplied query/body fields are strictly bounded. Customers cannot submit
tenant/customer IDs, read another customer's response or use aggregate routes
without the separate grant. Below the configured minimum, the aggregate returns
the response count and no NPS; this is not a promise of count suppression.

## Retention and operations

Reads stop exposing responses at retain_until. Physical deletion requires the
operator command, with credentials authorized for the selected tenant/org:

    go run ./cmd/customer-survey-retention -tenant TENANT_UUID -organization ORG -limit 100

DATABASE_URL is read only by the process. The command deletes at most 1–1000
expired responses and prints {"deleted":N}. It has a ten-second deadline.
Repeat under the project's approved schedule until deleted=0; failed or
interrupted runs may be repeated because only expired rows are eligible.
No scheduler, backup erasure, legal retention policy or production SLA is
inferred. Record successful runs and failures using the project's monitoring.
Definitions contain no individual answer and are retained separately.

Customer deletion/restriction, backups, audit access, runtime database grants,
load budgets and deployment/restore remain project gates. The reference
database account in tests is synthetic and is not a production role design.

## Evidence and scope

The focal fixture runs real Go HTTP, RS256/JWKS verification, PostgreSQL, Next
JWE sessions and Chromium desktop/mobile, Firefox and WebKit. It verifies
normal submit, duplicate-click fencing, zero score, a lost POST response,
GET-only recovery, permission denial, minimum response threshold and NPS.
Separate probes verify concurrent commits, deadline/lock ordering and actual
database restart. Retention checks bound deletion and preserve other scopes.

The library's broader integration, dependency admission and signed final release
gates remain independent. No 48/48 or production approval follows from this pack.

V402319: the numerical NPS owner is now an explicitly declared PostHog MIT adaptation. See docs/posthog-nps.md and third_party/posthog-nps/source-lock.json. Original caller glue remains AUTHORED; no corporate authorship is assigned to the API, database, UI or privacy policy. The API still returns unrounded float64 with the existing explicit reporting minimum.
