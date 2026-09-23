#!/usr/bin/env python3
"""Prepare/verify a pinned offline-install recipe. Never execute or admit pnpm."""
from __future__ import annotations
import argparse
import base64
import binascii
import os
from pathlib import Path
import shutil
import sys
import tempfile
from select_artifact import (SelectionError, require, digest, json_bytes, decode,
    no_reparse, read_bounded, load_policy, acquisition_binding, verify, safe_name)

NODE_SHA256 = '5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5'
PROFILES = {'enterprise-web': {'package.json': 'acc4e4bd62c13cf48e071e933266d82d067cab849f92474d71a533ebf6508951', 'pnpm-lock.yaml': 'd6c73eb82a56c7e8d21ee89ffb402d0adb2f0c15db9c33bb8ec7035775119928', 'pnpm-workspace.yaml': 'f52887318fe556dbdb3574622f5c6248483e843edc3b61ad98b6ae3b351571ba'}, 'playwright': {'package.json': 'eabb3e705b59759748724b2be3e8720184d26b570ab51e6c8a2ee8a0bdf41b90', 'pnpm-lock.yaml': '63eec1e3d5bac29965e850bf751269c7a8b317929b570205a1e223544f80a470'}, 'lighthouse': {'package.json': 'dcdb7fd7c346f73543643070f14c15578a3f425245b24a5398b314d26366b7bb', 'pnpm-lock.yaml': 'c690da3a314e5275d68c4ea4a88e700dbd9a58b0da6fcbc25aec2bdc486dd9aa'}}
RECIPE = 'install-plan.json'
NOTICE = 'BLUEOAK-NOTICE.md'
QRCODE_NOTICE = 'QRCODE-NOTICE.md'

SEMVER_NOTICE = 'SEMVER-UTILS-NOTICE.md'

RETAINED_NOTICE = 'PNPM-RETAINED-NOTICES.md'

NEXT_PATH_NOTICE = 'NEXT-PATH-MPL-SOURCE.md'
NEXT_PATH_DATA_SHA256 = '581c15a6e63ecc37d3c4f4fde4e701ed25532ffc141c519dd88919baf8685dd7'
NEXT_PATH_REGION = {'start': 7619794, 'end': 7620528, 'sha256': 'd8b58764b8a66de031a13ad8b76cccbb936e49f69c7df6d6bda306ced71fab83'}

CATALOG_SHA256 = 'e4a906897aeae3dec44d63beb40514329353f41859c8d380c126df8b2d599c1e'

SEMVER_LICENSE = 'Copyright 2013 AJ ONeal\n\nThis is open source software; you can redistribute it and/or modify it under the\nterms of either:\n\n   a) the "MIT License"\n   b) the "Apache-2.0 License"\n\nMIT License\n\n   Permission is hereby granted, free of charge, to any person obtaining a copy\n   of this software and associated documentation files (the "Software"), to deal\n   in the Software without restriction, including without limitation the rights\n   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n   copies of the Software, and to permit persons to whom the Software is\n   furnished to do so, subject to the following conditions:\n\n   The above copyright notice and this permission notice shall be included in all\n   copies or substantial portions of the Software.\n\n   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\n   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\n   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\n   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\n   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\n   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\n   SOFTWARE.\n\nApache-2.0 License Summary\n\n   Licensed under the Apache License, Version 2.0 (the "License");\n   you may not use this file except in compliance with the License.\n   You may obtain a copy of the License at\n\n     http://www.apache.org/licenses/LICENSE-2.0\n\n   Unless required by applicable law or agreed to in writing, software\n   distributed under the License is distributed on an "AS IS" BASIS,\n   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n   See the License for the specific language governing permissions and\n   limitations under the License.\n'
SEMVER_LICENSE_SHA256 = '80b98c1b20edfc51abd2c802ee7a1d3c5561151d108a24f19c207b6935beaed8'
SEMVER_REGION = {'start': 1347913, 'end': 1350984, 'region_sha256': 'aed20138b14658d55690a793c31b50d4f97e3c0d5260e323b7a74fb8279d9cf2', 'marker': '// ../../../../../setup-pnpm/store/v11/links/@/semver-utils/1.1.4/bbbf6bcd6e2cad096d8c5e4cd94c8815535d76e88f7c54b6d35929d441dd7625/node_modules/semver-utils/semver-utils.js', 'source_member_sha256': '1f85bceb6f2e6cdbefadd5bf9969db8637c754c9d6553f3d7663dca6ade7c293'}

