import {readFileSync} from "node:fs";
import {expect,it} from "vitest";
import {catalogCanonical,catalogSHA,catalogCommand,commandPayload,receiptMatches} from "./authoring";
const cases=JSON.parse(readFileSync("deploy/catalog/authoring-goldens.json","utf8")) as {name:string;payload:Record<string,unknown>;canonical:string;sha256:string}[];
it.each(cases)("matches typed Go $name hashes without narrowing int64 values",async c=>{
 const body=c.name==="media"?c.payload:commandPayload(catalogCommand.parse({action:c.name,organization_id:"org",...c.payload}));
 expect(catalogCanonical(body)).toBe(c.canonical);expect(await catalogSHA(catalogCanonical(body))).toBe(c.sha256);
});
it("normalizes source instants to the same Go wire representation",()=>{
 const source=cases.find(c=>c.name==="source")!;
 const c=catalogCommand.parse({action:"source",organization_id:"org",...source.payload,valid_from:"2026-09-12T00:00:00.000Z",valid_until:"2026-10-12T00:00:00.120Z"});
 expect(catalogCanonical(commandPayload(c))).toBe(source.canonical);
});
it("never accepts a receipt for another actor, payload or draft",()=>{
 const p={command_id:"command",kind:"publish" as const,request_sha256:"a".repeat(64),resource_id:"draft",snapshot_sha256:"b".repeat(64)};
 const receipt={...p,actor:"reviewer",generation:"1",replay:false,effective_price_book_id:"book"};
 expect(receiptMatches(receipt,p,"reviewer")).toBe(true);
 for(const change of [{actor:"other"},{request_sha256:"c".repeat(64)},{resource_id:"other"},{snapshot_sha256:"c".repeat(64)}])expect(receiptMatches({...receipt,...change},p,"reviewer")).toBe(false);
});
it("rejects unsupported transport shape and invalid source before a write",()=>{
 const source=cases.find(c=>c.name==="source")!;
 const input={action:"source",organization_id:"org",...source.payload};
 expect(catalogCommand.safeParse({...input,extra:true}).success).toBe(false);
 expect(()=>catalogCanonical({amount:42})).toThrow();
 expect(()=>catalogCanonical({"non-ascii-é":"value"})).toThrow();
});
