# V292 — reparación de privacidad del logger con política cerrada

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Fecha: 2026-09-07. Clasificación: mantenimiento correctivo aislado, T2803/T2809.
Estado: regresión local FAIL-386 corregida en 0.2.0; **CANDIDATE**, sin admisión de producto.
No reabre decisiones del perfil 67 packs/742 archivos ni agrega servicios/dependencias.

## Qué cambia realmente

GO-OBSERVABILITY-CORE conserva dos archivos AUTHORED. La serialización ahora
invoca Go log/slog JSONHandler del toolchain fijado 1.26.7; no se copió un nuevo
logger empresarial ni se atribuye nuestra política a Google.
El contador corregido V291 se conserva; el histograma no recibe promoción.

- Sólo IDs de eventos y campos exactos revisados; valores enum cerrado, bool o int64 acotado.
- Rechazo total antes del sink si hay mensaje libre, valor no aprobado, objeto
  anidado, campo extra/ausente, callback, tipo incorrecto o nivel inválido.
- Política y enums copiados al construir; snapshot de valores primitivos antes de emitir.
- JSONHandler real produce registros útiles y serializa writes. TryEmit comunica
  rechazo/filtro/configuración ausente/entrega incierta sin reflejar datos rechazados.
- Error o panic del sink no filtra su texto; no retry automático ni promesa de entrega.
- Migración incompatible explícita 0.1.x→0.2.0: NewLogger sin política no emite y
  Redact retorna cero campos. Usar NewJSONLogger/NewPolicyLogger y comprobar
  TryEmit. No upgrade silencioso ni aprobación por “todo verde”.

La prevención no consiste sólo en apagar logs: un evento permitido con tres
campos recorre política → TryEmit → JSONHandler → bytes → parser JSON y conserva
los seis campos esperados (incluidos time/level/msg); heartbeat sin campos también.
El sink recibe 320 registros concurrentes completos, además de los casos negativos.

## Autoridad y procedencia

Consultadas el 2026-09-07; no se presentan como ELITE_REFERENCE admitida de este pack:

