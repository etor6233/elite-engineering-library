# Entry latency and checkpoint interruption V326

2026-09-08. Library maintenance from checkpoint81, T2801/T2808 and the checkpoint subset of planned TEST-04. No product expansion, new source/runtime admission or final archive.

## Distance to the agreed objective

Roadmap §10 has one formally completed task (T2800) and ten open tasks with substantial partial work. No effort weights or complete acceptance evidence support a global completion percentage. Pack/file/test counts measure inventory or bounded verification, not fraction of the requested system delivered. The base is substantially built; material integration, preparation, operations and final artifact acceptance remain. The current42readiness diagnostics are conditions, not42equal-sized development tasks. User request to keep advancing preserves the existing business/data/authority boundaries.

## Controlled measurement

The unchanged canonical lifecycle harness ran against the current root and an independent thirteen-file byte-identical test copy. Each run exercised both NEW and EXISTING,55checks,28actual subprocesses and8expected negative exits. One warmup per location was excluded; five measured pairs ran in alternating full/clean order, at concurrency1. All12runs passed. These are repetitions of one55-check suite, not660independent test cases. The thirteen-file copy is sufficient for this test, not a complete distributable library.

Total harness wall time uses Python perf_counter_ns; per-stage figures are sums of the harness's existing subprocess Stopwatches. Entry stages comprise compose, bridge, advisory, blocked-readiness and first checkpoint. They exclude shared bootstrap materialization and parent-side file preparation; the full wall measurement includes all setup, modes, recovery tests and process startup. NEW always precedes EXISTING within the unchanged harness, so the difference is not a causal comparison of modes. OS cache was not flushed; all outliers retained. No p95/p99, SLO, product-build-time or optimization claim from five observations per location. No token, billing or actual agent-comprehension measurement.

| Location / measured phase | n | Minimum ms | Median ms | Maximum ms |
|---|---:|---:|---:|---:|
| full / complete harness | 5 | 9923.5 | 9994.8 | 10473.1 |
| full / NEW entry_subprocess_ms | 5 | 1914.0 | 2019.0 | 2154.0 |
| full / NEW resume_ms | 5 | 112.0 | 120.0 | 149.0 |
| full / EXISTING entry_subprocess_ms | 5 | 2199.0 | 2276.0 | 2387.0 |
| full / EXISTING resume_ms | 5 | 101.0 | 131.0 | 153.0 |
| clean / complete harness | 5 | 9637.8 | 10028.3 | 10785.0 |
| clean / NEW entry_subprocess_ms | 5 | 1918.0 | 2001.0 | 2204.0 |
| clean / NEW resume_ms | 5 | 115.0 | 129.0 | 146.0 |
| clean / EXISTING entry_subprocess_ms | 5 | 2072.0 | 2177.0 | 2292.0 |
| clean / EXISTING resume_ms | 5 | 115.0 | 122.0 | 151.0 |

The two whole-harness medians are nearly identical; this does not establish a speedup. BENCH-01 remains planned for its distinct hypothesis about actual indexed-agent work/context reduction. BENCH-02 records this observed tooling latency and independent-copy equivalence.

## Environment

```json
{
  "cpu": [
    {
      "Name": "Intel(R) Core(TM) Ultra 9 285K",
      "NumberOfCores": 24,
      "NumberOfLogicalProcessors": 24,
      "MaxClockSpeed": 3700
    }
  ],
  "os": {
    "Caption": "Microsoft Windows 11 Pro",
    "Version": "10.0.26200",
    "BuildNumber": "26200"
  },
  "physicalMemoryBytes": 205372039168,
  "powerScheme": "Power Scheme GUID: 381b4222-f694-41f0-9685-ff5bb260df2e  (Balanced)",
  "controls": "No affinity, priority, power, cache, thermal or background-load changes. Firmware/microcode/temperatures not measured."
}
```

## Context inventory boundary

{
  "all_markdown_files": 777,
  "all_markdown_bytes": 16337996,
  "must_read_refs": 12,
  "must_read_owner_bytes": 1236859,
  "claim": "file-byte inventory only; indexed excerpts actually read, tokens, billing and comprehension are not measured"
}

The777workspace Markdown files include local maintenance records excluded by the release inventory; do not confuse this denominator with760portable Markdown at checkpoint81. Owner file bytes are not necessarily read in full and are not token savings.

## Four injected checkpoint failure boundaries

The fault probe uses a fresh15-file EXECUTION-VALIDATOR1.3.0 materialization and synthetic EXISTING fixtures. A child intercepts only the event writer boundary. Two cases raise an injected write error; two terminate the actual child process using os._exit after writing/flushing the specified prefix. Real filesystem bytes and unmodified canonical validation/checkpoint functions are used.

| Boundary | Child outcome | Initial resume | Recovery observed |
|---|---|---|---|
| Before write | expected error78 | rejected2 | normal second checkpoint completes |
| Partial write | expected error78 | rejected2 | corrupt bytes preserved; exact known fixture prefix restored; second checkpoint completes |
| Partial write + process termination | expected exit79 | rejected2 | same bounded recovery with previous event intact |
| Full write + fsync + process termination | expected exit79 | accepted0 | existing second checkpoint recognized; duplicate append rejected without modification |

All four end with exactly two valid events and byte-identical first events and evidence files. Partial-log append attempts reject without mutation. This is not power-loss durability, disk/controller failure, concurrent-writer safety, crash atomicity of materializer/bridge, or automatic repair authority for real logs. Recovery time of these synthetic fixtures is not a project RTO. TEST-04 remains planned for its broader materialization/interruption scope; TEST-45 captures this checkpoint subset.

