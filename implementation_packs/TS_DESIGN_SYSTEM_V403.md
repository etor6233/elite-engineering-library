# TS_DESIGN_SYSTEM_V403.md

## 1. Metadata

```yaml
pack_id: "TS-DESIGN-SYSTEM-V403"
pack_version: "0.3.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "V403 dark console and guided input overlay, local fixtures; visual approval and target acceptance remain separate."
stacks: ["React 19.2.8", "Next 16.3.4", "TypeScript 7.0.2", "Python 3.14"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.15.2", "TS-FRANCHISE-JOURNEY-PORTALS 0.21.0", "GO-CONNECTED-DOCUMENT-REFERENCE 0.1.0", "V403 cloud-run-transport overlay from the full composition"]
incompatible_with: ["Parallel active design systems", "Unreviewed target revisions"]
upstream_sources: []
license_expression: "LicenseRef-Workspace-Owner"
verified_at: "2026-09-14"
```

## 2. Scope and provenance

AUTHORED integration glue reuses existing admitted React/Next and owner contracts. No Adobe, Google icon, xAI, Tesla, SpaceX or DogeOS code/assets incorporated. Candidate identity/license observations do not admit their runtime. T2804 owner extension; V402/337 source and ZIPs remain immutable.

## 3. Integration contract

Compose this bundle outside the canonical library with its companion pack. Run `python experience_overlay/apply_overlay.py --target . --report qualification/ui-overlay.json`. All source and before hashes must match before writes. Run again with `--verify`; rollback is automatic for ordinary write failures, not a crash-atomic multi-file transaction. The single selected CSS/token authority replaces the previous style entry only in the consumer. No parallel design runtime.

Reference catalogue is opt-in `ELITE_EXPERIENCE_CATALOGUE=1`; never enable the synthetic operations route in production. Config files default empty. Bind reviewed tenant/organization/subject references in the consumer; absence is BUSINESS_CONFIGURATION_REQUIRED, not a credential claim. The quote, lead assignment, availability and checklist evidence POST boundaries repeat the reviewed selection check. Other added selectors guide UX while unchanged domain owners remain the authority for actual object/organization/business validation; the configuration is never a permission grant. Functions at /experience/reference are read-only maintenance specifications and never auth roles. Consumers bind their own owner records using the C contract binder; no maintenance state or permissions are copied into the product. /experience is a concise task entry; /franchise uses authorized sidebar areas and focus panes that retain visited operation state. Mobile navigation is a native modal dialog; exact code entry is keyboard/manual only, bound to a complete authorized supply order. Camera decoding is NOT_ADMITTED and physical reader execution is NOT_RUN. The invoice reference connects the existing document owner receive/process/review/decision flow; its binary transport is the same cloud-run-transport from owner D, not duplicated here. This UI-only bundle must be composed with D before building. See docs/DOCUMENT_EXPERIENCE_V403.md for the exact four-field, JPEG/PDF, reviewer and fixture limits. No new OCR or fiscal model is introduced.

## 4. Verification and limits

`pnpm install --offline --frozen-lockfile --ignore-scripts`; `pnpm exec tsc --noEmit`; focused tests in `src/design/experience.v403.test.ts`, `src/platform/experience/selection.v403.test.ts`, original command route tests; `pnpm exec next build --webpack`. Portable browser harness: `python experience_overlay/qualification/run_qualification.py --target . --evidence <local-directory> --candidate` after the standalone build. It creates ephemeral loopback keys outside source, exercises actual Next BFF against AUTHORED HTTP fixtures, and retains a hash-bound receipt. Run without --candidate to compare the captured candidates; comparison is never human approval. Browser fixture receipts and candidate screenshots are linked by FRANCHISE_EXPERIENCE_PACK_PLAN_V403. Automated contrast/reflow/keyboard do not prove screen-reader, physical phone or human acceptance. Cloud execution is NOT_RUN here.

## 5. Exact file manifest

```text
CREATE experience_overlay/files/src/app/globals.css
CREATE experience_overlay/files/src/design/tokens.v403.json
CREATE experience_overlay/files/src/design/contrast.v403.ts
CREATE experience_overlay/files/src/design/experience.v403.test.ts
CREATE experience_overlay/files/src/components/experience-ui.tsx
```

## 6. Materialization blocks

### FILE: `experience_overlay/files/src/app/globals.css`
```yaml
block_id: "TS-DESIGN-SYSTEM-V403:experience_overlay/files/src/app/globals.css:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "aca4eca8260b778880a10e8135d63a16efe484114fa76d348e280fae18a18db3"
variables: []
secrets_allowed: false
```
````css
/* GENERATED from src/design/tokens.v403.json. One active token authority. */
:root {
  --ink: #f2f3f5;
  --muted: #b0b4bc;
  --surface: #101113;
  --panel: #191b1f;
  --line: #6b717c;
  --accent: #e8eaee;
  --accent-strong: #c9cdd5;
  --on-accent: #101113;
  --warning: #ffd38b;
  --warning-bg: #35270f;
  --danger: #ffb6b6;
  --danger-bg: #3d1c20;
  --success: #a3e8c0;
  --success-bg: #173324;
  --muted-bg: #25282e;
  --focus: #a9c6ff;
  --ink-soft: #2b2e35;
  color-scheme: dark; font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}
