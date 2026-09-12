# V402 — snapshot de conversión exacta integrado

GO-EXACT-FX-SNAPSHOT-ACCOUNTING0.1.0 tiene23archivos reconstruibles;22seleccionados cuando la licencia MIT idéntica ya la aporta BC amount. ExchangeExact y FindLast derivan del snapshot oficial BCApps2eae56d704a1fd035d104f333602aea7091b7749. RoundMinor es AUTHORED representación exacta con perfil explícito NEAREST_TIES_AWAY_FROM_ZERO; no se atribuye a la implementación interna de AL Round. La discrepancia del ejemplo oficial negativo se conserva en su expediente.

El perfil JSON inmutable fija tenant/organización, base, fechas, precisión, fuente y tasas. API y host opcional crean un recibo de conversión usando accounting, idempotency_record y outbox existentes. PostgreSQL con OIDC RS256/JWKS:12solicitudes concurrentes producen un solo recibo; replay, respuesta perdida, hash/actor/organización divergentes y expiración mientras espera el outbox verificados.525vectores Decimal independientes y18negativos; fuente/aritmética/fuzz/rebuild específicos en manifest. Se reutiliza el fuzz aritmético399113/3s por identidad de implementación; no se afirma que las suites AL se ejecutaron.

El recibo no crea asientos ni dinero. La conexión explícita con posting FX, cuando sea exigida por el claim contable más amplio, sigue pendiente; no cerrar GO-FX-CORE íntegro por el nombre del pack. Source, source-test y licencia están fijados en el mismo commit. Global G5–G8 de la composición siguen separados; ninguna cuenta live es necesaria para los fixtures.

| Recibo | SHA-256 |
|---|---|
| `fx-connected-manifest.json` | `28d356a229c53e8e47eebb55b6b0f00ed1ae81587d5026fabd81467fe618395d` |
| `fx-connected-admission.json` | `1ad0ff701abba8fae369f36f12413c2d1275362db2b93b9acc77c28e8c8a2f81` |
| `fx-social-integration.json` | `8ad182224a8a0683dba5edb5287f01c15ff405bace850670c5d9622faac87d94` |
| `communications-reconstruction.json` | `9c08438a62bfcfd32df138557cdfbaca77b0329274d16e1f3ea382d8926671c5` |
