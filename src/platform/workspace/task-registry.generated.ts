// AUTHORED projection from the SHA-locked U2 registry and actual franchise pane admission.
// Presence is not authorization: server filters installation, modules, features and session.
export const TASK_REGISTRY_SHA256 = "a4d3ffe72e6dbed550d31f6e9913838b0f8d53ec10bf8e088437610b69fcbf21";
export interface WorkspaceTask { id:string; href:string; label:string; group:string; anyPermissions:readonly string[]; allPermissions:readonly string[]; features:readonly string[]; modules:readonly string[]; icon:string; taskLabel?:string }
export const WORKSPACE_GROUPS = [
  {
    "id": "my-work",
    "label": "Mi trabajo"
  },
  {
    "id": "commercial",
    "label": "Comercial"
  },
  {
    "id": "operations",
    "label": "Operación"
  },
  {
    "id": "documents",
    "label": "Documentos"
  },
  {
    "id": "care",
    "label": "Atención"
  },
  {
    "id": "people-network",
    "label": "Personas y red"
  },
  {
    "id": "analysis",
    "label": "Análisis"
  },
  {
    "id": "settings",
    "label": "Configuración"
  }
] as const;
export const WORKSPACE_TASKS: readonly WorkspaceTask[] = [
  {
    "id": "route:/admin/catalog",
    "href": "/admin/catalog",
    "label": "Catálogo",
    "group": "commercial",
    "anyPermissions": [
      "catalog:read"
    ],
    "allPermissions": [],
    "features": [
      "catalog_editor"
    ],
    "modules": [],
    "icon": "/admin/catalog"
  },
  {
    "id": "route:/admin",
    "href": "/admin",
    "label": "Resumen",
    "group": "analysis",
    "anyPermissions": [
      "admin:read"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/admin"
  },
  {
    "id": "route:/admin/surveys",
    "href": "/admin/surveys",
    "label": "Encuestas",
    "group": "analysis",
    "anyPermissions": [
      "surveys:read"
    ],
    "allPermissions": [],
    "features": [
      "customer_surveys"
    ],
    "modules": [],
    "icon": "/admin/surveys"
  },
  {
    "id": "route:/customer/appointments",
    "href": "/customer/appointments",
    "label": "Mis turnos",
    "group": "care",
    "anyPermissions": [
      "customer:self"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/customer/appointments"
  },
  {
    "id": "route:/customer/handovers",
    "href": "/customer/handovers",
    "label": "Mis entregas",
    "group": "care",
    "anyPermissions": [
      "customer:self"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/customer/handovers"
  },
  {
    "id": "route:/customer",
    "href": "/customer",
    "label": "Mis pedidos",
    "group": "commercial",
    "anyPermissions": [
      "customer:self"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/customer"
  },
  {
    "id": "route:/customer/quotes",
    "href": "/customer/quotes",
    "label": "Mis cotizaciones",
    "group": "commercial",
    "anyPermissions": [
      "customer:self"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/customer/quotes"
  },
  {
    "id": "route:/customer/surveys",
    "href": "/customer/surveys",
    "label": "Responder encuesta",
    "group": "care",
    "anyPermissions": [
      "customer:self"
    ],
    "allPermissions": [],
    "features": [
      "customer_surveys"
    ],
    "modules": [],
    "icon": "/customer/surveys"
  },
  {
    "id": "route:/dashboard",
    "href": "/dashboard",
    "label": "Mi panel",
    "group": "my-work",
    "anyPermissions": [],
    "allPermissions": [],
    "features": [
      "role_workspace"
    ],
    "modules": [],
    "icon": "/dashboard"
  },
  {
    "id": "route:/factory",
    "href": "/factory",
    "label": "Producción",
    "group": "operations",
    "anyPermissions": [
      "factory:read"
    ],
    "allPermissions": [],
    "features": [
      "factory_portal"
    ],
    "modules": [
      "procurement"
    ],
    "icon": "/factory"
  },
  {
    "id": "route:/franchise/stored-value",
    "href": "/franchise/stored-value",
    "label": "Saldos y beneficios",
    "group": "commercial",
    "anyPermissions": [
      "stored_value:read"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/franchise/stored-value"
  },
  {
    "id": "route:/franchise/whatsapp",
    "href": "/franchise/whatsapp",
    "label": "Respuestas",
    "group": "care",
    "anyPermissions": [
      "whatsapp:approve"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/franchise/whatsapp"
  },
  {
    "id": "route:/guide",
    "href": "/guide",
    "label": "Guía de trabajo",
    "group": "people-network",
    "anyPermissions": [],
    "allPermissions": [],
    "features": [
      "role_workspace"
    ],
    "modules": [],
    "icon": "/guide"
  },
  {
    "id": "route:/guide/training",
    "href": "/guide/training",
    "label": "Capacitación",
    "group": "people-network",
    "anyPermissions": [
      "training:learn",
      "training:review"
    ],
    "allPermissions": [],
    "features": [
      "training_portal"
    ],
    "modules": [],
    "icon": "/guide/training"
  },
  {
    "id": "route:/help/library",
    "href": "/help/library",
    "label": "Centro de ayuda",
    "group": "people-network",
    "anyPermissions": [
      "help:read",
      "help:write",
      "help:publish"
    ],
    "allPermissions": [],
    "features": [
      "help_cms"
    ],
    "modules": [],
    "icon": "/help/library"
  },
  {
    "id": "route:/help",
    "href": "/help",
    "label": "Ayuda",
    "group": "care",
    "anyPermissions": [],
    "allPermissions": [],
    "features": [],
    "modules": [],
    "icon": "/help"
  },
  {
    "id": "route:/network",
    "href": "/network",
    "label": "Red y acuerdos",
    "group": "people-network",
    "anyPermissions": [
      "network:admin",
      "franchise:write"
    ],
    "allPermissions": [],
    "features": [
      "network_portal"
    ],
    "modules": [],
    "icon": "/network"
  },
  {
    "id": "route:/supply",
    "href": "/supply",
    "label": "Compras",
    "group": "operations",
    "anyPermissions": [
      "supply:read",
      "supply:factory-read"
    ],
    "allPermissions": [],
    "features": [
      "supply_portal"
    ],
    "modules": [],
    "icon": "/supply"
  },
  {
    "id": "route:/warranty",
    "href": "/warranty",
    "label": "Garantías",
    "group": "care",
    "anyPermissions": [
      "warranty:read",
      "warranty:factory-read",
      "warranty:self"
    ],
    "allPermissions": [],
    "features": [
      "warranty_portal"
    ],
    "modules": [],
    "icon": "/warranty"
  },
  {
    "id": "franchise:daily",
    "href": "/franchise?task=daily",
    "label": "Gestión diaria",
    "group": "my-work",
    "anyPermissions": [
      "lead:read",
      "handover:manage",
      "resource:manage",
      "availability:read",
      "availability:manage",
      "appointment:manage"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [
      "crm"
    ],
    "icon": "Gestión diaria",
    "taskLabel": "Gestión diaria"
  },
  {
    "id": "franchise:agenda",
    "href": "/franchise?task=agenda",
    "label": "Agenda",
    "group": "my-work",
    "anyPermissions": [
      "appointment:manage"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [
      "crm"
    ],
    "icon": "Agenda",
    "taskLabel": "Agenda"
  },
  {
    "id": "franchise:orders",
    "href": "/franchise?task=orders",
    "label": "Pedidos",
    "group": "commercial",
    "anyPermissions": [
      "inventory:allocate",
      "payment:create",
      "handover:manage",
      "admin:read"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [
      "crm"
    ],
    "icon": "Pedidos",
    "taskLabel": "Pedidos"
  },
  {
    "id": "franchise:deliveries",
    "href": "/franchise?task=deliveries",
    "label": "Entregas",
    "group": "operations",
    "anyPermissions": [
      "handover:manage"
    ],
    "allPermissions": [],
    "features": [],
    "modules": [
      "crm"
    ],
    "icon": "Entregas",
    "taskLabel": "Entregas"
  }
];
// These pages are verified in this overlay's input composition; missing extensions are excluded.
export const INSTALLED_WORKSPACE_ROUTES: readonly string[] = [
  "/admin",
  "/admin/catalog",
  "/admin/surveys",
  "/customer",
  "/customer/appointments",
  "/customer/handovers",
  "/customer/quotes",
  "/customer/surveys",
  "/dashboard",
  "/factory",
  "/franchise",
  "/franchise/stored-value",
  "/franchise/whatsapp",
  "/guide",
  "/guide/training",
  "/help",
  "/help/library",
  "/network",
  "/supply",
  "/warranty"
];
