# Google Lighthouse Web Quality Gate

## 1. Metadata

```yaml
pack_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE"
pack_version: "0.1.5"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa Google Lighthouse 13.4.1 oficial con lock exacto, ejecuta cinco auditorías secuenciales sobre una URL pública y falla cerrado por versión, warnings, origen, número de auditorías, mínimos y medianas de rendimiento, accesibilidad automatizada, buenas prácticas y SEO."
stacks: ["Google Lighthouse 13.4.1", "chrome-launcher 1.2.1", "Chromium 151.0.7922.34 rev 1234", "Node.js >=22.19", "pnpm 11.25.0", "PowerShell 7"]
compatible_with: ["MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.x", "TYPESCRIPT-GO-API-WEB-BRIDGE 0.2.x", "MARKDOWN-COMPOSITOR-CORE 0.2.x"]
incompatible_with: ["credentials in URLs", "remote HTTP targets", "single-run performance claims", "human accessibility claims from automation", "production claims from a local lab"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/GoogleChrome/lighthouse/tree/1d58f5b06d28e3419b38817a6c7488ec4413c67d", "https://registry.npmjs.org/lighthouse/-/lighthouse-13.4.1.tgz"]
verified_at: "2026-09-08"
```

## 2. Applicability

Use after a production build exists and before accepting a public web surface. The pack provides an isolated, pinned Google runtime plus deterministic local policy glue; it does not modify the web application's package manifest or lock. Adopt it when the project needs fast regression evidence for Lighthouse performance, automated accessibility, best practices and SEO.

Reject it as proof of assistive-technology interoperability, real-user performance, load capacity, security, authenticated roles or production readiness. Remote targets require HTTPS; only loopback may use HTTP. URLs with credentials or fragments fail before a browser starts.

## 3. Architecture contract

The exact Lighthouse npm graph remains in `google_lighthouse_web_quality_gate`. `run-quality-gate.mjs` launches an explicitly selected Chromium through Google's `chrome-launcher`, using a unique profile directory, and invokes the official Lighthouse API five times sequentially. Raw reports use exclusive creation; the validator requires the same target origin, exact runtime version, zero warnings, at least 140 audits, all four categories and configurable minimum/median floors. Cleanup is part of PASS: browser termination or profile removal failure makes the command fail.

The default lab policy requires performance median `>=0.90` and minimum `>=0.85`, accessibility `1.00`, best-practices `>=0.90`, and SEO `1.00`. Five sequential runs follow Google's documented variability guidance. A project may raise thresholds through a reviewed policy change, but may not lower or omit them silently. Reports can contain page data and stay in the explicit project evidence directory; targets containing credentials are prohibited.

Failure modes are non-zero exit with an exact reason: dependency/runtime drift, invalid target, browser absence, origin escape, warning, incomplete report, weak score, output collision or cleanup failure. Rollback removes this isolated directory and its profile entry from the composition plan; it changes no application runtime or database.

On Windows, a deeply composed path can place Lighthouse's nested CommonJS package boundary beyond the legacy path limit and make Node misclassify upstream modules. The verifier detects the risky component-root length, copies only its exact manifests/policy/runner/validator into a short uniquely named temporary capsule, performs a frozen offline pnpm install, executes without altering Node's module semantics, and removes the validated capsule with bounded retries in `finally`. The workaround changes neither upstream identities nor reports and fails closed if cleanup cannot be proven.

## 4. Exact file manifest

```text
CREATE google_lighthouse_web_quality_gate/package.json
CREATE google_lighthouse_web_quality_gate/pnpm-lock.yaml
CREATE google_lighthouse_web_quality_gate/source-lock.json
CREATE google_lighthouse_web_quality_gate/LICENSE.lighthouse.txt
CREATE google_lighthouse_web_quality_gate/quality-policy.json
CREATE google_lighthouse_web_quality_gate/lib/validate-reports.mjs
CREATE google_lighthouse_web_quality_gate/run-quality-gate.mjs
CREATE google_lighthouse_web_quality_gate/tests/validate-reports.test.mjs
CREATE google_lighthouse_web_quality_gate/verify_contract.ps1
CREATE google_lighthouse_web_quality_gate/README.md
```

## 5. Materialization blocks

### FILE: `google_lighthouse_web_quality_gate/package.json`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:package-json:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite isolated manifest for exact Google dependencies"
license: "LicenseRef-Workspace-Owner"
sha256: "dcdb7fd7c346f73543643070f14c15578a3f425245b24a5398b314d26366b7bb"
variables: []
secrets_allowed: false
```
````json
{
  "name": "elite-google-lighthouse-web-quality-gate",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "packageManager": "pnpm@11.25.0",
  "engines": {"node": ">=22.19"},
  "scripts": {
    "test": "node --test tests/*.test.mjs",
    "gate": "node run-quality-gate.mjs"
  },
  "dependencies": {
    "chrome-launcher": "1.2.1",
    "lighthouse": "13.4.1"
  }
}
````

### FILE: `google_lighthouse_web_quality_gate/pnpm-lock.yaml`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:pnpm-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "pnpm 11.19.0 exact resolution of Lighthouse 13.4.1 and chrome-launcher 1.2.1"
license: "LicenseRef-Workspace-Owner"
sha256: "c690da3a314e5275d68c4ea4a88e700dbd9a58b0da6fcbc25aec2bdc486dd9aa"
variables: []
secrets_allowed: false
```
````yaml
lockfileVersion: '9.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    dependencies:
      chrome-launcher:
        specifier: 1.2.1
        version: 1.2.1
      lighthouse:
        specifier: 13.4.1
        version: 13.4.1(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))

