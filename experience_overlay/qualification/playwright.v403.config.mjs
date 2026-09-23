import {target,evidence} from "./paths.mjs";
import {createRequire}from'node:module';
import {resolve}from'node:path';
const require=createRequire(process.env.PLAYWRIGHT_TEST_PACKAGE||resolve(target,'microsoft_playwright_browser_gate/package.json'));
const {defineConfig}=require('@playwright/test');
export default defineConfig({testDir:'./tests',workers:1,retries:0,timeout:45000,forbidOnly:true,reporter:[['list'],['json',{outputFile:process.env.UI_BROWSER_REPORT||resolve(evidence,'browser-tests.json')}]],outputDir:resolve(evidence,'browser-output'),snapshotPathTemplate:'{testDir}/visual-candidates/{arg}{ext}',use:{baseURL:'https://127.0.0.1:4313',ignoreHTTPSErrors:true,locale:'es-AR',browserName:'chromium',viewport:{width:1440,height:1000},trace:'retain-on-failure'},expect:{toHaveScreenshot:{maxDiffPixels:0,animations:'disabled'}}});