```json
{
  "schema": "elite-checkpoint-fault-probe/v1",
  "cases": [
    {
      "mode": "before_write",
      "child_exit": 78,
      "resume_after_interrupt": 2,
      "resume_after_recovery": 0,
      "original_event_sha256": "6dc543fbe7c85000f8ed8428063ef2dffa130bd3ba2c538e964a6febd5d4b667",
      "interrupted_log_sha256": "6dc543fbe7c85000f8ed8428063ef2dffa130bd3ba2c538e964a6febd5d4b667",
      "final_log_sha256": "a30f9d5a9b90c58eef51a054f23aed0161fe6e8be02eb9a3fbfe6df903003e51",
      "original_event_byte_exact": true,
      "evidence_files_unchanged": true,
      "final_event_count": 2,
      "elapsed_ms": 356.6929,
      "status": "PASS"
    },
    {
      "mode": "partial_error",
      "child_exit": 78,
      "resume_after_interrupt": 2,
      "resume_after_recovery": 0,
      "original_event_sha256": "f0705e8dfbdd848fff606d37c2b2f783486400431ad26e8fe8489eae692ba0e6",
      "interrupted_log_sha256": "ef9d9694850c193059a6f11c135a677878de2d67f3a0154a973971a80bf28973",
      "final_log_sha256": "088a414647d10113acc22ea4434f3c1b6531f1df503ab144eaa7aecef63affc4",
      "original_event_byte_exact": true,
      "evidence_files_unchanged": true,
      "final_event_count": 2,
      "elapsed_ms": 401.815,
      "status": "PASS"
    },
    {
      "mode": "partial_crash",
      "child_exit": 79,
      "resume_after_interrupt": 2,
      "resume_after_recovery": 0,
      "original_event_sha256": "433f2496b89ffb8716d7e8121c265a8a458c877415874583b40d266d9e08e530",
      "interrupted_log_sha256": "7d19b9b7ce0e4865cc7f74337752e87bab334cefb0581ab3f1bb3bd96b73f565",
      "final_log_sha256": "599b381645f3601136e3232fe960742385f6786a1671fab0a78508e8984c6ff5",
      "original_event_byte_exact": true,
      "evidence_files_unchanged": true,
      "final_event_count": 2,
      "elapsed_ms": 405.9913,
      "status": "PASS"
    },
    {
      "mode": "crash_after_fsync",
      "child_exit": 79,
      "resume_after_interrupt": 0,
      "resume_after_recovery": 0,
      "original_event_sha256": "21c3454c972e404df98f3448e79422836290567cc3855581b9454cfd3a844aec",
      "interrupted_log_sha256": "0b4856af53bfe197f679be3024aecafd841d17cfc421a8217a77e7ba68b585a5",
      "final_log_sha256": "0b4856af53bfe197f679be3024aecafd841d17cfc421a8217a77e7ba68b585a5",
      "original_event_byte_exact": true,
      "evidence_files_unchanged": true,
      "final_event_count": 2,
      "elapsed_ms": 378.699,
      "status": "PASS"
    }
  ],
  "limits": [
    "Synthetic injected write boundary and real child process exit; not power loss, disk/controller failure, concurrent writers or materializer/bridge crash atomicity",
    "No automatic repair of real project events; known-prefix recovery occurred only in owned synthetic fixtures with all damaged bytes retained"
  ]
}
```

## Protocol and source identity

```json
{
  "schema": "elite-entry-benchmark/v1",
  "source_manifest_sha256": "c408021fd795fec454eca8c3efc1ca52303db353ee504b18cbe89f3374fce377",
  "python": "3.14.4",
  "pwsh_sha256": "362a356ce7f0940ec74f73a8fc2c990a2cc24a38a11c90bbd8eca947110ad139",
  "warmups_per_location": 1,
  "samples_per_location": 5,
  "location_order": [
    "full",
    "clean",
    "full",
    "clean",
    "full",
    "clean",
    "full",
    "clean",
    "full",
    "clean",
    "full",
    "clean"
  ],
  "same_tests": true,
  "expected_checks": 55,
  "expected_calls": 28,
  "timeout_per_run_seconds": 90,
  "concurrency": 1,
  "clock": "time.perf_counter_ns + existing Stopwatch per subprocess",
  "cache_policy": "OS cache not flushed; explicit warmups excluded",
  "non_claims": [
    "no product build",
    "no production SLO or speedup claim",
    "no token/billing/model measurement",
    "no portable release admission for the thirteen-file test copy",
    "NEW precedes EXISTING inside the unchanged harness; do not infer causal mode speed difference"
  ]
}
```

```json
[
  {
    "path": "AGENT_SYSTEM_START.md",
    "sha256": "273e5c6a666c73d53facd25ce373eeed7dc141b82d5d4b4e21a963c973420771"
  },
  {
    "path": "AGENTS.md",
    "sha256": "503c64798d0a644c8818f7d0b2cdc9c762c1a7b820721d3bec8511ca3eaa89c7"
  },
  {
    "path": "INSTALL_AGENT_BRIDGE.ps1",
    "sha256": "fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b"
  },
  {
    "path": "materialize_markdown_pack.ps1",
    "sha256": "7ae045ff9b03ea22242d376d16b2a114f00e2f095217789eb92ca1f19e126356"
  },
  {
    "path": "START_ANY_PROJECT.md",
    "sha256": "c2c452e86d77b6e28f07411344690ec205c4fb0d1bf77473420ce5049ed191ab"
  },
  {
    "path": "VERIFY_EXECUTABLE_LIBRARY.ps1",
    "sha256": "f4980a513b4b3733d43809d4adc325502ca17b79fdaac57e2719b89c5a85e894"
  },
  {
    "path": "implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md",
    "sha256": "ab8a4c3fb4dbf39900dad8170a303ebfcfa56f622256e09b6cca610f339f99a3"
  },
  {
    "path": "implementation_packs/MARKDOWN_COMPOSITOR_CORE.md",
    "sha256": "a8505384bb3f9085d36f21f1b6471fe24cbd17b054fa87f82f1d24b7007ce9d5"
  },
  {
    "path": "implementation_packs/PROJECT_START_READINESS_VALIDATOR.md",
    "sha256": "f0e132178236937e676fb829050745dadcb8ece60f778bb6b3e5c9a1a3e8fe79"
  },
  {
    "path": "markdown_system/CAPABILITY_CATALOG.md",
    "sha256": "1105e21fce8479d53b6aef2eb500d7c7fbc352450a04328904253ba01c51eb7f"
  },
  {
    "path": "markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md",
    "sha256": "370b0ec3a753f1aa6921cda1a6017c92b6ca1579999803e2bd74cc3ebb0c5336"
  },
  {
    "path": "markdown_system/PROJECT_START_READINESS_GATE.md",
    "sha256": "38fd1139eae774329d7072d76ab368ac26c28f42cd8bf03d7f5772ea1349ce59"
  },
  {
    "path": "markdown_system/test_agent_entry_lifecycle.ps1",
    "sha256": "158b9268658736db55a06ae4f22574827478b56cde1211aa1e76636b1feff21f"
  }
]
```

## Raw measured runs

Each evidence file retains every subprocess output and exit expectation. Local fixtures are retained separately.

