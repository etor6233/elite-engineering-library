import{it,expect}from"vitest";
import{createElement as h}from"react";
import{renderToStaticMarkup}from"react-dom/server";
import{createHash}from"node:crypto";
import{readFileSync}from"node:fs";
import{privateMessage,controlledPrivateLabel}from"./private-catalog";
import messages from"./private-messages.json";
import{guideDisplay}from"./private-guide-display";
import{trainingCourseDisplay}from"./private-training-display";
import{ALL_GUIDES}from"@/platform/help/content";
import binding from"./private-training-binding.json";
import{courseView}from"@/platform/training/contract";
import{PrivateLocaleProvider}from"./private-provider";
import{RoleDashboard}from"@/components/role-dashboard";
import{TrainingWorkspace}from"@/components/training-workspace";
import{PrivateReadFailure}from"@/platform/backend/private-read-failure";
import{PortalPaging,portalPageHref}from"@/platform/backend/portal-paging";
import{minorAmountPresentation}from"./money";
const en={language:"en",locale:"en-US",timeZone:"America/Argentina/Buenos_Aires",source:"preference"}as const;
const render=(child:React.ReactNode)=>renderToStaticMarkup(h(PrivateLocaleProvider,{locale:en,children:child}));
function course(index:number){const row=binding.courses[index];if(!row)throw new Error("missing fixture course");return courseView.parse({course:row.source,articles:ALL_GUIDES.filter(g=>row.source.lessons.includes(g.id)).map(({id,version,title,paragraphs})=>({id,version,title,paragraphs})),profile_id:binding.profile_id,profile_revision:binding.profile_revision,profile_sha256:binding.profile_sha256,content_sha256:binding.content_sha256,method:"HUMAN_REVIEW_NO_GRANTS_V1"})}
it("retains all template parameters as literal text without HTML or recursive interpolation",()=>{
 const supplied="<img src=x onerror=alert(1)> {handover_state}";
 expect(privateMessage("en","p0629",{handover_id:supplied,handover_state:"accepted"})).toBe("Preparation retrieved: "+supplied+". Status accepted.");
 expect(privateMessage("es","p0629",{handover_id:"H1",handover_state:"accepted"})).toBe("Preparación recuperada: H1. Estado accepted.");
 expect(controlledPrivateLabel("en","Tenant-authored article")).toBe("Tenant-authored article");
 for(const value of Object.values(messages)){expect(value.es.trim()).not.toBe("");expect(value.en.trim()).not.toBe("");expect(value.en.match(/\{[a-z_]+\}/g)??[]).toEqual(value.es.match(/\{[a-z_]+\}/g)??[])}
});
it("translates every exact same-release guide while preserving the source and refusing modified or old content",()=>{
 const original=JSON.stringify(ALL_GUIDES);expect(ALL_GUIDES).toHaveLength(21);
 for(const guide of ALL_GUIDES){if(typeof guide.version!=="string")throw new Error("missing fixed guide version");const source={...guide,version:guide.version};const translated=guideDisplay(source,"en");expect(translated.displayLanguage).toBe("en");expect(translated.title).not.toBe(guide.title);expect(translated.paragraphs.every((p,i)=>p!==guide.paragraphs[i])).toBe(true);expect(guideDisplay(source,"es").paragraphs).toEqual(guide.paragraphs)}
 const known=ALL_GUIDES[0];expect(guideDisplay({...known,version:"0.0.1"},"en").title).toBe(known.title);
 const changed={...known,paragraphs:["Aceptar una cotización"]};expect(guideDisplay(changed,"en").paragraphs).toEqual(changed.paragraphs);expect(guideDisplay(changed,"en").displayLanguage).toBe("und");
 expect(guideDisplay({...known,summary:"Inventario"},"en").summary).toBe("Inventario");
 expect(JSON.stringify(ALL_GUIDES)).toBe(original);
});
it("binds five translated curricula to exact profile and course content without rewriting assessment evidence",()=>{
 const raw=readFileSync("deploy/training/reference.profile.json");expect(createHash("sha256").update(raw).digest("hex")).toBe(binding.profile_sha256);
 expect(createHash("sha256").update(readFileSync("training_content/help.bundle.json")).digest("hex")).toBe(binding.content_sha256);
 for(let i=0;i<5;i++){const value=course(i),before=JSON.stringify(value),view=trainingCourseDisplay(value,"en");expect(view.title).toBe(binding.courses[i]!.en.title);expect(view.prompts).toEqual(binding.courses[i]!.en.prompts);expect(JSON.stringify(value)).toBe(before);
  for(const altered of[{...value,profile_sha256:"0".repeat(64)},{...value,content_sha256:"0".repeat(64)},{...value,profile_revision:1},{...value,course:{...value.course,title:"Mi evaluación"}},{...value,course:{...value.course,prompts:[{id:"review",text:"Texto cambiado"}]}}])expect(trainingCourseDisplay(altered,"en").title).toBe(altered.course.title);
 }
});
it("renders private role links, training and error recovery in English without changing role grants or paging URLs",()=>{
 const sections=[{id:"guide",label:"Guía de uso",href:"/guide"},{id:"supply",label:"Suministro",href:"/supply"}],before=JSON.stringify(sections);
 const dashboard=render(h(RoleDashboard,{sections}));expect(dashboard).toContain("Available access");expect(dashboard).toContain('href="/supply"');expect(dashboard).not.toContain("Suministro");expect(JSON.stringify(sections)).toBe(before);
 const training=render(h(TrainingWorkspace,{courses:[course(0)],assessments:[],canLearn:true,canReview:false,scope:"tenant:user"}));expect(training).toContain("Register resources and recover a response");expect(training).toContain("Start practice");expect(training).not.toContain("Iniciar práctica");
 const failure=render(h(PrivateReadFailure,{path:"/customer"}));expect(failure).toContain("We could not load this information");expect(failure).toContain('href="/customer"');
 const href=portalPageHref("/customer",{orders_after:"opaque-reference"});const paging=render(h(PortalPaging,{label:"Orders",nextHref:href,firstHref:"/customer"}));expect(paging).toContain("View next page");expect(paging).toContain('href="/customer?orders_after=opaque-reference"');
});
it("formats safe minor units exactly in both languages without changing amounts or validity",()=>{
 const amount=9007199254740991;
 expect(minorAmountPresentation(amount,"ARS","en-US")).toEqual({amountLabel:"ARS 90,071,992,547,409.91",amountValid:true});
 expect(minorAmountPresentation(amount,"ARS","es-AR").amountLabel).toContain("90.071.992.547.409,91");
 expect(minorAmountPresentation(NaN,"ARS","en-US")).toEqual({amountLabel:"Amount cannot be verified; contact support.",amountValid:false});
 expect(minorAmountPresentation(100,"JPY","en-US").amountLabel).toBe("JPY 100");
});
