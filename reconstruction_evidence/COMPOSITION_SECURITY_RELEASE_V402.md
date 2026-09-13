# Composición: identidad, procedencia, SCA y lint — V402 / execution312

T2803 PROVEN_LOCAL para LIBRARY_INFRASTRUCTURE, excluyendo el alcance fiscal
reservado al penúltimo paso y Daybreak/libxml2 diferido. No READY global ni
seguridad productiva/ofensiva. Perfil109/1481; biblioteca204/2250.

La identidad/J5 probada en309/310 conserva exactamente su código: PKCE,
refresh/CAS/rotación/logout/revocación y bootstrap acotado con revisión de roles.
El IdP sigue siendo autoridad externa de cuentas. No se piden credenciales.

Se corrigió una regresión del lock: el grafo completo seleccionaba x/mod0.37,
afectado por GO-2026-6179/6180. Se restableció la versión0.40 ya admitida,
sin cambiar fuentes compiladas. Go mod verify PASS;320paquetes/2065archivos
compilados y los dos binarios API son byte-idénticos antes/después.
El grafo resuelve ahora x/mod0.40 y x/tools0.49. Advisory oficial:
https://pkg.go.dev/vuln/GO-2026-6179 y https://pkg.go.dev/vuln/GO-2026-6180 .
OSV2.5.1: cuatro grafos Go, tres locks pnpm y siete locks Python exactos,
73identidades Go/Python y cero hallazgos. ARCA NuGet queda para su etapa.
La admisión477identidades del instalador pnpm restringido311 se conserva.

DevSkim fue reconstruido desde el commit Microsoft firmado fijado y la
adaptación SharpCompress0.48.0 ya declarada:300tests y NuGet SCA0. El ejecutor
local corrigió un índice inválido que ocultaba errores cortos; el test focal
preserva comando y causa originales. El SDK10.0.400 anterior en Temp estaba
incompleto; se recuperó el mismo ZIP Microsoft, SHA512 fijado,5577archivos,
en destino nuevo. No se reemplazó evidencia anterior ni se cambió de versión.

Lint de1307archivos seleccionados:1842avisos originales,1841revisados como
no bloqueantes para este scope y1expediente fiscal explícitamente pendiente.
No hay avisos sin revisar ni suppressions en el scanner. Los hashes públicos,
fixtures, mínimosTLS1.2/1.3, callbacks constantes, UUIDv5 de identificación y
listeners locales tienen disposición por ocurrencia y SHA de fuente exacta.
Los UUID internos no son secretos ni autorización; RFC9562§5.5.
No se convierte un PASS de lint en análisis interprocedural u ofensivo.
El finding1539 de internal/fiscal/wsfeipc/provider.go se reserva para probar
su transporte UDS en ARCA; no se declara resuelto ni faltante de credenciales.

GO_ENTERPRISE_BACKEND_CORE0.4.6 y MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE0.1.1:
tres archivos revisados,1478anteriores intactos; cinco perfiles reconstruidos
exactamente. El mapa por archivo conserva1337AUTHORED/108ADAPTED/36VERBATIM,
sin inventar procedencia. FAIL921–926 resueltos localmente; evidencia RED
conservada. Próximo paso obligatorio: T2806 pipeline documental de referencia.