```json
[
  {
    "name": "00-full",
    "location": "full",
    "warmup": true,
    "wall_ms": 10648.505,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-7d143160648644c19b3800fbf0664878\\evidence.json",
    "evidence_sha256": "d00e675b77a81815549af0a0cbe3e7516c0c9bfe8db662c2e212503edf94a7b7",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2184,
        "resume_ms": 112,
        "recovered_resume_ms": 109,
        "advanced_resume_ms": 117
      },
      "EXISTING": {
        "entry_subprocess_ms": 2154,
        "resume_ms": 122,
        "recovered_resume_ms": 92,
        "advanced_resume_ms": 109
      }
    }
  },
  {
    "name": "00-clean",
    "location": "clean",
    "warmup": true,
    "wall_ms": 10807.3368,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-8946d5e992b24d4e9f5c5641510e3a8e\\evidence.json",
    "evidence_sha256": "626cf18b4b06ff93d3516ebe1850c57184e7c6b8ab35024339703dc5bcc055b1",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2156,
        "resume_ms": 135,
        "recovered_resume_ms": 102,
        "advanced_resume_ms": 93
      },
      "EXISTING": {
        "entry_subprocess_ms": 2188,
        "resume_ms": 142,
        "recovered_resume_ms": 110,
        "advanced_resume_ms": 98
      }
    }
  },
  {
    "name": "01-full",
    "location": "full",
    "warmup": false,
    "wall_ms": 9923.4873,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-d931e43721274151bfa8dc68fd669f3b\\evidence.json",
    "evidence_sha256": "1a317b91e11f80e0b327a3918c9e3df4fc4fa8d0d3a733a65444df45b7e02cab",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2096,
        "resume_ms": 130,
        "recovered_resume_ms": 85,
        "advanced_resume_ms": 95
      },
      "EXISTING": {
        "entry_subprocess_ms": 2199,
        "resume_ms": 131,
        "recovered_resume_ms": 97,
        "advanced_resume_ms": 101
      }
    }
  },
  {
    "name": "01-clean",
    "location": "clean",
    "warmup": false,
    "wall_ms": 10784.9692,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-3d55436865f04ae5ad15353689ce677f\\evidence.json",
    "evidence_sha256": "0be5b4707737a76ab0ac6b7783d538af7f7e30be503c438426b82a2c813bfa66",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2204,
        "resume_ms": 146,
        "recovered_resume_ms": 85,
        "advanced_resume_ms": 88
      },
      "EXISTING": {
        "entry_subprocess_ms": 2292,
        "resume_ms": 151,
        "recovered_resume_ms": 136,
        "advanced_resume_ms": 146
      }
    }
  },
  {
    "name": "02-full",
    "location": "full",
    "warmup": false,
    "wall_ms": 9994.7816,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-27921fc14fdb4c5296eb1571146b743e\\evidence.json",
    "evidence_sha256": "161dc8421d59eeeccbe0e50598455698b3dbf9591db9175371a17694675390ed",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 1914,
        "resume_ms": 120,
        "recovered_resume_ms": 73,
        "advanced_resume_ms": 77
      },
      "EXISTING": {
        "entry_subprocess_ms": 2276,
        "resume_ms": 101,
        "recovered_resume_ms": 116,
        "advanced_resume_ms": 131
      }
    }
  },
  {
    "name": "02-clean",
    "location": "clean",
    "warmup": false,
    "wall_ms": 10028.2927,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-e58ecad7aa60411ebbeaf614d77d6975\\evidence.json",
    "evidence_sha256": "65f158e3b766ab68ea62df28731262d033ed83885a52bf49effbb6160fe5ca15",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 1918,
        "resume_ms": 124,
        "recovered_resume_ms": 100,
        "advanced_resume_ms": 129
      },
      "EXISTING": {
        "entry_subprocess_ms": 2177,
        "resume_ms": 115,
        "recovered_resume_ms": 100,
        "advanced_resume_ms": 131
      }
    }
  },
  {
    "name": "03-full",
    "location": "full",
    "warmup": false,
    "wall_ms": 10006.9305,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-1c09d765735e4486a719d23c6a6e9cda\\evidence.json",
    "evidence_sha256": "fd026884cacfbbcc29ec700ae3734d27b77e85f4cf34fca562ca5caa903c034f",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2019,
        "resume_ms": 112,
        "recovered_resume_ms": 97,
        "advanced_resume_ms": 99
      },
      "EXISTING": {
        "entry_subprocess_ms": 2209,
        "resume_ms": 153,
        "recovered_resume_ms": 99,
        "advanced_resume_ms": 126
      }
    }
  },
  {
    "name": "03-clean",
    "location": "clean",
    "warmup": false,
    "wall_ms": 9637.8014,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-33fe5871c25f46acb893528dbda32491\\evidence.json",
    "evidence_sha256": "56299b3c98896e9e9d12a183f04c2f133598a3c88fc8892d219ccf8ee6869cfe",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2000,
        "resume_ms": 115,
        "recovered_resume_ms": 81,
        "advanced_resume_ms": 90
      },
      "EXISTING": {
        "entry_subprocess_ms": 2137,
        "resume_ms": 122,
        "recovered_resume_ms": 100,
        "advanced_resume_ms": 108
      }
    }
  },
  {
    "name": "04-full",
    "location": "full",
    "warmup": false,
    "wall_ms": 9930.1433,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-ac96ed383e7947ccbedfb74071364d54\\evidence.json",
    "evidence_sha256": "6090348b653a7987224bcd40992737ab388360d64a939b10293fbc277a3482ec",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 1999,
        "resume_ms": 113,
        "recovered_resume_ms": 98,
        "advanced_resume_ms": 115
      },
      "EXISTING": {
        "entry_subprocess_ms": 2280,
        "resume_ms": 111,
        "recovered_resume_ms": 87,
        "advanced_resume_ms": 121
      }
    }
  },
  {
    "name": "04-clean",
    "location": "clean",
    "warmup": false,
    "wall_ms": 9918.5381,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-aab6398a0ed44cbf9c11e894dac5525d\\evidence.json",
    "evidence_sha256": "6cfb5dd47c68368e603d0d2aea7421fb83bee30489cf1f86e18dda574b8d0d3d",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2001,
        "resume_ms": 134,
        "recovered_resume_ms": 76,
        "advanced_resume_ms": 85
      },
      "EXISTING": {
        "entry_subprocess_ms": 2072,
        "resume_ms": 128,
        "recovered_resume_ms": 113,
        "advanced_resume_ms": 104
      }
    }
  },
  {
    "name": "05-full",
    "location": "full",
    "warmup": false,
    "wall_ms": 10473.0698,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-7136624e68ee4eae87a0bcf5c89ed541\\evidence.json",
    "evidence_sha256": "e20d38b2f7bb5163627a4b253bdb07de21d72253d1920df8f82d96703d5fe774",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2154,
        "resume_ms": 149,
        "recovered_resume_ms": 132,
        "advanced_resume_ms": 140
      },
      "EXISTING": {
        "entry_subprocess_ms": 2387,
        "resume_ms": 140,
        "recovered_resume_ms": 79,
        "advanced_resume_ms": 91
      }
    }
  },
  {
    "name": "05-clean",
    "location": "clean",
    "warmup": false,
    "wall_ms": 10056.7645,
    "evidence_path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-entry-lifecycle-00c382209b4247c3bee4746947ab5500\\evidence.json",
    "evidence_sha256": "27c0f46b3b459e9416eed5346892929040d332cb42b5a75052e2bfe5ee5dfde2",
    "checks": 55,
    "expected_negative": 8,
    "phases": {
      "NEW": {
        "entry_subprocess_ms": 2017,
        "resume_ms": 129,
        "recovered_resume_ms": 88,
        "advanced_resume_ms": 95
      },
      "EXISTING": {
        "entry_subprocess_ms": 2194,
        "resume_ms": 120,
        "recovered_resume_ms": 83,
        "advanced_resume_ms": 91
      }
    }
  }
]
```

