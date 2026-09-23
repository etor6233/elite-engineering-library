// AUTHORED test fixture using Node's existing crypto/http standard library.
// Synthetic loopback identity only; never a deployable identity provider.
'use strict';
const http = require('node:http');
const crypto = require('node:crypto');
const args = process.argv.slice(2);
const ttl = args.length === 0 ? 600 : args.length === 2 && args[0] === '--ttl-seconds' && /^(?:[1-9][0-9]{0,3})$/.test(args[1]) ? Number(args[1]) : 0;
if (ttl < 1 || ttl > 1200) throw new Error('bounded fixture lifetime required');
const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', { modulusLength: 2048 });
const key = { ...publicKey.export({ format: 'jwk' }), kid: 'local-reference', use: 'sig', alg: 'RS256' };
let issuer;
const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json');
  if (req.method !== 'GET') { res.writeHead(405).end('{}'); return; }
  if (req.url === '/.well-known/openid-configuration') {
    res.end(JSON.stringify({ issuer, jwks_uri: issuer + '/jwks', authorization_endpoint: issuer + '/authorize',
      token_endpoint: issuer + '/token', response_types_supported: ['code'], subject_types_supported: ['public'],
      id_token_signing_alg_values_supported: ['RS256'] }));
  } else if (req.url === '/jwks') res.end(JSON.stringify({ keys: [key] }));
  else if (req.url === '/fixture-token') {
    const now = Math.floor(Date.now() / 1000);
    const claims = { iss: issuer, aud: 'elite-local-reference', sub: 'local-delivery-customer', iat: now, exp: now + 120,
      tenant_id: 'b0000000-0000-4000-8000-000000000001', organization_ids: ['b0000000-0000-4000-8000-000000000002'],
      permissions: ['order:create'] };
    const part = obj => Buffer.from(JSON.stringify(obj)).toString('base64url');
    const input = part({ alg: 'RS256', kid: key.kid, typ: 'JWT' }) + '.' + part(claims);
    res.end(JSON.stringify({ token: input + '.' + crypto.sign('RSA-SHA256', Buffer.from(input), privateKey).toString('base64url') }));
  } else res.writeHead(404).end('{}');
});
server.listen(0, '127.0.0.1', () => {
  issuer = 'http://127.0.0.1:' + server.address().port;
  process.stdout.write(JSON.stringify({ issuer }) + '\n');
});
process.stdin.on('data', () => server.close(() => process.exit(0)));
process.stdin.on('end', () => server.close(() => process.exit(0)));
setTimeout(() => server.close(() => process.exit(0)), ttl * 1000).unref();
