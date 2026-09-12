# Microsoft Business Central executable platform — V1

Fecha: 2026-08-26.

## Autoridad oficial

- Microsoft declara [BCApps](https://github.com/microsoft/BCApps) como repositorio único de System Application, Business Foundation, Base Application y aplicaciones first-party, bajo MIT.
- Release estable observada: [`releases/28.4/StrictMode`](https://github.com/microsoft/BCApps/releases/tag/releases/28.4/StrictMode), commit firmado `cb07eef12935e07dd4258c122d0e7b1920fc1223`, publicada 2026-08-07.
- La release contiene System Application, E-Document y External Storage Document Attachments, pero no contiene `src/Layers/W1/BaseApp`, `src/Apps/W1/APIV2` ni `src/Tests`. Por eso no se presenta como source release integral.
- El snapshot completo ya fijado de BCApps permanece `31a860b527f0dc72c7a44a255d7e7d403cfa4789`, 225.939.109 bytes, SHA-256 `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210`, clasificado como desarrollo condicionado.
- Microsoft publica [BcContainerHelper](https://github.com/microsoft/navcontainerhelper) para trabajar con Business Central Containers y documenta el [entorno de desarrollo basado en contenedores](https://learn.microsoft.com/dynamics365/business-central/dev-itpro/developer/devenv-running-container-development).

## Artefacto BcContainerHelper

PowerShell Gallery devolvió como versión estable vigente `6.1.16`, publicada `2026-08-21T10:56:10.073Z`:

```text
package bytes: 2886213
SHA-256: c20bae134a0cbd15927c062fd825234465ec8d189fe459688d282e11442bce10
Gallery SHA-512 base64: pGCpxgjcoUo8brljHY2ehsY0jyBUahhz4IXw97b9OVO5xE39Y6D74HQTWtCeDQoWrqewQJelN6f8Rzl/ouSfow==
LICENSE SHA-256: 9906940f61b1f0b533fa7d99baf55178b2808fbe113ea51dfbfad8572ccd5f2b
BcContainerHelper.psd1 SHA-256: 9667c6371b87d2c499607294a6c80dc18a20fd7e1964be17e2e5fdb327318892
BcContainerHelper.nuspec SHA-256: e3fddfd506394a9bd8373a8a369e6c596415cd316a672a2621cae3df921a049c
```

El paquete expandió 273 archivos, manifestó versión 6.1.16, importó correctamente en PowerShell 7.6.4 y exportó 309 comandos. Se verificó la presencia de `Get-BCArtifactUrl`, `New-BcContainer`, `Get-BcContainerAppInfo`, `Import-TestToolkitToBcContainer`, `Run-TestsInBcContainer` y `Run-AlPipeline`.

`Get-BCArtifactUrl -Type OnPrem -Country w1 -Version 28.4 -Select Latest` resolvió el artefacto oficial exacto observado:

```text
https://bcartifacts-exdbf9fwegejdqak.b02.azurefd.net/onprem/28.4.53241.0/w1
```

No se descargó ni ejecutó el artefacto Business Central. El host informó WinRM detenido, falta de permiso para modificar hosts y falta de autoridad para comandos Docker. No se elevaron privilegios ni se aceptó la EULA.

## Pack materializado

`MICROSOFT-BUSINESS-CENTRAL-EXECUTABLE-PLATFORM` 0.1.0 materializó seis archivos con hashes verificados. La adquisición offline contra el nupkg exacto pasó y conservó los hashes críticos.

```text
Materialized 6 files
BC_PLATFORM_CONFIG_VALID version=28.4.53241.0
BC_PLATFORM_RUNTIME_TEST_PASS positives=1 negatives=2 regression=1
BC_HELPER_ACQUIRED version=6.1.16
```

El runner sólo llama comandos oficiales Microsoft. La configuración, aprobación, hashes y receipt son orquestación `AUTHORED`; no contienen reglas empresariales ni se atribuyen a Microsoft.

## Fallos y condiciones

- `LIB-FAIL-096` y `LIB-FAIL-097`: patrón PowerShell `foreach |` falló y recurrió; el test del pack lo prohíbe estáticamente.
- `LIB-FAIL-098`: la plantilla vacía fue usada como fixture completo; el gate exigió identidad del proyecto.
- `LIB-FAIL-099`: una matriz no cruzó correctamente `pwsh -File`; el updater fue ejecutado en la misma sesión.
- El entorno local real permanece bloqueado por administrador, Windows containers/Docker, capacidad y aprobación del usuario.
- CRONUS/demo es desarrollo, no producción. Producción exige licencia Business Central y términos Microsoft.
- La localización argentina, fiscalidad, identidad, integraciones, datos, carga, seguridad y operación no se derivan automáticamente de W1.

## Admisión

`REBUILD_VERIFIED / CONDITIONED`. El agente puede componer y adquirir instantáneamente el módulo oficial sin Git y queda preparado para crear el entorno con aprobación. No existe evidencia de contenedor levantado ni de producción; afirmar lo contrario sería falso.
