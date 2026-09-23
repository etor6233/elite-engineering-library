# AUTHORED isolated launcher and exact source-manifest binding.
import hashlib
import json
from pathlib import Path
import sys
from types import ModuleType


def run():
    root = Path(__file__).resolve().parent
    lock_path = root/'engine-lock.json'
    raw = lock_path.read_bytes()
    if len(sys.argv) != 2 or len(raw) > 65536 or hashlib.sha256(raw).hexdigest() != sys.argv[1]:
        return 2
    lock = json.loads(raw)
    if lock['schema'] != 'elite.odoo-derived-source-lock.v1':
        return 2
    sources = {}
    for record in lock['files']:
        relative = Path(record['path'])
        if relative.is_absolute() or '..' in relative.parts:
            return 2
        path = (root/relative).resolve()
        if not path.is_relative_to(root) or path.is_symlink():
            return 2
        data = path.read_bytes()
        if len(data) != record['bytes'] or hashlib.sha256(data).hexdigest() != record['sha256']:
            return 2
        if record['path'] in sources:
            return 2
        sources[record['path']] = data
    # Do not use filesystem package/bytecode import after checking source. A
    # new __init__.py, .pyc or changed file must not replace verified bytes.
    package = ModuleType('odoo_loyalty')
    package.__path__ = []
    package.__package__ = 'odoo_loyalty'
    sys.modules['odoo_loyalty'] = package
    for name in ('engine', 'orm_contract', 'protocol'):
        relative = name + '.py'
        source = sources[relative]
        module = ModuleType('odoo_loyalty.' + name)
        module.__package__ = 'odoo_loyalty'
        module.__file__ = str(root/relative)
        sys.modules[module.__name__] = module
        setattr(package, name, module)
        exec(compile(source, module.__file__, 'exec', dont_inherit=True), module.__dict__)
    return package.protocol.main()


if __name__ == '__main__':
    try:
        raise SystemExit(run())
    except (OSError,ValueError,KeyError,TypeError):
        raise SystemExit(2)
