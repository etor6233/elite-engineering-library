# ARCA WSAA Official Client Samples

## 1. Metadata

```yaml
pack_id: "ARCA-WSAA-OFFICIAL-CLIENT-SAMPLES"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa clientes oficiales ARCA/AFIP de WSAA en PowerShell y C# con normalización CRLF→LF declarada, preserva su acuerdo de redistribución y añade adquisición byte-exacta hash-locked."
stacks: ["Windows PowerShell legacy sample", ".NET Framework 2.0 legacy sample", "PowerShell 7 acquisition"]
compatible_with: []
incompatible_with: []
license_expression: "LicenseRef-ARCA-WSAA-Example-Terms AND LicenseRef-Workspace-Owner"
upstream_sources: ["https://arca.gob.ar/ws/documentacion/wsaa.asp", "https://arca.gob.ar/ws/WSAA/ejemplos/dev-wsaa-cliente-powershell.zip", "https://arca.gob.ar/ws/WSAA/ejemplos/dev-wsaa-cliente-dotnet-cs.zip"]
verified_at: "2026-08-30"
```

Cinco archivos son `ADAPTED` únicamente por la normalización reversible CRLF→LF que exige el contenedor Markdown; su contenido procede de los ZIP oficiales ARCA/AFIP y los scripts de adquisición restauran y verifican los bytes originales. Sus README autorizan redistribuir, publicar o descargar el código total o parcialmente sin autorización y exigen evaluación/corrección técnica. Dos scripts son `AUTHORED` y sólo adquieren/verifican. El pack no afirma que samples legacy sean seguros o productivos: PowerShell usa `New-WebServiceProxy`, defaults locales, OpenSSL externo y archivos con ticket; C# requiere .NET Framework 2.0 y acepta password por CLI. Son autoridad protocolar y código oficial reutilizable, no adapter moderno admitido.

## 2. Applicability

Use to inspect or reconstruct the official WSAA signing/login flow and its exact redistribution terms. Use the acquisition script for reproducible source intake. Do not wire either sample directly into a production service. A modern adapter must isolate private keys, suppress token/sign output, pin endpoints/TLS, cache before expiry, singleflight refresh, redact evidence and pass ARCA homologation with project certificates.

## 3. Architecture contract

Materialized text is LF-normalized and separately attributed. Archive URL, byte length, archive SHA-256 and original selected-file SHA-256 fail closed in acquisition; materialized LF SHA-256 values fail closed in regression. Acquisition supports an offline cache, extracts into a unique temporary directory, rejects path escape, checks the official redistribution sentence and emits a non-secret source receipt. Known upstream conditions prevent a future update from silently promoting legacy code.

## 4. Exact file manifest

```text
CREATE official/arca/wsaa/powershell/README.md
CREATE official/arca/wsaa/powershell/wsaa-cliente.ps1
CREATE official/arca/wsaa/csharp/README.md
CREATE official/arca/wsaa/csharp/ClienteLoginCms.cs
CREATE official/arca/wsaa/csharp/app.config
CREATE tools/acquire-arca-wsaa-samples.ps1
CREATE tools/test-arca-wsaa-samples.ps1
```

## 5. Materialization blocks

### FILE: `official/arca/wsaa/powershell/README.md`
```yaml
block_id: "ARCA-WSAA:ps-readme:v1"
operation: CREATE
provenance: ADAPTED
source: "ARCA official archive dev-wsaa-cliente-powershell.zip; CRLF normalized to LF"
license: "LicenseRef-ARCA-WSAA-Example-Terms"
sha256: "f76dae5edf39e2602e0f414eb8fb675f7ed0b1db15cbfb3ad44b679ed061d705"
variables: []
secrets_allowed: false
```
````markdown
EJEMPLO DE APLICACION CLIENTE DEL WSAA
======================================

Ejemplo de cliente del WSAA (webservice de autenticacion y autorizacion). 
Consume el metodo LoginCms ejecutando desde la Powershell de Windows. 
Muestra en stdout el login ticket response. 


REQUERIMIENTOS
--------------

- Microsoft Powershell
- OpenSSL


USO
---

```
   wsaa-cliente.ps1 
      [ -Certificado unCertificado ] 
      [ -ClavePrivada unaClavePrivada ]
      [ -ServicioId unIdServicio ]
      [ -OutXml unArchivoXml ]
      [ -OutCms unArchivoXmlFirmado ]
      [ -WsaaWsdl unaURL ]
``` 

   
EJEMPLOS
--------

```
   .\wsaa-cliente.ps1 -Certificado unCert.crt -ClavePrivada unaPriv.key -ServicioId wsfe
```


ACUERDO DE USO
--------------

1. El Departamento de Soporte Tecnico de la AFIP (DeSoTe/AFIP), pone a disposicion
el siguiente codigo para su utilizacion con el WebService de Autenticacion y Autorizacion (WSAA)
de la AFIP.

2. El mismo puede ser re-distribuido, publicado o descargado en forma total o parcial, ya sea
en forma electronica, mecanica u optica, sin requerir la autorizacion de DeSoTe/AFIP. 

3. DeSoTe/AFIP no asume ninguna responsabilidad de los errores que pueda contener el codigo ni la
obligacion de subsanar dichos errores o informar de la existencia de los mismos.

4. DeSoTe/AFIP no asume ninguna responsabilidad que surja de la utilizacion del codigo, ya sea por
utilizacion ilegal de patentes, perdida de beneficios, perdida de informacion o cualquier otro
inconveniente.

5. Bajo ninguna circunstancia DeSoTe/AFIP podra ser indicada como responsable por consecuencias y/o
incidentes ya sean directos o indirectos que puedan surgir de la utilizacion del codigo.

