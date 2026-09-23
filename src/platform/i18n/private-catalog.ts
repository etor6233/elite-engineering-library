// AUTHORED es/en presentation catalog for local glue; no corporate code attribution.
import messages from "./private-messages.json";
import type {PrivateLanguage} from "./private-locale";
export type PrivateMessageKey=keyof typeof messages;
type Parameters={
 p0248:{resolution_action:string};
 p0255:{completed_at:string};
 p0267:{accepted_at:string};
 p0310:{request_id:string};
 p0629:{handover_id:string;handover_state:string};
 p0645:{evaluated_at:string;valid_until:string};
 p0646:{evaluated_at:string};
};
type Arguments<K extends PrivateMessageKey>=K extends keyof Parameters?[values:Parameters[K]]:[];
export function privateMessage<K extends PrivateMessageKey>(language:PrivateLanguage,key:K,...args:Arguments<K>):string{
 const value=messages[key][language];
 if(args.length===0)return value;
 const values=args[0] as Record<string,string>;
 return value.replace(/\{([a-z_]+)\}/g,(_,name:string)=>{if(typeof values[name]!=="string")throw new Error("missing display parameter");return values[name]});
}
export const privateTranslator=(language:PrivateLanguage)=><K extends PrivateMessageKey>(key:K,...args:Arguments<K>)=>privateMessage(language,key,...args);
export type PrivateTranslator=ReturnType<typeof privateTranslator>;
const known=new Map<string,PrivateMessageKey>(Object.entries(messages).filter(([,m])=>!m.es.includes("{")).map(([key,m])=>[m.es,key as PrivateMessageKey]));
// Call ONLY for locally controlled display labels, never CMS/provider/user content.
export function controlledPrivateLabel(language:PrivateLanguage,value:string):string{
 const key=known.get(value);return key===undefined?value:privateMessage(language,key as Exclude<PrivateMessageKey,keyof Parameters>);
}
