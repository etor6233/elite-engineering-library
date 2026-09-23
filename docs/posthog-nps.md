# PostHog NPS: bounded Go adaptation

This owner translates only calculateNPSFromRawData from PostHog's MIT core,
commit 6fafbb9081bd15e79448af5650e02a4f9ea435cc. It does not import PostHog's
product, enterprise code, dependencies, analytics runtime, or survey service.
Original utils.ts, utils.test.ts and root LICENSE are retained byte-for-byte
as non-executable evidence under third_party/posthog-nps. See source-lock.json.

The original sums 0..6 detractors, 7..8 passives, 9..10 promoters, then computes
(promoters - detractors) / total * 100. Existing SQL groups these same ranges;
the adapted Go owner takes validated grouped counts and derives passives from
total. It preserves zero-population behavior. int64 counts replace JS Number;
subtraction happens before float conversion, preserving the existing Go API.
The original toFixed(1) produces a display string. This adapter deliberately
returns unrounded float64 because changing the public API to one decimal would
be a compatibility change. Six meaningful typed equivalents of the original
nine cases are included; all nine original TS cases were separately executed.

The AUTHORED customerfeedback caller retains overflow-safe validation, explicit
project-owned minimum-response threshold, consent/configuration and access
boundaries. It rejects inconsistent/negative counts before calling Calculate;
the upstream legacy parser is not used. No sampling, legal consent, business
policy, or whole-PostHog certification is implied by this numerical adaptation.
Caller glue, fuzz oracle and source-lock metadata are local work. The algorithm
and six translated examples are ADAPTED MIT with the original notice retained.
No new runtime or module is required. Fixed official provenance is mandatory
before adoption; local fixture qualification does not certify production.