6. DeSoTe/AFIP no da ninguna garantia, expresa o implicita, de la utilidad del codigo, si el mismo es
correcto, o si cumple con los requerimientos de algun proposito en particular.

7. DeSoTe/AFIP puede realizar cambios en cualquier momento en el codigo sin previo aviso.

8. El codigo debera ser evaluado, verificado, corregido y/o adaptado por personal tecnico calificado
de las entidades que lo utilicen.

EL CODIGO FUENTE ES DISTRIBUIDO PARA EVALUACION, CON TODOS SUS ERRORES Y OMISIONES. LA
RESPONSABILIDAD DEL CORRECTO FUNCIONAMIENTO DEL MISMO YA SEA POR SI SOLO O COMO PARTE DE
OTRA APLICACION, QUEDA A CARGO DE LAS ENTIDADES QUE LO UTILICEN. LA UTILIZACION DEL CODIGO
SIGNIFICA LA ACEPTACION DE TODOS LOS TERMINOS Y CONDICIONES MENCIONADAS ANTERIORMENTE.
````

### FILE: `official/arca/wsaa/powershell/wsaa-cliente.ps1`
```yaml
block_id: "ARCA-WSAA:ps-client:v1"
operation: CREATE
provenance: ADAPTED
source: "ARCA official archive dev-wsaa-cliente-powershell.zip; CRLF normalized to LF"
license: "LicenseRef-ARCA-WSAA-Example-Terms"
sha256: "b7ca06a2fa210f9998f00e3a32378a79a48ae9b53c12aa8960752d9c199ddc2d"
variables: []
secrets_allowed: true
```
````powershell
#
# Ejemplo de cliente del WSAA (webservice de autenticacion y autorizacion). 
# Consume el metodo LoginCms ejecutando desde la Powershell de Windows. 
# Muestra en stdout el login ticket response.
#
# REQUISITOS: openssl
#
# Parametros de linea de comandos:
#
#   $Certificado: Archivo del certificado firmante a usar
#   $ClavePrivada: Archivo de clave privada a usar
#   $ServicioId: ID de servicio a acceder
#   $OutXml: Archivo TRA a crear
#   $OutCms: Archivo CMS a crear
#   $WsaaWsdl: URL del WSDL del WSAA
#
[CmdletBinding()]
Param(
   [Parameter(Mandatory=$False)]
   [string]$Certificado="myCert.crt",
	
   [Parameter(Mandatory=$False)]
   [string]$ClavePrivada="myPrivate.key",
   
   [Parameter(Mandatory=$False)]
   [string]$ServicioId="wsfe",
   
   [Parameter(Mandatory=$False)]
   [string]$OutXml="LoginTicketRequest.xml",   
   
   [Parameter(Mandatory=$False)]
   [string]$OutCms="LoginTicketRequest.xml.cms",   

   [Parameter(Mandatory=$False)]
   [string]$WsaaWsdl = "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL"    
)

$ErrorActionPreference = "Stop"

# PASO 1: ARMAR EL XML DEL TICKET DE ACCESO
$dtNow = Get-Date 
$xmlTA = New-Object System.XML.XMLDocument
$xmlTA.LoadXml('<loginTicketRequest><header><uniqueId></uniqueId><generationTime></generationTime><expirationTime></expirationTime></header><service></service></loginTicketRequest>')
$xmlUniqueId = $xmlTA.SelectSingleNode("//uniqueId")
$xmlGenTime = $xmlTA.SelectSingleNode("//generationTime")
$xmlExpTime = $xmlTA.SelectSingleNode("//expirationTime")
$xmlService = $xmlTA.SelectSingleNode("//service")
$xmlGenTime.InnerText = $dtNow.AddMinutes(-10).ToString("s")
$xmlExpTime.InnerText = $dtNow.AddMinutes(+10).ToString("s")
$xmlUniqueId.InnerText = $dtNow.ToString("yyMMddHHMM")
$xmlService.InnerText = $ServicioId
$seqNr = Get-Date -UFormat "%Y%m%d%H%S"
$xmlTA.InnerXml | Out-File $seqNr-$OutXml -Encoding ASCII

# PASO 2: FIRMAR CMS
openssl cms -sign -in $seqNr-$OutXml -signer $Certificado -inkey $ClavePrivada -nodetach -outform der -out $seqNr-$OutCms-DER

# PASO 3: ENCODEAR EL CMS EN BASE 64
openssl base64 -in $seqNr-$OutCms-DER -e -out $seqNr-$OutCms-DER-b64

# PASO 3: INVOCAR AL WSAA
try
{
   $cms = Get-Content $seqNr-$OutCms-DER-b64 -Raw
   $wsaa = New-WebServiceProxy -Uri $WsaaWsdl -ErrorAction Stop
   $wsaaResponse = $wsaa.loginCms($cms) 
   $wsaaResponse > $seqNr-loginTicketResponse.xml 
   $wsaaResponse
}
catch
{
   $errMsg = $_.Exception.Message
   $errMsg > $seqNr-loginTicketResponse-ERROR.xml 
   $errMsg
}
````

### FILE: `official/arca/wsaa/csharp/README.md`
```yaml
block_id: "ARCA-WSAA:cs-readme:v1"
operation: CREATE
provenance: ADAPTED
source: "ARCA official archive dev-wsaa-cliente-dotnet-cs.zip; CRLF normalized to LF"
license: "LicenseRef-ARCA-WSAA-Example-Terms"
sha256: "07f8d3930ae8bfee0661a5abc455ee1ded4ccac9f89cf46610429467275c01d5"
variables: []
secrets_allowed: false
```
````markdown
EJEMPLO DE APLICACION CLIENTE DEL WSAA
======================================

Ejemplo de cliente del WSAA (webservice de autenticacion y autorizacion). 
Consume el metodo LoginCms ejecutando desde linea de comandos de Windows. 
Muestra en stdout el login ticket response. 


