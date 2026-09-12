# UTF-8 checkpoint recovery V324

2026-09-08. Corrective library maintenance / FAIL-20260908-533. No product expansion or readiness promotion.

## Observed failure and containment

Checkpoint78 was rejected because open.blockers[4] exceeded300characters (445). The last durable event remained EVT-0077. Temporary Python scripts read UTF-8 JSON with an implicit Windows codepage, then wrote UTF-8, amplifying corruption. Canonical load_json and event parsing already use explicit UTF-8; no validator or length limit was weakened.

Preserved the four root JSON files and original event bytes before repair. Only208contract string values and1cursor blocker were restored. Every restoration strictly inverts cp1252/UTF-8 and reproduces the exact damaged string for its recorded number of layers. JSON keys, numeric values, IDs, decisions, evidence hashes and all77events remained unchanged during this recovery step. Subsequent evidence/hash/cursor changes are a separate checkpoint operation. No guessed character replacements or historical event rewrites.

## Regression and boundaries

Fresh canonical EXECUTION-VALIDATOR1.3.0 materialization:15files. Thirteen existing state tests plus two Unicode regression tests pass (15unique tests) under explicit UTF-8 and with UTF-8 mode deliberately disabled. The new cases round-trip Spanish, currency, punctuation, Japanese and emoji through actual load/checkpoint/event verification, and reproduce/reverse six layers of the observed Windows decoding failure. This is text-integrity evidence, not product acceptance.

Eighteen temporary V322/V323 scripts now read text with explicit UTF-8; original script bytes are preserved and corrected syntax compiles. Do not rerun their mutations. Future commands use Python -X utf8 and explicit read/write encoding. All four root JSON documents were inspected; only the two restored files had reversible corruption. Canonical implementation_packs and markdown_system Markdown had no observed amplification fingerprint; this bounded check is not a general natural-language encoding detector.

FAIL532 native provenance/licensing and readinessD-H remain open. T2809 still lacks real service identity/reporting/supervisor integration. No ZIP, runtime promotion, provider effects or100%completion claim.

## Reproducible recovery receipt

```json
{
  "files": [
    {
      "file": "PROJECT_ENGINEERING_CONTRACT.json",
      "before_sha256": "0368c1d6f0a34103c139dbd300c6f989c09284ed80cf39511058db594723b2ee",
      "after_sha256": "6ce3e852570fdf4a7c0dfc52f740309154f8796302b9d2ddbcf44804cf982cef",
      "changed_string_fields": 208
    },
    {
      "file": "PROJECT_EXECUTION_STATE.json",
      "before_sha256": "94a2fb2ea650780d95a3e7732361c8e2c7b3797eb91d15d80ebc56958a874906",
      "after_sha256": "9835adfd3dfe11daf95fe954a70478eaf3d21dd247e7270370e2b1192ed96d6d",
      "changed_string_fields": 1
    }
  ],
  "events_unchanged_sha256": "28244e3e06da78fe9134874513804028dd98ae799561052f22438450078c053a",
  "prior_event_count": 77,
  "script_fixes": [
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\checkpoint73.py",
      "before_sha256": "34e143a0ab0b66ffc74599202bfd26d2e5b2dc590f4b5ecfff677e8534e88871",
      "after_sha256": "2e3bc20060951ddbeb3c2b76058faa3fa05ca2cb5e345f57a20a4f462ea69ca1"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\checkpoint74.py",
      "before_sha256": "a4b3e62efb711100cf53b6be43715051c249ef1638cf08ad1a3ceb85c0040a31",
      "after_sha256": "4280e4f6c31ed8658559c2a2e21f26a7ffc37d6c9354e0732ef8d69e262201cb"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\checkpoint75.py",
      "before_sha256": "afc8a6c88078a16b41f5be475441dd140a900214a874cad2cc47be57b7eb28f5",
      "after_sha256": "f9cf055b850d1979feee6c33ef30302ae3d6f70548f72ee2da747061456515ea"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\encode-license-sidecars.py",
      "before_sha256": "ebfce2f83c70347539538a8dab9e63b8e50490b8d78a4d87af7fd629046b23f3",
      "after_sha256": "0bb1b43330a03c4e9db0de938fd961b129caaef923dddc65c686bd156c7df75e"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\finalize76.py",
      "before_sha256": "11da9e239810b592c4155743cfe5cfb6fc1958f6d2ad9a5dc394872673ba8352",
      "after_sha256": "74b8ab1787b43a167971e59708389aa432f8af71e537202c3473fab24ffe0f4f"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\fix-profile-counts.py",
      "before_sha256": "e2e250f850682ea96d4059b7b366369ee32fcde830b688cb55c0480388e8191b",
      "after_sha256": "65ae94239ad9c41d57d82fa118f6176a3d9c1e7b4870624d8805948478033000"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\integrate-candidate.py",
      "before_sha256": "d46de3ea62adf3eca389fcd83ac14547ab916d27665501051c54768fdf748018",
      "after_sha256": "f69a5ebdf9d7df97c1d3f5cd24568b56fb12e8acc8894e7a0e07f9f90aad88f2"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\prepare-intake.py",
      "before_sha256": "d2352e86a801dacecc06b7a6f04403b034dbf4ce2220aca6c301472ba0e0de0d",
      "after_sha256": "b84c8200f6c88caa45125cc3c23acbc43f8b46ee90c971bfd6324629c5723601"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\promote-pack.py",
      "before_sha256": "8644b9d3e3cbc63fe492429e43a58ecacdace1dcabae1480c03559e5e826d897",
      "after_sha256": "0837603b6abcd7cabd7002179cabc67f7598f5884d50aa15a418ce1336a5bc0b"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\qualify-and-answer.py",
      "before_sha256": "d05ced95d115901f10247825744836cc645c8852d2cf960cfb373d705345eedd",
      "after_sha256": "d5410714153a3d7712993b0eaa1ad809b987dd23ef9cbeee8a6988273ea5b6e2"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\record-qualification.py",
      "before_sha256": "d9a25545d1459824168614be83722fe0ab63f67bb07653d286ee4361c61273a3",
      "after_sha256": "8c42a951b74886412759a23a2c5ad270ac163528648faaf5444464d9b04c73e2"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v322-a5c006a784f04fc08764f4a3853f17a5\\verify-acquisition.py",
      "before_sha256": "346a34c85b6bbba9b0a25863c3e87205ae8456a8695e007c7dabe1dab6548355",
      "after_sha256": "4c398edf5c84682328c49fdb7f07fd45dfacad6737abe07b6d07aeca1c2a72b1"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v323-b72a92c4fe7c4907947b1e919e2469d0\\attestation-inputs.py",
      "before_sha256": "ac026b91137589e101edc574bc570172070bb91acda9aaa9f2e4f671143147e5",
      "after_sha256": "58d6b8bb7c5ec626118b3edf33012ac68acda138c40db541b2148ce16c7ba062"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v323-b72a92c4fe7c4907947b1e919e2469d0\\checkpoint78.py",
      "before_sha256": "673bb5059dc6f138ae68e12210335fd4fc400193d27c4f143c07dbf70e01be56",
      "after_sha256": "4ee2d2f05f9478c29799a430334c3252c7ec3be8e8c5759c097a8f77f936b978"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v323-b72a92c4fe7c4907947b1e919e2469d0\\record77.py",
      "before_sha256": "374c7774d2cb179c6705fe70cf233956a76c52a42b7a5c7d7e2af6e4f25be001",
      "after_sha256": "335ae2ef6008c103fff0a53175a71270f6cb4bb48eccbc2bc61e8e446c976d3e"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v323-b72a92c4fe7c4907947b1e919e2469d0\\research-manifests.py",
      "before_sha256": "f513edb21b88243cacd18dad681210016a344957cdd28ee8fadb3aaebab8737e",
      "after_sha256": "0d1bbc9dbb04e4e9d27871b3f4f2f83e7aed4aaa45f978fcdbcfc676b29fb53b"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v323-b72a92c4fe7c4907947b1e919e2469d0\\research-native-evidence.py",
      "before_sha256": "039bdd15e684a9111a68759f7740c1167e8a1cfaa3d6425225cf4bbc66af62ca",
      "after_sha256": "3dc6cd0cb34a46ee8f2049d894a8b3cfed3c2f3c5f3c893b595e78a165fc37c3"
    },
    {
      "path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v323-b72a92c4fe7c4907947b1e919e2469d0\\rustsec-stdlib-records.py",
      "before_sha256": "b946442f193fac6f7b98d6549bacbd3ca3f69b7520815be7ac431b828957176f",
      "after_sha256": "d93048d80597269eb60f524d9ee6086e751dc7c4e1d52a60a5ae2e4588da7f3f"
    }
  ]
}
```

## Per-field restoration audit

Reapplying UTF-8 encode/cp1252 decode `layers` times to `after` reproduces `before_utf8_sha256`. Original audit SHA256 fc848b8da9df03a684e301063e84b92e3e77e980bd71f767152ec308ebe04023.