QRCODE_BUNDLE = 'dist/pnpm.mjs'
QRCODE_BUNDLE_SHA256 = 'ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11'
QRCODE_MODULES = [{'path': 'vendor/QRCode/QRMode.js', 'source_sha256': '6b8ec04257a2d23b01e8189815292dca3651b38ea0a8f9c975b3c1d18dfb1b01', 'source_git_blob': '050c8a3037f60d1e034764035394c55bca34f3c3', 'start': 2065080, 'end': 2065663, 'region_sha256': '008dd02912868c88714e36bce44f2b472b7812d417183053b7fb3ae47d7126e5'}, {'path': 'vendor/QRCode/QR8bitByte.js', 'source_sha256': 'a67f0b2239db81b1fc1dfd8e169a879d7075dd79d0ae00dc155e9c3bac595891', 'source_git_blob': '94bf74f0e897de3d7c432d0d6f6422dcee8e471a', 'start': 2065664, 'end': 2066555, 'region_sha256': '3ec44883fe4fab2eb98f4f4ae841b836d19fb7fa47523fb4f019bbca23a77e60'}, {'path': 'vendor/QRCode/QRMath.js', 'source_sha256': '481fe65cd1a049a3cdd659ff20c45eb4e0cb2db285fa63a42478727e1b051667', 'source_git_blob': '8f4a0370ebb32a2e93ae053fb475732215a2cd02', 'start': 2066556, 'end': 2067839, 'region_sha256': 'b617520021a3eb20389c9483ab7fb517eaa5c5720e1058e60c5841302783a997'}, {'path': 'vendor/QRCode/QRPolynomial.js', 'source_sha256': '76eb786a451ceee003cb4279b7bc559e8a77321dad19ce11825a4d98d470b422', 'source_git_blob': '0c05f38ef324680c1cf53e57159567639177d546', 'start': 2067840, 'end': 2069843, 'region_sha256': '7789cf38de59770fe3377f63d3449aefc08ecea0ba9d3c65e47b037f5918a75c'}, {'path': 'vendor/QRCode/QRMaskPattern.js', 'source_sha256': 'b1f5a99876a31fccbfba89b973e11a4eb295f47b4b00e923814215309c0a725e', 'source_git_blob': 'f6fdeb53730c8b5aecec0182543955f112408df9', 'start': 2069844, 'end': 2070503, 'region_sha256': '7e34365450e4bf9a0a35f7075df7dc7ba336183767469919d695d14201fbf29c'}, {'path': 'vendor/QRCode/QRUtil.js', 'source_sha256': '4191ce852ba66124c4f1fc3bb1a507f667193c0466731339d8c1e66a19aa6bc5', 'source_git_blob': 'e5b7d5b3cc542b49373420ee2f3d3b1d14f2863b', 'start': 2070504, 'end': 2078555, 'region_sha256': '939ca3276f9e9c17531ea48823952ee064d383f64ab1c4e305d131b6622f67f9'}, {'path': 'vendor/QRCode/QRErrorCorrectLevel.js', 'source_sha256': 'd11a145632cea07057084190e86243b3054f30fc77256dc5ea0dc0e0cae54608', 'source_git_blob': '9b4b30099d033344b7b4da95f6774a36d132511c', 'start': 2078556, 'end': 2079113, 'region_sha256': 'ac3dbc0fa27e54a4f201e4541db90dc3630af98683beb04a757fef56b39c7631'}, {'path': 'vendor/QRCode/QRRSBlock.js', 'source_sha256': '78281d6a39b575a1078f1f70e7311e4a3c8b67e15e5468c25521b64d6ff6b931', 'source_git_blob': 'd150af174607997d03d501312c47a6969440a32a', 'start': 2079114, 'end': 2086288, 'region_sha256': 'db806a292b99a099c6ad6012716e4a06f36c3665a1c55f0c4cc1e12ee43b1a1e'}, {'path': 'vendor/QRCode/QRBitBuffer.js', 'source_sha256': '0b5de11b341f5dd92caf3e3a26469f86fa3eb9b3795db6a489e4d53d91ecb67d', 'source_git_blob': 'e2861f68d1b38a93e6630c3403c0769ba41af376', 'start': 2086289, 'end': 2087573, 'region_sha256': '80048565dce106ebc0d7d47d5a1238c43fe4cb5ea34c6eac572c2cb96fea7232'}, {'path': 'vendor/QRCode/index.js', 'source_sha256': '7377be90fc61a40268acf7f30d5bd89c2fca99c57ef5391623de8c151b8da7df', 'source_git_blob': '10eb8eb0a06aa50d0d3c508f886a7728f52e5d98', 'start': 2087574, 'end': 2099556, 'region_sha256': 'cd16b1735b8762154fb6e8eca5a9fb3d3578c12f15e8cbfccc8bcfacc241f6a8'}]
QRCODE_COMMIT = '90f66cf5c6b10bcb4358df96a9580f9eb383307b'
QRCODE_HEADER = '//---------------------------------------------------------------------\n// QRCode for JavaScript\n//\n// Copyright (c) 2009 Kazuhiko Arase\n//\n// URL: http://www.d-project.com/\n//\n// Licensed under the MIT license:\n//   http://www.opensource.org/licenses/mit-license.php\n//\n// The word "QR Code" is registered trademark of \n// DENSO WAVE INCORPORATED\n//   http://www.denso-wave.com/qrcode/faqpatent-e.html\n//\n//---------------------------------------------------------------------\n// Modified to work in node for this project (and some refactoring)\n//---------------------------------------------------------------------\n'
QRCODE_LICENSE = 'MIT License\n\nCopyright (c) 2009 Kazuhiko Arase\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n'
LICENSE_URL = 'https://blueoakcouncil.org/license/1.0.0'
BLUEOAK = {'chownr': {'version': '3.0.0', 'license': 'BlueOak-1.0.0', 'sha256': '4300e90fdd91ec7035047473c60f880251a9801bd786302729d4277751d3b948', 'path': 'dist/node_modules/chownr/package.json', 'source': 'https://github.com/isaacs/chownr/blob/8b9800ac5fe4da0b58bffc9c66dd618f0721472d/package.json'}, 'isexe': {'version': '4.0.0', 'license': 'BlueOak-1.0.0', 'sha256': '1889746f84d1e2a524353a2e31464995b808f24457f0dd55d6fd9da638f978ed', 'path': 'dist/node_modules/isexe/package.json', 'source': 'https://github.com/isaacs/isexe/blob/2e7df7dabc4f68e88cf6b32b9029225ba52c6b0c/package.json'}, 'minipass': {'version': '7.1.3', 'license': 'BlueOak-1.0.0', 'sha256': 'b8ab116187d5e9375d494f17437eb862cbee4b329c21f15412c6a5170ec2f55c', 'path': 'dist/node_modules/minipass/package.json', 'source': 'https://github.com/isaacs/minipass/blob/ab4b3b05d0d557ac6bb178f38501b11d0c96454e/package.json'}, 'tar': {'version': '7.5.22', 'license': 'BlueOak-1.0.0', 'sha256': '31b18f4fb83ba7183f54fa9e2cdca610897108c37ba80570c2c6fa7ccf112d97', 'path': 'dist/node_modules/tar/package.json', 'source': 'https://github.com/isaacs/node-tar/blob/2a22bfc5d3a432a606d9da0e2d87ba634aa3b1cb/package.json'}, 'yallist': {'version': '5.0.0', 'license': 'BlueOak-1.0.0', 'sha256': '1b9d47057ce39814531ff93f668823b4fa03e7d23945449c274a1ff6d4cc297f', 'path': 'dist/node_modules/yallist/package.json', 'source': 'https://github.com/isaacs/yallist/blob/a082c2dd872bb3cd8fb92359ed66605064957cd4/package.json'}}
EMPTY_FILES = ('config/user.npmrc', 'config/global.npmrc')
EMPTY_DIRS = ('config', 'home', 'home/appdata', 'home/localappdata', 'home/config',
              'home/data', 'home/pnpm', 'tmp')



