import {redirect} from "next/navigation";
import type {Route} from "next";
import {readSession,allowed} from "@/platform/auth/session";
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {WarehouseWorkspace} from "@/components/warehouse/warehouse-workspace";
export default async function WarehousePage(){const session=await readSession();if(!session)redirect('/api/auth/login?return_to=/operations/warehouse-bulk' as Route);const {locale}=await loadPrivateI18n();if(process.env.WAREHOUSE_WORKSPACE_ENABLED!=='true'||!allowed(session,'inventory:read'))return <><h1>{locale.language==='en'?'Warehouse':'Almacén'}</h1><p>{locale.language==='en'?'Access unavailable.':'Acceso no disponible.'}</p></>;return <WarehouseWorkspace language={locale.language}/>}
