# Contratos de funciones empresariales V403

Este directorio implementa un scaffold de contratos y validación para 11 funciones / 22 tareas. Es AUTHORED glue de biblioteca. No implementa once operaciones empresariales ni atribuye su código a Microsoft, DORA, HubSpot o xAI.

Las nueve referencias a packs existentes enlazan archivos y evidencia histórica de la biblioteca. Sus hashes se verifican contra el baseline 337. La consola no llama sus APIs; los límites de cada pack siguen vigentes.

## Lectura y frontend

business-functions.v403.json es el contrato consumido por la vista. Cada función contiene propósito, sistemas, tareas, inputs, outputs, permisos requeridos, aceptación, sources, owners, resume y binding de implementación. La ruta es #/functions/{id}. La cobertura de todas las secciones y tareas exige además prueba del frontend real; el validator de C sólo verifica sus vínculos declarados.

Grok Bot 101 es LEARNING_TRACK. Los demás nombres son BUSINESS_FUNCTION, nunca roles de autenticación. effective_grants permanece vacío y ninguna task concede permisos. Nombres propuestos se distinguen de marketing:read / marketing:request observados en el owner existente. El producto consumidor resuelve sus propios subjects, objetos, organizaciones y enforcement.

Los efectos están deshabilitados en este catálogo. Esa restricción no prohíbe operaciones de un producto que cuenten con autorización explícita y pruebas propias.

## Estados y DONE

- specification_status=SPEC_READY: el contrato local está definido.
- method_admission=OFFICIAL_METHOD_REFERENCED: método documental actual y claim limitado; no es admisión de código o arquitectura.
- galaxy_admission=WAITING_GALAXY_ADMISSION: sólo falta material futuro del evento. No impide usar las referencias documentales actuales.
- verification_status.target=NOT_TESTED: no hay prueba del negocio ni del target futuro.
- El JSON portable conserva scaffold=NOT_TESTED como declaración de entrada sin autoatestación. Los resultados ejecutados viven en receipts externos ligados al SHA exacto del contrato; la vista debe mostrar ambos alcances.
- PROVEN exige evidencia existente con result=PASS, claim_scope=FUNCTION_CONTRACT_SCAFFOLD, hash exacto y función incluida. Un hash acredita correspondencia, no la veracidad semántica de un informe escrito por un tercero.
- DONE del scaffold exige checks locales PASS, cobertura frontend independiente y cero defectos conocidos abiertos en ese alcance y revisión. CONDITIONED / NOT_TESTED no cuentan como PASS de una garantía requerida. Los owners canónicos conservan el triage; el validator no inventa un conteo de bugs.

## Owners y reanudación

El archivo entregado tiene owner_binding.scope=MAINTENANCE_FIXTURE. Referencia los seis owners existentes de este mantenimiento V403. Ningún proyecto consumidor debe heredar esas rutas, su identidad, progreso o aprobaciones.

Después de materializar el scaffold en el consumidor, preparar un owner-map JSON con claves progress, execution_state, execution_events, failures, dependency y freshness. Sus valores son rutas relativas de archivos YA existentes del proyecto. Resolver progress al tasks owner real (por ejemplo el PROJECT_TASKS_REF del consumidor), sin crear otro backlog.

Ejemplo de integración, con rutas elegidas por el consumidor:

    python roles/bind_consumer_owners.py --source roles/business-functions.v403.json --project-root . --owner-map owner-map.json --output roles/business-functions.consumer.json

El binder crea sólo un contrato nuevo, rechaza destinos existentes, rutas fuera de raíz, IDs/owners de mantenimiento y owners ausentes. No modifica state/events ni copia progreso. Resetea claims/evidencias de verificación y no concede permisos. Después ejecutar el validador de execution state del kit 1.3.1 del consumidor; OWNERS_BOUND_ONLY no autoriza resume.

El gate de contrato acepta --contract roles/business-functions.consumer.json. Sus referencias al kit y a los packs se resuelven contra --library-root. Una modificación material del contrato o de los refs invalida los receipts afectados.

## Comprobación local

Desde la raíz de un mantenimiento o consumidor con owners resueltos:

    python roles/validate_business_functions.py --project-root . --library-root "C:/ruta/biblioteca" --report roles/evidence/validation-new.json
    python roles/tests/test_business_functions.py --library-root "C:/ruta/biblioteca"

El report path debe ser nuevo. Los tests crean fixtures sintéticas aisladas bajo roles/evidence y no modifican los owners. Cubren privilegios implícitos, scope por objeto/organización, links, procedencia, Galaxy futuro, PROVEN no ligado, cambios de pack y herencia de owners. No ejercitan proveedores, BFF, usuarios humanos ni operaciones de negocio.

## Fuentes y procedencia

method-sources.v403.json contiene ocho fuentes oficiales actuales con consulta, claim y límite. Son enlaces y paráfrasis breves propias. No se adquirió código upstream ni se seleccionó Business Central por citar Dynamics 365.

evidence/source-observations.json conserva URL, fecha de consulta, HTTP200, tamaño y huella SHA-256 de la respuesta observada. No retiene cuerpos HTML: es evidencia histórica de observación, no snapshot reproducible ni firma editorial. Revalidar la fuente oficial si cambia una decisión material; no tratar el hash histórico como prueba de contenido actual.

library-bindings.v403.json sí permite comparar bytes locales de packs/evidencia existentes. Estas referencias no promueven el estado del producto consumidor.

## Reconstrucción

La guía canónica es markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md. La única fuente de código es implementation_packs/BUSINESS_FUNCTION_OPERATING_V403.md y usa los bloques y hashes estándar de la biblioteca.

Materializar con el script ya existente, sin crear otro compositor:

    pwsh -NoProfile -File materialize_markdown_pack.ps1 -PackFile implementation_packs/BUSINESS_FUNCTION_OPERATING_V403.md -Destination C:/ruta/ausente

El pack contiene ocho archivos. Conserva sus receipts y vuelve a ejecutar el gate con los owners del mantenimiento o consumidor. No copia state/events ni aprobaciones. Python usa sólo la biblioteca estándar; JSON, validator, binder y tests son AUTHORED glue. No se instala dependencia ni se atribuye código local a un proveedor.
