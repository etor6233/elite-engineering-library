"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useRef, useState, type FormEvent } from "react";
import { answerSchema, resultSchema, type SurveyAnswer, type SurveyDefinition, type SurveySubmission } from "@/platform/surveys/contract";

export function CustomerSurvey({ organization, definition, initialAnswer }: { organization: string; definition: SurveyDefinition; initialAnswer: SurveyAnswer | null }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [answer, setAnswer] = useState(initialAnswer); const [score, setScore] = useState(""); const [consent, setConsent] = useState(false);
 const [busy, setBusy] = useState(false); const [uncertain, setUncertain] = useState(false); const [message, setMessage] = useState("");
 const pending = useRef(false); const attempted = useRef<SurveySubmission | null>(null);
 const url = "/api/enterprise/surveys?organizationId=" + encodeURIComponent(organization) + "&surveyId=" + encodeURIComponent(definition.id);
 async function submit(event: FormEvent) {
  event.preventDefault(); if (pending.current || answer || uncertain || !consent || score === "" || !definition.accepting) return;
  const input: SurveySubmission = { score: Number(score), consent: true, consent_version: definition.consent_version };
  attempted.current = input; pending.current = true; setBusy(true); setMessage("");
  try {
   const response = await fetch(url, { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(input), signal: AbortSignal.timeout(7000) });
   if (!response.ok) throw new Error("confirmation unavailable");
   const value = resultSchema.parse(await response.json());
   if (value.answer.score !== input.score || value.answer.consent_version !== input.consent_version) throw new Error("receipt mismatch");
   setAnswer(value.answer); setMessage(t("p0281"));
  } catch { setUncertain(true); setMessage(t("p0282")); }
  finally { pending.current = false; setBusy(false); }
 }
 async function recover() {
  if (pending.current) return; pending.current = true; setBusy(true);
  try {
   const response = await fetch(url + "&view=response", { cache: "no-store", signal: AbortSignal.timeout(7000) });
   if (response.status === 404) { setMessage(t("p0283")); return; }
   if (!response.ok) throw new Error("read unavailable");
   const value = answerSchema.parse(await response.json());
   if (value.consent_version !== definition.consent_version) throw new Error("version mismatch");
   setAnswer(value); setUncertain(false); setMessage(t("p0284"));
  } catch { setMessage(t("p0285")); }
  finally { pending.current = false; setBusy(false); }
 }
 return <section lang="es" aria-labelledby="survey-question">
  <h1 id="survey-question" className="pageTitle">{definition.prompt}</h1>
  {answer ? <div className="notice"><p>{t("p0286")} <strong data-testid="survey-stored-score">{answer.score}</strong> {t("p0287")}</p><p>{t("p0288")}</p></div> :
   <form onSubmit={submit} aria-describedby="survey-result">
    <label>{t("p0289")}<select name="score" required value={score} disabled={busy || uncertain || !definition.accepting} onChange={event => setScore(event.target.value)}>
     <option value="">{t("p0290")}</option>{Array.from({ length: 11 }, (_, n) => <option key={n} value={n}>{n}</option>)}
    </select></label>
    <p>{definition.consent_notice}</p>
    <label><input name="consent" type="checkbox" checked={consent} required disabled={busy || uncertain || !definition.accepting} onChange={event => setConsent(event.target.checked)} />{t("p0291")}</label>
    <button type="submit" disabled={busy || uncertain || !definition.accepting || score === "" || !consent}>{busy ? t("p0292") : t("p0293")}</button>
    {!definition.accepting ? <p>{t("p0294")}</p> : null}
   </form>}
  {uncertain && !answer ? <button type="button" onClick={recover} disabled={busy}>{t("p0295")}</button> : null}
  <p id="survey-result" role="status" aria-live="polite">{message}</p>
  <details><summary>{t("p0296")}</summary><p>{t("p0297")}</p><p>{t("p0298")}</p><p>{t("p0031")}</p></details>
 </section>;
}

