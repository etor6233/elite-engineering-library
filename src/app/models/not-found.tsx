import Link from "next/link";
import {loadPublicLocale} from "@/platform/i18n/load-public-locale";
import {publicText as t} from "@/components/public-web/messages";
import s from "@/components/public-web/public-web.module.css";
export default async function MissingModel(){const locale=await loadPublicLocale();const en=locale.language==="en";return<div className={s.empty} data-public-state="empty"><h1>{en?"This model is not available.":"Este modelo no está disponible."}</h1><Link href="/models" className={s.secondary}>{t(locale.locale,"backModels")}</Link></div>;}
