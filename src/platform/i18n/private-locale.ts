// AUTHORED bounded es/en display preference; Intl remains the locale authority.
import{resolvePublicLocale,type PublicLanguage}from"./public-catalog";
export type PrivateLanguage=PublicLanguage;
export type PrivateLocale={language:PrivateLanguage;locale:string;timeZone:string;source:"preference"|"browser"|"configuration"};
export function privateLanguage(v:unknown):PrivateLanguage|null{return v==="es"||v==="en"?v:null}
// Narrow two-catalog preference, not a general RFC4647 matcher. Invalid ranges
// are ignored; q=0 excludes their base language. Oversize/unfulfillable headers
// are disregarded and the explicit configured catalog is used (RFC9110 12.5.4).
export function browserLanguage(header:string|null):PrivateLanguage|null{
 if(!header||header.length>512)return null;
 const ranges=header.split(",");if(ranges.length>16)return null;
 const candidates:{language:PrivateLanguage;quality:number;order:number}[]=[],excluded=new Set<PrivateLanguage>();
 for(const [order,raw]of ranges.entries()){
  const match=/^\s*([A-Za-z]{1,8}(?:-[A-Za-z0-9]{1,8})*)(?:\s*;\s*q=(0(?:\.[0-9]{0,3})?|1(?:\.0{0,3})?))?\s*$/i.exec(raw);if(!match)continue;
  try{const canonical=Intl.getCanonicalLocales(match[1]!)[0];if(!canonical)continue;const language=privateLanguage(new Intl.Locale(canonical).language);if(!language)continue;const quality=match[2]===undefined?1:Number(match[2]);if(quality===0)excluded.add(language);else candidates.push({language,quality,order})}catch{/* malformed tag has no authority */}
 }
 return candidates.filter(x=>!excluded.has(x.language)).sort((a,b)=>b.quality-a.quality||a.order-b.order)[0]?.language??null;
}
export function resolvePrivateLocale(configured:string,timeZone:string,preference:unknown,header:string|null):PrivateLocale{
 const baseline=resolvePublicLocale(configured,timeZone),chosen=privateLanguage(preference),browser=chosen?null:browserLanguage(header),language=chosen??browser??baseline.language;
 return{language,locale:language===baseline.language?baseline.locale:language==="en"?"en-US":"es-AR",timeZone:baseline.timeZone,source:chosen?"preference":browser?"browser":"configuration"};
}
