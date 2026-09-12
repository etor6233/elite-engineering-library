# Readiness → consumer profile renderer — V102

## 1. Resultado y procedencia

`PROJECT-START-READINESS-VALIDATOR` 0.4.0 añade un archivo de orquestación local `AUTHORED`: `render_consumer_profile.py`. No es código copiado ni atribuido a Google, Red Hat, Debezium, AWS, Microsoft u otra empresa. El runtime de mapping continúa siendo Google-origin CEL-Java 0.13.0, importado por la adaptación declarada del consumer Debezium/Quarkus; V102 sólo elimina una transcripción manual entre evidencia aprobada y configuración.

Antes de esta evidencia, la biblioteca contenía 428 Markdown. El inventario materializable pasó de 682 a 683 bloques: 549 `AUTHORED`, 30 `ADAPTED` y 104 `VERBATIM`. No cambió el inventario de 111 fuentes oficiales ni sus licencias.

## 2. Contrato ejecutable

El renderer:

1. vuelve a ejecutar la validación integral del readiness y bloquea cualquier record no `READY_TO_BUILD`;
2. selecciona exactamente un par `{class_id, variant_id}` autorizado como `READY_FOR_AUTOMATIC_STORAGE` y `OUTBOX_CDC_CONSUMER`;
3. exige los diez campos mapping exactos que el consumer materializado conoce;
4. acepta sólo el schema exacto de un template consumer nuevo, no reconocido, sin configuración, owners, approvals ni proofs;
5. rechaza claves JSON duplicadas, rutas no canónicas, traversal, symlinks, schema alterado, preconfiguración y output existente;
6. copia únicamente mapping ID/domain/version, expresión+SHA, output keys y ambos schema ID+SHA;
7. conserva `acknowledgeConditionedState=false`, todos los proofs en `false`, referencias operativas vacías y cero secretos.

El output no autoriza deploy ni persistencia. Endpoints, secret references, owners, approvals, migración, TLS/auth, corpus, carga, observabilidad, recovery y rollback siguen requiriendo evidencia del proyecto; el renderer no los inventa.

## 3. Reconstrucción y gates fechados

- source staging CPython 3.14.4: 40/40 unit/CLI PASS;
- materialización fresca desde el Markdown: 5/5 hashes exactos y 40/40 PASS;
- roundtrip source↔Markdown exacto para template, validator, renderer, suite y README;
- `VERIFY_LIBRARY_PASS`: 71 packs, 683 archivos materializables, 428 Markdown antes de esta evidencia, perfil readiness 1/5 y 36 perfiles válidos;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 71 packs/111 fuentes, readiness 40 tests, comparación dinámica de los diez campos mapping y schema completo renderer↔template consumer real PASS.

Hashes SHA-256 finales:

| Artefacto | Bytes | SHA-256 |
|---|---:|---|
| `implementation_packs/PROJECT_START_READINESS_VALIDATOR.md` | 83,533 | `618ad4fa695b5f197a30c4c6cb17f6dc337e33a8a7e836e34fc79efb2ee0efa3` |
| `project-readiness.template.json` reconstruido | — | `1cddd8f25aacfc3d36872619c22c631eab21c97e9f511029834d12aaa99aa5d1` |
| `validate_project_readiness.py` reconstruido | — | `bc1565c07afba1d819dced03e3a898f4761552af694b31f527a3024dbd958278` |
| `render_consumer_profile.py` reconstruido | — | `16978a72f46428f083062fa61062d37d806baefd4ae53d5fe8f8cd9694a2d53a` |
| `test_validate_project_readiness.py` reconstruido | — | `4ca12c15e94837a9c4d7ebcc6b07ba453358d4232f9f7b0b7a1357fa48db23ed` |
| `README.md` reconstruido | — | `6ea7065a1383486a88d2af0c149378ac065fb6f32aa87e5e03938494a6137608` |
| `markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md` | 1,576 | `ad130fb8063d6ad6bdba7067fba6332119847948a1a2fd12655a14e6a18b7830` |
| `VERIFY_EXECUTABLE_LIBRARY.ps1` | 118,197 | `f69bcb5d1939d9d4a52601fcae966191761bebd091d14c5d75945bb8dc3edbad` |

## 4. Límite honesto

V102 reduce tiempo y error de copia, pero no crea un mapping empresarial ni vuelve verdadero un corpus. La expresión, schemas, clases, variantes, métricas, aprobación y accesos deben provenir de autoridad/evidencia real del proyecto. Un `TRANSACTIONAL_DIRECT` no se fuerza dentro del consumer outbox; el renderer lo rechaza. El perfil generado sigue cerrado hasta que sus gates operativos reales sean demostrados.
