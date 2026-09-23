# ARCA WSAA and WSFEv1 generated clients

This component generates separate WCF clients from exact public ARCA WSAA and WSFEv1 WSDLs. It does not embed either WSDL or claim that generated code is an ARCA SDK.

Required inputs are an official .NET 10.0.400 SDK executable, a cache directory and an empty output directory. Homologation is the default and production metadata requires an explicit switch. The runner fixes `dotnet-svcutil` 8.0.0, WCF 10.0.652802 and patched cryptography 10.0.11, builds with warnings as errors and rejects vulnerable packages reported by current NuGet sources.

This is client generation only. It does not provide WSAA credentials, private-key custody, fiscal rules, authorization, idempotency, reconciliation or ARCA homologation.
