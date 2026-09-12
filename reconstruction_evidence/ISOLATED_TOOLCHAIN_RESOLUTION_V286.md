# V286 — herramientas aisladas y diagnóstico corto

2026-09-07. Mantenimiento T2801/T2808; no ampliación de producto ni promoción.
Código modificado: `VERIFY_EXECUTABLE_LIBRARY.ps1` y su regresión
`markdown_system/test_toolchain_resolution.ps1`, ambos **AUTHORED**.

## Cambio ejecutable

- NodeExecutable, PnpmExecutable, DockerExecutable y PsqlExecutable complementan
  PythonExecutable, GoExecutable y DotNetExecutable. Resolución por defecto intacta;
  todos los consumers directos de esos parámetros los reciben también en Audit.
- Python explícito inválido ya no se sustituye silenciosamente por python3. El
  fallback sólo existe cuando no hubo selección explícita.
- `-Mode Toolchain` ejecuta sólo probes de versión. No materializa packs, instala,
  arranca PostgreSQL/Docker ni ejecuta el audit. Su receipt declara
  `library_audit_executed=false`, `availability_is_admission=false`, `steps=[]` y
  `TOOLS_AVAILABLE|TOOLS_MISSING`, nunca un PASS de biblioteca. `-RequireAll`
  rechaza faltantes después de conservar evidencia; no sobrescribe receipts.
- Audit/Preflight/Foundation conservan sus gates. Este diagnóstico no los reemplaza.
  Los parámetros seleccionan llamadas directas del runner; no prometen configurar
  toolchains transitivos de scripts externos ni sus subprocesses.

Uso portable (rutas reales seleccionadas por cada proyecto, no copiadas de este host):

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Toolchain -GoExecutable $goPath -DotNetExecutable $dotnetPath -PsqlExecutable $psqlPath -EvidencePath $newReceiptPath
# Después: Preflight y gates del perfil con los mismos parámetros aplicables.
# Resolver pins/SCA/licencias antes de adoptar; no modificar PATH global.
```

## Evidencia ejecutada

46 aserciones PASS sobre funciones/declaraciones reales y child processes, también
desde copia limpia de ocho archivos idénticos. Incluyen siete rutas independientes,
selección inválida sin fallback, directorio rechazado, versiones mínimas/malformadas,
proceso fallido, consumers Audit, modo corto real, RequireAll y receipt inmutable.
Los stubs de versión de los tests de routing están identificados; no cuentan como
ejecución de Go/.NET/PostgreSQL. Selectores reales de distribución: 32 checks PASS.

Probe separado con ejecutables reales: **1.496 ms**, una observación local, no p95
ni benchmark de construcción. Resultado `TOOLS_MISSING missing=1`:

| Herramienta | Versión observada | Disponibilidad |
|---|---|---|
| PowerShell | 7.6.5 | sí |
| Python | 3.14.4 | sí |
| Go | 1.26.7 | sí, aislado V281 |
| .NET SDK | 10.0.400 | sí, aislado V127 |
| Node | 24.14.1 | sí |
| pnpm | 11.25.0 | sí por mínimo, NO aprobado contra pin 11.19.0 |
| Docker | no resuelto | no; no afirmar ausencia física universal |
| PostgreSQL psql | 18.6 | sí, aislado |

Antes V285 encontró cuatro faltantes por resolución por defecto. No estaban
demostrados ausentes: tres existían fuera del PATH. El nuevo resultado no repite
el Preflight completo ni aprueba builds, servidor PG, Docker, restore o SCA.

Hashes SHA-256 observados de ejecutables, no de sus árboles completos:

- Go: `5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc`.
- dotnet: `ab1b71fd3dd71062e074c9fab8312081a81b7f2b3e0327c48c4d249c8d1a3135`.
- psql: `1e23b7f9ac7649b4717ada9b1f4e1b5b66c343937b6a91583f63528b96a61503`.

Artefactos diagnósticos locales en
`C:/Users/NL/AppData/Local/Temp/elite-v286-cebbc9e894bc412faad1bfb95413262f/`:

- `real-tools.json`: `9c0b3e96e3d5a8c8eaa38f8c461a64c585f31bc99391d03d8f96a1a944d2663b`;
- `clean-tests-final.log`: `1ce317d57e585594d93584a0f21cc3bcde024762679e441fc04a319c2cdc679c`;
- `clean-tests.log`: primer intento, 46 PASS pero exit heredado incorrecto.

Son evidencia local temporal; la distribución lleva este reporte y la suite
reproducible, no esos runtimes ni rutas como dependencia portable.

## Fallos conservados

397: rg con glob de path Windows error 123; corregido mediante directorio literal
y -g. La búsqueda encuentra expedientes Microsoft ARCA V123–V129; no inventar paths.
398: prueba roja `Missing explicit tool parameter: NodeExecutable`; corregido en
runner con parámetros uniformes y fallback explícito controlado.
399: argv de fixture concatenaba parámetro/path; llamada directa ya pasaba, corregido
fixture separando elementos y preservando salida del error en la aserción.
400: suite emitía PASS pero heredaba exit 1 de RequireAll esperado; ahora exit 0
explícito sólo después de todas las aserciones y cleanup. Caller sigue fail-closed.
Los cuatro permanecen en el ledger; copia limpia final y 46 checks pasan con exit 0.
401: lectura masiva de owners truncada; parser la rechaza antes de editar. Se
reemplaza por patches exactos y lecturas de líneas/hashes acotadas, sin tomar
datos truncados como evidencia. Esta lección diagnóstica no suma un test de producto.

## Fuente y límite del claim

Microsoft documenta resolución de comandos y uso de rutas completas para elegir
un ejecutable: [Get-Command](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/get-command?view=powershell-7.6)
y [Command precedence](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_command_precedence?view=powershell-7.6),
consultados 2026-09-07. Son autoridad de semántica, NO autoría Microsoft de este
runner ni certificación de su diseño. La disponibilidad no demuestra versión
exacta admitida, integridad transitiva, licencia/SCA, soporte de OS o productividad.

No se cambian 160 packs/1433 archivos, perfil 67/742 ni los fallos abiertos 384–387.
Readiness continúa BLOCKED/44: todavía faltan assurance y journeys release-bound,
recorrido integral de referencia, equivalencias, privacidad/observabilidad,
integraciones/documentos/IA y release/recuperación del roadmap T2801–T2810.
La verificación estructural final se ejecuta después de sellar el checkpoint; no
se incluye en las 46 aserciones ni se equipara a un Audit/Foundation completo.
