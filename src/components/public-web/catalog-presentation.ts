import type { PublishedCatalog } from "@/platform/catalog/schema";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { publicText as t, type PublicTextKey } from "./messages";

export function modelClass(value:string,locale:string){const known:Record<string,PublicTextKey>={bicycle:"bicycle",motorcycle:"motorcycle",scooter:"scooter",utility:"utility",car:"car",other:"other"};return t(locale,known[value]??"other");}
export function modelDescription(specification:Record<string,unknown>){return typeof specification.description==="string"?specification.description:null;}
export function publicSpecifications(specification:Record<string,unknown>,locale:string){
  // Present only labeled, existing source fields. Never infer units or decode
  // unknown raw technical structures into unreviewed commercial claims.
  const fields:Record<string,PublicTextKey>={range:"range",weight:"weight",motor:"motor",battery:"battery",maximum_speed:"speed",frame:"frame",wheels:"wheels",brakes:"brakes",warranty:"warranty"};
  return Object.entries(fields).flatMap(([key,label])=>typeof specification[key]==="string"?[{label:t(locale,label),value:specification[key] as string}]:[]);
}
export function variantPrice(variant:PublishedCatalog["variants"][number],currency:string,locale:string){
  const amount=BigInt(variant.amount_minor_units);
  const display=minorAmountPresentation(amount<=BigInt(Number.MAX_SAFE_INTEGER)?Number(amount):NaN,currency,locale);
  return {amount:display.amountValid?display.amountLabel:t(locale,"priceUnavailable"),valid:display.amountValid,tax:t(locale,variant.tax_mode==="inclusive"?"included":variant.tax_mode==="exclusive"?"excluded":"notApplicable")};
}
export function startingPrice(catalog:PublishedCatalog,modelId:string,locale:string){
  const variants=catalog.variants.filter(item=>item.model_id===modelId);
  // Mixed tax modes are not comparable; send the reader to all explicit versions.
  if(!variants.length||new Set(variants.map(item=>item.tax_mode)).size!==1)return null;
  const lowest=variants.reduce((a,b)=>BigInt(a.amount_minor_units)<=BigInt(b.amount_minor_units)?a:b);
  const price=variantPrice(lowest,catalog.currency,locale);
  return price.valid?{...price,prefix:variants.length>1?t(locale,"priceFrom"):""}:null;
}
