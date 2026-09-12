# Microsoft Agents Human Oversight — reauditoría 2026-08-27 V1

## Decisión

`microsoft/agents-humanoversight` es código público oficial Microsoft bajo MIT y contiene una referencia útil de email con opciones, identidad devuelta por el conector Office 365, timeout y logging en Table Storage. No se admite para adopción inmediata: el repositorio fue archivado por Microsoft el 24 de agosto de 2026, la entrada HTTP no implementa Microsoft Entra ID y el propio README ordena considerar agregar Azure AD.

Estado: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

## Identidad fijada

- repositorio: <https://github.com/microsoft/agents-humanoversight>;
- commit firmado: `9eca38d78f7d26027a925ffbf99088d74e7ebcdd`;
- verificación GitHub: `verified=true`, `reason=valid`;
- tree: `2a6e51e38552bd7692c8f728fbefc46e2e43099a`;
- fecha del commit: `2025-04-15T07:47:04Z`;
- ZIP canónico por commit: 2.301.334 bytes, SHA-256 `b82c5b221bba29207f158315ebfafea7f51bce46bcf4f0043e18527cf6fd6b00`;
- tar.gz auxiliar auditado: 2.278.535 bytes, SHA-256 `f31abde6859c74d3f6975768f8f7be0c601c0dc9a2482dde30a7ddb5083fa5d1`;
- licencia MIT: 1.141 bytes, SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- estado GitHub observado el 2026-08-27: `archived=true`.

Artefactos principales:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `README.md` | 12.724 | `2b04ffc33ba1b9c72d45fee01118295fa61ed9e15c7ced8e1769f59e0f469f85` |
| `app/requirements.txt` | 77 | `30fca6fe63acd954131a5d5c502cd9d76f9aa4f956aec0d49dec0cf054d67a8d` |
| `app/human_oversight/approval.py` | 6.871 | `7bfdd8e046fc1347432cc63834acfc2505dd731aa1fbca58c3ae143e341b8c44` |
| `app/human_oversight/decorator.py` | 6.789 | `f5e79c64b20192b3b0e983522c8ef07c390eb702f50d067bd1b819aedf88e86c` |
| `deployment/logicapp.bicep` | 10.866 | `f47f2f16deb77eaeb20f2941f3558ad665806dee7fc06aed119915d20c0922c1` |
| `deployment/main.bicep` | 1.123 | `be9527e2e6f9dc1073dc2801b456c2f17f98472c9e73e58e178d70a7ea81a77d` |
| `app/tests/test_logic_app_integration.py` | 10.632 | `2defba3401b600aa13d535bb88fbbf6c0521f729bde12620c64476efdefcfb4c` |

## Qué demuestra realmente

El workflow recibe una solicitud HTTP, envía `Approve,Reject` mediante el conector Office 365 y toma el email del `responder` devuelto por ese conector. Registra correlation ID, estado, aprobador y timestamps en Azure Table Storage. El decorador bloquea la función cuando la decisión se rechaza, expira o falla.

Esto es una referencia oficial válida para:

- email interactivo con opciones;
- propagación de correlation ID;
- timeout fail-closed del decorador;
- captura del responder que devuelve el conector;
- logging operativo básico.

## Por qué no se admite

La frontera que inicia el workflow es un trigger HTTP `Request`; no contiene política Entra, issuer/audience/tenant/scopes ni autorización por recurso. El cliente sólo publica a `HO_LOGIC_APP_URL`; el README exige tratar esa URL como secreto, recomienda restricciones IP y dice que se agregue Azure AD para seguridad adicional. Además:

- el caller decide libremente `approverEmails`, `agentName`, `parameters` y `correlationId`;
- no se prueba que el responder esté autorizado para el recurso empresarial más allá de haber recibido el email;
- no hay nonce, replay gate, CAS/version ni idempotency key antes de ejecutar la función;
- Table Storage se usa como log, no como transición atómica conjunta con el efecto de negocio;
- el repositorio está archivado y no publica un lock completo de transitivas.

## Gates ejecutados

Sobre CPython 3.12.13 en venv aislado, la instalación mutable resolvió 89 paquetes. La suite oficial produjo:

```text
28 passed
2 skipped
1 RequestsDependencyWarning
```

Los dos tests de integración/live quedaron omitidos sin Logic App, emails y acción manual; por tanto no existe evidencia live en esta auditoría. La CI oficial también instala `pytest`/`pytest-cov` sin fijarlos.

OSV Scanner 2.5.1 sobre el freeze realmente instalado encontró 16 advisory IDs en cuatro paquetes: `python-dotenv 1.1.0`, `requests 2.32.3`, `semantic-kernel 1.28.0` y `werkzeug 3.1.1`. Algunos IDs son aliases del mismo advisory, pero el resultado no es cero y afecta pins directos oficiales.

## Condiciones preservadas

- `UP-FAIL-105`: repositorio archivado y frontera HTTP sin Entra/autorización por recurso;
- `UP-FAIL-106`: resolución no congelada y advisories vigentes;
- `UP-FAIL-107`: integración real omitida y falta de consistencia/idempotencia empresarial.

La fuente se conserva exacta para impedir adopción por título o marketing. Ningún perfil la selecciona y ningún agente puede presentarla como el adapter productivo de revisión humana.

## Fuentes oficiales

- <https://github.com/microsoft/agents-humanoversight>
- <https://github.com/microsoft/agents-humanoversight/tree/9eca38d78f7d26027a925ffbf99088d74e7ebcdd>
