# Firefox lifecycle and continuity audit V325

2026-09-08. Corrective library maintenance from checkpoint79; T2801/T2803/T2804/T2808. No product expansion or release promotion.

## Browser result and honest limit

FAIL423 remains DIAGNOSED. The preserved V294 log reports an exception while Browser.removeBrowserContext closes a context after assertions. The exact installed Firefox153.0/revision1538 still contains the corresponding SessionStore access without a missing-window guard. Its packaged Juggler handler calls context destruction, which removes the contextual identity and closes pages. This establishes the code path, not the scheduling cause of the historical failure.

The existing Playwright1.62.1 and Node24.20.0 ran48bounded context-lifecycle cases (empty context, immediate page, loaded page, two pages) and20cold browser launches with tracing/video and one/two pages. All68passed with retries0; each suite had a120second deadline. No sleeps, suppression, browser updates or production changes. A finite non-reproduction does not close the original intermittent defect. Reopen targeted work on a reproducible failing trace or an upstream fix tied to the exact browser artifact, then test the original lifecycle and connected journey.

The original Firefox trace path printed by the V294 log was not found under that stage; its current trace inventory contains a WebKit trace only. The preserved text log still proves the failure. Do not claim the old Firefox trace was inspected or infer that missing artifacts prove a fix.

