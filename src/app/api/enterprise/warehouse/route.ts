import {warehouseGateway} from "@/platform/warehouse/server";
export const runtime="nodejs";
const handle=(request:Request)=>process.env.WAREHOUSE_WORKSPACE_ENABLED==="true"?warehouseGateway(request):Response.json({code:"WAREHOUSE_NOT_CONFIGURED",effect:"NOT_ATTEMPTED"},{status:503,headers:{"cache-control":"no-store"}});
export const GET=handle;export const POST=handle;
