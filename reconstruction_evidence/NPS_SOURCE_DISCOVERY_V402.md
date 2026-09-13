# V402318 — T2801, frontera NPS pendiente de admisión

T2809 está cerrado en317; la revisión318 no cambia ningún archivo de producto.
T2801 permanece IN_PROGRESS. FAIL945 identifica NPS como negocio AUTHORED aún
no sustituido, ya señalado por LIBRARY_PROVENANCE_GAPS_V402. No se relabela.

Inspección oficial acotada: PostHog core MIT, fuera de ee. Commit de release
6fafbb9081bd15e79448af5650e02a4f9ea435cc, tag desktop-v0.61.382 publicado2026-09-12.
Los dos blobs de survey coinciden con el árbol actual inspeccionado por separado.
El tag no certifica nuestro adapter ni sustituye pruebas. Tres archivos obtenidos
en cuarentena se cotejan contra Gitblob/size y fijan SHA256 antes de incorporar.

La función real calculateNPSFromRawData está en utils.ts354–372; su selección
legacy340–350/379–382 y nueve casos originales utils.test.ts843–990 están
identificados. Aún no se ejecutó ni adaptó ese código. El original devuelve score
con un decimal; la API actual devuelve float64 sin redondeo. La adaptación debe
declarar ese delta y mantener el umbral de privacidad/configuración como glue.
No basta mover la fórmula anterior a otro archivo ni invocar el nombre PostHog.

| Archivo oficial fijado | SHA256 |
|---|---|
| [frontend/src/scenes/surveys/utils.ts](https://raw.githubusercontent.com/PostHog/posthog/6fafbb9081bd15e79448af5650e02a4f9ea435cc/frontend/src/scenes/surveys/utils.ts) | a43ed06c7eaffbe5bde3fe515177e512ebc4e6bce24dbef253ed100cbb1a08c1 |
| [frontend/src/scenes/surveys/utils.test.ts](https://raw.githubusercontent.com/PostHog/posthog/6fafbb9081bd15e79448af5650e02a4f9ea435cc/frontend/src/scenes/surveys/utils.test.ts) | 7b3b24fe2d6d4525113657629991e0d0056b2c8fcf91e6bd2a906b1c280bd565 |
| [LICENSE](https://raw.githubusercontent.com/PostHog/posthog/6fafbb9081bd15e79448af5650e02a4f9ea435cc/LICENSE) | 6d82d67dba42eb94ba10f1e986d2eec338c22fb7c5216c2c0ebdecd83d53a029 |

Siguiente: materializar gap gate; ejecutar los nueve casos fuente aislados con
Node24.20 ya admitido, derivación Go explícita y oráculo de agregados, G0–G8 y
notices; integrar desde el owner actual, probar resumen real y reconstrucción.
Después terminar la reconciliación T2801/G/H/384/385/532/804. TEST07/T2810 y
ARCA infra penúltima/dossier Daybreak último conservan su orden. Nada READY global.
