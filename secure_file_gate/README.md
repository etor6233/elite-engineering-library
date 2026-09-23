# Secure Local File Gate

This pack is authored orchestration around three exact official Windows x86-64 release binaries: Google Magika CLI 1.1.0, Cisco Talos ClamAV 1.5.4 and VirusTotal/Google YARA-X 1.20.0. It does not copy or relabel their source. Acquire the pinned assets through `OFFICIAL-UPSTREAM-ACQUISITION-CORE` and retain their Apache-2.0, GPL-2.0-only and BSD-3-Clause terms.

The policy template intentionally starts at `AWAITING_USER`, has no allowed content type and contains an unresolved YARA rules hash. The agent must obtain a real owner, approval ID/date, exact label/MIME/score pairs, limits and an approved ruleset before execution. The runner refuses unresolved configuration.

Prepare mutable ClamAV signatures and compiled project rules:

```powershell
python prepare_security_assets.py update-clamav --freshclam <freshclam.exe> --clamscan <clamscan.exe> --database <database-directory>
python prepare_security_assets.py compile-yara --yara <yr.exe> --source <approved-rules> --output <new-rules.yarc>
```

Do not trust `freshclam` exit status by itself. The preparer also makes `clamscan` load the official databases, enforce maximum age and scan a benign probe before writing a hash inventory receipt. YARA source provenance, review, false-positive testing, canary and rollback remain project approvals; this library never supplies invented detection rules.

Run the gate only inside a sandboxed worker with denied egress, bounded CPU/memory/time and separate quarantine/clean storage:

```powershell
python secure_local_file_gate.py --input <quarantined-file> --output <new-evidence-directory> --policy <approved-policy.json> --magika-exe <magika.exe> --clamscan-exe <clamscan.exe> --clamav-database <database-directory> --yara-exe <yr.exe> --compiled-rules <approved-rules.yarc>
```

`ADMITTED` means only that the selected official engines passed the approved limits and rules at that moment. It is not a universal claim that the file is safe or semantically correct, and the receipt always keeps `business_storage_authorized=false`. Conversion and field extraction occur later, with an approved analyzer/schema/corpus and field-level evidence.
