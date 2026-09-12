# V402 — corrección de dependencias del backend independiente

FAIL818: la composición ENTERPRISE_BACKEND omitía los paquetes internos channels/outbounddelivery usados por el checkout nuevo. La reconstrucción exacta anterior no había afirmado compilación de todos los perfiles. Se preserva el fallo original.

El perfil ahora selecciona GO_CHANNELS_CORE, el transporte/fence PostgreSQL de GO_PG_OUTBOUND_DELIVERY_FENCE y la representación/hash de GO_PG_CONTACT_CHANNEL_IDENTITY. Excluye archivos de app/conversación no requeridos por este backend. No agrega versión externa ni inventa un módulo mediante go get. La franquicia completa ya seleccionaba estos owners y sigue77/943.

La composición corregida39packs/531files pasó go build ./... con la red de módulos desactivada. Después se aplicó únicamente el hook opcional del host WhatsApp: su prueba en este backend sin el pack opcional pasó, bloqueando activación true y dejando false sin leer secretos. Ese hook sigue staging con el host; este cambio canónico sólo corrige la selección del plan.

El receipt conserva el inventario exacto y hashes de logs: `C:/Users/NL/AppData/Local/Temp/elite-v402-library-infra/backend-dependency-correction.json` SHA256 `b54a97be424bf323924cbddc81df25cd95a3b6ca2b5a95b7985cde2873137ef1`. No cierre global TEST02, SCA, operaciones o release.
