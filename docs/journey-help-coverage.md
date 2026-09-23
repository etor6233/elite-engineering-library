# Versioned help coverage — V380

Scope:15/15versioned operational guides already present in the V379 selected web sources.
Thirteen were extracted without changing any instruction/version; two were already shared.
This is coverage of existing guidance, not a claim that every possible user task has a guide.
In particular factory:read alone has no existing guide; it receives403, as do empty permissions
and quote:write without lead:read. No factory policy, CMS or training curriculum is invented.

GET resolves the existing JWE session and applies the page/section predicates below. Help is
generic PUBLIC content, not a command grant or tenant-confidential document store. Lead readers
may read lead guidance without acquiring assign/update permissions. Availability readers may
read cancellation guidance without acquiring manage; manage alone can read creation guidance.
Quote creation guidance also needs lead:read because the actual page loads its leads under that
permission. The notification feature flag gates only the WhatsApp article, not other agenda help.

| Session permissions | Available article IDs |
|---|---|
| (none) | (none;403) |
| factory:read | (none;403) |
| customer:self | handover-read-view, quote-acceptance-view |
| lead:read | lead-command-view, operation-sections-view |
| quote:write | (none;403) |
| lead:read, quote:write | lead-command-view, operation-sections-view, quote-create-view |
| availability:read | availability-cancel-view, operation-sections-view |
| availability:manage | availability-create-view, operation-sections-view |
| resource:manage | operation-sections-view, resource-create-view |
| appointment:manage | operation-sections-view, slot-create-view, whatsapp-status-view |
| inventory:allocate | operation-sections-view, order-operations-view |
| payment:create | operation-sections-view, order-operations-view |
| handover:manage | checklist-completion-view, checklist-publication-view, delivery-resolution-view, operation-sections-view, order-operations-view, return-operations-view |
| admin:read | operation-sections-view, order-operations-view |
| * | availability-cancel-view, availability-create-view, checklist-completion-view, checklist-publication-view, delivery-resolution-view, handover-read-view, lead-command-view, operation-sections-view, order-operations-view, quote-acceptance-view, quote-create-view, resource-create-view, return-operations-view, slot-create-view, whatsapp-status-view |

## Exact source extraction

| Article / version | Previous inline source | Original JSX SHA-256 |
|---|---|---|
| availability-cancel-view/1.0.0 | src/components/franchise-command-panel.tsx | 7b18e926cbcd8ab34b9a06fefdd5fa99e9162f964a3c22a9434d499e4bcbbd24 |
| availability-create-view/1.0.0 | src/components/franchise-command-panel.tsx | 3ea153fa5cc8cc8a9e36bad06ca38fa3581c62b7be72247fc826b0ade24041b8 |
| checklist-completion-view/1.0.0 | src/components/franchise-command-panel.tsx | 2eb9cc5986181d4016397cef06f8373ce6dc9820543f8ecdb6e7a2aa681a9cb5 |
| checklist-publication-view/1.0.0 | src/components/franchise-command-panel.tsx | a7487655d312417f4514358491fe0ef3d8b2972cb1113daa3cfaacf0af3f45f6 |
| delivery-resolution-view/1.0.0 | src/components/franchise-command-panel.tsx | 28126f161b2c40ab03836c14edefcd2530d9107bc2aa02d7b6f768bfcc85600a |
| handover-read-view/1.0.0 | src/app/customer/handovers/page.tsx | bcd57c885618d24d15922ee66ba495b2943c2b7aad506c5c7884b14d8a12737e |
| lead-command-view/1.0.0 | src/components/franchise-command-panel.tsx | 571367bf03b0d3ff229743a2ffde37c044edf4aabdd8bd05c9f927a74eec1dc8 |
| operation-sections-view/1.0.0 | src/app/franchise/page.tsx | f4ee283d9423bef7179d8535271d4ba970331409bd5165ddd232e64aceeadf5c |
| order-operations-view/1.1.0 | src/components/franchise-command-panel.tsx | d8954474a5b5f6d2f983ec4e755b5496f8617fea16cddc1dada7b27ce060b401 |
| quote-create-view/1.0.0 | src/components/franchise-command-panel.tsx | 2361437eaac0a930954e0852779dfc1ac0cc60d4f7ad336a9fa48fcce7cbf93d |
| resource-create-view/1.0.0 | src/components/franchise-command-panel.tsx | 2adbb54d47ad45c863bbf139fa81d39eb4780ba203a40e712dc00418252f4618 |
| return-operations-view/1.0.0 | src/components/franchise-command-panel.tsx | 890335861532191b4fb532c11a0c1c07998eb388853cf767b35094411f0b7119 |
| slot-create-view/1.0.0 | src/components/franchise-command-panel.tsx | db23a5340f56125be6842fa6d885995ad9bca2fae0b5c02191aa692711e43265 |

The13render regression oracles compare the full original inline HTML, including summary,
paragraph order, version, punctuation and CSS class, after excluding only the newly added help
link. JSX className becomes HTML class in the one classed details block; its original source
digest above is deliberately retained as history. A separate frozen HTML digest governs that
render test. Semantic guide changes require a new guide version and an explicit before/after
review; do not refresh regression digests merely to conceal an instruction change.

The real HTTPS Next/JWE browser suite checks15permission profiles times15exact article links
in each of four browsers (900direct read decisions), plus all15UI version links in each browser
(60navigations), recovery and absence of browser writes. These figures count checks, not global
controls. Original business handlers outside the extracted help blocks remain byte-identical.
No external IdP, database, provider, confidential data or training activity is used by this gate.