[data-brand="indigo"] {--accent:#c4c9ff;--accent-strong:#aab2ff;--focus:#bcc4ff;}
* { box-sizing: border-box; }
body { margin: 0; color: var(--ink); background: var(--surface); }
a { color: inherit; }
.skipLink { position: fixed; z-index: 100; left: 1rem; top: 1rem; padding: .7rem 1rem; border-radius: 8px; background: var(--ink); color: var(--on-accent); transform: translateY(-180%); }
.skipLink:focus { transform: translateY(0); }
header { background: var(--panel); border-bottom: 1px solid var(--line); }
.shell { width: min(1120px, calc(100% - 2rem)); margin: 0 auto; }
.headerRow { min-height: 68px; display: flex; align-items: center; justify-content: space-between; gap: 1.5rem; }
.brand { font-weight: 800; text-decoration: none; letter-spacing: -0.03em; }
nav { display: flex; flex-wrap: wrap; gap: 1rem; }
nav a { color: var(--muted); text-decoration: none; font-weight: 650; }
nav a:hover, nav a:focus-visible { color: var(--accent); text-decoration: underline; }
main { padding: 3rem 0 5rem; }
.hero { display: grid; grid-template-columns: minmax(0, 1.5fr) minmax(280px, .8fr); gap: 2rem; align-items: center; }
.eyebrow { color: var(--accent); font-weight: 800; text-transform: uppercase; font-size: .78rem; letter-spacing: .12em; }
h1 { font-size: clamp(2.4rem, 7vw, 5.5rem); line-height: .95; letter-spacing: -.06em; margin: .6rem 0 1.2rem; }
.pageTitle { font-size: clamp(2.2rem, 5vw, 4rem); }
h2 { font-size: clamp(1.5rem, 3vw, 2.4rem); letter-spacing: -.035em; }
p { line-height: 1.65; }
.lede { font-size: 1.18rem; color: var(--muted); max-width: 68ch; }
.panel, .card { background: var(--panel); border: 1px solid var(--line); border-radius: 18px; padding: 1.4rem; box-shadow: 0 14px 35px rgba(18, 34, 28, .06); }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1rem; margin-top: 1.5rem; }
.metric { font-size: 2rem; font-weight: 850; color: var(--accent); }
.muted { color: var(--muted); }
.button { display: inline-flex; align-items: center; justify-content: center; padding: .8rem 1rem; border-radius: 10px; border: 0; background: var(--accent); color: var(--on-accent); font-weight: 750; text-decoration: none; cursor: pointer; }
.button:hover, .button:focus-visible { background: var(--accent-strong); }
.actions { display: flex; flex-wrap: wrap; gap: .75rem; margin-top: 1.5rem; }
.primaryAction { display: inline-flex; padding: .85rem 1.1rem; border-radius: 10px; background: var(--accent); color: var(--on-accent); font-weight: 750; text-decoration: none; }
.primaryAction:hover, .primaryAction:focus-visible { background: var(--accent-strong); }
label { display: grid; gap: .4rem; font-weight: 700; }
input, select { width: 100%; border: 1px solid var(--line); border-radius: 9px; padding: .75rem; font: inherit; background: var(--panel); }
input:focus, select:focus { outline: 3px solid var(--focus); border-color: var(--accent); }
form { display: grid; gap: 1rem; }
.status { padding: .8rem; border-radius: 8px; background: var(--success-bg); color: var(--success); }
.notice { margin: 1rem 0 2rem; padding: 1rem; border: 1px solid var(--warning); border-radius: 10px; background: var(--warning-bg); color: var(--warning); line-height: 1.6; }
.error { background: var(--danger-bg); color: var(--danger); }
footer { padding: 2rem 0; border-top: 1px solid var(--line); color: var(--muted); background: var(--panel); }
code { background: var(--ink-soft); border-radius: 5px; padding: .12rem .3rem; }
.notificationStatus { border-top: 1px solid var(--line); margin-top: 1rem; padding-top: 1rem; overflow-wrap: anywhere; }
.notificationStatus dd { margin-inline-start: 0; }
@media (max-width: 760px) { .hero { grid-template-columns: 1fr; } .headerRow { align-items: flex-start; flex-direction: column; padding: 1rem 0; } }
@media (prefers-reduced-motion: reduce) { *, *::before, *::after { scroll-behavior: auto !important; transition: none !important; } }

