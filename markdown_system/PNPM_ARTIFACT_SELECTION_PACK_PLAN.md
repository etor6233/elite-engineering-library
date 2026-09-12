# pnpm Artifact Selection Pack Plan

Tooling separado para producir/verificar una proyección candidata de pnpm ya
adquirido. No ejecuta, instala ni admite runtime. No se incluye automáticamente
en perfiles de producto, ni cambia el transporte opaco de adquisición.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/PNPM_ARTIFACT_SELECTION_GATE.md",
      "packId": "PNPM-ARTIFACT-SELECTION-GATE",
      "version": "0.8.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Windows/CPython del proyecto, receipt de adquisición previamente validado,
hash esperado de ese receipt y archivo exacto obligatorios.9files;90tests históricos de selección/routing sin delta ejecutable,
proyección a target ausente y verify. Receipts con runtime_admitted=false;
consumers, licencias, monitoreo y evidencia operativa son gates posteriores.
Evidencia: reconstruction_evidence/PNPM_SELECTION_GATE_V334.md.

V335 añade prepare/verify de recetas para3consumers fijados. Entorno reemplazado, store/cache explícitos, workspace empresarial preservado y destinos fresh. No ejecución automática ni admisión runtime. Ver PNPM_CONSUMER_ROUTING_V335.md.

### V336 — renovación verificable de cache pnpm

Tres checks online desde cache vacía,17metadata oficiales completas con bytes/receipts,3revalidaciones offline sin verdicts históricos y3instalaciones finales offline PASS.275entradas entre3locks;278metadata inventariadas. FAIL561 corregido; no cambio de ejecutables/versiones/política. Procedimiento y límites: reconstruction_evidence/PNPM_CACHE_FRESHNESS_V336.md. Validez acotada a esta ejecución; renovación al reanudar, sin TTL o seguridad continua inferidos. Licencias/runtime/redistribución y10macrofrentes permanecen abiertos.

V338 notice coverage:473metadata hashes verified;468identities observed from456bundle identities+21external manifests+root (with overlap).22notice files and21embedded attribution paths inventoried;32exact draft texts preserved. BlueOak five external packages and QRCode vendored MIT require explicit notice treatment; nested works, semver-utils exact text and whole-bundle source/use obligations remain open. All442payload files unchanged; no runtime/redistribution admission. Evidence: reconstruction_evidence/PNPM_NOTICE_COVERAGE_V338.md.

V339:8fixed source manifests and7licence texts verified; four BlueOak source texts,3Noble locked build-input identities outside the original473graph,70node-gyp vendored files identical to fixed tree. Draft39preserves32previous texts. chownr declaration-only, semver-utils exact text and source/bundle/redistribution conditions remain open. All442selected files unchanged. See reconstruction_evidence/PNPM_NESTED_LICENSE_SOURCES_V339.md.

V340:103external Undici6files match fixed source; Undici7licences/headers bound separately. Yarn4.1.7 registry gitHead returns404 but official tag/source version exists at a distinct commit; preserve both, no published-byte equivalence. Root BSD and3MIT source comments retained; draft47preserves39previous texts,442selected files unchanged. FAIL575 scoped provenance remains unresolved. See reconstruction_evidence/PNPM_YARN_UNDICI_NOTICES_V340.md.

V341: five selected Undici7 AST units in two bundled modules match fixed-source syntax under declared identifier renames; nine initial controls plus eight actual-source mutations and two exact byte-range checks pass. No whole-module/runtime/redistribution claim; draft47/payload442 unchanged. Contract remains41passed/7pending of48 (approximately15% of rows, never global effort); ten macrofronts remain open. Evidence: reconstruction_evidence/PNPM_SCOPED_SOURCE_CORRESPONDENCE_V341.md.

V364: reconstruction_evidence/YARN_EMBEDDED_ASSET_PROVENANCE_V364.md fija un asset Yarn ESM completo (72463bytes decodificados) y5funciones AST contra source fijado;14controles negativos y extracción independiente. Fuente de paquetes no ejecutada.163Markdown de packs,442payload y47notices intactos. FAIL575/publicación completa,FAIL532/nativos y obligaciones de licencia siguen;43/48sin promoción.

V394: PNPM-ARTIFACT-SELECTION-GATE0.3.0/6files corrige FAIL782 (pins empresariales y Playwright obsoletos), entrega BLUEOAK-NOTICE.md ligado a5declaraciones exactas y verifica3recetas actuales;51tests y10mutaciones reales rechazadas. Payload442/22notices intactos; sin ejecución ni admisión runtime/redistribución. Evidencia: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md. Los pendientes BlueOak se reducen a entrega downstream/revisión restante; ausencia de LICENSE original chownr no se falsea. Otros permisos/source/publicación de pnpm continúan abiertos.

V395: PNPM-ARTIFACT-SELECTION-GATE0.4.0/6files entrega QRCODE-NOTICE.md con encabezado de autor/modificación exacto y texto MIT completo, además del aviso BlueOak intacto.59tests PASS,3recetas reales y12negativos rechazados;10source blobs/10regiones de bundle fijados por separado, sin claim de equivalencia completa. FAIL783 entrega local resuelto, FAIL784 metadata ADAPTED reconstruida;payload442/22notices intactos. No ejecución ni admisión runtime/redistribución;45/48. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397: PNPM-ARTIFACT-SELECTION-GATE0.6.0/7files entrega el conjunto completo de47textos retenidos/151253bytes en3recetas:141copias verificadas independientemente,12negativos reales y82tests PASS.22notices originales comparados contra payload y25evidencias conservan alcance; BlueOak/QRCode/semver y442payloadfiles intactos. Catálogo ADAPTED con términos por texto, sin algoritmos ni promoción de licencias/source/relinking/runtime.45/48; Daybreak diferido. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398: PNPM-ARTIFACT-SELECTION-GATE0.7.0/8files entrega fuente original next-path1.0.0, manifiesto y MPL completa en3recetas:9copias exactas/12negativos reales/90tests PASS. Commit oficial y3Git blobs verificados;4sentencias comparadas bajo adaptadores explícitos,10mutaciones rechazadas. No equivalencia runtime ni reproducibilidad pnpm. Colección47, suplementos anteriores y442payloadfiles intactos.45/48; Daybreak diferido. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

V401: contain_zip.py añade eliminación física de16módulos adm-zip y rechazo deZIP en candidato separado. Version0.8.0/9files. Original11.25.0 y candidato11.26.0 afectados;12.4.1 no admitido. Ver reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.
