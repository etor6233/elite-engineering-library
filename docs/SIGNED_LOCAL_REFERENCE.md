# Firma de la referencia local

AUTHORED glue sobre los owners admitidos Go, Next, pnpm restringido y OpenSSH.
No representa identidad corporativa, credenciales productivas ni un build service
aislado. La clave privada y el almacén de herramientas quedan fuera del producto.

El caller genera `.elite-local-build-inputs.json` desde el template, fija el
ejecutable Python y el inventario exacto materializado, y declara un `run_root`
externo ausente. Incluye configuración, inventario, driver y todos los manifests
de dependencias elegidos en `source_inputs_ref` del gate firmado. No contiene
secretos. El perfil del gate fija el SHA del driver `ci/build_signed_reference.ps1`.

Cada invocación ejecuta realmente `ci/build_local_reference.py`: copia limpia,
instalación offline restringida y GOCACHE independiente. Usa la misma ruta explícita
de workspace por el contrato de IDs RSC de Next16.3.4; archiva cada workspace terminado
con sus logs e inventario y nunca lo reutiliza como entrada compilada. Un lock
exclusivo permite exactamente dos builds por configuración. Fallos parciales se
conservan; requieren diagnóstico y un nuevo run_root, nunca borrado automático.

La proyección compara el árbol actual contra el recibo PASS y su inventario antes
y después de copiarlo al output del gate. El gate compara ambos ZIP deterministas,
escanea los manifests explícitos con OSV fijado y firma la procedencia. Su verificador
independiente exige `-TrustedAllowedSigners` desde una política confiada externa al
release recibido. El archivo de firmantes incluido en el ZIP no se autoriza solo.

El timestamp de ZIP es 1980-01-01 UTC. El claim de reproducibilidad se limita al
host/toolchain/arquitectura/workspace admitidos; no implica builds idénticos en
rutas arbitrarias o plataformas diferentes. Las fixtures de auth/build son públicas
y sólo locales. No se habilita deploy live con ellas.