button, textarea { font: inherit; }
button { min-height: 44px; padding: .7rem 1rem; border-radius: 10px; border: 1px solid var(--line); color:var(--ink); background:var(--panel); cursor:pointer; }
button:disabled { cursor:not-allowed; color:var(--muted); background:var(--muted-bg); }
:focus-visible { outline:3px solid var(--focus); outline-offset:3px; }
textarea { width:100%; min-height:6rem; border:1px solid var(--line); border-radius:9px; padding:.75rem; color:var(--ink); background:var(--panel); }
input[type="checkbox"] { width:1.4rem; height:1.4rem; accent-color:var(--accent); }
fieldset { min-width:0; border:0; padding:0; display:grid; gap:1rem; }
legend { font-weight:750; padding:0 0 .8rem; }
.badge { display:inline-flex; align-items:center; gap:.4rem; padding:.3rem .65rem; border-radius:999px; font-size:.875rem; font-weight:700; }
.badge-success { color:var(--success); background:var(--success-bg); }
.badge-warning { color:var(--warning); background:var(--warning-bg); }
.badge-danger { color:var(--danger); background:var(--danger-bg); }
.badge-muted { color:var(--muted); background:var(--muted-bg); }
.stack { display:grid; gap:1rem; }.cluster { display:flex; flex-wrap:wrap; align-items:center; gap:.75rem; }
.workspace { display:grid; grid-template-columns:240px minmax(0,1fr); gap:2rem; }
.workspace aside { border-right:1px solid var(--line); padding-right:1.25rem; }
.workspace aside a { display:block; padding:.75rem; border-radius:8px; text-decoration:none; }
.workspace aside a[aria-current="page"] { background:var(--accent); color:var(--on-accent); }
.workspace h1 { font-size:clamp(2rem,4vw,3.5rem); line-height:1.1; }
.toolbar { display:flex; flex-wrap:wrap; gap:1rem; align-items:end; }.toolbar label { flex:1; min-width:160px; }
.tableWrap { overflow-x:auto; border:1px solid var(--line); border-radius:12px; }
table { width:100%; border-collapse:collapse; background:var(--panel); }caption { text-align:left; padding:1rem; font-weight:750; }
th,td { text-align:left; padding:1rem; border-bottom:1px solid var(--line); overflow-wrap:anywhere; }th { background:var(--muted-bg); }
details { border:1px solid var(--line); border-radius:10px; padding:1rem; } summary { cursor:pointer; font-weight:750; }
.checkRow { display:flex; align-items:start; gap:.7rem; }.checkRow input { flex:none; }
.muted a { text-decoration:underline; }.srOnly { position:absolute; width:1px; height:1px; padding:0; margin:-1px; overflow:hidden; clip:rect(0,0,0,0); white-space:nowrap; border:0; }
.catalogue { overflow-wrap:anywhere; }.catalogue h2 { font-size:1.6rem; }.catalogue h3 { margin-top:0; }
.catalogue section { scroll-margin-top:1rem; }.catalogue .hero { padding:2rem 0; }.visualMark { border-radius:24px; background:var(--accent); color:var(--on-accent); padding:3rem; font-size:3.5rem; font-weight:800; line-height:1; }
.compactList { padding-left:1.25rem; }.compactList li { margin:.55rem 0; }
@media(max-width:760px) { .workspace {grid-template-columns:1fr;} .workspace aside {border-right:0;border-bottom:1px solid var(--line);padding:0 0 1rem;} .workspace aside nav {display:flex;} .workspace aside a {display:inline-block;} .tableWrap {max-width:100%;} .visualMark {font-size:2.5rem;padding:2rem;} }

pre {white-space:pre-wrap;overflow-wrap:anywhere;}

.journeyTable th:first-child {width:60%;}.areaCell {white-space:nowrap;}.journeyCards {display:none;list-style:none;margin:0;padding:0;}
@media(max-width:760px){.journeyTable {display:none;}.journeyCards {display:grid;gap:1rem;}.journeyCards h3 {margin-top:.75rem;}.journeyCards a {display:inline-flex;align-items:center;min-height:44px;}}

main {overflow-wrap:anywhere;}.stack > * {min-width:0;}

