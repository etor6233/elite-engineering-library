import"server-only";
import{loadPrivateLocale}from"./load-private-locale";
import{privateTranslator,controlledPrivateLabel}from"./private-catalog";
export async function loadPrivateI18n(){const locale=await loadPrivateLocale();return{locale,t:privateTranslator(locale.language),controlled:(value:string)=>controlledPrivateLabel(locale.language,value)}}
