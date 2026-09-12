# TypeScript Design System — V192

## Resultado estrecho

V192 materializa `TS-DESIGN-SYSTEM 0.1.0`, el design system accesible y mobile-first: tokens de color con **contraste WCAG AA verificado por Node**, tipografía, spacing, radios, breakpoints, y primitivos UI (Button, Card, Input, StatusBadge, Empty/Loading/Error) que usan los tokens como única fuente de verdad.

## Verificación de contraste (evidencia real, Node)

Fórmula de luminancia relativa W3C ejecutada con Node — **12/12 pares ≥4.5:1**:

| Par | Ratio |
|---|---|
| text/primary #111827 on #FFFFFF | 17.74:1 |
| text/secondary #4B5563 on #FFFFFF | 7.56:1 |
| text/muted #6B7280 on #FFFFFF | 4.83:1 |
| onPrimary #FFFFFF on #2563EB | 5.17:1 |
| onPrimary on #1D4ED8 | 6.70:1 |
| onDanger on #DC2626 | 4.83:1 |
| onSuccess on #15803D | 5.02:1 |
| onWarning on #B45309 | 5.02:1 |
| link #2563EB on surface | 4.95:1 |
| danger #B91C1C on surface | 6.19:1 |
| success #15803D on surface | 4.80:1 |
| warning #B45309 on surface | 4.81:1 |

Los fallos iniciales (verde #16A34A y ámbar #D97706 con texto blanco a 3.3:1) fueron detectados por la verificación y corregidos a tonos más oscuros. Evidencia de que el gate no se declara, se ejecuta.

## Autoridad y procedencia

- React 19.2.8 (MIT) — primitivos `AUTHORED`.
- WCAG 2.x (W3C) — fórmula de contraste, verificación real.

## Archivos materializados (5)

| Archivo | SHA-256 |
|---|---|
| src/design/tokens.ts | 15c7a09cde86cb404805502653f409544de34c244da769a855f85699f7e3d68f |
| src/design/contrast.ts | b2d90f600dff858818df9e47a2291a11d945382b035e78c619675e2f75828cf5 |
| src/design/cx.ts | b6ac9550a87795f64c5d08b383bca112ab685f276d9b575b909442893d54733f |
| src/design/tokens.test.ts | c24e3dd687fbb680ae9b9eb24e2d5bc6c6bc244d76c1fdb4d7e40eca8ba41738 |
| src/components/ui.tsx | c4e47f190cd5c05eb63e6f5b1ee159e8dd015b701e57bac79fe36470c23e444c |

SHA-256 del pack: `43b98851e5cd1816765327df48489b2f13f0d461dfc916589120aae1ae5e75f0`.

## Toolchain fijado

- Node v24.14.1 (verificación de contraste).

## Cobertura de invariantes

1. Todo par de texto ≥4.5:1 (WCAG AA) — verificado por Node.
2. Breakpoints ordenados (sm<md<lg).
3. Tokens como única fuente de color/tamaño (sin hardcode en componentes).
4. Componentes semánticos y accesibles (role/aria, focus, keyboard).

## Condiciones residuales

- Gate `pnpm typecheck/test/build` y Playwright (visual, responsive real): entorno.
- Juicio estético "moderno/agradable": verificación visual del proyecto.

V192 aporta el design system accesible; la verificación visual final es runtime del proyecto.
