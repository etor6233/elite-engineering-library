# Receipt — machine-local-home-path V403

Allowlisted JSON: `RECEIPT_MACHINE_LOCAL_HOME_PATH_V403.json`.

Former local staging name `RECEIPT_machine-local-home-path_V403.json` renamed to match `elite-public-release-input-policy/v1` `json_files` pattern `^[A-Z][A-Z0-9_]*\.json$`. Evidence not discarded.

Fix: image-default SDK profile home → `/app/.home` in `implementation_packs/FRANCHISE_CLOUD_EXECUTION_V403.md` with `CLOUDSDK_CONFIG=/app/.home/.config/gcloud` and `chown 1000:1000 /app/.home`.
