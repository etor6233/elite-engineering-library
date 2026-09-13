# V402 — conversión FX conectada al asiento durable

T2802, claim local acotado: snapshot FX admitido → recibo inmutable → borrador en el escritor contable existente → posting explícito → balance → reversión. GO-EXACT-FX-SNAPSHOT-ACCOUNTING0.2.0 y GO-ENTERPRISE-ACCOUNTING-LEDGER-API0.1.2. Siete archivos nuevos y cuatro modificados; 178packs/1804bloques en biblioteca, franquicia82/1039. No cierra T2802 entero, TEST02 ni READY global.

El contador autorizado aporta cuentas, período y descripción. El vínculo toma importe positivo, moneda local y fecha del recibo exacto; no inventa tasas, cuentas, ganancias/pérdidas, valuación ni política fiscal. La conversión conserva CONVERSION_RECEIPT_ONLY. El nuevo recibo histórico es FX_JOURNAL_PREPARED; current_status/current_version consulta el asiento actual. El posting y reversal conservan sus endpoints, permisos y algoritmo. No se agrega otro ledger. El SQL original de CreateJournal se extrae mecánicamente a createJournalInTx para compartir su transacción con el recibo, idempotency_record y outbox.

Prueba conectada única nueva: HTTP real con RS256/JWKS y PostgreSQL18.6, 62migraciones. EUR1000minor→USD1250minor; doce preparaciones concurrentes producen un solo asiento y recibo. GET recupera la respuesta perdida. Payload/actor/organización divergentes, segundo request key sobre el mismo recibo, permiso sólo lectura, source FX_CONVERSION falsificado, importe JSON inyectado/duplicado, cuenta inválida/inactiva, moneda no local y cero/negativo se rechazan. Una preparación cancelada con outbox bloqueado no deja asiento ni idempotencia; no se afirma haber observado directamente su espera SQL. Preparar no crea entradas contables. Posting requiere permiso separado, balance ±1250 y reversión existente da neto0. Replay tras retirar idempotencia transitoria devuelve el mismo recibo; historia inmutable. Período cerrado rechaza nueva preparación. Resultado final:2asientos y4entradas, sin skips.

Build completo y vet de paquetes afectados PASS. Reconstrucción exacta: franquicia82/1039, backend40/561, serverless57/786 y HTTP83/1047. Los tres perfiles Go compilan los paquetes contables/PG/HTTP y sus tests; ese comando ejecuta cero pruebas. No se repiten las suites de pagos, navegador o aritmética sin delta. Dos comentarios de fuente se corrigieron después del test para describir el nuevo alcance; no cambió código ejecutable. Los algoritmos BC y el fuzz aritmético previo se conservan byte-idénticos.

Downgrade ejecutado contra la base del ensayo: rechaza borrar el recibo existente y preserva su snapshot y las4entradas byte a byte. En una base independiente vacía, down→up PASS. FAIL822 fue un error del comprobador al asumir entry_no; el esquema real usa entry_id. Falló en SELECT antes del downgrade, se detuvo PostgreSQL y se conserva el log. Sólo se repitió ese probe corregido; ambos servidores propios terminaron. No se borra historia para forzar rollback.

## Admisión del claim y procedencia

| Gate | Evidencia y límite |
|---|---|
| G0 identidad | Mismos snapshots oficiales BCApps2eae56d704a1fd035d104f333602aea7091b7749 y owners exactos ya admitidos; manifest antes/después del glue |
| G1 licencia | MIT/avisos BC originales idénticos; siete nuevos archivos AUTHORED LicenseRef-Workspace-Owner, ninguna atribución a Microsoft |
| G2 autoridades | Precisión monetaria existente, cuentas/período elegidos explícitamente; no algoritmo fiscal/valuation nuevo |
| G3 arquitectura | Un único writer contable; vínculo inmutable y transacción común con idempotencia/outbox; actor y scope del recibo |
| G4 corrección | Journey conectado nuevo, concurrente y negativos materiales enumerados arriba |
| G5 seguridad | Permisos separados, scope, tampering y JSON exacto probados; SCA/seguridad global de composición se resuelven en T2803 y no se heredan de este PASS |
| G6 rendimiento/resiliencia | Concurrencia12→1, rollback por cancelación, recuperación durable tras idempotencia transitoria; no claim de carga target |
| G7 operación/evolución | Doc materializado, proyección estado actual, downgrade no destructivo probado, vet/build; despliegue/ops generales mantienen owners |
| G8 pack | Cuatro perfiles reconstruidos exactamente; siete archivos añadidos seleccionados por los planes; referencia externa1039archivos |

Los adaptadores de representación/SQL/HTTP son AUTHORED glue. ExchangeExact/FindLast y balance/post/reversal siguen con sus propias derivaciones; no se atribuye a la empresa upstream el código local. No se introducen dependencias. La ausencia de credenciales no tapa trabajo de seguridad/operación aún pendiente en la composición. El claim cerrado es posting de una conversión positiva a moneda local con cuentas explícitas, no un motor universal de diferencias de cambio.

## Recibos

Staging: C:/Users/NL/AppData/Local/Temp/elite-v402-library-infra.

| Archivo | SHA-256 |
|---|---|
| `FX_JOURNAL_CONNECTION_ADMISSION.candidate.md` | `0176bbd8dd61d8414af2e770507737814d3a7db1403456c42526dc54af4b850c` |
| `fx-journal-pg-first/result.json` | `b288cdacbfe13ae6db5d2d90a35f871e4cd87338d11d7d1949c7c566effaf371` |
| `fx-journal-pg-first/connected.log` | `a69430a68ebad82bea97dbe976e691a80b5e2ba0a276a8d23263491a48338a66` |
| `fx-journal-build-gates.json` | `53a44c110431b44fb65649f25a66c939966347780de4d66ee5e263ebad947597` |
| `fx-journal-downgrade/result.json` | `426265fd620426314b533d97bd15ef458c0c17d082cb09f6de66284e843056e5` |
| `fx-journal-downgrade-corrected/result.json` | `1e6f00353efbe802fa7f97537c3aae6d987aedadac11272bcb5aaa2ac6fa40e7` |
| `fx-journal-package-result.json` | `81ece03f4443265151a5b11efb79c99a91b6e632029a51e98dc6a14a494805b6` |
| `fx-journal-composition-stage.json` | `6a51adad8f7abaa34cf5b42a9afc2dcc70a50b95e6cf7c4198baa01b06ec5f28` |
| `fx-journal-reconstruction.json` | `9e21c0bc11d420b688ec7507f8f65236a889820c82d190e87dba09746bac5c78` |
