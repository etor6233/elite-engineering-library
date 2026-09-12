# V400 — encuestas conectadas de cliente

Estado: recorrido AUTHORED verificado en referencia local; admisión CONDITIONED.
No se modifica45/48 ni se declara lista la biblioteca. TEST02/03/07 permanecen
abiertos con los mismos requisitos. Daybreak/V386 sigue diferido por el usuario.

## Resultado concreto

GO-CUSTOMER-SURVEY-API0.1.0 y TS-CUSTOMER-SURVEY-PORTAL0.1.0 incorporan23archivos
nuevos: dominio, PostgreSQL, API, activación explícita, retención acotada, BFF,
formulario, lectura propia, resultados administrativos, pruebas y guía.
GO-ELECTROMOBILITY-APPLICATION1.9.4 monta la capability sólo al habilitarla.
GO-HTTP-METRICS-REFERENCE0.1.14 actualiza únicamente el vínculo al perfil padre.

La respuesta se liga al sujeto/tenant/org autenticado y a un cliente CRM activo.
Una respuesta por encuesta/persona; igual envío recupera recibo, diferente
conflicta. La base evalúa el plazo después de adquirir locks. Una definición
inactiva sin respuesta propia no se publica. La interfaz conserva cero, exige
aceptación explícita y ofrece recuperación GET después de una respuesta perdida.
El permiso surveys:read es independiente de customer:self.

## Evidencia ejecutada

- HTTP real, RS256/JWKS y PostgreSQL: identidad, scope, consentimiento, conflicto,
  24solicitudes concurrentes (1alta/23replays), respuesta perdida y cierre durante
  espera por lock. Reinicio real conserva el snapshot y permite recuperar3recibos.
- Next de producción con JWE y API real: Chromium escritorio/móvil, Firefox y
  WebKit,4proyectos PASS;8respuestas durables y8POST. Doble submit no duplica;
  descartar una respuesta HTTP después del commit se resuelve por GET, sin POST.
- Resultados por permiso, umbral explícito, NPS de2respuestas igual0; contenido
  HTML de la pregunta se muestra como texto. Capturas preservadas e inspección
  visual móvil realizada.
- Migración0054 up/down/up en base descartable vacía;53migraciones previas
  byte-idénticas. Retención oculta lecturas expiradas, borra como máximo1por
  invocación de prueba y conserva otro org/filas futuras. CLI real devuelve1/1/0.
- Go unitarios focales, vet y build del módulo PASS; gate nativo Go ejecuta
  baseline completo y11semillas/6219ejecuciones fuzz con presupuesto5s, PASS.
- Typecheck, build Next y38pruebas focales de contrato/BFF PASS. Dependencias,
  versiones, go.mod/go.sum y package/lock/workspace web no cambian.

No se reejecutan las21suites de cores sin delta. GO-SURVEYS-CORE0.1.2 permanece
intacto; el nuevo owner conecta el claim descrito sin atribuirlo a código de una
empresa ni certificar equivalencia empresarial completa.

## Assurance, procedencia y condiciones

G0 necesidad: T2802/T2804 y fila GO-SURVEYS-CORE de FRANCHISE_GAP_MAP.
G1 autoridad: fórmula de Bain y semántica oficial PostgreSQL; no código copiado.
G2 revisión: SQL parametrizado, locks cliente/definición, scopes de servidor,
gramática cerrada y cuerpos acotados; versión/score del recibo se contrastan.
G3 licencia:23archivos AUTHORED, LicenseRef-Workspace-Owner; dependencias heredadas.
G4 build y G5 pruebas: evidencia anterior, separada de aceptación productiva.
G6 reconstrucción: pendiente de anexar receipt exacto antes de cierre del lote.
G7 operación: CLI de retención probado; schedule, grants, carga, monitoring,
backup/erase y política real deben demostrarse en el target.
G8 admisión: CONDITIONED; no se declara USE_REUSABLE_PACK ni pase productivo.

La definición, pregunta de recomendación, aviso, fechas, mínimo de respuestas y
distribución autorizada del enlace son del proyecto. No incluye motor de campaña,
prueba de compra, representatividad, consentimiento legal, retención de backups
ni decisión comercial. Los fixtures son sintéticos. Los flags están apagados por
ausencia y no conceden permisos. El down migration destruye respuestas y no es
rollback productivo sin pérdida. La seguridad/licencia del grafo activo y release
firmado continúan bajo TEST03/TEST07; no se ocultan con este PASS.

## Correcciones y evidencia retenida

FAIL797: parámetro UUID/text del fixture se separó; r1 rechazado conservado.
FAIL798: r2 aceptaba201 después del plazo durante espera por lock; corregido en
la condición SQL posterior al lock, final1 rechaza404 y conserva0filas nuevas.
FAIL799/FAIL719: parámetros/rutas no observados causaron errores sólo de tooling.
FAIL800: next typegen añade dos imports generados en next-env.d.ts; comparación
detuvo publicación, delta explícito conservado sin alterar código productivo.

