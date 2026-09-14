# Portable Business Configuration Core

## 1. Metadata

```yaml
pack_id: "PBC-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un contrato portable y estricto, más validadores semánticos independientes en Go y Python, para describir una organización, sus unidades, capacidades, identidad y workflows sin imponer lenguaje de aplicación."
stacks:
  - "JSON Schema Draft 2020-12"
  - "Go 1.26.7 validator"
  - "Python 3.14.4 validator"
compatible_with:
  - "cualquier runtime con validador JSON Schema Draft 2020-12"
incompatible_with:
  - "runtimes que sólo implementan parcialmente Draft 2020-12 sin unevaluatedProperties"
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://json-schema.org/draft/2020-12"
verified_at: "2026-08-25"
```

`LicenseRef-Workspace-Owner` significa que los bloques son `AUTHORED` en este workspace y no copian código upstream, pero el propietario debe elegir una licencia de distribución antes de publicarlos. Hasta entonces el pack no puede alcanzar `REUSABLE_PACK` para terceros.

## 2. Applicability

### Resuelve

- configuración inicial de un negocio mediante un solo documento validable;
- selección de módulos sin crear ramas del core por cliente;
- descripción de sedes, franquicias, fábricas, depósitos, proveedores y canales;
- roles declarativos y flujos de trabajo configurables;
- base portable para que distintos adapters generen backend, frontend e infraestructura.

### Adoptar cuando

- el sistema debe adaptarse a distintos negocios mediante configuración;
- se necesita validar antes de iniciar tráfico;
- múltiples stacks deben interpretar el mismo contrato.

### Rechazar o extender cuando

- una regla fiscal, contable, financiera, de homologación o de marketplace exige semántica legal específica;
- se requiere ejecutar expresiones arbitrarias: este pack las prohíbe deliberadamente;
- el negocio no es multiunidad ni necesita configuración portable.

### Límite del claim

El schema prueba forma y restricciones locales. No prueba que un CUIT/VAT sea legalmente válido, que una moneda pueda usarse en un mercado, que un permiso sea seguro ni que un workflow preserve invariantes del dominio. Esas reglas pertenecen a extensions con tests.

## 3. Architecture contract

```text
business-profile.json (control plane, versionado)
             │ validación estricta
             ▼
     configuración admitida
       ├── backend adapters
       ├── frontend navigation/forms
       ├── authorization policy compiler
       ├── workflow compiler
       └── infrastructure profile
```

Invariantes:

- `tenant.id` es estable y globalmente único;
- códigos humanos son únicos dentro del documento y no son secretos;
- sólo se activan módulos enumerados;
- permisos usan namespace `recurso:acción`;
- un workflow sólo declara estados/transiciones; ejecutar una transición siempre requiere autorización y persistencia atómica en el adapter;
- ninguna credencial, token o clave privada entra en este archivo;
- campos desconocidos fallan por defecto.

Failure mode seguro: si el documento no valida, el servicio no inicia ni aplica una versión parcial.

Presupuesto: validación completa menor a 100 ms para documentos de hasta 1 MiB en hardware de desarrollo; el adapter debe medirlo y fijar un límite inferior si recibe configuración no confiable.

Rollback: conservar la última versión admitida y aplicar cambios con compare-and-swap. El schema por sí solo no implementa almacenamiento ni despliegue.

## 4. Exact file manifest

```text
CREATE contracts/business-profile.schema.json
CREATE config/business-profile.example.json
CREATE validators/python/business_profile_semantics.py
CREATE validators/python/test_business_profile_semantics.py
CREATE validators/go/go.mod
CREATE validators/go/profile/semantics.go
CREATE validators/go/profile/semantics_test.go
```

## 5. Materialization blocks

### FILE: `contracts/business-profile.schema.json`

```yaml
block_id: "PBC-CORE:business-profile-schema:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b8d8a8b199005301c6aa9d991cbce07b8e811e77563b2ab0d802390ec84c570d"
variables: []
secrets_allowed: false
```