packages:

  '@apm-js-collab/code-transformer-bundler-plugins@0.7.4':
    resolution: {integrity: sha512-nAfOeZPSUAQvJa1iFT/5oCrTm5YQhMMrfCNthNnaXHZiOQhu1KGuLoIx7HtbAi3wfwaBYLaICPIeenIaEwcXIg==}
    engines: {node: '>=18.0.0'}

  '@apm-js-collab/code-transformer@0.18.1':
    resolution: {integrity: sha512-u1Hb6bHjWtkSpiprwVP6YaHC1DTN4RAU3zYkUDUe7WMnJwdyU1pwTL9dFKiSJB9IiLue/EQovmyx6xhU7FFtAQ==}
    hasBin: true

  '@apm-js-collab/tracing-hooks@0.13.0':
    resolution: {integrity: sha512-mTvWz9rnQwx1U3h0XPTHaX7bgfkpipLLTQyjlC2cdhQpQEuoLT0AGzoydeoq2NxfEVv6fWOOETcSbb2nptleyw==}

  '@formatjs/ecma402-abstract@2.3.6':
    resolution: {integrity: sha512-HJnTFeRM2kVFVr5gr5kH1XP6K0JcJtE7Lzvtr3FS/so5f1kpsqqqxy5JF+FRaO6H2qmcMfAUIox7AJteieRtVw==}

  '@formatjs/fast-memoize@2.2.7':
    resolution: {integrity: sha512-Yabmi9nSvyOMrlSeGGWDiH7rf3a7sIwplbvo/dlz9WCIjzIQAfy1RMf4S0X3yG724n5Ghu2GmEl5NJIV6O9sZQ==}

  '@formatjs/icu-messageformat-parser@2.11.4':
    resolution: {integrity: sha512-7kR78cRrPNB4fjGFZg3Rmj5aah8rQj9KPzuLsmcSn4ipLXQvC04keycTI1F7kJYDwIXtT2+7IDEto842CfZBtw==}

  '@formatjs/icu-skeleton-parser@1.8.16':
    resolution: {integrity: sha512-H13E9Xl+PxBd8D5/6TVUluSpxGNvFSlN/b3coUp0e0JpuWXXnQDiavIpY3NnvSp4xhEMoXyyBvVfdFX8jglOHQ==}

  '@formatjs/intl-localematcher@0.6.2':
    resolution: {integrity: sha512-XOMO2Hupl0wdd172Y06h6kLpBz6Dv+J4okPLl4LPtzbr8f66WbIoy4ev98EBuZ6ZK4h5ydTN6XneT4QVpD7cdA==}

  '@jridgewell/sourcemap-codec@1.6.0':
    resolution: {integrity: sha512-T7jf+5zgsZHwNJ4lvQ7/aezbyk0nNX+zJVWpmHA7VYsEx7a7qr5Rg5IbtJFqkgze5Y2sruq1RUY8Q837Od7iFw==}

  '@opentelemetry/api-logs@0.220.0':
    resolution: {integrity: sha512-CmVa4ImJ+ynfrPMNaAXHET6Bhb44SwzmfyVJFq9ni2jgXJR/l7C6gfVFddNmHP+ZOkP9cf4f9DBe68qVLTHc9w==}
    engines: {node: '>=8.0.0'}

  '@opentelemetry/api@1.9.1':
    resolution: {integrity: sha512-gLyJlPHPZYdAk1JENA9LeHejZe1Ti77/pTeFm/nMXmQH/HFZlcS/O2XJB+L8fkbrNSqhdtlvjBVjxwUYanNH5Q==}
    engines: {node: '>=8.0.0'}

  '@opentelemetry/core@2.10.0':
    resolution: {integrity: sha512-/wNZ8twnEQQA4HoHu22+vcsdru6pWPWxW+7w+FlxT6Id7PE/WIbZmVKkte+PF72e0F2dnImFeHD2syyE1Mw6MQ==}
    engines: {node: ^18.19.0 || >=20.6.0}
    peerDependencies:
      '@opentelemetry/api': '>=1.0.0 <1.10.0'

  '@opentelemetry/instrumentation@0.220.0':
    resolution: {integrity: sha512-xQx3E2WxP1mDvKzxLxX+CTCtNLa560YJZ3087qYHerl2YmiKpv7AH+dAy7vmx+eVrZ5BwhfWUAVoKOoxCNHcpw==}
    engines: {node: ^18.19.0 || >=20.6.0}
    peerDependencies:
      '@opentelemetry/api': ^1.3.0

  '@opentelemetry/resources@2.10.0':
    resolution: {integrity: sha512-q6MMm2zhggzsHVNbabYwut+a6nbuQQe3URUoxaojM/8K1IBfwwPzvxIjNi2/lI1TFe+fMHMW9MWhrtDLEXEnkA==}
    engines: {node: ^18.19.0 || >=20.6.0}
    peerDependencies:
      '@opentelemetry/api': '>=1.3.0 <1.10.0'

  '@opentelemetry/sdk-trace-base@2.10.0':
    resolution: {integrity: sha512-GuYQQT7QD2EeO8lcZLRQzcbOyhqAzL+6WWTKTU9mSUBYBazkEDl+VrQcXQhbB08OWM9anD1aHleVadzulpOaUQ==}
    engines: {node: ^18.19.0 || >=20.6.0}
    peerDependencies:
      '@opentelemetry/api': '>=1.3.0 <1.10.0'

  '@opentelemetry/sdk-trace@2.10.0':
    resolution: {integrity: sha512-MfQGq3GRmTh5fM/y+OjaO0vj6+luCB1XO2gfXCalKCfgKw0eHL++sm75DNweC6ohlp+aFvACqeE0fYayqdRaoQ==}
    engines: {node: ^18.19.0 || >=20.6.0}
    peerDependencies:
      '@opentelemetry/api': '>=1.3.0 <1.10.0'

  '@opentelemetry/semantic-conventions@1.43.0':
    resolution: {integrity: sha512-eSYWTm620tTk45EKSedaUL8MFYI8hW164hIXsgIHyxu3VobUB3fFCu5t0hQby6OoWRPsG1KkKUG2M5UadiLiVg==}
    engines: {node: '>=14'}

  '@paulirish/trace_engine@0.0.65':
    resolution: {integrity: sha512-Qsm6F5C8xf6ZzQXbQc2+wcpe6sggfs/gvc/ytqSurdvYg3kyW0ECHCqE0CWBKZpqgjVfPNX9c7SCS3r2nEIRGg==}

  '@puppeteer/browsers@3.2.1':
    resolution: {integrity: sha512-KDz+3qDRdBAlRlMjmKyj6dEs33YHTk/xRHEENSXq6TNnhgoU15ruSHtEBeVF6OZ9tBDY55Se4P0nFMNsipzU9A==}
    engines: {node: '>=22.12.0'}
    hasBin: true
    peerDependencies:
      proxy-agent: '>=8.0.1'
      yauzl: ^2.10.0 || ^3.4.0
    peerDependenciesMeta:
      proxy-agent:
        optional: true
      yauzl:
        optional: true

  '@sentry/conventions@0.16.0':
    resolution: {integrity: sha512-fO9PLmHdVURcSPUpWCItWAtgKiMwGdJHbovoSEyLplX5sxs2ugvI4CBPTrkkgqhObnZOD0CnWBKDzSVQYBKEyQ==}
    engines: {node: '>=14'}

  '@sentry/core@10.71.0':
    resolution: {integrity: sha512-OIjT7rzcWJjUC6r3eBT3Td1j0afDBMkbbx9jTocSD+ZSfc25eEU7hoIPS0WvfeIOTIN3y8bfQnXavwMReaNVHQ==}
    engines: {node: '>=18'}

  '@sentry/node-core@10.71.0':
    resolution: {integrity: sha512-sxd0/ZW+Uda/17H0R7lB2Othm37VYcdwKdWdyHRizg872TfCX4UwTTPaoS2tMJAUjV2tVh83ZA0fsoC2q9SpjA==}
    engines: {node: '>=18'}
    peerDependencies:
      '@opentelemetry/api': ^1.9.0
      '@opentelemetry/core': ^1.30.1 || ^2.1.0
      '@opentelemetry/exporter-trace-otlp-http': '>=0.57.0 <1'
      '@opentelemetry/instrumentation': '>=0.57.1 <1'
      '@opentelemetry/sdk-trace-base': ^1.30.1 || ^2.1.0
    peerDependenciesMeta:
      '@opentelemetry/api':
        optional: true
      '@opentelemetry/core':
        optional: true
      '@opentelemetry/exporter-trace-otlp-http':
        optional: true
      '@opentelemetry/instrumentation':
        optional: true
      '@opentelemetry/sdk-trace-base':
        optional: true

  '@sentry/node@10.71.0':
    resolution: {integrity: sha512-bw2M/xkMu2+ATo6QWFmtTZTYp5LV1krt9/DTtYqtt4GmhXmXgdFPXCv6783AJI3HUoEMx2BcehhNbKzGrFXo/g==}
    engines: {node: '>=18'}

  '@sentry/opentelemetry@10.71.0':
    resolution: {integrity: sha512-YgeL0xTObKma3MuOrt+/6M/f6mo/Z08LHh3OxPomzQpgpCgCLGyJ/739cDSSH1LRjqWdPQtoia5asl5ZRdhWmw==}
    engines: {node: '>=18'}
    peerDependencies:
      '@opentelemetry/api': ^1.9.0
      '@opentelemetry/core': ^1.30.1 || ^2.1.0
      '@opentelemetry/sdk-trace-base': ^1.30.1 || ^2.1.0

  '@sentry/server-utils@10.71.0':
    resolution: {integrity: sha512-zdyShKNsghzPGVWVRzc7oybYTsRZdKgdPtq+dsksImo1b6n7ILlD7eYx1Q0eAOZoWkDqDmm+Ml6MQb76IT1eTA==}
    engines: {node: '>=18'}

  '@types/estree@1.0.9':
    resolution: {integrity: sha512-GhdPgy1el4/ImP05X05Uw4cw2/M93BCUmnEvWZNStlCzEKME4Fkk+YpoA5OiHNQmoS7Cafb8Xa3Pya8m1Qrzeg==}

  '@types/node@26.4.0':
    resolution: {integrity: sha512-faiGnoIrLH/V8cibOMEAZ8pMw6oXqSukl29ra4mN8GdaB2ZewzeaLj+INpV5N+Z1eKWzY+IzaIZH2EIR6YZRNQ==}

  ansi-colors@4.1.3:
    resolution: {integrity: sha512-/6w/C21Pm1A7aZitlI5Ni/2J6FFQN8i1Cvz3kHABAAbw93v/NlvKdVOqz7CCWz/3iv/JplRSEEZ83XION15ovw==}
    engines: {node: '>=6'}

  ansi-regex@5.0.1:
    resolution: {integrity: sha512-quJQXlTSUGL2LH9SUXo8VwsY4soanhgo6LNSm84E1LBcE8s3O0wpdiRzyR9z/ZZJMlMWv37qOOb9pdJlMUEKFQ==}
    engines: {node: '>=8'}

  ansi-regex@6.3.0:
    resolution: {integrity: sha512-WpDfL7NO6j7tH88IDBNVdUJxDh9nmCteAVW9dsep846XdwF4naCBK+/tGLX3KJgcpgMRXCFlTM2hKGoK9FsdrQ==}
    engines: {node: '>=12'}

  ansi-styles@4.3.0:
    resolution: {integrity: sha512-zbB9rCJAT1rbjiVDb2hqKFHNYLxgtk8NURxZ3IZwD3F6NtxbXZQCnnSi1Lkx+IDohdPlFp222wVALIheZJQSEg==}
    engines: {node: '>=8'}

  ansi-styles@6.2.3:
    resolution: {integrity: sha512-4Dj6M28JB+oAH8kFkTLUo+a2jwOFkuqb3yucU0CANcRRUbxS0cP0nZYCGjcc3BNXwRIsUVmDGgzawme7zvJHvg==}
    engines: {node: '>=12'}

  astring@1.9.0:
    resolution: {integrity: sha512-LElXdjswlqjWrPpJFg1Fx4wpkOCxj1TDHlSV4PlaRxHGWko024xICaa97ZkMfs6DRKlCguiAI+rbXv5GWwXIkg==}
    hasBin: true

  atomically@2.1.1:
    resolution: {integrity: sha512-P4w9o2dqARji6P7MHprklbfiArZAWvo07yW7qs3pdljb3BWr12FIB7W+p0zJiuiVsUpRO0iZn1kFFcpPegg0tQ==}

  axe-core@4.13.0:
    resolution: {integrity: sha512-UzGt8zg7Ny8djbYMhxl2zuEevVa7r2gJjYY5Lwr1xM7+XU2nd6CkIWFTVcCIbAP63vSz71NaVyyuSk9lHKcy0A==}
    engines: {node: '>=4'}

  chrome-launcher@1.2.1:
    resolution: {integrity: sha512-qmFR5PLMzHyuNJHwOloHPAHhbaNglkfeV/xDtt5b7xiFFyU1I+AZZX0PYseMuhenJSSirgxELYIbswcoc+5H4A==}
    engines: {node: '>=12.13.0'}
    hasBin: true

  chromium-bidi@17.0.2:
    resolution: {integrity: sha512-5v9GQFhTktFvotn/OFNJBmKLKRAb6n9r0bVCwf7sHgWc3/JryK0bj1nn93L3pHFrfgcsu6Be6EWsDi+1XHTGDg==}
    engines: {node: '>=20.19.0 <22.0.0 || >=22.12.0'}
    peerDependencies:
      devtools-protocol: '*'

  cjs-module-lexer@2.2.1:
    resolution: {integrity: sha512-Ca8swihM+/4yKecYHY52kgJd300hi2lADU/a1RxNTRe+RJ9jvqQlESpbz9DnG9mowez8qwXHB8qYdIUw9e+F5Q==}

  cliui@8.0.1:
    resolution: {integrity: sha512-BSeNnyus75C4//NQ9gQt1/csTXyo/8Sb+afLAkzAptFuMsod9HFokGNudZpi/oQV73hnVK+sR+5PVRMd+Dr7YQ==}
    engines: {node: '>=12'}

  cliui@9.0.1:
    resolution: {integrity: sha512-k7ndgKhwoQveBL+/1tqGJYNz097I7WOvwbmmU2AR5+magtbjPWQTS1C5vzGkBC8Ym8UWRzfKUzUUqFLypY4Q+w==}
    engines: {node: '>=20'}

  color-convert@2.0.1:
    resolution: {integrity: sha512-RRECPsj7iu/xb5oKYcsFHSppFNnsj/52OVTRKb4zP5onXwVF3zVmmToNcOfGC+CRDpfK/U584fMg38ZHCaElKQ==}
    engines: {node: '>=7.0.0'}

  color-name@1.1.4:
    resolution: {integrity: sha512-dOy+3AuW3a2wNbZHIuMZpTcgjGuLU/uBL/ubcZF9OXbDo8ff4O8yVp5Bf0efS8uEoYo5q4Fx7dY9OgQGXgAsQA==}

  configstore@7.1.0:
    resolution: {integrity: sha512-N4oog6YJWbR9kGyXvS7jEykLDXIE2C0ILYqNBZBp9iwiJpoCBWYsuAdW6PPFn6w06jjnC+3JstVvWHO4cZqvRg==}
    engines: {node: '>=18'}

  csp_evaluator@1.1.8:
    resolution: {integrity: sha512-EwOnfYuNbTytvbMKsLixTrRgnjOa0WZCxGy8A9nnSYAicrdwn+T/epU/yjgymmOxlgKnvH+8wXt+7p/8ak5Feg==}

  debug@4.4.3:
    resolution: {integrity: sha512-RGwwWnwQvkVfavKVt22FGLw+xYSdzARwm0ru6DhTVA3umU5hZc28V3kO4stgYryrTlLpuvgI9GiijltAjNbcqA==}
    engines: {node: '>=6.0'}
    peerDependencies:
      supports-color: '*'
    peerDependenciesMeta:
      supports-color:
        optional: true

  decimal.js@10.6.0:
    resolution: {integrity: sha512-YpgQiITW3JXGntzdUmyUR1V812Hn8T1YVXhCu+wO3OpS4eU9l4YdD3qjyiKdV6mvV29zapkMeD390UVEf2lkUg==}

  define-lazy-prop@2.0.0:
    resolution: {integrity: sha512-Ds09qNh8yw3khSjiJjiUInaGX9xlqZDY7JVryGxdxV7NPeuqQfplOpQ66yJFZut3jLa5zOwkXw1g9EI2uKh4Og==}
    engines: {node: '>=8'}

  devtools-protocol@0.0.1663043:
    resolution: {integrity: sha512-33aOY3ZnBP1dgZsshgaL+/XlsQleiFZgyUaDtdZkEa1nbZhVY1MoDeWjk+wxg25fU924l1ZJfoGNmjjeA/5s1w==}

  devtools-protocol@0.0.1666840:
    resolution: {integrity: sha512-gCcO42XCHKEs7Ag0S7aGYsnJ7hlgrO3qderYqeiY0Eqk+0GFfuvT13IA0hHreJTa2KCdDVyGMeOhdMNmrrTjVg==}

  dot-prop@9.0.0:
    resolution: {integrity: sha512-1gxPBJpI/pcjQhKgIU91II6Wkay+dLcN3M6rf2uwP8hRur3HtQXjVrdAK3sjC0piaEuxzMwjXChcETiJl47lAQ==}
    engines: {node: '>=18'}

  emoji-regex@10.6.0:
    resolution: {integrity: sha512-toUI84YS5YmxW219erniWD0CIVOo46xGKColeNQRgOzDorgBi1v4D71/OFzgD9GO2UGKIv1C3Sp8DAn0+j5w7A==}

  emoji-regex@8.0.0:
    resolution: {integrity: sha512-MSjYzcWNOA0ewAHpz0MxpYFvwg6yjy1NG3xteoqz644VCo/RPgnr1/GGt+ic3iJTzQ8Eu3TdM14SawnVUmGE6A==}

  enquirer@2.4.1:
    resolution: {integrity: sha512-rRqJg/6gd538VHvR3PSrdRBb/1Vy2YfzHqzvbhGIQpDRKIa4FgV/54b5Q1xYSxOOwKvjXweS26E0Q+nAMwp2pQ==}
    engines: {node: '>=8.6'}

  es-module-lexer@2.3.2:
    resolution: {integrity: sha512-poHGpORABojJJucnV9KbOavETW8lBVnphkW77ER5/BQ5Fz7oXSoCNek7IH3vR5nRjdsEz926ibFYX8KtLQmdyw==}

  escalade@3.2.0:
    resolution: {integrity: sha512-WUj2qlxaQtO4g6Pq5c29GTcWGDyd8itL8zTlipgECz3JesAiiOKotd8JU6otB3PACgG6xkJUyVhboMS+bje/jA==}
    engines: {node: '>=6'}

  escape-string-regexp@4.0.0:
    resolution: {integrity: sha512-TtpcNJ3XAzx3Gq8sWRzJaVajRs0uVxA2YAkdb1jm2YkPz4G6egUFAyA3n5vtEIZefPk5Wa4UXbKuS5fKkJWdgA==}
    engines: {node: '>=10'}

  esquery@1.7.0:
    resolution: {integrity: sha512-Ap6G0WQwcU/LHsvLwON1fAQX9Zp0A2Y6Y/cJBl9r/JbW90Zyg4/zbG6zzKa2OTALELarYHmKu0GhpM5EO+7T0g==}
    engines: {node: '>=0.10'}

  estraverse@5.3.0:
    resolution: {integrity: sha512-MMdARuVEQziNTeJD8DgMqmhwR11BRQ/cBP+pLtYdSTnf3MIO8fFeiINEbX36ZdNlfU/7A9f3gUw49B3oQsvwBA==}
    engines: {node: '>=4.0'}

  get-caller-file@2.0.5:
    resolution: {integrity: sha512-DyFP3BM/3YHTQOCUL/w0OZHR0lpKeGrxotcHWcqNEdnltqFwXVfhEBQ94eIo34AfQpo0rGki4cyIiftY06h2Fg==}
    engines: {node: 6.* || 8.* || >= 10.*}

  get-east-asian-width@1.6.0:
    resolution: {integrity: sha512-QRbvDIbx6YklUe6RxeTeleMR0yv3cYH6PsPZHcnVn7xv7zO1BHN8r0XETu8n6Ye3Q+ahtSarc3WgtNWmehIBfA==}
    engines: {node: '>=18'}

  graceful-fs@4.2.11:
    resolution: {integrity: sha512-RbJ5/jmFcNNCcDV5o9eTnBLJ/HszWV0P73bc+Ff4nS/rJj+YaS6IGyiOL0VoBYX+l1Wrl3k63h/KrH+nhJ0XvQ==}

  http-link-header@1.1.4:
    resolution: {integrity: sha512-xT3GPW6/ZbGuw4UvwHqErSCEjNUlwbQJuZn9/q5U4WEKfp2kENVCAlousG1zLxHeaQ/ffOHUNpWamvkbBW0eNw==}
    engines: {node: '>=6.0.0'}

  image-ssim@0.2.0:
    resolution: {integrity: sha512-W7+sO6/yhxy83L0G7xR8YAc5Z5QFtYEXXRV6EaE8tuYBZJnA3gVgp3q7X7muhLZVodeb9UfvjSbwt9VJwjIYAg==}

  import-in-the-middle@3.3.3:
    resolution: {integrity: sha512-AiohS3H80sXO6owEltjGX+glb7qXaDhBoJb9XcQVH4UI207xu/bDLUcadVKp7Qe576reg9yr/PXZjV5qx8gfbA==}
    engines: {node: '>=18'}

  intl-messageformat@10.7.18:
    resolution: {integrity: sha512-m3Ofv/X/tV8Y3tHXLohcuVuhWKo7BBq62cqY15etqmLxg2DZ34AGGgQDeR+SCta2+zICb1NX83af0GJmbQ1++g==}

  is-docker@2.2.1:
    resolution: {integrity: sha512-F+i2BKsFrH66iaUFc0woD8sLy8getkwTwtOBjvs56Cx4CgJDeKQeqfz8wAYiSb8JOprWhHH5p77PbmYCvvUuXQ==}
    engines: {node: '>=8'}
    hasBin: true

  is-fullwidth-code-point@3.0.0:
    resolution: {integrity: sha512-zymm5+u+sCsSWyD9qNaejV3DFvhCKclKdizYaJUuHA83RLjb7nSuGnddCHGv0hk+KY7BMAlsWeK4Ueg6EV6XQg==}
    engines: {node: '>=8'}

  is-wsl@2.2.0:
    resolution: {integrity: sha512-fKzAra0rGJUUBwGBgNkHZuToZcn+TtXHpeCgmkMJMMYx1sQDYaCSyjJBSCa2nH1DGm7s3n1oBnohoVTBaN7Lww==}
    engines: {node: '>=8'}

  jpeg-js@0.4.4:
    resolution: {integrity: sha512-WZzeDOEtTOBK4Mdsar0IqEU5sMr3vSV2RqkAIzUEV2BHnUfKGyswWFPFwK5EeDo93K3FohSHbLAjj0s1Wzd+dg==}

  js-library-detector@6.7.0:
    resolution: {integrity: sha512-c80Qupofp43y4cJ7+8TTDN/AsDwLi5oOm/plBrWI+iQt485vKXCco+yVmOwEgdo9VOdsYTuV0UlTeetVPTriXA==}
    engines: {node: '>=12'}

  legacy-javascript@0.0.1:
    resolution: {integrity: sha512-lPyntS4/aS7jpuvOlitZDFifBCb4W8L/3QU0PLbUTUj+zYah8rfVjYic88yG7ZKTxhS5h9iz7duT8oUXKszLhg==}

  lighthouse-logger@2.0.2:
    resolution: {integrity: sha512-vWl2+u5jgOQuZR55Z1WM0XDdrJT6mzMP8zHUct7xTlWhuQs+eV0g+QL0RQdFjT54zVmbhLCP8vIVpy1wGn/gCg==}

  lighthouse-stack-packs@1.12.3:
    resolution: {integrity: sha512-d8IsOpE83kbANgnM+Tp8+x6HcMpX9o2ITBiUERssgzAIFdZCQzs/f4k6D0DLQTE59enml9mbAOU52Wu35exWtg==}

  lighthouse@13.4.1:
    resolution: {integrity: sha512-fDu8lt3QLK/lTqIxtp1HkzQNJ32rsFHhbadYOepcMZFLgA8oINhxutMbMv8XXnpTOvZ0TXCo4JCk1LDTWaRLnA==}
    engines: {node: '>=22.19'}
    hasBin: true

  lodash-es@4.18.1:
    resolution: {integrity: sha512-J8xewKD/Gk22OZbhpOVSwcs60zhd95ESDwezOFuA3/099925PdHJ7OFHNTGtajL3AlZkykD32HykiMo+BIBI8A==}

  lookup-closest-locale@6.2.0:
    resolution: {integrity: sha512-/c2kL+Vnp1jnV6K6RpDTHK3dgg0Tu2VVp+elEiJpjfS1UyY7AjOYHohRug6wT0OpoX2qFgNORndE9RqesfVxWQ==}

  magic-string@0.30.21:
    resolution: {integrity: sha512-vd2F4YUyEXKGcLHoq+TEyCjxueSeHnFxyyjNp80yg0XV4vUhnDer/lvvlqM/arB5bXQN5K2/3oinyCRyx8T2CQ==}

  marky@1.3.0:
    resolution: {integrity: sha512-ocnPZQLNpvbedwTy9kNrQEsknEfgvcLMvOtz3sFeWApDq1MXH1TqkCIx58xlpESsfwQOnuBO9beyQuNGzVvuhQ==}

  meriyah@6.1.4:
    resolution: {integrity: sha512-Sz8FzjzI0kN13GK/6MVEsVzMZEPvOhnmmI1lU5+/1cGOiK3QUahntrNNtdVeihrO7t9JpoH75iMNXg6R6uWflQ==}
    engines: {node: '>=18.0.0'}

  mitt@3.0.1:
    resolution: {integrity: sha512-vKivATfr97l2/QBCYAkXYDbrIWPM2IIKEl7YPhjCvKlG3kE2gm+uBo6nEXK3M5/Ffh/FLpKExzOQ3JJoJGFKBw==}

  modern-tar@0.8.4:
    resolution: {integrity: sha512-gN54ddmyzEg10orwZ2u4OOv+bjpMWdIl5jIkodK97bMq8QBSL5c0D7YX0lT1Ooz+99S7+PvFbnxzdjgHo1r41g==}
    engines: {node: '>=18.0.0'}

  module-details-from-path@1.0.4:
    resolution: {integrity: sha512-EGWKgxALGMgzvxYF1UyGTy0HXX/2vHLkw6+NvDKW2jypWbHpjQuj4UMcqQWXHERJhVGKikolT06G3bcKe4fi7w==}

  ms@2.1.3:
    resolution: {integrity: sha512-6FlzubTLZG3J2a/NVCAleEhjzq5oxgHyaCU9yYXvcLsvoVaHJq/s5xXI6/XXP6tz7R9xAOtHnSO/tXtF3WRTlA==}

  open@8.4.2:
    resolution: {integrity: sha512-7x81NCL719oNbsq/3mh+hVrAWmFuEYUqrq/Iw3kUzH8ReypT9QQ0BLoJS7/G9k6N81XjW4qHWtjWwe/9eLy1EQ==}
    engines: {node: '>=12'}

  puppeteer-core@25.9.0:
    resolution: {integrity: sha512-U61rCwSMha62CA/Opy6tCx2Fx+ck7ouiKnbpEApzSoLYMoEu9F71nuFpHL55vmIt33/GYm6eKZVhH2ev0nAIeg==}
    engines: {node: '>=22.12.0'}

  require-directory@2.1.1:
    resolution: {integrity: sha512-fGxEI7+wsG9xrvdjsrlmL22OMTTiHRwAMroiEeMgq8gzoLC/PQr7RsRDSTLUg/bZAZtF+TVIkHc6/4RIKrui+Q==}
    engines: {node: '>=0.10.0'}

  require-in-the-middle@8.0.1:
    resolution: {integrity: sha512-QT7FVMXfWOYFbeRBF6nu+I6tr2Tf3u0q8RIEjNob/heKY/nh7drD/k7eeMFmSQgnTtCzLDcCu/XEnpW2wk4xCQ==}
    engines: {node: '>=9.3.0 || >=8.10.0 <9.0.0'}

  robots-parser@3.0.1:
    resolution: {integrity: sha512-s+pyvQeIKIZ0dx5iJiQk1tPLJAWln39+MI5jtM8wnyws+G5azk+dMnMX0qfbqNetKKNgcWWOdi0sfm+FbQbgdQ==}
    engines: {node: '>=10.0.0'}

  semifies@1.0.0:
    resolution: {integrity: sha512-xXR3KGeoxTNWPD4aBvL5NUpMTT7WMANr3EWnaS190QVkY52lqqcVRD7Q05UVbBhiWDGWMlJEUam9m7uFFGVScw==}

  source-map@0.6.1:
    resolution: {integrity: sha512-UjgapumWlbMhkBgzT7Ykc5YXUT46F0iKu8SGXq0bcwP5dz/h0Plj6enJqjz1Zbq2l5WaqYnrVbwWOWMyF3F47g==}
    engines: {node: '>=0.10.0'}

  speedline-core@1.4.3:
    resolution: {integrity: sha512-DI7/OuAUD+GMpR6dmu8lliO2Wg5zfeh+/xsdyJZCzd8o5JgFUjCeLsBDuZjIQJdwXS3J0L/uZYrELKYqx+PXog==}
    engines: {node: '>=8.0'}

  string-width@4.2.3:
    resolution: {integrity: sha512-wKyQRQpjJ0sIp62ErSZdGsjMJWsap5oRNihHhu6G7JVO/9jIB6UyevL+tXuOqrng8j/cxKTWyWUwvSTriiZz/g==}
    engines: {node: '>=8'}

  string-width@7.2.0:
    resolution: {integrity: sha512-tsaTIkKW9b4N+AEj+SVA+WhJzV7/zMhcSu78mLKWSk7cXMOSHsBKFWUs0fWwq8QyK3MgJBQRX6Gbi4kYbdvGkQ==}
    engines: {node: '>=18'}

  string-width@8.2.2:
    resolution: {integrity: sha512-GaPUh5gfdrYzqeVNZvUfT23vYYxXzKYidUcnMtJg/3rxRV63EFZy3k6xfKlmfeJD0176lnUV/Usr3XcwSvFzpg==}
    engines: {node: '>=20'}

  strip-ansi@6.0.1:
    resolution: {integrity: sha512-Y38VPSHcqkFrCpFnQ9vuSXmquuv5oXOKpGeT6aGrr3o3Gc9AlVa6JBfUSOCnbxGGZF+/0ooI7KrPuUSztUdU5A==}
    engines: {node: '>=8'}

  strip-ansi@7.2.0:
    resolution: {integrity: sha512-yDPMNjp4WyfYBkHnjIRLfca1i6KMyGCtsVgoKe/z1+6vukgaENdgGBZt+ZmKPc4gavvEZ5OgHfHdrazhgNyG7w==}
    engines: {node: '>=12'}

  stubborn-fs@2.0.0:
    resolution: {integrity: sha512-Y0AvSwDw8y+nlSNFXMm2g6L51rBGdAQT20J3YSOqxC53Lo3bjWRtr2BKcfYoAf352WYpsZSTURrA0tqhfgudPA==}

  stubborn-utils@1.0.2:
    resolution: {integrity: sha512-zOh9jPYI+xrNOyisSelgym4tolKTJCQd5GBhK0+0xJvcYDcwlOoxF/rnFKQ2KRZknXSG9jWAp66fwP6AxN9STg==}

  third-party-web@0.29.2:
    resolution: {integrity: sha512-fegtha91tq2DHphyoiBXVHjVi2YG9zFaRnboT9C28tO1en9Y3wJsfspuy40F+u5wl3hHVbw7cnd1b67kEGHb8g==}

  tldts-core@7.4.11:
    resolution: {integrity: sha512-CW3WN2rIIE/Of21mulhgnGOwoDyEFNygyIBOONSdyAuSATgMMUCpLeUlB+E8sAwA5xRV9hYPl+kyZ9citHCaKg==}

  tldts-icann@7.4.11:
    resolution: {integrity: sha512-1p+NDJ7FUYCliESmsQl9EW5Um8JIyFheiy6Y6ZoHho27d+TLAGFa0gMf1VOVy02vNY1aQo6xeKiaGj45+7P+PA==}

  tslib@2.8.1:
    resolution: {integrity: sha512-oJFu94HQb+KVduSUQL7wnpmqnfmLsOA/nAh6b6EH0wCEoK0/mPeXU6c3wKDV83MkOuHPRHtSXKKU99IBazS/2w==}

  type-fest@4.41.0:
    resolution: {integrity: sha512-TeTSQ6H5YHvpqVwBRcnLDCBnDOHWYu7IvGbHT6N8AOymcr9PJGjc1GTtiWZTYg0NCgYwvnYWEkVChQAr9bjfwA==}
    engines: {node: '>=16'}

  typed-query-selector@2.12.2:
    resolution: {integrity: sha512-EOPFbyIub4ngnEdqi2yOcNeDLaX/0jcE1JoAXQDDMIthap7FoN795lc/SHfIq2d416VufXpM8z/lD+WRm2gfOQ==}

  undici-types@8.3.0:
    resolution: {integrity: sha512-j375ScV60dom+YkPFIfTLcOiPxkN/buHz5GobjLhixFuANaNs3C9l4GmrWqejgXWJ7BbJcFYpTEUkS1Ge8bpZQ==}

  web-features@3.36.0:
    resolution: {integrity: sha512-B5W9HbyXT76soRG7Fn3AELZkuDQpmpSxszSmJZLIPUHYxIGkufNz0vcoEhfqnUF1HAfaducWIBaaKuusq9Am2g==}

  webdriver-bidi-protocol@0.4.2:
    resolution: {integrity: sha512-VSV+fzfChirL3e7jay2yUC7B4HQCGtEWEg/MSSQbK+qWbqeGlRLlXTzPpYr3XGUvbpDHumWZBJxgesg4N7dbtA==}

  when-exit@2.1.5:
    resolution: {integrity: sha512-VGkKJ564kzt6Ms1dbgPP/yuIoQCrsFAnRbptpC5wOEsDaNsbCB2bnfnaA8i/vRs5tjUSEOtIuvl9/MyVsvQZCg==}

  wrap-ansi@7.0.0:
    resolution: {integrity: sha512-YVGIj2kamLSTxw6NsZjoBxfSwsn0ycdesmc4p+Q21c5zPuZ1pl+NfxVdxPtdHvmNVOQ6XSYG4AUtyt/Fi7D16Q==}
    engines: {node: '>=10'}

  wrap-ansi@9.0.2:
    resolution: {integrity: sha512-42AtmgqjV+X1VpdOfyTGOYRi0/zsoLqtXQckTmqTeybT+BDIbM/Guxo7x3pE2vtpr1ok6xRqM9OpBe+Jyoqyww==}
    engines: {node: '>=18'}

  ws@7.5.13:
    resolution: {integrity: sha512-rsKI6xDBFVf4r/x8XyChGK04QR/XHroxs/jUcoWvtEZM8TPU/X/uIY9B1CsSzYws9ZJb/6bbBu7dPhFW00CAoA==}
    engines: {node: '>=8.3.0'}
    peerDependencies:
      bufferutil: ^4.0.1
      utf-8-validate: ^5.0.2
    peerDependenciesMeta:
      bufferutil:
        optional: true
      utf-8-validate:
        optional: true

  ws@8.21.3:
    resolution: {integrity: sha512-201TZ/kPWxoPr/OKWjquZR1SWKXcvxdH+e1xrx89b3YbmzLMFCLfnaG1HFIgWzJOEWZ7MvpK++odZufgYR50Rw==}
    engines: {node: '>=10.0.0'}
    peerDependencies:
      bufferutil: ^4.0.1
      utf-8-validate: '>=5.0.2'
    peerDependenciesMeta:
      bufferutil:
        optional: true
      utf-8-validate:
        optional: true

  xdg-basedir@5.1.0:
    resolution: {integrity: sha512-GCPAHLvrIH13+c0SuacwvRYj2SxJXQ4kaVTT5xgL3kPrz56XxkF21IGhjSE1+W0aw7gpBWRGXLCPnPby6lSpmQ==}
    engines: {node: '>=12'}

  y18n@5.0.8:
    resolution: {integrity: sha512-0pfFzegeDWJHJIAmTLRP2DwHjdF5s7jo9tuztdQxAhINCdvS+3nGINqPd00AphqJR/0LhANUS6/+7SCb98YOfA==}
    engines: {node: '>=10'}

  yargs-parser@21.1.1:
    resolution: {integrity: sha512-tVpsJW7DdjecAiFpbIB1e3qxIQsE6NoPc5/eTdrbbIC4h0LVsWhnoa3g+m2HclBIujHzsxZ4VJVA+GUuc2/LBw==}
    engines: {node: '>=12'}

  yargs-parser@22.0.0:
    resolution: {integrity: sha512-rwu/ClNdSMpkSrUb+d6BRsSkLUq1fmfsY6TOpYzTwvwkg1/NRG85KBy3kq++A8LKQwX6lsu+aWad+2khvuXrqw==}
    engines: {node: ^20.19.0 || ^22.12.0 || >=23}

  yargs@17.7.3:
    resolution: {integrity: sha512-GZtjxm/J/4TSxuL3FNYjCmLktBTnIw/rVmKSIyKeYAZpmJB2ig9VauCC5xsa82GNKVKDAqpOn3KVzNt0zmrU0g==}
    engines: {node: '>=12'}

  yargs@18.1.0:
    resolution: {integrity: sha512-2rAgRKu54VsHkqI0/tYkmluGXHD4KW7yZoycuqDQ15QOTnc2VVfy0nN/1eMhnQLO00A+dwtK20xuCnc1YGeUyg==}
    engines: {node: ^20.19.0 || ^22.12.0 || >=23}

  zod@3.25.76:
    resolution: {integrity: sha512-gzUt/qt81nXsFGKIFcC3YnfEAx5NkunCfnDlvuBSSFS02bcXu4Lmea0AFIUwbLWxWPx3d9p8S5QoaujKcNQxcQ==}

snapshots:

  '@apm-js-collab/code-transformer-bundler-plugins@0.7.4':
    dependencies:
      '@apm-js-collab/code-transformer': 0.18.1
      es-module-lexer: 2.3.2
      magic-string: 0.30.21
      module-details-from-path: 1.0.4

  '@apm-js-collab/code-transformer@0.18.1':
    dependencies:
      '@types/estree': 1.0.9
      astring: 1.9.0
      esquery: 1.7.0
      meriyah: 6.1.4
      semifies: 1.0.0
      source-map: 0.6.1

  '@apm-js-collab/tracing-hooks@0.13.0':
    dependencies:
      '@apm-js-collab/code-transformer': 0.18.1
      debug: 4.4.3
      module-details-from-path: 1.0.4
    transitivePeerDependencies:
      - supports-color

  '@formatjs/ecma402-abstract@2.3.6':
    dependencies:
      '@formatjs/fast-memoize': 2.2.7
      '@formatjs/intl-localematcher': 0.6.2
      decimal.js: 10.6.0
      tslib: 2.8.1

  '@formatjs/fast-memoize@2.2.7':
    dependencies:
      tslib: 2.8.1

  '@formatjs/icu-messageformat-parser@2.11.4':
    dependencies:
      '@formatjs/ecma402-abstract': 2.3.6
      '@formatjs/icu-skeleton-parser': 1.8.16
      tslib: 2.8.1

  '@formatjs/icu-skeleton-parser@1.8.16':
    dependencies:
      '@formatjs/ecma402-abstract': 2.3.6
      tslib: 2.8.1

  '@formatjs/intl-localematcher@0.6.2':
    dependencies:
      tslib: 2.8.1

  '@jridgewell/sourcemap-codec@1.6.0': {}

  '@opentelemetry/api-logs@0.220.0':
    dependencies:
      '@opentelemetry/api': 1.9.1

  '@opentelemetry/api@1.9.1': {}

  '@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1)':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/semantic-conventions': 1.43.0

  '@opentelemetry/instrumentation@0.220.0(@opentelemetry/api@1.9.1)':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/api-logs': 0.220.0
      import-in-the-middle: 3.3.3
      require-in-the-middle: 8.0.1
    transitivePeerDependencies:
      - supports-color

  '@opentelemetry/resources@2.10.0(@opentelemetry/api@1.9.1)':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/core': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/semantic-conventions': 1.43.0

  '@opentelemetry/sdk-trace-base@2.10.0(@opentelemetry/api@1.9.1)':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/core': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/resources': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/sdk-trace': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/semantic-conventions': 1.43.0

  '@opentelemetry/sdk-trace@2.10.0(@opentelemetry/api@1.9.1)':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/core': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/resources': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/semantic-conventions': 1.43.0

  '@opentelemetry/semantic-conventions@1.43.0': {}

  '@paulirish/trace_engine@0.0.65':
    dependencies:
      legacy-javascript: 0.0.1
      third-party-web: 0.29.2

  '@puppeteer/browsers@3.2.1':
    dependencies:
      modern-tar: 0.8.4
      yargs: 18.1.0

  '@sentry/conventions@0.16.0': {}

  '@sentry/core@10.71.0':
    dependencies:
      '@sentry/conventions': 0.16.0

  '@sentry/node-core@10.71.0(@opentelemetry/api@1.9.1)(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))(@opentelemetry/instrumentation@0.220.0(@opentelemetry/api@1.9.1))(@opentelemetry/sdk-trace-base@2.10.0(@opentelemetry/api@1.9.1))':
    dependencies:
      '@sentry/conventions': 0.16.0
      '@sentry/core': 10.71.0
      '@sentry/opentelemetry': 10.71.0(@opentelemetry/api@1.9.1)(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))(@opentelemetry/sdk-trace-base@2.10.0(@opentelemetry/api@1.9.1))
      import-in-the-middle: 3.3.3
    optionalDependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/core': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/instrumentation': 0.220.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/sdk-trace-base': 2.10.0(@opentelemetry/api@1.9.1)

  '@sentry/node@10.71.0(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/instrumentation': 0.220.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/sdk-trace-base': 2.10.0(@opentelemetry/api@1.9.1)
      '@sentry/conventions': 0.16.0
      '@sentry/core': 10.71.0
      '@sentry/node-core': 10.71.0(@opentelemetry/api@1.9.1)(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))(@opentelemetry/instrumentation@0.220.0(@opentelemetry/api@1.9.1))(@opentelemetry/sdk-trace-base@2.10.0(@opentelemetry/api@1.9.1))
      '@sentry/opentelemetry': 10.71.0(@opentelemetry/api@1.9.1)(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))(@opentelemetry/sdk-trace-base@2.10.0(@opentelemetry/api@1.9.1))
      '@sentry/server-utils': 10.71.0
      import-in-the-middle: 3.3.3
    transitivePeerDependencies:
      - '@opentelemetry/core'
      - '@opentelemetry/exporter-trace-otlp-http'
      - supports-color

  '@sentry/opentelemetry@10.71.0(@opentelemetry/api@1.9.1)(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))(@opentelemetry/sdk-trace-base@2.10.0(@opentelemetry/api@1.9.1))':
    dependencies:
      '@opentelemetry/api': 1.9.1
      '@opentelemetry/core': 2.10.0(@opentelemetry/api@1.9.1)
      '@opentelemetry/sdk-trace-base': 2.10.0(@opentelemetry/api@1.9.1)
      '@sentry/conventions': 0.16.0
      '@sentry/core': 10.71.0

  '@sentry/server-utils@10.71.0':
    dependencies:
      '@apm-js-collab/code-transformer-bundler-plugins': 0.7.4
      '@apm-js-collab/tracing-hooks': 0.13.0
      '@sentry/conventions': 0.16.0
      '@sentry/core': 10.71.0
      meriyah: 6.1.4
    transitivePeerDependencies:
      - supports-color

  '@types/estree@1.0.9': {}

  '@types/node@26.4.0':
    dependencies:
      undici-types: 8.3.0

  ansi-colors@4.1.3: {}

  ansi-regex@5.0.1: {}

  ansi-regex@6.3.0: {}

  ansi-styles@4.3.0:
    dependencies:
      color-convert: 2.0.1

  ansi-styles@6.2.3: {}

  astring@1.9.0: {}

  atomically@2.1.1:
    dependencies:
      stubborn-fs: 2.0.0
      when-exit: 2.1.5

  axe-core@4.13.0: {}

  chrome-launcher@1.2.1:
    dependencies:
      '@types/node': 26.4.0
      escape-string-regexp: 4.0.0
      is-wsl: 2.2.0
      lighthouse-logger: 2.0.2
    transitivePeerDependencies:
      - supports-color

  chromium-bidi@17.0.2(devtools-protocol@0.0.1666840):
    dependencies:
      devtools-protocol: 0.0.1666840
      mitt: 3.0.1
      zod: 3.25.76

  cjs-module-lexer@2.2.1: {}

  cliui@8.0.1:
    dependencies:
      string-width: 4.2.3
      strip-ansi: 6.0.1
      wrap-ansi: 7.0.0

  cliui@9.0.1:
    dependencies:
      string-width: 7.2.0
      strip-ansi: 7.2.0
      wrap-ansi: 9.0.2

  color-convert@2.0.1:
    dependencies:
      color-name: 1.1.4

  color-name@1.1.4: {}

  configstore@7.1.0:
    dependencies:
      atomically: 2.1.1
      dot-prop: 9.0.0
      graceful-fs: 4.2.11
      xdg-basedir: 5.1.0

  csp_evaluator@1.1.8: {}

  debug@4.4.3:
    dependencies:
      ms: 2.1.3

  decimal.js@10.6.0: {}

  define-lazy-prop@2.0.0: {}

  devtools-protocol@0.0.1663043: {}

  devtools-protocol@0.0.1666840: {}

  dot-prop@9.0.0:
    dependencies:
      type-fest: 4.41.0

  emoji-regex@10.6.0: {}

  emoji-regex@8.0.0: {}

  enquirer@2.4.1:
    dependencies:
      ansi-colors: 4.1.3
      strip-ansi: 6.0.1

  es-module-lexer@2.3.2: {}

  escalade@3.2.0: {}

  escape-string-regexp@4.0.0: {}

  esquery@1.7.0:
    dependencies:
      estraverse: 5.3.0

  estraverse@5.3.0: {}

  get-caller-file@2.0.5: {}

  get-east-asian-width@1.6.0: {}

  graceful-fs@4.2.11: {}

  http-link-header@1.1.4: {}

  image-ssim@0.2.0: {}

  import-in-the-middle@3.3.3:
    dependencies:
      cjs-module-lexer: 2.2.1
      es-module-lexer: 2.3.2
      module-details-from-path: 1.0.4

  intl-messageformat@10.7.18:
    dependencies:
      '@formatjs/ecma402-abstract': 2.3.6
      '@formatjs/fast-memoize': 2.2.7
      '@formatjs/icu-messageformat-parser': 2.11.4
      tslib: 2.8.1

  is-docker@2.2.1: {}

  is-fullwidth-code-point@3.0.0: {}

  is-wsl@2.2.0:
    dependencies:
      is-docker: 2.2.1

  jpeg-js@0.4.4: {}

  js-library-detector@6.7.0: {}

  legacy-javascript@0.0.1: {}

  lighthouse-logger@2.0.2:
    dependencies:
      debug: 4.4.3
      marky: 1.3.0
    transitivePeerDependencies:
      - supports-color

  lighthouse-stack-packs@1.12.3: {}

  lighthouse@13.4.1(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1)):
    dependencies:
      '@paulirish/trace_engine': 0.0.65
      '@sentry/node': 10.71.0(@opentelemetry/core@2.10.0(@opentelemetry/api@1.9.1))
      axe-core: 4.13.0
      chrome-launcher: 1.2.1
      configstore: 7.1.0
      csp_evaluator: 1.1.8
      devtools-protocol: 0.0.1663043
      enquirer: 2.4.1
      http-link-header: 1.1.4
      intl-messageformat: 10.7.18
      jpeg-js: 0.4.4
      js-library-detector: 6.7.0
      lighthouse-logger: 2.0.2
      lighthouse-stack-packs: 1.12.3
      lodash-es: 4.18.1
      lookup-closest-locale: 6.2.0
      open: 8.4.2
      puppeteer-core: 25.9.0
      robots-parser: 3.0.1
      speedline-core: 1.4.3
      third-party-web: 0.29.2
      tldts-icann: 7.4.11
      web-features: 3.36.0
      ws: 7.5.13
      yargs: 17.7.3
      yargs-parser: 21.1.1
    transitivePeerDependencies:
      - '@opentelemetry/core'
      - '@opentelemetry/exporter-trace-otlp-http'
      - bufferutil
      - proxy-agent
      - supports-color
      - utf-8-validate
      - yauzl

  lodash-es@4.18.1: {}

  lookup-closest-locale@6.2.0: {}

  magic-string@0.30.21:
    dependencies:
      '@jridgewell/sourcemap-codec': 1.6.0

  marky@1.3.0: {}

  meriyah@6.1.4: {}

  mitt@3.0.1: {}

  modern-tar@0.8.4: {}

  module-details-from-path@1.0.4: {}

  ms@2.1.3: {}

  open@8.4.2:
    dependencies:
      define-lazy-prop: 2.0.0
      is-docker: 2.2.1
      is-wsl: 2.2.0

  puppeteer-core@25.9.0:
    dependencies:
      '@puppeteer/browsers': 3.2.1
      chromium-bidi: 17.0.2(devtools-protocol@0.0.1666840)
      devtools-protocol: 0.0.1666840
      typed-query-selector: 2.12.2
      webdriver-bidi-protocol: 0.4.2
      ws: 8.21.3
    transitivePeerDependencies:
      - bufferutil
      - proxy-agent
      - utf-8-validate
      - yauzl

  require-directory@2.1.1: {}

  require-in-the-middle@8.0.1:
    dependencies:
      debug: 4.4.3
      module-details-from-path: 1.0.4
    transitivePeerDependencies:
      - supports-color

  robots-parser@3.0.1: {}

  semifies@1.0.0: {}

  source-map@0.6.1: {}

  speedline-core@1.4.3:
    dependencies:
      '@types/node': 26.4.0
      image-ssim: 0.2.0
      jpeg-js: 0.4.4

  string-width@4.2.3:
    dependencies:
      emoji-regex: 8.0.0
      is-fullwidth-code-point: 3.0.0
      strip-ansi: 6.0.1

  string-width@7.2.0:
    dependencies:
      emoji-regex: 10.6.0
      get-east-asian-width: 1.6.0
      strip-ansi: 7.2.0

  string-width@8.2.2:
    dependencies:
      get-east-asian-width: 1.6.0
      strip-ansi: 7.2.0

  strip-ansi@6.0.1:
    dependencies:
      ansi-regex: 5.0.1

  strip-ansi@7.2.0:
    dependencies:
      ansi-regex: 6.3.0

  stubborn-fs@2.0.0:
    dependencies:
      stubborn-utils: 1.0.2

  stubborn-utils@1.0.2: {}

  third-party-web@0.29.2: {}

  tldts-core@7.4.11: {}

  tldts-icann@7.4.11:
    dependencies:
      tldts-core: 7.4.11

  tslib@2.8.1: {}

  type-fest@4.41.0: {}

  typed-query-selector@2.12.2: {}

  undici-types@8.3.0: {}

  web-features@3.36.0: {}

  webdriver-bidi-protocol@0.4.2: {}

  when-exit@2.1.5: {}

  wrap-ansi@7.0.0:
    dependencies:
      ansi-styles: 4.3.0
      string-width: 4.2.3
      strip-ansi: 6.0.1

  wrap-ansi@9.0.2:
    dependencies:
      ansi-styles: 6.2.3
      string-width: 7.2.0
      strip-ansi: 7.2.0

  ws@7.5.13: {}

  ws@8.21.3: {}

  xdg-basedir@5.1.0: {}

  y18n@5.0.8: {}

  yargs-parser@21.1.1: {}

  yargs-parser@22.0.0: {}

  yargs@17.7.3:
    dependencies:
      cliui: 8.0.1
      escalade: 3.2.0
      get-caller-file: 2.0.5
      require-directory: 2.1.1
      string-width: 4.2.3
      y18n: 5.0.8
      yargs-parser: 21.1.1

  yargs@18.1.0:
    dependencies:
      cliui: 9.0.1
      escalade: 3.2.0
      get-caller-file: 2.0.5
      string-width: 8.2.2
      y18n: 5.0.8
      yargs-parser: 22.0.0

  zod@3.25.76: {}
