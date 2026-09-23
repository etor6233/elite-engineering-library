"use client";
import {experienceText} from "@/platform/experience/messages";
import {guideDisplay} from "@/platform/i18n/private-guide-display";
import {usePrivateI18n} from "@/platform/i18n/private-provider";
// Pure shared rendering: no effects, session state, operation or training progress.
type Guide = { id: string; version: string; title: string; summary: string; paragraphs: readonly string[]; inlineLead: boolean };
export function OperationalGuide({ guide, className, summaryLabel }: { guide: Guide; className?: string; summaryLabel?:string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const display=guideDisplay(guide,privateLocale.language);
  return <details lang={display.displayLanguage} className={className}>
    <summary>{summaryLabel?experienceText(summaryLabel,privateLocale.language):display.summary}</summary>
    <p>{guide.id}{t("p0280")}{guide.version}{guide.inlineLead ? " · " + display.paragraphs[0] : ""}</p>
    {display.paragraphs.slice(guide.inlineLead ? 1 : 0).map(text => <p key={text}>{text}</p>)}
    <a href={`/help?article=${encodeURIComponent(guide.id)}&version=${encodeURIComponent(guide.version)}`}>{t("p0143")}</a>
  </details>;
}