````json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.invalid/schemas/business-profile.schema.json",
  "title": "Portable Business Profile",
  "description": "Framework-neutral configuration contract. Secrets are forbidden.",
  "type": "object",
  "required": [
    "schemaVersion",
    "profileVersion",
    "tenant",
    "regionalSettings",
    "capabilities",
    "organizationUnits",
    "identity",
    "workflows"
  ],
  "properties": {
    "schemaVersion": {
      "const": "1.0.0"
    },
    "profileVersion": {
      "type": "integer",
      "minimum": 1
    },
    "tenant": {
      "$ref": "#/$defs/tenant"
    },
    "regionalSettings": {
      "$ref": "#/$defs/regionalSettings"
    },
    "capabilities": {
      "type": "array",
      "minItems": 1,
      "uniqueItems": true,
      "items": {
        "$ref": "#/$defs/capability"
      }
    },
    "organizationUnits": {
      "type": "array",
      "minItems": 1,
      "maxItems": 10000,
      "items": {
        "$ref": "#/$defs/organizationUnit"
      }
    },
    "identity": {
      "$ref": "#/$defs/identity"
    },
    "workflows": {
      "type": "array",
      "maxItems": 200,
      "items": {
        "$ref": "#/$defs/workflow"
      }
    },
    "customFields": {
      "type": "array",
      "maxItems": 500,
      "items": {
        "$ref": "#/$defs/customField"
      },
      "default": []
    }
  },
  "unevaluatedProperties": false,
  "$defs": {
    "uuid": {
      "type": "string",
      "format": "uuid",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
    },
    "code": {
      "type": "string",
      "minLength": 2,
      "maxLength": 64,
      "pattern": "^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$"
    },
    "permission": {
      "type": "string",
      "minLength": 3,
      "maxLength": 128,
      "pattern": "^[a-z][a-z0-9-]*:[a-z][a-z0-9-]*$"
    },
    "tenant": {
      "type": "object",
      "required": ["id", "code", "legalName", "displayName"],
      "properties": {
        "id": {
          "$ref": "#/$defs/uuid"
        },
        "code": {
          "$ref": "#/$defs/code"
        },
        "legalName": {
          "type": "string",
          "minLength": 1,
          "maxLength": 200
        },
        "displayName": {
          "type": "string",
          "minLength": 1,
          "maxLength": 120
        },
        "publicDomains": {
          "type": "array",
          "maxItems": 20,
          "uniqueItems": true,
          "items": {
            "type": "string",
            "format": "hostname",
            "maxLength": 253
          },
          "default": []
        }
      },
      "unevaluatedProperties": false
    },
    "regionalSettings": {
      "type": "object",
      "required": ["defaultLocale", "supportedLocales", "timeZone", "currency", "markets"],
      "properties": {
        "defaultLocale": {
          "type": "string",
          "pattern": "^[a-z]{2,3}(?:-[A-Z]{2})?$"
        },
        "supportedLocales": {
          "type": "array",
          "minItems": 1,
          "maxItems": 30,
          "uniqueItems": true,
          "items": {
            "type": "string",
            "pattern": "^[a-z]{2,3}(?:-[A-Z]{2})?$"
          }
        },
        "timeZone": {
          "type": "string",
          "minLength": 3,
          "maxLength": 64,
          "pattern": "^[A-Za-z_+-]+(?:/[A-Za-z0-9_+-]+)+$"
        },
        "currency": {
          "type": "string",
          "pattern": "^[A-Z]{3}$"
        },
        "markets": {
          "type": "array",
          "minItems": 1,
          "maxItems": 100,
          "uniqueItems": true,
          "items": {
            "type": "string",
            "pattern": "^[A-Z]{2}$"
          }
        }
      },
      "unevaluatedProperties": false
    },
    "capability": {
      "enum": [
        "public-web",
        "customer-portal",
        "admin-portal",
        "catalog",
        "crm",
        "sales",
        "inventory",
        "procurement",
        "supplier-management",
        "factory-tracking",
        "logistics",
        "franchise-management",
        "service-and-warranty",
        "payments",
        "marketplace-connectors",
        "advertising-connectors",
        "analytics",
        "notifications"
      ]
    },
    "organizationUnit": {
      "type": "object",
      "required": ["id", "code", "type", "name", "status"],
      "properties": {
        "id": {
          "$ref": "#/$defs/uuid"
        },
        "code": {
          "$ref": "#/$defs/code"
        },
        "type": {
          "enum": [
            "headquarters",
            "franchise",
            "store",
            "factory",
            "warehouse",
            "service-center",
            "supplier",
            "partner"
          ]
        },
        "name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 160
        },
        "status": {
          "enum": ["planned", "active", "suspended", "closed"]
        },
        "parentId": {
          "$ref": "#/$defs/uuid"
        },
        "market": {
          "type": "string",
          "pattern": "^[A-Z]{2}$"
        },
        "contact": {
          "type": "object",
          "properties": {
            "publicEmail": {
              "type": "string",
              "format": "email",
              "maxLength": 254
            },
            "publicPhone": {
              "type": "string",
              "pattern": "^\\+[1-9][0-9]{7,14}$"
            },
            "publicUrl": {
              "type": "string",
              "format": "uri",
              "pattern": "^https://"
            }
          },
          "unevaluatedProperties": false
        }
      },
      "unevaluatedProperties": false
    },
    "identity": {
      "type": "object",
      "required": ["providerMode", "roles"],
      "properties": {
        "providerMode": {
          "enum": ["oidc", "saml", "managed-external"]
        },
        "roles": {
          "type": "array",
          "minItems": 1,
          "maxItems": 100,
          "items": {
            "$ref": "#/$defs/role"
          }
        }
      },
      "unevaluatedProperties": false
    },
    "role": {
      "type": "object",
      "required": ["code", "name", "permissions"],
      "properties": {
        "code": {
          "$ref": "#/$defs/code"
        },
        "name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 100
        },
        "permissions": {
          "type": "array",
          "maxItems": 500,
          "uniqueItems": true,
          "items": {
            "$ref": "#/$defs/permission"
          }
        },
        "scope": {
          "enum": ["tenant", "organization-unit", "self"],
          "default": "tenant"
        }
      },
      "unevaluatedProperties": false
    },
    "workflow": {
      "type": "object",
      "required": ["code", "entityType", "initialState", "states", "transitions"],
      "properties": {
        "code": {
          "$ref": "#/$defs/code"
        },
        "entityType": {
          "$ref": "#/$defs/code"
        },
        "initialState": {
          "$ref": "#/$defs/code"
        },
        "states": {
          "type": "array",
          "minItems": 1,
          "maxItems": 100,
          "uniqueItems": true,
          "items": {
            "$ref": "#/$defs/code"
          }
        },
        "transitions": {
          "type": "array",
          "minItems": 1,
          "maxItems": 500,
          "items": {
            "$ref": "#/$defs/transition"
          }
        }
      },
      "unevaluatedProperties": false
    },
    "transition": {
      "type": "object",
      "required": ["code", "from", "to", "requiredPermission"],
      "properties": {
        "code": {
          "$ref": "#/$defs/code"
        },
        "from": {
          "$ref": "#/$defs/code"
        },
        "to": {
          "$ref": "#/$defs/code"
        },
        "requiredPermission": {
          "$ref": "#/$defs/permission"
        },
        "requiresReason": {
          "type": "boolean",
          "default": false
        }
      },
      "unevaluatedProperties": false
    },
    "customField": {
      "type": "object",
      "required": ["entityType", "code", "label", "valueType", "required"],
      "properties": {
        "entityType": {
          "$ref": "#/$defs/code"
        },
        "code": {
          "$ref": "#/$defs/code"
        },
        "label": {
          "type": "string",
          "minLength": 1,
          "maxLength": 100
        },
        "valueType": {
          "enum": ["string", "integer", "decimal", "boolean", "date", "datetime", "enum"]
        },
        "required": {
          "type": "boolean"
        },
        "enumValues": {
          "type": "array",
          "minItems": 1,
          "maxItems": 200,
          "uniqueItems": true,
          "items": {
            "type": "string",
            "minLength": 1,
            "maxLength": 100
          }
        },
        "publiclyVisible": {
          "type": "boolean",
          "default": false
        }
      },
      "allOf": [
        {
          "if": {
            "properties": {
              "valueType": {
                "const": "enum"
              }
            },
            "required": ["valueType"]
          },
          "then": {
            "required": ["enumValues"]
          },
          "else": {
            "not": {
              "required": ["enumValues"]
            }
          }
        }
      ],
      "unevaluatedProperties": false
    }
  }
}
````

