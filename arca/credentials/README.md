# Elite ARCA WSAA credential core

This library creates the official WSAA login ticket request, signs it as attached CMS/PKCS#7 with Microsoft platform cryptography, calls an injected `loginCms` transport and caches the returned token/sign in memory with single-flight refresh.

Production configuration must identify an ARCA-issued certificate already installed in an operating-system or managed certificate store. PFX bytes, passwords, token, sign and private keys are never configuration values, logs or receipts. The library never persists credentials.

The default invoice profile requests service `wsfe` for at most twelve hours. A generated WSAA proxy implements `IWsaaTransport`; fiscal issuance remains a separate durable boundary.
