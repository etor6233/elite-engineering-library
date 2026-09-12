# V402 — compatibilidad de proveedor Commerce / refunds

PASS del delta local de biblioteca. GO-OFFICIAL-RETURN-REFUND-WORKER pasa de
0.1.5 a 0.1.6: 25 archivos, dos fuentes existentes modificadas y dos archivos
nuevos (test y contrato). Conserva SUPPORTED_REFERENCE / REBUILD_VERIFIED /
CONDITIONED. No certifica producción ni cierra TEST02 completo.

## Corrección de integración

Commerce usa `mercadopago`; el owner histórico de refunds, su constraint0020,
registro de providers y validación SDK usan `mercado_pago`. La consulta anterior
no encontraba un pago Commerce nuevo al preparar una devolución. El cambio
productivo sólo amplía el filtro a ambas grafías exactas y proyecta
`mercadopago` → `mercado_pago` dentro de PostgresStore.Prepare, al construir una
proyección de devolución nueva. No cambia el pago fuente. La recuperación de
refunds existentes ocurre antes de esa consulta y devuelve su identidad intacta.

No se modifican IDs de pago/devolución, referencias provider, idempotency keys,
reglas de importe, concurrencia, estados parciales/completos, SDK ni dependencias.
El cambio es AUTHORED glue de compatibilidad. No deriva un algoritmo de negocio
de Stripe/Mercado Pago ni reclasifica el resto del owner como glue.

## Verificación exacta

- RED auténtico del código anterior: TestPostgresMercadoPagoProviderAliases/
  mercadopago falla con refund source mapping conflict. Nuevos tests sobre fuente
  anterior, sin reemplazar el store por un fake.
- GREEN: las dos grafías pasan con PostgreSQL18.6 real, 54 migraciones de la
  referencia hash-locked y el SDK oficial Mercado Pago1.14.0 usando requester
  local sin red. Comprueba dos refunds de1000 sobre pago2000: parcial conserva
  captured/version3 y total cambia a refunded/version4; incluye GET del refund
  persistido, referencias y keys intactas, selección ambigua rechazada y
  observaciones inmutables.
- Se ejecutaron también los tests existentes de Stripe, ambos SDK, processor y
  guard de DB. Diez tests principales, dos subtests de alias, cero skips. Vet del
  owner y build de todo el módulo PASS; PostgreSQL detenido al finalizar.
- Primer intento conservado como FAIL de fixture: dos pagos de ambigüedad tenían
  la misma referencia y chocaron con el UNIQUE. Se corrigieron los IDs sintéticos;
  no se cambió schema ni producto para tolerarlo. Esa corrida no cuenta como RED.
- El fixture reutiliza el seed/cleanup histórico con bypass de triggers sólo en
  DB descartable exclusiva `elite_refund_test_<unique>`/loopback. Durante los
  métodos del owner los triggers están activos. No prueba el journey de origen
  de autorización de devoluciones ni aceptación provider live.
- Reconstrucción materializer25/25 SHA exactos contra fuentes ya probadas. No se
  repite PostgreSQL después de materializar: sus bytes coinciden. go.mod/go.sum,
  provider.go, processor.go, model.go y las54migraciones permanecen idénticos.

## Gates del delta

G0/G1: identidad de owner y licencia local explícita; pines oficiales heredados
sin adquirir ni atribuir código externo nuevo. G2/G3: una traducción de nombre en
la frontera, conservando contrato durable e IDs; no owner paralelo. G4: RED/GREEN
del store con PG y SDK exactos. G5: tenant/keys/references y ambiguity guard
conservados; ninguna credencial, endpoint o dependencia nueva. G6: filtro de tres
literales y CASE constante, dentro de la consulta bloqueada existente; corpus
finito cubre las dos grafías, no afirma carga/seguridad productiva. G7: no
migración ni reescritura histórica; rollback a0.1.5 conserva datos pero vuelve a
rechazar pagos `mercadopago` nuevos en este flujo. G8: pack25/25 exacto, evidencia
y contrato de mapping materializable. No eleva la admisión del owner completo.

## Identidades

Pack anterior0.1.5 SHA `8d81308a434615f9d8e6de9225e0b2b0aaaae8e3186dff4c00f46b9294221ae0`; pack0.1.6 SHA `64cbc8f96c95c5b42a58e24516741bdc9cf42ce3ab46c2912c112725a1fff460`.

| Archivo de producto | Before SHA-256 | After SHA-256 |
|---|---|---|
| `return_refund_worker/docs/COMMERCE_PROVIDER_ALIAS.md` | `NEW` | `21d6e4003b6c911acdf690288114d1737e5af419296f15fa769792969dbe3e94` |
| `return_refund_worker/internal/refundworker/postgres.go` | `c320c84e90c43cefe67c2fbe0ef86582e06a595b26f9ebbbd45d2f4302dfaf85` | `06a0915c86c29d30516ad05a4ab283a1f4f38624224e7a7736e5737230648ae3` |
| `return_refund_worker/internal/refundworker/postgres_integration_test.go` | `1156f3214404c84d646d87033886869b814b339a0588ee327067e170726d0b36` | `f20113caa74aadd0a3fb28b94ecaf45c6cc49656dd2e7376bafd054b694e5d13` |
| `return_refund_worker/internal/refundworker/provider_alias_integration_test.go` | `NEW` | `4a3a9d2b9ef20ac08006926b3ed48a5d7009171f14dcc700f8e46cf906626fc9` |

Receipts fuera de distribución, bajo `C:/Users/NL/AppData/Local/Temp/elite-v402-library-infra`:

| Receipt | SHA-256 |
|---|---|
| `refund-alias-manifest.json` | `90bfb623d516c22b34edc4d2137811e43ac8a62e96894fea2b9f2154b04de52e` |
| `refund-alias.patch` | `c5590a9db25ee6f2ab1f93895725b67b4b45bd18be0b20bb7443e084f680648f` |
| `refund-alias-pg-1/red-alias.log` | `59be053c4eb8daed9afe587713cb7c96f19085c5766669f133407b18714a8588` |
| `refund-alias-pg-2/red-alias.log` | `2db46df9418a0ca0df5e81793dd0f13c36db9157b4e491734c578c62bc1493a7` |
| `refund-alias-pg-2/connected.log` | `805c54095cf8ef40698311fa1df4941b086f2edf840418a26538a8256127501b` |
| `refund-alias-pg-2/vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `refund-alias-pg-2/build.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `refund-alias-pg-2/result.json` | `c402425d90e1cee64ade3cbdd999545bc98eb77e8d91cad8bc14dd5ecc5388ea` |
| `refund-pack-roundtrip.json` | `bf49109f8917a62e97005e485f6b771c8e2068a019976fae3ee843d9ff99f333` |

Go1.26.8 y PostgreSQL18.6 mantienen el runtime lock usado en BC_EXACT_AMOUNT_ADAPTATION_V402.md; ningún gate live/fiscal/security bloqueado se considera ejecutado aquí.
