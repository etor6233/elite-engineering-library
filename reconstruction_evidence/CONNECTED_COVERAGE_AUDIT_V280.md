# V280 — cobertura seleccionada y admisión reabierta de observabilidad

Fecha: 2026-09-07. Mantenimiento/discovery, no expansión de producto.
Staging: elite-v280-741be8bd2ddd45d6b18747443f33991e.
Continuidad inicial: checkpoint V279 revisión 3, diez hashes y tres eventos PASS.

## Resultado del contraste de selección

Se leyó el JSON canónico FRANCHISE_COMPLETE_PACK_PLAN.md y se compararon pack_id
de los siguientes 21 packs citados en FRANCHISE_GAP_MAP.md. Ninguno está entre
los 67 seleccionados. No se deduce que falten 21 capacidades: puede haber una
implementación alternativa en commerce, supply, service, web u otro owner.
La equivalencia necesita demostrar los recorridos y no se asume por nombre.

Procedencia comprobada en metadata/bloques: los 21 packs son íntegramente
AUTHORED; 15 tienen upstream_sources igual a ["https://go.dev"]. Esa referencia
al lenguaje no identifica un ledger, POS o módulo de negocio publicado por
una compañía, ni demuestra equivalencia de calidad. Los otros seis tampoco
se convierten en código upstream por citar documentación adicional. Se debe
resolver la admisión de dominio y la integración antes de tratarlos como
código empresarial reutilizable; no se retira código ni se fabrica otra fuente.

| Pack no seleccionado | Revisión que debe hacerse antes de incorporarlo |
|---|---|
| GO-DASHBOARDS-CORE | comparar queries/portales/reporting existentes |
| GO-FX-CORE | fuente de tasas y reglas monetarias del proyecto |
| GO-GIFT-CARDS-CORE | saldo durable, autorización y reversión |
| GO-HELP-CENTER-CORE | ayuda frontend y versión del journey |
| GO-I18N-CORE | mensajes existentes, locales y experiencia por rol |
| GO-LOYALTY-CORE | política, earn/burn ligados a ventas y persistencia |
| GO-MARKETING-CORE | consentimiento, segmentación y efectos provider |
| GO-OBSERVABILITY-CORE | NO ADOPTAR: defectos reproducidos debajo |
| GO-ONBOARDING-CORE | flujo existente de alta y permisos |
| GO-PAYROLL-CORE | reglas de jurisdicción, fuente y aprobación |
| GO-POS-CORE | comercio/pagos/fiscalidad; no duplicar ventas |
| GO-PROMOTIONS-CORE | pricing/promociones ya existentes |
| GO-REFERRALS-CORE | vínculo con leads, atribución y antifraude |
| GO-REMINDERS-CORE | scheduler/jobs/notificaciones ya existentes |
| GO-REVIEWS-CORE | identidad, moderación y publicación |
| GO-SEO-CORE | SEO ya implementado en web pública |
| GO-SLO-CORE | señales/alertas reales del mismo journey |
| GO-SOCIAL-POSTING-CORE | contratos de publicación y aprobación de efectos |
| GO-SURVEYS-CORE | consentimiento, respuesta y reporting |
| GO-WAITLIST-CORE | agenda/capacidad/orden y notificación |
| GO-WARRANTY-CLAIMS-CORE | service/devoluciones existentes y sus invariantes |

Este listado es una reconciliación de selección, no una auditoría funcional
completa de esos 21 cores. GO-APP-WIRING sí está seleccionado, pero su claim
es el runtime conversacional: no ensambla automáticamente todos estos módulos.
La inspección de GO-LOYALTY-CORE muestra mapas/slice protegidos por mutex;
su conexión a compra/identidad se declara responsabilidad de composición.
No se debe confundir ese ledger en memoria con persistencia PostgreSQL.

## Cuatro contraejemplos ejecutados — código 0.1.0 exacto

Materialización: observability-audit, dos archivos canónicos.
observability.go SHA-256:
71bf24281aec336193df45ca161f9f05018995c2ad02255fce61ea4630bfecb4.
observability_test.go SHA-256:
272d4fae3a6c5e2031ceb402af6e0d08b89bd677f325486515b7ed378ff2fafb.

