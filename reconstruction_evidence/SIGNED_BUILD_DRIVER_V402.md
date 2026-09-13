# V402324 — driver de build firmado e identidad local

PROVEN_LOCAL del contrato de orquestación; T2810/TEST07 esperan el artefacto final.
Perfil116/1621 exacto; biblioteca205/2338:2010AUTHORED164ADAPTED164VERBATIM.
Cinco bloques nuevos de glue/tests/template/docs y una guía revisada;1615salidas
previas intactas.4rebuilds exactos1621/644/965/1621.12fixtures de contrato PASS;
3archivos de código DevSkim0hallazgos. No se repiten suites de negocio intactas.

El driver ejecuta dos veces el builder real, con instalaciones y caches nuevas,
misma ruta de workspace por Next16.3.4 y archivo de cada workspace terminado.
Lock exclusivo, máximo2builds por configuración; rechaza un intento parcial.
La proyección compara recibo PASS, inventario y corresponding source antes/después
de copiar. No firma artefactos históricos ni presenta fixtures como builds actuales.

Identidad Ed25519 local creada para el candidate autorizado, sin identidad
corporativa/productiva. Clave privada fuera del proyecto y sin imprimirse/copiarse.
Fingerprint público: SHA256:frxj2lcV7N0+H8tEJA1Ss8ToAD1KkyjBwVVHkKtPhMs.
FAIL955: OpenSSH agregó Administradores al archivo; la comprobación lo detectó.
Corrección DACL nativa conserva la misma clave y permite sólo usuario+SYSTEM;
verificación protegida en carpeta/clave/público/política y correspondencia PASS.
Set-Acl solicitaba un privilegio de auditoría no disponible; se modificó sólo DACL
con icacls. Evidencia anterior retenida, sin ampliar privilegios ni pedir secretos.

Sigue ARCA penúltimo: contratos/UDS/.NET/Go y fixtures; luego las dos builds finales,
SCA/SPDX/firma y aceptación durable NEW/EXISTING. Daybreak sólo expediente último.