def blueoak_notice(projection):
    """Link-form notice for five exact root declarations; no invented archive licence."""
    rows = []
    for name, binding in BLUEOAK.items():
        raw = read_bounded(Path(projection) / 'payload' / binding['path'], 1048576)
        require(digest(raw) == binding['sha256'], 'BlueOak manifest changed: ' + name)
        value = decode(raw)
        require(value.get('name') == name and value.get('version') == binding['version']
                and value.get('license') == 'BlueOak-1.0.0', 'BlueOak declaration changed: ' + name)
        rows.append('| ' + name + '@' + binding['version'] + ' | ' + binding['path']
                    + ' | ' + binding['sha256'] + ' | ' + binding['source'] + ' |')
    text = ('# Supplemental Blue Oak license notice\n\n'
            'License: Blue Oak Model License 1.0.0 (BlueOak-1.0.0).\n'
            'License link: ' + LICENSE_URL + '\n\n'
            'This link-form notice accompanies the following five exact package-root declarations '
            'in the selected pnpm 11.25.0 candidate. Keep this notice with any copy of these works.\n\n'
            '| Package | Manifest in payload | Manifest SHA-256 | Fixed author declaration |\n'
            '|---|---|---|---|\n' + '\n'.join(rows) + '\n\n'
            'This supplemental file is generated by local tooling; it is not an original file '
            'from the pnpm archive. chownr has a fixed author declaration but no named source '
            'license file in the retained tree. The yallist source license scopes packages under '
            'src to their own licenses; this notice does not relicense any nested work.\n\n'
            'This notice covers the listed declarations only. All original notices remain in '
            'the unchanged payload. Other pnpm licenses, source obligations, publishing provenance '
            'and runtime/redistribution admission remain separate and unresolved.\n')
    return text.encode('utf-8')