## Reproduction programs

Programs are AUTHORED maintenance diagnostics. Resolve the local library/runtime paths for the actual host; retain fixtures and never apply their synthetic recovery writes to a real project event log.

### benchmark_entry.py

```python
from pathlib import Path
import json,hashlib,subprocess,time,os,re,platform,statistics,sys
R=Path(os.environ['ELITE_LIBRARY_ROOT']);S=Path(__file__).parent
PWSH=Path(os.environ['ELITE_PWSH_EXECUTABLE'])
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
files=['AGENT_SYSTEM_START.md','AGENTS.md','INSTALL_AGENT_BRIDGE.ps1','materialize_markdown_pack.ps1','START_ANY_PROJECT.md','VERIFY_EXECUTABLE_LIBRARY.ps1','implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md','implementation_packs/MARKDOWN_COMPOSITOR_CORE.md','implementation_packs/PROJECT_START_READINESS_VALIDATOR.md','markdown_system/CAPABILITY_CATALOG.md','markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md','markdown_system/PROJECT_START_READINESS_GATE.md','markdown_system/test_agent_entry_lifecycle.ps1']
copy=S/'clean-library';copy.mkdir(exist_ok=False)
manifest=[]
for n in files:
 p=copy/n;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes((R/n).read_bytes());assert p.read_bytes()==(R/n).read_bytes();manifest.append(dict(path=n,sha256=sha(p)))
(S/'source-manifest.json').write_text(json.dumps(manifest,indent=2)+'\n',encoding='utf-8',newline='\n')
env=dict(os.environ);env['PATH']=str(Path(sys.executable).parent)+os.pathsep+env['PATH'];env['PYTHONUTF8']='1'
protocol=dict(schema='elite-entry-benchmark/v1',source_manifest_sha256=sha(S/'source-manifest.json'),python=platform.python_version(),pwsh_sha256=sha(PWSH),warmups_per_location=1,samples_per_location=5,location_order=['full','clean']*6,same_tests=True,expected_checks=55,expected_calls=28,timeout_per_run_seconds=90,concurrency=1,clock='time.perf_counter_ns + existing Stopwatch per subprocess',cache_policy='OS cache not flushed; explicit warmups excluded',non_claims=['no product build','no production SLO or speedup claim','no token/billing/model measurement','no portable release admission for the thirteen-file test copy','NEW precedes EXISTING inside the unchanged harness; do not infer causal mode speed difference'])
(S/'protocol.json').write_text(json.dumps(protocol,indent=2)+'\n',encoding='utf-8',newline='\n')
runs=[]
for round_index in range(6):
 for location,root in [('full',R),('clean',copy)]:
  name=f'{round_index:02d}-{location}';start=time.perf_counter_ns()
  try: p=subprocess.run([str(PWSH),'-NoProfile','-File',str(root/'markdown_system/test_agent_entry_lifecycle.ps1'),'-LibraryRoot',str(root),'-KeepEvidence'],env=env,capture_output=True,timeout=90,check=False)
  except subprocess.TimeoutExpired as e:
   (S/(name+'-timeout.txt')).write_text(str(e),encoding='utf-8');raise
  elapsed=(time.perf_counter_ns()-start)/1e6
  (S/(name+'.stdout.log')).write_bytes(p.stdout);(S/(name+'.stderr.log')).write_bytes(p.stderr)
  out=p.stdout.decode('utf-8');assert p.returncode==0,(name,p.returncode,out[-3000:])
  found=re.findall(r'^LIFECYCLE_EVIDENCE=(.+)\r?$',out,re.M);assert len(found)==1,(name,out)
  evidence=Path(found[0].strip())/'evidence.json';data=json.loads(evidence.read_text(encoding='utf-8'))
  assert data['status']=='PASS' and data['checks']==55 and len(data['calls'])==28 and data['product_readiness'] is False
  assert all(c['exit_code']==c['expected_exit'] for c in data['calls'])
  expected_negative=sum(c['expected_exit']!=0 for c in data['calls']);assert expected_negative==8
  phases={}
  for mode in ['NEW','EXISTING']:
   calls={c['id']:c for c in data['calls']}
   phases[mode]={'entry_subprocess_ms':sum(calls[mode+'-'+step]['elapsed_ms'] for step in ['compose','bridge','advisory','readiness-blocked','checkpoint']),'resume_ms':calls[mode+'-resume']['elapsed_ms'],'recovered_resume_ms':calls[mode+'-recovered']['elapsed_ms'],'advanced_resume_ms':calls[mode+'-resume-advanced']['elapsed_ms']}
  runs.append(dict(name=name,location=location,warmup=round_index==0,wall_ms=elapsed,evidence_path=str(evidence),evidence_sha256=sha(evidence),checks=55,expected_negative=expected_negative,phases=phases,calls=data['calls']))
  (S/(name+'.evidence.json')).write_bytes(evidence.read_bytes())
  with (S/'progress.jsonl').open('a',encoding='utf-8',newline='\n') as f:f.write(json.dumps({k:v for k,v in runs[-1].items() if k!='calls'})+'\n')
  print(json.dumps(dict(name=name,warmup=round_index==0,wall_ms=round(elapsed,1),status='PASS')),flush=True)
  assert all(sha(R/m['path'])==m['sha256'] and sha(copy/m['path'])==m['sha256'] for m in manifest),'source changed during benchmark'
def stats(values):return dict(n=len(values),min_ms=min(values),median_ms=statistics.median(values),max_ms=max(values),mean_ms=statistics.mean(values),stdev_ms=statistics.stdev(values))
summary={}
for location in ['full','clean']:
 selected=[r for r in runs if r['location']==location and not r['warmup']]
 summary[location]={'complete_harness_wall':stats([r['wall_ms'] for r in selected])}
 for mode in ['NEW','EXISTING']:
  summary[location][mode]={k:stats([r['phases'][mode][k] for r in selected]) for k in selected[0]['phases'][mode]}
state=json.loads((R/'PROJECT_EXECUTION_STATE.json').read_text(encoding='utf-8'));index={e['id']:e for e in state['evidence']}
context=[dict(id=i,path=index[i]['path'],bytes=(R/index[i]['path']).stat().st_size) for i in state['context']['must_read_refs']]
md=list(R.rglob('*.md'));surface=dict(all_markdown_files=len(md),all_markdown_bytes=sum(p.stat().st_size for p in md),must_read_refs=len(context),must_read_owner_bytes=sum(x['bytes'] for x in context),owners=context,claim='file-byte inventory only; indexed excerpts actually read, tokens, billing and comprehension are not measured')
(S/'results.json').write_text(json.dumps(dict(protocol=protocol,summary=summary,runs=runs,context_surface=surface),ensure_ascii=False,indent=2)+'\n',encoding='utf-8',newline='\n')
print('ENTRY_BENCHMARK_COMPLETE',flush=True)
```

