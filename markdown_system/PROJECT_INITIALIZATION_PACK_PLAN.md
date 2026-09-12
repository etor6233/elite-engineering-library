# Project Initialization Pack Plan

Este perfil materializa únicamente el tooling previo que adquiere fuentes oficiales sin Git o artefactos SDK PyPI hash-exactos, inicializa el harness oficial Spec Kit sin crear `.git` y ejecuta los probes documentales permitidos. No materializa todavía backend, frontend ni reglas del negocio. Debe ejecutarse después de `PROJECT_START_READINESS_GATE.md` y antes del perfil de producto elegido.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/OFFICIAL_UPSTREAM_ACQUISITION_CORE.md",
      "packId": "OFFICIAL-UPSTREAM-ACQUISITION-CORE",
      "version": "0.4.88",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/OFFICIAL_DOCUMENT_SDK_ARTIFACT_CORE.md",
      "packId": "OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE",
      "version": "0.3.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Las condiciones reconocidas son: PowerShell 7, Python 3.12 compatible, red o cache exacta para adquirir archives/artefactos/modelos, capacidad de disco, términos/licencias por upstream, cuentas/cuotas/costos del proveedor y fallos upstream consignados en evidencia. `acknowledgeConditions` reconoce el estado condicionado del pack pero no sustituye los approvals hash-linked de cada adquisición, ni habilita almacenamiento de campos REVESTEX o proveedores externos. El agente debe generar `PROJECT_EXTERNAL_SOURCE_LOCK.md`, completar los registros de selección/aprobación, probar accesos reales y cerrar corpus/ground truth antes de promover automatización documental.

V396: cerrado descubrimiento y entrega local de licencia original semver-utils1.1.4: MIT OR Apache-2.0 explícito, opción MIT completa y1839bytes originales retenidos. Core0.4.88/51files/196sources/25profiles añade sólo cuarentena npm exact4193bytes/SHA512;140checks y gates anteriores PASS. Planner0.5.0/6files/67tests;3recetas reales/12negativos PASS, BlueOak/QRCode y442payloadfiles intactos. Clave registry vencida documentada; no firma vigente, equivalencia build, ejecución de recetas ni admisión global pnpm.45/48sin promoción; V386/Daybreak sigue diferido. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.