def qrcode_notice(projection):
    """Deliver the retained vendor attribution and full MIT permission notice."""
    bundle = read_bounded(Path(projection) / 'payload' / QRCODE_BUNDLE, 20000000)
    require(digest(bundle) == QRCODE_BUNDLE_SHA256, 'QRCode bundle changed')
    require(len(QRCODE_MODULES) == 10, 'QRCode module inventory changed')
    rows = []
    for item in QRCODE_MODULES:
        require(0 <= item['start'] < item['end'] <= len(bundle), 'QRCode module range invalid')
        require(digest(bundle[item['start']:item['end']]) == item['region_sha256'], 'QRCode module changed')
        rows.append('| ' + item['path'] + ' | ' + item['source_sha256'] + ' | '
                    + item['region_sha256'] + ' |')
    text = ('# QRCode vendored attribution and MIT permission notice\n\n'
            'Component: qrcode-terminal 0.12.0, vendor/QRCode.\n'
            'Selected pnpm 11.25.0 bundle SHA-256: ' + QRCODE_BUNDLE_SHA256 + '\n'
            'Fixed source: https://github.com/gtanner/qrcode-terminal/tree/' + QRCODE_COMMIT + '/vendor/QRCode\n'
            'Permission text authority: https://opensource.org/license/mit\n\n'
            '## Retained author and modification notice\n\n' + QRCODE_HEADER + '\n'
            '## MIT permission notice\n\n' + QRCODE_LICENSE + '\n'
            '## Exact source and bundle-region inventory\n\n'
            '| Vendor path | Fixed source SHA-256 | Selected bundle region SHA-256 |\n'
            '|---|---|---|\n' + '\n'.join(rows) + '\n\n'
            'This is a locally assembled supplement. The source header is retained verbatim; '
            'the complete MIT permission text is supplied with the copyright identity from that header. '
            'It is not an original pnpm archive file. It does not replace the containing package\'s '
            'Apache license or any other notice. The path/hash inventory is not proof of a complete '
            'reproducible build or equivalent module bindings. Keep this notice with copies of the work. '
            'This planner does not distribute the payload or admit runtime/redistribution.\n')
    return text.encode('utf-8')


def semver_notice(projection):
    """Retain the exact dual-license file; select its full MIT option for this supplement."""
    bundle = read_bounded(Path(projection) / 'payload' / QRCODE_BUNDLE, 20000000)
    require(digest(bundle) == QRCODE_BUNDLE_SHA256, 'semver-utils bundle changed')
    item = SEMVER_REGION
    require(0 <= item['start'] < item['end'] <= len(bundle), 'semver-utils region invalid')
    require(digest(bundle[item['start']:item['end']]) == item['region_sha256'], 'semver-utils region changed')
    require(digest(SEMVER_LICENSE.encode('utf-8')) == SEMVER_LICENSE_SHA256, 'semver-utils original license changed')
    text = ('# semver-utils 1.1.4 original license notice\n\n'
            'Original npm artifact: https://registry.npmjs.org/semver-utils/-/semver-utils-1.1.4.tgz\n'
            'Artifact SHA-256: fa6980458971864f25f2afc7d7f6632b015d8b767296d48421e6ff13221bd2b3\n'
            'Original package/LICENSE SHA-256: ' + SEMVER_LICENSE_SHA256 + '\n'
            'Original package/semver-utils.js SHA-256: ' + item['source_member_sha256'] + '\n'
            'Selected pnpm bundle SHA-256: ' + QRCODE_BUNDLE_SHA256 + '\n'
            'Selected bundle region SHA-256: ' + item['region_sha256'] + '\n\n'
            'The original LICENSE explicitly permits either MIT or Apache-2.0. This supplement '
            'selects the complete MIT option and retains the entire original file below. '
            'The package manifest and registry declaration APACHEv2 remain unchanged as historical '
            'metadata; this is a locally assembled supplement, not a replacement pnpm archive file.\n\n'
            '## Verbatim original license\n\n' + SEMVER_LICENSE + '\n'
            '## Evidence scope\n\n'
            'The artifact was acquired in quarantine against exact SHA512 registry integrity. '
            'Its registry signature cryptographically verifies with a key expired on 2025-01-29; '
            'no current valid signature or trusted signing time is claimed. '
            'The source member and bundle region are separately fixed observations, not proof of '
            'complete source/build equivalence. Keep the copyright and permission notice with copies '
            'of the work. Other pnpm obligations and runtime/redistribution admission remain open.\n')
    return text.encode('utf-8')