### checkpoint_fault_child.py

```python
from pathlib import Path
import sys,os
kit,root,mode=map(str,sys.argv[1:]);sys.path.insert(0,kit)
import checkpoint_execution_state as checkpoint
root=Path(root);real_fdopen=os.fdopen
class InjectedWriter:
    def __init__(self,handle):self.handle=handle
    def __enter__(self):return self
    def __exit__(self,*args):self.handle.close()
    def write(self,payload):
        if mode=='before_write':raise OSError('INJECTED_BEFORE_WRITE')
        part=payload if mode=='crash_after_fsync' else payload[:len(payload)//2]
        self.handle.write(part);self.handle.flush();os.fsync(self.handle.fileno())
        if mode=='partial_error':raise OSError('INJECTED_PARTIAL_WRITE')
        os._exit(79)
    def flush(self):return self.handle.flush()
    def fileno(self):return self.handle.fileno()
checkpoint.os.fdopen=lambda *a,**kw:InjectedWriter(real_fdopen(*a,**kw))
try:
    checkpoint.append_checkpoint(root/'PROJECT_EXECUTION_STATE.json',root,root/'PROJECT_EXECUTION_EVENTS.jsonl','EVT-0002','STATE_ADVANCED','synthetic-test','Interrupted synthetic checkpoint')
except OSError as e:
    if str(e) not in ('INJECTED_BEFORE_WRITE','INJECTED_PARTIAL_WRITE'):raise
    print(str(e));sys.exit(78)
raise AssertionError('fault did not fire')
```

### checkpoint_fault_probe.py

```python
from pathlib import Path
import json,sys,subprocess,hashlib,time
S=Path(__file__).parent;kit=S/'fault-kit/engineering_execution_kit';sys.path.insert(0,str(kit))
import test_execution_state as fixtures
from checkpoint_execution_state import append_checkpoint
from validate_execution_state import validate_event_log,StateError
sha=lambda b:hashlib.sha256(b).hexdigest()
receipts=[]
for mode in ['before_write','partial_error','partial_crash','crash_after_fsync']:
    root=S/'fault-cases'/mode;root.mkdir(parents=True,exist_ok=False)
    fixture=fixtures.ExecutionStateTests();fixture.root=root;fixture.events=root/'PROJECT_EXECUTION_EVENTS.jsonl';fixture.files={}
    fixture._write('failure','PROJECT_FAILURE_LESSONS.md','# Synthetic failure ledger\n')
    state=fixture._base('EXISTING');p=root/'PROJECT_EXECUTION_STATE.json';events=fixture.events
    p.write_text(json.dumps(state,ensure_ascii=False,indent=2)+'\n',encoding='utf-8',newline='\n')
    append_checkpoint(p,root,events,'EVT-0001','BASELINE_CAPTURED','synthetic-test','Preserve original event')
    before=events.read_bytes();state['revision']=2;state['next_action']='Synthetic interrupted append, no product approval'
    p.write_text(json.dumps(state,ensure_ascii=False,indent=2)+'\n',encoding='utf-8',newline='\n')
    evidence_before={f:sha((root/f).read_bytes()) for f in ['PROJECT_FAILURE_LESSONS.md','evidence/inventory.txt','evidence/delta.txt']}
    start=time.perf_counter_ns()
    child=subprocess.run([sys.executable,'-X','utf8',str(S/'checkpoint_fault_child.py'),str(kit),str(root),mode],capture_output=True,timeout=20,check=False)
    expected=78 if mode in ['before_write','partial_error'] else 79;assert child.returncode==expected,(mode,child.returncode,child.stderr)
    after=events.read_bytes();assert after.startswith(before)
    saved=S/'fault-preserved'/mode;saved.mkdir(parents=True,exist_ok=False)
    (saved/'before.events.jsonl').write_bytes(before);(saved/'interrupted.events.jsonl').write_bytes(after);(saved/'state.json').write_bytes(p.read_bytes());(saved/'child.stdout').write_bytes(child.stdout);(saved/'child.stderr').write_bytes(child.stderr)
    args=[sys.executable,'-X','utf8',str(kit/'validate_execution_state.py'),str(p),'--project-root',str(root),'--events',str(events),'--level','resume']
    validation=subprocess.run(args,capture_output=True,check=False,timeout=20)
    assert validation.returncode==(0 if mode=='crash_after_fsync' else 2),(mode,validation.returncode,validation.stdout,validation.stderr)
    (saved/'validation.stdout').write_bytes(validation.stdout);(saved/'validation.stderr').write_bytes(validation.stderr)
    if mode=='before_write':
        assert after==before
        append_checkpoint(p,root,events,'EVT-0002','STATE_ADVANCED','synthetic-test','Recovered before-write failure')
    elif mode in ['partial_error','partial_crash']:
        # Preserve every corrupted byte before restoring only this owned synthetic log.
        try:append_checkpoint(p,root,events,'EVT-0002','STATE_ADVANCED','synthetic-test','Must reject corrupt tail')
        except StateError:pass
        else:raise AssertionError('corrupt tail accepted')
        assert events.read_bytes()==after
        assert len(validate_event_log(saved/'before.events.jsonl'))==1
        events.write_bytes(before)
        append_checkpoint(p,root,events,'EVT-0002','STATE_ADVANCED','synthetic-test','Recovered exact known prefix; interrupted bytes retained separately')
    else:
        assert len(validate_event_log(events))==2
        try:append_checkpoint(p,root,events,'EVT-0002','STATE_ADVANCED','synthetic-test','Must reject duplicate after response loss')
        except StateError:pass
        else:raise AssertionError('duplicate checkpoint accepted')
        assert events.read_bytes()==after
    final=subprocess.run(args,capture_output=True,check=False,timeout=20);assert final.returncode==0,(mode,final.stdout,final.stderr)
    assert len(validate_event_log(events))==2 and events.read_bytes().startswith(before)
    assert all(sha((root/f).read_bytes())==digest for f,digest in evidence_before.items())
    elapsed=(time.perf_counter_ns()-start)/1e6
    receipts.append(dict(mode=mode,child_exit=child.returncode,resume_after_interrupt=validation.returncode,resume_after_recovery=final.returncode,original_event_sha256=sha(before),interrupted_log_sha256=sha(after),final_log_sha256=sha(events.read_bytes()),original_event_byte_exact=True,evidence_files_unchanged=True,final_event_count=2,elapsed_ms=elapsed,status='PASS'))
(S/'fault-results.json').write_text(json.dumps(dict(schema='elite-checkpoint-fault-probe/v1',cases=receipts,limits=['Synthetic injected write boundary and real child process exit; not power loss, disk/controller failure, concurrent writers or materializer/bridge crash atomicity','No automatic repair of real project events; known-prefix recovery occurred only in owned synthetic fixtures with all damaged bytes retained']),indent=2)+'\n',encoding='utf-8',newline='\n')
print(json.dumps(receipts,indent=2))
```