REQUERIMIENTOS DE PLATAFORMA
----------------------------

- .NET Framework 2.0 Redistributable


USO
---

```
   clientelogincms_cs [ opciones ] ...
``` 

opciones:

- -s servicio      ID del servicio de negocio
- -w url           URL del WSDL del WSAA
- -c certif        Ruta del certificado (con clave privada)
- -p certifpwd     Password del certificado (con clave privada)
- -x IP:port       IP:port del proxy
- -y proxyusr      Usuario del proxy
- -z proxypwd      Password del proxy
- -v on|off        Salida detallada
- -?               Muestra ayuda de uso
   
EJEMPLOS
--------

```
   clientelogincms_cs -s wsfe -v on
```

```   
   clientelogincms_cs    
     -s "wsfe"
     -c "f:\wsaa_test\micertificado.p12"
     -p "pass_p12"
     -w "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL"
     -x "http://10.20.152.112:80"
     -y "user_proxy"
     -z "pass_proxy"
     -v on 

```

	 
NOTAS
-----

El certificado usado para firmar (especificado en la opcion -c) debe incluir la clave privada.
La password del certificado se especifica en la opcion -p
El certificado con clave privada en formato PKCS12 se puede generar de esta forma:

```
   openssl pkcs12 -export -in MiCertificado.crt -inkey claveprivada -out cert.p12
```

ACUERDO DE USO
--------------

1. El Departamento de Soporte Tecnico de la AFIP (DeSoTe/AFIP), pone a disposicion
el siguiente codigo para su utilizacion con el WebService de Autenticacion y Autorizacion (WSAA)
de la AFIP.

2. El mismo puede ser re-distribuido, publicado o descargado en forma total o parcial, ya sea
en forma electronica, mecanica u optica, sin requerir la autorizacion de DeSoTe/AFIP. 

3. DeSoTe/AFIP no asume ninguna responsabilidad de los errores que pueda contener el codigo ni la
obligacion de subsanar dichos errores o informar de la existencia de los mismos.

4. DeSoTe/AFIP no asume ninguna responsabilidad que surja de la utilizacion del codigo, ya sea por
utilizacion ilegal de patentes, perdida de beneficios, perdida de informacion o cualquier otro
inconveniente.

5. Bajo ninguna circunstancia DeSoTe/AFIP podra ser indicada como responsable por consecuencias y/o
incidentes ya sean directos o indirectos que puedan surgir de la utilizacion del codigo.

6. DeSoTe/AFIP no da ninguna garantia, expresa o implicita, de la utilidad del codigo, si el mismo es
correcto, o si cumple con los requerimientos de algun proposito en particular.

7. DeSoTe/AFIP puede realizar cambios en cualquier momento en el codigo sin previo aviso.

8. El codigo debera ser evaluado, verificado, corregido y/o adaptado por personal tecnico calificado
de las entidades que lo utilicen.

EL CODIGO FUENTE ES DISTRIBUIDO PARA EVALUACION, CON TODOS SUS ERRORES Y OMISIONES. LA
RESPONSABILIDAD DEL CORRECTO FUNCIONAMIENTO DEL MISMO YA SEA POR SI SOLO O COMO PARTE DE
OTRA APLICACION, QUEDA A CARGO DE LAS ENTIDADES QUE LO UTILICEN. LA UTILIZACION DEL CODIGO
SIGNIFICA LA ACEPTACION DE TODOS LOS TERMINOS Y CONDICIONES MENCIONADAS ANTERIORMENTE.
````

### FILE: `official/arca/wsaa/csharp/ClienteLoginCms.cs`
```yaml
block_id: "ARCA-WSAA:cs-client:v1"
operation: CREATE
provenance: ADAPTED
source: "ARCA official archive dev-wsaa-cliente-dotnet-cs.zip; CRLF normalized to LF"
license: "LicenseRef-ARCA-WSAA-Example-Terms"
sha256: "58694ce2cae209ad8386163396d336433961c091240cf1cee5c61f55e2746781"
variables: []
secrets_allowed: true
```
````csharp
using System;
using System.Collections;
using System.Collections.Generic;
using System.Data;
using System.Diagnostics;
using System.Text;
using System.Xml;
using System.Net;
using System.Security;
using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;
using System.IO;
using System.Runtime.InteropServices;

/// <summary>
/// Clase para crear objetos Login Tickets
/// </summary>
/// <remarks>
/// Ver documentacion: 
///    Especificacion Tecnica del Webservice de Autenticacion y Autorizacion
///    Version 1.0
///    Departamento de Seguridad Informatica - AFIP
/// </remarks>
class LoginTicket
{ 
    public UInt32 UniqueId; // Entero de 32 bits sin signo que identifica el requerimiento
    public DateTime GenerationTime; // Momento en que fue generado el requerimiento
    public DateTime ExpirationTime; // Momento en el que expira la solicitud
    public string Service; // Identificacion del WSN para el cual se solicita el TA
    public string Sign; // Firma de seguridad recibida en la respuesta
    public string Token; // Token de seguridad recibido en la respuesta
    public XmlDocument XmlLoginTicketRequest = null;
    public XmlDocument XmlLoginTicketResponse = null;
    public string RutaDelCertificadoFirmante;
    public string XmlStrLoginTicketRequestTemplate = "<loginTicketRequest><header><uniqueId></uniqueId><generationTime></generationTime><expirationTime></expirationTime></header><service></service></loginTicketRequest>";
    private bool _verboseMode = true;
    private static UInt32 _globalUniqueID = 0; // OJO! NO ES THREAD-SAFE

