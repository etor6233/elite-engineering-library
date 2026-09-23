import {mkdir, mkdtemp, rm, writeFile} from 'node:fs/promises';
import {tmpdir} from 'node:os';
import {resolve, join} from 'node:path';
import process from 'node:process';
import lighthouse from 'lighthouse';
import {launch} from 'chrome-launcher';
import {loadJson, validatePolicy, validateReports, validateTargetUrl} from './lib/validate-reports.mjs';

function argument(name) {
  const prefix = `--${name}=`;
  const values = process.argv.slice(2).filter(value => value.startsWith(prefix));
  if (values.length !== 1 || values[0].length === prefix.length) throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: exactly one ${prefix}<value> is required`);
  return values[0].slice(prefix.length);
}

const targetUrl = validateTargetUrl(argument('url')).href;
const chromePath = resolve(argument('chrome-path'));
const outputDirectory = resolve(argument('output-dir'));
const policyPath = resolve(argument('policy'));
const policy = validatePolicy(await loadJson(policyPath));
await mkdir(outputDirectory, {recursive: true});
const profile = await mkdtemp(join(tmpdir(), 'elite-lighthouse-profile-'));
let chrome;
let primaryError;
try {
  chrome = await launch({
    chromePath,
    userDataDir: profile,
    chromeFlags: ['--headless', '--disable-gpu', '--no-sandbox'],
    logLevel: 'silent'
  });
  const reports = [];
  for (let index = 0; index < policy.numberOfRuns; index += 1) {
    const result = await lighthouse(targetUrl, {
      port: chrome.port,
      output: 'json',
      logLevel: 'silent',
      onlyCategories: Object.keys(policy.categories)
    });
    if (!result?.lhr || typeof result.report !== 'string') throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: run ${index + 1} returned no JSON report`);
    reports.push(result.lhr);
    await writeFile(join(outputDirectory, `lighthouse-run-${index + 1}.json`), result.report, {encoding: 'utf8', flag: 'wx'});
  }
  const summary = validateReports(reports, policy, targetUrl);
  await writeFile(join(outputDirectory, 'quality-gate-summary.json'), `${JSON.stringify(summary, null, 2)}\n`, {encoding: 'utf8', flag: 'wx'});
  process.stdout.write(`GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE_PASS runs=${summary.numberOfRuns} performance_median=${summary.categories.performance.median} accessibility_median=${summary.categories.accessibility.median} best_practices_median=${summary.categories['best-practices'].median} seo_median=${summary.categories.seo.median}\n`);
} catch (error) {
  primaryError = error;
  throw error;
} finally {
  let cleanupError;
  try { if (chrome) await chrome.kill(); } catch (error) { cleanupError = error; }
  try { await rm(profile, {recursive: true, force: true, maxRetries: 20, retryDelay: 250}); } catch (error) { cleanupError ??= error; }
  if (!primaryError && cleanupError) throw new Error(`LIGHTHOUSE_QUALITY_GATE_FAILED: browser cleanup failed: ${cleanupError.message}`);
}
