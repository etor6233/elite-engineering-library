"""AUTHORED source-oracle projection; not an AL runtime or source admission."""
from pathlib import Path
import argparse,ast,hashlib,json,re

def main():
    parser=argparse.ArgumentParser()
    parser.add_argument('--project-root',type=Path,required=True)
    parser.add_argument('--output',type=Path,required=True)
    args=parser.parse_args()
    source=args.project_root/'docs/provenance/source/ServiceItemLine.Table.al'
    with (args.project_root/'docs/provenance/BC_WARRANTY_DERIVATION.json').open('rb') as stream:metadata=stream.read(32769)
    if len(metadata)>32768:raise SystemExit('Manifest exceeds bound')
    manifest=json.loads(metadata)
    with source.open('rb') as stream:raw=stream.read(2097153)
    if len(raw)>2097152:raise SystemExit('Source exceeds bound')
    if hashlib.sha256(raw).hexdigest()!=manifest['source']['source_file_sha256']:
        raise SystemExit('Source hash mismatch')
    start,end=manifest['source']['source_lines']
    lines=raw.splitlines(keepends=True);segment=b''.join(lines[start-1:end])
    if hashlib.sha256(segment).hexdigest()!=manifest['source']['segment_sha256']:
        raise SystemExit('Source segment mismatch')
    text=segment.decode('utf-8')
    expressions=[re.search(r'WarrantyParts := (.*);',text)[1],re.search(r'WarrantyLabor := (.*);',text)[1],re.search(r'if ("Warranty Starting Date \(Parts\)" > "Warranty Ending Date \(Parts\)") then',text)[1]]
    fields={'"Warranty Starting Date (Parts)"':'parts_start','"Warranty Ending Date (Parts)"':'parts_end','"Warranty Starting Date (Labor)"':'labor_start','"Warranty Ending Date (Labor)"':'labor_end','Date':'date'}
    compiled=[]
    for expression in expressions:
        for key,value in fields.items():expression=expression.replace(key,value)
        tree=ast.parse(expression,mode='eval')
        allowed=(ast.Expression,ast.BoolOp,ast.And,ast.Compare,ast.Name,ast.Load,ast.Gt,ast.GtE,ast.LtE)
        if not all(isinstance(node,allowed) for node in ast.walk(tree)) or not {n.id for n in ast.walk(tree) if isinstance(n,ast.Name)}<=set(fields.values()):
            raise SystemExit('Unsupported source expression')
        compiled.append(compile(tree,'fixed-official-AL-date-expression','eval'))
    a,b,c=compiled;vectors=[]
    for date in range(5):
        for ps in range(4):
            for pe in range(4):
                for ls in range(4):
                    for le in range(4):
                        env={'date':date,'parts_start':ps,'parts_end':pe,'labor_start':ls,'labor_end':le}
                        bad=bool(eval(c,{'__builtins__':{}},env));parts=False if bad else bool(eval(a,{'__builtins__':{}},env));labor=False if bad else bool(eval(b,{'__builtins__':{}},env))
                        vectors.append(dict(env,error=bad,parts=parts,labor=labor,any=parts or labor))
    output=(json.dumps(vectors,separators=(',',':'))+'\n').encode('utf-8')
    # Exclusive creation; never replace source, checked-in fixtures or receipts.
    with args.output.open('xb') as stream:stream.write(output)
    print(json.dumps({'status':'SOURCE_PROJECTION_ONLY','vectors':len(vectors),'sha256':hashlib.sha256(output).hexdigest(),'source_sha256':hashlib.sha256(raw).hexdigest()}))

if __name__=='__main__':main()
