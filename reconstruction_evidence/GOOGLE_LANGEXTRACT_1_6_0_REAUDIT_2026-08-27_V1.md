# Google LangExtract 1.6.0 — readmisión ejecutable V1

Fecha: 2026-08-27  
Resultado: **`SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`**  
Alcance: identidad, licencia, artefactos publicados, instalación, suites offline, resolución de dependencias, scan OSV y comportamiento fail-closed de los providers Gemini realtime/batch.  
No demuestra: precisión del negocio, OCR, corpus propio, proveedor live, privacidad, costo, carga, disponibilidad ni autorización para persistencia automática.

## 1. Autoridades oficiales

- repositorio: https://github.com/google/langextract
- release: https://github.com/google/langextract/releases/tag/v1.6.0
- anuncio Google Developers: https://developers.googleblog.com/introducing-langextract-a-gemini-powered-information-extraction-library/
- canal Security: https://github.com/google/langextract/security
- issue realtime sin texto: https://github.com/google/langextract/issues/508
- issue batch sin texto: https://github.com/google/langextract/issues/527
- PR realtime: https://github.com/google/langextract/pull/507
- PR unificado Gemini/OpenAI: https://github.com/google/langextract/pull/534

Los issues y PR indicados estaban abiertos al verificar el 2026-08-27. No se copió ni aplicó código de PR no mergeado.

## 2. Identidad y licencia fijadas

| Superficie | Valor verificado |
|---|---|
| release | `v1.6.0`, publicada `2026-07-02T06:23:27Z` |
| commit directo del tag | `62a25764745c0970f322b04bd89c844e5554bce1` |
| verificación GitHub del commit | `verified=true`, reason `valid` |
| archive oficial | 10,772,919 bytes; SHA-256 `06bf2d9f5098303bb0bc89c7e9d85ec4479b1476831ec5abaa44a0df02db5602` |
| licencia raíz | Apache-2.0; 11,358 bytes; SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| `pyproject.toml` | 4,429 bytes; SHA-256 `1c22d5236397bd799b132fd2041508f3ed6e21ba633b09b83b5026e0d924a32c` |
| Python | `>=3.10`; auditoría con CPython 3.12.13 aislado |
| wheel PyPI | `langextract-1.6.0-py3-none-any.whl`; 150,147 bytes; SHA-256 `ffa2d584c29ea25c73a2723cec756f21a76799e8ab76a1d9d92841e27504b057` |
| sdist PyPI | `langextract-1.6.0.tar.gz`; 142,637 bytes; SHA-256 `ac6412152b173bfb0d6b4e08bf6ace7319b1e734c4f9c6f33d39f0c5915e640d` |

El source contenía 49 archivos bajo `langextract/`; wheel y sdist contenían los mismos 47 módulos/artefactos distribuibles, con 47/47 hashes de payload iguales al source. Sólo omitieron dos README internos. Un wheel reconstruido localmente tuvo hash de archive distinto al publicado, por lo que no se afirma reproducibilidad byte a byte del contenedor wheel.

El commit no contiene `SECURITY.md` raíz. El canal oficial de GitHub Security se conserva, pero esa página no sustituye una política versionada dentro del source.

## 3. Reconstrucción y suites oficiales

Instalación exacta desde el source fijado con extras `openai,test` en venv aislado:

```text
Python 3.12.13
pip check: PASS
pytest -ra -m "not live_api and not requires_pip and not integration"
718 collected
697 passed
21 deselected
90 warnings
duration: 25.40s
```

Las 21 pruebas excluidas requieren APIs/credenciales/servicios o instalación plugin/integración. Las 90 warnings se retienen como superficies de deprecación/compatibilidad; no se convirtieron en un falso PASS limpio.

## 4. Dependencias y vulnerabilidades conocidas

- resolución auditada: 68 requirements instalados;
- `requirements-audit.txt` SHA-256 `1c975b51705a8bc3dd3b381b1ddf4efb96487ea43173b31392eda6746ea1aeaa`;
- Google OSV-Scanner 2.5.1 exacto, binario SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`;
- scan repetido exit `0`;
- reporte SHA-256 `2bc4d70f4ffa455edb9e826ba8ca7ddf9e71c720bfe81144bdc3d8be1ab7d623`;
- cero advisories conocidos en esa consulta fechada.

Upstream no publica un dependency lock. El resultado OSV de una resolución fechada no autoriza reutilizar versiones futuras ni demuestra ausencia de vulnerabilidades desconocidas.

## 5. Gate adversarial de integridad

Código oficial relevante:

- `langextract/providers/gemini.py:371` devuelve directamente `ScoredOutput(score=1.0, output=response.text)`;
- `langextract/providers/gemini_batch.py:622` ejecuta `_extract_text(resp) or ""`.

Un probe offline aislado, sin credenciales ni red, simuló una respuesta realtime sin texto y una respuesta batch con `blockReason=SAFETY` sin texto. Resultado reproducido:

```text
REALTIME_EMPTY_SUCCESS_REPRODUCED score=1.0 output=None
BATCH_EMPTY_SUCCESS_REPRODUCED output='' block_reason_discarded=true
LANGEXTRACT_V1_6_0_FAIL_CLOSED_GATE=FAIL
```

Probe SHA-256: `7c11c263a60ad31308fbb3ac3b23c59247a53cd7e1bccae1859d39b473e0d7c5`.

Esto es material para integridad de datos: un bloqueo/refusal puede aparecer como éxito vacío y perder su causa. Que 697 tests pasen y el scan fechado no encuentre advisories no compensa ese comportamiento.

## 6. Decisión de admisión

LangExtract 1.6.0 queda fijado para que el agente pueda adquirir y estudiar código oficial de Google sobre extracción estructurada y grounding, pero:

- no se compone en un sistema;
- no autoriza almacenamiento automático;
- no se parchea desde PRs abiertos;
- no se presenta como extracción “100% precisa”;
- no se promueve por prestigio, licencia, estrellas, cantidad de tests ni ausencia fechada de advisories.

La readmisión de una release oficial posterior exige, como mínimo: regresiones realtime y batch para no-text/refusal; resultados parciales explícitos; corpus propio por clase/campo; evidencia por campo; abstención/revisión; reconciliación; precisión/recall y costos de error; malware/prompt injection/privacidad/tenant isolation; proveedor live, carga, recovery, canary y rollback.

## 7. Integración en Elite

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.43;
- source lock: 89 IDs únicos;
- nuevo perfil `google-grounded-text-extraction` con 10 entradas obligatorias y 10 bloqueantes;
- adquisición: 8 negativos PASS;
- perfiles: 14 válidos, 5 negativos y 1 positivo PASS;
- pack: 23 archivos materializados;
- gate global: `VERIFY_LIBRARY_PASS` con 46 packs / 417 archivos / 312 Markdown y 23 perfiles;
- auditor ejecutable: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=89 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1`;
- lecciones: 459 fallos locales retenidos; `UP-FAIL-121` conserva el defecto upstream abierto; cuatro temporales de auditoría fueron eliminados tras preservar hashes y resultados.

Conclusión: **código oficial de alto valor, identidad/licencia verificadas y accesible al agente; no es código listo para persistencia empresarial confiable en esta release**.
