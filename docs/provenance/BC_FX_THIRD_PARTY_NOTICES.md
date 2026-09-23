# BC FX adapted source notice

Copyright (c) Microsoft Corporation.

The two arithmetic/selection files in internal/bcfx are adapted from the fixed BCApps source documented in BC_FX_DERIVATION.md and licensed under the exact MIT notice in licenses/Microsoft-BCApps-MIT.txt. Local integration files are AUTHORED and are not Microsoft code. No AL runtime or upstream suite execution is claimed.

Source correction FX-DOC-ROUND-001: the official System.Round page inspected for this candidate (https://learn.microsoft.com/en-us/dynamics365/business-central/dev-itpro/developer/methods-auto/system/system-round-method, last-updated2025-01-28) has an inconsistent negative example: -1234.56789 with precision1 is printed as -1234, while its nearest rule and other negative examples imply -1235. The root candidate comment asserting AL built-in equality was corrected before publication. This profile explicitly selects nearest/ties-away-from-zero and verifies it with an independent Python Decimal oracle; it makes no AL runtime-equivalence claim for that contradiction. No arbitrary exclusion of negative amounts is introduced.
