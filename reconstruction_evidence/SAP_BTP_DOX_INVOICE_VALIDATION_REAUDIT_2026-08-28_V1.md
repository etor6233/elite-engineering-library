# SAP BTP DOX Invoice Validation — reauditoría 2026-08-28 V1

## Veredicto

`SAP-samples/btp-cap-dox-invoice-validation` queda `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`. Contiene código oficial útil de roles, asignación de validator, correcciones e historial de decisión, pero no constituye un componente reproducible, seguro ni listo para composición inmediata. No se copió ni se atribuyó a SAP ningún arreglo local.

## Identidad y licencia exactas

- Repositorio oficial activo: `SAP-samples/btp-cap-dox-invoice-validation`.
- Branch `main`; commit firmado/verificado `990f4d8650c12d39c9a84974091c1b248a549ab6`; tree `baf1a25d096d6ed294d96d23b248f70d94dc7650`; fecha 2026-06-11.
- Cero tags y cero releases al corte.
- Archive exacto: 3.666.469 bytes; SHA-256 `4dc23ccc7c1eef78513bf287cce7be2ea5de4a1a2ea38406b22ea02f8b472765`; 190 entries, 143 files, 46 dirs y cero paths inseguros.
- Apache-2.0: `LICENSE`, 11.357 bytes, SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`.
- `REUSE.toml`: 1.987 bytes, SHA-256 `1024807fe40a9f9f7daa751e19a979e71e5a49025a88352833d0bdc76e923390`. Aclara que Apache no concede las llamadas a productos SAP/terceros; sin licencias aplicables, esas llamadas quedan limitadas a evaluación interna no productiva/no comercial.

## Artefactos oficiales auditados

El source lock conserva SHA-256 exacto de 13 archivos: README, REUSE, manifests root/API/router/UI, servicio y autorización CAP, router CSRF, requirements y server snapshot, `mta.yaml` y documentación de despliegue. El lock no contiene una reinterpretación ni una implementación inventada.

## Reproducibilidad y builds

- No existen lockfiles source, workflows CI ni tests automatizados; sólo un archivo HTTP manual.
- Root declara Node `^22`; la documentación oficial dice Node 20 y el router resuelto declara soporte hasta Node 20.
- Runtime auditado: Node 22.23.2/npm 10.9.8. Asset Windows SHA-256 `1177b4137ba5adaa56354ae40f1080c7450e8ae09cecb47da459d1c52ac99f97`; SHASUMS oficial SHA-256 `778ac5b2fcdbd68d9c0ae9f4310674faa3af0910bd0d18e7f6597787c40a3e39`.
- `npm run win:setup` falla en `generate:types`; la ejecución explícita posterior del generador pasa, pero no convierte el setup limpio en reproducible.
- `npm run build:cf` del API falla porque `cds-ts` no pertenece al grafo instalado; la guía exige tooling global sin pin.
- UI `tsc` + Vite pasa, 1.860 módulos, con chunk principal aproximado de 1,95 MB y warning de tamaño.
- npm audit de grafos generados: root 0, API 0, UI 5 (4 moderate, 1 high), router 17 (7 low, 5 moderate, 5 high). Upstream también resuelve paquetes SAP declarados fuera de soporte.
- Los cuatro lockfiles fueron generados sólo para auditoría, no atribuidos a SAP: root 138 paquetes/SHA `6d50f90f…91a1`; API 208/`d7169304…77e6`; UI 228/`7069e207…c516`; router 236/`f2518e3c…e2ed`.

## Snapshot service Python

- Upstream fija 38 dependencias directas y Python 3.11, pero no un lock transitivo; `uvloop==0.17.0` impide instalación Windows.
- En WSL2/Linux, un entorno 3.11.16 sin seed instaló los 38 pins; `uv pip check`, compileall e imports `model`/`server` pasaron.
- Fixture de auditoría SHA-256 `3c31ddc4a0cc62cc9865137fc60e5aaa5a55b2018f54b0894380e31c6b1a9721`.
- `uv.lock`: 41 paquetes, 69.443 bytes, SHA-256 `56d5da51bc8d4c9a3d761c92b6d9389279b80cfd07407073389b06d7be20ad3d`.
- `uv audit --locked --python-version 3.11 --python-platform linux`: 40 auditados, 160 advisories, 14 paquetes afectados. La mayor concentración está en pypdf (73), Pillow (37), Starlette (14) y Jinja2 (10); no se parcheó upstream.

## Contrato humano: valor real y bloqueantes

Valor oficial comprobado:

- identidades `managed`, users por email y roles por proyecto;
- validator actual asociado a invoice;
- tablas de correcciones/versiones con reason;
- forwarding a rol distinto;
- accept/reject final limitado a accounting/current validator;
- UI de correcciones e historial.

Bloqueantes comprobados:

- cero ETag, CAS u optimistic concurrency y cero lease/expiry de assignment;
- `setCV(projectId, invoiceId)` no demuestra que la invoice pertenezca al project solicitado;
- cero tests automatizados del flujo;
- CSRF desactivado en rutas API y snapshot;
- destinos `strictSSL:false`;
- snapshot FastAPI sin autenticación/autorización, payload base64 sin límites/gates de malware y escrituras concurrentes a `rendered.html`, `rendered.pdf` y `snapshot.pdf` compartidos.

## Admisión

No se materializa ni compone este source en un sistema. Queda como evidencia oficial exacta para comparar futuras revisiones SAP. Una reevaluación exige release fijada, locks oficiales, builds/tests/SCA verdes, autenticación y aislamiento del snapshot, TLS/CSRF seguros, request binding, CAS, lease, corpus autorizado y pruebas live con cuentas/licencias/costo aprobados por el usuario.
