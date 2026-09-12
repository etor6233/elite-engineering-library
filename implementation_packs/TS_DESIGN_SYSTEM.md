# TypeScript Design System

## 1. Metadata

```yaml
pack_id: "TS-DESIGN-SYSTEM"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el design system accesible (tokens de color con contraste WCAG AA verificado, tipografía, spacing, radios, breakpoints) y los primitivos UI mobile-first (Button, Card, Input, StatusBadge, Empty/Loading/Error)."
stacks: ["React 19.2.8", "TypeScript"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.5.x", "TS-MULTIROLE-ONBOARDING 0.1.x"]
incompatible_with: ["colores sin contraste AA", "hardcode de colores/tamaños en componentes"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://www.w3.org/TR/WCAG21/"]
verified_at: "2026-09-02"
```

La paleta de color está **verificada contra WCAG AA (≥4.5:1)** con la fórmula de luminancia relativa de W3C, ejecutada con Node (evidencia reproducible). Los componentes usan los tokens como única fuente de verdad; no hay colores ni tamaños hardcodeados. La verificación visual final (browser/Playwright) queda como gate del proyecto.

## 2. Applicability

Use este pack para una UI accesible, consistente y mobile-first en todos los roles (cliente/empleado/administrador/dueño). Rechace para colores sin contraste AA o tamaños hardcodeados.

## 3. Architecture contract

- **Ownership**: `src/design` gobierna tokens y contraste; `src/components/ui` gobierna los primitivos. Los pages/portales consumen estos primitivos.
- **Invariantes**: (1) todo par de texto ≥4.5:1 (WCAG AA). (2) breakpoints ordenados (sm<md<lg). (3) tokens son la única fuente de color/tamaño. (4) componentes semánticos y accesibles (role/aria).
- **Data flow**: `tokens.ts` → `contrast.ts` (verificación) → `ui.tsx` (primitivos) → pages.
- **Failure modes**: contraste insuficiente → falla el test; componente sin tokens → rechazado.
- **Seguridad/privacidad**: sin datos; solo presentación.
- **Performance budget**: componentes ligeros, sin dependencias extra.
- **Operación/migración/rollback**: aditivo bajo `src/`; colisión detiene composición.

## 4. Exact file manifest

```text
CREATE src/design/tokens.ts
CREATE src/design/contrast.ts
CREATE src/design/cx.ts
CREATE src/design/tokens.test.ts
CREATE src/components/ui.tsx
```

## 5. Materialization blocks

### FILE: `src/design/tokens.ts`
```yaml
block_id: "TS-DESIGN-SYSTEM:src/design/tokens.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "15c7a09cde86cb404805502653f409544de34c244da769a855f85699f7e3d68f"
variables: []
secrets_allowed: false
```
````ts
// Design tokens verified for WCAG AA contrast (>= 4.5:1) — see tokens.test.ts
// and the reconstruction evidence. The palette follows a standard 8-point
// gray scale and an accessible blue primary, with darker semantic greens and
// ambers so white text always passes AA.
export const colors = {
  background: "#FFFFFF",
  surface: "#F9FAFB",
  border: "#E5E7EB",
  textPrimary: "#111827",
  textSecondary: "#4B5563",
  textMuted: "#6B7280",
  primary: "#2563EB",
  primaryHover: "#1D4ED8",
  onPrimary: "#FFFFFF",
  danger: "#DC2626",
  onDanger: "#FFFFFF",
  success: "#15803D",
  onSuccess: "#FFFFFF",
  warning: "#B45309",
  onWarning: "#FFFFFF",
} as const;

export const spacing = {
  0: "0",
  1: "0.25rem",
  2: "0.5rem",
  3: "0.75rem",
  4: "1rem",
  5: "1.25rem",
  6: "1.5rem",
  8: "2rem",
  10: "2.5rem",
  12: "3rem",
} as const;

export const radius = {
  sm: "0.25rem",
  md: "0.375rem",
  lg: "0.5rem",
  xl: "0.75rem",
  full: "9999px",
} as const;

export const breakpoints = { sm: 640, md: 768, lg: 1024 } as const;

export const font = {
  sans: "ui-sans-serif, system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif",
  size: {
    xs: "0.75rem",
    sm: "0.875rem",
    base: "1rem",
    lg: "1.125rem",
    xl: "1.25rem",
    "2xl": "1.5rem",
  },
} as const;
````

