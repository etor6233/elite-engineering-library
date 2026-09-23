// AUTHORED display adapter. Immutable source content and training hashes are unchanged.
import{ALL_GUIDES}from"@/platform/help/content";
import{controlledPrivateLabel}from"./private-catalog";
import type{PrivateLanguage}from"./private-locale";
type Guide={id:string;version:string;title:string;paragraphs:readonly string[];summary?:string;inlineLead?:boolean};
export function guideDisplay(source:Guide,language:PrivateLanguage){
 const original=ALL_GUIDES.find(x=>x.id===source.id&&x.version===source.version);
 const known=original!==undefined&&original.title===source.title&&(source.summary===undefined||("summary"in original&&source.summary===original.summary))&&original.paragraphs.length===source.paragraphs.length&&original.paragraphs.every((p,i)=>p===source.paragraphs[i]);
 const translate=(text:string)=>known?controlledPrivateLabel(language,text):text;
 return{id:source.id,version:source.version,title:translate(source.title),paragraphs:source.paragraphs.map(translate),summary:translate(source.summary??source.title),inlineLead:source.inlineLead??false,displayLanguage:known?language:"und"};
}
