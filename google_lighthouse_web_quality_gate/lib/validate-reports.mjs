import {readFile} from 'node:fs/promises';

function fail(message) {
  throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: ${message}`);
}

function median(values) {
  const ordered = [...values].sort((a, b) => a - b);
  const middle = Math.floor(ordered.length / 2);
  return ordered.length % 2 ? ordered[middle] : (ordered[middle - 1] + ordered[middle]) / 2;
}

export function validatePolicy(policy) {
  if (policy?.schemaVersion !== '1.0') fail('unsupported policy schema');
  if (!/^\d+\.\d+\.\d+$/.test(policy.expectedLighthouseVersion ?? '')) fail('expectedLighthouseVersion must be exact');
  if (!Number.isInteger(policy.numberOfRuns) || policy.numberOfRuns < 3 || policy.numberOfRuns > 9 || policy.numberOfRuns % 2 === 0) {
    fail('numberOfRuns must be an odd integer from 3 through 9');
  }
  if (!Number.isInteger(policy.maxRunWarnings) || policy.maxRunWarnings < 0) fail('maxRunWarnings must be a non-negative integer');
  if (!Number.isInteger(policy.minimumAuditCount) || policy.minimumAuditCount < 1) fail('minimumAuditCount must be positive');
  const names = Object.keys(policy.categories ?? {});
  if (names.join(',') !== 'performance,accessibility,best-practices,seo') fail('all four categories are required in canonical order');
  for (const name of names) {
    const rule = policy.categories[name];
    for (const key of ['median', 'minimum']) {
      if (!Number.isFinite(rule?.[key]) || rule[key] < 0 || rule[key] > 1) fail(`${name}.${key} must be between 0 and 1`);
    }
    if (rule.minimum > rule.median) fail(`${name}.minimum cannot exceed median`);
  }
  return policy;
}

export function validateTargetUrl(rawUrl) {
  let url;
  try { url = new URL(rawUrl); } catch { fail('target URL is invalid'); }
  if (url.username || url.password) fail('credentials are forbidden in target URLs');
  const loopback = url.hostname === '127.0.0.1' || url.hostname === 'localhost' || url.hostname === '::1';
  if (url.protocol !== 'https:' && !(url.protocol === 'http:' && loopback)) fail('target must use HTTPS or loopback HTTP');
  if (url.hash) fail('URL fragments are not accepted as audit targets');
  return url;
}

export function validateReports(reports, policy, expectedUrl) {
  validatePolicy(policy);
  const target = validateTargetUrl(expectedUrl);
  if (!Array.isArray(reports) || reports.length !== policy.numberOfRuns) fail(`expected ${policy.numberOfRuns} reports`);
  const scoreSets = Object.fromEntries(Object.keys(policy.categories).map(name => [name, []]));
  const runs = [];
  for (const [index, report] of reports.entries()) {
    if (report?.lighthouseVersion !== policy.expectedLighthouseVersion) fail(`run ${index + 1} Lighthouse version drifted`);
    if (!Array.isArray(report.runWarnings) || report.runWarnings.length > policy.maxRunWarnings) fail(`run ${index + 1} has disallowed warnings`);
    const requested = validateTargetUrl(report.requestedUrl);
    const final = validateTargetUrl(report.finalUrl);
    if (requested.origin !== target.origin || final.origin !== target.origin) fail(`run ${index + 1} escaped the target origin`);
    const auditCount = Object.keys(report.audits ?? {}).length;
    if (auditCount < policy.minimumAuditCount) fail(`run ${index + 1} audit count ${auditCount} is below ${policy.minimumAuditCount}`);
    const scores = {};
    for (const name of Object.keys(policy.categories)) {
      const score = report.categories?.[name]?.score;
      if (!Number.isFinite(score) || score < 0 || score > 1) fail(`run ${index + 1} category ${name} has no numeric score`);
      scoreSets[name].push(score);
      scores[name] = score;
    }
    runs.push({run: index + 1, auditCount, warnings: report.runWarnings.length, scores});
  }
  const categories = {};
  for (const [name, values] of Object.entries(scoreSets)) {
    const actualMedian = median(values);
    const actualMinimum = Math.min(...values);
    const rule = policy.categories[name];
    if (actualMedian < rule.median) fail(`${name} median ${actualMedian} is below ${rule.median}`);
    if (actualMinimum < rule.minimum) fail(`${name} minimum ${actualMinimum} is below ${rule.minimum}`);
    categories[name] = {median: actualMedian, minimum: actualMinimum, values};
  }
  return {
    schemaVersion: '1.0',
    status: 'PASS',
    targetOrigin: target.origin,
    lighthouseVersion: policy.expectedLighthouseVersion,
    numberOfRuns: reports.length,
    categories,
    runs,
    limits: {
      automatedLabOnly: true,
      provesHumanAssistiveTechnology: false,
      provesProductionPerformance: false,
      provesLoadOrSecurity: false
    }
  };
}

export async function loadJson(path) {
  return JSON.parse(await readFile(path, 'utf8'));
}
