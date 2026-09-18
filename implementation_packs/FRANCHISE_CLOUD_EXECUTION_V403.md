# Franchise cloud execution extension V403 — portable controller home

## 1. Metadata

```yaml
pack_id: "FRANCHISE-CLOUD-EXECUTION-V403"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Portable franchise cloud controller home under /app/.home; no machine-local /home/<user> or Windows profile paths in distributed materialization"
stacks: ["Python 3.14+ standard library", "Linux amd64 candidate OCI target", "Google Cloud SDK runtime image candidate"]
compatible_with: ["CONTAINER-PACKAGING-CORE", "PORTABLE-CI-GATE-RUNNER", "SECURE-OPS-DELIVERY-CORE"]
incompatible_with: ["image-default SDK profile home", "<USERPROFILE> or C:\\Users\\… or C:/Users/… absolute paths in distributed packs", "production inferred from local fixtures"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["Exact official SDK image pins remain in cloud image locks; no outside product code copied here"]
verified_at: "2026-09-18"
```

## 2. Applicability

Use when materializing the FinOps controller OCI recipe for the V403 cloud extension. The controller image must not depend on the image-default SDK profile home or any operator workstation profile. Local workshop copies may live under `<WORKSHOP>/<LIBRARY_ROOT>` or `<USERPROFILE>`; only the portable `/app/.home` contract below is distributable.

## 3. Architecture contract

The controller runs as `USER 1000:1000` with `HOME=/app/.home` and `CLOUDSDK_CONFIG=/app/.home/.config/gcloud`. Root setup creates `/app/.home/.config/gcloud` and `chown -R 1000:1000 /app/.home` to match the existing GCLOUD_RUNTIME image UID. Reject image-default SDK profile homes, `/home/<username>`, `C:\Users\…`, and `C:/Users/…` in any distributed block. Other V403 cloud files remain conditioned on their owning extension plan; this pack fixes only the portable home surface.

## 4. Exact file manifest

```text
CREATE cloud/Dockerfile.controller
```

## 5. Materialization blocks

### FILE: `cloud/Dockerfile.controller`

```yaml
block_id: "FRANCHISE-CLOUD-EXECUTION-V403:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local orchestration/configuration/test glue over unchanged admitted owners; exact official method and candidate pins declared"
license: "LicenseRef-Workspace-Owner"
sha256: "3e4697df48cfd84c7431569102bf0b578e3f7eeb36b85bceeefa0faf09c70982"
variables: []
secrets_allowed: false
```

````text
# AUTHORED packaging only; exact official SDK image remains a target candidate.
ARG GCLOUD_RUNTIME
FROM ${GCLOUD_RUNTIME}
USER root
RUN test -x /usr/lib/google-cloud-sdk/platform/bundledpythonunix/bin/python3 && \
    ln -s /usr/lib/google-cloud-sdk/platform/bundledpythonunix/bin/python3 /usr/local/bin/python3 && \
    mkdir -p /app /app/.home/.config/gcloud && chown -R 1000:1000 /app/.home
COPY cloud/ /app/cloud/
COPY production_admission_gate/ /app/production_admission_gate/
ENV HOME=/app/.home \
    CLOUDSDK_CONFIG=/app/.home/.config/gcloud \
    CLOUDSDK_PYTHON=/usr/lib/google-cloud-sdk/platform/bundledpythonunix/bin/python3 \
    CLOUDSDK_CORE_DISABLE_USAGE_REPORTING=true \
    CLOUDSDK_COMPONENT_MANAGER_DISABLE_UPDATE_CHECK=true \
    PYTHONDONTWRITEBYTECODE=1
WORKDIR /app
USER 1000:1000
ENTRYPOINT ["python3", "-B", "/app/cloud/finops_cycle.py"]
````