def retained_notices(projection):
    """Deliver the reviewed evidence collection without converting it into license admission."""
    raw = read_bounded(Path(__file__).with_name('retained-notice-catalog.json'), 524288)
    require(digest(raw) == CATALOG_SHA256, 'retained notice catalog changed')
    catalog = decode(raw)
    require(set(catalog) == {'schema', 'status', 'release_ready', 'parent_manifest_sha256', 'entries'},
            'retained notice catalog fields changed')
    require(catalog['schema'] == 'elite-pnpm-retained-notice-catalog/v1'
            and catalog['status'] == 'RETAINED_TEXTS_ONLY' and catalog['release_ready'] is False,
            'retained notice scope inflated')
    require(catalog['parent_manifest_sha256'] == '27d666e86d769c368439b1cc20e0fa72a42f9366973216b42ced5434f003fa9d',
            'retained notice parent changed')
    require(isinstance(catalog['entries'], list) and len(catalog['entries']) == 47,
            'retained notice inventory changed')
    seen = set(); total = 0; originals = 0
    out = [b'# Retained pnpm notice and license-text evidence\n\n'
           b'47 exact text copies: 22 original selected-payload notices and 25 retained research texts.\n'
           b'This collection preserves evidence; it is not a complete recursive SBOM, license clearance, '
           b'source offer or source/relinking fulfillment. Runtime and redistribution admission remain open.\n'
           b'Texts can repeat across versions; 47 is not a count of uniquely cleared dependencies. '
           b'Original pnpm payload, BlueOak, QRCode and semver-utils supplements remain separate and unchanged.\n\n'
           b'Authority and qualification records: PNPM_NOTICE_COVERAGE_V338.md, '
           b'PNPM_NESTED_LICENSE_SOURCES_V339.md, PNPM_YARN_UNDICI_NOTICES_V340.md, '
           b'SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md. Tagged source and published-byte identity '
           b'must not be conflated. Follow the exact scope attached to each retained text.\n\n']
    for entry in catalog['entries']:
        require(set(entry) == {'path', 'source_locator', 'scope', 'origin', 'payload_path', 'bytes', 'sha256', 'content_base64'},
                'retained notice entry fields changed')
        path = safe_name(entry['path'])
        require(path.casefold() not in seen, 'duplicate retained notice path'); seen.add(path.casefold())
        require(type(entry['bytes']) is int and 0 < entry['bytes'] <= 65536, 'retained notice byte budget')
        require(isinstance(entry['content_base64'], str), 'retained notice encoding required')
        try:
            content = base64.b64decode(entry['content_base64'], validate=True)
        except (ValueError, binascii.Error) as exc:
            raise SelectionError('retained notice encoding invalid') from exc
        require(base64.b64encode(content).decode('ascii') == entry['content_base64'], 'noncanonical notice encoding')
        require(len(content) == entry['bytes'] and digest(content) == entry['sha256'], 'retained notice text changed')
        total += len(content); require(total <= 200000, 'retained collection byte budget')
        require(all(isinstance(entry[k], str) and entry[k] and '\n' not in entry[k] and '\r' not in entry[k]
                    for k in ('source_locator', 'scope')), 'retained notice provenance required')
        if entry['origin'] == 'ORIGINAL_SELECTED_PAYLOAD':
            require(path.startswith('original/') and entry['payload_path'] == path[len('original/'):],
                    'original notice path mismatch')
            actual = read_bounded(Path(projection) / 'payload' / safe_name(entry['payload_path']), 65536)
            require(actual == content, 'original payload notice changed'); originals += 1
        else:
            require(entry['origin'] == 'RETAINED_RESEARCH_EVIDENCE' and entry['payload_path'] is None
                    and not path.startswith('original/'), 'retained notice origin inflated')
        header = ('## ' + path + '\n\nOrigin: ' + entry['origin'] + '\nScope: ' + entry['scope']
                  + '\nSource evidence locator: ' + entry['source_locator'] + '\nSHA-256: '
                  + entry['sha256'] + '\nBytes: ' + str(entry['bytes']) + '\n\nBEGIN VERBATIM TEXT\n\n').encode('utf-8')
        out.extend([header, content, b'\n\nEND VERBATIM TEXT\n\n'])
    require(originals == 22 and total == 151253, 'retained notice coverage changed')
    result = b''.join(out); require(len(result) <= 262144, 'retained notice output budget')
    return result


