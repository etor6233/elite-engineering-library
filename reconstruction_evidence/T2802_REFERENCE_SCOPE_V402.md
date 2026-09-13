# V402 — alcance explícito de subcapacidades de referencia

T2802 sigue IN_PROGRESS. Esta decisión no reduce las48superficies ni excluye J1–J5.
El blueprint ENTERPRISE_FULL_STACK_BLUEPRINT.md, secciones4–5, gobierna el alcance.
No se infiere obligación de implantar todos los cores históricos que no fueron
seleccionados por el perfil. Usuario autoriza NONE_WITH_REASON con evidencia.

| Subcapacidad | Decisión y fundamento |
|---|---|
| GO-PAYROLL-CORE | NONE_WITH_REASON: No wage contract, payroll run, employment deduction or payslip journey in J1–J5. Resource employee/contractor identifies appointment staffing, not employment/payroll accounting. No payroll claim or existing isolated-core promotion. |
| GO-POS-CORE | NONE_WITH_REASON: Reference J1 uses server-price quotation/order and hosted provider checkout or approved stored value. No cash drawer, change, terminal peripherals or offline retail journey selected. Those capabilities need explicit target scope. |
| GO-PROMOTIONS-CORE | NONE_WITH_REASON: The reference uses versioned price books and the already admitted named loyalty discount program. General coupon stacking/campaign coupon codes are not required by J1–J5; no equivalence to the isolated promotions core claimed. |
| GO-REFERRALS-CORE | NONE_WITH_REASON: Lead source attribution exists; a rewarded referral program, invite code and reward policy are not required by J1–J5. No inferred reward/financial rule. |
| GO-REVIEWS-CORE | NONE_WITH_REASON: Reference feedback is the connected private survey owner. Public review publication/moderation/ranking is not selected in J1–J5; survey is not claimed equivalent to public reviews. |
| GO-WAITLIST-CORE | NONE_WITH_REASON: J1/J4 select explicit capacity-controlled appointments. Full capacity rejects admission. No FIFO waiting offer/promotion policy is selected; no waitlist equivalence claim. |

Trigger de reapertura por fila: selección explícita en el blueprint del target;
requiere admisión de fuente/capability y pruebas antes de incorporar el core.
No se promocionan los cores aislados por esta clasificación.

J1 y J4 ya tienen garantía conectada canónica V402. J2 conserva un hueco concreto:
las transiciones PO/fábrica/stock no ligan cantidad solicitada, manifiesto serial,
ASN y recepción parcial. J3 permite draft, pero faltan release/reversión aprobada
y proyección/feed conectados. Se implementan ambos; no se excluyen para cerrar.
J5 y ayuda/evaluación pertenecen a T2804; seguridad/operación y cuentas externas
conservan sus gates separados. Hashes y selección: T2802_REFERENCE_SCOPE_V402.json.