    /// <summary>
    /// Construye un Login Ticket obtenido del WSAA
    /// </summary>
    /// <param name="argServicio">Servicio al que se desea acceder</param>
    /// <param name="argUrlWsaa">URL del WSAA</param>
    /// <param name="argRutaCertX509Firmante">Ruta del certificado X509 (con clave privada) usado para firmar</param>
    /// <param name="argPassword">Password del certificado X509 (con clave privada) usado para firmar</param>
    /// <param name="argProxy">IP:port del proxy</param>
    /// <param name="argProxyUser">Usuario del proxy</param>''' 
    /// <param name="argProxyPassword">Password del proxy</param>
    /// <param name="argVerbose">Nivel detallado de descripcion? true/false</param>
    /// <remarks></remarks>
    public string ObtenerLoginTicketResponse(string argServicio, string argUrlWsaa, string argRutaCertX509Firmante, SecureString argPassword, string argProxy, string argProxyUser, string argProxyPassword, bool argVerbose)
    {
        const string ID_FNC = "[ObtenerLoginTicketResponse]";
        this.RutaDelCertificadoFirmante = argRutaCertX509Firmante;
        this._verboseMode = argVerbose;
        CertificadosX509Lib.VerboseMode = argVerbose;
        string cmsFirmadoBase64 = null;
        string loginTicketResponse = null;
        XmlNode xmlNodoUniqueId = default(XmlNode);
        XmlNode xmlNodoGenerationTime = default(XmlNode);
        XmlNode xmlNodoExpirationTime = default(XmlNode);
        XmlNode xmlNodoService = default(XmlNode);

        // PASO 1: Genero el Login Ticket Request
        try
        {
            _globalUniqueID += 1;

            XmlLoginTicketRequest = new XmlDocument();
            XmlLoginTicketRequest.LoadXml(XmlStrLoginTicketRequestTemplate);

            xmlNodoUniqueId = XmlLoginTicketRequest.SelectSingleNode("//uniqueId");
            xmlNodoGenerationTime = XmlLoginTicketRequest.SelectSingleNode("//generationTime");
            xmlNodoExpirationTime = XmlLoginTicketRequest.SelectSingleNode("//expirationTime");
            xmlNodoService = XmlLoginTicketRequest.SelectSingleNode("//service");
            xmlNodoGenerationTime.InnerText = DateTime.Now.AddMinutes(-10).ToString("s");
            xmlNodoExpirationTime.InnerText = DateTime.Now.AddMinutes(+10).ToString("s");
            xmlNodoUniqueId.InnerText = Convert.ToString(_globalUniqueID);
            xmlNodoService.InnerText = argServicio;
            this.Service = argServicio;

            if (this._verboseMode) Console.WriteLine(XmlLoginTicketRequest.OuterXml);
        }
        catch (Exception excepcionAlGenerarLoginTicketRequest) 
        {
            throw new Exception(ID_FNC + "***Error GENERANDO el LoginTicketRequest : " + excepcionAlGenerarLoginTicketRequest.Message + excepcionAlGenerarLoginTicketRequest.StackTrace);
        }

        // PASO 2: Firmo el Login Ticket Request
        try
        {
            if (this._verboseMode) Console.WriteLine(ID_FNC + "***Leyendo certificado: {0}", RutaDelCertificadoFirmante);

            X509Certificate2 certFirmante = CertificadosX509Lib.ObtieneCertificadoDesdeArchivo(RutaDelCertificadoFirmante, argPassword);

            if (this._verboseMode)
            {
                Console.WriteLine(ID_FNC + "***Firmando: ");
                Console.WriteLine(XmlLoginTicketRequest.OuterXml);
            }

            // Convierto el Login Ticket Request a bytes, firmo el msg y lo convierto a Base64
            Encoding EncodedMsg = Encoding.UTF8;
            byte[] msgBytes = EncodedMsg.GetBytes(XmlLoginTicketRequest.OuterXml);
            byte[] encodedSignedCms = CertificadosX509Lib.FirmaBytesMensaje(msgBytes, certFirmante);
            cmsFirmadoBase64 = Convert.ToBase64String(encodedSignedCms);
        }
        catch (Exception excepcionAlFirmar)
        {
            throw new Exception(ID_FNC + "***Error FIRMANDO el LoginTicketRequest : " + excepcionAlFirmar.Message);
        }

        // PASO 3: Invoco al WSAA para obtener el Login Ticket Response
        try
        {
            if (this._verboseMode)
            {
                Console.WriteLine(ID_FNC + "***Llamando al WSAA en URL: {0}", argUrlWsaa);
                Console.WriteLine(ID_FNC + "***Argumento en el request:");
                Console.WriteLine(cmsFirmadoBase64);
            }

            ClienteLoginCms_CS.Wsaa.LoginCMSService servicioWsaa = new ClienteLoginCms_CS.Wsaa.LoginCMSService();
            servicioWsaa.Url = argUrlWsaa;

            // Veo si hay que salir a traves de un proxy
            if (argProxy != null)
            {
                servicioWsaa.Proxy = new WebProxy(argProxy, true);
                if (argProxyUser != null)
                {
                    NetworkCredential Credentials = new NetworkCredential(argProxyUser, argProxyPassword);
                    servicioWsaa.Proxy.Credentials = Credentials;
                }
            }

            loginTicketResponse = servicioWsaa.loginCms(cmsFirmadoBase64);

            if (this._verboseMode)
            {
                Console.WriteLine(ID_FNC + "***LoguinTicketResponse: ");
                Console.WriteLine(loginTicketResponse);
            }

        }
        catch (Exception excepcionAlInvocarWsaa)
        {
            throw new Exception(ID_FNC + "***Error INVOCANDO al servicio WSAA : " + excepcionAlInvocarWsaa.Message);
        }

        // PASO 4: Analizo el Login Ticket Response recibido del WSAA
        try
        {
            XmlLoginTicketResponse = new XmlDocument();
            XmlLoginTicketResponse.LoadXml(loginTicketResponse);

            this.UniqueId = UInt32.Parse(XmlLoginTicketResponse.SelectSingleNode("//uniqueId").InnerText);
            this.GenerationTime = DateTime.Parse(XmlLoginTicketResponse.SelectSingleNode("//generationTime").InnerText);
            this.ExpirationTime = DateTime.Parse(XmlLoginTicketResponse.SelectSingleNode("//expirationTime").InnerText);
            this.Sign = XmlLoginTicketResponse.SelectSingleNode("//sign").InnerText;
            this.Token = XmlLoginTicketResponse.SelectSingleNode("//token").InnerText;
        }
        catch (Exception excepcionAlAnalizarLoginTicketResponse)
        {
            throw new Exception(ID_FNC + "***Error ANALIZANDO el LoginTicketResponse : " + excepcionAlAnalizarLoginTicketResponse.Message);
        }
        return loginTicketResponse;
    }
}