```json
[
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "project",
      "objective"
    ],
    "after": "Completar el roadmap existente con evidencia conectada y reutilización NEW/EXISTING; este contrato no certifica una franquicia.",
    "layers": 4,
    "before_utf8_sha256": "462e1a565a7f9d736d76bfbc3e00652909cdb7545187c148960494a639624cab"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "project",
      "scope",
      1
    ],
    "after": "La aceptación de todo el perfil requiere sus journeys; no queda cubierta sólo por TEST-01",
    "layers": 4,
    "before_utf8_sha256": "01922fa5cece0492eef8124127c0ec535ef0ab63a7405f5b2004615a63e92bd1"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "project",
      "out_of_scope",
      0
    ],
    "after": "Despliegue real, cuentas, gastos, campañas, mensajes o datos privados sin autoridad",
    "layers": 4,
    "before_utf8_sha256": "be1386ae3538d4d225acb16eedf0cd13450cbc9848c9a748a632e42a9f6805e2"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "implementation_assurance",
      "blockers",
      1
    ],
    "after": "21 cores conciliados V293:12 solapamientos parciales,8 sin equivalente demostrado y observabilidad CANDIDATE; integración/admisión target pendientes",
    "layers": 4,
    "before_utf8_sha256": "147648c691f07526057e4699570c1093c349a1ac45f926d4d3e26347afe251ab"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "implementation_assurance",
      "blockers",
      2
    ],
    "after": "Observabilidad0.3.0: contador V291, logger V292 e histograma V301 reparados localmente; política/integración/operación target pendientes, CANDIDATE",
    "layers": 4,
    "before_utf8_sha256": "c2775e1a325079f8f4fe34a4c429856f5fc3b7df67ef3475bed988ab91df0a81"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "implementation_assurance",
      "blockers",
      3
    ],
    "after": "Journeys del perfil T2802-T2809 y aceptación release-bound pendientes",
    "layers": 4,
    "before_utf8_sha256": "ccd40f018c1a6cf6a7387677436c1999c563331a579e672a4aae55a9ed0fc0b8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "requirements",
      0,
      "statement"
    ],
    "after": "Entrada y reanudación NEW/EXISTING sin dañar cambios existentes; LIB-R01/R08, T2801.",
    "layers": 4,
    "before_utf8_sha256": "2bb588d309d4a953f4ce0186db60a4756b559c92200bae27217fdbe4b0f23bba"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "requirements",
      1,
      "statement"
    ],
    "after": "Procedencia/admisión verificables, no equivalencia por fama; LIB-R02/R07, T2801/T2803.",
    "layers": 4,
    "before_utf8_sha256": "4e811b1a247b0be251b7406f0352f8fcebb19f938e5c2f6a9ce4eb5460f07d6c"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "requirements",
      2,
      "statement"
    ],
    "after": "Fallos diagnosticables y recuperación segura; LIB-R06, T2803/T2809.",
    "layers": 4,
    "before_utf8_sha256": "f8f1da16b05e1b142cc0664338b3f8b63ebb83ae597b8609869a7e5c0ec84454"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "requirements",
      3,
      "statement"
    ],
    "after": "Distribución reproducible con instrucciones utilizables, sin heredar aprobaciones; LIB-R09, T2810.",
    "layers": 4,
    "before_utf8_sha256": "07581e396cba0c33abb31a32722c80e092f3bb4fcf863d4ff49c074e7787d4e9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "requirements",
      4,
      "statement"
    ],
    "after": "Historia humana raw separada de interpretación y automatización; LIB-R10, T2805/T2807 con T2803/T2806.",
    "layers": 4,
    "before_utf8_sha256": "b369f3a8ee7c4b16dbdc46c81b49e3c80937b3673c685b2ddd2d9f1e71ad7cab"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      0,
      "stimulus"
    ],
    "after": "Colisión o alteración de artefacto al reanudar",
    "layers": 4,
    "before_utf8_sha256": "4ad9712615879fabf1fcf9b28637088a798d662ff22393046ae187eb87abafcf"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      0,
      "response"
    ],
    "after": "Rechazar sin sobrescritura ajena y explicar recuperación",
    "layers": 4,
    "before_utf8_sha256": "e844467886b88b2eb2d56d98f03f55175fc4beef366b9246c9b98202d901b9ab"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      0,
      "measure"
    ],
    "after": "Exit no cero esperado, hashes originales y cadena sin modificación no autorizada",
    "layers": 4,
    "before_utf8_sha256": "4f92a71d8fa04923fb347876a491adc7519760a69177f1389d73fca64c7384c6"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      1,
      "source"
    ],
    "after": "Importación histórica solicitada por el usuario",
    "layers": 4,
    "before_utf8_sha256": "fb5159f73e8f9889909014b44a4fce0b69991a3c8625a3d6abcf92d340eaff35"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      1,
      "environment"
    ],
    "after": "Corpus autorizado futuro; ningún chat privado usado en este mantenimiento",
    "layers": 4,
    "before_utf8_sha256": "08fb199f8ea8ed96d1ae6bc79f8993961ab2697d3d59c098bc6809df01cd1440"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      1,
      "response"
    ],
    "after": "Conservar fuente; separar hechos e interpretación; cero efectos live",
    "layers": 4,
    "before_utf8_sha256": "c58e6820804243b817fa310293a607c629bcbcc368628e6d472296a8068adabc"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "quality_scenarios",
      1,
      "measure"
    ],
    "after": "Comparación byte a byte, lineage por derivación y ausencia de jobs/envíos/tools",
    "layers": 4,
    "before_utf8_sha256": "8aa3781bba90bb0cfea0ab739de02ed449392c51af172c5a6c1d81f334d1d9bb"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "context"
    ],
    "after": "Owners Markdown existentes -> materialización en staging -> herramientas/verificadores -> evidencia -> checkpoint. No subsistema paralelo.",
    "layers": 4,
    "before_utf8_sha256": "6a59ec6124387cb0312419f689243672a60bd14b929733ac738180bbd7c44250"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "components",
      0,
      "name"
    ],
    "after": "Biblioteca canónica",
    "layers": 4,
    "before_utf8_sha256": "4c1f954de34c9a54cb5f3fe30e7406df8aad2540ee4411bf8514387fb4364482"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "data_flows",
      2,
      "from"
    ],
    "after": "Histórico autorizado futuro",
    "layers": 4,
    "before_utf8_sha256": "3cc9ad84969cacd76272042070b4b3dae06fd1d5cec6f9acb282efc496e0bfa4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "data_flows",
      2,
      "data"
    ],
    "after": "Original separado e interpretación trazable, no live inbox",
    "layers": 4,
    "before_utf8_sha256": "5285e5cf0279ba433e09fed9e536c3ffc66a9315a7fa468b12b7d0c966a8744a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "decisions",
      0,
      "context"
    ],
    "after": "Continuidad del plan del usuario sin rediseño ni gates ficticios",
    "layers": 4,
    "before_utf8_sha256": "3ae3b884c2929fbade8c092e90cae3418261a152276f99c52ca7d5051af67947"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "decisions",
      0,
      "decision"
    ],
    "after": "Mantener roadmap único y añadir trazabilidad requisito-prueba-evidencia en el contrato canónico",
    "layers": 4,
    "before_utf8_sha256": "25b6a610b4f616493f18d3f3201592875d8993c597ddd1fe31f93b834ccf1c15"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "decisions",
      0,
      "consequences",
      0
    ],
    "after": "No sustituye readiness ni aceptación de producto",
    "layers": 4,
    "before_utf8_sha256": "f3ad832c4a0d54526a28364c99ccb4e19fcb52ae6d3c01842a2e40b0dc2ce2d7"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "architecture",
      "decisions",
      0,
      "rollback_trigger"
    ],
    "after": "Contradicción con owner canónico: registrar fallo y corregir el delta, sin reescribir evidencia",
    "layers": 4,
    "before_utf8_sha256": "2d2c686885cb6f1f056a90a5d767c08dc5190195b4045b9f8728a270c881595c"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "api",
      0,
      "name"
    ],
    "after": "CLI de validación",
    "layers": 4,
    "before_utf8_sha256": "46e37a6bc83e53eeca4b56d45c11edc7dbd19d10f258aa0cbdeb637599d67127"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "api",
      0,
      "input"
    ],
    "after": "Contrato/estado y nivel explícito",
    "layers": 4,
    "before_utf8_sha256": "4d802ce9d5cbb5a55cb981cc98ff263406b056f0f67bf713c7b8a8df2f1ca803"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "api",
      0,
      "output"
    ],
    "after": "Resultado y diagnóstico",
    "layers": 4,
    "before_utf8_sha256": "bb431bb1a0389c9f98092fb0fa2905fce93556a8c37a05c3d79d648adde8a99e"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "api",
      0,
      "errors",
      0
    ],
    "after": "plan válido no equivale a evidence aprobado",
    "layers": 4,
    "before_utf8_sha256": "4fc596e9bf8a1ca811105a0fcdf8b782b5baef6fe98844ab8460893df7ab0c2f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "data",
      0,
      "invariants",
      1
    ],
    "after": "Original raw no se reescribe por normalización; retención y acceso deben definirse",
    "layers": 4,
    "before_utf8_sha256": "e148128f5477c29f2b6748297831811855504ab9b4b5319149ea819d59fb4485"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "data",
      0,
      "invariants",
      2
    ],
    "after": "Un hash demuestra identidad de bytes, no verdad semántica ni exhaustividad del proveedor",
    "layers": 4,
    "before_utf8_sha256": "5564dea40713a34f01d1d0a3bd55832c7569258820931d8de5b1587329196f90"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "lifecycle",
      0,
      "allowed",
      0
    ],
    "after": "plan validado sin promoción",
    "layers": 4,
    "before_utf8_sha256": "832fe45baf014ad4d805c30ed27fb4f56b2e029a7f36214cf290356e5a6ba0f7"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "contracts",
      "lifecycle",
      0,
      "allowed",
      1
    ],
    "after": "evidence sólo después de pruebas y artefacto/release exactos",
    "layers": 4,
    "before_utf8_sha256": "794a240020b4fd337fb08285dba421f9ab5fa4c79db4644dbfb0d767013feb8e"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "risks",
      0,
      "failure_or_threat"
    ],
    "after": "Promover por estructura/reputación o incorporar código incompatible",
    "layers": 4,
    "before_utf8_sha256": "dea54b530321d6b27e75c4f722b82114d77b168ecafecf00e1e3484c611d3aa9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "risks",
      0,
      "controls",
      0
    ],
    "after": "Admisión oficial por claim",
    "layers": 4,
    "before_utf8_sha256": "3b249473bcd6556153cafba37c73d7a6ba33d3ae827aaae5a2537bcc3278aaa5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "risks",
      1,
      "failure_or_threat"
    ],
    "after": "Perder cambios, filtrar datos o reenviar históricos",
    "layers": 4,
    "before_utf8_sha256": "cbcff8b949fc9a36e557551dbfe8604efb0ce16f1ba6c51b8e494f1204e8be50"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "risks",
      1,
      "controls",
      1
    ],
    "after": "Preservación del original y no backfill al runtime live",
    "layers": 4,
    "before_utf8_sha256": "8f678a3302f42d2fbe65de50f6cfd086e3b0dbbc03cd5d356d0e8da3075cfe39"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "risks",
      1,
      "controls",
      2
    ],
    "after": "Fallo y recuperación trazables",
    "layers": 4,
    "before_utf8_sha256": "5d488812465447c43409a9c978b0b53078260131ed813b566b405e9654f2da47"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "slos",
      0,
      "window"
    ],
    "after": "Cada ejecución y candidato de release",
    "layers": 4,
    "before_utf8_sha256": "e8e0dbd62bfb8d39353e136345a8bd72dea75f9803f4f159896ee333516dcb18"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "slos",
      0,
      "budget"
    ],
    "after": "Cero tolerancia en estas invariantes; no es disponibilidad de producción",
    "layers": 4,
    "before_utf8_sha256": "934b9097a27304f8597a2dfc6f75b5473f105d530c35e5c3c6c7726f032087f6"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      0,
      "setup"
    ],
    "after": "Fixtures sintéticos NEW/EXISTING, PowerShell 7.6.5 y Python 3.14.4; V287.",
    "layers": 4,
    "before_utf8_sha256": "23408d2b67a6f71ce048055c832b1c7f10aa24678ff9576e9df737a835f426ec"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      1,
      "oracle"
    ],
    "after": "Cada claim tiene fuente estrecha, licencia, revisión, limitaciones y prueba; ningún CANDIDATE promovido por reputación.",
    "layers": 4,
    "before_utf8_sha256": "b4260997ac7aa90a132ea85c8dee402a79956ea5d8622a52fd4577ef370cd251"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      2,
      "setup"
    ],
    "after": "Tooling y grafo exactos del candidato; equivalencia 385 abierta; privacidad local 386 corregida V292, integración pendiente.",
    "layers": 4,
    "before_utf8_sha256": "543fe528532055528f2ca2f12f71e7be05d0daa45f515dbd29f00e6ff3df16a9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      2,
      "action"
    ],
    "after": "Ejecutar SCA, límites de rutas, datos sensibles y trust boundaries según T2803.",
    "layers": 4,
    "before_utf8_sha256": "a758d8030feeaf126a8d6b48089cdc7793f44b252fe2db2d57b2ead3d9e0d59d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      2,
      "oracle"
    ],
    "after": "Sin findings bloqueantes ni PII en observabilidad; no reutilizar un PASS histórico sobre un grafo diferente.",
    "layers": 4,
    "before_utf8_sha256": "e98724d015536aab2695c0673f38881e6c120d96b26ab68635f97759831244c8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      3,
      "setup"
    ],
    "after": "Manifest/checkpoint y procesos de materialización; staging aislado.",
    "layers": 4,
    "before_utf8_sha256": "d60088817b507413741d77a0d168e4755331953f62b2a597bcb0bf7e800b1c06"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      3,
      "action"
    ],
    "after": "Inyectar interrupción y fallo de escritura en recorridos afectados; repetir regresiones y verificar reinicio.",
    "layers": 4,
    "before_utf8_sha256": "a31819e9397dcec69968856fef6e3ec5cc0773003e8c77c66ca9b44e0f61b9cb"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      3,
      "oracle"
    ],
    "after": "Sin promoción, pérdida o doble efecto; originales y fallos conservados; V284/V287 sólo cubren subconjuntos.",
    "layers": 4,
    "before_utf8_sha256": "ec5a9d09913cce0cf172894c37fa63f4d1f39260a736cd9c742604cf666eb005"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      4,
      "setup"
    ],
    "after": "Eventos y métricas del recorrido CLI y futuro montaje runtime T2809.",
    "layers": 4,
    "before_utf8_sha256": "46cb0ea7fb971f1d231271ff4d117e1f6cb68ab72244b0bc65808663de4cb8a5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      4,
      "action"
    ],
    "after": "Probar correlación, diagnóstico comprensible, redacción, monotonicidad y señales de fallo del sink.",
    "layers": 4,
    "before_utf8_sha256": "6126c8c547491272b77479d9fa1932552f2421a3e087f0c3cf6b40305660eda0"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      4,
      "oracle"
    ],
    "after": "Verifica privacidad/operación integrada del target; reutiliza TEST-10 contador y TEST-11 política local sin confundirlos con exportación, retención o salud del sink.",
    "layers": 4,
    "before_utf8_sha256": "b881a34be31214f1ed2a0c6d4de6c8dd27ca34edbcaf357d28711b7805fcd0b5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      5,
      "setup"
    ],
    "after": "Release candidato y snapshot previo de instalación NEW/EXISTING.",
    "layers": 4,
    "before_utf8_sha256": "a760233d7be6577ee1194d113ee636de00f3ad7de3d78bca588f55292766e532"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      5,
      "action"
    ],
    "after": "Probar fallo parcial, restauración exacta y reanudación desde el artefacto distribuido.",
    "layers": 4,
    "before_utf8_sha256": "8c6fe5cf0216d6ff4174b0e3bf92ef81db7cddb69627d8e19b13cd865f1a2f22"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      5,
      "oracle"
    ],
    "after": "Archivo previo y cadena de eventos preservados; recuperación medible; no extrapolar a PostgreSQL/PITR.",
    "layers": 4,
    "before_utf8_sha256": "0917b2bc0fe62bdd510327a2ddc3897bcb44b470deac18fdeb0cb6fdd7982565"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      6,
      "action"
    ],
    "after": "Dos builds independientes, comparación byte a byte, SCA, notices/SBOM, firma y verificación independiente.",
    "layers": 4,
    "before_utf8_sha256": "03877f6e39f13afc514e797d9bcb212d047613c273fd651d1580ba1b4bebeafc"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      6,
      "oracle"
    ],
    "after": "ZIP reproducible y sin expedientes locales ni aprobaciones heredadas; sólo después de cerrar gates materiales.",
    "layers": 4,
    "before_utf8_sha256": "ac5c720be776525e19933b411b283b523a9519746a91ffd1965128a9666d7b11"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      7,
      "setup"
    ],
    "after": "Usuario/agente NEW/EXISTING y guía/CLI de la misma revisión.",
    "layers": 4,
    "before_utf8_sha256": "17b85c17eab7cb336eb5f7f8732a98b492b20c635640440b7ccae6845d6d537d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      7,
      "action"
    ],
    "after": "Recorrer selección, diagnóstico de entrada incompleta, ayuda y recuperación con artefacto exacto.",
    "layers": 4,
    "before_utf8_sha256": "720dcc8f3e0dd78540e84c786d038ca76dffe4b59b40a49aad747c0e733b1325"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      7,
      "oracle"
    ],
    "after": "Instrucciones ejecutables y errores accionables; versión de ayuda consistente; sin lectura masiva innecesaria.",
    "layers": 4,
    "before_utf8_sha256": "388448e24d477fcd169c56cada47e672a921531e9d85b181f2bc9a3997960698"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      8,
      "setup"
    ],
    "after": "Corpus histórico autorizado y contrato oficial del canal: todavía no disponibles.",
    "layers": 4,
    "before_utf8_sha256": "69be64755770d818f7fd3c88d831ca6525f7ee65f1821e54dbb2bf59243c2434"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      8,
      "action"
    ],
    "after": "Verificar bytes originales, lineage derivado, aislamiento, importación sin envíos/tools y empalme histórico/live.",
    "layers": 4,
    "before_utf8_sha256": "ff0e9730c44d4476d579fa8f62b6aab439e6b6136f3c2def34cdb78712b101d4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      8,
      "oracle"
    ],
    "after": "Cero modificación del original por transformaciones; duplicados y conflictos trazables; ambigüedad no se inventa; cero efectos de backfill.",
    "layers": 4,
    "before_utf8_sha256": "e8137dbe58d184c7239da962515efc6dd9cee7fe7f1caeeb91aac6cc1178b0bf"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      9,
      "action"
    ],
    "after": "Negativos/overflow/concurrencia, oráculo big.Int y fuzz nativo acotado; vet/build.",
    "layers": 4,
    "before_utf8_sha256": "2a59077958aed169a994b23fffc085c36dc92013803a336cf4f5c6fe3263ee02"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      9,
      "oracle"
    ],
    "after": "Rechazo sin mutación/panic; contador monotónico. No prueba privacidad, race ni operación.",
    "layers": 4,
    "before_utf8_sha256": "6dae1ab8bc6692ed525589052976e703c51a61d3c392f6b28c362983660ad79f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      10,
      "setup"
    ],
    "after": "GO-OBSERVABILITY-CORE 0.2.0 aislado, política sintética cerrada, Go 1.26.7 JSONHandler y reconstrucción2/2.",
    "layers": 4,
    "before_utf8_sha256": "78dd1bbf4c01e14e798483b2e1c8dd704ca3fee7ea8fc65184954359a64edf22"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      10,
      "action"
    ],
    "after": "Baseline tres exposiciones, política→TryEmit→JSON real, negativos/snapshot/concurrencia/errores y fuzz dos targets; vet/build.",
    "layers": 4,
    "before_utf8_sha256": "b6981a54c0e91a901a266dd7659471d09f2ae3e3c195a709f992ebe62c7df982"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      10,
      "oracle"
    ],
    "after": "Eventos aprobados llegan íntegros; entradas no admitidas no llegan al sink ni aparecen en errores. No prueba política empresarial, exportación, race, SCA ni operación target.",
    "layers": 4,
    "before_utf8_sha256": "52971b13b6f593c5839e502e475b261fb7c4e75eb8d7e793d6d60391bc35f0fd"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      11,
      "setup"
    ],
    "after": "V294 reconstrucción canónica 67/742, PostgreSQL18.6/52 migraciones, Next/BFF/OIDC/JWKS sintéticos reales.",
    "layers": 4,
    "before_utf8_sha256": "6c4b00af971e67f9bd5121756745752f61e9041c2bce5d00899f1cd6e259e320"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      11,
      "action"
    ],
    "after": "4 proyectos navegador x3 fases: aceptación, concurrencia, duplicados, permisos, respuesta perdida, GET503 y reinicio Next/recreación API/pool.",
    "layers": 4,
    "before_utf8_sha256": "a26ffce3a6a56447ce578ff47409f6177ae0c123b30466d8037884f7c4a1261f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      11,
      "oracle"
    ],
    "after": "Pedido único y referencias durables, 3 pedidos/3 aceptaciones/6 eventos por tenant; rechazo no crea pedidos. No prueba stock/pago/entrega, WCAG integral, proveedor ni producción.",
    "layers": 4,
    "before_utf8_sha256": "a5b2797298c0cb71353ac8e9f8ea3c68a2fcacb5581996642e80754a0cdf74a8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      12,
      "setup"
    ],
    "after": "V294 mismos pedidos del navegador, Commerce PostgreSQL existente; fixture sintético sin worker.",
    "layers": 4,
    "before_utf8_sha256": "504460bd395481586601f62aabeecf55522b34a80b2314f4eafa085053852f73"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      12,
      "action"
    ],
    "after": "Dos pedidos compiten por stock; clave/importe/moneda/org de solicitud pago; lectura con pool nuevo y snapshot antes/después de reinicio PostgreSQL.",
    "layers": 4,
    "before_utf8_sha256": "3c78e0108b8121f80036074526f95be8030e1442775a96efe1ed93bf5ad96364"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      12,
      "oracle"
    ],
    "after": "Una reserva y una solicitud created, dos eventos asociados; snapshot idéntico. No UI operativa, ACL HTTP Commerce, cobro, proveedor, PITR ni DR.",
    "layers": 4,
    "before_utf8_sha256": "ea2e6261d0a4105d1beff914b7b9c91070b442d6d44e4cc9d8da34a89f5506da"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      20,
      "setup"
    ],
    "after": "V302 fuente canónica4/4 recompuesta; Next construido, Go1.26.7, PostgreSQL18.6 sintético, cuatro proyectos Playwright.",
    "layers": 4,
    "before_utf8_sha256": "5cc1bae9978e48f81a744dcaca85aa66cede2f241aebf872483ab145a77e70a2"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      20,
      "action"
    ],
    "after": "Baseline UI rojo; matriz52 fases y agenda4,105 web PASS/1 SKIP, Go test/vet/build. Rol, consulta503, recuperación GET, permisos HTTP y cero efectos prohibidos.",
    "layers": 4,
    "before_utf8_sha256": "b7102daa3a500b67be51b5b283631b921d5e2938ab5b1a6a8b78cad54a622b15"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      20,
      "oracle"
    ],
    "after": "Cuatro roles sólo acceden a secciones permitidas; pedidos sobreviven fallo independiente, sin detalle privado ni POST de recuperación; HTTP403 y cero checklist. FAIL457 sigue abierto.",
    "layers": 4,
    "before_utf8_sha256": "d70987fa5870b304616b88a421d60a2ccb722eb2991c540217e0acd72f5a4ee4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      21,
      "setup"
    ],
    "after": "V303 fuente canónica4/4 recompuesta, Next construido, Go1.26.7/PostgreSQL18.6/OIDC sintético, Playwright cuatro proyectos.",
    "layers": 4,
    "before_utf8_sha256": "a1a4bd5d31991497fb107fe3f18d3469a49017898075f848162b430380df3593"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      21,
      "action"
    ],
    "after": "Rojo post-commit V302; tres resoluciones por proyecto con respuesta perdida, JSON vacío, consulta503, storage ausente, carrera200/409, replay/permisos.56 fases, agenda4, Go/web/build.",
    "layers": 4,
    "before_utf8_sha256": "6e5256fc7cb5e5adc62ba8c0e7f34c072a3732b23049163bc4ea0057b57dbf4c"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      21,
      "oracle"
    ],
    "after": "Resoluciones3, sucesor1, autorizaciones2 y eventos3 por tenant; marker sin notas, consulta sin reenvío y estado recuperado. No efectos downstream completados ni readmisión del sender genérico FAIL457.",
    "layers": 4,
    "before_utf8_sha256": "18f413c1f5234886e88b956db69c78f4ba81529d34af0e1d1a6b64c08d260b3e"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      22,
      "setup"
    ],
    "after": "V304 fuente canónica8/8, web construida, Go/PostgreSQL local e issuer sintético, cuatro proyectos Playwright.",
    "layers": 4,
    "before_utf8_sha256": "929a94ec487a5bb46e3c830580421fef7190b9be5cd2896f7eeec001f4f1ef80"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      22,
      "action"
    ],
    "after": "Asignación y tres transiciones de lead; respuesta perdida, JSON vacío, consulta503, storage sin POST, concurrencia200/409, permisos/actor falsificado y rollback outbox.60 fases más agenda4, PostgreSQL3 y Go/web.",
    "layers": 4,
    "before_utf8_sha256": "8c8ad7f2be8dbdd7b517242e6508cbac52e28be53b75c4bca1a632da3e0db1d0"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      22,
      "oracle"
    ],
    "after": "Lead lost/version5, cuatro eventos con actor HTTP por tenant, consulta sin reenvío y ausencia de actualización parcial; métodos sin actor rechazados. No readmisión del sender restante ni producción.",
    "layers": 4,
    "before_utf8_sha256": "5f0526be6299a59fd1c8a4bc35965abf3d2f51b6eeb02e5ed7b7ad05a6f4c488"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      23,
      "setup"
    ],
    "after": "V305 seis archivos canónicos recompuestos; misma web y backend verificados por SHA. PostgreSQL18.6/OIDC local y Playwright cuatro proyectos.",
    "layers": 4,
    "before_utf8_sha256": "34eceb41ad8b225e18a783e7ec067e4a88b2c81026d337fe6f37bd0e25d7d09a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      23,
      "action"
    ],
    "after": "Cancelación working/unavailable y rechazo por cita activa; pérdida/JSON/GET503/storage/carrera/permisos. Coherencia publicación/reserva con jornada actual; carreras cancel/block-v-book.64 fases, agenda4, focales10/PG3 y Go/web.",
    "layers": 4,
    "before_utf8_sha256": "f18fa2bcecb93f5e4287bac2393e98119c0a39817cae5da144f4537bda28ade0"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      23,
      "oracle"
    ],
    "after": "Cuatro cancelaciones/version2/eventos con actor por tenant, quinta active/version1 sin evento; dos cupos retirados y dos reservas HTTP rechazadas sin efectos. Lock cooperativo evita éxito simultáneo de reserva y bloqueo/cancelación en las pruebas; no producción.",
    "layers": 4,
    "before_utf8_sha256": "8f3726d7ff1397e86b62299adba837a9806ef5523913224196f19929c7fc3499"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      24,
      "setup"
    ],
    "after": "V306: cotización con clave UI/BFF/HTTP, actor autenticado y consulta scoped.",
    "layers": 4,
    "before_utf8_sha256": "efa71c551168e88536363107c1d1cf11f72589d9bf1f3e20ad9c3a718ac112c8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      24,
      "oracle"
    ],
    "after": "3 cotizaciones/keys/eventos/actores por proyecto; snapshot reiniciado idéntico; no reemisión automática ni producción.",
    "layers": 4,
    "before_utf8_sha256": "ad242c449121d152dbbacb458943f8d576f9e31a689541bbf05a69b74bd6e880"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      25,
      "setup"
    ],
    "after": "V307: alta de intervalos con referencia durable y preparación explícita tras consulta.",
    "layers": 4,
    "before_utf8_sha256": "b8671469f87e607a7884d1e07f063c04689c1c553435009304b646ef19e597a5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      25,
      "action"
    ],
    "after": "72 fases/4;110 web PASS/1 SKIP; carrera8, PG3, aislamiento, corrupción, rollback y reinicio; rebuild9/9.",
    "layers": 4,
    "before_utf8_sha256": "b89ff48d18c24552037deb00171cf6e5032a34feb1fae5a6c71239cea0941945"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      25,
      "oracle"
    ],
    "after": "1 creación/7 replay en carrera;3 intervalos/keys/eventos/actores por proyecto; cancelado no resucita, sin cambios de agenda.",
    "layers": 4,
    "before_utf8_sha256": "b027166f71423f653f93776477b57ce7140325dfb09264fb69222a884711a0ab"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      26,
      "action"
    ],
    "after": "76 fases/4 más agenda4;111 web PASS/1 SKIP; PG3 carrera8/rollback/aislamiento; rebuild9/9 y reinicio.",
    "layers": 4,
    "before_utf8_sha256": "8468bc904b43e3a4b456a729741c0cff2f70ef9820418d5d502c0db60f4afe67"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      26,
      "oracle"
    ],
    "after": "1 creación/7 replay;3 recursos/keys/eventos/actores y6 habilidades por proyecto; recurso inactivo no se reactiva.",
    "layers": 4,
    "before_utf8_sha256": "a8abafe3a30f67c83ab6e21f00d44fc29702f66ed673e9260b90cb95b1093814"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      27,
      "setup"
    ],
    "after": "V309 nueve archivos canónicos recompuestos; UI/BFF/HTTP/PG con fixtures y jornadas explícitas.",
    "layers": 4,
    "before_utf8_sha256": "7bebe14305b3bd7b91a04f8bf51450e3e1c5634849f453bf6fc24a9761723c69"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      27,
      "action"
    ],
    "after": "80 fases/4 proyectos, agenda4,112 web PASS/1 SKIP, PG3 carrera8/atomicidad/lectura llena y cerrada; reinicio real del cluster sintético.",
    "layers": 4,
    "before_utf8_sha256": "e8abf30d60d8588449c8fea294a9fe1d3e78c5ae8611593126b054a53de931c8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      27,
      "oracle"
    ],
    "after": "1 creación/7 replay concurrentes;3 cupos/keys/eventos/actores por proyecto. Lookup recupera2/2 ocupado y cerrado sin reabrir ni depender del listado público. Rebuild9/9; no cambios comerciales ni producción.",
    "layers": 4,
    "before_utf8_sha256": "14bfbb607ca60c160a4daf049735586370a0aa509ac4423e12ccb020f367fafa"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      28,
      "setup"
    ],
    "after": "V310 diez archivos canónicos recompuestos; UI/BFF/HTTP/PG con fixtures sintéticos.",
    "layers": 4,
    "before_utf8_sha256": "ea0fe0d1a51a69f7e836168625645dc8afbe1996045268476ad58c351558d5fe"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      28,
      "action"
    ],
    "after": "84 fases/4 proyectos,agenda4,113 web PASS/1 SKIP,PG3 carrera8/atomicidad/inmutabilidad; reinicio real del cluster sintético.",
    "layers": 4,
    "before_utf8_sha256": "6c91cbb29602c551339cdc4d6160dde13bdd2aa868c5663a5b4d9e58131e6009"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      28,
      "oracle"
    ],
    "after": "1 publicación/7 conflictos concurrentes;3 versiones/6 ítems/3 eventos/3 actores por proyecto. Consulta exacta confirma contenido por huella, conserva versión previa y bloqueo si diverge. Rebuild10/10; no cambios comerciales ni producción.",
    "layers": 4,
    "before_utf8_sha256": "9f3e29437133e3800926d18b1f3de25416440d813451ce762a450129da1b815f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      29,
      "setup"
    ],
    "after": "V311 diez archivos canónicos recompuestos; fixtures explícitos de entregas prepared, UI/BFF/HTTP/PG.",
    "layers": 4,
    "before_utf8_sha256": "92ff33e859902ff6a3111c1d326dd3f36feaa839c8eed7dc8b44f47ef319b28f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      29,
      "action"
    ],
    "after": "88 fases/4 proyectos,agenda4,114 web PASS/1 SKIP,PG3 carrera8/atomicidad/inmutabilidad/estado posterior; reinicio real del cluster sintético.",
    "layers": 4,
    "before_utf8_sha256": "ec0320975fe96e2c4e77f33c2e620c7c6d2cf9902eaf8be785de15769f703282"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      29,
      "oracle"
    ],
    "after": "1 completado/7 conflictos concurrentes;3 entregas/6 respuestas/3 eventos/3 actores por proyecto. Consulta exacta compara huella y conserva actor/fecha aunque estado actual sea accepted/rejected. Rebuild10/10; sin regla comercial nueva ni producción.",
    "layers": 4,
    "before_utf8_sha256": "f02bb6a036f38b33f01c29869ecfdcc1abbad2043e1d143847e1f683793fdb61"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      30,
      "setup"
    ],
    "after": "V312 diez archivos recompuestos; PostgreSQL sintético, portal/BFF/HTTP, kit Go fuzz separado admitido.",
    "layers": 4,
    "before_utf8_sha256": "fe2797d78541a00a4b4d257681a984ebd9a5cd6eb032f1d7aace8ac972111eb6"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      30,
      "oracle"
    ],
    "after": "3 recepciones/3 decisiones/11 solicitudes/6 eventos por proyecto; receive y decide concurrentes1 éxito/7 conflictos. Scope cruzado0efectos, consulta conserva actor/evidencia incluso sin lista. Rebuild10/10 y snapshot idéntico; ningún efecto downstream ejecutado.",
    "layers": 4,
    "before_utf8_sha256": "b78d83778789ebb1782de162c027194e3da22c08f4dd76cdefde8ac9bbfee910"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      31,
      "setup"
    ],
    "after": "V313 reconstrucción67/745; dos módulos Go oficiales actualizados, cinco archivos canónicos y proyecciones SCA trazables.",
    "layers": 4,
    "before_utf8_sha256": "3820a3555c4c6d3cccba5318a4d811b10c2252e5d9a1b833732046805afac4ed"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      31,
      "oracle"
    ],
    "after": "GO-2026-5970/6179/6180 resueltos en locks; scan21 fuentes/465 ocurrencias/347 paquetes/0findings.47 wheels verificadas,134 requisitos compatibles.5/5 rebuild y2/2 binarios idénticos entre candidate1 y sucesor; aislamiento de refund probado. No runtime/deploy/security integral.",
    "layers": 4,
    "before_utf8_sha256": "ffffa2671706a8921bf33efcab775b1cad494d5587f10e9520a42dce970c29f4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      32,
      "action"
    ],
    "after": "Reproducir salida prematura, verificar drenaje/respuesta/admisión/deadline/listener/startup, ciclos inválidos y Go test/vet/build.",
    "layers": 4,
    "before_utf8_sha256": "72043c950eac16e6550af43778755fe5584ca5f565328a05ee51cd1a2a5608c2"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      32,
      "oracle"
    ],
    "after": "Red falla por retorno antes del drenaje; cinco tests raíz x3 PASS en ambos árboles,4/4 reconstrucción. Deadline fuerza cierre y devuelve fallo; listener error preservado tras respuesta. No señales OS/target/hijacked.",
    "layers": 4,
    "before_utf8_sha256": "ec34260c63879f01397abf338a257a461a9ad3dcc3d84864fef36e69d826ea3a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      33,
      "action"
    ],
    "after": "Reproducir28 bypass de referencias y error de acceso, verificar12 tests y4 invocaciones CLI; rebuild canónico2/2.",
    "layers": 4,
    "before_utf8_sha256": "6ab50e6b55f83f18543b48cefef31c22f69f82d0226ec7e7c02a690b91460b98"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      33,
      "oracle"
    ],
    "after": "Ambos modos exigen rutas SBOM/provenance; files rechaza ausencia/directorio/escape/NUL/acceso, completo PASS de presencia. CLI2errores exactos al omitir ambos; no aprobación semántica.",
    "layers": 4,
    "before_utf8_sha256": "2c610ec892c52a43ec5bc5831c3e6d2036fd93419664d13c03f17f7fc0de816b"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      34,
      "setup"
    ],
    "after": "V316 core0.4.4/workers0.3.1 compuestos67/745; dos DB dedicadas con fundación0001 y publisher sintético.",
    "layers": 4,
    "before_utf8_sha256": "f55fdd178685282dd478d3615195892527cf3531dea4c4a1464fef3f53c40e16"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      34,
      "action"
    ],
    "after": "Seis rojos SQL de expiración/generación/row lock y tres publisher; RemainingLease y retry real de store con clave estable;33 pruebas por árbol,4/4 rebuild y Go gates.",
    "layers": 4,
    "before_utf8_sha256": "51657c636301e79a06569e1117c402789a933aee3fdd2623387bb0bd71d92766"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      34,
      "oracle"
    ],
    "after": "Confirmar/reprogramar sólo con generación y lease vigente tras lock; publisher sin presupuesto no se ejecuta. PostgreSQL conserva retry y2 publicados al reintentar con EventID estable/gen1→2. No exactly-once remoto ni supervisor.",
    "layers": 4,
    "before_utf8_sha256": "33730157368c57a1b84a438bee0690db87a3c120524a87387efe1669d99ab64a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      35,
      "setup"
    ],
    "after": "V317 composición67/745, Go1.26.7 observado y árbol íntegro, DBs loopback, web exacto y Playwright0.1.23.",
    "layers": 4,
    "before_utf8_sha256": "70b297bac51404c5a66d3e63100da2560c4bbd909a4823b9b72007d7e3b9fb6c"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      35,
      "oracle"
    ],
    "after": "Identidad y atribución exacta sólo por receipts;2 paquetes OSV por versión,0matches y control2/41;33outbox+15HTTP+33PG+6refund,115web/1skip,92browser+4agenda y467038fuzz. Esperar GET RSC antes de cambiar cookies,6POST y lector rechazado.",
    "layers": 4,
    "before_utf8_sha256": "d9770169d8b4e483500b41fc348430ae07ab95a8ae9f6b5b0156215798310935"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      36,
      "setup"
    ],
    "after": "V318 refund0.1.3 aislado, Go1.26.7 observado y DB sintética exclusiva; causas privadas sintéticas.",
    "layers": 4,
    "before_utf8_sha256": "2edc7219b5fc0c749eb12ddeab3fce8b9681a3389d8ce5eb73c7bf81fdb0f477"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      36,
      "action"
    ],
    "after": "Reproducir divulgación en startup real y log de Step; corregir salida a códigos fijos;6tests por árbol,6PostgreSQL y3/3 rebuild con Go gates.",
    "layers": 4,
    "before_utf8_sha256": "8868690fe40da0c11c25324702f70d8cd9602277f15b69bab4ef058919b7a141"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      36,
      "oracle"
    ],
    "after": "No llamar Error() ni imprimir marcador privado; preservar exit1, códigos operativos fijos e idle silencioso. Binario final replica el probe original sin divulgación. Retry y resultados durables intactos.",
    "layers": 4,
    "before_utf8_sha256": "834839c260c47efaf1b210490fde312f7a36784c0bce10b4138c7e889a06da4b"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      37,
      "setup"
    ],
    "after": "V319 Windows aislado; Node24.20.0 firmado e independiente, Go1.26.7 observado, pnpm11.19.0 exacto, PostgreSQL sintético.",
    "layers": 4,
    "before_utf8_sha256": "55f8bfa3dda8d91f78195ee99bf67de1734ac2913f4f511d71c7c646c88be2d9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      37,
      "oracle"
    ],
    "after": "22CVEs del control Node,0 del candidato en snapshots exactos; ZIP145/4/9 rechazado; controles vacíos/timeout no se admiten como gate;2selecciones115PASS/1SKIP,23fases por4proyectos+agenda4;746/746;256GET/8arranques sin regresión local.",
    "layers": 4,
    "before_utf8_sha256": "7d9f8b0cd049ef08f8c06b5a6254b13e3ed3075a167f592adf937a0b57f41554"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      38,
      "action"
    ],
    "after": "Preservar matching oficial, limitar lectura/ejecución, reproducir bridge sustituido y artefacto diferente;30tests sobre selección final/rebuilt,4CLI,10/10 identidad y SCA CLI2/0.",
    "layers": 4,
    "before_utf8_sha256": "4895a246462fdb0a8dab914531aee0355189b4e03a3d035b2d45c60b31b3496a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "tests",
      38,
      "oracle"
    ],
    "after": "Vacío/alteración/stale/future/env hostil/licencias faltantes no producen PASS; timeout/output real finitos, causa privada estática y receipt no sobrescrito; control vulnerable del motor true, currentNodePASS; G0–G8 instancia ready=true sin monitor/producción.",
    "layers": 4,
    "before_utf8_sha256": "4fffc524477f5ae4c00c64fcc3505a279616f0bc94d661315c1e21706741705a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "hypothesis"
    ],
    "after": "La entrada por mapas y materialización selectiva reduce trabajo repetido sin perder controles",
    "layers": 4,
    "before_utf8_sha256": "665579561a0982b48440a82053d44625190263ee17453a4ca6eff2e7d7861122"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "environment",
      "toolchains"
    ],
    "after": "PowerShell 7.6.5/Python 3.14.4 en V287; registrar exactos en nueva medición",
    "layers": 4,
    "before_utf8_sha256": "7446665d1796849744ceb6bcc100438ac062abe54b3c1a65a68efb52c6381aed"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "environment",
      "hardware"
    ],
    "after": "No caracterizado; pendiente para comparación",
    "layers": 4,
    "before_utf8_sha256": "8bc6edb7717687895a3b43ba9c9c2faf1cfa6358b5e462935f7a9cf9673d9340"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "workload",
      "candidate"
    ],
    "after": "NEW/EXISTING de maintenance-entry más contexto efectivo de agente",
    "layers": 4,
    "before_utf8_sha256": "ad970cd14329e0dda4894d4f239a98358c20a359310eebdf5cd0df81816662b4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "correctness_gate"
    ],
    "after": "TEST-01 debe pasar antes de comparar; no extrapolar tiempos de procesos a tiempo de construcción del negocio",
    "layers": 4,
    "before_utf8_sha256": "3c877c8e435817ff1e39180b9e5cdb4a0fe47e66d6196a6d3c5f687bf0faadf9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "metrics",
      2
    ],
    "after": "Tokens sólo si medibles",
    "layers": 4,
    "before_utf8_sha256": "acaf03c5a8e6b30393ed2e08a26db3f2105b3b026f2a5b609433c4bd6a4c3674"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "baseline"
    ],
    "after": "OPEN: falta comparación controlada",
    "layers": 4,
    "before_utf8_sha256": "d377d941fc0e4ff42c5b0a19a0d17e372dfbf8436b55feb48d3fe62f960131e4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "benchmarks",
      0,
      "candidate"
    ],
    "after": "V287 observó 9507/8997 ms de procesos; no constituye benchmark comparativo ni medición de tokens",
    "layers": 4,
    "before_utf8_sha256": "c0ed7b798ff4b6c987786bbdd669220a1c6945dc1f83b74bc8d070af64f733bc"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "release",
      "rollout",
      1
    ],
    "after": "Probar artefacto exacto y reconstrucción independiente",
    "layers": 4,
    "before_utf8_sha256": "0a6f3a76936f0706dac33a95781054e1a245c46f8a02d0b8c50274a938009ef4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "release",
      "rollout",
      2
    ],
    "after": "Distribuir únicamente tras gates T2810",
    "layers": 4,
    "before_utf8_sha256": "06bc427bcd5e68b7850bab34a32df29def693c028a56fb4465b6de8ec6ced8f2"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "release",
      "gates",
      0
    ],
    "after": "Plan válido",
    "layers": 4,
    "before_utf8_sha256": "2bb23f8ba5eac8b44e8a361a2fe7c05f20146b05643d9f49a38e83c055578ff5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "release",
      "rollback"
    ],
    "after": "Preservar snapshot previo y original, restaurar sólo cambios propios verificados y repetir gates; no reset destructivo",
    "layers": 4,
    "before_utf8_sha256": "2b7f6a0cb8c86bc49ad8ae95d72eb94e6ea63b2d465693888f380b98ce9b162f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "operations",
      "alerts",
      0
    ],
    "after": "Todo fallo no esperado se registra inmediatamente; 384/385 y acceso 405 abiertos. 386/387 corregidos localmente; observabilidad no admitida para producción.",
    "layers": 4,
    "before_utf8_sha256": "4961eb85154eedb3017d1928b325cb20161285fa046f6782ed2d86b40b6ed412"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "codex",
      "deliverables",
      1
    ],
    "after": "Fuente canónica y evidencia reconstruible",
    "layers": 4,
    "before_utf8_sha256": "2ccd9114338fcc41820b038ba627bb0e4db785410608729683b81f61f8fb4531"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "codex",
      "deliverables",
      2
    ],
    "after": "ZIP final sólo tras aceptación",
    "layers": 4,
    "before_utf8_sha256": "43664e836007a04e7be10f4de530d3f31c8801e279ce028d156e704ce2719a5f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "codex",
      "forbidden_assumptions",
      1
    ],
    "after": "No afirmar precisión semántica perfecta ni todas las conversaciones disponibles",
    "layers": 4,
    "before_utf8_sha256": "0e96b4434cb9416da7557cccfc964e22fd1548fd08bc39e22c27f2a1c94dafeb"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "codex",
      "forbidden_assumptions",
      2
    ],
    "after": "No presentar AUTHORED como código Microsoft/Google/NASA",
    "layers": 4,
    "before_utf8_sha256": "faf0463eafccfee96c27d10378bee5ae0a120bae7c62a20db93cd044e1dbab80"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "codex",
      "forbidden_assumptions",
      3
    ],
    "after": "No afirmar coste cero, tokens ahorrados o producción por esta validación",
    "layers": 4,
    "before_utf8_sha256": "a66d1cf608abc75ad28fcdbe421b5999c02808be244c1eb976f9f24fe7bad111"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      0,
      "claim"
    ],
    "after": "Reporte histórico V287: 55 checks entrada CLI NEW/EXISTING y copia limpia; referencias a receipts originales. No prueba de release ni de negocio.",
    "layers": 4,
    "before_utf8_sha256": "bce88bb143a11bd60594c73447013cb830abd4e3aad0078223278a79f25e3c06"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      0,
      "environment"
    ],
    "after": "PowerShell 7.6.5 / Python 3.14.4 / Windows, fixtures sintéticos",
    "layers": 4,
    "before_utf8_sha256": "aca15bb4550986e7ef2aeb63eb3d4e8b04fe206e89bf5ad2dd4fd7ae0262b98a"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      1,
      "claim"
    ],
    "after": "Corrección local de contador FAIL-387, reconstrucción2/2 y pruebas focales; CANDIDATE permanece.",
    "layers": 4,
    "before_utf8_sha256": "1718ce45aea7ec7fbb984b6a7536f1e5a3d8a09f373bdbaa368b385970bf85e8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      1,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317), go test/vet/build y GO-NATIVE-FUZZ-GATE 0.1.0",
    "layers": 4,
    "before_utf8_sha256": "b4e20517d5ae0f9d3d2369fb9e28e1cf6184dee9412ccdffa764bca977ecc2d2"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      1,
      "environment"
    ],
    "after": "Windows, módulo local sin dependencias, GOMAXPROCS=2 en fuzz",
    "layers": 4,
    "before_utf8_sha256": "0da083748521d08465b68334cd01ac861f007b1ae4e16a878acce33a16b19615"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      2,
      "claim"
    ],
    "after": "Corrección local FAIL-386 bajo política cerrada, JSON oficial Go y reconstrucción2/2; CANDIDATE permanece.",
    "layers": 4,
    "before_utf8_sha256": "516f2bce4a968c0f0b049485af6eb5a82d5605b63b5a056935436a61445e01ed"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      2,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) test/vet/build y GO-NATIVE-FUZZ-GATE 0.1.0",
    "layers": 4,
    "before_utf8_sha256": "5b60406f975676042827980a9648bc768175255e511b09d20968a33d1b66c0e2"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      2,
      "environment"
    ],
    "after": "Windows, módulo aislado sin dependencias nuevas, fixtures sintéticos",
    "layers": 4,
    "before_utf8_sha256": "742f71885c5904f56ebf4c5ed70616584be2f0200caf4bab19c028d7a2182b1b"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      3,
      "claim"
    ],
    "after": "V294 aceptación browser→pedido único, UX/GET/ayuda y continuación repositorio stock/payment-created; reinicio PostgreSQL idéntico. No proveedor/cobro/UI operativa; FAIL423 abierto.",
    "layers": 4,
    "before_utf8_sha256": "a9fa88bee022e1e5c0680ff23640795ed410f97a19e013fe5971d58b00f0c974"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      3,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / Playwright1.62.1 / Next exact lock",
    "layers": 4,
    "before_utf8_sha256": "d7b326dbca52cabe11412a8ff93a1c5bb90c159b53f0222ad22c2d46c0f8ed29"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      3,
      "environment"
    ],
    "after": "Windows, PostgreSQL18.6, fixture sintético loopback",
    "layers": 4,
    "before_utf8_sha256": "16623d14dd420ad7a1e3f364622bc1378621f9cc1e244c0cb4922ec5dad8ea29"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      4,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PowerShell7.6.5 / Microsoft DPAPI",
    "layers": 4,
    "before_utf8_sha256": "8829ce6dcb32a8bd44b22a2d8ed57740dba2468731f343b5b8c517534800fb2d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      5,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Playwright1.62.1 / Next lock",
    "layers": 4,
    "before_utf8_sha256": "9202b8cb9c3fea554b4b0135cb34f41331f041965e986b573df3606c7c7f2149"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      6,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / pgx5.10.0",
    "layers": 4,
    "before_utf8_sha256": "55ff97d63e7bee0c4a4676171e3049572815f08a1715568f438e2a917253b010"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      7,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Playwright1.62.1 / Next lock",
    "layers": 4,
    "before_utf8_sha256": "9202b8cb9c3fea554b4b0135cb34f41331f041965e986b573df3606c7c7f2149"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      8,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / pgx5.10.0 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "9d3dc53a2bcb3eaecbe4aabe45d4a95ff798ca2e1334540a2b4aea3996fe5a17"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      9,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      10,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / GO_NATIVE_FUZZ_GATE0.1.0 / PowerShell7",
    "layers": 4,
    "before_utf8_sha256": "29ca4a94fdb37494dff89fe11e3a2ec1dff7e607b4c4803e996d84848589fec0"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      11,
      "claim"
    ],
    "after": "Permisos de secciones y recuperación de consultas probados localmente; no sender de comandos, creador inicial ni operación productiva.",
    "layers": 4,
    "before_utf8_sha256": "fbd940fcd6d3d1f85c0f367f92b181a05b1b3f750bb5ced8483b2da789f2676d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      11,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      11,
      "environment"
    ],
    "after": "Windows loopback, issuer y datos sintéticos, composición canónica sin servicios externos",
    "layers": 4,
    "before_utf8_sha256": "e5fd56f9c5f1434c26841dcb9ad6c0e699500f046843c1118673980eec55f719"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      12,
      "claim"
    ],
    "after": "Tres decisiones de resolución de discrepancias conectadas y recuperables localmente; no cierre del sender genérico ni producto.",
    "layers": 4,
    "before_utf8_sha256": "04b8619ac9e24dcbb57dc123c7380fc6a79bad12b8c5644de2963b854d012f9b"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      12,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      12,
      "environment"
    ],
    "after": "Windows loopback, datos y OIDC sintéticos, fuente canónica y web compilada, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "d5b94806740a366ccd663df775f541066af2115b2893d5d128f3a7671d4cf3d5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      13,
      "claim"
    ],
    "after": "Asignación y estados de lead recuperables con actor HTTP y atomicidad local demostrados; no cierre de comandos restantes ni producción.",
    "layers": 4,
    "before_utf8_sha256": "c2d2066e37463d7f485dbfbd076af7117e492f26c2111e20515337b9ac251706"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      13,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      13,
      "environment"
    ],
    "after": "Windows loopback, datos y OIDC sintéticos, fuente canónica y web compilada, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "d5b94806740a366ccd663df775f541066af2115b2893d5d128f3a7671d4cf3d5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      14,
      "claim"
    ],
    "after": "Cancelación de disponibilidad recuperable y coherencia de cupos/reservas en owners locales demostradas; no agenda completa ni readmisión de otros comandos.",
    "layers": 4,
    "before_utf8_sha256": "926477ffbef9babde681e8ca1933a6f14f8ab77786f89a9fc945fcf58f93b073"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      14,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      14,
      "environment"
    ],
    "after": "Windows loopback, datos y OIDC sintéticos, fuente canónica y web compilada, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "d5b94806740a366ccd663df775f541066af2115b2893d5d128f3a7671d4cf3d5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      15,
      "claim"
    ],
    "after": "V306: cotización con clave UI/BFF/HTTP, actor autenticado y consulta scoped. 3 cotizaciones/keys/eventos/actores por proyecto; snapshot reiniciado idéntico; no reemisión automática ni producción. No cierre integral, multi-tab, SCA, carga, restore/PITR ni target productivo.",
    "layers": 4,
    "before_utf8_sha256": "3700a044d4c493c5a47d7b8a39b2d1b1318358e8a2b25c5ead4b57e840ec99ef"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      15,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      15,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos, packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "dd79e485bb32ef5fea5df55053ff2915c2e84efdb1b711382a31ecea5c7ed136"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      16,
      "claim"
    ],
    "after": "V307: alta de intervalos con referencia durable y preparación explícita tras consulta. 1 creación/7 replay en carrera;3 intervalos/keys/eventos/actores por proyecto; cancelado no resucita, sin cambios de agenda. No cierre integral, multi-tab, SCA, carga, restore/PITR ni target productivo.",
    "layers": 4,
    "before_utf8_sha256": "76e97b1d00849da2715c7d49cb931bcf76ed6cdde8524ab8ae6146fe74c4ed0f"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      16,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      16,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos, packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "dd79e485bb32ef5fea5df55053ff2915c2e84efdb1b711382a31ecea5c7ed136"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      17,
      "claim"
    ],
    "after": "V308: alta de recursos con actor y consulta del estado/habilidades vigentes. 1 creación/7 replay;3 recursos/keys/eventos/actores y6 habilidades por proyecto; recurso inactivo no se reactiva. No cierre integral, multi-tab, SCA, carga, restore/PITR ni target productivo.",
    "layers": 4,
    "before_utf8_sha256": "5cfe394efac0f0663642c71e719a17119d4d5d8daaa19d5b22f211353e1d86e5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      17,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      17,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos, packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "dd79e485bb32ef5fea5df55053ff2915c2e84efdb1b711382a31ecea5c7ed136"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      18,
      "claim"
    ],
    "after": "Publicación idempotente y recuperación scoped de capacidad con actor atómico; no agenda completa, multi-tab, SCA/carga/restore del target ni readmisión de checklist/retornos.",
    "layers": 4,
    "before_utf8_sha256": "b7fef7801855eb3befa2153325d33575ef435e9327e16a593a664b16b6bdf5bf"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      18,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      18,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos y packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "7c87e367577a87afa50223d31e824e6acfc7f83a3a93f40ad2012784fffc26e8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      19,
      "claim"
    ],
    "after": "Publicación y recuperación exacta de checklist inmutable con actor atómico; no flujo completo de entrega/retornos, multi-tab, SCA/carga/restore del target ni aceptación de negocio.",
    "layers": 4,
    "before_utf8_sha256": "f78b9ce41c372d0285f6ee90650dff019fa605be5e0002673f1275ab3cc97f33"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      19,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      19,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos y packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "7c87e367577a87afa50223d31e824e6acfc7f83a3a93f40ad2012784fffc26e8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      20,
      "claim"
    ],
    "after": "Recuperación exacta de checklist completada con respuestas/actor inmutables y estado actual; no regla de creación de entrega, retornos completos, multi-tab, SCA/carga/restore del target ni aceptación de negocio.",
    "layers": 4,
    "before_utf8_sha256": "e69afe568b4cfc3e8faf0a2f5a979820533c2ba8ccbb6f43482e23ab4a467325"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      20,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1",
    "layers": 4,
    "before_utf8_sha256": "4c51971bf7dd87b04b32a437977b6bff17d2d13f0180e5ffd2ef98f136dd4842"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      20,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos y packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "7c87e367577a87afa50223d31e824e6acfc7f83a3a93f40ad2012784fffc26e8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      21,
      "claim"
    ],
    "after": "Recuperación exacta de recepción/decisión de devoluciones, scope completo y evidencia inmutable; no ejecución downstream, multi-tab, SCA/DAST/carga/restore target ni aceptación de negocio.",
    "layers": 4,
    "before_utf8_sha256": "33d7e8172bb926cc2e6a51dd1809aa8380a15e781bae5f85b5829fbd6c58d50b"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      21,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / Next16.3.2 / Playwright1.62.1 / GO_NATIVE_FUZZ_GATE0.1.0",
    "layers": 4,
    "before_utf8_sha256": "d907dc1b6df033371fd5c6fb23ba71fd0028094a4e474df2ec395fba62444f32"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      21,
      "environment"
    ],
    "after": "Windows loopback, fixtures sintéticos y packs canónicos, sin efectos externos",
    "layers": 4,
    "before_utf8_sha256": "7c87e367577a87afa50223d31e824e6acfc7f83a3a93f40ad2012784fffc26e8"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      22,
      "claim"
    ],
    "after": "Actualización de dependencias fijadas y cobertura de graphs/pins de la composición; no binarios externos, stdlib/runtime, monitoreo persistente, SAST/DAST/carga ni autorización productiva.",
    "layers": 4,
    "before_utf8_sha256": "6bb77909b9bb2ea2c796c092aff32f4c27ba1a3a6bdd3c43d2678e2e3dd79981"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      22,
      "environment"
    ],
    "after": "Windows loopback y fuentes oficiales; sin instalación de wheels, cuentas, mensajes, deploy o scheduled task",
    "layers": 4,
    "before_utf8_sha256": "31514894daaef81ea38e0d334fbd6fc4622169c8472822e88e3d0f9b6032ad5d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      23,
      "claim"
    ],
    "after": "Espera y deadline del drenaje HTTP compartido por ambos hosts; no prueba supervisor, señal OS real, transacción de negocio, race ni operabilidad integral.",
    "layers": 4,
    "before_utf8_sha256": "bac3506bdc6f0a5b5748d98bc9e7094645be642e0fd0a26e9a7cb863a6b203a9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      23,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / net/http TCP / compositor canónico",
    "layers": 4,
    "before_utf8_sha256": "dff438b9c6fd5ddf5b932606b7bb6c91fdd47b31d2058d02795108f6b54556c4"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      24,
      "claim"
    ],
    "after": "Control de referencias y presencia del gate operativo; no contenido, hashes, firmas, frescura ni admisión productiva.",
    "layers": 4,
    "before_utf8_sha256": "6a2fb129f160120ff7ed89ea21c95e0f731bc5dc957551bcac0a8abda0ba9d7c"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      24,
      "tool"
    ],
    "after": "CPython3.14.4 / unittest / CLI / compositor canónico",
    "layers": 4,
    "before_utf8_sha256": "fed38e581ec7e3196029e8a2a4254b1697a6bca6e705165534119ec9a2492791"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      25,
      "claim"
    ],
    "after": "Fencing de estado local del outbox y presupuesto cooperativo del publisher; no efecto remoto exactly-once, supervisor, alertas, retención ni producción.",
    "layers": 4,
    "before_utf8_sha256": "e0545aa0021a3371a0fcecc18093476ecc7d3682399aadbd127a71f4200be28e"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      25,
      "tool"
    ],
    "after": "Go (atribución histórica exacta retirada; ver V317) / PostgreSQL18.6 / compositor canónico",
    "layers": 4,
    "before_utf8_sha256": "561c2d7e8a8dbd2c908a00003eb5bac47d07bb7ce2e4eb92405f2c88b5a818c6"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      25,
      "environment"
    ],
    "after": "Windows loopback, dos bases exclusivas, publisher sintético; sin cuentas ni datos reales",
    "layers": 4,
    "before_utf8_sha256": "3144275175106d898f8a7e541b229ebbea0941e8ba58376044e6d9248fe720bd"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      26,
      "claim"
    ],
    "after": "Identidad Go y cobertura suplementaria stdlib/toolchain; revalidación del head y sincronización de harness. No certifica snapshots históricos, otros runtimes, monitor ni producción.",
    "layers": 4,
    "before_utf8_sha256": "d76cef6c39c9f84cba711e0fc717ec0e833ba977a9768f5f4285709e0097af72"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      26,
      "environment"
    ],
    "after": "Windows loopback con datos sintéticos; sin providers live",
    "layers": 4,
    "before_utf8_sha256": "a3ebe4845de584a58de8b3ff3ddacd6717c8f72c12d438e79473f8cbee91a67d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      27,
      "claim"
    ],
    "after": "Privacidad de salida explícita del host refund; no prueba logs de terceros, delivery/retención/alertas, supervisor o producción.",
    "layers": 4,
    "before_utf8_sha256": "16010becba653eae0a0acaca3c056e59e9ece86fb36042f2e1055b2fe77087e1"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      27,
      "tool"
    ],
    "after": "Go1.26.7 observado / PostgreSQL18.6 / compositor canónico",
    "layers": 4,
    "before_utf8_sha256": "af0badbea20a0cd5f329f411cc9a6ddf378728a06c16eae94f4c7d2e99963f8d"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      27,
      "environment"
    ],
    "after": "Windows loopback, procesos hijos con entorno mínimo y datos sintéticos",
    "layers": 4,
    "before_utf8_sha256": "92790e41833c25f337f99acb36ef0644cab7fb0522f575fc1ea94b1e38e6cbe5"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      28,
      "claim"
    ],
    "after": "Identidad y evaluación acotada de Node24.20.0 independiente; rechazo del ZIP y cobertura falsa. No SCA integral, reusable checker, monitoring ni promoción.",
    "layers": 4,
    "before_utf8_sha256": "88ef6e50f49966739ac05a5d2ad5b5cdda71bcd242269f2db1d57ca69c70a625"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      29,
      "claim"
    ],
    "after": "Adapter acotado del motor oficial Node para instancia Windows exacta y snapshots vigentes; no CLI raw, SCA integral ni monitor/producción.",
    "layers": 4,
    "before_utf8_sha256": "5149fa95bba72580f8686e620f8eeb0e901e367330b0b25e0ce6017610d203f9"
  },
  {
    "file": "PROJECT_ENGINEERING_CONTRACT.json",
    "path": [
      "evidence",
      29,
      "environment"
    ],
    "after": "Windows temporal aislado, fuentes y fixtures públicos, sin providers/datos reales ni instalación global",
    "layers": 4,
    "before_utf8_sha256": "d2dc1c116c9bfa133f03c53ab504ba6fb7ea5b343bd0e2c4ba35ab0903b63481"
  },
  {
    "file": "PROJECT_EXECUTION_STATE.json",
    "path": [
      "open",
      "blockers",
      4
    ],
    "after": "V303 corrige resolver de discrepancias; FAIL457 conserva defectos del sender genérico en otros comandos. Falta recuperación por owner antes de readmisión integral.",
    "layers": 6,
    "before_utf8_sha256": "2f99ebb1932a643a562b2dcb5c15069d2f78596c80ef5ffdc71394075193b742"
  }
]
```

