import{redirect}from"next/navigation";
import type{Route}from"next";
import{readSession}from"@/platform/auth/session";
import{financeEntry}from"@/platform/finance-vnext/gateway";
import{financeConfig}from"@/platform/finance-vnext/server";
import{FinanceWorkspace}from"@/components/finance-next/finance-workspace";
export default async function FinancePage(){const session=await readSession();if(!session)redirect('/api/auth/login?return_to=/admin/finance' as Route);let config;try{config=financeConfig()}catch{return <><h1>Finanzas</h1><p>Este espacio no está disponible.</p></>};if(session.tenantId!==config.tenant||!session.organizations.includes(config.organization)||!financeEntry(session))return <><h1>Finanzas</h1><p>No tenés acceso a este espacio.</p></>;return <FinanceWorkspace/>}