Se añadió sólo un test adversarial temporal AUTHORED, no código de producto.
Go 1.26.7, Windows, GO111MODULE=off, `go test -count=1 -v` desde
observability-audit/internal/observability. Log: observability-audit.log.

| Caso | Esperado para el claim anunciado | Observado |
|---|---|---|
| Counter.Inc(5), después Inc(-1) | contador no decreciente | baja de 5 a 4 |
| Redact con campo email | no exponer PII si se promete protección general | correo sintético sin cambio |
| Redact con metadata.password | secreto anidado no emitido | secreto sintético sin cambio |
| Logger.Info con email en mensaje | no emitir PII si se promete protección general | Msg llega al sink intacto |

Resultado: **cinco tests históricos PASS y cuatro pruebas adversariales FAIL**,
exit 1. No es una filtración productiva observada; todo el input fue sintético.
Monotonía contradice una invariante explícita; los tres restantes delimitan el
claim amplio de redacción, no afirman que exista detector universal de PII.

Reproducción mínima de los cuatro casos (mismos inputs y aserciones):

```go
package observability
import "testing"
func TestV280CounterMonotonic(t *testing.T) {
 var c Counter
 before := c.Inc(5)
 if after := c.Inc(-1); after < before { t.Fatal("counter decreased") }
}
func TestV280EmailRedaction(t *testing.T) {
 input := "synthetic@example.invalid"
 if got := Redact(map[string]any{"email": input}); got["email"] == input { t.Fatal("email survived") }
}
func TestV280NestedSecretRedaction(t *testing.T) {
 input := map[string]any{"password": "synthetic-not-a-real-secret"}
 got := Redact(map[string]any{"metadata": input})
 if nested, ok := got["metadata"].(map[string]any); ok && nested["password"] == input["password"] { t.Fatal("nested secret survived") }
}
func TestV280MessageRedaction(t *testing.T) {
 input := "contact synthetic@example.invalid"
 var observed Entry
 l := NewLogger(LevelInfo, func(e Entry) { observed = e })
 l.Info(input, nil)
 if observed.Msg == input { t.Fatal("message survived") }
}
```

## Contención y procedencia

El código original no se corrigió ni se atribuyó a un tercero. Se preservaron
sus hashes e histórico V219; admisión pasa de CONDITIONED a CANDIDATE, lo que
el compositor canónico rechaza aunque acknowledgeConditions=true. No lo usa
ninguno de los 51 perfiles de composición actuales según búsqueda exacta del ID.
La corrección o reemplazo permanece pendiente bajo T2809, con admisión y pruebas.
La condición no se resuelve añadiendo un Collector detrás de un log que ya filtró.

Fuentes oficiales consultadas como método, no código importado ni aval del pack:

- Google/Android, [Log Info Disclosure](https://developer.android.com/privacy-and-security/risks/log-info-disclosure): riesgo de datos sensibles en logs; su contexto es Android, no se extrapolan configuraciones de logcat al backend Go.
- OpenTelemetry, [Handling sensitive data](https://opentelemetry.io/docs/security/handling-sensitive-data/): limitar datos recogidos y aplicar controles explícitos; no garantiza reconocer cualquier PII ni respalda este logger local.

## Cierre y límites de esta auditoría

REUSABLE_CODE_READINESS_ROADMAP.md §10 incorpora T2801–T2810 como desglose del
hito existente. Nueve bloques asignan exactamente las 48 superficies normativas;
T2810 es validación/distribución conjunta. T2800 sólo cierra este contraste y
reproducción focal. No hay porcentaje ni estimación de esfuerzo demostrados.

Controles ejecutados: CAPABILITY_ALLOCATION_48_OF_48_EXACT_PASS (conjunto exacto,
sin duplicados) y EXPECTED_CANDIDATE_COMPOSITION_REJECTION_PASS (perfil temporal
con acknowledgeConditions=true rechazado por el compositor canónico). No se
desactivaron aserciones ni se etiquetaron como PASS los cuatro contraejemplos.

Pendientes: readiness/assurance de mantenimiento, equivalencias de los cores,
recorridos integrales, código admisible para observabilidad, operación de host,
providers y configuración del target. La estructura conserva 160 packs/1433
archivos, perfil 67/742, 51 perfiles y pasa a 713 Markdown. No se crea ZIP,
cuenta, suscripción, envío externo, dependencia ni pack nuevo.