## Regression source

Run with the freshly materialized engineering_execution_kit path as the first argument. `-X utf8=0` is used only for the negative environment regression.

```python
from pathlib import Path
import sys,json,unittest
KIT=Path(sys.argv.pop(1))
sys.path.insert(0,str(KIT))
import test_execution_state as base
from checkpoint_execution_state import append_checkpoint
from validate_execution_state import load_json,validate_event_log,state_sha256

class UnicodeRecoveryTests(base.ExecutionStateTests):
    # Inherit all existing execution-state regressions as well.
    def test_non_ascii_checkpoint_is_exact(self):
        value='Reanudación: aprobación, genérico, €; español — 日本語 😀'
        state=self._base(); state['next_action']=value
        p=self.root/'PROJECT_EXECUTION_STATE.json'
        p.write_text(json.dumps(state,ensure_ascii=False)+'\n',encoding='utf-8',newline='\n')
        raw=p.read_bytes()
        self.assertEqual(load_json(p)['next_action'],value)
        append_checkpoint(p,self.root,self.events,'EVT-0001','BASELINE_CAPTURED','agente',value)
        self.assertEqual(p.read_bytes(),raw)
        event=validate_event_log(self.events)[0]
        self.assertEqual(event['summary'],value)
        self.assertEqual(event['state_sha256'],state_sha256(state))

    def test_windows_wrong_decode_reproduces_and_reverses_exactly(self):
        original='genérico; recuperación; readmisión'
        damaged=original
        for _ in range(6): damaged=damaged.encode('utf-8').decode('cp1252')
        self.assertNotEqual(original,damaged)
        repaired=damaged
        for _ in range(6): repaired=repaired.encode('cp1252').decode('utf-8')
        self.assertEqual(original,repaired)
        self.assertEqual(json.loads(json.dumps({'text':original},ensure_ascii=False).encode('utf-8').decode('utf-8'))['text'],original)

if __name__=='__main__': unittest.main()
```

