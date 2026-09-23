# Microsoft AVM secure SFTP landing

This materializes a fail-closed SFTP landing account from Microsoft Azure Verified Modules, not a custom SFTP server. The composition pins `avm/res/storage/storage-account:0.33.0`, enables hierarchical namespace and SFTP, accepts only an SSH public key, grants the uploader only create/write (`cw`) in one quarantine container, disables Azure Storage Shared Key and public blob access, denies the public network, and creates a Blob private endpoint with private DNS.

The landing account intentionally uses soft delete rather than blob versioning or WORM. Microsoft documents that versioning is unavailable with hierarchical namespace and immutable storage is unavailable while SFTP is enabled. A project must therefore copy the exact uploaded bytes, after close/reconciliation and SHA-256 calculation, into the separately selected retained-original store before extraction or persistence. Compose the existing Azure retained-evidence lane and durable document pipeline for that boundary.

Run `pwsh ./verify_contract.ps1 -AllowNetwork` to download the exact Microsoft Bicep CLI `0.46.1`, verify its published SHA-256, restore the pinned AVM module, and compile this source. No Azure resource is created by the verifier. A live deployment remains blocked until the project supplies an Azure subscription, region, private network/DNS, Log Analytics workspace, SSH key ownership/rotation, cost approval, downstream identity, Event Grid/durable reconciliation, malware gate, retained store, backup/restore and a real SFTP transfer test.

