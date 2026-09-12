# Evidencia V237 — identidad de canal PostgreSQL conectada al lead

Fecha: 2026-09-04  
Pack: `GO-PG-CONTACT-CHANNEL-IDENTITY 0.1.0`  
Pack SHA-256: `961c5e78619d0be9abd2498018e5f69c310795fc89d81bbc2155598d52acf916`  
Plan completo SHA-256 tras selección: `083c74b90a121f0549de6060c6658e9276264e2efdebd9486416c777026bff07`

## Claim demostrado

La composición ya no depende de un resolver ficticio para unir contacto y conversación. Ocho archivos `AUTHORED` materializan:

- binding exacto tenant+canal+identidad externa hacia un lead CRM existente;
- HMAC-SHA-256 con clave mínima de 32 bytes; la identidad externa no se persiste en claro;
- organización obtenida por join con el lead, nunca desde texto ni salida del modelo;
- alta, actualización de política/PII y revocación mediante versión esperada;
- lead y subject inmutables para impedir reasignación silenciosa;
- ledger de decisiones append-only con request hash y replay exacto;
- conflicto explícito ante dos altas concurrentes;
- implementación real de `conversationruntime.ContactResolver`;
- E2E lead/binding PostgreSQL → canal → Responses tool → gateway autenticado/idempotente → cotización → turno/reply replayable.

## Autoridades y procedencia

- NIST SP 800-63 Rev. 4 gobierna solamente el claim estrecho de sujeto único, registro de evidencia/consentimientos y asociaciones de identidad: <https://pages.nist.gov/800-63-4/sp800-63a/accounts/> y <https://pages.nist.gov/800-63-4/sp800-63/model/>.
- PostgreSQL 18 gobierna aislamiento serializable, concurrencia y failure handling: <https://www.postgresql.org/docs/18/mvcc.html>.
- Los ocho archivos son `AUTHORED`; no son código textual de NIST, PostgreSQL, Meta, Google, TikTok ni OpenAI.

## Gates reproducidos

- round-trip del pack: 8/8 SHA-256 idénticos;
- `gofmt` sobre cinco archivos Go: no-op;
- PostgreSQL 18.6, base limpia `elite_contact_v237`: 47/47 migraciones `up`;
- SQL focal 0047 con `ON_ERROR_STOP=1`: PASS; ledger inmutable y ausencia de identidad cruda;
- suites focales `internal/app`, `internal/contactidentity`, `internal/platform/postgres`: PASS dos veces consecutivas;
- suite completa: 49 paquetes, `go test -count=1 ./...` PASS;
- `go vet ./...`: PASS;
- `go build ./...`: PASS;
- `VERIFY_LIBRARY_PASS`: 149 packs, 1311 archivos, 656 Markdown, 44 perfiles; perfil completo 57 packs/620 archivos.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 149 packs, 121 fuentes y 14 adapters de proveedor; los gates live declarados por cada pack permanecen explícitamente omitidos cuando faltan cuenta, costo o target autorizado.

## Fallos convertidos en regresión

- `LIB-FAIL-1847`: un ledger inmutable no se limpia con DELETE; los tests usan identidades únicas y la base descartable se reinicia por ciclo controlado.
- `LIB-FAIL-1848`: UUID único no basta si `tenant_code` queda constante; ambas claves se hacen únicas.
- La suite doble demuestra que la corrección no depende de una primera ejecución limpia.

## Condiciones conservadas

El pack no prueba propiedad live de un número/cuenta, método de verificación del proveedor, consentimiento jurídico, rotación de HMAC, entrega outbound, IdP, carga, pentest ni producción. Esas decisiones se cargan en el proyecto y no se infieren por coincidencia de teléfono, email, nombre o texto.
