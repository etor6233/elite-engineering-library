from __future__ import annotations
import pathlib,tempfile,unittest
from validate_packaging import validate
ROOT=pathlib.Path(__file__).resolve().parents[1]
class PackagingTests(unittest.TestCase):
    def test_canonical_packaging(self)->None: validate(ROOT)
    def test_root_runtime_is_rejected(self)->None:
        with tempfile.TemporaryDirectory() as folder:
            target=pathlib.Path(folder)
            (target/"deploy").mkdir();(target/"packaging").mkdir()
            (target/"Dockerfile.api").write_text((ROOT/"Dockerfile.api").read_text().replace("USER nonroot:nonroot","USER 0"),encoding="utf-8")
            for name in ("compose.yaml","postgres-migrate.sh"):(target/"deploy"/name).write_text((ROOT/"deploy"/name).read_text(),encoding="utf-8")
            with self.assertRaisesRegex(ValueError,"non-root"):validate(target)
    def test_literal_secret_is_rejected(self)->None:
        with tempfile.TemporaryDirectory() as folder:
            target=pathlib.Path(folder);(target/"deploy").mkdir()
            (target/"Dockerfile.api").write_text((ROOT/"Dockerfile.api").read_text(),encoding="utf-8")
            (target/"deploy/postgres-migrate.sh").write_text((ROOT/"deploy/postgres-migrate.sh").read_text(),encoding="utf-8")
            compose=(ROOT/"deploy/compose.yaml").read_text().replace("POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?set POSTGRES_PASSWORD outside source control}","POSTGRES_PASSWORD: hardcoded")
            (target/"deploy/compose.yaml").write_text(compose,encoding="utf-8")
            with self.assertRaisesRegex(ValueError,"literal database secret"):validate(target)
if __name__=="__main__":unittest.main()
