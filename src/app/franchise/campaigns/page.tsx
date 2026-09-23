import {createHash} from "node:crypto";
import {redirect,notFound} from "next/navigation";
import type {Route} from "next";
import {readSession} from "@/platform/auth/session";
import {campaignPermissions} from "@/platform/campaigns/contract";
import {CampaignWorkspace} from "@/components/campaigns/campaign-workspace";
export default async function CampaignPage(){
 if(process.env.CAMPAIGN_WORKSPACE_ENABLED!=="true")notFound();
 const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/franchise/campaigns" as Route);
 const permissions=campaignPermissions(s.permissions);if(!permissions.read)return <><h1>Campañas</h1><p>No tenés acceso a este espacio.</p></>;
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,[...s.organizations].sort(),s.subject,[...s.permissions].sort()])).digest("hex");
 return <CampaignWorkspace key={scope} scope={scope} subject={s.subject} permissions={permissions}/>;
}