/// <summary>
/// Libreria de utilidades para manejo de certificados
/// </summary>
/// <remarks></remarks>
class CertificadosX509Lib
{
    public static bool VerboseMode = false;

    /// <summary>
    /// Firma mensaje
    /// </summary>
    /// <param name="argBytesMsg">Bytes del mensaje</param>
    /// <param name="argCertFirmante">Certificado usado para firmar</param>
    /// <returns>Bytes del mensaje firmado</returns>
    /// <remarks></remarks>
    public static byte[] FirmaBytesMensaje(byte[] argBytesMsg, X509Certificate2 argCertFirmante)
    {
        const string ID_FNC = "[FirmaBytesMensaje]";
        try
        {
            // Pongo el mensaje en un objeto ContentInfo (requerido para construir el obj SignedCms)
            ContentInfo infoContenido = new ContentInfo(argBytesMsg);
            SignedCms cmsFirmado = new SignedCms(infoContenido);

            // Creo objeto CmsSigner que tiene las caracteristicas del firmante
            CmsSigner cmsFirmante = new CmsSigner(argCertFirmante);
            cmsFirmante.IncludeOption = X509IncludeOption.EndCertOnly;

            if (VerboseMode) Console.WriteLine(ID_FNC + "***Firmando bytes del mensaje...");

            // Firmo el mensaje PKCS #7
            cmsFirmado.ComputeSignature(cmsFirmante);

            if (VerboseMode) Console.WriteLine(ID_FNC + "***OK mensaje firmado");

            // Encodeo el mensaje PKCS #7.
            return cmsFirmado.Encode();
        }
        catch (Exception excepcionAlFirmar)
        {
            throw new Exception(ID_FNC + "***Error al firmar: " + excepcionAlFirmar.Message);
        }
    }

    /// <summary>
    /// Lee certificado de disco
    /// </summary>
    /// <param name="argArchivo">Ruta del certificado a leer.</param>
    /// <returns>Un objeto certificado X509</returns>
    /// <remarks></remarks>
    public static X509Certificate2 ObtieneCertificadoDesdeArchivo(string argArchivo, SecureString argPassword)
    {
        const string ID_FNC = "[ObtieneCertificadoDesdeArchivo]";
        X509Certificate2 objCert = new X509Certificate2();
        try
        {
            if (argPassword.IsReadOnly())
            {
                objCert.Import(File.ReadAllBytes(argArchivo), argPassword, X509KeyStorageFlags.PersistKeySet);
            }
            else
            {
                objCert.Import(File.ReadAllBytes(argArchivo));
            }
            return objCert;
        }
        catch (Exception excepcionAlImportarCertificado)
        {
            throw new Exception(ID_FNC + "***Error al leer certificado: " + excepcionAlImportarCertificado.Message);
        }
    }
}

/// <summary>
/// Clase principal
/// </summary>
/// <remarks></remarks>
class ProgramaPrincipal
{
    // Valores por defecto, globales en esta clase
    const string DEFAULT_URLWSAAWSDL = "https://wsaahomo.afip.gov.ar/ws/services/LoginCms?WSDL";
    const string DEFAULT_SERVICIO = "wsfe";
    const string DEFAULT_CERTSIGNER = "c:\\MiCertificadoConClavePrivada.pfx";
    const string DEFAULT_PROXY = null;
    const string DEFAULT_PROXY_USER = null;
    const string DEFAULT_PROXY_PASSWORD = null;
    const bool DEFAULT_VERBOSE = true;