def next_path_notice(projection):
    """Supply original covered source, manifest, license and bounded comparison qualifications."""
    raw = read_bounded(Path(__file__).with_name('next-path-source-evidence.json'), 65536)
    require(digest(raw) == NEXT_PATH_DATA_SHA256, 'next-path source evidence changed')
    data = decode(raw)
    require(data['schema'] == 'elite-next-path-source-evidence/v1'
            and data['package'] == 'next-path' and data['version'] == '1.0.0'
            and data['license'] == 'MPL-2.0' and data['runtime_admitted'] is False
            and data['status'] == 'SOURCE_AND_LICENSE_DELIVERY_ONLY', 'next-path evidence scope inflated')
    require([entry['path'] for entry in data['files']] == ['index.js', 'package.json', 'LICENSE'],
            'next-path source inventory changed')
    bundle = read_bounded(Path(projection) / 'payload' / QRCODE_BUNDLE, 20000000)
    require(digest(bundle) == QRCODE_BUNDLE_SHA256, 'next-path selected bundle changed')
    region = NEXT_PATH_REGION
    require(0 <= region['start'] < region['end'] <= len(bundle)
            and digest(bundle[region['start']:region['end']]) == region['sha256'],
            'next-path bundle region changed')
    proof = data['correspondence']
    require(proof['statements_compared'] == 4 and len(proof['negative_probes']) == 10
            and proof['source_code_executed'] is False and proof['runtime_equivalence_claimed'] is False
            and proof['whole_pnpm_build_proven'] is False and proof['npm_tarball_byte_identity_proven'] is False,
            'next-path comparison scope inflated')
    out = [('# next-path 1.0.0 — MPL source and license\n\n'
            'Covered component license: Mozilla Public License 2.0 (MPL-2.0).\n'
            'The original source files below are made available under that license; the full license is included. '
            'No additional restriction is imposed here on the rights granted for that covered source.\n'
            'Official fixed source: https://github.com/sholladay/next-path/tree/' + data['commit'] + '\n'
            'Selected pnpm bundle SHA-256: ' + QRCODE_BUNDLE_SHA256 + '\n'
            'Selected module region SHA-256: ' + region['sha256'] + '\n\n'
            '## Source and transformation scope\n\n'
            'The fixed package manifest publishes index.js as its main and sole files entry. '
            'All four module statements were compared structurally after only the explicitly listed adapters. '
            'Ten changed-source controls were rejected. This is source correspondence under those adapters, '
            'not runtime equivalence, an npm tarball byte match or a reproducible whole-pnpm build.\n'
            'Generated packaging changes: one __commonJS wrapper, two named top-level const declarations emitted as var, '
            'and the following identifier/loader mappings (emitted name to original name): '
            + ', '.join(k + ' -> ' + v for k, v in proof['identifier_renames'].items()) + '.\n'
            '__require(path) -> require(path) is a conditional loader adapter; actual host resolution remains unproven. '
            'The function body, property names, literals, argument order and export were preserved by the comparison. '
            'No additional business or algorithm change is claimed. Other pnpm components keep their own licenses.\n\n'
            '## Original files — complete, unchanged bytes\n\n').encode('utf-8')]
    for entry in data['files']:
        require(type(entry['bytes']) is int and 0 < entry['bytes'] <= 32768, 'next-path source byte budget')
        try:
            content = base64.b64decode(entry['content_base64'], validate=True)
        except (ValueError, binascii.Error) as exc:
            raise SelectionError('next-path source encoding invalid') from exc
        require(len(content) == entry['bytes'] and digest(content) == entry['sha256'], 'next-path original source changed')
        out.extend([('### ' + entry['path'] + '\n\nSource: ' + entry['source_url'] + '\nGit blob: '
                     + entry['git_blob'] + '\nSHA-256: ' + entry['sha256'] + '\nBytes: '
                     + str(entry['bytes']) + '\n\nBEGIN ORIGINAL FILE\n\n').encode('utf-8'),
                    content, b'\n\nEND ORIGINAL FILE\n\n'])
    out.append(b'This local delivery does not authorize redistribution of the entire pnpm bundle, '
               b'waive other source/relinking obligations, or admit any runtime. '
               b'No supplied source file is installed or executed by the planner.\n')
    result = b''.join(out); require(len(result) <= 65536, 'next-path notice output budget')
    return result

def paths(project, projection, node, store, cache, target):
    result = [no_reparse(x) for x in (project, projection, node, store, cache, target)]
    project, projection, node, store, cache, target = result
    require(all(p.is_dir() for p in (project, projection, store, cache)), 'input directory missing')
    require(node.is_file(), 'Node executable missing')
    for i, a in enumerate(result):
        for b in result[i + 1:]:
            require(a != b and a not in b.parents and b not in a.parents, 'overlapping inputs/output rejected')
    return result


def consumer_files(consumer, project):
    require(consumer in PROFILES, 'unsupported consumer')
    policy = PROFILES[consumer]
    require(not os.path.lexists(project / 'node_modules'), 'fresh consumer without node_modules required')
    for name, expected in policy.items():
        require(digest(read_bounded(project / name, 2097152)) == expected, 'consumer input changed: ' + name)
    for ancestor in (project, *project.parents):
        for name in ('.npmrc', '.pnpmfile.cjs', 'pnpmfile.cjs', '.pnpmfile.mjs', 'pnpmfile.mjs'):
            require(not os.path.lexists(ancestor / name), 'ambient project configuration rejected')
    if 'pnpm-workspace.yaml' not in policy:
        require(not os.path.lexists(project / 'pnpm-workspace.yaml'), 'unsupported workspace configuration')
    require(not os.path.lexists(project / 'pnpm-workspace.yml'), 'unsupported workspace extension')
    return dict(policy)


