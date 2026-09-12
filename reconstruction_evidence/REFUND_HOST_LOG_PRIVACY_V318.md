# Privacidad del log del host de reembolsos — V318, 2026-09-08

Mantenimiento T2803/T2809; FAIL509 corregido en el owner
GO_OFFICIAL_RETURN_REFUND_WORKER0.1.3. Tres archivos AUTHORED: main, regresión
nueva y README. No nuevo upstream, dependencia, cambio de retry o regla de refund.
El método reutiliza la salida estática ya aplicada en GO_ELECTROMOBILITY_APPLICATION;
no incorpora el logger CANDIDATE de observabilidad ni atribuye el fix a un SDK.

## Defecto reproducido y corrección

El binario refund de V317/Go1.26.7 recibe un DATABASE_URL sintético de puerto
inválido: falla antes de conectar y devuelve exit1, pero log.Fatal(err) publica
PRIVATE-SYNTHETIC-DATABASE. La ruta del loop también imprime el error completo
de Step, que puede transportar causas de DB/provider. No se afirma una fuga
real de credenciales: todas las entradas usadas fueron sintéticas.

El host ahora emite REFUND_HOST_FAILED al terminar por fallo, conserva exit1,
y REFUND_STEP_FAILED ante error del paso. Causas arbitrarias no se formatean en
estas salidas. El loop conserva comportamiento de idle/retry y consulta del
resultado durable; no copia errores a otra tabla ni añade logging de payloads.
La ayuda del mismo pack explica los códigos y la consulta operativa pendiente.

## Regresiones y evidencia conectada

Antes del fix se extrajo el bloque existente a logStepError sin modificar su
comportamiento, sólo para poder probar la frontera del host. Dos tests raíz
fallaron: un error cuyo Error() deja una marca fue formateado y escapó al log;
el proceso hijo con main real divulgó configuración sintética. El log del rojo
se conserva. No se toma un fallo de compilación como prueba de privacidad.

Después del fix, las dos regresiones ×3 pasan en candidato y reconstrucción:
6PASS por árbol. Verifican que no se invoque Error(), no aparezca el marcador,
se emita el código fijo y nil/ErrNoWork/wrapped ErrNoWork permanezcan silenciosos.
Dos casos de startup usan proceso hijo, entorno mínimo y deadline10s: falta de
configuración y DSN inválido deben conservar exit1 sin publicar la causa.

Go1.26.7 observado por ejecutable/compiler/GOROOT/SHA256/GOTOOLCHAIN=local;
test/vet/build -mod=readonly y mod verify PASS en ambos árboles del módulo.
6 tests raíz PostgreSQL del refund PASS en base local exclusiva protegida por
guard. El flujo de asignación/parcial/total/ambigüedad conserva sus oráculos.
El binario final conserva buildinfo1.26.7; repite el sondeo sintético original:
exit1, código fijo y marcador privado ausente.

Reconstrucción3/3 idéntica desde el owner; composición67/746. El nuevo test
eleva el inventario a1437 archivos; notices AUTHORED1199/ADAPTED133/VERBATIM105.
Los locks, migraciones y providers son idénticos; no se repiten navegadores
por un cambio aislado del log de este proceso. V317 conserva su evidencia
por snapshot, no se transfiere como una ejecución de V318.

## Límites y continuidad

Se corrige exclusivamente la salida explícita del host: no prueba ausencia
de datos en logs de terceros, entrega/retención/alertas del target, trazas,
supervisor, permisos de operación o SCA integral. El detalle durable existente
mantiene sus propios controles; no se inventa un canal nuevo para la causa.
El guard de DB usa la base de refund aislada, no el audit compartido V296.
Sin provider live, pagos, redrive, ZIP/release ni cambios comerciales.

TEST37/EVID28: plan PASS; evidence integral BLOCKED. D/canal histórico,
liberación comercial y target/IdP/providers pendientes; ARCA diferida.
Continuar T2809 por owners y cerrar sólo lo demostrado.

Staging `%LOCALAPPDATA%/Temp/elite-v318-993950e79c854c4c9b26938962ffbdd3`.
probe.py conserva el binario before; prepare-red.py/canonical.py son mutaciones
de una sola ejecución. Reproducir en destino ausente con composición canónica,
Go exacto y dos regresiones TestRefundHost; rebuilt-gates.ps1 describe la suite
y la DB opt-in. No DROP ni datos empresariales.

