# Enterprise Web BFF

Frontend y BFF opcional para el perfil empresarial Go/PostgreSQL. No contiene migrations, acceso SQL, persistencia, workers ni reglas transaccionales.

## Composición

Se materializa con `TS-OIDC-PORTAL-ADAPTER`; el backend debe exponer las APIs públicas y las queries protegidas declaradas por el perfil empresarial.

Configuración no secreta: `BUSINESS_CONFIG_FILE`, `APP_BASE_URL`, `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `ENTERPRISE_API_BASE_URL`, `ENTERPRISE_TENANT_CODE` y `ENTERPRISE_ORGANIZATION_CODE`. `AUTH_SESSION_SECRET` y `OIDC_CLIENT_SECRET` se entregan mediante el mecanismo de secretos elegido, nunca en el repositorio.

## Verificación

```powershell
pnpm install --frozen-lockfile --offline
pnpm typecheck
pnpm test
pnpm build
pnpm licenses:report
```

La composición genera un nonce impredecible por request mediante `src/proxy.ts`, fuerza render dinámico y aplica una política CSP sin `unsafe-inline` ni `unsafe-eval` en producción. Google CSP Evaluator 1.1.8 verifica la política y Microsoft Playwright comprueba en cuatro navegadores que cada script lleva el nonce de su respuesta y que el siguiente request recibe otro valor.

Producción sigue condicionada al IdP, HTTPS/edge, comportamiento del CDN/WAF, protección CSRF donde corresponda, antiabuso distribuido, accesibilidad con navegador/AT, carga representativa, observabilidad y rollout/rollback del proyecto. El gate local no sustituye la repetición sobre el edge productivo.

## Imágenes

Esta referencia no transforma imágenes en runtime. `next.config.ts` establece `images.unoptimized: true`; `pnpm-workspace.yaml` excluye la dependencia opcional `sharp` y el lock congelado conserva esa selección. `/_next/image` responde404, comprobado en el recorrido de agenda con cuatro navegadores.

Antes de agregar optimización de imágenes, admitir la canalización elegida, versiones exactas, componentes nativos, licencias, avisos de seguridad y pruebas de rendimiento. Quitar estas dos opciones no equivale a admitir el grafo anterior. La omisión de Sharp no sustituye los gates del resto de dependencias ni la aceptación del proyecto.