def recipe(consumer, project, projection, node, store, cache, target, acquisition_receipt, acquisition_sha256):
    require(os.name == 'nt', 'routing supported only on Windows')
    project, projection, node, store, cache, target = paths(project, projection, node, store, cache, target)
    acquisition_receipt = no_reparse(acquisition_receipt)
    require(target != acquisition_receipt and target not in acquisition_receipt.parents, 'receipt/output overlap')
    raw = read_bounded(acquisition_receipt, 1048576)
    require(digest(raw) == acquisition_sha256, 'acquisition identity changed')
    policy = load_policy()
    acquisition_binding(raw, policy)
    selection = verify(projection, policy)
    require(selection['acquisition_receipt_sha256'] == acquisition_sha256, 'projection acquisition mismatch')
    require(digest(read_bounded(node, 150000000)) == NODE_SHA256, 'Node identity changed')
    files = consumer_files(consumer, project)
    system = no_reparse(os.environ.get('SystemRoot', ''))
    require(system.is_dir() and system.name.casefold() == 'windows', 'trusted Windows SystemRoot required')
    environment = {
        'SystemRoot': str(system), 'WINDIR': str(system),
        'PATH': str(node.parent) + os.pathsep + str(system / 'System32'),
        'TEMP': str(target / 'tmp'), 'TMP': str(target / 'tmp'),
        'HOME': str(target / 'home'), 'USERPROFILE': str(target / 'home'),
        'APPDATA': str(target / 'home/appdata'), 'LOCALAPPDATA': str(target / 'home/localappdata'),
        'XDG_CONFIG_HOME': str(target / 'home/config'), 'XDG_DATA_HOME': str(target / 'home/data'),
        'PNPM_HOME': str(target / 'home/pnpm'), 'CI': 'true', 'NEXT_TELEMETRY_DISABLED': '1',
    }
    argv = [str(node), str(projection / 'payload/bin/pnpm.mjs'), 'install',
            '--offline', '--frozen-lockfile', '--ignore-scripts', '--ignore-pnpmfile',
            '--package-import-method=copy', '--verify-store-integrity',
            '--store-dir=' + str(store), '--config.cache-dir=' + str(cache), '--config.userconfig=' + str(target / EMPTY_FILES[0]),
            '--config.globalconfig=' + str(target / EMPTY_FILES[1])]
    if 'pnpm-workspace.yaml' not in files:
        argv.append('--ignore-workspace')
    return {'schema': 'elite-pnpm-install-plan/v6', 'consumer': consumer,
            'inputs': {'project': str(project), 'projection': str(projection), 'node': str(node),
                       'store': str(store), 'cache': str(cache), 'acquisition_receipt': str(acquisition_receipt),
                       'acquisition_receipt_sha256': acquisition_sha256},
            'consumer_files': files, 'node_sha256': NODE_SHA256,
            'next_path_source': {'path': NEXT_PATH_NOTICE, 'sha256': digest(next_path_notice(projection)),
                                 'data_sha256': NEXT_PATH_DATA_SHA256, 'scope': 'SOURCE_AND_LICENSE_DELIVERY_ONLY'},
            'retained_notices': {'path': RETAINED_NOTICE, 'sha256': digest(retained_notices(projection)),
                                 'catalog_sha256': CATALOG_SHA256, 'copies': 47, 'scope': 'RETAINED_TEXTS_ONLY'},
            'semver_notice': {'path': SEMVER_NOTICE, 'sha256': digest(semver_notice(projection)),
                              'scope': 'ORIGINAL_SEMVER_UTILS_LICENSE_MIT_OPTION'},
            'qrcode_notice': {'path': QRCODE_NOTICE, 'sha256': digest(qrcode_notice(projection)),
                              'scope': 'QRCODE_VENDOR_ATTRIBUTION_AND_MIT_PERMISSION'},
            'supplemental_notice': {'path': NOTICE, 'sha256': digest(blueoak_notice(projection)),
                                    'license_url': LICENSE_URL, 'scope': 'FIVE_PINNED_ROOT_DECLARATIONS'},
            'selection_receipt_sha256': digest(read_bounded(projection / 'selection-receipt.json', 4096)),
            'cwd': str(project), 'argv': argv, 'environment': environment,
            'environment_mode': 'REPLACE_NOT_MERGE', 'shell': False,
            'runtime_admitted': False, 'redistribution_admitted': False, 'executed': False,
            'conditions': ['Verify this plan and its externally retained SHA256 immediately before use.',
                           'Separate runtime/license admission required before execution.',
                           'Keep NEXT-PATH-MPL-SOURCE.md, PNPM-RETAINED-NOTICES.md, BLUEOAK-NOTICE.md, QRCODE-NOTICE.md and SEMVER-UTILS-NOTICE.md with copies of their listed works; this plan does not distribute the payload.',
                           'Offline and ignore-scripts flags are not an operating-system sandbox.',
                           'Trusted local inputs must remain immutable through any later use.',
                           'No scripts, remote resolution, cloning, global installs or arbitrary CLI arguments in this recipe.']}


