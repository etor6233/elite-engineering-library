"use client";
import{createContext,useContext,useMemo,type ReactNode}from"react";
import{privateTranslator,controlledPrivateLabel}from"./private-catalog";
import type{PrivateLocale}from"./private-locale";
const fallback:PrivateLocale={language:"es",locale:"es-AR",timeZone:"UTC",source:"configuration"};
const Context=createContext<PrivateLocale>(fallback);
export function PrivateLocaleProvider({locale,children}:{locale:PrivateLocale;children:ReactNode}){return <Context.Provider value={locale}>{children}</Context.Provider>}
export function usePrivateI18n(){const locale=useContext(Context);return useMemo(()=>({locale,t:privateTranslator(locale.language),controlled:(value:string)=>controlledPrivateLabel(locale.language,value)}),[locale])}