## Evidence hashes

- results.json: 7af44169aeb00db70a6a9b9eecdf00a1b649df2dc452af8ef7906eed577c1dd9
- fault-results.json: c50ed93d87c4712f0be93c8040e91ad5839b57c646a9829fa1db941514371fb4
- environment.json: 9cff088ec76db5631d2b19f15548504a72525aaca0cadab538f66bd2e0b2310f
- protocol.json: 8e6fe7d8b7f1bd60db411b862616302ccc3a3ffd600c1068dc5bd7183bd743ca
- source-manifest.json: c408021fd795fec454eca8c3efc1ca52303db353ee504b18cbe89f3374fce377

## Portable-example correction / FAIL536

The initial library gate rejected a machine-specific home path in the embedded benchmark program. The original executed script and raw measurements remain retained; the example now uses ELITE_LIBRARY_ROOT and ELITE_PWSH_EXECUTABLE and the invoking Python directory. Its syntax is checked; the twelve-run performance sample is from the original script, not a claim that a different host was benchmarked. Set those variables to observed paths and run from a fresh diagnostic directory. Original script SHA256 09b700888a447df8957e4d405dcaf669f3f4bdb43d9b270fdc5a172f889a1ab5; portable example SHA256 42d5bee42b4380025343ecf66bd96f2af0e956482415e03afbba47cdf809361b. The distribution rule was not relaxed.

## Continued materializer interruption verification — TEST-04

Eight additional cases inject exception or actual child termination after complete or half-file writes at files1and14of15. The temporary instrumented copy changes only the selected fault boundary; canonical code remains unchanged. Every case preserves the existing synthetic user file/BOM and partial output, rejects reentry into the occupied destination, and reconstructs15/15exact files using the unchanged canonical script in a new destination. Half-file cases additionally check exact source prefixes. No partial output is deleted or promoted.

Together with the four checkpoint cases, this executes TEST-04 for direct materializer/checkpoint interruption and recovery:12boundaries with process restart and exact preservation/reconstruction. Earlier statements that its materializer subset was planned describe the initial checkpoint-only cut; this continuation supersedes them for the tested scope. Bridge crash atomicity, concurrent writers, power loss, final-release restore and target recovery remain unproven. T2801/T2808/T2810 are not completed by one test.

```json
{
  "complete_file_boundaries": {
    "canonical_materializer_sha256": "7ae045ff9b03ea22242d376d16b2a114f00e2f095217789eb92ca1f19e126356",
    "instrumented_materializer_sha256": "f23ffff5327a226e7d9fed20097a855af2070736edaf2cc820689f81fb8400fe",
    "pack_sha256": "ab8a4c3fb4dbf39900dad8170a303ebfcfa56f622256e09b6cca610f339f99a3",
    "injection": "Exact write call retained, marker plus injected exception or pause added immediately after files1/14; real child terminated only at marker. No canonical script edits.",
    "cases": [
      {
        "mode": "error",
        "after_file": 1,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 1,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2178.5669,
        "status": "PASS"
      },
      {
        "mode": "kill",
        "after_file": 1,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 1,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2127.8394,
        "status": "PASS"
      },
      {
        "mode": "error",
        "after_file": 14,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 14,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2085.6955,
        "status": "PASS"
      },
      {
        "mode": "kill",
        "after_file": 14,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 14,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 1914.9187,
        "status": "PASS"
      }
    ],
    "limits": [
      "Synthetic exception after a successful file write, not physical disk failure",
      "Partial destination intentionally retained; recovery uses a new destination without deletion or overwrite",
      "No proof of whole-tree crash atomicity, bridge crash safety, concurrent writers or power loss"
    ]
  },
  "partial_file_boundaries": {
    "canonical_materializer_sha256": "7ae045ff9b03ea22242d376d16b2a114f00e2f095217789eb92ca1f19e126356",
    "instrumented_materializer_sha256": "396961c0b3145214188976c3c54270463e6b655a2f13ebc3fef3641513571b87",
    "pack_sha256": "ab8a4c3fb4dbf39900dad8170a303ebfcfa56f622256e09b6cca610f339f99a3",
    "injection": "Synthetic half-file write at files1/14 then injected exception or marked pause; real child terminated at marker. Other file writes retain canonical call. No canonical edits.",
    "cases": [
      {
        "mode": "error",
        "after_file": 1,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 1,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2510.8085,
        "status": "PASS"
      },
      {
        "mode": "kill",
        "after_file": 1,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 1,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2473.2325,
        "status": "PASS"
      },
      {
        "mode": "error",
        "after_file": 14,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 14,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2418.6016,
        "status": "PASS"
      },
      {
        "mode": "kill",
        "after_file": 14,
        "interrupted_exit": 1,
        "reentry_exit": 1,
        "rebuild_exit": 0,
        "partial_files": 14,
        "recovered_files": 15,
        "original_preserved": true,
        "partial_preserved": true,
        "recovered_byte_parity": true,
        "elapsed_ms": 2319.3817,
        "status": "PASS"
      }
    ],
    "limits": [
      "Synthetic partial-byte write, not physical disk failure",
      "Partial destination intentionally retained; recovery uses a new destination without deletion or overwrite",
      "No proof of whole-tree crash atomicity, bridge crash safety, concurrent writers or power loss"
    ]
  }
}
```

