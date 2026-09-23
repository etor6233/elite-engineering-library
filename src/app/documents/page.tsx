import {redirect} from "next/navigation";
import type {Route} from "next";
import {allowed,readSession} from "@/platform/auth/session";
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {documentConfig} from "@/platform/documents-vnext/server";
import {documentId} from "@/platform/documents-vnext/contract";
import {DocumentWorkspace} from "@/components/documents-next/document-workspace";
// Operational entry independent of /experience. The BFF and backend recheck every request.
export default async function DocumentsPage({searchParams}:{searchParams:Promise<{document?:string|string[]}>}){
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/documents" as Route);
 const {locale}=await loadPrivateI18n(),en=locale.language==="en";let config;
 try{config=documentConfig()}catch{return <><h1>{en?"Documents":"Documentos"}</h1><p>{en?"This workspace is unavailable.":"Este espacio no está disponible."}</p></>}
 if(session.tenantId!==config.tenant||!session.organizations.includes(config.organization)||!["documents:read","documents:write","documents:review","documents:process"].some(permission=>allowed(session,permission)))return <><h1>{en?"Documents":"Documentos"}</h1><p>{en?"Your account does not have access to these documents.":"Tu cuenta no tiene acceso a estos documentos."}</p></>;
 const initial=(await searchParams).document;
 if(initial!==undefined&&(typeof initial!=="string"||!documentId.safeParse(initial).success))return <><h1>{en?"Documents":"Documentos"}</h1><p>{en?"The review link is invalid.":"El enlace de revisión no es válido."}</p><a href="/documents">{en?"Open documents":"Abrir documentos"}</a></>;
 return <DocumentWorkspace language={locale.language} {...(initial?{initialDocumentId:initial}:{})}/>;
}
