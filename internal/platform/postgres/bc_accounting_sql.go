// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// Adapted from microsoft/BCApps at 2eae56d704a1fd035d104f333602aea7091b7749,
// GenJnlPostReverse.Codeunit.al:209-257. Local PostgreSQL representation delta.
// The source negates signed G/L amounts; this boundary maps the result to the
// library's existing nonnegative debit/credit columns. Original entry links,
// journal state, tenant authorization and transaction remain the caller's owner.
// Negative source correction columns and turnover are not claimed equivalent.
package postgres

const reverseJournalLinesBCSQL = `insert into accounting.journal_line(tenant_id,journal_id,line_no,account_code,description,debit_minor_units,credit_minor_units)
select tenant_id,$3,line_no,account_code,'Reversal: '||description,
       greatest(credit_minor_units-debit_minor_units,0),
       greatest(debit_minor_units-credit_minor_units,0)
from accounting.journal_line where tenant_id=$1 and journal_id=$2`
