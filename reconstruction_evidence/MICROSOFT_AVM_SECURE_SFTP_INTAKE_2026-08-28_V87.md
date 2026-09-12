# Microsoft AVM Secure SFTP Intake — V87

Fecha: 2026-08-28  
Alcance: puerta SFTP privada y compilable basada en módulos oficiales Microsoft Azure Verified Modules, sin despliegue ni afirmaciones de persistencia automática.

## Fuentes oficiales fijadas

- repositorio [`Azure/bicep-registry-modules`](https://github.com/Azure/bicep-registry-modules/tree/32f1df57c75a48b5f0519083b7048f8fc00c33d8), commit firmado `32f1df57c75a48b5f0519083b7048f8fc00c33d8`, MIT;
- módulo AVM Storage Account `0.33.0` y módulo Local User `0.1.0` del mismo commit;
- prueba WAF E2E oficial del módulo Storage Account;
- [`Bicep 0.46.1`](https://github.com/Azure/bicep/releases/tag/v0.46.1), ejecutable Windows x64 de 117.685.696 bytes y SHA-256 `441d3d6094513acaa8a9be0b96b754d174bc0c19064c7be04e35b87e67b66250`;
- documentación Microsoft de [SFTP para Azure Blob Storage](https://learn.microsoft.com/en-us/azure/storage/blobs/secure-file-transfer-protocol-support-connect) y [almacenamiento inmutable](https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-storage-overview).

## Bytes Microsoft incorporados sin modificación

| Archivo | SHA-256 |
|---|---|
| `upstream/LICENSE` | `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383` |
| `upstream/local-user/main.bicep` | `efd933ceeb703f6d71b6a8f185ff41f38e0734014f8f35930c1b2afb9c059adc` |
| `upstream/waf-aligned/main.test.bicep` | `8f4e0c67787dfbb829824121cb813415cafd85d9ed2f15740c2b5ad9de4eab24` |

Cinco archivos adicionales son composición, configuración, lock y verificación Elite declarados `AUTHORED`; no se atribuyen a Microsoft.

## Contrato implementado

- StorageV2 `Standard_ZRS`, HNS y SFTP habilitados;
- autenticación del uploader exclusivamente por clave SSH; contraseña deshabilitada;
- Azure Storage Shared Key y acceso público a blobs deshabilitados;
- red pública deshabilitada, firewall `Deny`/bypass `None`, Private Endpoint Blob y DNS privado;
- permiso del uploader exactamente `cw` —crear/escribir—, sin listar, leer ni borrar;
- TLS 1.2, infrastructure encryption, soft delete y diagnósticos de lectura/escritura/borrado;
- cuenta retained-original separada obligatoria antes de extracción o persistencia.

La separación de retención no es una invención local: Microsoft documenta que blob versioning no es compatible con HNS y que immutable storage no es compatible mientras SFTP está habilitado. El landing SFTP recibe en cuarentena; un pipeline posterior debe demostrar cierre de transferencia, SHA-256, seguridad y copia exacta a la retención elegida.

## Reconstrucción y pruebas

- `MICROSOFT-AVM-SECURE-SFTP-INTAKE` 0.1.0 materializó 8/8 archivos;
- round-trip del pack: 8/8 hashes exactos;
- verifier estático: PASS;
- descarga de Bicep 0.46.1, validación de hash, restore AVM y compilación: `PASS_PINNED_0.46.1`;
- mutación `allowSshPassword: true`: rechazada;
- mutación de permisos a `rcwdl`: rechazada;
- `VERIFY_LIBRARY_PASS`: 61 packs, 580 archivos materializables y perfil SFTP 1/8;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 61 packs, 111 upstreams centralizados y un intake AVM SFTP materializado/verificado.

## Límite honesto

El verificador no crea recursos Azure. Un proyecto sigue bloqueado hasta aportar suscripción, región, red/subnet, DNS privado, Log Analytics, propietario y rotación de clave SSH, costos, identidades posteriores, Event Grid/orquestación durable, malware, retained-original, backup/restore y una transferencia SFTP DEV real con cierre y reconciliación. Este pack cierra código reutilizable para el landing; no convierte por sí solo un archivo recibido en dato empresarial autorizado.
