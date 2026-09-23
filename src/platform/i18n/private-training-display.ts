// AUTHORED source-bound display only; assessment source and hashes stay unchanged.
import type{PrivateLanguage}from"./private-locale";
import binding from"./private-training-binding.json";
import type{CourseView}from"@/platform/training/contract";
export function trainingCourseDisplay(source:CourseView,language:PrivateLanguage){
 const record=binding.courses.find(x=>x.source.id===source.course.id);
 const same=record&&record.source.title===source.course.title&&record.source.role===source.course.role&&JSON.stringify(record.source.lessons)===JSON.stringify(source.course.lessons)&&record.source.prompts.length===source.course.prompts.length&&record.source.prompts.every((p,i)=>p.id===source.course.prompts[i]?.id&&p.text===source.course.prompts[i]?.text);
 const known=source.profile_sha256===binding.profile_sha256&&source.content_sha256===binding.content_sha256&&source.profile_id===binding.profile_id&&source.profile_revision===binding.profile_revision&&source.method==="HUMAN_REVIEW_NO_GRANTS_V1"&&same;
 return known&&language==="en"?{title:record.en.title,prompts:record.en.prompts,displayLanguage:"en"}:{title:source.course.title,prompts:source.course.prompts,displayLanguage:known?"es":"und"};
}