    /// <summary>
    /// Funcion Main (consola)
    /// </summary>
    /// <param name="args">Argumentos de linea de comandos</param>
    /// <returns>0 si termin� bien, valores negativos si hubieron errores</returns>
    /// <remarks></remarks>
    public static int Main(string[] args)
    {
        const string ID_FNC = "[Main]";

        string strUrlWsaaWsdl = DEFAULT_URLWSAAWSDL;
        string strIdServicioNegocio = DEFAULT_SERVICIO;
        string strRutaCertSigner = DEFAULT_CERTSIGNER;
        SecureString strPasswordSecureString = new SecureString();
        string strProxy = DEFAULT_PROXY;
        string strProxyUser = DEFAULT_PROXY_USER;
        string strProxyPassword = DEFAULT_PROXY_PASSWORD;
        bool blnVerboseMode = DEFAULT_VERBOSE;

        MostrarVersion();

        if (args.Length == 0)
        {
            ExplicarUso();
            return -1;
        }

        // Analizo argumentos de linea de comandos
        for (int i = 0; i <= args.Length - 1; i++)
        {
            string argumento;
            argumento = args[i];

            if (String.Compare(argumento, "-w", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� la URL del WSDL del WSAA");
                    return -1;
                }
                else
                {
                    strUrlWsaaWsdl = args[i + 1];
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-s", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� el ID del servicio de negocio");
                    return -1;
                }
                else
                {
                    strIdServicioNegocio = args[i + 1];
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-c", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� ruta del certificado firmante");
                    return -1;
                }
                else
                {
                    strRutaCertSigner = args[i + 1];
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-p", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� password del certificado firmante");
                    return -1;
                }
                else
                {
                    foreach (char c in args[i+1]) strPasswordSecureString.AppendChar(c);
                    strPasswordSecureString.MakeReadOnly();
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-x", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� IP:port del proxy");
                    return -1;
                }
                else
                {
                    strProxy = args[i + 1];
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-y", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� usuario del proxy");
                    return -1;
                }
                else
                {
                    strProxyUser = args[i + 1];
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-z", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� password del proxy");
                    return -1;
                }
                else
                {
                    strProxyPassword = args[i + 1];
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-v", true) == 0)
            {
                if (args.Length < (i + 2))
                {
                    Console.WriteLine("Error: no se especific� modo: on|off");
                    return -1;
                }
                else
                {
                    blnVerboseMode = (String.Compare(args[i + 1], "on", true) == 0 ? true : false);
                    i = i + 1;
                }
            }

            else if (String.Compare(argumento, "-?", true) == 0)
            {
                ExplicarUso();
                return 0;
            }

            else
            {
                System.Reflection.Assembly assembly = System.Reflection.Assembly.GetExecutingAssembly();
                FileVersionInfo fvi = FileVersionInfo.GetVersionInfo(assembly.Location);

                Console.WriteLine("Error: argumento desconocido: {0}", argumento);
                Console.WriteLine("Para obtener ayuda: {0} -?", fvi.ProductName);
                return -2;
            }
        }

        // Argumentos OK, entonces procesar normalmente...

        LoginTicket objTicketRespuesta = null;
        string strTicketRespuesta = null;

        try
        {
            if (blnVerboseMode)
            {
                Console.WriteLine(ID_FNC + "***Servicio a acceder: {0}", strIdServicioNegocio);
                Console.WriteLine(ID_FNC + "***URL del WSAA: {0}", strUrlWsaaWsdl);
                Console.WriteLine(ID_FNC + "***Ruta del certificado: {0}", strRutaCertSigner);
                Console.WriteLine(ID_FNC + "***Modo verbose: {0}", blnVerboseMode);
            }
            objTicketRespuesta = new LoginTicket();
            if (blnVerboseMode) Console.WriteLine(ID_FNC + "***Accediendo a {0}", strUrlWsaaWsdl);
            strTicketRespuesta = objTicketRespuesta.ObtenerLoginTicketResponse(strIdServicioNegocio, strUrlWsaaWsdl, strRutaCertSigner, strPasswordSecureString, strProxy, strProxyUser, strProxyPassword, blnVerboseMode);
        }
        catch (Exception excepcionAlObtenerTicket)
        {
            Console.WriteLine(ID_FNC + "***EXCEPCION AL OBTENER TICKET: " + excepcionAlObtenerTicket.Message);
            return -10;
        }
        return 0;
    }

    /// <summary>
    /// Explica el uso del comando
    /// </summary>
    /// <remarks></remarks>
    public static void ExplicarUso()
    {
        System.Reflection.Assembly assembly = System.Reflection.Assembly.GetExecutingAssembly();
        FileVersionInfo fvi = FileVersionInfo.GetVersionInfo(assembly.Location);

        Console.WriteLine("");
        Console.WriteLine("Uso: {0} [opciones]...", fvi.ProductName);
        Console.WriteLine("");
        Console.WriteLine("opciones:");
        Console.WriteLine("");
        Console.WriteLine("  -s servicio      ID del servicio de negocio");
        Console.WriteLine("                   Valor por defecto: " + DEFAULT_SERVICIO);
        Console.WriteLine("");
        Console.WriteLine("  -c certif        Ruta del certificado (con clave privada)");
        Console.WriteLine("                   Valor por defecto: " + DEFAULT_CERTSIGNER);
        Console.WriteLine("");
        Console.WriteLine("  -p certifpwd     Password del certificado (con clave privada)");
        Console.WriteLine("                   Valor por defecto: sin password");
        Console.WriteLine("");
        Console.WriteLine("  -x IP:port       IP:port del proxy");
        Console.WriteLine("                   Valor por defecto: sin proxy");
        Console.WriteLine("");
        Console.WriteLine("  -y proxyusr      Usuario del proxy");
        Console.WriteLine("                   Valor por defecto: sin usuario proxy");
        Console.WriteLine("");
        Console.WriteLine("  -z proxypwd      Password del proxy");
        Console.WriteLine("                   Valor por defecto: sin password proxy");
        Console.WriteLine("");
        Console.WriteLine("  -w url           URL del WSDL del WSAA");
        Console.WriteLine("                   Valor por defecto: " + DEFAULT_URLWSAAWSDL);
        Console.WriteLine("");
        Console.WriteLine("  -v on|off        Reportes detallados de la ejecuci�n");
        Console.WriteLine("");
        Console.WriteLine("  -?               Esta ayuda");
    }

    public static void MostrarVersion()
    {
        System.Reflection.Assembly assembly = System.Reflection.Assembly.GetExecutingAssembly();
        FileVersionInfo fvi = FileVersionInfo.GetVersionInfo(assembly.Location);

        Console.WriteLine("Aplicacion: {0}", fvi.ProductName);
        Console.WriteLine("Version   : {0}", fvi.FileVersion);
    }
}
````

