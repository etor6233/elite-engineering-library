# ARCA: referencia local y conexión de homologación

El artefacto contiene los dos workers Go, el worker SOAP/UDS .NET, un host de
fixtures explícito y el runtime Microsoft .NET/ASP.NET 10.0.11 fijado. Los clientes
WSAA/WSFE se generan desde los dos WSDL oficiales fijados; sus fuentes y receipt
quedan en `arca/generated-source`. No contiene certificados ni secretos de usuario.

## Ensayo local

Usar únicamente una base PostgreSQL de fixtures con las migraciones del mismo
artefacto. El API debe apuntar a esa misma base. No conectar el fixture fiscal a
datos reales. La identidad fiscal y las operaciones de punto de venta/factura
son las del contrato HTTP materializado y mantienen autorización por organización.

En tres terminales PowerShell abiertas en la raíz del artefacto:

```powershell
# Terminal 1: host de homologación simulada; genera su certificado efímero local.
$env:ELITE_ARCA_FIXTURE = '1'
$env:ARCA_WSFE_SOCKET = Join-Path $env:TEMP 'elite-arca-reference.sock'
& .\arca\dotnet\dotnet.exe .\arca\fixture\Elite.Arca.Wsfe.Fixture.dll $env:ARCA_WSFE_SOCKET
```

```powershell
# Terminal 2: configurar DATABASE_URL para la base local de fixtures.
$env:ARCA_WSFE_SOCKET = Join-Path $env:TEMP 'elite-arca-reference.sock'
$env:FISCAL_WORKER_ID = 'reference-fiscal-1'
& .\arca\go\arca-fiscal-worker.exe
```

```powershell
# Terminal 3: la misma DATABASE_URL y el mismo socket.
$env:ARCA_WSFE_SOCKET = Join-Path $env:TEMP 'elite-arca-reference.sock'
$env:FISCAL_PARAMETER_WORKER_ID = 'reference-parameters-1'
& .\arca\go\arca-parameter-worker.exe
```

El fixture reproduce una respuesta perdida después de autorizar; el worker Go
consulta el comprobante y reconcilia su estado durable sin volver a emitirlo.
También acepta la nota de crédito asociada y las siete consultas de parámetros.
Sus valores son datos sintéticos para el ensayo, no reglas fiscales aplicables.
Los procesos interactivos se detienen con Ctrl+C; la aceptación automatizada
cierra el fixture por su endpoint de control privado UDS y verifica PostgreSQL
después de un reinicio.

## Cuando el usuario aporte sus credenciales

La conexión de homologación usa el mismo socket y los mismos workers Go. Reemplazar
el host de fixtures por `arca/worker/Elite.Arca.Wsfe.Worker.dll` ejecutado con el
runtime incluido. Sus entradas son `ARCA_ENVIRONMENT=homologation`,
`ARCA_WSFE_SOCKET`, `ARCA_CERTIFICATE_THUMBPRINT` y
`ARCA_CERTIFICATE_STORE=CurrentUser|LocalMachine`; el certificado con clave privada
debe existir en el almacén indicado. El CUIT/punto de venta pertenece a la
configuración de la organización en PostgreSQL, nunca al código generado.

Estado de la conexión externa: `CONDITIONED_USER_CREDENTIALS`. El worker falla
cerrado si faltan sus entradas. Este artefacto no habilita facturación productiva;
las reglas fiscales, autorización del contribuyente y aceptación real se aplican
cuando se materialice el proyecto concreto.

## Reconstrucción

`ci/build_signed_reference.ps1` llama al builder local, que incluye
`ci/build_arca_reference.py`. Los inputs externos del build se fijan por SHA-256:
SDK/runtime, PowerShell, archivo original svcutil, nueve parches NuGet oficiales,
los dos WSDL y cinco paquetes de runtime. Cada build crea su propia caché NuGet
vacía y restaura sólo desde el directorio local comprobado. No ejecutar un
`dotnet tool restore` sobre el manifiesto histórico svcutil8.0.0: el generador
admitido es la composición adaptada definida en `svcutil-adaptation.lock.json`.
La adaptación es glue local declarado; no es una publicación oficial de Microsoft.

El release conserva licencias originales, fuentes correspondientes y el grafo
de dependencias. La generación offline no afirma haber consultado vulnerabilidades:
el gate de release ejecuta SCA sobre los manifiestos explícitos de esta revisión.