### FILE: `config/business-profile.example.json`

```yaml
block_id: "PBC-CORE:business-profile-example:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "99729bdaac62e349e1905f13503049aae34305bc28240676ca08f856ea27e2f8"
variables: []
secrets_allowed: false
```

````json
{
  "schemaVersion": "1.0.0",
  "profileVersion": 1,
  "tenant": {
    "id": "018f4d4a-7b36-7a21-8d10-2f4c54c28a01",
    "code": "electromobility-group",
    "legalName": "Electromobility Group S.A.",
    "displayName": "Electromobility Group",
    "publicDomains": ["example.com"]
  },
  "regionalSettings": {
    "defaultLocale": "es-AR",
    "supportedLocales": ["es-AR", "en-US"],
    "timeZone": "America/Argentina/Buenos_Aires",
    "currency": "ARS",
    "markets": ["AR"]
  },
  "capabilities": [
    "public-web",
    "customer-portal",
    "admin-portal",
    "catalog",
    "crm",
    "sales",
    "inventory",
    "procurement",
    "supplier-management",
    "factory-tracking",
    "logistics",
    "franchise-management",
    "service-and-warranty",
    "marketplace-connectors",
    "advertising-connectors",
    "analytics",
    "notifications"
  ],
  "organizationUnits": [
    {
      "id": "018f4d4a-7b36-7a21-8d10-2f4c54c28a02",
      "code": "hq-ar",
      "type": "headquarters",
      "name": "Casa central Argentina",
      "status": "active",
      "market": "AR",
      "contact": {
        "publicEmail": "contacto@example.com",
        "publicPhone": "+541155555555",
        "publicUrl": "https://example.com"
      }
    },
    {
      "id": "018f4d4a-7b36-7a21-8d10-2f4c54c28a03",
      "code": "franchise-cordoba",
      "type": "franchise",
      "name": "Franquicia Córdoba",
      "status": "planned",
      "parentId": "018f4d4a-7b36-7a21-8d10-2f4c54c28a02",
      "market": "AR"
    }
  ],
  "identity": {
    "providerMode": "oidc",
    "roles": [
      {
        "code": "tenant-admin",
        "name": "Administrador del tenant",
        "scope": "tenant",
        "permissions": [
          "catalog:read",
          "catalog:write",
          "customer:read",
          "customer:write",
          "inventory:read",
          "inventory:write",
          "order:read",
          "order:transition",
          "organization:manage"
        ]
      },
      {
        "code": "franchise-operator",
        "name": "Operador de franquicia",
        "scope": "organization-unit",
        "permissions": [
          "catalog:read",
          "customer:read",
          "customer:write",
          "inventory:read",
          "order:read",
          "order:transition"
        ]
      },
      {
        "code": "customer",
        "name": "Cliente",
        "scope": "self",
        "permissions": [
          "catalog:read",
          "order:read"
        ]
      }
    ]
  },
  "workflows": [
    {
      "code": "vehicle-order",
      "entityType": "order",
      "initialState": "draft",
      "states": [
        "draft",
        "submitted",
        "confirmed",
        "allocated",
        "in-transit",
        "ready-for-delivery",
        "delivered",
        "cancelled"
      ],
      "transitions": [
        {
          "code": "submit",
          "from": "draft",
          "to": "submitted",
          "requiredPermission": "order:transition"
        },
        {
          "code": "confirm",
          "from": "submitted",
          "to": "confirmed",
          "requiredPermission": "order:transition"
        },
        {
          "code": "allocate",
          "from": "confirmed",
          "to": "allocated",
          "requiredPermission": "inventory:write"
        },
        {
          "code": "dispatch",
          "from": "allocated",
          "to": "in-transit",
          "requiredPermission": "order:transition"
        },
        {
          "code": "receive",
          "from": "in-transit",
          "to": "ready-for-delivery",
          "requiredPermission": "order:transition"
        },
        {
          "code": "deliver",
          "from": "ready-for-delivery",
          "to": "delivered",
          "requiredPermission": "order:transition"
        },
        {
          "code": "cancel-before-allocation",
          "from": "submitted",
          "to": "cancelled",
          "requiredPermission": "order:transition",
          "requiresReason": true
        }
      ]
    }
  ],
  "customFields": [
    {
      "entityType": "vehicle-model",
      "code": "battery-chemistry",
      "label": "Química de batería",
      "valueType": "enum",
      "required": true,
      "enumValues": ["lfp", "nmc"],
      "publiclyVisible": true
    },
    {
      "entityType": "vehicle-model",
      "code": "nominal-range-km",
      "label": "Autonomía nominal (km)",
      "valueType": "integer",
      "required": true,
      "publiclyVisible": true
    }
  ]
}
````

