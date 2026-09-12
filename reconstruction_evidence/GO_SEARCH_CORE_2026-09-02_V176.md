# Go Search Core — V176

## Resultado estrecho

V176 materializa `GO-SEARCH-CORE 0.1.0`, un índice de búsqueda gobernado por tenant: full-text PostgreSQL (`tsvector`/GIN + `websearch_to_tsquery` + `ts_rank`, configuración `simple` fija con peso title A / body B), filtros de `kind`/`facets`, almacenamiento de embeddings `real[]` y búsqueda vectorial por coseno en un core Go stdlib-only. Negación por defecto de tenant en toda query.

No afirma recall universal de un modelo de embeddings ni escalado ANN. Para escala se fijan Faiss/Meta `7059eaf7da7eddda62e71367e684d4bdedd7f94f` y DiskANN/Microsoft `860cf47bc11b6c2b938818773831a9374553d976` (MIT) como backends de retrieval autorizados, no materializados en este pack.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED`, sin dependencias de terceros.
- PostgreSQL 18.6 full-text search (PostgreSQL License) — migración `AUTHORED` sobre `tsvector`/GIN/`ts_rank`.
- Backends ANN autorizados (condicionados): Faiss/Meta y DiskANN/Microsoft, ya pinnados en `NLP_RAG_RETRIEVAL_DATA.md`.

## Archivos materializados (11)

| Archivo | SHA-256 |
|---|---|
| internal/search/model.go | 955eb3492c284cbcc69590c7a98a3398b8d5d34ded6e7c831976cea328cda901 |
| internal/search/cosine.go | ed136fa4504cf069a83f60cc07bf35ad094ca405c68dfb542f65a6ae8222312c |
| internal/search/fulltext.go | 0018e2cc8c396a49e0390cd99f4a12a73775e7a457f9702465fd26a4c5ea8a38 |
| internal/search/store.go | 6c080f74820ccfa0679786706d2d13d04f4a9c9f368f4e2ac0abee7fcd8a5a44 |
| internal/search/memory.go | 765217bf4874b67cdf23bbb8bc6bc01ac2573c1a1f80b736a010951050c91609 |
| internal/search/model_test.go | ba2875badd4e8608976b41071704af391c129b9bffa78a96e9d4db33cb55f79b |
| internal/search/cosine_test.go | 5c86909c0615ea7a0b545324b7c04a8ec9c5465883534998ca553cdf56189be7 |
| internal/search/memory_test.go | a3a5fad0ab1530699fb2fdbd970ff70ee8c4dfcaab73239fc40142a538c7f0cc |
| db/migrations/0040_search.up.sql | 1c407427f5d9cbeae0d800e50db3edd208e91030695fabf9c2945c3dbf43c8d1 |
| db/migrations/0040_search.down.sql | a6ed0dbcd81c88d0dcc541d3830cc94ee231778d54a9f503cf6cac194889cd52 |
| db/tests/0040_search.test.sql | 513527780cf4956f73f4dd236e78128cc5e31ca4b586908e481d12adfb767995 |

SHA-256 del pack: `e58163b599179b3c2fd133b3cbebee61ef1348d4dd0472b8d3bd0a2fced803a4` (32.393 bytes).

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).
- PostgreSQL 18.6 x86-64 (cluster aislado, puerto 55441, `initdb` limpio).

## Verificación ejecutada

- `go test ./internal/search/ -count=1`: **11/11 PASS** (validación con negativos, coseno ortogonal/idéntico/opuesto/dim/cero, ranking title>body, aislamiento tenant/org, filtro kind/facets, top-K/umbral vectorial, delete not-found/denied).
- `go test ./... -count=1` (aifoundation + search): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 11/11 bloques reproducen byte a byte los archivos probados; `go test` sobre el árbol materializado PASS.
- PostgreSQL 18.6 real (initdb → start → up → test → down): `ON_ERROR_STOP=1` exit 0. Cinco aserciones DO pasan: full-text 'scooter' en (tenant-a,org-1) → 2 con `d1` top por peso A; 'urbano' → 1; aislamiento tenant-b=1/org-2=1; filtro `facets @> '{"category":"movilidad"}'` → 3; round-trip embedding `real[]` dim 3 valor 1.0. Down deja `to_regnamespace('search') is null = true`.

## Cobertura de invariantes probada

1. `tenant_id` obligatorio (deny-by-default): query/delete sin tenant → `ErrDenied`.
2. Full-text determinista con configuración `simple` fija y ranking por peso (title A > body B).
3. Aislamiento multi-tenant y por organización en toda búsqueda.
4. Filtros de `kind` y `facets` aplicados sin ampliar scope.
5. Embedding `real[]` con dimensión consistente; vector cero o dimensión distinta no produce score inventado.
6. Migración atómica up/down sin estado parcial.

## Condiciones residuales

- Adapter `pgx` (PostgreSQL Store) y su prueba de integración: pendientes (driver no disponible sin red/cache en este host).
- Modelo de embeddings real y su recall/corpus: decisión de proyecto.
- Faiss/DiskANN para escala ANN: expediente separado si el workload lo exige.
- Indexación automática desde el catálogo/conocimiento (pipeline de republicación): pendiente.
- RAG-AGENTS (consumidor) y CACHE: aún no materializados.

V176 añade la superficie `SEARCH` ejecutable; no declara producción ni ANN a escala.
