// AUTHORED local trusted-stdin shutdown bridge to Next's existing SIGTERM owner.
// The Windows Job Object still bounds the whole descendant tree on failure.
'use strict';
const path = require('node:path');
const fs = require('node:fs');
if (process.argv.length !== 3 || !path.isAbsolute(process.argv[2])) throw new Error('absolute server entry required');
const root = path.dirname(process.argv[2]);
const links = JSON.parse(fs.readFileSync(path.join(root, 'elite-runtime-links.json'), 'utf8'));
function local(rel) {
  if (typeof rel !== 'string' || !rel.startsWith('node_modules/') || rel.includes('\\') || rel.includes(':') || rel.split('/').some(p => !p || p === '.' || p === '..')) throw new Error('invalid runtime link');
  return path.join(root, ...rel.split('/'));
}
for (const [name, dest] of Object.entries(links)) {
  const link = local(name), target = local(dest);
  if (fs.existsSync(link) || !fs.statSync(target).isDirectory()) throw new Error('occupied link or absent target');
  // ZIPs preserve files, not the empty scope directories that precede links.
  fs.mkdirSync(path.dirname(link), { recursive: true });
  fs.symlinkSync(target, link, 'junction');
}
require(process.argv[2]);
let requested = false;
function stop() {
  if (requested) return;
  requested = true;
  process.emit('SIGTERM', 'SIGTERM');
}
let input = '';
process.stdin.setEncoding('utf8');
process.stdin.on('data', chunk => {
  input += chunk;
  if (input.length > 16) throw new Error('invalid local shutdown command');
  if (input === 'STOP\n') stop();
});
process.stdin.on('end', stop);