````

### FILE: `google_lighthouse_web_quality_gate/source-lock.json`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite receipt of exact official Google/GitHub/npm identities"
license: "LicenseRef-Workspace-Owner"
sha256: "82f69486cbd9cc83092762b63f24a9c86806dd45c57746b2cf1ecdcdc96c0a6d"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": "1.0",
  "verifiedAt": "2026-08-29",
  "source": {
    "owner": "GoogleChrome",
    "repository": "lighthouse",
    "tag": "v13.4.1",
    "commit": "1d58f5b06d28e3419b38817a6c7488ec4413c67d",
    "tree": "405d16b2cdaac5c3ba3823b720031386a618165d",
    "commitSignatureVerified": false,
    "releaseId": 356969559,
    "publishedAt": "2026-07-20T20:10:47Z",
    "archiveUrl": "https://github.com/GoogleChrome/lighthouse/archive/1d58f5b06d28e3419b38817a6c7488ec4413c67d.zip",
    "archiveBytes": 73569399,
    "archiveSha256": "f526c6719e496519fc85986226364ef4e6274b1b0e2d8f64a5b8e22fc51510e3",
    "licenseExpression": "Apache-2.0",
    "licenseUrl": "https://raw.githubusercontent.com/GoogleChrome/lighthouse/v13.4.1/LICENSE",
    "licenseBytes": 11358,
    "licenseSha256": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "packageManifestSha256": "9a2acb5f3c0908842aacec66eb0bdd8307edd73469098e7eb186b5b79dfad1cf"
  },
  "npm": {
    "name": "lighthouse",
    "version": "13.4.1",
    "tarballUrl": "https://registry.npmjs.org/lighthouse/-/lighthouse-13.4.1.tgz",
    "tarballBytes": 3562319,
    "tarballSha256": "110759ba9e863c024e214e9b08ed2b0344d89b492286227235d4dfb990dc3e54",
    "integrity": "sha512-fDu8lt3QLK/lTqIxtp1HkzQNJ32rsFHhbadYOepcMZFLgA8oINhxutMbMv8XXnpTOvZ0TXCo4JCk1LDTWaRLnA==",
    "shasum": "40dedda9675b988ffd1b409678666a274977aff7",
    "unpackedBytes": 18983250,
    "nodeEngine": ">=22.19",
    "attestationsUrl": "https://registry.npmjs.org/-/npm/v1/attestations/lighthouse@13.4.1",
    "attestationsBytes": 14734,
    "attestationsSha256": "b1bd42f081b45ae2c6929be8084452a6ddb6f8ac47832e869862fe6adc0641f7",
    "attestationStatements": 2,
    "attestedSha512Hex": "7c3bbc96ddd02cafe54ea231b69d4793340d277dabb051e16da75839ea5c31914b800f2820d871bad31b32ff175e7a533af6744d70a8e090a4d4b0d359a44b9c",
    "resolvedProductionDependencies": 119,
    "pnpmAudit": {"date": "2026-08-29", "vulnerabilities": 0}
  },
  "browserContract": {
    "provider": "Microsoft Playwright",
    "packageVersion": "1.62.1",
    "browser": "Chromium",
    "browserVersion": "151.0.7922.34",
    "revision": "1234",
    "executableSha256Windows": "409805a16d6416087e6b2f778df1cf8f7bbb267d6b99f6b5bb0a618eace234f2"
  }
}
````

