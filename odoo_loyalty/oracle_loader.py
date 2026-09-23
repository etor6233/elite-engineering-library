# AUTHORED test-only loading glue. Exact source functions retain Odoo LGPL attribution.
import ast
import hashlib
import importlib.util
from pathlib import Path
from types import SimpleNamespace
from collections import defaultdict

ROOT = Path(__file__).parent
_float_path = ROOT/'upstream/odoo/tools/float_utils.py'
_spec = importlib.util.spec_from_file_location('odoo_original_float_utils', _float_path)
_float_module = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(_float_module)
original_float_round = _float_module.float_round


def original_functions():
    path = ROOT/'upstream/addons/sale_loyalty/models/sale_order.py'
    tree = ast.parse(path.read_bytes())
    names = {'_get_point_changes', '_get_real_points_for_coupon', '_program_check_compute_points'}
    nodes = [node for node in ast.walk(tree) if isinstance(node, ast.FunctionDef) and node.name in names]
    if len(nodes) != len(names):
        raise ValueError('source oracle missing functions')
    scope = {'defaultdict':defaultdict,'float_round':original_float_round,
             '_':lambda text, **kwargs: text % kwargs if kwargs else text}
    exec(compile(ast.Module(body=nodes,type_ignores=[]), str(path), 'exec'), scope)
    return SimpleNamespace(**{name:scope[name] for name in names})