| Receipt | SHA256 |
|---|---|
| canonical.py | 3fafefbc357ceb43349e668f57b355001e4d69b8d30f219661a167ce69d5709f |
| checkpoint.py | e7adac09ae2de688874411aea92be7694083eed1d480c310ca6675f903cbec4e |
| finalize.py | b8e0d08c5e2428170803af8dd86bb520f0ee5c779fdb83af2294fd6d211d56e0 |
| host-green.log | fa98b266963f9ed5b8d11f19f3461e4aa2d020eca4af446875af35ad957b7263 |
| host-red.log | 4a1c3be4afb14f0231104d14692c58a0433063e2c1b862b01c61e84607309c1f |
| main-before.go | a1efc8563b655ba34bdd4a8cfb4314b7b6d4e51871374d4215f70e24bdad6384 |
| prepare-red.py | f5b6339c9260adc5e7aa424b58dea6990cd70fd5a40c5f4b74396d0c009d2167 |
| probe.py | 01104a1c788eb62db09537a1296aad80059dbdb8a99e00e10c76bdf8417eeb07 |
| rebuild-identity.json | fb1b1e473d7d56c4b29ce8d5070147070d797cec134a779f503e6a8860d28852 |
| rebuilt-binary.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-buildinfo.log | a14bcb2b103bf4dba88c7e2cbb1a47e7a687fa87c8dfe56c4739318808f52b76 |
| rebuilt-gates.ps1 | d715604cfb5f6980d6f4e029bfdeb885f6bbd2854fbd60818a38c66fde916ce1 |
| rebuilt-host.log | 58e307226abab1f9763517918cd83bc9c206f097cb0f03f0d9e5313d4856d42a |
| rebuilt-mod-verify.log | acdb6a6a98dc2c297af31bb3538319778be8c8ef263cd0a8c9799de9c4998533 |
| rebuilt-postgres.log | c9171e3f2a9ab00520af3dd23b8d5ee87cc7b6ea1e135fc220d184bb488c69c6 |
| rebuilt-test.log | bf7b817e5d5272eebac40ba4381c16db054f173276d4777d6628d67b52a42e78 |
| rebuilt-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| refund-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| refund-mod-verify.log | acdb6a6a98dc2c297af31bb3538319778be8c8ef263cd0a8c9799de9c4998533 |
| refund-test.log | a3f4f41db609f537ef97a5370b9b7b23696f9d1d26b6a761d0a149d741364b1a |
| refund-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| startup-after.json | 908e76a51f2f0cd68976aa72b4f52d7195d1560878a2b73b599d281f28c52fed |
| startup-after.log | 33986421b68193e750a5d4dda3fbd674c25225963055717364d322d66f1ea2cc |
| startup-before.json | 93148f4f7b9583363b0a7690688fe8e79178841f230353c68794c9778c975302 |
| startup-before.log | b73773607c4a90744d83f09e2a91c6cd2bf66d441c8561e2197eb4b41c2127be |
| toolchain.json | 2345511cded485899c61f18587de4374e84e607efbd6f8e04e28c8e63fc6ce27 |
| rebuilt/return_refund_worker/cmd/return-refund-worker/main.go | 1b3af4d08058ecb7a9034d923f7bbe33d1806de02e16fb503a8011f3c857598d |
| rebuilt/return_refund_worker/cmd/return-refund-worker/main_test.go | 36302b9886acad8d5bf25ccaa9b311356a19d100209883eddd21948db21a3045 |
| rebuilt/return_refund_worker/README.md | 59c45c3694e32c63dbe01e473ead71dacaa9026da1313bd495a3c7ec90ca9bff |

## Cierre general

VERIFY_LIBRARY_PASS160 packs/1437 archivos/751 Markdown y51 perfiles; franquicia67/746. Checkpoint63 validado antes del verificador; TEST37/EVID28 plan PASS e integral BLOCKED. Log library-final.log SHA256:58193898318e222d6bc7774c5efb70fa12b0598b9377219ae8307f9f4d8d42a6. Checkpoint64 conserva89 refs y12 must_read; continuar por owners sin promover el conjunto.
