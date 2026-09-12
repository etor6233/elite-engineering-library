# Cloudflare Pingora Network Edge — reconstrucción 2026-08-27 V1

## Veredicto

`cloudflare-pingora-0.8.1` queda `PINNED_CANDIDATE_CONDITIONED`. Es código público oficial de Cloudflare bajo Apache-2.0, útil como referencia real de proxy, load balancing, límites, timeout y hashing consistente. No es un pack composable ni una autorización de despliegue: un test HTTP/2 oficial falla de forma reproducible, el release no publica `Cargo.lock`, la resolución fechada contiene advisories, el workflow del release figura fallido y el ZIP exige preservar dos symlinks POSIX.

## Identidad oficial fijada

| Superficie | Valor verificado |
|---|---|
| repositorio | `https://github.com/cloudflare/pingora` |
| release | `0.8.1`, publicado 2026-06-04 |
| página del release | `https://github.com/cloudflare/pingora/releases/tag/0.8.1` |
| commit | `719ef6cd54e40b530127751bab6c1afc5ae815a8` |
| objeto tag anotado | `8782f3498b061e23d9fe72e6876f72e9a58cabe5` |
| firma | tag y commit sin firma según API oficial de GitHub |
| ZIP oficial por commit | 3.673.628 bytes; SHA-256 `b93a3c6520ff6c542a636b2015709ccf46493aae3e82f4175cafe09717c4c33a` |
| tar.gz oficial | 3.403.806 bytes; SHA-256 `a8e72b383726064aeed08a7c807ef5fe3faeda788fb1ac38de296a524fd8b9cb` |
| licencia | Apache-2.0; `LICENSE` SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| manifest raíz | `Cargo.toml` SHA-256 `a86fd294d9122f84fb4d7f15406fb46268e52028c6f14ae0948b0257d72cd7b1` |

El lock de Elite fija además doce paths Cloudflare críticos —workflow, manifests de crates, server HTTP/2, integración proxy y APIs de load balancing/limits/timeout/ketama— con sus hashes individuales. El código upstream no se copia dentro de los bloques Markdown; se adquiere externamente sólo mediante source profile y aprobación hash-linked.

## Toolchains exactos de la auditoría

| Herramienta | Identidad |
|---|---|
| rustup | 1.29.0, instalador Linux 20.838.840 bytes, SHA-256 `4acc9acc76d5079515b46346a485974457b5a79893cfb01112423c89aeb5aa10` |
| Rust/Cargo | 1.98.0, instalación temporal aislada |
| manifest Rust stable | SHA-256 `3f7d139b73bbbd0004ef6e58b430831c68cdad2b1f64ee2eb35d54c09199489a` |
| CMake | Kitware 4.4.3, tar 64.872.980 bytes, SHA-256 `d6c83076c575bc00b823522ac974bda66d0af05d6ddc30e739c12385cf32c6cc` |
| OSV Scanner | Google 2.5.1, binario SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6` |

El release no publica `Cargo.lock`. La auditoría generó una resolución compatible con MSRV, SHA-256 `6ab25fbbe0c259353eca8f9cb6627b0cb3129f1338138a331eb1eb4aa48fb2fa`, con 479 paquetes. Es evidencia fechada, no un lock upstream ni una recomendación de adopción.

## Resultados ejecutables

- 71 tests focales PASS en TinyUFO, error, ketama, limits, LRU, timeout y doc tests.
- `pingora-core`: 337 PASS, 1 FAIL, 2 ignored.
- fallo reproducible: `protocols::http::v2::server::test::test_req_header_no_eos_empty_data_with_eos`, con `Reset(StreamId(1), CANCEL, Remote)`.
- el resto de core, HTTP y load-balancing seleccionado pasó al excluir únicamente el caso identificado; esto no borra el fallo.
- integración oficial `pingora-proxy/tests/test_basic.rs`: 15/15 PASS después de preparar el origin OpenResty conforme al fixture oficial.
- el workflow GitHub del propio release 0.8.1 figura fallido en la matriz Rust 1.91.1: `https://github.com/cloudflare/pingora/actions/runs/26974601898`; el workflow previo de la rama de release figura verde en `https://github.com/cloudflare/pingora/actions/runs/26973692713`. Ninguno sustituye los gates del target.