### FILE: `google_lighthouse_web_quality_gate/LICENSE.lighthouse.txt`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:lighthouse-license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleChrome/lighthouse/v13.4.1/LICENSE"
license: "Apache-2.0"
sha256: "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30"
variables: []
secrets_allowed: false
```
````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `google_lighthouse_web_quality_gate/quality-policy.json`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:quality-policy:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite fail-closed policy applying official Google Lighthouse output"
license: "LicenseRef-Workspace-Owner"
sha256: "68760a01e0d252427855cdce221c66cde4d5cdecc695bf1cd9e311d39bc983a2"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": "1.0",
  "expectedLighthouseVersion": "13.4.1",
  "numberOfRuns": 5,
  "maxRunWarnings": 0,
  "minimumAuditCount": 140,
  "categories": {
    "performance": {"median": 0.90, "minimum": 0.85},
    "accessibility": {"median": 1.00, "minimum": 1.00},
    "best-practices": {"median": 0.90, "minimum": 0.90},
    "seo": {"median": 1.00, "minimum": 1.00}
  }
}
````

### FILE: `google_lighthouse_web_quality_gate/lib/validate-reports.mjs`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:report-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite boundary validation around official Lighthouse reports"
license: "LicenseRef-Workspace-Owner"
sha256: "1256a637eade8918fcf8ba57b307d11f5924c13ce517c08428a34eac52c74464"
variables: []
secrets_allowed: false
```
````javascript
import {readFile} from 'node:fs/promises';

