# ARCA WSFE local Unix-socket worker

This process is the narrow boundary between the Go fiscal owner and the Microsoft-generated WSAA/WSFE clients. It listens only on an absolute Unix-domain socket, accepts a strict snake-case JSON contract, invokes homologation endpoints and returns only safe result fields. It has no TCP listener and never returns or logs Token, Sign, certificate material, private keys or provider messages.

Required runtime inputs are `ARCA_ENVIRONMENT=homologation`, `ARCA_WSFE_SOCKET`, `ARCA_CERTIFICATE_THUMBPRINT` and `ARCA_CERTIFICATE_STORE=CurrentUser|LocalMachine`. Provision the socket directory with least privilege before startup and ensure the socket path does not exist. The process creates the socket and sets mode `0660` on Unix; Kestrel does not create or secure the parent directory.

Production is deliberately blocked. Before promotion, generate current official clients, prove certificate association and current ARCA parameter tables, run the owner/worker reconciliation flow against homologation, approve tax/accounting policy, and execute deployment, monitoring and rollback gates for the selected runtime.
