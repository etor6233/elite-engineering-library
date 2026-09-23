# Garantía de servicio en la composición de referencia

La garantía conserva los términos de la cotización aceptada. El pedido, pago,
entrega, caso de servicio, aprobación y stock mantienen sus escritores existentes.
Las comparaciones inclusivas de fechas de piezas/mano de obra son la adaptación
de `ServiceItemLine.CheckWarranty` identificada en
`docs/provenance/BC_WARRANTY_DERIVATION.json`. Configuración, consentimiento,
calendario, transporte y coordinación transaccional son AUTHORED; este módulo no
reproduce todo Business Central ni se atribuye ese trabajo a Microsoft.

El perfil configura explícitamente términos, zona horaria, días, exclusiones,
organizaciones y duración de la reserva pendiente. El último día incluido se calcula como fecha inicial + `*_duration_days`; el predicado incluye ambos extremos. El ejemplo es un fixture de
infraestructura. No contiene credenciales ni afirma obligaciones legales de un
país. La conciliación seleccionada es un comprobante interno de servicio, con
referencias al costo de inventario. No crea un reembolso, nómina o documento fiscal.

Desde la raíz materializada, un perfil se genera en un destino ausente:

```powershell
$source = 'deploy/warranty/profile.reference.json'
$sourceHash = (Get-FileHash -LiteralPath $source -Algorithm SHA256).Hash.ToLowerInvariant()
go run ./cmd/warranty-profile --input $source --sha256 $sourceHash --output deploy/warranty/active
$activation = Get-Content -LiteralPath deploy/warranty/active/activation.json -Raw | ConvertFrom-Json
$env:WARRANTY_ENABLED = 'true'
$env:WARRANTY_TENANT_ID = $activation.tenant_id
$env:WARRANTY_ORGANIZATION_ID = $activation.organization_id
$env:WARRANTY_PROFILE_SHA256 = $activation.profile_sha256
$env:WARRANTY_PROFILE_FILE = (Resolve-Path -LiteralPath deploy/warranty/active/profile.json).Path
```

El CLI admite overrides explícitos `--tenant-id`, `--organization-id` y
`--factory-organization-id`; vuelve a validar mediante el mismo loader del host.
Conserva el hash de origen y los cambios en `activation.json`. Ambos archivos se
reconstruyen con bytes idénticos en directorios distintos; la activación contiene
una ruta relativa. No reutiliza ni sobrescribe destinos existentes. Un fallo de
escritura conserva el destino parcial para inspección.

Con `WARRANTY_ENABLED=false` no se leen perfiles. Activado exige ruta absoluta,
archivo regular de hasta32768bytes, hash y scope exactos, migraciones0069/0070 y
los guards habilitados. La política de términos ofrecidos se vuelve inmutable;
un cambio de precio/cliente/variante requiere otra cotización. El cliente reconoce
el hash y versión antes de aceptar. Una configuración posterior no modifica el
perfil vendido ni inventa una anulación por reembolso.

Todos los endpoints requieren `organization_id` y bearer verificado. La API
reduce la autoridad al permiso y organización de esa ruta, incluso con roles
que poseen varias organizaciones. Las versiones viajan como strings decimales;
se rechazan claves duplicadas, campos desconocidos y objetos mayores de32768bytes.

| Operación | Ruta POST | Permiso |
| --- | --- | --- |
| Ofrecer términos | `/v1/franchise/warranty/quotes/{id}/terms` | `warranty:offer` |
| Consentimiento | `/v1/customer/warranty/quotes/{id}/acknowledgements` | `warranty:self` |
| Activación | `/v1/franchise/warranty/handovers/{id}/activation` | `warranty:activate` |
| Abrir caso | `/v1/customer/warranty/claims` o `/v1/franchise/warranty/claims` | `warranty:self` o `warranty:request` |
| Diagnóstico | `/v1/franchise/warranty/claims/{id}/diagnosis` | `warranty:diagnose` |
| Plan/reservas | `/v1/franchise/warranty/claims/{id}/plan` | `warranty:plan` |
| Decisión | `/v1/franchise/warranty/claims/{id}/decision` | `warranty:approve` |
| Trabajo/repuestos | `/v1/franchise/warranty/claims/{id}/work` | `warranty:work` |
| Calidad | `/v1/franchise/warranty/claims/{id}/quality` | `warranty:quality` |
| Aceptación | `/v1/customer/warranty/claims/{id}/acceptance` | `warranty:self` |
| Conciliación | `/v1/factory/warranty/claims/{id}/reconciliation` | `warranty:reconcile` |
| Cancelación previa al consumo | `/v1/franchise/warranty/claims/{id}/cancellation` | `warranty:cancel` |

Las lecturas de cotización/activación tienen rutas GET de cliente y franquicia.
El perfil actual se consulta en `/v1/franchise/warranty/profile`. El caso y su
último recibo se consultan en `/v1/{customer|franchise|factory}/warranty/claims/{id}`;
`command_id` recupera el recibo exacto de un comando. Permisos de lectura:
`warranty:self`, `warranty:read` o `warranty:factory-read`. La fábrica se toma del
perfil vendido; el cliente debe coincidir con el propietario del caso.

La apertura liga unidad, cliente y cita de servicio confirmada/completada. El
diagnóstico conserva causa y evidencia. Un plan por caso admite hasta16líneas
de repuestos enteros y trabajo de servicio; usa reservas `demand_kind=service`.
La decisión exige un sujeto diferente y payload exacto. No se autoaprueba. Un
plan fuera de cobertura puede rechazarse; no reserva piezas ni permite aprobación.
La aprobación convierte la reserva pendiente en hold durable de reparación.

El trabajo consume el stock con el writer FIFO/specific existente en la misma
transacción Serializable del recibo. La calidad la registra un sujeto diferente
del técnico. Un resultado fallido permite otra revisión con evidencia de
corrección; no permite aceptación. El cliente acepta el resultado aprobado y la
fábrica confirma el comprobante enlazado antes de cerrar. Nuevas necesidades de
repuestos después de un plan o consumos ya terminados requieren un caso de
seguimiento; este perfil no agrega líneas silenciosamente a un plan aprobado.

Ante respuesta perdida, consultar el comando guardado antes de decidir un retry.
El replay exige actor y payload originales. El rechazo o la cancelación anterior
al consumo liberan las reservas y conservan la decisión histórica. No revierten
stock ya consumido ni alteran dinero. El caso genérico no puede omitir los recibos
conectados. Los downgrades de garantía rechazan datos existentes; una base vacía
puede bajar/subir ambas migraciones. El frontend por rol se compone por su owner.
