# TypeScript Multi-Role Onboarding — V185

## Resultado estrecho

V185 materializa `TS-MULTIROLE-ONBOARDING 0.1.0`, el panel multi-rol simplificado y el onboarding por rol (dueño/administrador/empleado/cliente) sobre el BFF web: navegación por permiso, dashboard de espacio de trabajo y guías de uso, gateados por la sesión OIDC.

El rol es sólo un conjunto curado de permisos para el onboarding; **nunca** un atajo de autorización. La visibilidad de cada sección proviene de `readSession()` y `visibleSections(permissions)`; el backend revalida token/permiso/organización.

## Estado honesto

`implementation: RECONSTRUCTIBLE` (no `REBUILD_VERIFIED`): los seis archivos materializan byte a byte con SHA-256, pero el gate de typecheck/test/build (`pnpm install --frozen-lockfile --offline && pnpm typecheck && pnpm test && pnpm build`) requiere el perfil web completo + `node_modules`, no disponible en este host (sin red, pnpm bloqueado). No se declara verificado sin correr ese gate.

## Autoridad y procedencia

- Código `AUTHORED` sobre los patrones del bridge (`readSession`/`allowed`, typed routes, nav por permiso).
- Next.js 16.3.2 / React 19.2.8 (dependencias del perfil web).

## Archivos materializados (6)

| Archivo | SHA-256 |
|---|---|
| src/platform/roles/role-visibility.ts | b9a90455f567013b5a0195c4e9c33b19fe52ee1f1b059f4aee06247f647de30e |
| src/platform/roles/role-visibility.test.ts | 434e701578dcfcb2b2953d46dd64a6bf0e2597dee92328dfb58b62108a03b985 |
| src/components/role-dashboard.tsx | e7cc110ff72b3aee84c11e6ca6269cde2b3aec21ed985d3dc08a3160b5f92a80 |
| src/app/dashboard/page.tsx | ae63e7f010ee3b68f80284e63df5451c5f25ec24789b9a5136c8ba668a07d9af |
| src/app/guide/page.tsx | aa3e6fd16f3d43ab42212a4a42292a76102d255d8d3a30cc0730bb4bc1391341 |
| src/app/guide/[role]/page.tsx | 548519e49768fc35ab589560d9f0a6ea8d26696745935129b4912b3ac220d8c8 |

SHA-256 del pack: `fbe0384d9cb25933c027854995005ec8fba9609eb6091a813200a3abc7996c8f`.

## Verificación ejecutada

- Round-trip de hashes: 6/6 bloques reproducen byte a byte los archivos.
- `role-visibility.test.ts`: 6 aserciones (wildcard, permiso exacto, sin ampliar scope, guía accesible, roles con pasos) — corren bajo `pnpm test` del perfil web (pendiente de ejecutar).

## Cobertura de invariantes (declaradas, pendientes de gate)

1. Visibilidad = permiso exacto o wildcard; nunca se amplía.
2. La guía de uso queda accesible a cualquier rol autenticado.
3. El rol no concede permiso: lo concede la sesión.

## Condiciones residuales

- `REBUILD_VERIFIED`: `pnpm install --frozen-lockfile --offline && pnpm typecheck && pnpm test && pnpm build` sobre el perfil web completo.
- Playwright/Lighthouse sobre el panel multi-rol: gates independientes del proyecto.
- IdP/backend reales, accesibilidad AT, RUM/load y producción: condicionados.

V185 añade el panel y onboarding multi-rol; no declara verificado sin el gate pnpm.