Official source references consulted: [Microsoft pinned Juggler implementation](https://github.com/microsoft/playwright/blob/26a9e470a7b3c7822084b09fb7f13902c5f37b51/browser_patches/firefox/juggler/TargetRegistry.js) and [Mozilla SessionStore source](https://searchfox.org/mozilla-central/source/browser/components/sessionstore/SessionStore.sys.mjs). The Mozilla page is contextual/current source, not an immutable identity of the installed runtime. Packaged bytes were read directly from the two local omni.ja archives and hashed separately; no installed browser file was patched.

## Cursor reconciliation, not another product fix

FAIL457 is already REGRESSION_PROVEN with the inventoried sender family closed in V312. Its20current portal files reconstruct byte-identically to the V321 tested composition; the obsolete generic sender/message is absent from the current component. Therefore the cursor must stop describing that repaired family as open. T2804 remains pending for its wider scope, including multi-tab and downstream execution; no task is completed by this reconciliation. Existing native FAIL532 remains explicitly open.

## Tool availability

Using previously observed isolated paths resolves .NET10.0.400 and psql18.6. The current routing probe sees7/8tools; Docker remains unresolved. Node24.20.0 and pnpm11.25.0 are selected explicitly. Version availability is not SDK dependency/security admission or proof of a running database/container; no installations, paid resources or global PATH changes occurred. Full Preflight must use these same paths after checkpointing the changed owners.

## Local receipt

```json
{
  "probe": {
    "kind": "result",
    "completed": 48,
    "failed": 0,
    "status": "PASS",
    "claim": "bounded context lifecycle probe, not proof that historical Firefox race is fixed"
  },
  "cold_probe": {
    "kind": "result",
    "completed": 20,
    "failed": 0,
    "retries": 0,
    "status": "PASS",
    "claim": "cold-start/tracing/video lifecycle sample; no historical-race fix claim"
  },
  "identity": {
    "kind": "identity",
    "node": "v24.20.0",
    "node_sha256": "5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5",
    "playwright": "1.62.1",
    "firefox_path": "C:\\Users\\NL\\AppData\\Local\\ms-playwright\\firefox-1538\\firefox\\firefox.exe",
    "firefox_sha256": "9571c9a49d6d4de3e4c3b14e2e4bd6521e5b879f4ab6c95143e2098cb66e6618",
    "budget": {
      "iterations": 12,
      "scenarios": 4,
      "timeout_ms": 120000
    },
    "retries": 0
  },
  "toolchains": [
    {
      "id": "pwsh-7",
      "required_for": "library",
      "path": "C:\\Users\\NL\\.cache\\codex-runtimes\\codex-primary-runtime\\dependencies\\native\\powershell\\pwsh.exe",
      "version": "PowerShell 7.6.5",
      "available": true
    },
    {
      "id": "python-3.12+",
      "required_for": "quality, licenses, document intelligence",
      "path": "C:\\Python314\\python.exe",
      "version": "Python 3.14.4",
      "available": true
    },
    {
      "id": "go-1.26.7",
      "required_for": "portable backend foundation",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v281-30cf0957089d499db7bf9b2d1669b7db\\official-toolchain\\go\\bin\\go.exe",
      "version": "go version go1.26.7 windows/amd64",
      "available": true
    },
    {
      "id": "dotnet-10+",
      "required_for": "WCF/SOAP client generation lane",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-dotnet-v127\\dotnet.exe",
      "version": "10.0.400",
      "available": true
    },
    {
      "id": "node-24+",
      "required_for": "web/BFF foundation",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v319-952f1100edd1424f842acd958e3e1f75\\node-only\\node.exe",
      "version": "v24.20.0",
      "available": true
    },
    {
      "id": "pnpm-11.25.0",
      "required_for": "web/BFF frozen install",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077\\pnpm-wrapper.ps1",
      "version": "11.25.0",
      "available": true
    },
    {
      "id": "docker",
      "required_for": "PostgreSQL integration/recovery and optional Business Central Windows-container gates",
      "path": null,
      "version": "",
      "available": false
    },
    {
      "id": "psql",
      "required_for": "direct PostgreSQL migration/restore gates",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-franchise-build-f3200751e6f644c9ae8c08a62d6b14e5\\postgresql-18.6\\pgsql\\bin\\psql.exe",
      "version": "psql (PostgreSQL) 18.6",
      "available": true
    }
  ],
  "portal_parity": [
    {
      "path": "docs/whatsapp-status-operations.md",
      "sha256": "5619c6d4d1791107592dccd99334d3758039d7e7200b5672fa2299036ece62a8"
    },
    {
      "path": "src/components/appointment-notification-status.tsx",
      "sha256": "8a1942f6968e4abbc4718686d212f5d0b1f6812c178e117484958ddfdf250678"
    },
    {
      "path": "src/components/customer-appointment-actions.tsx",
      "sha256": "296457c4bef4d7fbea3835678025330640326228df7c0fa3fafea4506659e261"
    },
    {
      "path": "src/components/customer-handover-actions.tsx",
      "sha256": "c1f264d9a83d58f7c34976e8479b107ef8b1d472357123ebd97c9b87f380e5af"
    },
    {
      "path": "src/components/customer-quote-actions.tsx",
      "sha256": "ace08dab683e987f925382351e98c9b9bc363aaac01be3bf0075a1c0c5616e46"
    },
    {
      "path": "src/components/franchise-command-panel.tsx",
      "sha256": "5424f0ef08974699bc8b55cdeba2a3fde2fc0dceb44b43e432c3d5c3644eb9fb"
    },
    {
      "path": "src/platform/notifications/status-contract.test.ts",
      "sha256": "1c6b7eacbb277fe3c66e635852c9ffbd6c0b65a98e8892e0028d145339464f62"
    },
    {
      "path": "src/platform/notifications/status-contract.ts",
      "sha256": "501130c6f091a1e007d4c6af585cb8acb0ea2fab95f2eadc786329db41ca3fd7"
    },
    {
      "path": "src/app/franchise/page.tsx",
      "sha256": "f6e1121a6bbe33d98f6798aed24bb1f2e915d522cde9128eca01f090c7b0ef54"
    },
    {
      "path": "src/app/locations/appointment-form.tsx",
      "sha256": "978443bed71c02c300f2daf9d708b2999a0575dcdc865e8e8a7b4e725a69a861"
    },
    {
      "path": "src/app/locations/page.tsx",
      "sha256": "2ae4cc75528c3d6cb50c57c3d25331ad26f331c1525fe52bd365f9677cf68680"
    },
    {
      "path": "src/app/customer/appointments/page.tsx",
      "sha256": "f4f68ad8c458a4107332174300ed52189b7d2586ad083707e3a351128803e9a0"
    },
    {
      "path": "src/app/customer/handovers/page.tsx",
      "sha256": "9e889f9da062581429fb7c7e0903036b2967ec461a4ad8588e633b04d5882ea7"
    },
    {
      "path": "src/app/customer/quotes/page.tsx",
      "sha256": "a776ade4087bb40d8db326dd411864f22f1acf98704ed29376d2ff500c8390aa"
    },
    {
      "path": "src/app/api/enterprise/appointments/route.test.ts",
      "sha256": "f3a01f762c855aafdaef85a386374f94e944cf8fdfe9640fa85b5da5b7c8b53c"
    },
    {
      "path": "src/app/api/enterprise/appointments/route.ts",
      "sha256": "78e6effcc7de35e9ae3fcb00599666609c4429baa0fe95a75664765eef00f120"
    },
    {
      "path": "src/app/api/enterprise/franchise/commands/route.test.ts",
      "sha256": "5e190d894b1e396b6d99e86e321b2216cc2cf64b24c21ea7b1f160b5fe0f086e"
    },
    {
      "path": "src/app/api/enterprise/franchise/commands/route.ts",
      "sha256": "c48ac2f8c73b027b93e24e1e26ff0701309ba53993c8156a20d62519cebf9f56"
    },
    {
      "path": "src/app/api/enterprise/franchise/notifications/route.test.ts",
      "sha256": "b1cac0a7bc3605e88dab51b33d0f044495de55cd5de52079cdba1b7afdf25238"
    },
    {
      "path": "src/app/api/enterprise/franchise/notifications/route.ts",
      "sha256": "ebfe5ab0ca8a45c3601666e2248e894bf1acd7bbf2a7a61e0c8cf813be438455"
    }
  ],
  "source_inspection": [
    {
      "archive": "C:\\Users\\NL\\AppData\\Local\\ms-playwright\\firefox-1538\\firefox\\omni.ja",
      "entry": "chrome/juggler/content/TargetRegistry.js",
      "sha256": "35db3ad21a5267019ea27cf9191ff60fd8658b7151ec0d24083566a63cf6e278",
      "matches": []
    },
    {
      "archive": "C:\\Users\\NL\\AppData\\Local\\ms-playwright\\firefox-1538\\firefox\\omni.ja",
      "entry": "chrome/juggler/content/protocol/BrowserHandler.js",
      "sha256": "8963a31c877509cde244bfb978910775c0d8bef0cf7d0a510aca12727a2bcfa5",
      "matches": [
        {
          "line": 67,
          "snippet": "    return {browserContextId: browserContext.browserContextId};\n  }\n\n  async ['Browser.removeBrowserContext']({browserContextId}) {\n    if (!this._enabled)\n      throw new Error('Browser domain is not enabled');\n    await this._targetRegistry.browserContextForId(browserContextId).destroy();\n    this._createdBrowserContextIds.delete(browserContextId);\n  }\n\n  dispose() {\n    helper.removeListeners(this._eventListeners);\n    for (const [target, session] of this._attachedSessions)\n      this._dispatcher.destroySession(session);\n    this._attachedSessions.clear();\n    for (const browserContextId of this._createdBrowserContextIds) {\n      const browserContext = this._targetRegistry.browserContextForId(browserContextId);\n      if (browserContext.removeOnDetach)\n        browserContext.destroy();\n    }\n    this._createdBrowserContextIds.clear();\n  }\n\n  _shouldAttachToTarget(target) {\n    if (this._createdBrowserContextIds.has(target._browserContext.browserContextId))\n      return true;\n    return this._attachToDefaultContext && target._browserContext === this._targetRegistry.defaultContext();\n  }"
        }
      ]
    },
    {
      "archive": "C:\\Users\\NL\\AppData\\Local\\ms-playwright\\firefox-1538\\firefox\\browser\\omni.ja",
      "entry": "modules/sessionstore/SessionStore.sys.mjs",
      "sha256": "36ada4d858c1a9b3e0eca482b4b45d464f84f991cb151536058e5d19a93c1b10",
      "matches": [
        {
          "line": 646,
          "snippet": "    return SessionStoreInternal.getClosedWindowData();\n  },\n\n  maybeDontRestoreTabs(aWindow) {\n    SessionStoreInternal.maybeDontRestoreTabs(aWindow);\n  },\n\n  undoCloseWindow: function ss_undoCloseWindow(aIndex) {\n    return SessionStoreInternal.undoCloseWindow(aIndex);\n  },\n\n  forgetClosedWindow: function ss_forgetClosedWindow(aIndex) {\n    return SessionStoreInternal.forgetClosedWindow(aIndex);\n  },\n\n  getCustomWindowValue(aWindow, aKey) {\n    return SessionStoreInternal.getCustomWindowValue(aWindow, aKey);\n  },\n\n  setCustomWindowValue(aWindow, aKey, aStringValue) {\n    SessionStoreInternal.setCustomWindowValue(aWindow, aKey, aStringValue);\n  },\n\n  deleteCustomWindowValue(aWindow, aKey) {\n    SessionStoreInternal.deleteCustomWindowValue(aWindow, aKey);\n  },\n\n  getCustomTabValue(aTab, aKey) {"
        },
        {
          "line": 647,
          "snippet": "  },\n\n  maybeDontRestoreTabs(aWindow) {\n    SessionStoreInternal.maybeDontRestoreTabs(aWindow);\n  },\n\n  undoCloseWindow: function ss_undoCloseWindow(aIndex) {\n    return SessionStoreInternal.undoCloseWindow(aIndex);\n  },\n\n  forgetClosedWindow: function ss_forgetClosedWindow(aIndex) {\n    return SessionStoreInternal.forgetClosedWindow(aIndex);\n  },\n\n  getCustomWindowValue(aWindow, aKey) {\n    return SessionStoreInternal.getCustomWindowValue(aWindow, aKey);\n  },\n\n  setCustomWindowValue(aWindow, aKey, aStringValue) {\n    SessionStoreInternal.setCustomWindowValue(aWindow, aKey, aStringValue);\n  },\n\n  deleteCustomWindowValue(aWindow, aKey) {\n    SessionStoreInternal.deleteCustomWindowValue(aWindow, aKey);\n  },\n\n  getCustomTabValue(aTab, aKey) {\n    return SessionStoreInternal.getCustomTabValue(aTab, aKey);"
        },
        {
          "line": 4916,
          "snippet": "    closedWinData.groups = Cu.cloneInto(abbreviatedGroups, {});\n  },\n\n  maybeDontRestoreTabs(aWindow) {\n    // Don't restore the tabs if we restore the session at startup\n    this._windows[aWindow.__SSi]._maybeDontRestoreTabs = true;\n  },\n\n  isLastRestorableWindow() {\n    return (\n      Object.values(this._windows).filter(winData => !winData.isPrivate)\n        .length == 1 &&\n      !this._closedWindows.some(win => win._shouldRestore || false)\n    );\n  },\n\n  undoCloseWindow: function ssi_undoCloseWindow(aIndex) {\n    if (!(aIndex in this._closedWindows)) {\n      throw Components.Exception(\n        \"Invalid index: not in the closed windows\",\n        Cr.NS_ERROR_INVALID_ARG\n      );\n    }\n    // reopen the window\n    let state = { windows: this._removeClosedWindow(aIndex) };\n    delete state.windows[0].closedAt; // Window is now open.\n\n    // If any saved tab groups are in the closed window, convert the saved tab"
        }
      ]
    }
  ]
}
```

## Reproduction sources

These are local synthetic diagnostic programs, not reusable product implementations. Use the existing pinned tool paths or resolve them from the actual project toolchain record.

### probe.cjs

```javascript
const {createRequire}=require('node:module');
const {resolve,join}=require('node:path');
const fs=require('node:fs'); const crypto=require('node:crypto');
const req=createRequire('<LOCALAPPDATA>/Temp/elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077/rebuilt/microsoft_playwright_browser_gate/package.json');
const {firefox}=req('@playwright/test');
const sha=p=>crypto.createHash('sha256').update(fs.readFileSync(p)).digest('hex');
const output=join(__dirname,'probe-results.jsonl'); const fd=fs.openSync(output,'wx');
const record=x=>fs.writeSync(fd,JSON.stringify(x)+'\n');
const deadline=setTimeout(()=>{record({status:'DEADLINE_EXCEEDED'});process.exit(2)},120000);
(async()=>{
  record({kind:'identity',node:process.version,node_sha256:sha(process.execPath),playwright: req('@playwright/test/package.json').version,firefox_path:firefox.executablePath(),firefox_sha256:sha(firefox.executablePath()),budget:{iterations:12,scenarios:4,timeout_ms:120000},retries:0});
  let failed=0,completed=0;
  for(const scenario of ['empty-context','immediate-page','loaded-page','two-pages']){
    const browser=await firefox.launch({headless:true,timeout:20000});record({kind:'browser',scenario,version:browser.version()});
    try{
      for(let i=0;i<12;i++){
        const start=performance.now();let context;
        try {
          context=await browser.newContext();
          if(scenario!=='empty-context') {
            const page=await context.newPage();
            if(scenario==='loaded-page'||scenario==='two-pages') {
              await page.setContent('<main><h1>Fixture local</h1><button>Aceptar</button></main>');
              if(await page.getByRole('heading',{name:'Fixture local'}).textContent()!=='Fixture local')throw new Error('SEMANTIC_ASSERTION_FAILED');
            }
            if(scenario==='two-pages') await context.newPage();
          }
          await context.close();
          completed++;record({scenario,iteration:i,status:'PASS',elapsed_ms:performance.now()-start});
        }catch(error){failed++;record({scenario,iteration:i,status:'FAIL',error:String(error),stack:error.stack,elapsed_ms:performance.now()-start});break;}
      }
    }finally{try{await browser.close()}catch(error){failed++;record({scenario,status:'BROWSER_CLOSE_FAIL',error:String(error)})}}
  }
  record({kind:'result',completed,failed,status:failed?'FAIL':'PASS',claim:'bounded context lifecycle probe, not proof that historical Firefox race is fixed'});
  process.exitCode=failed?1:0;
})().catch(error=>{record({status:'HARNESS_FAIL',error:String(error),stack:error.stack});process.exitCode=2}).finally(()=>{clearTimeout(deadline);fs.closeSync(fd)});
```

### cold-probe.cjs

```javascript
const {createRequire}=require('node:module');const {join}=require('node:path');const fs=require('node:fs');
const req=createRequire('<LOCALAPPDATA>/Temp/elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077/rebuilt/microsoft_playwright_browser_gate/package.json');
const {firefox}=req('@playwright/test');
const fd=fs.openSync(join(__dirname,'cold-probe-results.jsonl'),'wx');const record=x=>fs.writeSync(fd,JSON.stringify(x)+'\n');
const deadline=setTimeout(()=>{record({status:'DEADLINE_EXCEEDED'});process.exit(2)},120000);
(async()=>{let completed=0,failed=0;
 for(let i=0;i<20;i++){
  let browser;const start=performance.now();
  try{
   browser=await firefox.launch({headless:true,timeout:20000});
   const context=await browser.newContext({recordVideo:{dir:join(__dirname,'cold-artifacts',String(i))}});
   await context.tracing.start({screenshots:true,snapshots:true,sources:true});
   const page=await context.newPage();await page.setContent('<main><h1>Reinicio sintético</h1></main>');
   if(await page.getByRole('heading').textContent()!=='Reinicio sintético')throw new Error('SEMANTIC_ASSERTION_FAILED');
   if(i%2===1){const other=await context.newPage();await other.setContent('<main>Segunda pestaña</main>')}
   await context.tracing.stop({path:join(__dirname,'cold-artifacts',String(i),'trace.zip')});
   await context.close();completed++;record({iteration:i,pages:i%2+1,status:'PASS',elapsed_ms:performance.now()-start});
  }catch(error){failed++;record({iteration:i,status:'FAIL',error:String(error),stack:error.stack});break}
  finally{if(browser)try{await browser.close()}catch(error){failed++;record({iteration:i,status:'BROWSER_CLOSE_FAIL',error:String(error)})}}
 }
 record({kind:'result',completed,failed,retries:0,status:failed?'FAIL':'PASS',claim:'cold-start/tracing/video lifecycle sample; no historical-race fix claim'});process.exitCode=failed?1:0;
})().catch(error=>{record({status:'HARNESS_FAIL',error:String(error)});process.exitCode=2}).finally(()=>{clearTimeout(deadline);fs.closeSync(fd)});
```

## Evidence hashes

- probe-results.jsonl: 44efca1215fc66bb7de43711c60f0fb886d94d1a745a07c29144c68e56b9cd7b
- cold-probe-results.jsonl: 2167fd23c97af6f74d9415d87f3de9952ec3a8eab5729b0248955bfbd7efd0f7
- toolchain.json: e3101de745d259a20975188eeb136af3f3416b9ef12c09b00398caf5f1747db0
- runtime-source-inspection.json: 5c4c0eacbd15ee38f024407dee3158df78584aa2b0cbb2be73a31a23cce0d27e
- audit-receipt.json: 797b224defab033eadd6f2a50af4e8848b8811c6b8e8f59c33f7db9f68057c60

## Full Preflight after checkpoint80

Actual process exit0 is informational; result remains BLOCKED because Docker is unresolved. All151audit steps PASS, including full structural reconstruction161packs/1453files/760Markdown/52profiles and franchise67/746. Explicit isolated .NET/psql paths eliminate two path-resolution gaps in this invocation; no install or production/runtime admission. Readiness still BLOCKED42 and native FAIL532/Firefox423 remain open.

```json
{
  "status": "BLOCKED",
  "steps_passed": 151,
  "missing_tools": [
    "docker"
  ],
  "source_receipt_sha256": "b263115c26320ecf3ffe14e364c949634dd62d66333b6b3b0312544806cf0bac",
  "log_sha256": "70bf1716e9b3ef6540e67dc1e4aa17c625d9988808c8a9d7859b6808d7e6f4f4",
  "structural_inventory": {
    "packs": 161,
    "materialized_files": 1453,
    "markdown_files": 760,
    "profiles": 52,
    "franchise_packs": 67,
    "franchise_files": 746
  },
  "toolchains": [
    {
      "id": "pwsh-7",
      "required_for": "library",
      "path": "C:\\Users\\NL\\.cache\\codex-runtimes\\codex-primary-runtime\\dependencies\\native\\powershell\\pwsh.exe",
      "version": "PowerShell 7.6.5",
      "available": true
    },
    {
      "id": "python-3.12+",
      "required_for": "quality, licenses, document intelligence",
      "path": "C:\\Python314\\python.exe",
      "version": "Python 3.14.4",
      "available": true
    },
    {
      "id": "go-1.26.7",
      "required_for": "portable backend foundation",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v281-30cf0957089d499db7bf9b2d1669b7db\\official-toolchain\\go\\bin\\go.exe",
      "version": "go version go1.26.7 windows/amd64",
      "available": true
    },
    {
      "id": "dotnet-10+",
      "required_for": "WCF/SOAP client generation lane",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-dotnet-v127\\dotnet.exe",
      "version": "10.0.400",
      "available": true
    },
    {
      "id": "node-24+",
      "required_for": "web/BFF foundation",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v319-952f1100edd1424f842acd958e3e1f75\\node-only\\node.exe",
      "version": "v24.20.0",
      "available": true
    },
    {
      "id": "pnpm-11.25.0",
      "required_for": "web/BFF frozen install",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077\\pnpm-wrapper.ps1",
      "version": "11.25.0",
      "available": true
    },
    {
      "id": "docker",
      "required_for": "PostgreSQL integration/recovery and optional Business Central Windows-container gates",
      "path": null,
      "version": "",
      "available": false
    },
    {
      "id": "psql",
      "required_for": "direct PostgreSQL migration/restore gates",
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-franchise-build-f3200751e6f644c9ae8c08a62d6b14e5\\postgresql-18.6\\pgsql\\bin\\psql.exe",
      "version": "psql (PostgreSQL) 18.6",
      "available": true
    }
  ],
  "steps": [
    {
      "id": "library-structural-and-materialization",
      "status": "PASS",
      "elapsed_ms": 64997.0
    },
    {
      "id": "project-readiness-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 939.0
    },
    {
      "id": "project-readiness-syntax-and-70-regressions",
      "status": "PASS",
      "elapsed_ms": 3028.0
    },
    {
      "id": "engineering-execution-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 753.0
    },
    {
      "id": "engineering-execution-syntax-and-24-regressions",
      "status": "PASS",
      "elapsed_ms": 822.0
    },
    {
      "id": "capability-gap-resolution-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 823.0
    },
    {
      "id": "capability-gap-resolution-syntax-and-nine-regressions",
      "status": "PASS",
      "elapsed_ms": 1558.0
    },
    {
      "id": "microsoft-devskim-adapted-sast-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 814.0
    },
    {
      "id": "microsoft-devskim-adapted-sast-static-and-negatives",
      "status": "PASS",
      "elapsed_ms": 942.0
    },
    {
      "id": "gitlab-opengrep-signed-sast-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 845.0
    },
    {
      "id": "gitlab-opengrep-signed-sast-static-and-negative",
      "status": "PASS",
      "elapsed_ms": 1726.0
    },
    {
      "id": "portable-signed-release-evidence-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 1026.0
    },
    {
      "id": "portable-signed-release-evidence-gate-vector-and-tamper",
      "status": "PASS",
      "elapsed_ms": 1095.0
    },
    {
      "id": "official-upstream-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 1036.0
    },
    {
      "id": "official-upstream-lock-negative-suite",
      "status": "PASS",
      "elapsed_ms": 9761.0
    },
    {
      "id": "official-source-profile-suite",
      "status": "PASS",
      "elapsed_ms": 4393.0
    },
    {
      "id": "official-opaque-artifact-transport-suite",
      "status": "PASS",
      "elapsed_ms": 5959.0
    },
    {
      "id": "official-upstream-lock-124",
      "status": "PASS",
      "elapsed_ms": 992.0
    },
    {
      "id": "microsoft-kiota-openapi-client-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 828.0
    },
    {
      "id": "microsoft-kiota-openapi-client-gate-source-contract",
      "status": "PASS",
      "elapsed_ms": 32.0
    },
    {
      "id": "go-native-fuzz-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 822.0
    },
    {
      "id": "go-native-fuzz-gate-source-contract",
      "status": "PASS",
      "elapsed_ms": 19.0
    },
    {
      "id": "microsoft-playwright-browser-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 856.0
    },
    {
      "id": "microsoft-playwright-browser-gate-source-contract",
      "status": "PASS",
      "elapsed_ms": 27.0
    },
    {
      "id": "microsoft-playwright-browser-gate-frozen-offline-install",
      "status": "PASS",
      "elapsed_ms": 741.0
    },
    {
      "id": "microsoft-playwright-browser-gate-runtime",
      "status": "PASS",
      "elapsed_ms": 9607.0
    },
    {
      "id": "google-lighthouse-web-quality-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 927.0
    },
    {
      "id": "google-lighthouse-web-quality-gate-source-contract",
      "status": "PASS",
      "elapsed_ms": 22.0
    },
    {
      "id": "google-lighthouse-web-quality-gate-frozen-offline-install",
      "status": "PASS",
      "elapsed_ms": 6066.0
    },
    {
      "id": "google-lighthouse-web-quality-gate-runtime-contract",
      "status": "PASS",
      "elapsed_ms": 28644.0
    },
    {
      "id": "aws-powertools-idempotent-sqs-batch-materialization",
      "status": "PASS",
      "elapsed_ms": 880.0
    },
    {
      "id": "aws-powertools-idempotent-sqs-batch-contracts",
      "status": "PASS",
      "elapsed_ms": 106.0
    },
    {
      "id": "aws-lambda-durable-execution-materialization",
      "status": "PASS",
      "elapsed_ms": 749.0
    },
    {
      "id": "aws-lambda-durable-execution-contracts",
      "status": "PASS",
      "elapsed_ms": 99.0
    },
    {
      "id": "aws-textractor-official-materialization",
      "status": "PASS",
      "elapsed_ms": 684.0
    },
    {
      "id": "aws-textractor-official-contracts",
      "status": "PASS",
      "elapsed_ms": 92.0
    },
    {
      "id": "transparency-dev-tessera-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 774.0
    },
    {
      "id": "transparency-dev-tessera-static-contracts",
      "status": "PASS",
      "elapsed_ms": 98.0
    },
    {
      "id": "official-document-sdk-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 822.0
    },
    {
      "id": "official-document-sdk-acquisition-suite",
      "status": "PASS",
      "elapsed_ms": 994.0
    },
    {
      "id": "official-document-sdk-lock-14",
      "status": "PASS",
      "elapsed_ms": 660.0
    },
    {
      "id": "microsoft-azure-di-official-invoice-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 702.0
    },
    {
      "id": "microsoft-azure-di-official-invoice-sample-regression",
      "status": "PASS",
      "elapsed_ms": 145.0
    },
    {
      "id": "microsoft-azure-di-official-invoice-fixture-regression",
      "status": "PASS",
      "elapsed_ms": 810.0
    },
    {
      "id": "microsoft-azure-cu-official-invoice-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 698.0
    },
    {
      "id": "microsoft-azure-cu-official-invoice-sample-regression",
      "status": "PASS",
      "elapsed_ms": 312.0
    },
    {
      "id": "microsoft-azure-cu-official-invoice-python-syntax",
      "status": "PASS",
      "elapsed_ms": 78.0
    },
    {
      "id": "microsoft-azure-cu-official-binary-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 686.0
    },
    {
      "id": "microsoft-azure-cu-official-binary-regression",
      "status": "PASS",
      "elapsed_ms": 69.0
    },
    {
      "id": "microsoft-azure-cu-official-binary-syntax-0",
      "status": "PASS",
      "elapsed_ms": 72.0
    },
    {
      "id": "microsoft-azure-cu-official-binary-syntax-1",
      "status": "PASS",
      "elapsed_ms": 78.0
    },
    {
      "id": "microsoft-azure-cu-official-binary-syntax-2",
      "status": "PASS",
      "elapsed_ms": 62.0
    },
    {
      "id": "microsoft-azure-cu-official-analyzer-copy-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 634.0
    },
    {
      "id": "microsoft-azure-cu-official-analyzer-copy-regression",
      "status": "PASS",
      "elapsed_ms": 103.0
    },
    {
      "id": "microsoft-azure-cu-official-analyzer-copy-syntax-0",
      "status": "PASS",
      "elapsed_ms": 73.0
    },
    {
      "id": "microsoft-azure-cu-official-analyzer-copy-syntax-1",
      "status": "PASS",
      "elapsed_ms": 66.0
    },
    {
      "id": "microsoft-azure-cu-official-analyzer-copy-syntax-2",
      "status": "PASS",
      "elapsed_ms": 62.0
    },
    {
      "id": "google-cloud-docai-official-process-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 652.0
    },
    {
      "id": "google-cloud-docai-official-process-regression",
      "status": "PASS",
      "elapsed_ms": 122.0
    },
    {
      "id": "google-cloud-docai-official-custom-extractor-syntax",
      "status": "PASS",
      "elapsed_ms": 77.0
    },
    {
      "id": "google-cloud-docai-official-fixture-regression",
      "status": "PASS",
      "elapsed_ms": 938.0
    },
    {
      "id": "google-cloud-docai-official-validator-compile",
      "status": "PASS",
      "elapsed_ms": 69.0
    },
    {
      "id": "google-document-ai-official-lifecycle-materialization",
      "status": "PASS",
      "elapsed_ms": 662.0
    },
    {
      "id": "google-document-ai-official-lifecycle-static-regression",
      "status": "PASS",
      "elapsed_ms": 133.0
    },
    {
      "id": "dapr-official-outbox-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 925.0
    },
    {
      "id": "dapr-official-outbox-source-regression",
      "status": "PASS",
      "elapsed_ms": 680.0
    },
    {
      "id": "dapr-official-outbox-acquisition-regression",
      "status": "PASS",
      "elapsed_ms": 687.0
    },
    {
      "id": "microsoft-pg-durable-human-handoff-materialization",
      "status": "PASS",
      "elapsed_ms": 922.0
    },
    {
      "id": "microsoft-pg-durable-python-syntax",
      "status": "PASS",
      "elapsed_ms": 102.0
    },
    {
      "id": "microsoft-avm-secure-sftp-intake-materialization",
      "status": "PASS",
      "elapsed_ms": 951.0
    },
    {
      "id": "microsoft-avm-secure-sftp-intake-contract",
      "status": "PASS",
      "elapsed_ms": 890.0
    },
    {
      "id": "secure-email-mime-core-materialization",
      "status": "PASS",
      "elapsed_ms": 893.0
    },
    {
      "id": "secure-email-mime-core-contract",
      "status": "PASS",
      "elapsed_ms": 1164.0
    },
    {
      "id": "aws-ses-immutable-email-receiver-materialization",
      "status": "PASS",
      "elapsed_ms": 874.0
    },
    {
      "id": "aws-ses-immutable-email-receiver-contract",
      "status": "PASS",
      "elapsed_ms": 1399.0
    },
    {
      "id": "aws-guardduty-immutable-release-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 848.0
    },
    {
      "id": "aws-guardduty-immutable-release-gate-contract",
      "status": "PASS",
      "elapsed_ms": 1061.0
    },
    {
      "id": "aws-guardduty-magika-idp-dispatch-gate-materialization",
      "status": "PASS",
      "elapsed_ms": 971.0
    },
    {
      "id": "aws-guardduty-magika-idp-dispatch-gate-contract",
      "status": "PASS",
      "elapsed_ms": 1222.0
    },
    {
      "id": "aws-idp-immutable-evaluation-handoff-materialization",
      "status": "PASS",
      "elapsed_ms": 776.0
    },
    {
      "id": "aws-idp-immutable-evaluation-handoff-contract",
      "status": "PASS",
      "elapsed_ms": 935.0
    },
    {
      "id": "aws-idp-evaluation-decision-worker-materialization",
      "status": "PASS",
      "elapsed_ms": 850.0
    },
    {
      "id": "aws-idp-evaluation-decision-worker-contract",
      "status": "PASS",
      "elapsed_ms": 978.0
    },
    {
      "id": "aws-idp-postgres-persistence-boundary-materialization",
      "status": "PASS",
      "elapsed_ms": 812.0
    },
    {
      "id": "aws-idp-postgres-persistence-boundary-contract",
      "status": "PASS",
      "elapsed_ms": 976.0
    },
    {
      "id": "debezium-postgres-outbox-runtime-materialization",
      "status": "PASS",
      "elapsed_ms": 747.0
    },
    {
      "id": "debezium-postgres-outbox-runtime-contract",
      "status": "PASS",
      "elapsed_ms": 1080.0
    },
    {
      "id": "debezium-postgres-inbox-consumer-materialization",
      "status": "PASS",
      "elapsed_ms": 904.0
    },
    {
      "id": "debezium-postgres-inbox-consumer-contract",
      "status": "PASS",
      "elapsed_ms": 783.0
    },
    {
      "id": "google-cel-document-mapping-materialization",
      "status": "PASS",
      "elapsed_ms": 764.0
    },
    {
      "id": "official-business-central-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 790.0
    },
    {
      "id": "official-business-central-runtime-suite",
      "status": "PASS",
      "elapsed_ms": 977.0
    },
    {
      "id": "microsoft-markitdown-local-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 966.0
    },
    {
      "id": "microsoft-markitdown-local-profile-contract",
      "status": "PASS",
      "elapsed_ms": 24.0
    },
    {
      "id": "microsoft-markitdown-local-wrapper-unit",
      "status": "PASS",
      "elapsed_ms": 388.0
    },
    {
      "id": "document-pipeline-routing-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 821.0
    },
    {
      "id": "document-pipeline-routing-unit",
      "status": "PASS",
      "elapsed_ms": 432.0
    },
    {
      "id": "secure-local-file-ingestion-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 840.0
    },
    {
      "id": "secure-local-file-ingestion-unit",
      "status": "PASS",
      "elapsed_ms": 581.0
    },
    {
      "id": "microsoft-durable-document-orchestration-materialization",
      "status": "PASS",
      "elapsed_ms": 877.0
    },
    {
      "id": "microsoft-durable-document-orchestration-static-contract",
      "status": "PASS",
      "elapsed_ms": 43.0
    },
    {
      "id": "aws-textract-document-runtime-materialization",
      "status": "PASS",
      "elapsed_ms": 805.0
    },
    {
      "id": "aws-textract-static-lock-and-fail-closed-profile",
      "status": "PASS",
      "elapsed_ms": 23.0
    },
    {
      "id": "aws-enterprise-storage-email-materialization",
      "status": "PASS",
      "elapsed_ms": 809.0
    },
    {
      "id": "aws-enterprise-object-lock-contract",
      "status": "PASS",
      "elapsed_ms": 15.0
    },
    {
      "id": "aws-enterprise-storage-email-go-gates",
      "status": "PASS",
      "elapsed_ms": 19727.0
    },
    {
      "id": "aws-secure-quarantine-intake-materialization",
      "status": "PASS",
      "elapsed_ms": 679.0
    },
    {
      "id": "aws-secure-quarantine-intake-fail-closed-contract",
      "status": "PASS",
      "elapsed_ms": 21.0
    },
    {
      "id": "aws-secure-quarantine-intake-go-gates",
      "status": "PASS",
      "elapsed_ms": 3497.0
    },
    {
      "id": "google-cloud-storage-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 781.0
    },
    {
      "id": "google-cloud-storage-fail-closed-contract",
      "status": "PASS",
      "elapsed_ms": 28.0
    },
    {
      "id": "google-cloud-storage-go-gates",
      "status": "PASS",
      "elapsed_ms": 38274.0
    },
    {
      "id": "amazon-spapi-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 741.0
    },
    {
      "id": "amazon-spapi-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 32.0
    },
    {
      "id": "amazon-easyship-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 796.0
    },
    {
      "id": "amazon-easyship-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 17.0
    },
    {
      "id": "amazon-fulfillment-delivery-evidence-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 798.0
    },
    {
      "id": "amazon-fulfillment-delivery-evidence-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 12.0
    },
    {
      "id": "amazon-supply-sources-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 943.0
    },
    {
      "id": "amazon-supply-sources-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 14.0
    },
    {
      "id": "amazon-mli-inventory-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 953.0
    },
    {
      "id": "amazon-mli-inventory-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 16.0
    },
    {
      "id": "amazon-external-inventory-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 894.0
    },
    {
      "id": "amazon-external-inventory-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 26.0
    },
    {
      "id": "google-merchant-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 928.0
    },
    {
      "id": "google-merchant-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 15.0
    },
    {
      "id": "google-ads-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 794.0
    },
    {
      "id": "google-ads-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 12.0
    },
    {
      "id": "meta-ads-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 755.0
    },
    {
      "id": "meta-ads-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 13.0
    },
    {
      "id": "meta-lead-reconciliation-materialization",
      "status": "PASS",
      "elapsed_ms": 757.0
    },
    {
      "id": "meta-lead-reconciliation-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 12.0
    },
    {
      "id": "tiktok-ads-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 650.0
    },
    {
      "id": "tiktok-ads-hash-lock-contract",
      "status": "PASS",
      "elapsed_ms": 13.0
    },
    {
      "id": "tiktok-lead-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 912.0
    },
    {
      "id": "tiktok-lead-official-contract-and-fail-closed-profile",
      "status": "PASS",
      "elapsed_ms": 26.0
    },
    {
      "id": "meta-whatsapp-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 1163.0
    },
    {
      "id": "meta-whatsapp-source-and-provenance-contract",
      "status": "PASS",
      "elapsed_ms": 16.0
    },
    {
      "id": "meta-whatsapp-offline-unit-contract",
      "status": "PASS",
      "elapsed_ms": 562.0
    },
    {
      "id": "firebase-push-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 908.0
    },
    {
      "id": "firebase-push-source-and-module-contract",
      "status": "PASS",
      "elapsed_ms": 12.0
    },
    {
      "id": "mercadolibre-marketplace-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 858.0
    },
    {
      "id": "mercadolibre-current-http-authority-contract",
      "status": "PASS",
      "elapsed_ms": 15.0
    },
    {
      "id": "azure-blob-evidence-adapter-materialization",
      "status": "PASS",
      "elapsed_ms": 779.0
    },
    {
      "id": "azure-blob-evidence-lock-and-authority-contract",
      "status": "PASS",
      "elapsed_ms": 17.0
    },
    {
      "id": "azure-content-understanding-runtime-materialization",
      "status": "PASS",
      "elapsed_ms": 803.0
    },
    {
      "id": "azure-content-understanding-runtime-lock-contract",
      "status": "PASS",
      "elapsed_ms": 11.0
    },
    {
      "id": "google-document-ai-runtime-materialization",
      "status": "PASS",
      "elapsed_ms": 871.0
    },
    {
      "id": "google-document-ai-runtime-lock-contract",
      "status": "PASS",
      "elapsed_ms": 11.0
    },
    {
      "id": "strict-document-field-evaluation-pack-materialization",
      "status": "PASS",
      "elapsed_ms": 873.0
    },
    {
      "id": "strict-document-field-evaluation-python-syntax",
      "status": "PASS",
      "elapsed_ms": 92.0
    }
  ]
}
```
