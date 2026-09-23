import base64,hashlib,hmac,json,tempfile,time,unittest
from pathlib import Path
from test_whatsapp_cloud import profile
from whatsapp_cloud import build_reply_payload,normalize_verified_webhook,send_bridge,recover_send_bridge

class ConversationContractTests(unittest.TestCase):
    def request(self):
        start=int(time.time())-10
        return {"kind":"text_reply","recipient":"5491112345678","text":"Texto exacto <revisado>","source_message_id":"wamid.fixture","last_inbound_at":start,"window_expires_at":start+86400}
    def test_exact_text_payload_window_and_utf8(self):
        request=self.request()
        recipient,payload=build_reply_payload(profile(),request)
        self.assertEqual(recipient,request["recipient"])
        self.assertEqual(payload,{"messaging_product":"whatsapp","recipient_type":"individual","to":recipient,"type":"text","text":{"body":request["text"]}})
        for change in [{"text":"\ud800"},{"window_expires_at":request["window_expires_at"]+1},{"last_inbound_at":1,"window_expires_at":86401},{"template_name":"invented"}]:
            with self.subTest(change=list(change)),self.assertRaises((ValueError,PermissionError,UnicodeError)):
                build_reply_payload(profile(),{**request,**change})
    def test_recovery_requires_bound_original_evidence_without_transport(self):
        request=self.request()
        with tempfile.TemporaryDirectory() as tmp:
            folder=Path(tmp)/"send"
            frame={"schema":"elite-whatsapp-send-bridge/v1","binding_sha256":"a"*64,"profile":profile(),"request":request,"access_token":"synthetic","output_directory":str(folder)}
            sent=send_bridge(json.dumps(frame).encode(),lambda *args:(200,{},b'{"messages":[{"id":"wamid.original"}]}'))
            original={"schema":"elite-whatsapp-recover-send/v1","binding_sha256":"a"*64,"profile":profile(),"request":request,"receipt":base64.b64encode((folder/"SEND_RECEIPT.json").read_bytes()).decode(),"response":base64.b64encode((folder/"provider-response.json").read_bytes()).decode()}
            self.assertEqual(recover_send_bridge(json.dumps(original).encode()),sent)
            changed={**original,"request":{**request,"text":"otro texto"}}
            with self.assertRaises(ValueError):recover_send_bridge(json.dumps(changed).encode())
            changed={**original,"response":base64.b64encode(b'{"messages":[{"id":"wamid.forged"}]}').decode()}
            with self.assertRaises(ValueError):recover_send_bridge(json.dumps(changed).encode())

if __name__=="__main__":unittest.main()
