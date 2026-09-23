import s from "./public-web.module.css";
import {loadPublicLocale} from "@/platform/i18n/load-public-locale";
import {publicText as t} from "./messages";
export default async function PublicLoading(){const locale=await loadPublicLocale();return<div className={s.skeleton} role="status" aria-busy="true" data-public-state="loading"><p>{t(locale.locale,"loadModels")}</p><span className={s.skeletonTitle} aria-hidden="true"/><span className={s.skeletonPanel} aria-hidden="true"/></div>;}