/* V403 focus pattern: one active task, no operation markers are unmounted on navigation. */
[hidden]{display:none!important;}.taskWorkspace{display:grid;gap:1.5rem}.taskNavigation{display:flex;gap:.5rem;overflow-x:auto;flex-wrap:nowrap;padding:.25rem 0}.taskNavigation button{white-space:nowrap;border:0;background:transparent;font-weight:650;color:var(--muted)}.taskNavigation button[aria-current="page"]{background:var(--muted-bg);color:var(--ink)}.taskBody{min-width:0}.taskBody h2{font-size:1.5rem}.taskBody>div>.grid{margin-top:0}.taskBody details{margin-top:1rem}.experienceApp{max-width:980px;margin:0 auto}.experienceTopline{display:flex;justify-content:space-between;align-items:center;color:var(--muted);font-size:.875rem;margin:0 0 3rem}.experienceTopline a{text-decoration:none}.experienceApp h1{font-size:clamp(2.7rem,5vw,4.5rem);line-height:1.06;font-weight:650;letter-spacing:-.05em;max-width:750px;margin:0 0 1rem}.experienceIntro{color:var(--muted);font-size:1.08rem;margin:0 0 2.4rem}.nextTask{display:grid;grid-template-columns:1fr 220px;align-items:center;gap:2rem;background:var(--panel);border:1px solid var(--line);border-radius:22px;padding:2.5rem;margin:2.5rem 0 1.3rem}.taskEyebrow{font-size:.72rem;letter-spacing:.13em;font-weight:750;color:var(--accent)}.nextTask h2{font-size:1.8rem;line-height:1.15;margin:.9rem 0}.nextTask p{color:var(--muted);margin:.75rem 0 1.6rem}.nextTask .primaryAction{gap:2.5rem;align-items:center}.taskSteps{border-left:1px solid var(--line);padding-left:2rem}.taskSteps p{display:flex;gap:1rem;color:var(--ink);font-size:.9rem;margin:1.2rem 0}.taskSteps span{color:var(--muted);font-variant-numeric:tabular-nums}.quickTasks{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem}.quickTasks a{display:flex;justify-content:space-between;gap:1rem;padding:1.2rem .2rem;border-bottom:1px solid var(--line);font-size:.95rem;color:var(--ink)}.experienceOperation{max-width:620px;margin:0 auto}.experienceOperation>a:first-child{display:inline-block;margin-bottom:2rem;color:var(--muted);font-size:.9rem;text-decoration:none}.experienceOperation h1{font-size:2.8rem;line-height:1.08}.experienceOperation [role="region"]>form{background:var(--panel);border:1px solid var(--line);border-radius:20px;padding:1.8rem}.experienceOperation fieldset>legend{display:none}.experienceOperation form button:not([type="button"]){background:var(--accent);color:var(--on-accent);border-color:var(--accent);margin-top:.5rem}.experienceOperation [role="status"]:empty{display:none}.experienceOperation [role="region"]>details{margin-top:1rem}.experienceOperation label{font-size:.95rem}.experienceOperation small{font-weight:400;color:var(--muted)}
@media(max-width:760px){.experienceApp h1{font-size:2.75rem}.experienceTopline{margin-bottom:2.3rem}.nextTask{grid-template-columns:1fr;padding:1.5rem;gap:1rem;margin-top:2rem}.nextTask h2{font-size:1.6rem}.taskSteps{display:flex;gap:1.5rem;border-left:0;border-top:1px solid var(--line);padding:1rem 0 0;justify-content:space-between}.taskSteps p{display:grid;gap:.3rem;margin:0;font-size:.8rem}.quickTasks{grid-template-columns:1fr;gap:0}.quickTasks a{min-height:64px;align-items:center}.experienceOperation h1{font-size:2.15rem}.experienceOperation [role="region"]>form{padding:1.2rem}.experienceOperation .checkRow{align-items:center}.taskNavigation{margin-inline:-.25rem}.taskBody .card{padding:1.1rem}}

.focusTitle{font-size:clamp(2rem,4vw,3rem);line-height:1.12;margin-bottom:1.7rem}.branchContext{font-size:.85rem;color:var(--muted);margin:0 0 1rem}.taskPicker{max-width:400px;width:100%;font-size:.9rem}.referenceDetails{border:0;padding:0;margin:.5rem 0 1rem!important;font-size:.8rem;color:var(--muted)}.referenceDetails summary{font-weight:500}.taskBody .card{box-shadow:none}.taskBody .grid{grid-template-columns:minmax(0,1fr);max-width:760px}.taskBody [role="status"]:empty{display:none}