### materializer_fault_probe.py

```python
from pathlib import Path
import os,sys,subprocess,time,json,hashlib
S=Path(__file__).parent;R=Path(os.environ['ELITE_LIBRARY_ROOT']);pwsh=Path(os.environ['ELITE_PWSH_EXECUTABLE'])
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
source=(R/'materialize_markdown_pack.ps1').read_text(encoding='utf-8');needle='  [IO.File]::WriteAllBytes($target, $block.Bytes)';assert source.count(needle)==1
instrumented=source.replace("$ErrorActionPreference = 'Stop'","$ErrorActionPreference = 'Stop'\n$script:probeWriteCount = 0",1).replace(needle,needle+'''\n  $script:probeWriteCount++
  if ($script:probeWriteCount -eq [int]$env:ELITE_FAULT_AFTER) {
    [IO.File]::WriteAllText($env:ELITE_FAULT_MARKER, [string]$script:probeWriteCount)
    if ($env:ELITE_FAULT_MODE -eq 'error') { throw 'INJECTED_AFTER_FILE_WRITE' }
    [Threading.Thread]::Sleep(30000)
  }''')
script=S/'materializer-instrumented.ps1';script.write_text(instrumented,encoding='utf-8',newline='\n')
expected={p.relative_to(S/'fault-kit').as_posix():sha(p) for p in (S/'fault-kit').rglob('*') if p.is_file() and '__pycache__' not in p.parts}
assert len(expected)==15,len(expected)
results=[]
for mode,after in [('error',1),('kill',1),('error',14),('kill',14)]:
    name=f'{mode}-{after}';case=S/'materializer-cases'/name;case.mkdir(parents=True,exist_ok=False)
    project=case/'project';project.mkdir();owned=project/'existing-user-code.txt';owned.write_bytes(b'\xef\xbb\xbfOriginal synthetic user bytes\r\n');owned_sha=sha(owned)
    partial=project/'partial-tooling';marker=case/'boundary.txt';env=dict(os.environ,ELITE_FAULT_MODE=mode,ELITE_FAULT_AFTER=str(after),ELITE_FAULT_MARKER=str(marker))
    args=[str(pwsh),'-NoProfile','-File',str(script),'-PackFile',str(R/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',str(partial)]
    start=time.perf_counter_ns()
    with (case/'injected.log').open('wb') as out:
        child=subprocess.Popen(args,env=env,stdout=out,stderr=subprocess.STDOUT,creationflags=getattr(subprocess,'CREATE_NO_WINDOW',0))
        deadline=time.monotonic()+10
        while not marker.exists() and child.poll() is None and time.monotonic()<deadline:time.sleep(.025)
        if not marker.exists():
            if child.poll() is None:child.kill()
            child.wait(timeout=10);raise AssertionError(f'{name}: boundary not reached')
        if mode=='kill':child.kill()
        code=child.wait(timeout=10)
    assert code!=0 and int(marker.read_text(encoding='utf-8'))==after
    if mode=='error':assert 'INJECTED_AFTER_FILE_WRITE' in (case/'injected.log').read_text(encoding='utf-8')
    current={p.relative_to(partial).as_posix():sha(p) for p in partial.rglob('*') if p.is_file()}
    assert len(current)==after and all(expected[n]==h for n,h in current.items())
    assert sha(owned)==owned_sha
    normal=[str(pwsh),'-NoProfile','-File',str(R/'materialize_markdown_pack.ps1'),'-PackFile',str(R/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',str(partial)]
    retry=subprocess.run(normal,capture_output=True,timeout=15,check=False);assert retry.returncode!=0
    (case/'occupied-destination.log').write_bytes(retry.stdout+retry.stderr)
    assert current=={p.relative_to(partial).as_posix():sha(p) for p in partial.rglob('*') if p.is_file()}
    recovered=project/'recovered-tooling';normal[-1]=str(recovered)
    rebuild=subprocess.run(normal,capture_output=True,timeout=15,check=False);assert rebuild.returncode==0,(name,rebuild.stderr)
    (case/'rebuild.log').write_bytes(rebuild.stdout+rebuild.stderr)
    actual={p.relative_to(recovered).as_posix():sha(p) for p in recovered.rglob('*') if p.is_file()}
    assert actual==expected and sha(owned)==owned_sha
    results.append(dict(mode=mode,after_file=after,interrupted_exit=code,reentry_exit=retry.returncode,rebuild_exit=rebuild.returncode,partial_files=after,recovered_files=len(actual),original_preserved=True,partial_preserved=True,recovered_byte_parity=True,elapsed_ms=(time.perf_counter_ns()-start)/1e6,status='PASS'))
(S/'materializer-fault-results.json').write_text(json.dumps(dict(canonical_materializer_sha256=sha(R/'materialize_markdown_pack.ps1'),instrumented_materializer_sha256=sha(script),pack_sha256=sha(R/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),injection='Exact write call retained, marker plus injected exception or pause added immediately after files1/14; real child terminated only at marker. No canonical script edits.',cases=results,limits=['Synthetic exception after a successful file write, not physical disk failure','Partial destination intentionally retained; recovery uses a new destination without deletion or overwrite','No proof of whole-tree crash atomicity, bridge crash safety, concurrent writers or power loss']),indent=2)+'\n',encoding='utf-8',newline='\n')
print(json.dumps(results,indent=2))
```

