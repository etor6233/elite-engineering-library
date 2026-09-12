# Human Review Leader Code Reaudit — V4

Date: 2026-08-28  
Scope: Microsoft Durable Task JS stable runtime and exact human-interaction examples  
Governing result: `RUNTIME_CONDITIONED`; both human-review examples `REJECTED_COMPONENT`

## Exact official source

- Owner/repository: Microsoft, `microsoft/durabletask-js`.
- Stable release/tag: `v0.4.0`, published 2026-07-31.
- Exact commit: `051bd5f4bb8a6a4c3408865625847abfc4659115`.
- Tag/commit status: lightweight tag over an unsigned commit.
- Current repository status: active, MIT; signed current `main` commit `910299294372c54986d96a0d9012d1082db840fc` dated 2026-08-27 is newer than the stable release and was not substituted for it.
- Archive: 1,372,198 bytes; SHA-256 `9b0b01aa6aab53f5865897ccf6dec9ab98dd847e358a147466e2d237208a9ba5`.
- License SHA-256: `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`.
- Root lock SHA-256: `913a7950f8a582d4b90d39fca099017d92081c5c999bdce323301cac6dc412bc`.
- Root manifest SHA-256: `c03d1d7277b52981ec3f08f78aeba1b88e6e1c2f21593e0da5dc0221c08dca10`.
- Official repository: https://github.com/microsoft/durabletask-js
- Stable source: https://github.com/microsoft/durabletask-js/tree/051bd5f4bb8a6a4c3408865625847abfc4659115

## Exact artifacts

| Artifact | SHA-256 | Exact behavior |
|---|---|---|
| `examples/hello-world/human_interaction.ts` | `50c38988307508b43dcace0c2f9da8c496796d104d9ee5129c3b6674fbf210e5` | purchase event, 24-hour timer and client-supplied approver string |
| `examples/azure-managed/human-interaction/index.ts` | `5e3ea4372353f49361957d51aa5668ab72d05a49addae4a0a857375ab8fb4c3f` | external event vs timer, custom status and managed/local scheduler clients |
| `examples/azure-managed/human-interaction/README.md` | `8d28fc6467037b781f6f59918fe018363d5be6b6a92547731e54b479cc40a5e2` | emulator/cloud run instructions; smoke searches only a completion string |

## Executed verification

- Node 24.14.1 satisfies the published `>=22` engine.
- npm 11.19.0 `ci` rejected the exact lock while trying to resolve `@types/node` 22.20.1; the manifest actually declares `^22.0.0`, so this was not misreported as an exact upstream pin.
- npm 10.9.2, isolated through `npx`, installed the exact lock: 560 packages, exit 0. It did not change global npm or authorize a generated replacement lock.
- Complete workspace build: PASS.
- ESLint: PASS.
- Unit suites: 98/98, 1,503/1,503 tests PASS — durable-functions 152, core 1,176, Azure managed 109, export-history 66.
- Search of the unit/e2e test trees found no test targeting either human-interaction example.
- npm audit over the exact lock: 8 findings — 4 high, 3 moderate, 1 low — across `@babel/core`, OpenTelemetry packages, `brace-expansion`, `js-yaml`, `protobufjs` and `tar`.
- The Docker/DTS emulator smoke and Azure Managed DTS were not run: they require Docker or an Azure resource/account and do not add human identity to the sample event.

## Contract audit

| Requirement | Managed example | Classic example | Result |
|---|---|---|---|
| durable wait/timeout | `yield whenAny([approvalEvent, timeout])` | declares event/timer | present only as orchestration primitive |
| human identity | event is `{approved:boolean}` | `approver` is typed by the same CLI client | not authenticated |
| mandatory reason | absent | absent | absent |
| field corrections | absent | absent | absent |
| assignment lease/expiry | timer only, no assignment | timer only | incomplete |
| CAS/version conflict | absent | absent | absent |
| request-bound decision evidence | orchestration history only; no signed decision receipt | string result only | absent |
| dual control | absent | absent | absent |
| correct race use | yes | calls `whenAny(tasks)` without `yield` and compares Task objects | classic path not proven |

Azure identity can authenticate the scheduler client in cloud mode; it does not authenticate the person who made the business decision. The local default explicitly uses `Authentication=None`.

## Admission decision

The runtime is legitimate Microsoft code and materially tested, so its exact event/timer primitives remain useful research. The two examples are not an immediately reusable enterprise approval implementation. Elite does not patch the missing `yield`, inject its own identity schema or combine the sample with authored persistence while attributing that result to Microsoft. Both examples remain rejected until Microsoft publishes a signed, clean, directly tested flow and the target proves authenticated identity, reason, corrections, conflict handling, durable evidence, recovery, load and operations.
