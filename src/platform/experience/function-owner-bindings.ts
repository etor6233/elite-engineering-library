// AUTHORED routing glue. Functions are topics, never authorization roles.
import bindings from './function-owner-bindings.v403.json';
export function functionOwnerBinding(taskId:string) { return bindings.tasks.find(task=>task.id===taskId); }
export function functionOwnerDestination(taskId:string):string|null {
 const entry=functionOwnerBinding(taskId)?.entry;
 return entry && 'route' in entry && typeof entry.route==='string' && entry.route.startsWith('/') && !entry.route.startsWith('//') ? entry.route : null;
}
