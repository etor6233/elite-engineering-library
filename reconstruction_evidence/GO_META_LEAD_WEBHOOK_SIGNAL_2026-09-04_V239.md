# Evidencia V239 — señal Meta Lead Ads firmada y durable

Fecha: 2026-09-04  
Pack: `GO-META-LEAD-WEBHOOK-SIGNAL 0.1.0`  
Pack SHA-256: `27e60c84a69af72cd5d9678db17a5525905b3e43a2b893b8c73d877f9050162c`  
Plan completo SHA-256 tras selección: `7d1c38cae8f4c7253b68c12745d2764ea08a63d23a8140ec152a6be0c95b7451`

## Fuentes oficiales exactas

- `fbsamples/messenger-platform-samples@354ee221ac1d081cc6105a1515a8468cc44f6710`, tree `a711d923d3797d202cdeaf415d7de088da0ee424`, source ZIP 7.674.719 bytes/SHA-256 `65fc927d030d97ba7b8cc8f891fc5369a89b820e73326ee12703b9052af5202c`. Los paths focales implementan `X-Hub-Signature-256`, compare constante, challenge y regresiones de seguridad. Licencia restringida al uso con APIs Meta/Facebook.
- `fbsamples/lead-ads-webhook-sample@c0843165ad0f6acd82731e0c9da96a22b4162b35`, tree `93c3556abd9bee22bb92db66673e75d67e0f2449`, source ZIP 7.604.183 bytes/SHA-256 `8da8485a8df58aa7ebf9e004e9ab3ebb44e019bc426be6b535ab0646785a67c8`, MIT. El Postman/controller gobierna sólo el envelope `page/entry/changes/leadgen` y el handoff por `leadgen_id`; el repositorio está archivado y su controller no se adopta como seguridad productiva.

## Claim demostrado

Nueve archivos reconstruibles aportan endpoint ejecutable Go, challenge, firma SHA-256 sobre bytes crudos, recorrido de todas las entradas/cambios, límite de body, lote y signals PostgreSQL append-only, dedupe/replay tenant+organización scoped, divergencia fail-closed y outbox `meta-lead.signal-received`. El ACK `200` sólo ocurre después del commit. El signal queda `pending_retrieval`; no se presenta como lead completo, consentimiento ni venta.

Procedencia: dos archivos `ADAPTED` y siete `AUTHORED`. No contienen código textual de producto Meta ni afirman que el glue local sea Meta-authored.

## Gates reproducidos

- materialización/round-trip: 9/9 SHA-256 idénticos;
- tests focales leadstream/PostgreSQL/command: dos ejecuciones PASS;
- firma ausente/alterada, challenge erróneo, content type, tamaño, JSON, objeto, batch, replay, divergencia y fallo de store: PASS;
- PostgreSQL 18.6 limpio `elite_meta_v239_a`: 49/49 migraciones y SQL 0049 PASS;
- full suite: 51 paquetes `go test -count=1 ./...`, `go vet ./...` y `go build ./...` PASS;
- `VERIFY_LIBRARY_PASS`: 151 packs, 1.331 archivos, 660 Markdown, 44 perfiles; franquicia 59/640;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 151 packs/121 fuentes del lock central y todos los componentes focales offline; gates live explícitamente omitidos por falta de autorización/target, no simulados.

## Condiciones que no se simulan

App/business/cuenta, permisos Lead Retrieval, Page/form ownership, suscripción `leadgen`, test lead oficial, HTTPS/CDN/WAF, delivery/retry reales, cuotas, política de PII/consentimiento y la ejecución signal→SDK fetch→import→promotion→journey deben probarse en el proyecto. TikTok conserva un gate independiente.
