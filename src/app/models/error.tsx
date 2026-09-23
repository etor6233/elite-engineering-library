"use client";
import {usePublicWebLocale} from "@/components/public-web/public-shell";
import {publicText as t} from "@/components/public-web/messages";
import s from "@/components/public-web/public-web.module.css";
// Do not render or log the server error payload into the public page.
export default function ErrorBoundary({reset}:{error:Error&{digest?:string};reset:()=>void}){const locale=usePublicWebLocale();return<div className={s.empty} data-public-state="error"><h1>{t(locale,"catalogueError")}</h1><p>{t(locale,"catalogueErrorBody")}</p><button type="button" className={s.secondary} onClick={reset}>{t(locale,"retry")}</button></div>;}
