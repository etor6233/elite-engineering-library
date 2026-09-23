import { expect, test } from "vitest";
import { createElement } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { createHash } from "node:crypto";
import {OperationalGuide as V402Guide} from "./fixtures/operational-guide.v402";
import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "./content";

// Historical V379 digests retained as provenance only. V402 had already added i18n/lang; do not call those old HTML bytes the current semantic baseline.
const original = [
  {
    "id": "return-operations-view",
    "sha256": "890335861532191b4fb532c11a0c1c07998eb388853cf767b35094411f0b7119",
    "className": null
  },
  {
    "id": "checklist-completion-view",
    "sha256": "2eb9cc5986181d4016397cef06f8373ce6dc9820543f8ecdb6e7a2aa681a9cb5",
    "className": null
  },
  {
    "id": "checklist-publication-view",
    "sha256": "a7487655d312417f4514358491fe0ef3d8b2972cb1113daa3cfaacf0af3f45f6",
    "className": null
  },
  {
    "id": "slot-create-view",
    "sha256": "db23a5340f56125be6842fa6d885995ad9bca2fae0b5c02191aa692711e43265",
    "className": null
  },
  {
    "id": "resource-create-view",
    "sha256": "2adbb54d47ad45c863bbf139fa81d39eb4780ba203a40e712dc00418252f4618",
    "className": null
  },
  {
    "id": "availability-create-view",
    "sha256": "3ea153fa5cc8cc8a9e36bad06ca38fa3581c62b7be72247fc826b0ade24041b8",
    "className": null
  },
  {
    "id": "quote-create-view",
    "sha256": "2361437eaac0a930954e0852779dfc1ac0cc60d4f7ad336a9fa48fcce7cbf93d",
    "className": null
  },
  {
    "id": "delivery-resolution-view",
    "sha256": "28126f161b2c40ab03836c14edefcd2530d9107bc2aa02d7b6f768bfcc85600a",
    "className": null
  },
  {
    "id": "lead-command-view",
    "sha256": "571367bf03b0d3ff229743a2ffde37c044edf4aabdd8bd05c9f927a74eec1dc8",
    "className": null
  },
  {
    "id": "availability-cancel-view",
    "sha256": "7b18e926cbcd8ab34b9a06fefdd5fa99e9162f964a3c22a9434d499e4bcbbd24",
    "className": null
  },
  {
    "id": "order-operations-view",
    "sha256": "d8954474a5b5f6d2f983ec4e755b5496f8617fea16cddc1dada7b27ce060b401",
    "className": null
  },
  {
    "id": "operation-sections-view",
    "sha256": "29fd0212dffe93893892c469c07a6a2a557b6a938d0d7fa2ce9569ddd241932b",
    "className": "card"
  },
  {
    "id": "handover-read-view",
    "sha256": "bcd57c885618d24d15922ee66ba495b2943c2b7aad506c5c7884b14d8a12737e",
    "className": null
  }
] as const;
test.each(original)("inline $id preserves exact V402 default rendering, paragraphs, version and link", row => {
 const html=renderToStaticMarkup(createElement(OperationalGuide,{guide:OPERATIONAL_GUIDES[row.id],...(row.className?{className:row.className}:{})}));
 const expectedLink=`<a href="/help?article=${row.id}&amp;version=${OPERATIONAL_GUIDES[row.id].version}">Abrir esta guía en Ayuda</a>`;
 expect(html).toContain(expectedLink);
 const baseline=renderToStaticMarkup(createElement(V402Guide,{guide:OPERATIONAL_GUIDES[row.id],...(row.className?{className:row.className}:{})}));
 expect(html).toBe(baseline);
});