### FILE: `official/arca/wsaa/csharp/app.config`
```yaml
block_id: "ARCA-WSAA:cs-config:v1"
operation: CREATE
provenance: ADAPTED
source: "ARCA official archive dev-wsaa-cliente-dotnet-cs.zip; CRLF normalized to LF"
license: "LicenseRef-ARCA-WSAA-Example-Terms"
sha256: "494ad889f15c76aaeffcc07717f59ac2371f97bba0d10fec61ea2bb8c1aecc19"
variables: []
secrets_allowed: false
```
````xml
<?xml version="1.0" encoding="utf-8" ?>
<configuration>
    <configSections>
        <sectionGroup name="applicationSettings" type="System.Configuration.ApplicationSettingsGroup, System, Version=2.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" >
            <section name="ClienteLoginCms_CS.Properties.Settings" type="System.Configuration.ClientSettingsSection, System, Version=2.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" requirePermission="false" />
        </sectionGroup>
    </configSections>
    <applicationSettings>
        <ClienteLoginCms_CS.Properties.Settings>
            <setting name="ClienteLoginCms_CS_Wsaa_LoginCMSService" serializeAs="String">
                <value>https://wsaahomo.afip.gov.ar/ws/services/LoginCms</value>
            </setting>
        </ClienteLoginCms_CS.Properties.Settings>
    </applicationSettings>
</configuration>
````

### FILE: `tools/acquire-arca-wsaa-samples.ps1`
```yaml
block_id: "ARCA-WSAA:acquire:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed acquisition of official ARCA URLs"
license: "LicenseRef-Workspace-Owner"
sha256: "035588dd44099e47c47cae77733b60c28279bbe9524222999a8744d75d20faaf"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$CacheDirectory,
  [Parameter(Mandatory=$true)][string]$OutputDirectory,
  [switch]$Offline
)
$ErrorActionPreference='Stop'
$utf8=[Text.UTF8Encoding]::new($false)
$sources=@(
  @{Name='powershell';Url='https://arca.gob.ar/ws/WSAA/ejemplos/dev-wsaa-cliente-powershell.zip';Archive='powershell.zip';Bytes=48293;Sha256='6d2a36105dafef967cc3e92a001e9f0f8746135406bbd99e0e3195ca7f162514';Root='dev-wsaa-cliente-powershell';Files=@{
    'README.md'='2303fe9b36fb2fa5ffc1fddcef1422fea1e76d43730049f0b25b4627c3b96715';
    'source/wsaa-cliente.ps1'='edd4b7d840cf873d1913177fbb8e5beaac266f53770cf5d893a2e554f5d2db6f'
  }}
  @{Name='csharp';Url='https://arca.gob.ar/ws/WSAA/ejemplos/dev-wsaa-cliente-dotnet-cs.zip';Archive='csharp.zip';Bytes=64195;Sha256='f46c8ba58968abf91a96f617d23de256b3c30a29f8f55c3aa9072c88865fbb73';Root='dev-wsaa-cliente-dotnet-cs';Files=@{
    'README.md'='c985c5e1e195fe193e1b47605eeb93b24536d49c463e235618fe6f1aa3d3f436';
    'source/ClienteLoginCms_CS/ClienteLoginCms.cs'='6813ae36b017624eca84acca09f7d2dccbb0eacb4be935efbcf6264e2bffea0b';
    'source/ClienteLoginCms_CS/app.config'='7f84a6c529f2a134e2d9d350f735008548376b312018b3c0192d88bcf0a7965e'
  }}
)
function Assert-Child([string]$Root,[string]$Candidate){$rootPath=[IO.Path]::GetFullPath($Root);$path=[IO.Path]::GetFullPath($Candidate);if(-not $path.StartsWith($rootPath+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){throw "path escaped root: $path"};$path}
function Hash([string]$Path){(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()}
$cache=[IO.Path]::GetFullPath($CacheDirectory);$output=[IO.Path]::GetFullPath($OutputDirectory)
if($cache-eq$output){throw 'cache and output must differ'}
[IO.Directory]::CreateDirectory($cache)|Out-Null
if(Test-Path -LiteralPath $output){if((Get-ChildItem -Force -LiteralPath $output).Count-ne 0){throw 'output must be absent or empty'}}else{[IO.Directory]::CreateDirectory($output)|Out-Null}
foreach($source in $sources){
  $archive=Assert-Child $cache (Join-Path $cache $source.Archive)
  if(-not(Test-Path -LiteralPath $archive)){if($Offline){throw "offline archive missing: $($source.Archive)"};$client=[Net.Http.HttpClient]::new();try{$bytes=$client.GetByteArrayAsync($source.Url).GetAwaiter().GetResult();[IO.File]::WriteAllBytes($archive,$bytes)}finally{$client.Dispose()}}
  if((Get-Item -LiteralPath $archive).Length-ne$source.Bytes-or(Hash $archive)-ne$source.Sha256){throw "archive identity mismatch: $($source.Name)"}
  $extract=Assert-Child $cache (Join-Path $cache ('.extract-'+$source.Name+'-'+[guid]::NewGuid().ToString('N')));[IO.Directory]::CreateDirectory($extract)|Out-Null
  try{
    [IO.Compression.ZipFile]::ExtractToDirectory($archive,$extract)
    $root=Assert-Child $extract (Join-Path $extract $source.Root)
    foreach($relative in $source.Files.Keys){$input=Assert-Child $root (Join-Path $root $relative);if(-not(Test-Path -LiteralPath $input -PathType Leaf)-or(Hash $input)-ne$source.Files[$relative]){throw "selected file mismatch: $($source.Name)/$relative"};$target=Assert-Child $output (Join-Path $output ($source.Name+'/'+$relative));[IO.Directory]::CreateDirectory([IO.Path]::GetDirectoryName($target))|Out-Null;[IO.File]::WriteAllBytes($target,[IO.File]::ReadAllBytes($input))}
    $terms=[IO.File]::ReadAllText((Join-Path $root 'README.md'));if($terms-notmatch'puede ser re-distribuido, publicado o descargado'){throw "redistribution terms missing: $($source.Name)"}
  }finally{if(Test-Path -LiteralPath $extract){[IO.Directory]::Delete($extract,$true)}}
}
$receipt=[ordered]@{schema='elite-arca-wsaa-source/v1';authority='ARCA/AFIP official WSAA documentation';retrieved_utc=[DateTimeOffset]::UtcNow.ToString('O');sources=@($sources|ForEach-Object{[ordered]@{name=$_.Name;url=$_.Url;archive_sha256=$_.Sha256;archive_bytes=$_.Bytes}});production_admitted=$false}
[IO.File]::WriteAllText((Join-Path $output 'SOURCE_RECEIPT.json'),($receipt|ConvertTo-Json -Depth 5)+"`n",$utf8)
'ARCA_WSAA_ACQUISITION_PASS sources=2 selected_files=5 production_admitted=false'
````