### FILE: `validators/python/business_profile_semantics.py`

```yaml
block_id: "PBC-CORE:python-semantic-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cd8f4e09674c340060742c692dc5807b5df6d8a048ed903138e1ae84a1620e50"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import json
import sys
from pathlib import Path
from typing import Any


CAPABILITY_BY_NAMESPACE = {
    "catalog": "catalog",
    "customer": "crm",
    "inventory": "inventory",
    "order": "sales",
    "organization": None,
}


def _error(errors: list[dict[str, str]], code: str, path: str) -> None:
    errors.append({"code": code, "path": path})


def _duplicates(values: list[str]) -> set[str]:
    seen: set[str] = set()
    duplicates: set[str] = set()
    for value in values:
        if value in seen:
            duplicates.add(value)
        seen.add(value)
    return duplicates


def validate(profile: dict[str, Any]) -> list[dict[str, str]]:
    errors: list[dict[str, str]] = []
    regional = profile.get("regionalSettings", {})
    if regional.get("defaultLocale") not in regional.get("supportedLocales", []):
        _error(errors, "DEFAULT_LOCALE_NOT_SUPPORTED", "/regionalSettings/defaultLocale")

    units = profile.get("organizationUnits", [])
    unit_ids = [unit.get("id", "") for unit in units]
    unit_codes = [unit.get("code", "") for unit in units]
    for duplicate in sorted(_duplicates(unit_ids)):
        _error(errors, "DUPLICATE_ORGANIZATION_ID", f"/organizationUnits/{duplicate}")
    for duplicate in sorted(_duplicates(unit_codes)):
        _error(errors, "DUPLICATE_ORGANIZATION_CODE", f"/organizationUnits/{duplicate}")

    parents = {unit.get("id", ""): unit.get("parentId") for unit in units}
    for unit_id, parent_id in parents.items():
        if parent_id and parent_id not in parents:
            _error(errors, "ORGANIZATION_PARENT_NOT_FOUND", f"/organizationUnits/{unit_id}/parentId")
    visiting: set[str] = set()
    visited: set[str] = set()

    def visit(unit_id: str) -> None:
        if unit_id in visited:
            return
        if unit_id in visiting:
            _error(errors, "ORGANIZATION_CYCLE", f"/organizationUnits/{unit_id}/parentId")
            return
        visiting.add(unit_id)
        parent_id = parents.get(unit_id)
        if parent_id in parents:
            visit(parent_id)
        visiting.remove(unit_id)
        visited.add(unit_id)

    for unit_id in parents:
        visit(unit_id)

    capabilities = set(profile.get("capabilities", []))
    roles = profile.get("identity", {}).get("roles", [])
    for duplicate in sorted(_duplicates([role.get("code", "") for role in roles])):
        _error(errors, "DUPLICATE_ROLE_CODE", f"/identity/roles/{duplicate}")
    granted: set[str] = set()
    for index, role in enumerate(roles):
        for permission in role.get("permissions", []):
            granted.add(permission)
            namespace = permission.partition(":")[0]
            if namespace not in CAPABILITY_BY_NAMESPACE:
                _error(errors, "UNKNOWN_PERMISSION_NAMESPACE", f"/identity/roles/{index}/permissions/{permission}")
                continue
            required = CAPABILITY_BY_NAMESPACE[namespace]
            if required and required not in capabilities:
                _error(errors, "PERMISSION_REQUIRES_DISABLED_CAPABILITY", f"/identity/roles/{index}/permissions/{permission}")

    workflows = profile.get("workflows", [])
    for duplicate in sorted(_duplicates([workflow.get("code", "") for workflow in workflows])):
        _error(errors, "DUPLICATE_WORKFLOW_CODE", f"/workflows/{duplicate}")
    for workflow_index, workflow in enumerate(workflows):
        states = workflow.get("states", [])
        state_set = set(states)
        base = f"/workflows/{workflow_index}"
        if len(states) != len(state_set):
            _error(errors, "DUPLICATE_WORKFLOW_STATE", f"{base}/states")
        if workflow.get("initialState") not in state_set:
            _error(errors, "INITIAL_STATE_NOT_FOUND", f"{base}/initialState")
        transition_codes: set[str] = set()
        for transition_index, transition in enumerate(workflow.get("transitions", [])):
            transition_base = f"{base}/transitions/{transition_index}"
            code = transition.get("code", "")
            if code in transition_codes:
                _error(errors, "DUPLICATE_TRANSITION_CODE", f"{transition_base}/code")
            transition_codes.add(code)
            if transition.get("from") not in state_set:
                _error(errors, "TRANSITION_FROM_STATE_NOT_FOUND", f"{transition_base}/from")
            if transition.get("to") not in state_set:
                _error(errors, "TRANSITION_TO_STATE_NOT_FOUND", f"{transition_base}/to")
            if transition.get("requiredPermission") not in granted:
                _error(errors, "TRANSITION_PERMISSION_NOT_GRANTED", f"{transition_base}/requiredPermission")

    custom_fields = profile.get("customFields", [])
    keys = [f"{field.get('entityType', '')}:{field.get('code', '')}" for field in custom_fields]
    for duplicate in sorted(_duplicates(keys)):
        _error(errors, "DUPLICATE_CUSTOM_FIELD", f"/customFields/{duplicate}")
    return sorted(errors, key=lambda item: (item["path"], item["code"]))


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: business_profile_semantics.py <profile.json>", file=sys.stderr)
        return 2
    profile = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
    errors = validate(profile)
    print(json.dumps({"valid": not errors, "errors": errors}, separators=(",", ":")))
    return 0 if not errors else 1


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `validators/python/test_business_profile_semantics.py`

```yaml
block_id: "PBC-CORE:python-semantic-validator-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9f36738a195f96c6e6648c4000eff82ca675176f9dffe022e97987d4bc4bd5cc"
variables: []
secrets_allowed: false
```

````python
import copy
import json
import unittest
from pathlib import Path

from business_profile_semantics import validate


EXAMPLE = Path(__file__).parents[2] / "config" / "business-profile.example.json"


class SemanticValidationTests(unittest.TestCase):
    def setUp(self) -> None:
        self.profile = json.loads(EXAMPLE.read_text(encoding="utf-8"))

    def codes(self, profile: dict) -> set[str]:
        return {error["code"] for error in validate(profile)}

    def test_example_is_semantically_valid(self) -> None:
        self.assertEqual(validate(self.profile), [])

    def test_rejects_missing_parent_and_cycle(self) -> None:
        missing = copy.deepcopy(self.profile)
        missing["organizationUnits"][1]["parentId"] = "missing"
        self.assertIn("ORGANIZATION_PARENT_NOT_FOUND", self.codes(missing))

        cyclic = copy.deepcopy(self.profile)
        cyclic["organizationUnits"][0]["parentId"] = cyclic["organizationUnits"][1]["id"]
        self.assertIn("ORGANIZATION_CYCLE", self.codes(cyclic))

    def test_rejects_disabled_capability_and_orphan_state(self) -> None:
        disabled = copy.deepcopy(self.profile)
        disabled["capabilities"].remove("inventory")
        self.assertIn("PERMISSION_REQUIRES_DISABLED_CAPABILITY", self.codes(disabled))

        orphan = copy.deepcopy(self.profile)
        orphan["workflows"][0]["transitions"][0]["to"] = "unknown"
        self.assertIn("TRANSITION_TO_STATE_NOT_FOUND", self.codes(orphan))


if __name__ == "__main__":
    unittest.main()
````

### FILE: `validators/go/go.mod`

```yaml
block_id: "PBC-CORE:go-semantic-validator-module:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "501c29f124cd2949609d9ba553e988c2bde60d285f2c46de99ae9f258b228cde"
variables: []
secrets_allowed: false
```

````text
module elite.local/business-profile-validator

go 1.26.0
````

### FILE: `validators/go/profile/semantics.go`

```yaml
block_id: "PBC-CORE:go-semantic-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9beaec812667cc1a998bccfc1ae9d207465e874aabba9739496eb41cad66a8dd"
variables: []
secrets_allowed: false
```

````go
package profile

import (
	"fmt"
	"sort"
	"strings"
)

type Profile struct {
	RegionalSettings struct {
		DefaultLocale    string   `json:"defaultLocale"`
		SupportedLocales []string `json:"supportedLocales"`
	} `json:"regionalSettings"`
	Capabilities      []string           `json:"capabilities"`
	OrganizationUnits []OrganizationUnit `json:"organizationUnits"`
	Identity          Identity           `json:"identity"`
	Workflows         []Workflow         `json:"workflows"`
	CustomFields      []CustomField      `json:"customFields"`
}

type OrganizationUnit struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	ParentID string `json:"parentId"`
}

type Identity struct {
	Roles []Role `json:"roles"`
}
type Role struct {
	Code        string   `json:"code"`
	Permissions []string `json:"permissions"`
}
type Workflow struct {
	Code         string       `json:"code"`
	InitialState string       `json:"initialState"`
	States       []string     `json:"states"`
	Transitions  []Transition `json:"transitions"`
}
type Transition struct {
	Code               string `json:"code"`
	From               string `json:"from"`
	To                 string `json:"to"`
	RequiredPermission string `json:"requiredPermission"`
}
type CustomField struct {
	EntityType string `json:"entityType"`
	Code       string `json:"code"`
}
type ValidationError struct {
	Code string `json:"code"`
	Path string `json:"path"`
}

var capabilityByNamespace = map[string]string{
	"catalog": "catalog", "customer": "crm", "inventory": "inventory",
	"order": "sales", "organization": "",
}

func Validate(value Profile) []ValidationError {
	errors := make([]ValidationError, 0)
	add := func(code, path string) { errors = append(errors, ValidationError{Code: code, Path: path}) }
	if _, ok := stringSet(value.RegionalSettings.SupportedLocales)[value.RegionalSettings.DefaultLocale]; !ok {
		add("DEFAULT_LOCALE_NOT_SUPPORTED", "/regionalSettings/defaultLocale")
	}

	parents := make(map[string]string, len(value.OrganizationUnits))
	unitIDs, unitCodes := map[string]struct{}{}, map[string]struct{}{}
	for _, unit := range value.OrganizationUnits {
		if _, exists := unitIDs[unit.ID]; exists {
			add("DUPLICATE_ORGANIZATION_ID", "/organizationUnits/"+unit.ID)
		}
		if _, exists := unitCodes[unit.Code]; exists {
			add("DUPLICATE_ORGANIZATION_CODE", "/organizationUnits/"+unit.Code)
		}
		unitIDs[unit.ID], unitCodes[unit.Code], parents[unit.ID] = struct{}{}, struct{}{}, unit.ParentID
	}
	for unitID, parentID := range parents {
		if parentID != "" {
			if _, exists := parents[parentID]; !exists {
				add("ORGANIZATION_PARENT_NOT_FOUND", "/organizationUnits/"+unitID+"/parentId")
			}
		}
	}
	visiting, visited := map[string]bool{}, map[string]bool{}
	var visit func(string)
	visit = func(unitID string) {
		if visited[unitID] {
			return
		}
		if visiting[unitID] {
			add("ORGANIZATION_CYCLE", "/organizationUnits/"+unitID+"/parentId")
			return
		}
		visiting[unitID] = true
		if parentID := parents[unitID]; parentID != "" {
			if _, exists := parents[parentID]; exists {
				visit(parentID)
			}
		}
		visiting[unitID], visited[unitID] = false, true
	}
	for unitID := range parents {
		visit(unitID)
	}

	capabilities := stringSet(value.Capabilities)
	granted, roleCodes := map[string]struct{}{}, map[string]struct{}{}
	for roleIndex, role := range value.Identity.Roles {
		if _, exists := roleCodes[role.Code]; exists {
			add("DUPLICATE_ROLE_CODE", "/identity/roles/"+role.Code)
		}
		roleCodes[role.Code] = struct{}{}
		for _, permission := range role.Permissions {
			granted[permission] = struct{}{}
			namespace := strings.SplitN(permission, ":", 2)[0]
			required, known := capabilityByNamespace[namespace]
			path := fmt.Sprintf("/identity/roles/%d/permissions/%s", roleIndex, permission)
			if !known {
				add("UNKNOWN_PERMISSION_NAMESPACE", path)
			} else if required != "" {
				if _, enabled := capabilities[required]; !enabled {
					add("PERMISSION_REQUIRES_DISABLED_CAPABILITY", path)
				}
			}
		}
	}

	workflowCodes := map[string]struct{}{}
	for workflowIndex, workflow := range value.Workflows {
		base := fmt.Sprintf("/workflows/%d", workflowIndex)
		if _, exists := workflowCodes[workflow.Code]; exists {
			add("DUPLICATE_WORKFLOW_CODE", "/workflows/"+workflow.Code)
		}
		workflowCodes[workflow.Code] = struct{}{}
		states := stringSet(workflow.States)
		if len(states) != len(workflow.States) {
			add("DUPLICATE_WORKFLOW_STATE", base+"/states")
		}
		if _, exists := states[workflow.InitialState]; !exists {
			add("INITIAL_STATE_NOT_FOUND", base+"/initialState")
		}
		transitionCodes := map[string]struct{}{}
		for transitionIndex, transition := range workflow.Transitions {
			transitionBase := fmt.Sprintf("%s/transitions/%d", base, transitionIndex)
			if _, exists := transitionCodes[transition.Code]; exists {
				add("DUPLICATE_TRANSITION_CODE", transitionBase+"/code")
			}
			transitionCodes[transition.Code] = struct{}{}
			if _, exists := states[transition.From]; !exists {
				add("TRANSITION_FROM_STATE_NOT_FOUND", transitionBase+"/from")
			}
			if _, exists := states[transition.To]; !exists {
				add("TRANSITION_TO_STATE_NOT_FOUND", transitionBase+"/to")
			}
			if _, exists := granted[transition.RequiredPermission]; !exists {
				add("TRANSITION_PERMISSION_NOT_GRANTED", transitionBase+"/requiredPermission")
			}
		}
	}
	customKeys := map[string]struct{}{}
	for _, field := range value.CustomFields {
		key := field.EntityType + ":" + field.Code
		if _, exists := customKeys[key]; exists {
			add("DUPLICATE_CUSTOM_FIELD", "/customFields/"+key)
		}
		customKeys[key] = struct{}{}
	}
	sort.Slice(errors, func(i, j int) bool {
		if errors[i].Path == errors[j].Path {
			return errors[i].Code < errors[j].Code
		}
		return errors[i].Path < errors[j].Path
	})
	return errors
}

func stringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
````

### FILE: `validators/go/profile/semantics_test.go`

```yaml
block_id: "PBC-CORE:go-semantic-validator-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "f07b9904d328b052026cff63741005ef83aed2afd5a5a40184cd66c99051c71d"
variables: []
secrets_allowed: false
```

````go
package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func example(t *testing.T) Profile {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "..", "config", "business-profile.example.json"))
	if err != nil {
		t.Fatal(err)
	}
	var value Profile
	if err := json.Unmarshal(content, &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func codes(errors []ValidationError) map[string]struct{} {
	result := make(map[string]struct{}, len(errors))
	for _, validationError := range errors {
		result[validationError.Code] = struct{}{}
	}
	return result
}

func TestExampleIsSemanticallyValid(t *testing.T) {
	if errors := Validate(example(t)); len(errors) != 0 {
		t.Fatalf("unexpected errors: %+v", errors)
	}
}

func TestRejectsCycleDisabledCapabilityAndOrphanState(t *testing.T) {
	value := example(t)
	value.OrganizationUnits[0].ParentID = value.OrganizationUnits[1].ID
	value.Capabilities = remove(value.Capabilities, "inventory")
	value.Workflows[0].Transitions[0].To = "unknown"
	actual := codes(Validate(value))
	for _, expected := range []string{"ORGANIZATION_CYCLE", "PERMISSION_REQUIRES_DISABLED_CAPABILITY", "TRANSITION_TO_STATE_NOT_FOUND"} {
		if _, exists := actual[expected]; !exists {
			t.Errorf("missing %s in %+v", expected, actual)
		}
	}
}

func remove(values []string, target string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != target {
			result = append(result, value)
		}
	}
	return result
}
````

## 6. Configuration surface

| Campo | Tipo | Default seguro | Validación | Secreto | Mutabilidad | Efecto |
|---|---|---|---|---|---|---|
| `schemaVersion` | semver fijado | ninguno | `1.0.0` | no | sólo migración | selecciona contrato |
| `profileVersion` | entero | ninguno | `>=1` | no | monotónica | control de concurrencia |
| `tenant` | objeto | ninguno | estricto | no | id inmutable | boundary primario |
| `regionalSettings` | objeto | ninguno | locales/ISO/IANA sintácticos | no | versionada | presentación y defaults |
| `capabilities` | enum[] | ninguno | única/no vacía | no | versionada | módulos habilitados |
| `organizationUnits` | objeto[] | ninguno | forma estricta | no | versionada | topología empresarial |
| `identity.roles` | objeto[] | ninguno | permisos namespaced | no | versionada/gated | autorización compilada |
| `workflows` | objeto[] | `[]` | forma estricta | no | versionada/gated | máquinas de estado |
| `customFields` | objeto[] | `[]` | tipos cerrados | no | versionada | extensibilidad de datos |

Validaciones semánticas que debe añadir el adapter antes de admitir:

- `defaultLocale` pertenece a `supportedLocales`;
- `parentId` referencia una unidad existente y no crea ciclos;
- `code` e `id` son únicos por colección;
- estados y transiciones referencian estados existentes;
- `initialState` pertenece a `states`;
- roles no conceden permisos desconocidos por los módulos activos;
- no se desactiva un módulo con datos o workflows activos sin plan de migración.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| JSON Schema | Draft 2020-12 | contrato | especificación pública | diseño/runtime | <https://json-schema.org/draft/2020-12> |
| Go | `1.26.7` | validador semántico A | BSD-3-Clause | build/runtime de tooling | <https://go.dev/dl/> |
| Python | `3.14.4` | validador semántico B | PSF-2.0 | build/runtime de tooling | <https://www.python.org/downloads/> |
| `jsonschema` | `4.26.0`, sólo evidence externo | metaschema/estructura | MIT | verificación; no materializado | <https://pypi.org/project/jsonschema/4.26.0/> |

El validador concreto pertenece al adapter de stack y debe registrar versión, licencia y soporte real de `format` y `unevaluatedProperties`. Este pack no incorpora una dependencia implícita.

## 8. Apply order

Workspace vacío:

1. Crear exactamente los dos archivos del manifest.
2. Sustituir el `$id` de dominio inválido por un URI estable del proyecto, sin cambiar el dialecto.
3. Elegir un validador Draft 2020-12 y fijarlo en el lockfile del adapter; `jsonschema 4.26.0` es la combinación verificada, no una dependencia silenciosa del pack.
4. Activar assertion explícita de `format`; no asumir que el validador la habilita.
5. Ejecutar validación estructural del ejemplo.
6. Ejecutar al menos uno de los validadores semánticos materializados; en CI ejecutar ambos y exigir resultados equivalentes.
7. Rechazar arranque o publicación ante cualquier error.

Workspace existente:

1. Comparar schemas y producir un reporte de compatibilidad.
2. No sobrescribir configuración activa.
3. Migrar una copia y validar datos reales.
4. Aplicar con versionado optimista y rollback a la versión previa.

## 9. Verification

Gates mínimos del adapter:

```text
STRUCTURE-01  schema valida contra el metaschema 2020-12
STRUCTURE-02  example valida contra el schema con format assertions
NEGATIVE-01   una propiedad desconocida es rechazada
NEGATIVE-02   un campo secreto conocido (password/clientSecret/privateKey) es rechazado
SEMANTIC-01   referencias y unicidad son válidas
SEMANTIC-02   workflows no contienen estados huérfanos
SEMANTIC-03   jerarquía organizacional no contiene ciclos
SECURITY-01   permisos desconocidos o fuera de capabilities son rechazados
PERF-01       documento de límite válido se valida dentro del presupuesto
```

Casos negativos obligatorios que el implementation adapter debe materializar:

- agregar `tenant.clientSecret`;
- usar un `parentId` inexistente;
- duplicar `organizationUnit.code`;
- declarar transición a un estado inexistente;
- definir `enumValues` para un custom field no enum;
- desactivar `inventory` mientras existe permiso `inventory:write`.

Éxito esperado: todos los positivos admitidos, todos los negativos rechazados con ruta y código de error estable.

## 10. Reconstruction evidence

Estado actual: `REBUILD_VERIFIED / CONDITIONED`, evidencia `PBC-20260824-V1`.

- 7 archivos materializados con SHA-256 verificado;
- metaschema Draft 2020-12 y ejemplo: PASS con `jsonschema 4.26.0` y format assertions;
- Python 3.14.4: 3 tests positivos/negativos PASS;
- Go 1.26.7: tests positivos/negativos PASS y `gofmt` sin diff;
- ambos adapters producen el mismo resultado para los casos compartidos.

Condiciones restantes: fijar y licenciar el validador estructural en cada producto, ampliar el corpus diferencial/property a todos los códigos de error, medir el documento límite de 1 MiB y definir/migrar reglas legales, fiscales o específicas del negocio. `CONDITIONED` permite composición sólo con acknowledgement y gates; no equivale a adopción automática para terceros.
- elegir licencia de distribución del workspace.