### materializer_partial_fault_probe.py

```python
from pathlib import Path
import os,sys,subprocess,time,json,hashlib
S=Path(__file__).parent;R=Path(os.environ['ELITE_LIBRARY_ROOT']);pwsh=Path(os.environ['ELITE_PWSH_EXECUTABLE'])
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
source=(R/'materialize_markdown_pack.ps1').read_text(encoding='utf-8');needle='  [IO.File]::WriteAllBytes($target, $block.Bytes)';assert source.count(needle)==1
instrumented=source.replace("$ErrorActionPreference = 'Stop'","$ErrorActionPreference = 'Stop'\n$script:probeWriteCount = 0",1).replace(needle,'''  $script:probeWriteCount++
  if ($script:probeWriteCount -eq [int]$env:ELITE_FAULT_AFTER) {
    $half = [Math]::Floor($block.Bytes.Length / 2)
    $partialBytes = [byte[]]::new($half)
    [Array]::Copy($block.Bytes, $partialBytes, $half)
    [IO.File]::WriteAllBytes($target, $partialBytes)
    [IO.File]::WriteAllText($env:ELITE_FAULT_MARKER, [string]$script:probeWriteCount)
    if ($env:ELITE_FAULT_MODE -eq 'error') { throw 'INJECTED_PARTIAL_FILE_WRITE' }
    [Threading.Thread]::Sleep(30000)
  } else {
    [IO.File]::WriteAllBytes($target, $block.Bytes)
  }''')
script=S/'materializer-partial-instrumented.ps1';script.write_text(instrumented,encoding='utf-8',newline='\n')
expected={p.relative_to(S/'fault-kit').as_posix():sha(p) for p in (S/'fault-kit').rglob('*') if p.is_file() and '__pycache__' not in p.parts}
assert len(expected)==15,len(expected)
results=[]
for mode,after in [('error',1),('kill',1),('error',14),('kill',14)]:
    name=f'partial-{mode}-{after}';case=S/'materializer-cases'/name;case.mkdir(parents=True,exist_ok=False)
    project=case/'project';project.mkdir();owned=project/'existing-user-code.txt';owned.write_bytes(b'\xef\xbb\xbfOriginal synthetic user bytes\r\n');owned_sha=sha(owned)
    partial=project/'partial-tooling';marker=case/'boundary.txt';env=dict(os.environ,ELITE_FAULT_MODE=mode,ELITE_FAULT_AFTER=str(after),ELITE_FAULT_MARKER=str(marker))
    args=[str(pwsh),'-NoProfile','-File',str(script),'-PackFile',str(R/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',str(partial)]
    start=time.perf_counter_ns()
    with (case/'injected.log').open('wb') as out:
        child=subprocess.Popen(args,env=env,stdout=out,stderr=subprocess.STDOUT,creationflags=getattr(subprocess,'CREATE_NO_WINDOW',0))
        deadline=time.monotonic()+10
        while not marker.exists() and child.poll() is None and time.monotonic()<deadline:time.sleep(.025)
        if not marker.exists():
            if child.poll() is None:child.kill()
            child.wait(timeout=10);raise AssertionError(f'{name}: boundary not reached')
        if mode=='kill':child.kill()
        code=child.wait(timeout=10)
    assert code!=0 and int(marker.read_text(encoding='utf-8'))==after
    if mode=='error':assert 'INJECTED_PARTIAL_FILE_WRITE' in (case/'injected.log').read_text(encoding='utf-8')
    current={p.relative_to(partial).as_posix():sha(p) for p in partial.rglob('*') if p.is_file()}
    assert len(current)==after
    damaged=[n for n,h in current.items() if expected[n]!=h];assert len(damaged)==1
    original=(S/'fault-kit'/damaged[0]).read_bytes();assert (partial/damaged[0]).read_bytes()==original[:len(original)//2]
    assert sha(owned)==owned_sha
    normal=[str(pwsh),'-NoProfile','-File',str(R/'materialize_markdown_pack.ps1'),'-PackFile',str(R/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),'-Destination',str(partial)]
    retry=subprocess.run(normal,capture_output=True,timeout=15,check=False);assert retry.returncode!=0
    (case/'occupied-destination.log').write_bytes(retry.stdout+retry.stderr)
    assert current=={p.relative_to(partial).as_posix():sha(p) for p in partial.rglob('*') if p.is_file()}
    recovered=project/'recovered-tooling';normal[-1]=str(recovered)
    rebuild=subprocess.run(normal,capture_output=True,timeout=15,check=False);assert rebuild.returncode==0,(name,rebuild.stderr)
    (case/'rebuild.log').write_bytes(rebuild.stdout+rebuild.stderr)
    actual={p.relative_to(recovered).as_posix():sha(p) for p in recovered.rglob('*') if p.is_file()}
    assert actual==expected and sha(owned)==owned_sha
    results.append(dict(mode=mode,after_file=after,interrupted_exit=code,reentry_exit=retry.returncode,rebuild_exit=rebuild.returncode,partial_files=after,recovered_files=len(actual),original_preserved=True,partial_preserved=True,recovered_byte_parity=True,elapsed_ms=(time.perf_counter_ns()-start)/1e6,status='PASS'))
(S/'materializer-partial-fault-results.json').write_text(json.dumps(dict(canonical_materializer_sha256=sha(R/'materialize_markdown_pack.ps1'),instrumented_materializer_sha256=sha(script),pack_sha256=sha(R/'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md'),injection='Synthetic half-file write at files1/14 then injected exception or marked pause; real child terminated at marker. Other file writes retain canonical call. No canonical edits.',cases=results,limits=['Synthetic partial-byte write, not physical disk failure','Partial destination intentionally retained; recovery uses a new destination without deletion or overwrite','No proof of whole-tree crash atomicity, bridge crash safety, concurrent writers or power loss']),indent=2)+'\n',encoding='utf-8',newline='\n')
print(json.dumps(results,indent=2))
```

## Integrated gate and FAIL536

VERIFY_LIBRARY at checkpoint83 passed161packs/1453files/761Markdown/52profiles after the portable-path correction; log SHA256 89cf5b562529d0f7d758b0b4755e27830ec1bcc252ebc7c01c50aae066dff2ac. FAIL536 is closed for that correction and its rejected attempt remains preserved. The materializer continuation adds evidence only, not packs or runtimes. Final validation follows checkpoint84.
