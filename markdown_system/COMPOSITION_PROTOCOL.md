# Composition Protocol

## 1. Entradas

- `PROJECT_BLUEPRINT.md`;
- `PROJECT_AUTHORITY_MAP.md`;
- catálogo de capacidades;
- restricciones del workspace existente;
- expedientes de licencias y versiones actuales.

## 2. Selección

Para cada journey:

1. dividir en capacidades verticales;
2. consultar autoridades y riesgos;
3. buscar packs con `implementation >= RECONSTRUCTIBLE`; preferir `REBUILD_VERIFIED` y permitir `CONDITIONED` cuando sus condiciones puedan probarse automáticamente;
4. descartar incompatibilidades de stack, licencia, operación o dominio;
5. comparar construir/adoptar/servicio externo;
6. seleccionar la composición mínima y registrar alternativas rechazadas.

La ausencia de pack no se oculta: se crea un trabajo de admisión o implementación, no código improvisado presentado como élite.

La ausencia de `REUSABLE_PACK` no congela el proyecto. El agente puede construir el pack faltante, usar un `REBUILD_VERIFIED / CONDITIONED` satisfaciendo sus condiciones o implementar desde contratos autoridad como trabajo nuevo `AUTHORED`, siempre con tests y procedencia honesta.

## 3. Resolución de conflictos

Prioridad:

```text
invariante y requisito del proyecto
> seguridad/licencia/regulación
> contrato de integración
> pack específico de capacidad
> pack de foundation
> preferencia de framework
```

Dos packs que escriben el mismo path necesitan merge plan explícito. Dos sources of truth para la misma entidad están prohibidos salvo sincronización y ownership formalizados.

## 4. `PROJECT_PACK_PLAN.md`

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/PACK.md",
      "packId": "PACK-ID",
      "version": "0.0.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Este bloque es el input ejecutable del compositor y usa exactamente sus nombres `camelCase`. `acknowledgeConditions` registra una selección deliberada; no demuestra que las condiciones estén satisfechas. `stackDecision`, claims usados, evidencia de licencia, desviaciones, colisiones, gates y capacidades abiertas se registran en secciones Markdown adyacentes del `PROJECT_PACK_PLAN.md`, no como claves no soportadas del manifiesto ejecutable.

## 5. Materialización

La implementación de referencia es `implementation_packs/MARKDOWN_COMPOSITOR_CORE.md`. Se materializa primero con `materialize_markdown_pack.ps1`; luego recibe un plan JSON o Markdown con JSON fenced. PowerShell es tooling de bootstrap, no decisión de stack del producto.

1. verificar hashes/pins/licencias;
2. resolver variables sin escribir secretos;
3. crear archivos en orden del pack;
4. aplicar formatting sólo después de preservar semántica;
5. generar registro `MATERIALIZATION_RECORD.md` con block IDs → paths → hashes;
6. ejecutar gates rápidos;
7. construir un vertical slice antes de expandir módulos;
8. ejecutar gates completos y capturar evidencia.

## 6. Adaptación empresarial

La configuración cambia composición y política declarativa. No autoriza a omitir módulos necesarios. Ejemplo:

```text
modules.payments=false
→ no materializar UI/API/worker de pagos
→ pedidos no pueden declarar estado `paid`
→ acceptance journeys se recalculan
→ adapters de payment quedan incompatibles
```

Todo flag debe modificar coherentemente UI, API, dominio, datos, permisos, telemetría y tests; esconder un menú no desactiva una capability.

## 7. Cierre

El agente entrega:

- archivos y decisiones;
- packs/bloques usados;
- licencias/notices;
- tests ejecutados y no ejecutados;
- claims demostrados;
- riesgos y capacidades abiertas;
- comando siguiente reproducible.
