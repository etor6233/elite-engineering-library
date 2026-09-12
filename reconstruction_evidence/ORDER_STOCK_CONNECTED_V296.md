# V296 — pedido existente → reserva HTTP → recuperación

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Fecha: 2026-09-07. Owners: T2802/T2803/T2804, roadmap §10.
PASS local acotado, CONDITIONED; no cierre de pagos/entregas ni producción.
V295 y ARCA diferida permanecen; no cuentas, cobros ni proveedores activados.

## Código y prueba real

Commerce 0.4.0 añade GET /v1/commerce/orders: snapshot PostgreSQL Repeatable Read,
read-only, organización y permiso operativo. Lee los owners existentes de pedidos,
líneas, stock, pago y handover; no duplica inventario ni expone identidad del cliente
o referencias del proveedor. Fallo de lectura: 503.

Portal 0.12.0: selección de unidad por serie, versiones internas, BFF estricto con
Origin/sesión/organización/permiso y revalidación Go. Reserva con la transacción
existente. AllocateStockAs agrega actor_subject desde el principal verificado al
mismo outbox/commit; sin soporte auditado rechaza, no acepta actor del JSON.
Los callers internos anteriores siguen compatibles, pero requieren revisión propia
de auditoría. Error incierto/conflicto bloquea reenvío hasta lectura nueva.
Ayuda/práctica/soporte order-operations-view/1.0.0.

11 archivos modificados, cero archivos de producto nuevos, cuatro owners:
GO-COMMERCE-PRICING-PAYMENT-API 0.4.0 (5), TS-FRANCHISE-JOURNEY-PORTALS 0.12.0 (4),
GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.3 (harness),
MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.9 (suite AUTHORED).
Main 1.9.1 sólo aclara compatibilidad; conserva código V295.
Reconstrucción en destino ausente: 67 packs/743 archivos, 11/11 byte-idénticos.

PostgreSQL18.6/52 migraciones, Go1.26.7, Next/Playwright de los locks existentes.
Cuatro proyectos (Chromium escritorio/móvil, Firefox, WebKit), seis fases:
aceptación, GET fallido, reinicio; reserva operativa, GET operativo 503, reinicio.
**24 fases PASS**, 12 operativas nuevas; gate final 63,42s.
Los mismos tres pedidos se crean por UI de aceptación, no por inserts de pedidos.
Next se reinicia y API/pool se recrean entre fases; no es proceso main Go productivo.

Por tenant: dos asignaciones, dos reservas canónicas y dos eventos con actor
operator; tercer pedido sin stock. Pruebas: pérdida de respuesta después de commit,
GET sin segundo POST, teclado, reenvío, carrera entre pedidos por una unidad
(200/409), versión vieja, línea ajena, unidad ausente, Origin/organización ajenos,
sesión ausente, lector y permiso BFF falsificado rechazado nuevamente por Go.
Consulta PG: aislamiento tenant/org, estado de pago del fixture y límite explícito.
Tests web104PASS/1SKIP opt-in; Go test ./..., vet/build PASS; cinco negativos HTTP
del reader PASS; integración PG focal PASS. Build Next PASS. Capturas de cuatro
proyectos, ancho móvil sin desborde e inspección visual móvil; no WCAG integral.

Reinicio PG controlado: snapshots ordenados de pedidos/líneas/stock/reservas/pagos/
handovers/outbox idénticos. Cluster propio detenido. No PITR/DR ni restore de backup.

## Receipts y reproducción

Staging: %LOCALAPPDATA%/Temp/elite-v296-a380dd37465046c4a166bf68cdc2baec.
Fallos y pruebas intermedias preservados; datos e identidades exclusivamente sintéticos.

| Receipt | SHA-256 |
|---|---|
| operation-audited-final.log | c2b850a0b2b4b559ad2057483bbbe3071885595d4c7223b26cea1763abdfa017 |
| audited-projection.log | 45390f57b9073e412979a3c42da5ef45c972ade7c0968542f7dc8dfbd21b073f |
| snapshot-after-audited.json, igual a before | 50b349f4e9718bb40c4c3df2668e0f9cfd429e64e6bcdbdab5e45397b09c9558 |
| web-tests.log | 41eaf2eaa29926b3544ea74990602bb677d15118e99286fa8afe2be1f6b3658c |
| rebuilt-final-build.log | 937177c22f8798aa951861f027e5e1c827897e77dde0e26314d37d3fc0bef698 |

Componer FRANCHISE_COMPLETE_PACK_PLAN en destino vacío; instalar ambos locks
offline/frozen (Playwright anidado con --ignore-workspace); construir Next; aplicar
52 migraciones a una base descartable loopback elite_confirmation_*.
Configurar TEST_DATABASE_URL de esa base, ELITE_WEB_ROOT absoluto,
ELITE_QUOTE_E2E=1 y ELITE_ORDER_E2E=1. Nunca usar un target real.

```text
go test ./internal/platform/httpapi -run '^(TestQuoteAcceptanceBrowserPostgres|TestOperationReadFailsClosedWithoutScopeOrReader)$' -count=1 -v
go test ./internal/platform/postgres -run '^TestCommercePriceOrderAllocationPaymentFlow$' -count=1 -v
```

## Procedencia y pendientes

Delta AUTHORED, cero código upstream nuevo. Método consultado el 2026-09-07:
[PostgreSQL18 snapshot](https://www.postgresql.org/docs/18/transaction-iso.html) y
[Microsoft Playwright UI/API](https://playwright.dev/docs/api-testing).
Esas fuentes no aprueban ni certifican el código local.

- Vista limitada: 25 pedidos, 100 líneas/estados por pedido, 200 unidades; truncated
  bloquea reservas. Falta paginación para mayor volumen, carga y operación real.
- Pago por UI/adapter/reconciliación pendiente. Created/authorized no es captured.
  Selección de proveedor/cuenta/sandbox y custodia siguen en la guía V295; poseer
  credenciales no termina un adapter. Ningún worker de cobro corrió en estas pruebas.
- Falta creación inicial de handover para pedido ordinario: los inserts productivos
  encontrados son de excepción/cambio; los demás son fixtures. No copiar la regla
  de devolución para cerrar una venta. Conectar inicio autorizado, stock/pago,
  checklist, aceptación, efectos finales y recuperación requiere ese cierre local.
- Antes de efectos reales, elegir proveedor y definir la regla de liberación de
  entrega. No se asume pago total, crédito ni deuda aprobada. ARCA permanece diferida.
- T2801–T2810, readiness41/assurance, privacidad, monitor, capacitación integral,
  accesibilidad, carga, despliegue/restore y aceptación empresarial siguen pendientes.
  FAIL423 previo no se cierra por un PASS posterior.
- FAIL437 (puerto),438 (fixture con NULL unique),439 (selector Next) y440 (actor)
  corregidos para este scope. Registro raíz y ledger conservan la historia.
