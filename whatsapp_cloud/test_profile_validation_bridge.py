import base64, hashlib, json, subprocess, sys, unittest
from pathlib import Path
from test_whatsapp_cloud import profile

class ProfileValidationBridgeTests(unittest.TestCase):
    def test_blocked_default_and_exact_fixture(self):
        root=Path(__file__).resolve().parent
        for good, raw in [(False,(root/"provider-profile.template.json").read_bytes()),(True,(json.dumps(profile(),ensure_ascii=False,indent=2)+"\n").encode("utf-8"))]:
            with self.subTest(accepted=good):
                frame=json.dumps({"schema":"elite-whatsapp-validate-profile/v1","profile":base64.b64encode(raw).decode()}).encode()
                result=subprocess.run([sys.executable,"-I","-B",str(root/"whatsapp_cloud.py"),"--validate-profile-bridge"],input=frame,capture_output=True)
                if good:
                    self.assertEqual(result.returncode,0)
                    self.assertEqual(json.loads(result.stdout),{"schema":"elite-whatsapp-profile-validation/v1","profile_sha256":hashlib.sha256(raw).hexdigest()})
                else:
                    self.assertNotEqual(result.returncode,0)
                    self.assertEqual(result.stdout,b"")

if __name__=="__main__":unittest.main()
