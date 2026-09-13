# V402323 — contrato de firma verificado con confianza externa

PROVEN_LOCAL del gate0.2.0. No se firma aún el producto actual ni se cierraTEST07.
Biblioteca205/2333, procedencia2005AUTHORED164ADAPTED164VERBATIM; perfil116/1616.
Tres archivos previos cambian y se añade1testAUTHORED;1612previos exactos.

Corrige separadoresWindows delZIP y rechaza backslash,drive y segmentos ambiguos.
verify_release exige -TrustedAllowedSigners fuera del release recibido; compara
el hash de esa política externa, identidad y namespaceelite-release-v1. Una
política/clave sustituta dentro de un candidato coherente ya no se autoautoriza.
La clave privada permanece fuera del producto; tests usan sólo claves sintéticas.

Build verificaSHA/tamaño/reparse de los inputs antes y después de ambas builds,
fija el perfil consumido y verifica estabilidad antes de firmar. OSV usa exactamente
dependency_manifest_refs, incluyendo owners anidados elegidos por el proyecto;
el caller debe declarar el conjunto completo. Los resultados finales serán reales,
no los JSON/SPDX sintéticos empleados para probar el verificador.

4casos de helpers reales,1ZIP separadores,8casos de firma/verificador PASS; vector
público OpenSSH original/tamper PASS con AuthenticodeMicrosoft válido.1build
deliberadamente fallida confirma evidencia retenida y ausencia de publicación,
sin ejecutarSCA.3rebuilds1616/1616/11 exactos;3archivos lint6hallazgos de pins/
digest revisados,0sintriage/0nuevos blockers. Fuente/app/dependencias sin cambios.

FAIL953 RESOLVED_LOCAL_GATE_CONTRACT. FAIL954 fue un fixture incorrecto: ZipInfo
normalizaba nombresWindows al construir/mostrar; ahora se fija y verifica el nombre
original serializado antes de esperar rechazo. Los dos intentos previos permanecen.
No se relajó el gate ni se presentó el fixture como release/SCA del producto.

Sigue la preparación del driver/identidad/destino y ARCA penúltimo; las dos builds
finales deben contener esa cohorte. Daybreak/libxml2 sólo expediente al final.
