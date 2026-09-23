"""AUTHORED independent Decimal oracle; creates only an absent output file."""
from decimal import Decimal, localcontext, ROUND_HALF_UP
from pathlib import Path
import argparse
import json

parser = argparse.ArgumentParser()
parser.add_argument("--output", type=Path, required=True)
args = parser.parse_args()
cases = [(str(n), str(d), dec, q)
         for n in [-123456789, -103, -101, -25, -5, -1, 0, 1, 5, 25, 101, 103, 123456789]
         for d in [1, 2, 4, 5, 8, 10, 20, 100]
         for dec, q in [(0, 1), (2, 1), (2, 5), (3, 25), (9, 1)]]
cases += [(n, d, 0, 1) for n, d in [
    ("9223372036854775807", "1"), ("-9223372036854775808", "1"),
    ("18446744073709551615", "2"), ("-18446744073709551617", "2"),
    ("-123456789", "100000")]]
vectors = []
with localcontext() as context:
    context.prec = 200
    for n, d, decimals, quantum in cases:
        units = Decimal(n) / Decimal(d) * Decimal(10) ** decimals / Decimal(quantum)
        result = int(units.quantize(Decimal("1"), rounding=ROUND_HALF_UP) * Decimal(quantum))
        vectors.append({"numerator": n, "denominator": d, "decimals": decimals,
                        "precision": str(quantum), "want": str(result),
                        "overflow": not (-(2 ** 63) <= result < 2 ** 63)})
document = {"oracle": "Python decimal 200 digits, ROUND_HALF_UP, independent of Go math/big quotient/remainder",
            "vectors": vectors}
with args.output.open("x", encoding="utf-8", newline="\n") as stream:
    stream.write(json.dumps(document, indent=2) + "\n")
