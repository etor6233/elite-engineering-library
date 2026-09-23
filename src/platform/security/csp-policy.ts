const noncePattern = /^[A-Za-z0-9+/]+={0,2}$/;

export function buildContentSecurityPolicy(nonce: string, development: boolean): string {
  if (nonce.length < 16 || !noncePattern.test(nonce)) throw new Error("CSP nonce must be base64 and at least 16 characters");
  return `
    default-src 'self';
    script-src 'self' 'nonce-${nonce}' 'strict-dynamic'${development ? " 'unsafe-eval'" : ""};
    style-src 'self' 'nonce-${nonce}';
    img-src 'self' blob: data:;
    font-src 'self';
    connect-src 'self';
    object-src 'none';
    base-uri 'self';
    form-action 'self';
    frame-ancestors 'none';
    upgrade-insecure-requests;
  `.replace(/\s{2,}/g, " ").trim();
}