## Logs

- final-regression-utf8.log: SHA256 9cdf7fb5050e0aacdd63d6e959d3a99de46cc17e72524796e36526b88555a6b2;15tests,exit0.
- final-regression-locale.log: SHA256 9cdf7fb5050e0aacdd63d6e959d3a99de46cc17e72524796e36526b88555a6b2;15tests,exit0.

## Integrated verification and readiness reconciliation

Checkpoint78 appended after recovery and resume validated103evidence/78events, state SHA2563148cc089c0c33dbf736ddffa15c104f3a209ff99ec3546860b8808ed67b83c5. Full VERIFY_LIBRARY then passed161packs/1453files/759Markdown/52profiles; log SHA256 9702350241255ad70151d78ade853be6a5219233ed00e5f9b1dd4c24bf7a7458.

FAIL534 corrected a stale C/PENDING summary against existing C/ANSWERED JSON and V287 answer; removed locally closed FAIL386 from the open list while retaining CANDIDATE operational limits; added existing FAIL532 native evidence gap to dependency blockers. No D–H answers were inferred. The unchanged canonical readiness validator0.6.1 ran on the current root and returned exit2/BLOCKED; its exact stdout bytes replace the stale report after preserving before. The only diagnostic delta is the newly represented native blocker. This is corrected traceability, not improved readiness or task completion.

```json
{
  "before_errors": 41,
  "after_errors": 42,
  "validator_exit": 2,
  "added_errors": [
    "blockers.open_dependency_blockers must be empty"
  ],
  "files": [
    {
      "path": "PROJECT_READINESS_RECORD.md",
      "before_sha256": "f9226b5188441a36c105f601deb6229e83cf36696da36765e70f676e61a99e44",
      "after_sha256": "6f85c20d342b057e101350c8d1e971e2d128c5aecb4b651ddd4533d722f715af"
    },
    {
      "path": "PROJECT_READINESS_GATE.json",
      "before_sha256": "dd700ebad49c2718ebc1fdfae6efb301441e39092f84f53bcfb6fb392a4f234e",
      "after_sha256": "090653052db29e24045f5421321e5b351d0f9e79471c609b5d04796fbedb14a3"
    },
    {
      "path": "PROJECT_READINESS_REPORT.json",
      "before_sha256": "0151bf7341b47284e9bd02c762dad3ee46f146f659dc0e595ca0aba88aed14c8",
      "after_sha256": "0f4492023826cbb83ec0ff98c5e2f9c95b361cc51f401675ef6bb93d6dd6405b"
    }
  ]
}
```
