// AUTHORED exact installation additions for NEXT-UNIFIED-R2. Feature flags do not install routes.
import {INSTALLED_WORKSPACE_ROUTES} from "./task-registry.generated";
export const CONNECTED_WORKSPACE_ROUTES:readonly string[]=[...INSTALLED_WORKSPACE_ROUTES,"/documents","/franchise/campaigns","/intelligence/activity"];
