"use client";
import {useEffect,useState} from "react";
import {publicText as t} from "./messages";
import s from "./public-web.module.css";
// Missing published media never falls back to another product or to a concept.
export function PublicMedia({src,alt,width,height,locale,priority=false}:{src:string;alt:string;width:number;height:number;locale:string;priority?:boolean}){
  const [failed,setFailed]=useState(false);
  useEffect(()=>{setFailed(false);},[src]);
  if(failed)return<span className={s.noMedia} role="img" aria-label={`${alt}: ${t(locale,"noMedia")}`} data-media-state="error">{t(locale,"noMedia")}</span>;
  return<img src={src} alt={alt} width={width} height={height} loading={priority?"eager":"lazy"} fetchPriority={priority?"high":"auto"} onError={()=>setFailed(true)}/>;
}
