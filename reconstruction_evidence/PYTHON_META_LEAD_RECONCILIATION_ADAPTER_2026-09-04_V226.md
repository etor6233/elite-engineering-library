# Python Meta Lead Reconciliation Adapter — evidencia V226

## Resultado

```yaml
pack: PYTHON-META-LEAD-RECONCILIATION-ADAPTER
version: 0.1.0
authority: SUPPORTED_REFERENCE
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
runtime: CPython 3.14.4 Windows x86-64
sdk: facebook-business 26.0.1
graph_api: v26.0
tests: 10/10 PASS
compileall: PASS
pip_check: PASS
materialization: 8/8 PASS
```

## Procedencia exacta

- SDK oficial: `facebook/facebook-python-business-sdk`, release 26.0.1, commit `788f363d15b1269ab5efb7cd00fb5e3b133cd99b`.
- Wheel oficial: `facebook_business-26.0.1-py3-none-any.whl`, SHA-256 `41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca`.
- `LeadgenForm.get_leads` conserva el edge `/leads`; `LeadgenForm.get_test_leads` conserva `/test_leads`. La prueba inspeccionó ambos métodos en el SDK instalado exacto.
- Los campos admitidos derivan de constantes de `Lead` en la misma revisión.
- El wrapper es `ADAPTED`; validación, normalización, staging y receipts son locales. No se presenta como código escrito por Meta.
- La licencia de plataforma Meta y las licencias del grafo de 18 wheels permanecen obligatorias.

## Invariantes demostrados

1. El perfil empieza bloqueado y exige app/business, acceso a leads, Page/form, test contract, PII, retención, cuota/costo y reconciliación.
2. `TEST` llama sólo `get_test_leads`; `LIVE` llama sólo `get_leads`.
3. La respuesta exacta exportada por el SDK queda separada del lote normalizado.
4. El batch conserva tenant/organization fuera del candidato nullable para aislar también la evidencia rechazada.
5. Un campo de formulario no aprobado conserva raw pero rechaza el candidato.
6. Un campo con cardinalidad distinta de uno conserva todos los valores en raw y se rechaza como ambiguo.
7. Ningún candidato concede contacto ni escribe al CRM automáticamente.
8. Form ID y PII no aparecen en el receipt.
9. Un límite alcanzado declara truncamiento; no se presenta como reconciliación completa.
10. Error de provider o output existente falla sin publicar un resultado parcial.

## Hashes de archivos materializados

| Archivo | SHA-256 |
|---|---|
| `fetch_leads.py` | `5ab4269967b9d12065ea007a49206ef503b4c6e41638b19de0eb579920b71876` |
| `provider-profile.template.json` | `f3417cf84c7b4865863cea224bb5e30283aa82bfc5285303bf67df3c7b396833` |
| `README.md` | `39d7bd2931409e3fcfacffec0509861e49931a72f3a5ec5a11a0ec535e95b791` |
| `requirements-direct.in` | `b0fe2f0d83c899fd9986a9990b9102b613774bd912aaffd5219736b014051c88` |
| `requirements-windows-py314.lock` | `a4b4b65d2a19be4bb51892bcb25a4e8127f7b2c16b8f49a2f40b319797314e19` |
| `retrieval.template.json` | `c11aee0e9c8d59619a63d32bbf1d048114494582345375d9340bc80b3976e462` |
| `sdk-artifact.lock.json` | `03eb5d7bd845cd02d42045e1d74ee307976cc31eadc5b8c4925b222f86c5a887` |
| `test_fetch_leads.py` | `0fcb84c03bfc99c7fec1b484b2a19b822f56ac0e56e84547ad83b6ca99ed20b8` |

Pack Markdown SHA-256: `3ba9176d98ea0a2dc0b5020ccaa063de77ccf79159f518b93eaf87e3f221e03f`.

## Límites retenidos

- No hubo credenciales ni cuenta Meta real; acceso, completitud, retención provider, cuotas y costos siguen condicionados.
- No se admitió un schema público exacto del webhook Meta; este pack recupera y reconcilia por SDK, no recibe notificaciones.
- El lote todavía debe entrar al owner PostgreSQL central y pasar la promoción explícita V225.
- No prueba consentimiento, contacto, conversación, cita, venta ni postback de conversiones.
