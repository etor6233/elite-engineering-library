# Admin/ops V403 — ampliación de T2804 sobre público R3

Revisión: `V403-ADMIN-OPS-0.1.0-r3`. Pack `TS-ADMIN-OPS-EXPERIENCE-V403`, versión `0.1.0`.
Implementación autorizada por el usuario. **Aprobación visual: PENDING_USER_APPROVAL.**
No registra APPROVE del shell público, no cambia el gate de producto ni autoriza producción.

## Autoridad y composición

Este plan amplía [FRANCHISE_EXPERIENCE_PACK_PLAN_V403.md](FRANCHISE_EXPERIENCE_PACK_PLAN_V403.md)
y [PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md](PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md).
La secuencia es **UI 0.3.0 → público R3 → admin R3**. No crea otro design system ni roadmap.
El expediente de autorización es `V403-ADMIN-PASS2-PLAN-01`; sus bytes permanecen históricos.

El pack crea únicamente `admin_ops_overlay/**`. Su aplicación posterior exige una composición
pública R3 exacta y un **destino ausente** fuera de la base y del bundle. No aplicar sobre V402,
el workspace histórico, ni una franquicia real. El manifiesto compuesto declara los 8 reemplazos
existentes y los 7 archivos nuevos; cada uno tiene before/after SHA. No reinterpretar como permiso
para editar el manifiesto protegido del público R3.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [{
    "path": "implementation_packs/TYPESCRIPT_ADMIN_OPS_EXPERIENCE_V403.md",
    "packId": "TS-ADMIN-OPS-EXPERIENCE-V403",
    "version": "0.1.0",
    "acknowledgeConditions": true,
    "files": ["*"],
    "variables": {}
  }]
}
```

## Recorrido de reconstrucción

1. Materializar el Markdown con `materialize_markdown_pack.ps1` en un directorio nuevo.
2. Ejecutar `python <bundle>/admin_ops_overlay/prepare_bundle.py`, según §6 del pack. Conserva exactamente el
   manifiesto, incluso sus finales de línea; transportar en base64 no lo reescribe.
3. Ejecutar `apply_admin_overlay.py` con `--base`, `--bundle`, `--destination` y el digest
   confiable de `composed-manifest.json`:
   `f2d4b6e42e12c3d084e6bd9347be216dd34417e98dbf5e584bb9c3106fceae6a`.
4. Conservar `ADMIN_OVERLAY_MATERIALIZATION_RECEIPT.json` del destino. Su PASS acredita
   reconstrucción de fuentes; no es una segunda firma, un deploy ni un PASS de producción.
5. Reutilizar locks y toolchain del perfil. Construir con el comando del expediente y ejecutar
   el harness incluido cuando exista un delta o un entorno nuevo que lo justifique.

Los paths absolutos de owners dentro del manifiesto identifican las fuentes auditadas. Los
archivos de producto se resuelven desde el `--base` explícito mediante rutas relativas fijadas.
No se requieren secretos para reconstruir fuentes o preparar fixtures.

## Alcance y preservación

El tratamiento nuevo se aplica a `/franchise` y a `/experience`. La demo de operaciones sigue
siendo opt-in. La autorización de servidor, aislamiento y BFF conservan sus owners. Las pruebas
operativas usan el BFF real con backend sintético en memoria: no acreditan durabilidad Go/PG,
aceptación real del cliente ni cobro live.

Se preservan bytes públicos R3, layout, tokens, locks y configuración; no hay CSS global nuevo.
El chrome conserva su tratamiento anterior fuera de las rutas autorizadas. Las comparaciones
ejecutadas cubren 5 rutas públicas y 6 privadas fuera del alcance, en PC y viewport móvil.
`/admin` y `/customer` no se declaran rediseñadas.

El delta contiene presentación y composición **AUTHORED**, y adaptación declarada de harness
local. No incorpora componentes, fuentes, iconos ni assets de xAI u otras webs observadas.
Los métodos Fluent/WAI-ARIA orientan decisiones; citar métodos no atribuye origen del código.
Los notices y locks heredados siguen siendo obligatorios; no se adquieren dependencias nuevas.

## Evidencia y parada

[ADMIN_OPS_EXPERIENCE_V403_R3.md](../reconstruction_evidence/ADMIN_OPS_EXPERIENCE_V403_R3.md)
enlaza build, reconstrucción, pruebas, galería y recibo de esta revisión.

La anterior aprobación condicional admin quedó revocada para cierre visual; el recibo PHASE2
histórico no se cambia. Cada candidata nueva requiere aprobación humana explícita. Viewport
simulado no es teléfono físico; automatización no equivale a lector de pantalla ni aceptación.
Cloud, producción, cámara QR/OCR y checkout live quedan fuera del claim de este overlay.

**Parar para control visual después de este expediente. Sin commits ni push en esta fase.**
Árbol de ampliación post-V402 ≠ bytes del ZIP de cierre; el ZIP sigue siendo la instantánea
canónica del READY original. El checkpoint 337 permanece COMPLETE e intacto.