.compactAppHeader>.shell{min-height:68px;display:flex;align-items:center;justify-content:space-between;gap:1rem}.appHeaderActions{display:flex;gap:.4rem;align-items:center}.compactAppHeader details{position:relative;border:0;padding:.7rem .6rem;font-size:.85rem}.compactAppHeader summary{font-weight:600}.compactAppHeader details[open]>nav,.compactAppHeader details[open]>form{position:absolute;right:0;top:100%;width:min(340px,calc(100vw - 2rem));padding:1.2rem;background:var(--panel);border:1px solid var(--line);border-radius:12px;z-index:30;box-shadow:0 12px 32px #12221c18}.compactAppHeader nav{display:grid}.compactAppHeader nav a{padding:.6rem}.compactAppHeader details p{font-size:.9rem}.experienceApp h1{font-size:clamp(2rem,4vw,3rem);margin:0 0 1.5rem;line-height:1.1}.experienceApp .nextTask{margin-top:0}.experienceOperation>a:first-child{margin-bottom:1.3rem}.experienceOperation .experienceIntro{margin-bottom:1.5rem}
@media(max-width:760px){.compactAppHeader>.shell{min-height:64px;gap:.6rem}.compactAppHeader .brand{font-size:.85rem;max-width:180px}.appHeaderActions{gap:0}.compactAppHeader details{font-size:.8rem;padding:.65rem .4rem}main:has(.experienceApp),main:has(.experienceOperation){padding-top:1.5rem}.experienceOperation h1{margin-bottom:1rem}.experienceApp .taskSteps{padding-top:.8rem}.experienceApp .nextTask{margin-bottom:.7rem}}

.taskBody form+form{margin-top:1.6rem}.taskBody form>button{justify-self:start}.taskBody form+button,.taskBody [role="status"]+button{margin-top:1rem}.taskBody .primaryAction{border:0}.taskBody .grid article>h2{margin-bottom:.5rem}.taskBody .grid article>p{margin:.5rem 0}
@media(max-width:760px){.taskBody form>button{width:100%;justify-self:stretch}}

label:has(>input[type="checkbox"]){display:flex;align-items:flex-start;gap:.7rem}.publicCatalogue h1{font-size:clamp(2.2rem,4.5vw,3.6rem)}.publicCatalogue .card{display:grid;grid-template-columns:minmax(180px,.8fr) minmax(280px,1.2fr);column-gap:3rem;padding:2rem;align-content:start;box-shadow:none}.publicCatalogue .card>form{grid-column:2;grid-row:1 / span 3}.publicCatalogue .card>h2{margin:.5rem 0;font-size:2.4rem}.publicCatalogue .card>p{margin:0;color:var(--muted)}.publicCatalogue button[type="submit"]{justify-self:start;background:var(--accent);color:var(--on-accent);border-color:var(--accent)}.publicCatalogue [role="status"]:empty{display:none}
@media(max-width:760px){.publicCatalogue .card{display:block;padding:1.4rem}.publicCatalogue .card>p{margin:.5rem 0 1.5rem}.publicCatalogue button[type="submit"]{width:100%;justify-self:stretch}}

/* V403 UI 0.3: one token system, dark console shell, native modal mobile navigation. */
input, select {color:var(--ink)}
.panel,.card {border-radius:12px;box-shadow:none;border-color:var(--muted-bg)}
.applicationSidebar {position:fixed;inset:0 auto 0 0;width:272px;padding:0 18px;display:flex;flex-direction:column;background:var(--surface);border-right:1px solid var(--muted-bg);overflow-y:auto;z-index:35}
.applicationBrand {display:flex;min-height:88px;align-items:center;justify-content:space-between;gap:12px;padding:8px 12px}
.applicationBrand>a {font-size:1rem;font-weight:650;letter-spacing:-.03em;text-decoration:none}
.applicationNavigation {display:grid;gap:4px;margin-top:20px}
.applicationNavigation button,.applicationNavigation a {display:flex;align-items:center;gap:12px;min-height:46px;padding:10px 14px;color:var(--muted);font-size:.91rem;font-weight:500;text-align:left;text-decoration:none;border:0;border-radius:8px;background:transparent}
.applicationNavigation svg,.applicationBottom svg {width:20px;height:20px;flex:none}
.applicationNavigation [aria-current="page"] {background:var(--muted-bg);color:var(--ink)}
.applicationNavigation a:hover,.applicationNavigation button:hover {background:var(--panel);color:var(--ink)}
.applicationSpaces {margin-top:24px;border:0;padding:12px 14px;color:var(--muted);font-size:.85rem}
.applicationSpaces summary {font-weight:500}.applicationSpaces nav{display:grid;gap:4px;margin-top:12px}.applicationSpaces a{padding:8px 0;font-weight:400}
.applicationBottom {display:grid;gap:4px;margin-top:auto;padding:30px 12px 24px;font-size:.88rem}
.applicationBottom>a {display:flex;align-items:center;gap:12px;padding:10px 0;text-decoration:none;color:var(--muted)}
.applicationBottom details {border:0;padding:10px 0;color:var(--muted)}.applicationBottom summary{font-weight:500}.applicationBottom form {margin-top:16px;font-size:.85rem}.applicationBottom form p{font-size:.8rem}
.applicationMobileHeader {display:none}.applicationDrawer:not([open]){display:none}
body:has(.applicationSidebar)>div {min-height:100vh}
body:has(.applicationSidebar) main.shell {width:auto;max-width:none;margin-left:272px;padding:48px clamp(24px,5vw,96px) 80px}
body:has(.applicationSidebar) main> * {max-width:1100px}
body:has(.applicationSidebar) h1 {font-size:1.8rem;font-weight:600;letter-spacing:-.04em;line-height:1.2;margin:0 0 30px}
.branchContext {margin:0 0 10px;font-size:.8rem}
.taskWorkspace {gap:20px}.taskNavigation {padding:0;gap:4px}.taskNavigation button{min-height:38px;padding:8px 13px;font-size:.85rem;border-radius:7px}
.taskPicker {max-width:280px;font-size:.8rem;font-weight:500;gap:8px}
.taskPicker select {font-size:.9rem;padding:10px 12px;min-height:44px}
.taskBody .card {padding:24px;border:0;background:var(--panel)}
.taskBody .card h2 {font-size:1.2rem;font-weight:550}
.taskBody .grid {max-width:760px}.taskBody details{border-color:var(--muted-bg)}
.taskBody form {gap:18px}.taskBody label{font-size:.85rem;font-weight:500}.taskBody input,.taskBody select{padding:11px 12px}
.taskBody form>button {font-size:.85rem;min-height:42px}
.taskBody .referenceDetails {font-size:.75rem}
.experienceApp {max-width:940px;margin:0}.experienceApp .nextTask {border:0;border-radius:12px;padding:28px 32px;margin-top:0;grid-template-columns:1fr 190px}
.nextTask h2 {font-size:1.5rem;font-weight:550}.nextTask p{font-size:.9rem}.nextTask .primaryAction {font-size:.9rem;padding:12px 18px}
.taskSteps {border-color:var(--muted-bg)}.taskEyebrow{color:var(--muted);font-weight:600;font-size:.68rem}
.quickTasks {gap:24px}.quickTasks a{font-size:.88rem;font-weight:450;border-color:var(--muted-bg);text-decoration:none}
.experienceOperation {max-width:650px;margin:0}.experienceOperation>a:first-child{font-size:.8rem;margin-bottom:24px}
.experienceOperation .experienceIntro{font-size:.9rem}.experienceOperation [role="region"]>form{border:0;border-radius:12px;padding:28px}
.experienceOperation label{font-size:.88rem;font-weight:500}.experienceOperation small{font-size:.8rem}
.experienceOperation [role="region"]>details{font-size:.85rem;border-color:var(--muted-bg)}
.iconButton {width:44px;height:44px;padding:10px;border:0;background:transparent;font-size:28px;font-weight:300;display:grid;place-content:center}
.iconButton svg{width:22px;height:22px}
.publicCatalogue .card {display:block;max-width:760px;padding:32px}.publicCatalogue .card>h2{font-size:2rem}.publicCatalogue .modelInquiry {border:0;padding:0;margin-top:28px;max-width:540px}
.modelInquiry>summary {display:inline-flex;align-items:center;min-height:44px;padding:12px 18px;border-radius:8px;background:var(--accent);color:var(--on-accent);font-size:.9rem;font-weight:600;list-style:none}
.modelInquiry>summary::-webkit-details-marker {display:none}.modelInquiry[open]>summary{margin-bottom:24px}.modelInquiry form{padding-top:4px}.modelInquiry label{font-size:.9rem;font-weight:500}
@media(max-width:900px){
 .applicationSidebar {display:none}
 .applicationMobileHeader {display:flex;align-items:center;gap:14px;height:64px;padding:0 16px;background:var(--surface);border-bottom:1px solid var(--muted-bg)}
 .applicationMobileHeader>span {font-size:.85rem;font-weight:600;letter-spacing:-.025em}
 .applicationDrawer {position:fixed;inset:0;margin:0;width:100%;height:100dvh;max-height:100dvh;max-width:none;border:0;padding:0 20px;background:var(--surface);color:var(--ink);display:flex;flex-direction:column}
 .applicationDrawer::backdrop {background:var(--surface)}.applicationDrawer .applicationBrand{min-height:72px;padding:8px 0}.applicationDrawer .applicationNavigation{margin-top:20px}.applicationDrawer .applicationNavigation button,.applicationDrawer .applicationNavigation a{min-height:52px;font-size:1rem}.applicationDrawer .applicationBottom{padding-bottom:24px}
 body:has(.applicationSidebar) main.shell {margin:0;padding:28px 22px 56px;width:100%}
 body:has(.applicationSidebar) h1{font-size:1.55rem;margin-bottom:24px}
 .taskBody .card{padding:20px 18px}.taskPicker{max-width:none}.taskNavigation{max-width:100%;margin:0}.taskBody form>button{width:100%}
 .experienceApp .nextTask{padding:24px 20px;grid-template-columns:1fr;gap:16px}.nextTask h2{font-size:1.4rem}.nextTask p{font-size:.88rem}.taskSteps{border-top:1px solid var(--muted-bg)}
 .experienceOperation [role="region"]>form{padding:22px 18px}.experienceOperation .experienceIntro{font-size:.88rem}.experienceOperation>a:first-child{margin-bottom:20px}.experienceOperation .checkRow{align-items:flex-start}
 .quickTasks{gap:0}.quickTasks a{min-height:60px}
 .publicCatalogue .card{padding:24px}.publicCatalogue .modelInquiry>summary{width:100%;justify-content:center}
}

.contactWork {border:0;padding:0}.contactWork>summary{display:inline-flex;min-height:42px;align-items:center;border:1px solid var(--line);border-radius:8px;padding:10px 14px;font-size:.85rem;font-weight:550}.contactWork[open]>summary{margin-bottom:22px}.contactWork>.taskWorkspace{margin-top:4px}.taskWorkspace>.focusTitle{margin-bottom:0!important}.unitCodePicker{padding:20px 0;max-width:650px}.unitCodeEntry{display:flex;gap:10px;margin-top:8px}.unitCodeEntry input{min-width:0}.unitCodePicker>.muted{font-size:.8rem;margin:10px 0}.unitCodePicker [role="status"]{font-size:.9rem}.unitCodePicker button{white-space:nowrap}

.documentWorkspace {display:grid;gap:20px;max-width:680px}.documentWorkspace .surface{padding:26px;background:var(--panel);border-radius:12px}.documentWorkspace h2{font-size:1.2rem;font-weight:550;margin:0 0 20px}.documentWorkspace header{background:transparent;border:0;margin-bottom:20px}.documentWorkspace header p{font-size:.85rem;color:var(--muted)}.documentWorkspace label{font-size:.88rem;font-weight:500}.documentWorkspace .muted{font-size:.85rem}.documentWorkspace article>button,.documentWorkspace article>form,.documentWorkspace article>div{margin-top:22px}.documentWorkspace form{margin-top:22px}.documentWorkspace form>button{justify-self:start}.documentWorkspace dl{display:grid;gap:16px}.documentWorkspace dl>div{display:grid;grid-template-columns:1fr 2fr;gap:16px}.documentWorkspace dt{color:var(--muted);font-size:.85rem}.documentWorkspace dd{margin:0;overflow-wrap:anywhere}.documentWorkspace details{margin-top:24px;border-color:var(--muted-bg);font-size:.85rem}.documentWorkspace>button{justify-self:start}
@media(max-width:900px){.taskBody form>button,.contactWork>summary{min-height:44px}.documentWorkspace .surface{padding:22px 18px}.documentWorkspace form>button,.documentWorkspace>button{width:100%;justify-self:stretch}.documentWorkspace dl>div{grid-template-columns:1fr;gap:5px}}
````

### FILE: `experience_overlay/files/src/design/tokens.v403.json`
```yaml
block_id: "TS-DESIGN-SYSTEM-V403:experience_overlay/files/src/design/tokens.v403.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "756567004c1dc2d3afbd7d7d9f166621ffe4be195f8486a61210e3a27d9a02e7"
variables: []
secrets_allowed: false
```
````json
{
  "revision": "V403-UI-0.3.0",
  "default": {
    "ink": "#f2f3f5",
    "muted": "#b0b4bc",
    "surface": "#101113",
    "panel": "#191b1f",
    "line": "#6b717c",
    "accent": "#e8eaee",
    "accent-strong": "#c9cdd5",
    "on-accent": "#101113",
    "warning": "#ffd38b",
    "warning-bg": "#35270f",
    "danger": "#ffb6b6",
    "danger-bg": "#3d1c20",
    "success": "#a3e8c0",
    "success-bg": "#173324",
    "muted-bg": "#25282e",
    "focus": "#a9c6ff",
    "ink-soft": "#2b2e35"
  },
  "brands": {
    "forest": {},
    "indigo": {
      "accent": "#c4c9ff",
      "accent-strong": "#aab2ff",
      "focus": "#bcc4ff"
    }
  }
}
````

### FILE: `experience_overlay/files/src/design/contrast.v403.ts`
```yaml
block_id: "TS-DESIGN-SYSTEM-V403:experience_overlay/files/src/design/contrast.v403.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "750dd1bb000dec0035e6b0763b6c7b9c120040d5fc0a19ca0a3a68d6acc02f85"
variables: []
secrets_allowed: false
```
````ts
// AUTHORED verification glue implementing W3C relative luminance.
export function contrast(foreground: string, background: string): number {
    const luminance = (hex: string) => {
        if (!/^#[a-f0-9]{6}$/i.test(hex))
            throw new Error("Six-digit sRGB required");
        const rgb = [1, 3, 5].map(i => parseInt(hex.slice(i, i + 2), 16) / 255).map(c => c <= .04045 ? c / 12.92 : ((c + .055) / 1.055) ** 2.4);
        return rgb[0]! * 0.2126 + rgb[1]! * 0.7152 + rgb[2]! * 0.0722;
    };
    const a = luminance(foreground), b = luminance(background);
    return (Math.max(a, b) + .05) / (Math.min(a, b) + .05);
}
````

### FILE: `experience_overlay/files/src/design/experience.v403.test.ts`
```yaml
block_id: "TS-DESIGN-SYSTEM-V403:experience_overlay/files/src/design/experience.v403.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "36a5d3a46f9e56376e650d3dd9a47e4b3c8d458c523e999bd169a9daabde1420"
variables: []
secrets_allowed: false
```
````ts
import { describe, it, expect } from "vitest";
import { readFileSync } from "node:fs";
import { contrast } from "./contrast.v403";
import tokens from "./tokens.v403.json";
const pairs = [["ink", "panel"], ["ink", "surface"], ["muted", "panel"], ["muted", "surface"], ["muted", "muted-bg"], ["on-accent", "accent"], ["on-accent", "accent-strong"], ["accent", "panel"], ["accent", "surface"], ["success", "success-bg"], ["warning", "warning-bg"], ["danger", "danger-bg"], ["ink", "muted-bg"], ["ink", "ink-soft"], ["warning", "panel"]];
describe("V403 actual text pairs", () => {
    for (const [brand, delta] of Object.entries(tokens.brands)) {
        const t = { ...tokens.default, ...delta } as Record<string, string>;
        for (const [a, b] of pairs)
            it(`${brand}: ${a}/${b} >=4.5`, () => expect(contrast(t[a!]!, t[b!]!)).toBeGreaterThanOrEqual(4.5));
    }
    it("retains the exact historical defect as a regression oracle", () => { expect(contrast("#dc2626", "#fee2e2")).toBeCloseTo(3.9534, 3); expect(contrast(tokens.default.danger, tokens.default["danger-bg"])).toBeGreaterThan(4.5); });
    it("the generated page stylesheet contains the selected tokens", () => {
        const css = readFileSync("src/app/globals.css", "utf8");
        for (const [k, v] of Object.entries(tokens.default))
            expect(css).toContain(`--${k}: ${v};`);
        expect(css).not.toContain("#DC2626");
    });
});
// WCAG non-text controls/focus use 3:1, independently from the 4.5:1 text rule.
describe("V403 selected control and focus boundaries", () => { for (const [brand, delta] of Object.entries(tokens.brands)) {
    const values = { ...tokens.default, ...delta } as Record<string, string>;
    for (const token of ["line", "focus"])
        for (const background of ["surface", "panel"])
            it(`${brand} ${token}/${background} >=3`, () => expect(contrast(values[token]!, values[background]!)).toBeGreaterThanOrEqual(3));
} });
````

### FILE: `experience_overlay/files/src/components/experience-ui.tsx`
```yaml
block_id: "TS-DESIGN-SYSTEM-V403:experience_overlay/files/src/components/experience-ui.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "09389bb2f5048b758f9bb34e9ae638d75d76aefbe590be2eb52377d9f35a4f07"
variables: []
secrets_allowed: false
```
````tsx
"use client";
import { usePrivateI18n } from "@/platform/i18n/private-provider";
import { experienceText } from "@/platform/experience/messages";
// AUTHORED presentation glue; native HTML semantics, existing React dependency pin.
import type { ReactNode, ButtonHTMLAttributes } from "react";
export type Tone = "success" | "warning" | "danger" | "muted";
export function StatusBadge({ tone = "muted", children }: {
    tone?: Tone;
    children: ReactNode;
}) { return <span className={`badge badge-${tone}`}>{children}</span>; }
export function Button({ children, className = "", ...props }: ButtonHTMLAttributes<HTMLButtonElement>) { return <button className={`button ${className}`} {...props}>{children}</button>; }
export type ExperienceState = "loading" | "empty" | "error" | "forbidden" | "slow" | "uncertain";
const states = { loading: ["Cargando información", "Estamos consultando tus registros."], empty: ["Todavía no hay resultados", "Probá otro filtro o iniciá una nueva operación."], error: ["No pudimos cargar la información", "Tus cambios no se descartaron. Volvé a consultar."], forbidden: ["Esta operación requiere otro permiso", "Pedí acceso al responsable de tu organización."], slow: ["La consulta está tardando", "Podés seguir trabajando en otra tarea. Conservamos tu información."], uncertain: ["Necesitamos comprobar el resultado", "La operación puede haberse registrado. Consultá su estado antes de volver a enviarla."] } as const;
export function StatePanel({ state, onRetry }: {
    state: ExperienceState;
    onRetry?: () => void;
}) {
    const { locale } = usePrivateI18n();
    const tx = (text: string) => experienceText(text, locale.language);
    const [title, body] = states[state];
    return <section className="card" aria-live="polite" aria-busy={state === "loading" || state === "slow"} role={state === "error" || state === "forbidden" ? "alert" : "status"}>
    <h3>{tx(title)}</h3>
    <p>{tx(body)}</p>{onRetry && ["error", "uncertain", "slow"].includes(state) ? <button onClick={onRetry} type="button">{tx("Consultar estado")}</button> : null}</section>;
}
````

