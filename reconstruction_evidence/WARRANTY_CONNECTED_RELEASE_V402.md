# V402 — garantía vendida y reclamo conectado, publicación canónica

Claim PROVEN_LOCAL de T2802: términos/consentimiento de cotización→pedido→pago SDK
fixture→handover aceptado/recibo→activación de garantía; después cita/diagnóstico→
plan/repuestos→aprobación humana→trabajo/calidad→aceptación cliente→conciliación
interna de fábrica. T2802 conserva la conciliación J2/J3 y de subcapacidades;
T2804 conserva la UI de garantía/capacitación. No se declara READY global.

Dos packs nuevos: GO-BC-WARRANTY-COVERAGE-ADAPTER0.1.0,8bloques, y
GO-CONNECTED-WARRANTY-CLAIM0.1.0,26bloques. Seis owners existentes actualizados:
aplicación1.13.0, fulfillment0.5.0, aprobación0.4.0, initial-handover0.4.0,
checkout0.3.0 e inventario0.17.0. Son33archivos nuevos y8modificados; una licencia
MIT compartida se selecciona una sola vez. Biblioteca182packs/1911bloques;
franquicia86packs/1145archivos.67migraciones seleccionadas,0069/0070 añadidas.

## Procedencia y política

BCApps2eae56d704a1fd035d104f333602aea7091b7749 aporta sólo la comparación inclusiva
de fechas de partes/mano de obra y su error de período de partes invertido,
ADAPTED con ordinales/errores Go declarados. El archivo oficial completo y licencia
MIT viajan VERBATIM con SHA; no se ejecuta AL ni se atribuye el workflow a Microsoft.
1280vectores se proyectan desde tres expresiones AL fijadas, no una suite BC completa.
El glue de términos, PG, transporte/host/materializador y fixtures es AUTHORED.

El perfil fija explícitamente días/exclusiones y zona horaria. Final=start+días,
ambos extremos inclusivos; la fecha de inicio es el día local de entrega aceptada.
No se deduce duración legal. El perfil vendido sigue gobernando aunque cambie la
configuración actual. Un plan acotado por caso; repuestos adicionales requieren
un caso de seguimiento. La conciliación es acuse interno, payment_created=false.

## Evidencia funcional reutilizada por identidad

- Oferta/ack inmutables ligados a cotización y aceptación; activación ligada a
  cliente/unidad/handover/política vendida. Rollback, concurrencia y replay probados.
- Cuatro casos J4 HTTP: cierre completo y tres cancelaciones por rechazo humano,
  reparación aprobada antes de consumir stock y exclusión de falla. Reviewer
  separado; API genérica no saltea la historia conectada.
- Stock real: dos capas5@100 y5@120; se consumen7 con costo exacto740 y quedan3
  disponibles/0reservadas. Writer FIFO existente; outbox fallido revierte consumo.
- Calidad fallida exige evidencia correctiva; cliente/fábrica exactos. Dos respuestas
  HTTP perdidas después de commit recuperadas por GET actor/hash. Sin reenvío oculto.
- Lease real mínimo60s vence al COMMIT: decisión/paso/evento/consumo revierte;
  rechazo posterior libera reserva. Este probe usa handover/cita históricos SQL
  declarados, no suma otro journey J1.
- Host opt-in valida perfil/hash/scope/tamaño y guards activos. JSON/int64/scope/
  autenticación negativos; el verifier de HTTP es fixture explícito, no JWT proof.
- Downgrade poblado rechaza sin pérdida; base vacía down/up PASS sin CASCADE.
- CLI emite dos perfiles/activaciones byte-idénticos y no sobrescribe destino.
  Fuzz parser/calendario103752ejecuciones y predicado527548, presupuesto finito.
  No SAST/DAST/carga ni prueba económica de proveedor real.

Receipts previos completos: WARRANTY_CONNECTED_TERMS_V402.md,
WARRANTY_CONNECTED_CLAIM_V402.md y WARRANTY_INTERFACE_AND_PORTABILITY_V402.md.
Los rojos originales se conservan; sus correcciones están en esta fuente exacta.
No se repitieron estas suites para empaquetar.

## G0–G8 para el claim local

| Gate | Prueba y límite |
|---|---|
| G0 identidad | Receipts oficiales fijados por commit/Gitblob/SHA; archivo completo seleccionado y derivación trazable |
| G1 licencia | MIT exacta compartida; notices del producto separan ADAPTED, VERBATIM y AUTHORED |
| G2 autoridades | Predicate BC estrecho admitido; servicio/aprobación/FIFO/precision existentes gobiernan sus efectos |
| G3 arquitectura | Owners únicos, una transacción/outbox, binding físico de términos y caso; sin otro stock/ledger |
| G4 correctitud | J1/J4, costo740, policy vendida, negativos/rollback/expiry y recuperación de arriba |
| G5 seguridad del delta | Scope/actor/payload/permisos/cuerpos/perfil y aprobación separada; composición SCA/identidad sigue T2803 |
| G6 resiliencia | Reserva/locks, lease al COMMIT, respuesta perdida y rollback; carga general sigue T2809 |
| G7 evolución | Configuración versionada, materializador determinista, historia inmutable, down/up seguro y runbook |
| G8 reconstrucción | Cuatro perfiles exactos86/1145,40/567,57/793,87/1153; backend/serverless compilan su cierre de imports, cero tests runtime repetidos; rebuild final del perfil con metadata admitida |

REBUILD_VERIFIED/CONDITIONED sólo para ese perfil explícito y la composición
compatible. Las condiciones globales siguen visibles; no se disfraza código
faltante como falta de credenciales. La UI pertenece al siguiente owner T2804.

Manifest de archivos/fuentes/delta/receipts: WARRANTY_CONNECTED_RELEASE_V402.json.
Referencia externa: <LOCALAPPDATA>/Temp/elite-v402-library-infra/warranty-admitted-reference.
Destino durable y release firmado siguen T2810; producción no autorizada.
