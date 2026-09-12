# Kubernetes Deploy/Rollback Admission — V112

Fecha: 2026-08-29. Resultado: `REBUILD_VERIFIED / CONDITIONED`.

## Autoridad oficial exacta

- Kubernetes `v1.37.0`, release GitHub `377265881`, publicada `2026-08-26T16:29:12Z`.
- tag object `157e582fcc3ebba3c22b16721f49d6890f784c1f`; commit `f54c212e3a2f75d674b717a9b29052b20b60aefc`; tree `95340e46b29b51d6779c5bf402fc254171ce9786`.
- ZIP exacto del commit: 70.912.415 bytes; SHA-256 `19ef12d9f8027d6212171adaf8f99175189eee68f88c336dff8f27fb72689efc`.
- `LICENSE` Apache-2.0 SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`; `CHANGELOG/CHANGELOG-1.37.md` SHA-256 `0a895234e7308c404a4c3e02aff4b0cecf4d05c99c850495cce35f41fd0bae66`.
- `kubectl.exe` Windows amd64: 63.456.768 bytes; SHA-256 publicado y verificado `4721b614a67bb4932a0369e61f4a323d8c6ca00943d3a2ff14837c124da06f0e`; runtime `v1.37.0`, commit exacto, tree clean, Go 1.26.6, kustomize 5.8.1.
- Fuentes: `https://github.com/kubernetes/kubernetes/releases/tag/v1.37.0`, `https://kubernetes.io/docs/concepts/workloads/controllers/deployment/`, `https://dl.k8s.io/release/v1.37.0/bin/windows/amd64/kubectl.exe.sha256`.

El tag anotado y el commit aparecen unsigned; no se afirma firma. La fuente y el binario oficial no se embeben en secure operations: el pack de adquisición los fija y el proyecto los adquiere bajo aprobación.

## Implementación local declarada

`SECURE-OPS-DELIVERY-CORE` 1.0.0 añadió tres archivos `AUTHORED`: template, validador y regresiones de `DEPLOY_ROLLBACK`. La semántica exige siete ejecuciones ordenadas y hash-bound sobre un Kubernetes `Deployment` real: preflight de cluster/context/namespace/autorización; aplicación de digest inmutable; canary; rollout completo; rollback realmente ejecutado; restauración del digest anterior con salud e invariantes; y compatibilidad de migración expand/contract, forward/backward y rollback-safe.

No se atribuye este adaptador a Kubernetes. `kubectl rollout status`/`undo` con exit cero, un dry-run o un fixture no satisfacen el control.

## Reconstrucción y gates

- acquisition core 0.4.75: 25/25 archivos; 121/121 fuentes, nueve negativos; 16 perfiles y cinco negativos/un positivo; perfil Kubernetes 1 y enterprise 22 PASS;
- secure operations 1.0.0: 34/34 archivos, comparación staging↔materialización exacta;
- adapter Kubernetes: 14/14 regresiones PASS;
- grupo productivo: 87 casos, 86 PASS/1 skip de symlink de plataforma;
- `VERIFY_LIBRARY_PASS`: 73 packs/738 archivos, 440 Markdown, backend 16/123, web 5/69;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 73 packs/121 fuentes, exit 0.

## Fallos convertidos en memoria

`LIB-FAIL-1222` a `LIB-FAIL-1226` registran el tarball incompatible con el contrato ZIP, contextos de patch rechazados, frontera array del updater, parámetro materializer supuesto y comparación contaminada por bytecode. `UP-FAIL-201` conserva firma y alcance reales de Kubernetes.

## Límite honesto

No se accedió a un cluster, workload, registry, servicio, migración ni datos de proyecto. V112 prueba maquinaria reusable y fail-closed; no prueba que un proyecto esté desplegado o listo para producción.
