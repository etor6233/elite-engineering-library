import test from 'node:test';
import assert from 'node:assert/strict';
import {validatePolicy, validateReports, validateTargetUrl} from '../lib/validate-reports.mjs';

const policy = {
  schemaVersion: '1.0',
  expectedLighthouseVersion: '13.4.1',
  numberOfRuns: 5,
  maxRunWarnings: 0,
  minimumAuditCount: 2,
  categories: {
    performance: {median: 0.90, minimum: 0.85},
    accessibility: {median: 1, minimum: 1},
    'best-practices': {median: 0.90, minimum: 0.90},
    seo: {median: 1, minimum: 1}
  }
};

function report(performance = 0.95) {
  return {
    lighthouseVersion: '13.4.1',
    requestedUrl: 'http://127.0.0.1:4181/',
    finalUrl: 'http://127.0.0.1:4181/',
    runWarnings: [],
    audits: {one: {}, two: {}},
    categories: {
      performance: {score: performance},
      accessibility: {score: 1},
      'best-practices': {score: 0.96},
      seo: {score: 1}
    }
  };
}

test('accepts five governed reports and exposes explicit limits', () => {
  const summary = validateReports([report(0.91), report(0.92), report(0.95), report(0.97), report(0.99)], policy, 'http://127.0.0.1:4181/');
  assert.equal(summary.status, 'PASS');
  assert.equal(summary.categories.performance.median, 0.95);
  assert.equal(summary.limits.provesHumanAssistiveTechnology, false);
});

test('rejects credentials, remote HTTP and fragments', () => {
  assert.throws(() => validateTargetUrl('https://user:secret@example.com/'));
  assert.throws(() => validateTargetUrl('http://example.com/'));
  assert.throws(() => validateTargetUrl('https://example.com/#private'));
});

test('rejects missing categories, even run counts and version ranges', () => {
  assert.throws(() => validatePolicy({...policy, numberOfRuns: 4}));
  assert.throws(() => validatePolicy({...policy, expectedLighthouseVersion: '^13.4.1'}));
  const incomplete = structuredClone(policy);
  delete incomplete.categories.seo;
  assert.throws(() => validatePolicy(incomplete));
});

test('rejects warnings, origin escapes, missing audits and version drift', () => {
  const warnings = report(); warnings.runWarnings = ['warning'];
  assert.throws(() => validateReports([warnings, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  const escape = report(); escape.finalUrl = 'https://example.com/';
  assert.throws(() => validateReports([escape, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  const few = report(); few.audits = {one: {}};
  assert.throws(() => validateReports([few, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  const drift = report(); drift.lighthouseVersion = '13.4.2';
  assert.throws(() => validateReports([drift, report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
});

test('rejects a weak single run and a weak median', () => {
  assert.throws(() => validateReports([report(0.84), report(), report(), report(), report()], policy, 'http://127.0.0.1:4181/'));
  assert.throws(() => validateReports([report(0.86), report(0.87), report(0.88), report(0.96), report(0.98)], policy, 'http://127.0.0.1:4181/'));
});