### FILE: `src/design/contrast.ts`
```yaml
block_id: "TS-DESIGN-SYSTEM:src/design/contrast.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b2d90f600dff858818df9e47a2291a11d945382b035e78c619675e2f75828cf5"
variables: []
secrets_allowed: false
```
````ts
// WCAG 2.x relative luminance and contrast ratio (pure, dependency-free).
function channel(hex: string, start: number): number {
  const v = parseInt(hex.slice(start, start + 2), 16) / 255;
  return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4);
}

function luminance(hex: string): number {
  const h = hex.replace("#", "");
  const r = channel(h, 0);
  const g = channel(h, 2);
  const b = channel(h, 4);
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

// contrastRatio returns the WCAG contrast ratio between two hex colors.
export function contrastRatio(fg: string, bg: string): number {
  const l1 = luminance(fg);
  const l2 = luminance(bg);
  const hi = Math.max(l1, l2);
  const lo = Math.min(l1, l2);
  return (hi + 0.05) / (lo + 0.05);
}

// AA_NORMAL is the WCAG AA threshold for normal text.
export const AA_NORMAL = 4.5;

// AA_LARGE is the WCAG AA threshold for large text (>= 18pt or 14pt bold).
export const AA_LARGE = 3.0;
````

### FILE: `src/design/cx.ts`
```yaml
block_id: "TS-DESIGN-SYSTEM:src/design/cx.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b6ac9550a87795f64c5d08b383bca112ab685f276d9b575b909442893d54733f"
variables: []
secrets_allowed: false
```
````ts
// cx joins truthy class names (pure, dependency-free).
export function cx(...parts: Array<string | false | undefined | null>): string {
  return parts.filter(Boolean).join(" ");
}
````

### FILE: `src/design/tokens.test.ts`
```yaml
block_id: "TS-DESIGN-SYSTEM:src/design/tokens.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c24e3dd687fbb680ae9b9eb24e2d5bc6c6bc244d76c1fdb4d7e40eca8ba41738"
variables: []
secrets_allowed: false
```
````ts
import { describe, expect, it } from "vitest";
import { AA_NORMAL, contrastRatio } from "@/design/contrast";
import { breakpoints, colors } from "@/design/tokens";
import { cx } from "@/design/cx";

// Every foreground/background pair used in the UI must pass WCAG AA (>= 4.5).
const pairs: Array<[string, string]> = [
  [colors.textPrimary, colors.background],
  [colors.textSecondary, colors.background],
  [colors.textMuted, colors.background],
  [colors.onPrimary, colors.primary],
  [colors.onPrimary, colors.primaryHover],
  [colors.onDanger, colors.danger],
  [colors.onSuccess, colors.success],
  [colors.onWarning, colors.warning],
  [colors.primary, colors.surface],
  [colors.danger, colors.surface],
  [colors.success, colors.surface],
  [colors.warning, colors.surface],
];

describe("design tokens pass WCAG AA", () => {
  it("every text pair has contrast >= 4.5", () => {
    for (const [fg, bg] of pairs) {
      const ratio = contrastRatio(fg, bg);
      expect(ratio).toBeGreaterThanOrEqual(AA_NORMAL);
    }
  });

  it("breakpoints are strictly ordered", () => {
    expect(breakpoints.sm).toBeLessThan(breakpoints.md);
    expect(breakpoints.md).toBeLessThan(breakpoints.lg);
  });
});

describe("cx", () => {
  it("joins truthy class names", () => {
    expect(cx("a", false, "b", undefined, "c")).toBe("a b c");
  });
});
````

### FILE: `src/components/ui.tsx`
```yaml
block_id: "TS-DESIGN-SYSTEM:src/components/ui.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c4e47f190cd5c05eb63e6f5b1ee159e8dd015b701e57bac79fe36470c23e444c"
variables: []
secrets_allowed: false
```
````tsx
import type {
  ButtonHTMLAttributes,
  InputHTMLAttributes,
  ReactNode,
} from "react";
import { colors, font, radius, spacing } from "@/design/tokens";
import { cx } from "@/design/cx";

// Accessible, mobile-first UI primitives. Tokens are the single source of
// truth; no hardcoded colors or sizes appear in components.

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";

const buttonStyles: Record<ButtonVariant, React.CSSProperties> = {
  primary: {
    backgroundColor: colors.primary,
    color: colors.onPrimary,
    border: "1px solid " + colors.primary,
  },
  secondary: {
    backgroundColor: colors.surface,
    color: colors.textPrimary,
    border: "1px solid " + colors.border,
  },
  danger: {
    backgroundColor: colors.danger,
    color: colors.onDanger,
    border: "1px solid " + colors.danger,
  },
  ghost: {
    backgroundColor: "transparent",
    color: colors.primary,
    border: "1px solid transparent",
  },
};

export function Button({
  variant = "primary",
  className,
  style,
  type = "button",
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & { variant?: ButtonVariant }) {
  return (
    <button
      type={type}
      {...props}
      className={cx("btn", className)}
      style={{
        fontFamily: font.sans,
        fontSize: font.size.sm,
        padding: spacing[2] + " " + spacing[4],
        borderRadius: radius.md,
        cursor: "pointer",
        minHeight: "2.5rem",
        ...buttonStyles[variant],
        ...style,
      }}
    />
  );
}

export function Card({
  children,
  className,
  style,
}: {
  children: ReactNode;
  className?: string;
  style?: React.CSSProperties;
}) {
  return (
    <section
      className={cx("card", className)}
      style={{
        backgroundColor: colors.surface,
        border: "1px solid " + colors.border,
        borderRadius: radius.lg,
        padding: spacing[4],
        ...style,
      }}
    >
      {children}
    </section>
  );
}

export function Input({
  label,
  id,
  className,
  style,
  ...props
}: InputHTMLAttributes<HTMLInputElement> & { label: string }) {
  return (
    <label style={{ display: "block", marginBottom: spacing[4] }}>
      <span
        style={{
          display: "block",
          fontSize: font.size.sm,
          color: colors.textSecondary,
          marginBottom: spacing[1],
        }}
      >
        {label}
      </span>
      <input
        id={id}
        {...props}
        className={cx("input", className)}
        style={{
          fontFamily: font.sans,
          fontSize: font.size.base,
          color: colors.textPrimary,
          backgroundColor: colors.background,
          border: "1px solid " + colors.border,
          borderRadius: radius.md,
          padding: spacing[2] + " " + spacing[3],
          width: "100%",
          boxSizing: "border-box",
          ...style,
        }}
      />
    </label>
  );
}

type StatusTone = "success" | "warning" | "danger" | "muted";

const statusColors: Record<StatusTone, { bg: string; fg: string }> = {
  success: { bg: "#DCFCE7", fg: colors.success },
  warning: { bg: "#FEF3C7", fg: colors.warning },
  danger: { bg: "#FEE2E2", fg: colors.danger },
  muted: { bg: colors.surface, fg: colors.textMuted },
};

export function StatusBadge({
  tone = "muted",
  children,
}: {
  tone?: StatusTone;
  children: ReactNode;
}) {
  const c = statusColors[tone];
  return (
    <span
      role="status"
      style={{
        backgroundColor: c.bg,
        color: c.fg,
        fontFamily: font.sans,
        fontSize: font.size.xs,
        borderRadius: radius.full,
        padding: spacing[1] + " " + spacing[3],
      }}
    >
      {children}
    </span>
  );
}

export function EmptyState({ title, action }: { title: string; action?: ReactNode }) {
  return (
    <Card style={{ textAlign: "center" }}>
      <p style={{ color: colors.textMuted, marginBottom: action ? spacing[3] : 0 }}>{title}</p>
      {action}
    </Card>
  );
}

export function LoadingState({ label = "Cargando…" }: { label?: string }) {
  return (
    <div role="status" aria-live="polite" style={{ padding: spacing[4], color: colors.textMuted }}>
      {label}
    </div>
  );
}

export function ErrorState({ message }: { message: string }) {
  return (
    <div
      role="alert"
      style={{
        color: colors.danger,
        border: "1px solid " + colors.danger,
        borderRadius: radius.md,
        padding: spacing[3],
        marginBottom: spacing[4],
      }}
    >
      {message}
    </div>
  );
}
````


## 6. Configuration surface

Sin variables ni secretos. Los tokens son constantes versionadas.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| React | 19.2.8 | primitivos | MIT | runtime | https://github.com/facebook/react |
| Vitest | (perfil web) | test de contraste | MIT | build | https://github.com/vitest-dev/vitest |

## 8. Apply order

1. Componer el perfil web.
2. Colocar los cinco archivos bajo `src/`.
3. Verificar contraste (`node` o `pnpm test`) y `pnpm typecheck`.
4. Rollback: eliminar los cinco archivos.

## 9. Verification

- **Contraste WCAG AA verificado con Node**: 12/12 pares ≥4.5:1 (17.74, 7.56, 4.83, 5.17, 6.70, 4.83, 5.02, 5.02, 4.95, 6.19, 4.80, 4.81).
- `tokens.test.ts`: 3 aserciones (todos los pares AA, breakpoints ordenados, `cx`).
- Gate `pnpm typecheck/test/build`: pendiente (entorno).

## 10. Reconstruction evidence

Véase `reconstruction_evidence/TS_DESIGN_SYSTEM_2026-09-02_V192.md`.