Stage relativo: Temp/elite-v399-8041bdb5480e4648b437b93f98d667e8; pointer local elite-v399-current.txt.
Los archivos temporales no son autoridad portable ni claim de runtime permanente.
La biblioteca conserva las fuentes exactas en los packs; el record local liga
receipts y fuentes. Reconstrucción canónica y checkpoint se anexan al completar.

Autoridades consultadas2026-09-11:
[Bain: medición NPS](https://www.netpromotersystem.com/about/measuring-your-net-promoter-score/)
y [PostgreSQL18: locks](https://www.postgresql.org/docs/18/explicit-locking.html).

## Hashes de receipts

| Evidencia local | SHA-256 |
|---|---|
| feedback-final1/result.json | 3994c0c3701d241132ea9d3e19d708bcf79053d19e29a1fdcd8b84a7e0051ad2 |
| feedback-final1/connected.log | 54b441c14a32aee67e0c4b510e9e7466789e4def2c1e13683fc961f3ba61c9d7 |
| feedback-final1/browser.log | 879b4e8e9f0c8531f3be2dee2c69c06e03240642aeb2cdaa065f4c0284a0eb94 |
| feedback-final1/recovered.log | 94d23bf88aad7e4689174986d86ee4fe872cbb28988b58253e12fdff2ecf25c0 |
| feedback-final1/vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| feedback-final1/build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| feedback-final1/retention-cli-0.log | e43379283a0c68fe06ecb2d10869806f35c3589efeecc8432197831b85815224 |
| feedback-final1/retention-cli-1.log | e43379283a0c68fe06ecb2d10869806f35c3589efeecc8432197831b85815224 |
| feedback-final1/retention-cli-2.log | 006ff62d9ca57eeb4b0d3e88e15574c0e56d7ecb9144bc45709846ef575f26b1 |
| feedback-web-unit.log | b59bf4b9ed991ec6c0ddf535ae5899cc44c1153b72f93145e31092aee9bc7632 |
| feedback-web-types.log | 99d7c2179cdd316467c51f37901025a0753453664d741214ba56fc631fa6c0bc |
| feedback-web-build.log | 917c1d4ee89d458ce805c9cad57b6ec82bd429664a9cc5f5a492cdcdc172aeb1 |
| survey-fuzz.log | deec48598b335c030ab744c075af7627c6fac9927ace541e56df486bf9b79a56 |
| survey-fuzz-receipt.json | ac43d502eb3a236754d2d902039718aee05a2409091751d887981eb3f1673b4c |
| next-env-generated-delta.json | abe05d431a4e85b2c1b500dea0792a8ba5bd8275d58c23f28cb6c7e70d0e8aff |

## Fuentes nuevas y activación

| Archivo | SHA-256 |
|---|---|
| internal/customerfeedback/service.go | 1d11d2c4a1f8388c18ba46e7510266f2b4c8bce8654cf23dfc28fc6aa7642b79 |
| internal/customerfeedback/service_test.go | 0aee215524094374b0976caaa7ca8426b5c2cd178dffa9318374b017c478b916 |
| db/migrations/0054_customer_feedback.up.sql | 666b7ebb74a7851d9cd791d6937ae020fbea0c1b904a8d8a1cb684481875f347 |
| db/migrations/0054_customer_feedback.down.sql | 2177f10b87b0a0fe196c986687f3e9ad71c5d97ce12c0ce166b26620acd8d9c0 |
| internal/platform/postgres/customerfeedback.go | 9be6e41b375ba1fe28259a167e12de0940b2d2603ec0e218e658f0bbf860d1aa |
| internal/platform/httpapi/customerfeedback.go | e6560b6e36912494580ea5cff17a863039d8b70280cd7a6777a090890f86630d |
| internal/platform/httpapi/customerfeedback_test.go | 4647abffd3a756221a9e73ad6241d42ba587599961613a6a7e526feba22abdf1 |
| internal/platform/httpapi/customerfeedback_integration_test.go | a82893a0c0ad50cdb6ba4e76fcc4985a7c7328a28cce12c0edccdc1dcefcf1a0 |
| internal/platform/httpapi/customerfeedback_browser_test.go | cecf6fcd1fc5caf26ff8c4cb3a12724dd1085af2566f32345255e7885c86bdf6 |
| internal/platform/httpapi/customerfeedback_fuzz_test.go | 348d0994d660894e4f15836e9efedd8c85d41b17d65d9ab14a129d933be06e01 |
| cmd/electromobility-api/customerfeedback_activation.go | 8f0d7811b5e75b56e0c219fbb06a9f41d3e6ff05e4f5c8ed46ba4fae74fad0ef |
| cmd/electromobility-api/customerfeedback_activation_test.go | e50299271f65cefc82764643dde517fa7f97aab0f8914b40dfe8309c1db013ba |
| cmd/customer-survey-retention/main.go | f492d6493d93e1d79fc85f901541657d35c58184883cf09d6eb52bf7dd7b6c68 |
| cmd/customer-survey-retention/main_test.go | 744090c0f7f009c4fc0bab147c278d8949f9d408c31046463db417a3f2421552 |
| docs/customer-surveys.md | 1bfa151b9d9c3d9f09f8982d44cfbe3743b74c4c89ccc4fbda7c34cb78f7890b |
| src/platform/surveys/contract.ts | eba644ae21718d6033387ba27f96e9c4fda5c37ebcddd299f252a5e414351d2f |
| src/platform/surveys/contract.test.ts | cdc6631a3cf220477bcedaff9326beda591f37e2a1a5e54d7334c02f6ce6c1d9 |
| src/app/api/enterprise/surveys/route.ts | 3070da7d8a67af77bc3d8050edc12b16f28a2a0fee3a76edf67db0938afc9941 |
| src/app/api/enterprise/surveys/route.test.ts | f28c869694c0672d515ee78d34d819d3e7ec884c5e193f3b7d58b994f1e42db5 |
| src/components/customer-survey.tsx | c1f4bb1f4253e6ade51d514fc2665fb0b9439806ed4752978c7eeac76ffc0fc8 |
| src/app/customer/surveys/page.tsx | 812e3b4181ac3c8692885dee5091aada69fafb8894d3b0e0412f450b1377785c |
| src/app/admin/surveys/page.tsx | 2cd4348f92104fb2d616363d15f4eab83115882e1f6ae44b84c3c8179e9b819b |
| microsoft_playwright_browser_gate/tests/customer-survey.spec.mjs | e8c2b5a70290b4d7a5035990e6055cda34e8ee24905915f6faa82d9576f94476 |
| cmd/electromobility-api/main.go | 65f4ab673ddb7f8a7ca95d22325ce4b13ad30302b28f898d6b9f167f7d1e625f |

## Reconstrucción final del lote

PASS:828/828archivos desde71packs reconstruidos en destino ausente.23fuentes
nuevas y activación coinciden exactamente con bytes probados;803fuentes previas
intactas. Sólo main y el lock de vínculo al perfil cambian entre las805previas.
next-env.d.ts conserva template canónico; sus2imports generados se documentan.
Receipt survey-reconstruction.json SHA-256 d467d821bdc20c9af0755eea6df0ab1eeeb949e38efe313ca93ca0843e275014.
G6 de este lote queda verificado; no elimina condiciones G7/G8 ni TEST02/03/07.

## Corrección final de atomicidad — FAIL802

Vigente: GO-CUSTOMER-SURVEY-API0.1.1, TS-CUSTOMER-SURVEY-PORTAL0.1.1 y
GO-HTTP-METRICS-REFERENCE0.1.15. Sólo los2SQL nuevos agregan BEGIN/COMMIT;
los22archivos de código/activación/doc restantes conservan bytes probados.
La misma migración provocada a fallar dejaba2tablas; corregida deja0. Up/down/up
válido también pasa. No cambio de negocio, esquema resultante ni dependencias.
Preflight269 verifica la base anterior completa; este delta se acredita con su
regresión SQL y reconstrucción exacta, sin relabelar ni repetir164gates intactos.
Receipt feedback-atomic1/result.json SHA-256 dea57a4e67d99dde18e4826ae0aa75d786d63eec3f451ee1062ab499d2eec88f.

## Cierre del lote V400 — checkpoint270

Preflight269:164pasos ejecutados PASS;56perfiles,167packs/1587archivos/
850Markdown y1330AUTHORED/150ADAPTED/107VERBATIM. Disponibilidad BLOCKED sólo
Docker ausente; el circuito PostgreSQL directo se probó con servidor real.
La corrección SQL posterior tiene su regresión red/green y reconstrucción final
828/828:825archivos idénticos a la reconstrucción inicial,2SQL transaccionales y
1lock de vínculo actualizado. No se repiten/relabelan164gates por esta corrección.
Franquicia70/820, referenciaHTTP71/828, web7/134.
Controles45passed; TEST02/03/07blocked. Cierre del circuito delimitado, no48/48.

Final reconstruction SHA-256 af12cbc03f19eab437e34fc791946a56ffb2ef6efd5a069154164dc2706a9f55.
Preflight269 receipt SHA-256 c6d00c6bf7dc7a959aade788bb7388db2418d1f80dc0896a3a138abdbaeb6fae.