def verify_plan(target, expected_sha256):
    target = no_reparse(target)
    raw = read_bounded(target / RECIPE, 65536)
    require(digest(raw) == expected_sha256, 'install plan identity changed')
    value = decode(raw)
    inp = value['inputs']
    expected = recipe(value['consumer'], inp['project'], inp['projection'], inp['node'], inp['store'], inp['cache'],
                      target, inp['acquisition_receipt'], inp['acquisition_receipt_sha256'])
    require(raw == json_bytes(expected), 'install plan changed')
    found_files, found_dirs = set(), set()
    for root, dirs, files in os.walk(target, followlinks=False):
        for name in dirs + files:
            p = no_reparse(Path(root) / name)
            rel = p.relative_to(target).as_posix()
            if p.is_dir():
                require(rel in EMPTY_DIRS, 'unexpected plan directory')
                found_dirs.add(rel)
            else:
                require(rel in (*EMPTY_FILES, RECIPE, NOTICE, QRCODE_NOTICE, SEMVER_NOTICE, RETAINED_NOTICE, NEXT_PATH_NOTICE), 'unexpected plan file')
                found_files.add(rel)
                if rel == NEXT_PATH_NOTICE:
                    require(digest(read_bounded(p, 65536)) == expected['next_path_source']['sha256'], 'next-path source delivery changed')
                if rel == RETAINED_NOTICE:
                    require(digest(read_bounded(p, 262144)) == expected['retained_notices']['sha256'], 'retained notice delivery changed')
                if rel == SEMVER_NOTICE:
                    require(digest(read_bounded(p, 16384)) == expected['semver_notice']['sha256'], 'semver-utils notice changed')
                if rel == QRCODE_NOTICE:
                    require(digest(read_bounded(p, 16384)) == expected['qrcode_notice']['sha256'], 'QRCode notice changed')
                if rel == NOTICE:
                    require(digest(read_bounded(p, 16384)) == expected['supplemental_notice']['sha256'], 'supplemental notice changed')
                if rel in EMPTY_FILES:
                    require(read_bounded(p, 0) == b'', 'plan configuration is not empty')
    require(found_files == {*EMPTY_FILES, RECIPE, NOTICE, QRCODE_NOTICE, SEMVER_NOTICE, RETAINED_NOTICE, NEXT_PATH_NOTICE} and found_dirs == set(EMPTY_DIRS), 'plan files/directories missing')
    return value


def create_plan(consumer, project, projection, node, store, cache, target, acquisition_receipt, acquisition_sha256):
    target = no_reparse(target)
    require(target.parent.is_dir() and not os.path.lexists(target), 'target must be absent with existing parent')
    value = recipe(consumer, project, projection, node, store, cache, target, acquisition_receipt, acquisition_sha256)
    raw = json_bytes(value)
    stage = Path(tempfile.mkdtemp(prefix='.pnpm-plan-', dir=target.parent))
    try:
        for name in EMPTY_DIRS:
            (stage / name).mkdir(exist_ok=True)
        for name in EMPTY_FILES:
            (stage / name).write_bytes(b'')
        (stage / NOTICE).write_bytes(blueoak_notice(projection))
        (stage / QRCODE_NOTICE).write_bytes(qrcode_notice(projection))
        (stage / SEMVER_NOTICE).write_bytes(semver_notice(projection))
        (stage / RETAINED_NOTICE).write_bytes(retained_notices(projection))
        (stage / NEXT_PATH_NOTICE).write_bytes(next_path_notice(projection))
        (stage / RECIPE).write_bytes(raw)
        # Recipe paths refer to the final target. Publication never replaces a competitor.
        no_reparse(target.parent)
        os.rename(stage, target)
    finally:
        if stage.exists():
            require(stage.parent == target.parent and stage.name.startswith('.pnpm-plan-'), 'cleanup scope rejected')
            shutil.rmtree(stage)
    verify_plan(target, digest(raw))
    return digest(raw)


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest='command', required=True)
    make = sub.add_parser('prepare')
    make.add_argument('--consumer', choices=sorted(PROFILES), required=True)
    for name in ('project', 'projection', 'node', 'store', 'cache', 'target', 'acquisition-receipt', 'acquisition-sha256'):
        make.add_argument('--' + name, required=True)
    check = sub.add_parser('verify')
    check.add_argument('--target', required=True)
    check.add_argument('--sha256', required=True)
    args = parser.parse_args(argv)
    try:
        if args.command == 'prepare':
            sha = create_plan(args.consumer, args.project, args.projection, args.node, args.store, args.cache,
                              args.target, args.acquisition_receipt, args.acquisition_sha256)
        else:
            verify_plan(args.target, args.sha256)
            sha = args.sha256
        print('PNPM_INSTALL_PLAN_PASS sha256=' + sha + ' executed=false runtime_admitted=false')
        return 0
    except (SelectionError, OSError, ValueError, UnicodeError, KeyError, TypeError, RecursionError) as exc:
        print('PNPM_INSTALL_PLAN_BLOCKED: ' + str(exc), file=sys.stderr)
        return 2

if __name__ == '__main__':
    raise SystemExit(main())
