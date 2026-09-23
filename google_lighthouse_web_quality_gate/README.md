# Google Lighthouse Web Quality Gate

This isolated module installs Google Lighthouse `13.4.1` and `chrome-launcher` `1.2.1` from the exact pnpm lock. It runs the four official Lighthouse categories five times against one public URL, requires no warnings or origin escape, then enforces both median and minimum score floors.

## Use

```powershell
pnpm install --ignore-workspace --frozen-lockfile
pwsh -NoProfile -File ./verify_contract.ps1
pwsh -NoProfile -File ./verify_contract.ps1 -RunTarget -TargetUrl http://127.0.0.1:4181/ -ChromePath <exact-chromium-executable> -OutputDirectory ./lighthouse-results
```

The default policy requires performance median `>=0.90` and minimum `>=0.85`, accessibility `1.00`, best practices `>=0.90`, and SEO `1.00`. All five raw reports and one summary receipt remain under the explicit output directory.

On Windows, when the component path would push Lighthouse's nested CommonJS package boundary beyond the legacy path limit, the verifier copies only its exact manifests, policy, runner and validator into a uniquely named short temporary capsule, performs a frozen offline pnpm install, and executes there without changing Node's module semantics. The capsule is validated as a child of the system temporary directory and removed with bounded retries in `finally`; unsafe or incomplete cleanup fails closed.

This gate is automated lab evidence only. It does not prove assistive-technology interoperability, real-user performance, load capacity, offensive security, authenticated roles, deployment, canary, rollback, or business acceptance. Remote targets require HTTPS; only loopback may use HTTP. URLs containing credentials or fragments fail closed.

## Frozen offline installation with pnpm 11.25.0

Use the exact admitted pnpm artifact and observed Node runtime from the project tool lock. An offline store needs both package bytes and the registry metadata used by supply-chain policy checks. ERR_PNPM_NO_OFFLINE_META is a failed gate even if node_modules was linked. In an authorized public-registry preparation step, run the same frozen install online to populate metadata, then repeat --offline --frozen-lockfile --ignore-workspace. Preserve the lock hash and both receipts; never disable trust or release-age policy to force PASS. V321 verified this sequence with unchanged dependency lockfiles.
