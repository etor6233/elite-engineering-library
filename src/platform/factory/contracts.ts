// AUTHORED transport contract; state names mirror the existing operations owner.
import { z } from "zod";
export const factoryID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,199}$/);
export const factoryUnitSchema=z.object({id:factoryID,organization_id:factoryID,purchase_order_id:factoryID,variant_id:factoryID,serial_number:z.string().min(1).max(200),state:z.enum(["planned","assembly","quality","released","shipped","received","rejected"])});
export type FactoryUnit=z.infer<typeof factoryUnitSchema>;
