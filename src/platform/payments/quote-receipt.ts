// AUTHORED validation glue for the existing Go Quote receipt; no commercial rule.
export function isQuoteAcceptanceReceipt(value: unknown, organizationId: string, quote: {id:string;version:number}): boolean {
  if (!value || typeof value !== "object" || Array.isArray(value)) return false;
  const receipt=value as Record<string,unknown>;
  return receipt.id===quote.id && receipt.organization_id===organizationId && receipt.state==="accepted" &&
    Number.isSafeInteger(quote.version) && quote.version>0 && Number.isSafeInteger(quote.version+1) &&
    receipt.version===quote.version+1 && typeof receipt.order_id==="string" && receipt.order_id.length>0 && receipt.order_id.length<=128;
}
