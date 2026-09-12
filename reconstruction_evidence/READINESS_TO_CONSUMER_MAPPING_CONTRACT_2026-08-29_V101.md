# Readiness-to-consumer mapping contract — evidencia V101

Fecha: 2026-08-29  
Estado: `REUSABLE_PACK / REBUILD_VERIFIED`

## Procedencia

`PROJECT-START-READINESS-VALIDATOR` 0.3.0 continúa siendo orquestación `AUTHORED` de Elite Engineering Library. No es código atribuido a Red Hat, Google ni otra empresa. El runtime de mapping sigue siendo el artefacto oficial Google-origin CEL-Java 0.13.0 dentro del consumer Debezium/Quarkus condicionado; V101 sólo elimina traducción manual entre una decisión aprobada y su configuración.

## Brecha cerrada

V100 exigía hashes e identidad por clase, pero no conservaba los bytes de la expresión, las claves exactas de salida ni el modo arquitectónico. Un agente habría tenido que completar esos valores después del readiness y podía introducir una autoridad distinta.

V101 exige, para cada clase con storage automático:

- `persistence_mode`: `TRANSACTIONAL_DIRECT` u `OUTBOX_CDC_CONSUMER`;
- exactamente los diez campos `mapping*` presentes en `debezium_postgres_inbox_consumer/project-profile.template.json`, sin override adicional;
- `mappingDomainType == class_id`;
- expresión UTF-8 trimmed/single-line, máximo 4.096 caracteres, con SHA-256 recomputado y coincidente;
- mapping ID, semantic version, schema IDs, tres hashes y output-key CSV bajo las mismas regex que el constructor Java;
- evidence paths de corpus, evaluación y aprobación;
- database/observability/backup/privacy para modo directo; además outbox/broker para CDC consumer;
- cero configuración de persistencia latente en clases `REVIEW_ONLY` o no autorizadas.

## Gates observados

- source tree: CPython py_compile y 33/33 unit/CLI PASS;
- materialización Markdown: 4/4 hashes exactos y segundo 33/33 PASS;
- fixture positivo outbox/CDC y fixture positivo direct: PASS;
- negativos: campo extra, gramática/ID/version/schema/output key, expression SHA/newline, cross-class, capabilities faltantes y dormant review config: PASS;
- auditor global compara dinámicamente `MAPPING_FIELDS` de readiness con los campos `mapping*` del perfil materializado del consumer: igualdad exacta PASS;
- `VERIFY_LIBRARY_PASS`: 71 packs/682 archivos/427 Markdown antes de esta evidencia y 36 perfiles;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 71 packs/111 fuentes, readiness 33 tests y contrato cruzado readiness↔consumer PASS.

Hashes focales:

- pack Markdown: 69.494 bytes, SHA-256 `71bdc20b89fc54d8e3a6da4d3868e520f0d0da61267d41962f0256627edd8e00`;
- template cerrado: `1cddd8f25aacfc3d36872619c22c631eab21c97e9f511029834d12aaa99aa5d1`;
- validator: `bc1565c07afba1d819dced03e3a898f4761552af694b31f527a3024dbd958278`;
- tests: `5f1c7dec4eebd3abbb8c0f86925871ff17926cf70172bfa9bef3faf67571888f`;
- README materializado: `93771a9205ba5d7bade1e7bc537dc539eab8058ffb7bf16f95164cc5eb132a04`.

## Límites

El gate no crea una expresión, schema, corpus, aprobación ni infraestructura. Exige que existan y coincidan antes de abrir el build. Tampoco convierte el modo directo en código ya implementado para todo dominio: el pack/arquitectura seleccionado debe demostrar su propia transacción. `OUTBOX_CDC_CONSUMER` se apoya en la lane V99 condicionada y conserva todos sus gates Kafka/PostgreSQL/live.

Memoria al cierre focal: 1.105 fallos locales + 185 condiciones upstream = 1.290 IDs únicos, cero abiertos locales.
