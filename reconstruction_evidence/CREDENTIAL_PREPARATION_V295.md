# V295 — preparación de credenciales y diferimiento fiscal

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Fecha: 2026-09-07. Mantenimiento sobre checkpoint21/V294; no otro roadmap.
Owners T2801/T2803 y dependencias T2805–T2809. No cuenta/contrato/gasto/mensaje
productivo, credencial real ni trámite fiscal creado. ARCA diferida por usuario.

## Resultado ejecutable y alcance

1. PROJECT_SECRETS_TEMPLATE reemplazado por guía única: clasificación/titularidad,
   orden de decisiones, DEV/PROD, consumers reales, Config vs entorno, alternativas,
   custodia, fuentes oficiales, pruebas, rotación/revocación/recuperación y errores.
   Referenciada por protocolo de entrada/readiness y accelerator existentes.
2. GO-ELECTROMOBILITY-APPLICATION 1.9.1, AUTHORED, CONDITIONED: ARCA_ENABLED
   false/ausente omite FiscalModule sin leer socket; true valida socket y conserva
   permisos; valores inválidos rechazan. No montar workers/sidecars fiscales.
   No se elimina schema ni se habilitan emisión o notas fiscales.
3. project_local_secrets.ps1 / test_project_local_secrets.ps1: utilidad AUTHORED
   sobre DPAPI/.NET de Microsoft, sólo DEV Windows, no bóveda productiva. Prompt
   oculto, CSPRNG, DACL, sin sobreescritura, Check sin valores/hashes, proceso hijo
   con entorno limitado y stdout/stderr descartados; no logging ni supervisor.
   Se distribuyen con allowlists exactas y suite en VERIFY_LIBRARY.

## Evidencia local

Staging propio conservado:
`%LOCALAPPDATA%\Temp\elite-v295-7adb00832ea64cfd81f826f9b827bc45`.
No dependencias nuevas ni SDK/runtime actualizados. Go1.26.7 local fijado.

| Prueba | Resultado y límite |
|---|---|
| Composición del perfil canónico | 67 packs/743 archivos; main, test y runtime doc reconstruidos 3/3 byte-idénticos; no frontend cambiado |
| Config fiscal | Ocho casos: absent/false/stray socket, flag inválido/whitespace, true sin socket/relativo/válido. Fiscal 404 desactivado, 401 sin identidad al activarlo; ruta independiente 204 en mux sintético |
| Config providers | Diez casos: ausente/registro vacío/malformado/inline/referencia ausente/secreto ausente/corto/válido/duplicado/trailing. Ausente no lee secretos, errores no contienen valor sintético |
| Suites Go focales | 22 tests PASS +18 subcasos PASS, tres SKIP de integración que requieren PostgreSQL; no contados como PASS conectado. Paquetes main/providerintegration/app/contactidentity/domainbind |
| Go vet/build ./... | PASS en árbol reconstruido; no reemplaza seguridad ofensiva ni assurance |
| Custodia Windows | 16 checks PASS: missing/nombre inválido/duplicados/control de proceso, escritura cifrada, lectura, sobreescritura rechazada, CSPRNG32, hijo recibe sólo secreto elegido, padre intacto, config secreta/prod rechazadas, corrupción/tipo/ACL amplia rechazados |
| Correspondencia automática | 69 pares consumer/nombre documentados. Incluye Go/TS, referencias de perfiles y interpolación YAML; revisión manual adicional de Config, SDK auth, workers, release y alternativas documentales |

El script `verify-v295.ps1` conserva selección extraída del plan, cotejo de hashes,
inventario y comandos. Regex no demuestra exhaustividad semántica de cualquier
proyecto: excluye fixtures/tests/C# fiscal diferido y no infiere constructor/env.
El módulo refund anidado reveló MERCADO_PAGO_ACCESS_TOKEN/REFUND_WORKER_ID y
STRIPE_SECRET_KEY reales para devolución, no loader del main de cobro. Se corrigió
la guía. Sus errores raw/logging requieren revisión antes de habilitarlo live.

