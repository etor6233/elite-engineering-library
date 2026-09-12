# All Implementation Packs Materialization — 2026-08-26 V14

## Estado gobernante

Este snapshot sucede a V13 para la biblioteca fuente actual. Conserva 37 packs y eleva el total materializable de 335 a 337 archivos porque `OFFICIAL-UPSTREAM-ACQUISITION-CORE` añade dos perfiles ejecutables. El pack de adquisición pasa de 12 a 14 archivos y de 63 a 66 fuentes exactas.

## Ampliación de ingesta de archivos

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.12 incorpora Cisco Talos ClamAV 1.5.4, VirusTotal YARA-X 1.20.0 y Google safearchive `f7ce9d7…`. El perfil `secure-file-ingestion-leaders` selecciona seis fuentes; `secure-archive-linux` selecciona una; el perfil documental base selecciona 27. La evidencia exacta es `LEADER_LOCAL_FILE_SECURITY_2026-08-26_V1.md`.

ClamAV source y Windows ZIP tienen SHA y firma GPG Cisco verificados; `freshclam` y database tests pasaron. YARA-X version/rule scan/compile pasaron. safearchive pasó tres paquetes en Linux/WSL2 y retiene quince fallos Windows, por lo que nunca se selecciona automáticamente fuera de Linux. El main actual de Microsoft Content Processing fue re-evaluado y rechazado como actualización porque conserva fallos Processor/API/Jest.

El template documental ahora obliga a cerrar canales, original inmutable, quarantine/lifecycle, content type, malware/signature freshness, YARA rulesets, archives/bombs/traversal/symlinks, PII, provenance, idempotencia, DLQ, replay, reconciliación, restore y pruebas de admisión antes de programar.

## Conteos y límites

- packs: 37;
- archivos materializables: 337;
- Markdown fuente esperado al cierre: 221;
- evidencias: 93;
- source archives exactos: 66;
- perfiles oficiales: 5;
- perfil secure-file: 6;
- perfil archive Linux: 1;
- perfil documental base: 27;
- SDK artifacts documentales: 5;
- provider adapters: 7.

El conteo no promueve código condicionado. La biblioteca sigue `NOT_READY_UNDER_EXPANDED_USER_STANDARD`: adquirir fuentes no demuestra un pipeline real ni exactitud sobre documentos del proyecto. Foundation V12 conserva la última ejecución integral de los archivos de aplicación no modificados; V14 exige Audit actual del pack/perfiles y verifier global antes de considerarse cerrado.
