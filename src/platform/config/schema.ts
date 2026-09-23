import { z } from "zod";

const identifier = z.string().regex(/^[a-z][a-z0-9_]*$/, "use snake_case identifiers");
const permission = z.string().regex(/^\*$|^[a-z][a-z0-9_]*(?::[a-z][a-z0-9_]*){1,2}$/);

const fieldSchema = z.discriminatedUnion("type", [
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("string"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("number"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("boolean"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("date"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("select"), required: z.boolean(), options: z.array(identifier).min(1) })
]);

const workflowSchema = z.object({
  initial: identifier,
  states: z.array(identifier).min(1),
  transitions: z.array(z.object({ from: identifier, to: identifier, permission }))
});

export const businessConfigSchema = z.object({
  schemaVersion: z.literal("1.0.0"),
  business: z.object({
    id: z.string().regex(/^[a-z][a-z0-9-]*$/),
    name: z.string().min(1),
    defaultLocale: z.string().min(2),
    defaultMarket: z.string().length(2),
    supportEmail: z.email()
  }),
  markets: z.array(z.object({
    code: z.string().length(2),
    name: z.string().min(1),
    currency: z.string().length(3),
    locales: z.array(z.string().min(2)).min(1),
    timeZone: z.string().min(1),
    taxMode: z.enum(["internal", "external"])
  })).min(1),
  organizationTypes: z.array(z.object({
    id: identifier,
    label: z.string().min(1),
    allowedParents: z.array(identifier)
  })).min(1),
  roles: z.array(z.object({
    id: identifier,
    label: z.string().min(1),
    permissions: z.array(permission).min(1)
  })).min(1),
  modules: z.record(identifier, z.object({ enabled: z.boolean() })),
  workflows: z.record(identifier, workflowSchema),
  customFields: z.record(identifier, z.array(fieldSchema)),
  integrations: z.array(z.object({
    id: identifier,
    provider: identifier,
    enabled: z.boolean(),
    mode: z.enum(["sandbox", "production"]),
    capabilities: z.array(identifier).min(1),
    credentialRefEnv: z.string().regex(/^[A-Z][A-Z0-9_]+$/)
  })),
  features: z.record(identifier, z.boolean())
}).superRefine((config, context) => {
  const unique = (values: string[], path: (string | number)[], label: string) => {
    if (new Set(values).size !== values.length) {
      context.addIssue({ code: "custom", message: `${label} must be unique`, path });
    }
  };

  unique(config.markets.map((item) => item.code), ["markets"], "market codes");
  unique(config.organizationTypes.map((item) => item.id), ["organizationTypes"], "organization type ids");
  unique(config.roles.map((item) => item.id), ["roles"], "role ids");
  unique(config.integrations.map((item) => item.id), ["integrations"], "integration ids");

  const market = config.markets.find((item) => item.code === config.business.defaultMarket);
  if (!market) {
    context.addIssue({ code: "custom", message: "defaultMarket must reference an existing market", path: ["business", "defaultMarket"] });
  } else if (!market.locales.includes(config.business.defaultLocale)) {
    context.addIssue({ code: "custom", message: "defaultLocale must belong to defaultMarket", path: ["business", "defaultLocale"] });
  }

  const organizationTypeIds = new Set(config.organizationTypes.map((item) => item.id));
  config.organizationTypes.forEach((type, index) => {
    type.allowedParents.forEach((parent) => {
      if (!organizationTypeIds.has(parent)) {
        context.addIssue({ code: "custom", message: `unknown parent organization type: ${parent}`, path: ["organizationTypes", index, "allowedParents"] });
      }
      if (parent === type.id) {
        context.addIssue({ code: "custom", message: "organization type cannot parent itself", path: ["organizationTypes", index, "allowedParents"] });
      }
    });
  });

  for (const [name, workflow] of Object.entries(config.workflows)) {
    unique(workflow.states, ["workflows", name, "states"], `${name} states`);
    const states = new Set(workflow.states);
    if (!states.has(workflow.initial)) {
      context.addIssue({ code: "custom", message: "initial state must be declared", path: ["workflows", name, "initial"] });
    }
    workflow.transitions.forEach((transition, index) => {
      if (!states.has(transition.from) || !states.has(transition.to)) {
        context.addIssue({ code: "custom", message: "transition references an unknown state", path: ["workflows", name, "transitions", index] });
      }
      if (transition.from === transition.to) {
        context.addIssue({ code: "custom", message: "self transitions are not allowed", path: ["workflows", name, "transitions", index] });
      }
    });
  }

  for (const [entity, fields] of Object.entries(config.customFields)) {
    unique(fields.map((field) => field.id), ["customFields", entity], `${entity} custom field ids`);
  }

  const moduleDependencies: Record<string, string[]> = {
    payments: ["orders"],
    fulfillment: ["orders", "inventory"],
    service: ["catalog"],
    procurement: ["inventory"],
    integrations: ["catalog"]
  };
  for (const [moduleName, dependencies] of Object.entries(moduleDependencies)) {
    if (config.modules[moduleName]?.enabled) {
      dependencies.forEach((dependency) => {
        if (!config.modules[dependency]?.enabled) {
          context.addIssue({ code: "custom", message: `${moduleName} requires enabled module ${dependency}`, path: ["modules", moduleName] });
        }
      });
    }
  }
});

export type BusinessConfig = z.infer<typeof businessConfigSchema>;
export type WorkflowConfig = z.infer<typeof workflowSchema>;
export type CustomField = z.infer<typeof fieldSchema>;