OpenResty se obtuvo de paquetes oficiales 1.31.1.1. SHA-256: core `404a8fc4598732145783457299b680a2fc29c5da5f8d5254b51fa05befa7671b`, OpenSSL 3 `eeef574d9babb62e93ac391dbf4f64b62d8fd5ab5f609134d16df3cc69299317`, PCRE2 `7416d3131a711202a5c48397239a522c567e7ac22cebb0ea07f8cfb53a0caa5c` y zlib `ed6029fe8b6002b293090cb9a1fc56f80ddc87a4d70dcb8ebb6860606f9b8027`.

## Dependencias y advisories

OSV Scanner terminó no-cero sobre el lock generado: 11 advisories en ocho package-version.

| Paquete | Versión | Advisory IDs |
|---|---:|---|
| `daemonize` | 0.5.0 | `RUSTSEC-2025-0069` |
| `derivative` | 2.2.0 | `RUSTSEC-2024-0388` |
| `lru` | 0.16.4 | `RUSTSEC-2026-0253` |
| `paste` | 1.0.15 | `RUSTSEC-2024-0436` |
| `protobuf` | 2.28.0 | `RUSTSEC-2024-0437`, `GHSA-2gh3-rmm4-6rq5` |
| `quinn-proto` | 0.11.14 | `RUSTSEC-2026-0185`, `GHSA-4w2j-m93h-cj5j` |
| `rustls-pemfile` | 2.2.0 | `RUSTSEC-2025-0134` |
| `time` | 0.3.45 | `RUSTSEC-2026-0009`, `GHSA-r6v5-fh4h-64xc` |

El reporte JSON tiene 113.583 bytes y SHA-256 `a60785256f7c1edec3ab9d030114c037d4df3a7641433fedabac340189e59a7c`. La presencia de un advisory no demuestra por sí sola exploitability del producto final, pero prohíbe omitirlo o promover este grafo sin lock propio, reachability, actualización compatible y repetición integral.

## Fidelidad del archive y plataforma

El ZIP oficial contiene exactamente dos directorios representados como symlinks Unix, ambos con modo `0120777`:

| Path | Target exacto |
|---|---|
| `pingora/tests/keys` | `../../pingora-core/tests/keys` |
| `pingora-proxy/tests/utils/conf/origin/conf/keys` | `../../keys` |

`.NET ZipFile` sobre Windows los proyectó como archivos regulares con el target como contenido. Esto originó `LIB-FAIL-434`. `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.42 ahora:

1. declara los dos links y sus targets en el lock;
2. valida modo Unix, inventario y target 1:1 antes de extraer;
3. rechaza Windows antes de descargar o crear destino;
4. en Linux/macOS reconstruye los symlinks y verifica `LinkType`/`LinkTarget`;
5. registra `archive_symlinks_verified` en el receipt;
6. rechaza targets que escapen del root y conserva una regresión negativa.

La suite reconstruida desde Markdown pasó `UPSTREAM_ACQUISITION_TEST_PASS negatives=8`, `SOURCE_PROFILE_TEST_PASS valid=13 negatives=5 positives=1`, y el probe Windows demostró rechazo antes de crear destino.

## Condiciones para cualquier proyecto

Antes de adoptar se requieren, como mínimo: Linux/macOS aprobado; Rust/C/CMake/TLS fijados; lock y SBOM revisados; resolución de advisories; suite upstream total verde; TLS/certificados/mTLS; header/body/protocol policy; discovery/health/retry/idempotencia; overload/rate/concurrency; aislamiento/autorización; tráfico hostil; load/soak/perfiles de memoria; observabilidad sin datos sensibles; canary, rollback, recovery e incident owner.

Las cifras públicas de tráfico o latencia de Cloudflare describen su entorno. No son evidencia transferible ni objetivo demostrado para otro sistema.

## Resultado de admisión

- código oficial de empresa líder: **sí**;
- identidad/licencia y artefactos exactos: **sí**;
- adquisición inmediata fiel en Windows: **no, rechazo deliberado**;
- fuente condicionada disponible para un target Linux/macOS: **sí**;
- composición o despliegue inmediato: **no**;
- `REUSABLE_PACK`: **no**.
