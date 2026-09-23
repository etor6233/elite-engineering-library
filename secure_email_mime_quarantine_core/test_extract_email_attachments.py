import base64
import json
import tempfile
import unittest
from email.message import EmailMessage
from pathlib import Path

from extract_email_attachments import IntakeRejected, Limits, extract, sha256


ALLOWED = {"application/pdf", "text/csv", "image/jpeg"}


def email_bytes(items=(("invoice.pdf", "application", "pdf", b"PDF"),)):
    msg = EmailMessage()
    msg["From"] = "sender@example.com"
    msg["To"] = "intake@example.com"
    msg["Subject"] = "documents"
    msg.set_content("attached")
    for name, main, sub, data in items:
        msg.add_attachment(data, maintype=main, subtype=sub, filename=name, cte="base64")
    return msg.as_bytes()


def receipt(raw):
    return {
        "receipt_version": 1,
        "decision": "ADMIT",
        "tenant_id": "tenant-a",
        "raw_sha256": sha256(raw),
        "source_object": {"bucket": "raw-bucket", "key": "mail/1", "version_id": "v1", "sequencer": "001"},
        "retained_original": {"uri": "s3://retained/mail/1", "version_id": "locked-v1"},
    }


class ExtractTests(unittest.TestCase):
    def run_extract(self, raw=None, rec=None, limits=Limits(), allowed=ALLOWED):
        raw = raw or email_bytes()
        rec = rec or receipt(raw)
        root = tempfile.TemporaryDirectory()
        self.addCleanup(root.cleanup)
        out = Path(root.name) / "out"
        result = extract(raw, rec, out, allowed, limits)
        return result, out

    def rejected(self, code, raw=None, rec=None, limits=Limits(), allowed=ALLOWED):
        with self.assertRaises(IntakeRejected) as ctx:
            self.run_extract(raw, rec, limits, allowed)
        self.assertEqual(code, ctx.exception.code)

    def test_positive_is_hash_named_atomic_and_blocks_storage(self):
        result, out = self.run_extract()
        item = result["attachments"][0]
        self.assertEqual(b"PDF", (out / item["quarantine_name"]).read_bytes())
        self.assertFalse(result["automatic_storage_authorized"])
        self.assertEqual("PENDING", item["security_decision"])
        self.assertTrue((out / "manifest.json").is_file())

    def test_filename_traversal_is_never_a_path(self):
        raw = email_bytes((("../../secret.pdf", "application", "pdf", b"PDF"),))
        result, out = self.run_extract(raw, receipt(raw))
        self.assertNotIn("secret", result["attachments"][0]["quarantine_name"])
        self.assertEqual(2, len(list(out.iterdir())))

    def test_multiple_same_names_do_not_collide(self):
        raw = email_bytes((("same.pdf", "application", "pdf", b"A"), ("same.pdf", "application", "pdf", b"B")))
        result, _ = self.run_extract(raw, receipt(raw))
        self.assertEqual(2, len({x["quarantine_name"] for x in result["attachments"]}))

    def test_manifest_is_deterministic(self):
        raw = email_bytes()
        a, _ = self.run_extract(raw, receipt(raw))
        b, _ = self.run_extract(raw, receipt(raw))
        self.assertEqual(a, b)

    def test_receipt_must_admit(self):
        raw = email_bytes(); rec = receipt(raw); rec["decision"] = "REJECT"
        self.rejected("RECEIPT_NOT_ADMITTED", raw, rec)

    def test_raw_hash_must_match(self):
        raw = email_bytes(); rec = receipt(raw); rec["raw_sha256"] = "0" * 64
        self.rejected("RAW_SHA256_MISMATCH", raw, rec)

    def test_source_version_is_required(self):
        raw = email_bytes(); rec = receipt(raw); rec["source_object"]["version_id"] = ""
        self.rejected("INVALID_RECEIPT_VERSION_ID", raw, rec)

    def test_retained_original_is_required(self):
        raw = email_bytes(); rec = receipt(raw); rec.pop("retained_original")
        self.rejected("INVALID_RECEIPT_RETAINED_ORIGINAL", raw, rec)

    def test_raw_size_limit(self):
        raw = email_bytes()
        self.rejected("RAW_SIZE_LIMIT", raw, receipt(raw), Limits(max_raw_bytes=10))

    def test_attachment_count_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"A"), ("b.pdf", "application", "pdf", b"B")))
        self.rejected("ATTACHMENT_COUNT_LIMIT", raw, receipt(raw), Limits(max_attachments=1))

    def test_attachment_size_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"AB"),))
        self.rejected("ATTACHMENT_SIZE_LIMIT", raw, receipt(raw), Limits(max_attachment_bytes=1))

    def test_total_size_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"AA"), ("b.pdf", "application", "pdf", b"BB")))
        self.rejected("ATTACHMENT_TOTAL_SIZE_LIMIT", raw, receipt(raw), Limits(max_total_attachment_bytes=3))

    def test_content_type_allowlist(self):
        raw = email_bytes((("a.exe", "application", "octet-stream", b"MZ"),))
        self.rejected("ATTACHMENT_CONTENT_TYPE_REJECTED", raw, receipt(raw))

    def test_non_base64_rejected(self):
        msg = EmailMessage(); msg.set_content("body"); msg.add_attachment("plain", subtype="csv", filename="a.csv", cte="quoted-printable")
        raw = msg.as_bytes()
        self.rejected("ATTACHMENT_TRANSFER_ENCODING_NOT_BASE64", raw, receipt(raw))

    def test_no_attachments_rejected(self):
        msg = EmailMessage(); msg.set_content("body"); raw = msg.as_bytes()
        self.rejected("NO_ATTACHMENTS", raw, receipt(raw))

    def test_existing_output_rejected(self):
        raw = email_bytes(); rec = receipt(raw)
        with tempfile.TemporaryDirectory() as td:
            out = Path(td) / "out"; out.mkdir()
            with self.assertRaises(IntakeRejected) as ctx: extract(raw, rec, out, ALLOWED)
            self.assertEqual("OUTPUT_ALREADY_EXISTS", ctx.exception.code)

    def test_invalid_tenant_rejected(self):
        raw = email_bytes(); rec = receipt(raw); rec["tenant_id"] = "../tenant"
        self.rejected("INVALID_RECEIPT_TENANT_ID", raw, rec)

    def test_part_limit(self):
        raw = email_bytes((("a.pdf", "application", "pdf", b"A"),))
        self.rejected("MIME_PART_LIMIT", raw, receipt(raw), Limits(max_parts=1))


if __name__ == "__main__":
    unittest.main(verbosity=2)
