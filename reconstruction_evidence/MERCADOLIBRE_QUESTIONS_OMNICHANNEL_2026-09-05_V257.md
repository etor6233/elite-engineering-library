# Evidencia V257 — Mercado Libre Questions en el ingreso omnicanal

Fecha: 2026-09-05. Resultado de biblioteca: `PASS`. Estado productivo del provider/target: `CONDITIONED`.

## Autoridad y procedencia

- Autoridad oficial vigente: Mercado Libre, *Manage questions and answers*, actualización publicada 2026-01-15.
- Claims estrechos usados: topic `questions`; `GET /questions/{id}?api_version=4`; búsqueda por `seller_id` con estado `UNANSWERED`; texto vacío permitido cuando el estado es `BANNED`; identidad de seller/buyer/item y datos de contacto que el contrato muestra para Vehicles.
- Mercado Libre recomienda asistencia semiautomática para responder. V257 no implementa `POST /answers` ni atribuye una respuesta automática a la empresa.
- `GO-MERCADOLIBRE-MARKETPLACE-ADAPTER 0.2.0` es `AUTHORED` contra el contrato HTTP oficial. No es un SDK de Mercado Libre: el SDK Go oficial archivado continúa rechazado como no funcional/no mantenido.
- Los dos archivos nuevos de `GO-OMNICHANNEL-LEAD-INGRESS 0.2.0` son uno `ADAPTED` al contrato oficial y uno `AUTHORED` de pruebas. No contienen bytes copiados del provider.

## Vertical materializado

```text
notification topic=questions
→ resource allowlisted /questions/{id}
→ GET v4 autenticado
→ seller esperado obligatorio
→ RawEvent hash-bound + payload provider preservado
→ candidato pending_policy o rechazo durable
→ PostgreSQL lead_ingress_raw + lead_candidate + outbox atómicos
→ promoción CRM/contacto separada
→ respuesta sólo mediante política, aprobación y outbound fence
```

Se añadió reconciliación `GET /questions/search?seller_id=...&status=UNANSWERED&api_version=4`. Estados `BANNED`, `DELETED`, `DISABLED`, `UNDER_REVIEW`, fecha inválida, buyer ambiguo o estado desconocido no producen una identidad inventada: conservan evidencia rechazada. Los datos personales posibles de Vehicles nunca cambian `contact_eligibility=pending_policy`.

## Evidencia ejecutada

- Toolchain: Go 1.26.7 windows/amd64, `go.exe` SHA-256 `5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc`.
- Adapter Mercado Libre: materialización 10/10, 10 tests y `go vet ./...` PASS.
- Ingreso omnicanal: materialización 12/12; tests focales y vet de `internal/leadstream` PASS.
- `SECURE-OPS-DELIVERY-CORE 1.1.1`: 37/37; 101 casos, 100 PASS/1 skip; catálogo de admisión fija el adapter 0.2.0.
- `VERIFY_LIBRARY_PASS`: 159 packs, 1.392 archivos materializables; perfil franquicia 66 packs/701 archivos sin colisiones.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 159 packs, 121 fuentes, 16 adapters y todos los gates offline aplicables PASS. OpenGrep continúa correctamente bloqueado porque Cosign 3.1.3 no pasa su SCA vigente.

## Condiciones que este PASS no inventa

Continúan requiriendo evidencia del proyecto: aplicación/cuenta/seller/token/scopes, configuración y autenticidad del callback, pruebas sandbox/live permitidas, cuotas/costo, retención/privacidad, mapping con CRM, política de contacto, operación humana, respuesta real, reconciliación provider, carga, seguridad, recovery y aceptación empresarial. Ninguna de esas condiciones se declara cerrada por fixtures o tests offline.
