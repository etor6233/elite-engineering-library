import {describe,it,expect} from 'vitest';
import catalogue from './business-functions.v403.json';
import bindings from './function-owner-bindings.v403.json';
import {functionOwnerBinding,functionOwnerDestination} from './function-owner-bindings';
describe('11 functions use existing owners without creating roles',()=>{
 it('binds exactly every existing task',()=>{expect(bindings.tasks.map(t=>t.id).sort()).toEqual(catalogue.functions.flatMap(f=>f.tasks.map(t=>t.id)).sort());expect(new Set(bindings.tasks.map(t=>t.function_id)).size).toBe(11);expect(bindings.tasks).toHaveLength(22)});
 it('never converts proposed permissions or event titles into grants',()=>{expect(bindings.grants).toEqual([]);for(const task of bindings.tasks){expect(task.grants).toEqual([]);expect(task.galaxy_method).toBe('WAITING_GALAXY_ADMISSION');if(task.runtime_claim==='PROJECT_RECORD_WORKFLOW')expect(task.required_product_permissions).toEqual([])}});
 it('routes product tasks to local owner surfaces and project tasks to records',()=>{for(const task of bindings.tasks){if(task.runtime_claim==='PROJECT_RECORD_WORKFLOW'){expect(functionOwnerDestination(task.id)).toBeNull();expect('records' in task.entry).toBe(true)}else{expect(functionOwnerDestination(task.id)).toMatch(/^\/(franchise)(?:[/?]|$)/);expect(task.owner_refs.length).toBeGreaterThan(0)}}});
 it('does not invent a destination for an unknown task',()=>{expect(functionOwnerBinding('unknown')).toBeUndefined();expect(functionOwnerDestination('unknown')).toBeNull()});
});
