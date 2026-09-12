# V287 — entrada conectada del agente: NEW / EXISTING

2026-09-07. T2801, mantenimiento de biblioteca; **no certificación de franquicia**.
Test canónico: `markdown_system/test_agent_entry_lifecycle.ps1`, AUTHORED,
SHA-256 `158b9268658736db55a06ae4f22574827478b56cde1211aa1e76636b1feff21f`.

## Recorrido ejecutado, no piezas aisladas

```text
inventario vacío o archivo existente con BOM
→ reconstruir compositor 0.2.2 + execution validator 1.3.0 desde Markdown
→ componer readiness 0.6.1 en subdirectorio nuevo, sin pisar producto
→ instalar bridge Codex/Claude
→ generar asesoramiento A y ejecutar readiness real: BLOCKED
→ capturar inventario/delta/evidencia por hash y checkpoint 1
→ reanudar con validación real
→ rechazar colisión, checkpoint duplicado y artefacto alterado
→ restaurar bytes exactos y reanudar
→ checkpoint 2 conserva evento 1 y archivos existentes
```

**55 checks PASS**, 28 procesos reales, ocho salidas negativas esperadas. Se ejecutó
primero contra la fuente canónica y después desde copia limpia con **13/13 archivos
idénticos**, reconstruyendo otra vez los tres packs desde Markdown. Sin mocks en
este recorrido. Python 3.14.4 / PowerShell 7.6.5; no instalación, Git, red, servidor,
cuenta de proveedor ni producto nuevo.

| Aserción por modo | Resultado |
|---|---|
| Inventario antes de scaffolding | NEW vacío; EXISTING un archivo sintético |
| Composición y cuatro archivos de bridge | presentes, receipt de materialización real |
| Gate inicial | exit 2 y reporte BLOCKED, sin aprobaciones heredadas |
| Estado/eventos | checkpoint y resume PASS |
| Composición sobre destino ocupado | exit 1; bytes de todo el tooling intactos |
| Bridge repetido | AGENTS.md idéntico |
| Checkpoint duplicado | exit 2; cadena sin cambios |
| Artefacto con hash alterado | exit 2, evidence hash mismatch |
| Restauración y nuevo checkpoint | resume PASS; dos eventos; evento inicial idéntico |
| Archivo previo EXISTING | SHA/BOM conservados |

Reproducir sin servicios externos:

```powershell
./markdown_system/test_agent_entry_lifecycle.ps1 -KeepEvidence
```

Sin KeepEvidence limpia exclusivamente su directorio temporal propio. El test
está integrado en VERIFY_LIBRARY.ps1; el distribuidor admite sólo su nombre exacto.
No agrega un framework, subsistema paralelo ni otro gate de aprobación.

## Evidencia local preservada

- Primera ejecución final: `C:/Users/NL/AppData/Local/Temp/elite-entry-lifecycle-df8f87f815ac44e7905c659d05a437a3/evidence.json`, SHA-256 `e862c12f6fcfbca2059ef99d8d1319fa7c7d8a2f01cfe567ef4576f5b2f85ff0`.
- Copia limpia: `C:/Users/NL/AppData/Local/Temp/elite-entry-lifecycle-430749bda6af4b4ca977db533f89f579/evidence.json`, SHA-256 `77691f9ce4fae8e0f8695932c76e5324d66e624aa16768c42c540ef4ab3ec9cb`.
- Suma de tiempos de procesos: 9.507 ms y 8.997 ms respectivamente. No incluye toda
  la instrumentación ni significa que el agente construya un producto en 9 s;
  no mide tokens, modelo, descarga, compilación del negocio, p95 o p99.
- Copia fuente y log: `C:/Users/NL/AppData/Local/Temp/elite-v287-e6942108180440fa83d0f6acb29cd9f1/`.

La distribución lleva el test reproducible y este reporte, no los fixtures ni
sus checkpoints/aprobaciones. Ambos fixtures conservan `status=BLOCKED`,
`release.status=NONE` y `product_readiness=false`.

## Fallos y correcciones propias

- FAIL-402: Get-Command devolvió dos aplicaciones Python; la conversión a string
  concatenó las rutas. Selección única explícita, sin cambiar PATH; prueba completa
  y copia limpia posteriores PASS. Staging fallido conservado en
  `elite-entry-lifecycle-9bd3dca7ce9e4ca5bf11a12fede09f64` dentro de Temp.
- FAIL-403: fixture NEW declaraba delta_scope_ref de EXISTING; validator rechazó
  correctamente. Fixture distingue modos y no debilita el contrato. Staging previo
  `elite-entry-lifecycle-9226fa02b54446b5bf45928adb1d974d`; ensayo final PASS.

## Preparación real y alcance pendiente

Ronda C generada y explicada; ANSWERED sólo para mantenimiento sin efectos fiscales
o comerciales, con trigger de reapertura ante cualquier proyecto que los realice.
No se eligió un contador, certificado, regla de stock/dinero ni jurisdicción.
Reporte real pasa de 44 a 43 observaciones y sigue **BLOCKED**. Before SHA-256
`20e7982a87738cb37437d344498a06d6cdc615bcf29aca6136c9c9434756a30a`;
after `0a545a7f2b855940a86be9c36ef42f630a5fff7e8fb0cb18c3f11db1e4bab330`.
El before completo se conserva en staging/readiness-before.json y sus observaciones
en V285. Esta reducción no es un porcentaje de terminación.

LIB-R01 gana evidencia integrada de entrada/reanudación/recuperación del tooling.
No se marca completo T2801: faltan las rondas D–H, assurance, authority-map generado
y journeys release-bound del mantenimiento. Tampoco cierra T2802–T2810, los fallos
384–387, el recorrido comercial/documents/IA, la operación ni el ZIP final.
El inventario de producto no cambia: 160 packs/1433 archivos, perfil 67/742.

## Autoridad de método, no falsa autoría

[Google SRE, Release Engineering](https://sre.google/sre-book/release-engineering/)
separa pruebas continuas y pruebas del artefacto empaquetado, con trazabilidad de
revisión; [NASA, Systems Engineering Handbook](https://www.nasa.gov/reference/system-engineering-handbook-appendix/)
separa verificación, validación e integración. Consultados 2026-09-07. Este ensayo
aplica ese criterio a herramientas concretas de la biblioteca; no copia Rapid,
no es código de Google/NASA ni equivale a su aceptación del diseño. La prueba de
un ZIP final exacto sigue pendiente bajo T2810.