### FILE: `tools/test-arca-wsaa-samples.ps1`
```yaml
block_id: "ARCA-WSAA:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact-byte and condition regression"
license: "LicenseRef-Workspace-Owner"
sha256: "924126c9a3ca869f3c97e89d1794b614612ab79cad939336a2a4d8bd9b0d8499"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]param([string]$Root=(Join-Path $PSScriptRoot '..'))
$ErrorActionPreference='Stop';$rootPath=(Resolve-Path -LiteralPath $Root).Path
$expected=@{
  'official/arca/wsaa/powershell/README.md'='f76dae5edf39e2602e0f414eb8fb675f7ed0b1db15cbfb3ad44b679ed061d705';
  'official/arca/wsaa/powershell/wsaa-cliente.ps1'='b7ca06a2fa210f9998f00e3a32378a79a48ae9b53c12aa8960752d9c199ddc2d';
  'official/arca/wsaa/csharp/README.md'='07f8d3930ae8bfee0661a5abc455ee1ded4ccac9f89cf46610429467275c01d5';
  'official/arca/wsaa/csharp/ClienteLoginCms.cs'='58694ce2cae209ad8386163396d336433961c091240cf1cee5c61f55e2746781';
  'official/arca/wsaa/csharp/app.config'='494ad889f15c76aaeffcc07717f59ac2371f97bba0d10fec61ea2bb8c1aecc19'
}
foreach($relative in $expected.Keys){$path=Join-Path $rootPath $relative;if(-not(Test-Path -LiteralPath $path -PathType Leaf)){throw "missing $relative"};$hash=(Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant();if($hash-ne$expected[$relative]){throw "hash mismatch $relative"}}
$tokens=$null;$errors=$null;[Management.Automation.Language.Parser]::ParseFile((Join-Path $rootPath 'official/arca/wsaa/powershell/wsaa-cliente.ps1'),[ref]$tokens,[ref]$errors)|Out-Null;if($errors.Count-ne0){throw 'official PowerShell sample has parser errors'}
$terms=[IO.File]::ReadAllText((Join-Path $rootPath 'official/arca/wsaa/powershell/README.md'));if($terms-notmatch'puede ser re-distribuido, publicado o descargado'-or$terms-notmatch'TODOS SUS ERRORES Y OMISIONES'){throw 'official use terms not preserved'}
$sample=[IO.File]::ReadAllText((Join-Path $rootPath 'official/arca/wsaa/powershell/wsaa-cliente.ps1'));foreach($condition in @('New-WebServiceProxy','myPrivate.key','loginTicketResponse.xml')){if($sample-notmatch[regex]::Escape($condition)){throw "expected upstream condition disappeared: $condition"}}
$csharp=[IO.File]::ReadAllText((Join-Path $rootPath 'official/arca/wsaa/csharp/ClienteLoginCms.cs'));foreach($contract in @('SignedCms','X509Certificate2','loginCms')){if($csharp-notmatch[regex]::Escape($contract)){throw "C# contract missing: $contract"}}
'ARCA_WSAA_SAMPLE_PASS files=5 parser=pass redistribution_terms=preserved production_admitted=false'
````

## 6. Configuration surface

Acquisition accepts cache/output paths and `-Offline`; it has no secret. Official samples accept certificate/private-key/password inputs and therefore must run only in an isolated evaluation environment with non-production credentials. Never commit generated TRA/CMS/LoginTicketResponse files.

## 7. Dependency bill

| Dependency | Pin | Use | License/terms | Source |
|---|---|---|---|---|
| ARCA PowerShell ZIP | 48,293 bytes / SHA `6d2a...514` | exact official WSAA sample | redistribution agreement in README | `arca.gob.ar` |
| ARCA C# ZIP | 64,195 bytes / SHA `f46c...b73` | exact official WSAA sample | redistribution agreement in README | `arca.gob.ar` |
| PowerShell | 7+ | acquisition/test only | MIT | Microsoft |
| Official sample runtimes | legacy Windows PowerShell/OpenSSL and .NET Framework 2.0 | evaluation only | platform terms | upstream README |

## 8. Apply order

Materialize independently as source evidence. Run the exact-byte test, acquire from cache or official HTTPS and compare receipts. Do not add this pack to the product runtime profile. Build the modern WSAA/WSFE adapter separately, with a project-approved certificate and ARCA homologation gate.

## 9. Verification

Require 7/7 reconstruction, five LF-materialized SHA matches, five original-file SHA matches during acquisition, both archive size/hash locks, preserved redistribution terms, PowerShell parser pass, expected legacy-condition checks, offline acquisition and source receipt with `production_admitted=false`. Live WSAA is prohibited without project authority, test certificate and secret controls.

## 10. Reconstruction evidence

Recorded in `reconstruction_evidence/ARCA_WSAA_OFFICIAL_SAMPLES_2026-08-30_V122.md`.
