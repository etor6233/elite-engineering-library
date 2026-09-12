# Operational evidence gate — V315, 2026-09-08

Mantenimiento T2803/T2809, FAIL-20260908-503. SECURE-OPS-DELIVERY-CORE1.1.2,
código AUTHORED y admisión CONDITIONED. Dos archivos del owner reconstruidos
idénticos desde Markdown; composición67/745. No dependencia ni regla nueva de
negocio. El contrato ya exige SBOM/provenance y el modo files comprobar archivos.

## Fallo reproducido y corrección

El validador sólo inspeccionaba referencias nonempty_string y no exigía las dos
rutas de supply_chain.sbom.artifact y provenance.artifact. Omitir/null/vacío/
placeholder/número/lista/objeto saltaba el control. Se reprodujeron28 falsos PASS:
siete variantes ×dos artefactos ×dos modos, contra control completo.

Ahora ambos campos son obligatorios en plan y files; las siete referencias
del modo files se recorren por nombre, sin omitir valores inválidos. Fallos
OSError/ValueError/RuntimeError al resolver o inspeccionar devuelven error
del campo, sin traceback ni volcado del path. Se conserva restricción de raíz
y archivo regular. La CLI emite FAIL/exit2 al faltar las referencias.

El primer rojo tuvo29 fallos, uno debido a un oráculo demasiado específico:
NUL en Windows/CPython3.14.4 ya devolvía file-not-found, no abortó. Ese test se
corrigió para exigir rechazo y se agregó PermissionError inyectado explícito.
El rojo refinado confirma28 falsos PASS y una excepción por acceso; ambos logs
se conservan. No se atribuye un problema de permisos real a esa inyección.

## Verificación

- Doce tests en candidato y doce en reconstrucción PASS: seis anteriores más
  seis sobre el límite de evidencia. Incluyen28 variantes negativas, control
  completo, ausencia/directorio, escape a archivo existente, NUL y fallo de acceso.
- CLI reconstruida: expediente completo PASS/0; dos referencias omitidas/nulas
  FAIL/2 con exactamente dos errores; NUL FAIL/2; template sin resolver FAIL/2.
- Identidad2/2 de validator/tests; canonical.ps1 actualiza owner1.1.2 y los tres
  planes que lo seleccionan. No cambia el conteo de archivos ni otros módulos.
- No consumidores de los mensajes antiguos encontrados en production_admission_gate;
  no se transfieren suites productivas como si se hubieran ejecutado en este corte.

## Alcance y límites

Los archivos sintéticos contienen {} para probar presencia. Ese PASS no verifica
schema, contenido, hashes, firmas, frescura, SBOM real, SLSA o aprobación. El gate
productivo separado sigue exigiendo sus ocho receipts ligados al release/target.
No hay monitor runtime, alertas, supervisor, despliegue o recuperación demostrado.
No cierra T2803/T2809 ni cambia ronda D/canal histórico, decisión comercial, V295
o ARCA diferida. Sin ZIP/promoción ni porcentaje inferido. TEST34/EVID25 enlazan
el corte; assurance plan PASS, evidence integral BLOCKED.

## Reproducción y receipts

Staging `%LOCALAPPDATA%/Temp/elite-v315-13cbdf9f4f3b4d18b51a0d76e6ffefc5`.
CPython3.14.4 Windows. Componer FRANCHISE_COMPLETE_PACK_PLAN en destino ausente;
ejecutar `python -m unittest -v test_validate_operational_readiness` desde
ops/readiness. verify.py conserva cuatro invocaciones de CLI con
`--require-evidence-files --root <fixture>`. Red.py registra el primer fallo;
baseline conserva bytes anteriores. No repetir scripts de mutación sobre staging.

| Receipt | SHA-256 |
|---|---|
| red.py | 63a3f30bf04b0c55c5dc68a81d2e660f78dc5cc032ac6d9bc78aa48a1ed06015 |
| readiness-red.log | 78b9f118b300e787999187e2b6aa6f691115b78fecbcdbacd885e2cec040f5b6 |
| readiness-red-refined.log | 9a676d5303ebb8599d36113d36abc816a18a435a043af7cc1d3f173d6a62685c |
| readiness-green.log | f4219e2cebbf26b400aa65f22623072d16e92f4bf0c4af53b7e760d1c1b629c7 |
| readiness-rebuilt.log | 9431db2c4d762df30122f3b2513d589b6379be1b41ab9f44ca54deb62a0e0488 |
| verify.py | f89bc9b794f2d46e30db9d4bcb79d049102a50c3d1158a8dacba8cd94acb5855 |
| cli-receipts.json | 00052964652b18dafe183227b3845f530a44781b8e8f159719df4f8f98370870 |
| canonical.ps1 | 2fcfbae94230ae788c00690e747637fc2fd539889964bef32b09a5a091bd6a15 |
| compose-rebuilt.log | fde6191ed35adebedd7a17b6b911a6b47a382f361535ccd0df529770c49059a1 |
| baseline-pack.md | f9fbd06b522a7f4ea9248fe2fee29b3a23a1c8aef5f1d2c78bd2ed33dcf6c8a5 |
| rebuilt/ops/readiness/validate_operational_readiness.py | 950e8ff3aa63dea22b10c6bccf0501a5b57b3139fb0df51c8e7d65d79ce799ac |
| rebuilt/ops/readiness/test_validate_operational_readiness.py | 7030fa5a1294d237eb8ffdafdc97fb6d092bc260198c7ecb82b9e2fa41907098 |

## Cierre general

VERIFY_LIBRARY_PASS160 packs/1436 archivos/748 Markdown y51 perfiles; franquicia67/745. Checkpoint57 validado antes del verificador; TEST34/EVID25 plan PASS y evidence integral BLOCKED. Log library-final.log SHA256: 0501c23e3b1fec04af65cc42cce63c22a79d873bde6175a831b70c6971bd2237.