| Fuente oficial | Claim usado y límite |
|---|---|
| [Google Android — Log Info Disclosure](https://developer.android.com/privacy-and-security/risks/log-info-disclosure) | Evitar datos sensibles y contenido impredecible; preferir constantes. Aplicamos ese principio al diseño, no las instrucciones Android/R8 a Go. |
| [OpenTelemetry — Handling sensitive data](https://opentelemetry.io/docs/security/handling-sensitive-data/) | Minimizar atributos y revisar qué es sensible según contexto; OTel no lo decide automáticamente. No se instala Collector ni se afirma anonimización. |
| [Go — JSONHandler.Handle](https://pkg.go.dev/log/slog#JSONHandler.Handle) | Serialización JSON y una llamada serializada a Writer por registro. No sanitiza objetos arbitrarios; sólo le enviamos primitivas ya aprobadas. |

Código estándar Go: BSD-3-Clause, adquirido y verificado en el expediente V281;
no cambia el pin de runtime. json_handler.go SHA-256: a10503560615180e018a44a0ebd5ab0c4d0a79f74ab7ee0248ca24277da259df.
Nuevas reglas, límites de 128/32/64 y tests: AUTHORED / LicenseRef-Workspace-Owner.
La página HTML OTel se recuperó después de fallar su representación index.md
(FAIL-414); no se inventó contenido ni acceso al formato fallido.

## Evidencia ejecutada

Staging propio: Temp/elite-v292-f9a8893d99eb4168a997a60d53694b3c.
Sólo fixtures sintéticos. GOTOOLCHAIN=local, GOPROXY=off, GOMAXPROCS=2 para fuzz.

| Paso | Resultado observado |
|---|---|
| Baseline exacto 0.1.1 + tres probes V280 | 3 FAIL esperados: email superior, secreto anidado y mensaje; reproducción de FAIL-386 |
| Candidato 0.2.0 | 16 tests unitarios PASS + 11 semillas explícitas de dos targets fuzz |
| Política negativa | Campos/arbitrarios/ciclos/formatters, límites de configuración y tipos rechazados sin salida |
| Política positiva | JSON real, campos exactos, niveles, heartbeat, snapshot y 320 writes concurrentes PASS |
| Error del writer / panic | Error constante, sin texto privado ni retry, lock liberado; segunda llamada posible |
| Reconstrucción Markdown en destino ausente | 2/2 archivos byte-idénticos; mismos 16 tests y semillas PASS |
| Vet / build reconstruidos | Ambos exit 0 |
| Fuzz candidato, 10 s por target | Logger 1.407.447; contador 1.639.625 ejecuciones, PASS |
| Fuzz reconstruido, 10 s por target | Logger 1.392.989; contador 1.641.402 ejecuciones, PASS |
| Composición con acknowledgeConditions=true | Rechazo CANDIDATE esperado, sin crear destino |

Los totales de fuzz son ejecuciones, no casos únicos ni cobertura exhaustiva.
El cache del contador aporta nueve entradas baseline, siete semillas explícitas.
Pruebas concurrentes no sustituyen race detector. No se ejecutó race, carga,
seguridad ofensiva, SCA nuevo, despliegue ni suite integrada de toda la biblioteca.

## Identidad de artefactos y receipts

| Artefacto | SHA-256 |
|---|---|
| GO_OBSERVABILITY_CORE.md | 82bf34505c3fbdbdb9359f659f40b3f8bc32089a51e14495dd1f9260508d8c03 |
| observability.go | cff998c833279e3d3fc317b7ff91831d52e404f6aab3bdd0b257c6d9c4f2989b |
| observability_test.go | d0b6cbacee3a94b1f34737754427e279914b698786d89788804ae2886624dd2c |

| Receipt local | SHA-256 |
|---|---|
| privacy-before.log | fc733fc51121270961a5c0ff11fd3915557ea9093ecbb5360cd6b18355247bb5 |
| candidate-tests.log | c946eb8236526ca2496b4523baf92fbd2c75abe7c44fa851ad1ba5145d2e9c3e |
| candidate-fuzz.log | 1973813085fc361dd28e7df90778bded6ed8c85cfff3a40cdaa49128dd2c6950 |
| candidate-fuzz-receipt.json | 5aeafd7b6fef3238bdd9964d33509fbb884ae3789ff62ff4ff5d656c8c07ce61 |
| rebuilt-tests.log | fa14fe9e1b9c6810251eefce66a407eb085cad09bcd1e9fb7fa9c10d36c2bc4f |
| rebuilt-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-fuzz.log | 0f13a74eda82c0c735ff61c97abe721cc2defe4670e0d9e3c19fc97437d01e8d |
| rebuilt-fuzz-receipt.json | 5aeafd7b6fef3238bdd9964d33509fbb884ae3789ff62ff4ff5d656c8c07ce61 |
| candidate-rejection.log | 70cc25c23717c523352cc57110d0c5e8dc2006c5cbb815d20e5c8f7bd14b6971 |
| fuzz-profile.json | 6a006ec0f2fe9c7f8092f8016022c114f994cf320a761a780ee626b95e8394ec |

Los logs vacíos de vet/build tienen hash de archivo vacío; éxito confirmado por
exit codes separados con abort ante fallo, no inferido de ausencia de texto.
Para reproducir: materializar el pack en módulo aislado Go 1.26.7; ejecutar
go test ./... -count=1, go vet ./..., go build ./... y GO_NATIVE_FUZZ_GATE
con FuzzPolicyLoggerClosedValues y FuzzCounterMonotonic, 10s cada uno.
Los tests y semillas portables están dentro del Markdown, no sólo en Temp.

## Lo que NO queda cerrado

Política debe ser revisada: permitir un identificador/enum/rango sensible puede
divulgarlo; tampoco bloquea canales encubiertos ni logging externo al componente.
No se autoriza original crudo de conversaciones/documentos como registro de log:
retained-original y derivados de IA siguen separados y sin alteraciones forzadas.

Sink síncrono: puede bloquear o escribir parcialmente; prohíbe reentrada.
Falta montaje del host/exportador y política real, alertas/retención, latencia,
gates de seguridad/integración y validación del histograma. FAIL-386 se cierra
para las rutas reproducidas y el nuevo contrato cerrado, no como privacidad
universal. T2809 permanece abierto. Readiness y assurance integral siguen BLOCKED.
El avance real es un defecto de código reparado y reconstruible, no un porcentaje
nuevo ni certificación de que todo el sistema está terminado.
