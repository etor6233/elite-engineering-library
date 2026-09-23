# Contact task permissions — U0

This documents the existing boundary; it grants no permissions and changes no role assignments.

- Listing/discovering contacts in `/franchise`: `lead:read` and the session organization are required by the existing Go owner.
- Quoting an available contact: `lead:read` for discovery plus `quote:write` for its action.
- Assigning/changing a contact: `lead:read` for discovery plus `lead:assign` / `lead:update` for the respective action.
- A write-only principal can be valid for an explicitly scoped integration/BFF operation; that permission alone does not grant listing or a browser task that needs listing.

The target role/profile must declare these combinations where a complete browser task is REQUIRED. No extra grant is inferred from role names or Galaxy functions. Server tenant/organization/object/subject checks remain authoritative. Four executed role cases establish the current page behavior; they do not certify every target role configuration.

CUIP007 is closed as an existing policy boundary, with no code privilege change. See the U0 receipt and official Next.js authentication guidance: https://nextjs.org/docs/app/guides/authentication.