function fail(message) {
  throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: ${message}`);
}

function median(values) {
  const ordered = [...values].sort((a, b) => a - b);
  const middle = Math.floor(ordered.length / 2);
  return ordered.length % 2 ? ordered[middle] : (ordered[middle - 1] + ordered[middle]) / 2;
}

export function validatePolicy(policy) {
  if (policy?.schemaVersion !== '1.0') fail('unsupported policy schema');
  if (!/^\d+\.\d+\.\d+$/.test(policy.expectedLighthouseVersion ?? '')) fail('expectedLighthouseVersion must be exact');
  if (!Number.isInteger(policy.numberOfRuns) || policy.numberOfRuns < 3 || policy.numberOfRuns > 9 || policy.numberOfRuns % 2 === 0) {
    fail('numberOfRuns must be an odd integer from 3 through 9');
  }
  if (!Number.isInteger(policy.maxRunWarnings) || policy.maxRunWarnings < 0) fail('maxRunWarnings must be a non-negative integer');
  if (!Number.isInteger(policy.minimumAuditCount) || policy.minimumAuditCount < 1) fail('minimumAuditCount must be positive');
  const names = Object.keys(policy.categories ?? {});
  if (names.join(',') !== 'performance,accessibility,best-practices,seo') fail('all four categories are required in canonical order');
  for (const name of names) {
    const rule = policy.categories[name];
    for (const key of ['median', 'minimum']) {
      if (!Number.isFinite(rule?.[key]) || rule[key] < 0 || rule[key] > 1) fail(`${name}.${key} must be between 0 and 1`);
    }
    if (rule.minimum > rule.median) fail(`${name}.minimum cannot exceed median`);
  }
  return policy;
}

export function validateTargetUrl(rawUrl) {
  let url;
  try { url = new URL(rawUrl); } catch { fail('target URL is invalid'); }
  if (url.username || url.password) fail('credentials are forbidden in target URLs');
  const loopback = url.hostname === '127.0.0.1' || url.hostname === 'localhost' || url.hostname === '::1';
  if (url.protocol !== 'https:' && !(url.protocol === 'http:' && loopback)) fail('target must use HTTPS or loopback HTTP');
  if (url.hash) fail('URL fragments are not accepted as audit targets');
  return url;
}

export function validateReports(reports, policy, expectedUrl) {
  validatePolicy(policy);
  const target = validateTargetUrl(expectedUrl);
  if (!Array.isArray(reports) || reports.length !== policy.numberOfRuns) fail(`expected ${policy.numberOfRuns} reports`);
  const scoreSets = Object.fromEntries(Object.keys(policy.categories).map(name => [name, []]));
  const runs = [];
  for (const [index, report] of reports.entries()) {
    if (report?.lighthouseVersion !== policy.expectedLighthouseVersion) fail(`run ${index + 1} Lighthouse version drifted`);
    if (!Array.isArray(report.runWarnings) || report.runWarnings.length > policy.maxRunWarnings) fail(`run ${index + 1} has disallowed warnings`);
    const requested = validateTargetUrl(report.requestedUrl);
    const final = validateTargetUrl(report.finalUrl);
    if (requested.origin !== target.origin || final.origin !== target.origin) fail(`run ${index + 1} escaped the target origin`);
    const auditCount = Object.keys(report.audits ?? {}).length;
    if (auditCount < policy.minimumAuditCount) fail(`run ${index + 1} audit count ${auditCount} is below ${policy.minimumAuditCount}`);
    const scores = {};
    for (const name of Object.keys(policy.categories)) {
      const score = report.categories?.[name]?.score;
      if (!Number.isFinite(score) || score < 0 || score > 1) fail(`run ${index + 1} category ${name} has no numeric score`);
      scoreSets[name].push(score);
      scores[name] = score;
    }
    runs.push({run: index + 1, auditCount, warnings: report.runWarnings.length, scores});
  }
  const categories = {};
  for (const [name, values] of Object.entries(scoreSets)) {
    const actualMedian = median(values);
    const actualMinimum = Math.min(...values);
    const rule = policy.categories[name];
    if (actualMedian < rule.median) fail(`${name} median ${actualMedian} is below ${rule.median}`);
    if (actualMinimum < rule.minimum) fail(`${name} minimum ${actualMinimum} is below ${rule.minimum}`);
    categories[name] = {median: actualMedian, minimum: actualMinimum, values};
  }
  return {
    schemaVersion: '1.0',
    status: 'PASS',
    targetOrigin: target.origin,
    lighthouseVersion: policy.expectedLighthouseVersion,
    numberOfRuns: reports.length,
    categories,
    runs,
    limits: {
      automatedLabOnly: true,
      provesHumanAssistiveTechnology: false,
      provesProductionPerformance: false,
      provesLoadOrSecurity: false
    }
  };
}

export async function loadJson(path) {
  return JSON.parse(await readFile(path, 'utf8'));
}
````

### FILE: `google_lighthouse_web_quality_gate/run-quality-gate.mjs`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite orchestration of official Lighthouse and chrome-launcher APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "3a99c67e54fe4c5f7db16ba8ade037e24a3be9f4980fc8d0c31c300d2b135d4c"
variables: []
secrets_allowed: false
```
````javascript
import {mkdir, mkdtemp, rm, writeFile} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import {resolve, join} from 'node:path';
import process from 'node:process';
import lighthouse from 'lighthouse';
import {launch} from 'chrome-launcher';
import {loadJson, validatePolicy, validateReports, validateTargetUrl} from './lib/validate-reports.mjs';

function argument(name) {
  const prefix = `--${name}=`;
  const values = process.argv.slice(2).filter(value => value.startsWith(prefix));
  if (values.length !== 1 || values[0].length === prefix.length) throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: exactly one ${prefix}<value> is required`);
  return values[0].slice(prefix.length);
}

const targetUrl = validateTargetUrl(argument('url')).href;
const chromePath = resolve(argument('chrome-path'));
const outputDirectory = resolve(argument('output-dir'));
const policyPath = resolve(argument('policy'));
const policy = validatePolicy(await loadJson(policyPath));
await mkdir(outputDirectory, {recursive: true});
const profile = await mkdtemp(join(tmpdir(), 'elite-lighthouse-profile-'));
let chrome;
let primaryError;
try {
  chrome = await launch({
    chromePath,
    userDataDir: profile,
    chromeFlags: ['--headless', '--disable-gpu', '--no-sandbox'],
    logLevel: 'silent'
  });
  const reports = [];
  for (let index = 0; index < policy.numberOfRuns; index += 1) {
    const result = await lighthouse(targetUrl, {
      port: chrome.port,
      output: 'json',
      logLevel: 'silent',
      onlyCategories: Object.keys(policy.categories)
    });
    if (!result?.lhr || typeof result.report !== 'string') throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: run ${index + 1} returned no JSON report`);
    reports.push(result.lhr);
    await writeFile(join(outputDirectory, `lighthouse-run-${index + 1}.json`), result.report, {encoding: 'utf8', flag: 'wx'});
  }
  const summary = validateReports(reports, policy, targetUrl);
  await writeFile(join(outputDirectory, 'quality-gate-summary.json'), `${JSON.stringify(summary, null, 2)}\n`, {encoding: 'utf8', flag: 'wx'});
  process.stdout.write(`GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_PASS runs=${summary.numberOfRuns} performance_median=${summary.categories.performance.median} accessibility_median=${summary.categories.accessibility.median} best_practices_median=${summary.categories['best-practices'].median} seo_median=${summary.categories.seo.median}\n`);
} catch (error) {
  primaryError = error;
  throw error;
} finally {
  let cleanupError;
  try { if (chrome) await chrome.kill(); } catch (error) { cleanupError = error; }
  try { await rm(profile, {recursive: true, force: true, maxRetries: 20, retryDelay: 250}); } catch (error) { cleanupError ??= error; }
  if (!primaryError && cleanupError) throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: browser cleanup failed: ${cleanupError.message}`);
}
````

### FILE: `google_lighthouse_web_quality_gate/tests/validate-reports.test.mjs`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:validator-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite negative and positive contract suite"
license: "LicenseRef-Workspace-Owner"
sha256: "08f82f9d548b95312d3144f688deb3097a9af0313987d5214c1ca4f19c460787"
variables: []
secrets_allowed: false
```
````javascript
import test from 'node:test';
import assert from 'node:assert/strict';
import {validatePolicy, validateReports, validateTargetUrl} from '../lib/validate-reports.mjs';

const policy = {
  schemaVersion: '1.0',
  expectedLighthouseVersion: '13.4.1',
  numberOfRuns: 5,
  maxRunWarnings: 0,
  minimumAuditCount: 2,
  categories: {
    performance: {median: 0.90, minimum: 0.85},
    accessibility: {median: 1, minimum: 1},
    'best-practices': {median: 0.90, minimum: 0.90},
    seo: {median: 1, minimum: 1}
  }
};

function report(performance = 0.95) {
  return {
    lighthouseVersion: '13.4.1',
    requestedUrl: 'http://127.0.0.1:4181/',
    finalUrl: 'http://127.0.0.1:4181/',
    runWarnings: [],
    audits: {one: {}, two: {}},
    categories: {
      performance: {score: performance},
      accessibility: {score: 1},
      'best-practices': {score: 0.96},
      seo: {score: 1}
    }
  };
}

test('accepts five governed reports and exposes explicit limits', () => {
  const summary = validateReports([report(0.91), report(0.92), report(0.95), report(0.97), report(0.99)], policy, 'http://127.0.0.1:4181/');
  assert.equal(summary.status, 'PASS');
  assert.equal(summary.categories.performance.median, 0.95);
  assert.equal(summary.limits.provesHumanAssistiveTechnology, false);
});

test('rejects credentials, remote HTTP and fragments', () => {
  assert.throws(() => validateTargetUrl('https://user:secret@example.com/'));
  assert.throws(() => validateTargetUrl('http://example.com/'));
  assert.throws(() => validateTargetUrl('https://example.com/#private'));
});

test('rejects missing categories, even run counts and version ranges', () => {
  assert.throws(() => validatePolicy({...policy, numberOfRuns: 4}));
  assert.throws(() => validatePolicy({...policy, expectedLighthouseVersion: '^13.4.1'}));
  const incomplete = structuredClone(policy);
  delete incomplete.categories.seo;
  assert.throws(() => validatePolicy(incomplete));
});

test('rejects warnings, origin escapes, missing audits and version drift', () => {
  const warnings = report(); warnings.runWarnings = ['warning'];
  assert.throws(() => validateReports([warnings, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  const escape = report(); escape.finalUrl = 'https://example.com/';
  assert.throws(() => validateReports([escape, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  const few = report(); few.audits = {one: {}};
  assert.throws(() => validateReports([few, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  const drift = report(); drift.lighthouseVersion = '13.4.2';
  assert.throws(() => validateReports([drift, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
});

test('rejects a weak single run and a weak median', () => {
  assert.throws(() => validateReports([report(0.84), report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  assert.throws(() => validateReports([report(0.86), report(0.87), report(0.88), report(0.96), report(0.98)], policy, 'http://127.0.0.1:4181/'));
});
````

### FILE: `google_lighthouse_web_quality_gate/verify_contract.ps1`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite runtime, source, license and target verifier"
license: "LicenseRef-Workspace-Owner"
sha256: "8ac6b3ed5bf18231573f75ba92c02e050d4582909852f2ba4c2deb89768a3529"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [switch] $RunTarget,
  [string] $TargetUrl = '',
  [string] $ChromePath = '',
  [string] $OutputDirectory = ''
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$root = [IO.Path]::GetFullPath($PSScriptRoot)

function Fail([string] $Message) { throw "GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_FAILED: $Message" }

$nodeVersion = (& node --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $nodeVersion -notmatch '^v(?<major>\d+)\.') { Fail 'Node.js is unavailable' }
if ([int]$Matches.major -lt 22) { Fail "Node.js >=22.19 is required; actual=$nodeVersion" }
$pnpmVersion = (& pnpm --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $pnpmVersion -ne '11.25.0') { Fail "pnpm 11.25.0 is required; actual=$pnpmVersion" }
$lighthouseVersion = (& pnpm --dir $root --ignore-workspace exec lighthouse --version | Out-String).Trim()
if ($LASTEXITCODE -ne 0 -or $lighthouseVersion -ne '13.4.1') { Fail "Lighthouse 13.4.1 is required; actual=$lighthouseVersion" }

$sourceLock = Get-Content -Raw -LiteralPath (Join-Path $root 'source-lock.json') | ConvertFrom-Json -Depth 30
if ($sourceLock.source.commit -cne '1d58f5b06d28e3419b38817a6c7488ec4413c67d' -or $sourceLock.npm.tarballSha256 -cne '110759ba9e863c024e214e9b08ed2b0344d89b492286227235d4dfb990dc3e54') { Fail 'source identity drifted' }
$licensePath = Join-Path $root 'LICENSE.lighthouse.txt'
if ((Get-Item -LiteralPath $licensePath).Length -ne 11358 -or (Get-FileHash -Algorithm SHA256 -LiteralPath $licensePath).Hash.ToLowerInvariant() -cne 'cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30') { Fail 'Apache-2.0 license bytes drifted' }

& node --test (Join-Path $root 'tests/validate-reports.test.mjs')
if ($LASTEXITCODE -ne 0) { Fail 'negative policy suite failed' }

if ($RunTarget) {
  if ([string]::IsNullOrWhiteSpace($TargetUrl)) { Fail 'TargetUrl is required with RunTarget' }
  if ([string]::IsNullOrWhiteSpace($ChromePath)) { Fail 'ChromePath is required with RunTarget' }
  $resolvedChrome = [IO.Path]::GetFullPath($ChromePath)
  if (-not (Test-Path -LiteralPath $resolvedChrome -PathType Leaf)) { Fail "Chrome executable is missing: $resolvedChrome" }
  if ([string]::IsNullOrWhiteSpace($OutputDirectory)) { $OutputDirectory = Join-Path $root 'lighthouse-results' }
  $resolvedOutput = [IO.Path]::GetFullPath($OutputDirectory)
  if (Test-Path -LiteralPath $resolvedOutput) { Fail "output directory must not exist: $resolvedOutput" }
  $executionRoot = $root
  $capsule = $null
  try {
    if ($IsWindows -and $root.Length -gt 110) {
      $capsule = Join-Path ([IO.Path]::GetTempPath()) ('elite-lh-run-' + [guid]::NewGuid().ToString('N'))
      [void](New-Item -ItemType Directory -Path $capsule)
      [void](New-Item -ItemType Directory -Path (Join-Path $capsule 'lib'))
      foreach ($relative in @('package.json', 'pnpm-lock.yaml', 'quality-policy.json', 'run-quality-gate.mjs', 'lib/validate-reports.mjs')) {
        [IO.File]::Copy((Join-Path $root $relative), (Join-Path $capsule $relative), $false)
      }
      & pnpm --dir $capsule --ignore-workspace install --frozen-lockfile --offline
      if ($LASTEXITCODE -ne 0) { Fail 'failed to build the Windows short-path execution capsule' }
      $executionRoot = $capsule
    }
    & node (Join-Path $executionRoot 'run-quality-gate.mjs') "--url=$TargetUrl" "--chrome-path=$resolvedChrome" "--output-dir=$resolvedOutput" "--policy=$(Join-Path $executionRoot 'quality-policy.json')"
    if ($LASTEXITCODE -ne 0) { Fail 'target quality gate failed' }
  } finally {
    if ($capsule) {
      $temp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd('\')
      $resolvedCapsule = [IO.Path]::GetFullPath($capsule)
      if (-not $resolvedCapsule.StartsWith($temp + '\', [StringComparison]::OrdinalIgnoreCase) -or (Split-Path -Leaf $resolvedCapsule) -notlike 'elite-lh-run-*') {
        Fail 'refused unsafe Lighthouse execution-capsule cleanup'
      }
      $cleanupError = $null
      for ($attempt = 1; $attempt -le 5; $attempt++) {
        try {
          Remove-Item -LiteralPath $resolvedCapsule -Recurse -Force -ErrorAction Stop
          $cleanupError = $null
          break
        } catch {
          $cleanupError = $_
          Start-Sleep -Milliseconds 200
        }
      }
      if ($cleanupError -or (Test-Path -LiteralPath $resolvedCapsule)) { Fail 'failed to remove the Lighthouse execution capsule' }
    }
  }
}

Write-Output ('GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_CONTRACT_PASS runtime=13.4.1 tests=5 target=' + $(if ($RunTarget) {'5-runs'} else {'SKIPPED'}))
````

### FILE: `google_lighthouse_web_quality_gate/README.md`
```yaml
block_id: "GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite operator instructions"
license: "LicenseRef-Workspace-Owner"
sha256: "cf3ce29e9d2fe84eb06d05b42300d8ca2b646880732e6ec97b4ca55234567721"
variables: []
secrets_allowed: false
```
````markdown
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
````

## 6. Configuration surface

| Variable/input | Type | Safe default | Validation | Secret | Mutability/effect |
|---|---|---|---|---|---|
| `TargetUrl` | absolute URL | none | exactly one; HTTPS or loopback HTTP; no credentials/fragments | no | per audit target |
| `ChromePath` | file path | none | existing executable; project resolves it from its pinned browser runtime | no | per machine |
| `OutputDirectory` | absent directory | `<gate>/lighthouse-results` | must not exist; reports use exclusive writes | possibly sensitive output | per evidence run |
| `quality-policy.json` | JSON policy | 5 runs and governed floors | exact schema/version; odd 3–9 runs; four categories required | no | reviewed project policy |

No credentials, cookies or tokens belong in URLs, the pack or reports. Authenticated flows require a separate project-specific Playwright/identity gate.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| `lighthouse` | `13.4.1` | official audits and report model | Apache-2.0 | verification runtime | GoogleChrome/npm |
| `chrome-launcher` | `1.2.1` | explicit isolated Chromium lifecycle | Apache-2.0 | verification runtime | GoogleChrome/npm |
| Chromium | `151.0.7922.34` rev `1234` | audited browser in verified profile | see Playwright browser distribution | external verification runtime | Microsoft Playwright browser lock |
| Node.js | `>=22.19` | required by Lighthouse 13 | upstream terms | toolchain | nodejs.org |
| pnpm | `11.25.0` | frozen exact install | MIT | build/verification | pnpm.io |

The materialized lock resolves 119 production packages. Dated `pnpm audit` found zero vulnerabilities; dependency monitoring must rerun the project's update contract rather than silently changing the lock.

## 8. Apply order

1. Materialize into an empty destination or compose through the current web profile.
2. Resolve the project-selected Chromium executable; the verified composition uses Microsoft Playwright Chromium revision 1234.
3. Run `pnpm install --ignore-workspace --frozen-lockfile`; use offline mode when the exact store is already seeded.
4. Run `verify_contract.ps1` for source/runtime/policy regressions.
5. Build and start the target on loopback or select its HTTPS URL.
6. Run `verify_contract.ps1 -RunTarget ...`; retain five reports plus the summary receipt.
7. If any gate fails, stop promotion, record the failure and correct the target or reviewed policy. Never suppress an audit or warning to manufacture PASS.

For an existing workspace, keep this module isolated. Rollback removes the plan entry and directory; no database or application migration is involved.

## 9. Verification

Expected commands:

```powershell
pnpm install --ignore-workspace --frozen-lockfile --offline
pnpm audit --prod --audit-level=low
pwsh -NoProfile -File ./verify_contract.ps1
pwsh -NoProfile -File ./verify_contract.ps1 -RunTarget -TargetUrl http://127.0.0.1:4181/ -ChromePath <pinned-chromium> -OutputDirectory <new-evidence-dir>
```

Admission requires exact Lighthouse `13.4.1`, exact license/source identities, five policy tests, five target runs, no warnings, no origin escape, thresholds satisfied and cleanup complete. The verified enterprise web result is performance median/minimum `0.98/0.97`, accessibility `1/1`, best practices `1/0.96` and SEO `1/1`.

Negative coverage rejects credentials, remote HTTP, fragments, missing categories, even run counts, version ranges/drift, warnings, origin escape, incomplete audits, weak single runs and weak medians. Provider production, RUM, load, authenticated roles, assistive technologies, offensive security, deployment, canary and rollback remain separate project gates.

## 10. Reconstruction evidence

Governed evidence: `reconstruction_evidence/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_2026-08-29_V105.md`.

The clean reconstruction fixes Google release `v13.4.1`/commit `1d58f5b06d28e3419b38817a6c7488ec4413c67d`, archive/package/license hashes, Node `24.14.1`, pnpm `11.19.0`, 119 dependencies, zero `pnpm audit` findings, ten files, five policy tests, five real enterprise-web reports and the Windows long-path regression. The evidence also retains all failures and explicit non-claims.

## V321 — official package-manager security update

pnpm11.25.0 replaces vulnerable11.19.0 for this consumer; dependency lock bytes remain unchanged. Published bundle473packages/0advisories, registry ECDSA and publish/SLSA attestations verified, exact commit6d90c71efdffbc909b499490b64c66badc720327. Node24.20.0 standalone observed;115webPASS/1existingSKIP, build,92browser phases+agenda4,Playwright runtime4,Lighthouse policy5, local installation comparison. Eight changed materialized files across three owners; no upstream matcher or browser version changed. See reconstruction_evidence/PNPM_BUNDLE_SECURITY_V321.md. Native6files unchanged and transitive native SCA remains separate; product admission/monitor/target still blocked.
