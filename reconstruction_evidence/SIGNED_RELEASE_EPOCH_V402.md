# V402326 — fecha UTC del release y destino durable

FAIL972 RESOLVED_LOCAL. El gate0.2.0 interpretaba la fecha fija del ZIP en UTC−03 y entregaba315543600 al builder, que exige315532800. Gate0.2.1 usa AssumeUniversal. La regresión evalúa la asignación real del gate y comprueba offset0/epoch315532800 en Windows UTC−03.5casos del builder,1ZIP y8verificaciones firmadas con fixtures PASS; no son una firma del producto.

Un rebuild completo independiente coincide en1653archivos;1651permanecen idénticos. Los únicos2bloques cambiados son glue/test AUTHORED del gate. Los4planes que seleccionan el gate fijan0.2.1; los inventarios derivan de sus bloques canónicos. Sin cambios de dependencias ni algoritmos de negocio.

Destino durable creado fuera de la biblioteca: Desktop/Elite Franchise Reference V402/reference-staging. BridgeBoth y materialización NEW PASS; fuente durable actualizada por delta exacto326. Dos builds completos, SCA, SPDX, firma externa y aceptación siguen en curso; no se declara TEST07 ni readiness final.
