# Dependencias y licencias

Las dependencias directas están fijadas en `package.json` y el grafo completo en `pnpm-lock.yaml`.

| Dependencia directa | Versión | Licencia declarada por el paquete |
|---|---:|---|
| jose | 6.2.10 | MIT |
| Next.js | 16.3.4 | MIT |
| openid-client | 6.8.5 | MIT |
| React / React DOM | 19.2.8 | MIT |
| Google SafeValues | 1.2.0 | Apache-2.0 |
| Google CSP Evaluator | 1.1.8 | Apache-2.0 |
| server-only | 0.0.1 | MIT |
| Zod | 4.4.3 | MIT |
| TypeScript | 7.0.2 | Apache-2.0 |
| Vitest | 4.1.11 | MIT |

Google CSP Evaluator se usa sólo como dependencia de desarrollo y su propio README declara que no es un producto oficial de Google ni ofrece garantía. Regenerar el inventario efectivo con `pnpm licenses:report`. El grafo transitivo observado también contiene Apache-2.0, MIT, ISC, BSD-3-Clause, 0BSD, CC-BY-4.0, MPL-2.0 y componentes binarios de imagen con términos adicionales. Este resumen no reemplaza conservar notices, revisar artefactos realmente distribuidos ni una revisión legal para el modo de entrega elegido.

No se asignó una licencia de distribución al código propio del adapter: esa decisión pertenece al propietario del repositorio. El pack se compone con el overlay OIDC y el adaptador de licencia para producir el inventario legal efectivo.

## FX direct-base adaptation

`internal/bcfx/exchange.go` and `selection.go` adapt Microsoft BCApps commit `2eae56d704a1fd035d104f333602aea7091b7749` under MIT. The exact copyright/license is `licenses/Microsoft-BCApps-MIT.txt`; derivation and source hashes are in `docs/provenance/BC_FX_DERIVATION.md` and `BC_FX_SOURCE_LOCK.json`. `rounding.go` and snapshot/accounting/API/host glue are AUTHORED; no AL built-in equivalence or upstream runtime execution is claimed.

## Odoo stored-value calculation

The optional gift-card/loyalty module adapts selected Odoo Community19.0 code at99edb6dd82b7b560930c00b03b694ba700785370 under LGPL-3.0-only. Complete license, copyright, source, local changes and replacement instructions are preserved in odoo_loyalty/ and docs/provenance/ODOO_STORED_VALUE_NOTICES.md. Go/SQL/portal glue is AUTHORED with its own declared provenance. No complete Odoo runtime or corporate authorship for that glue is claimed. The Next.js direct-version notice above is aligned to the existing16.3.4 package/lock; no dependency update was performed by this notice correction.

## Connected document reference

When the AWS Textract/document capability is selected (these dependencies are not implied in a web-only profile), AWS SDK for Go v2 Textract1.45.0 and its15fixed modules: Apache-2.0; complete original LICENSE/NOTICE and per-module hashes are in `aws_textract_runtime/licenses/selected-notices.json`. Preserve all four original texts on redistribution. The Microsoft public invoice fixture uses its exact MIT notice at `azure_document_intelligence_official_invoice/upstream/LICENSE.txt`. Original fixture bytes are acquired by fixed commit/hash, not included here. Document pipeline/test/configuration glue is AUTHORED locally and is not attributed to AWS, Microsoft or Google.