No se repitió el navegador/PG de V294 sin delta frontend; se conserva como evidencia
histórica, no prueba nueva del arranque main sin ARCA contra IdP/PG reales. Las
pruebas de ruta independiente aquí son sintéticas; no cierre de todos los journeys.
Los fixtures DPAPI de cada test se retiran al finalizar; sólo fixtures propios,
ningún dato del usuario. Staging de composición/evidencia se conserva.

## Receipts

| Archivo en staging | SHA-256 |
|---|---|
| rebuilt-compose.log | 719f5342483c593e35bbfd6909ff86e3b4ecaf1145e0b32f457e20de20047fdf |
| credential-consumer-inventory.json | 9a48458b327ad7fc6d9c2a1400b905347dc54384ebecfe7c6c8d5b60bbf93a79 |
| rebuilt-config-tests.log | b91c268b757e13a76e3c173052bca909e6efa39863df799f85abb8d65460819e |
| rebuilt-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| local-secrets-tests.log | 280a460983752fb0dd3eff6408f067ebc6e6abed529dd37217cedd79bfaf6542 |
| rebuilt/cmd/electromobility-api/main.go | 051d1a353c2c71b02e8d0c607235187b657f18e409982bb976128c0fabd97e37 |
| rebuilt/cmd/electromobility-api/main_test.go | 737cddd7114d0435120658c8a96aca778074f2b25a9ebde4d12ead39aee7b458 |
| rebuilt/docs/electromobility-api-runtime.md | 3d488453bf876e751d00bffa2a6e4c54982c2255c5a223a58f4cf323af17a061 |

## Fuentes actuales y bloqueos

Consultadas 2026-09-07, URLs enlazadas en guía canónica: Microsoft Export-Clixml
(DPAPI mismo usuario/equipo), AWS Secrets Manager/Textract/SES, Google Secret
Manager/ADC/Ads configuración/Developer Token/Merchant autorización, Azure Key
Vault/Foundry resource, Cloudflare tokens, Stripe keys, Mercado Pago credentials,
Amazon app registration, Firebase Admin y colección oficial Meta WhatsApp.
OpenAI Projects distingue spend alerts de hard spend limit: disponibilidad y
enforcement de cuenta se verifican antes de afirmar corte de gasto.

Meta leads/docs, TikTok create-subscription y ML auth no recuperados completos
por navegación pública: no se inventan menús/scopes/caducidad. No se atribuye a
proveedores este helper ni el diseño local. Sin nueva adquisición VERBATIM.
Provenance y permisos se conservan en packs; guías de alta no readmiten código.

## Pendientes concretos, sin cierre global

- Elegir proyecto/target/IdP y cuentas/proveedores que realmente se usarán. Completar
  **en la misma guía** sus recursos/permisos/menús oficiales y probes sin divulgar
  valores; no pedir todas las alternativas ni ARCA. Esto requiere decisión del titular.
- Obtener documentación oficial accesible de app Meta leads/TikTok/ML cuando se
  seleccionen; mantener gates existentes y adapters incompletos explícitos.
- Configurar gestor productivo/identidad temporal y wiring del target; completar
  callbacks/claims de ambos tokens y límites de logout/rotación sin inventarlos.
- Continuación local V294 T2802/T2804: portal operativo→permisos HTTP→stock/pago/
  entrega con recuperación. Guía de credenciales no sustituye ese recorrido.
- Assurance/readiness41, operación/privacidad/restore/rollout y gates del target
  siguen abiertos. No porcentaje nuevo, ZIP final ni producción certificada.

Fallos 430–436 registrados; 432 incluye comandos/patches fallidos, sin contarlos
como progreso. Shared ledger agrega LIB-FAIL2144–2148; no reescribe la historia.
El primer VERIFY_LIBRARY rechazó un home local literal en el ejemplo de guía;
se sustituyó por resolución desde la raíz de biblioteca. Se repite, sin presentar
ese rechazo como PASS ni alterar su detector.
El siguiente gate reconstruyó el corpus y detectó recuento de procedencia obsoleto:
se actualizó a AUTHORED1196/ADAPTED133/VERBATIM105/TOTAL1434, sin reclasificar código.
También se actualizó el owner FRANCHISE_PREFLIGHT_GAP, detectado obsoleto por el
control siguiente: 160/1434/728/51 y composición743, sin cambiar historia de V258.
Checkpoint y verificación global se registran después de sincronizar owners.
