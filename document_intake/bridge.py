"""AUTHORED channel-to-document quarantine bridge; no SMTP/SFTP server or OCR.

Call only with version-bound receipts emitted by the admitted intake owner.
Incoming documents remain QUARANTINED. Processing and independent approval
belong to the existing document owner, never this transport adapter.
"""
from dataclasses import dataclass
from pathlib import Path
import hashlib, json, re, urllib.request, urllib.error, urllib.parse, uuid

MAX_BYTES = 2 * 1024 * 1024
HEX = re.compile(r"^[0-9a-f]{64}$")


class IntakeError(Exception):
    pass


@dataclass(frozen=True)
class Profile:
    origin: str
    tenant: str
    organization: str
    profile_sha256: str
    mode: str
    fixture_loopback: bool = False

    def validate(self):
        parsed = urllib.parse.urlsplit(self.origin)
        local = self.fixture_loopback and parsed.scheme == 'http' and parsed.hostname in ('127.0.0.1', 'localhost', '::1')
        if (parsed.scheme != 'https' and not local) or parsed.path not in ('', '/') or parsed.query or parsed.fragment or parsed.username or parsed.password:
            raise IntakeError('ORIGIN_REJECTED')
        if not HEX.fullmatch(self.profile_sha256) or not self.tenant or not self.organization or self.mode not in ('FIXTURE', 'TYPED_FIXTURE', 'PROVIDER'):
            raise IntakeError('PROFILE_REJECTED')
        if self.fixture_loopback and self.mode == 'PROVIDER':
            raise IntakeError('PROVIDER_REQUIRES_TLS')


class _NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def sha(data):
    return hashlib.sha256(data).hexdigest()


def _identity(profile, channel, source, part, class_id):
    profile.validate()
    if channel not in ('email', 'sftp') or type(part) is not int or part < 1:
        raise IntakeError('CHANNEL_REJECTED')
    if not re.fullmatch(r'[a-z][a-z0-9-]{0,79}', class_id):
        raise IntakeError('CLASS_REJECTED')
    if not isinstance(source, dict) or any(not isinstance(source.get(k), str) or not source[k] or len(source[k]) > 1024 for k in ('bucket', 'key', 'version_id')):
        raise IntakeError('VERSION_BOUND_SOURCE_REQUIRED')
    # Content hash is deliberately NOT the identity: changing a retained version's
    # bytes must conflict, rather than silently create a second document.
    binding = json.dumps([profile.tenant, profile.organization, channel, source['bucket'], source['key'], source['version_id'], part], separators=(',', ':'))
    return str(uuid.uuid5(uuid.NAMESPACE_URL, 'elite-document-intake/v1:' + binding))


def _request(origin, token, path, method='GET', body=None, headers=None):
    if not token or '\r' in token or '\n' in token:
        raise IntakeError('AUTH_REQUIRED')
    request = urllib.request.Request(origin.rstrip('/') + path, data=body, method=method,
        headers={'Authorization': 'Bearer ' + token, **(headers or {})})
    try:
        response = urllib.request.build_opener(_NoRedirect).open(request, timeout=15)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        payload = response.read(MAX_BYTES + 1)
        if len(payload) > MAX_BYTES:
            raise IntakeError('RESPONSE_TOO_LARGE')
        return response.status, json.loads(payload)


def receive(profile, token, channel, source, part, class_id, data, expected_sha256):
    document_id = _identity(profile, channel, source, part, class_id)
    if not isinstance(data, bytes) or not 0 < len(data) <= MAX_BYTES or not HEX.fullmatch(expected_sha256) or sha(data) != expected_sha256:
        raise IntakeError('ORIGINAL_HASH_OR_SIZE')
    suffix = 'pdf' if data.startswith(b'%PDF-') else 'jpg' if data.startswith(b'\xff\xd8\xff') else None
    if suffix is None:
        raise IntakeError('ORIGINAL_TYPE')
    path = '/v1/documents/' + document_id

    def validate(value):
        if (value.get('document_id') != document_id or value.get('original_sha256') != expected_sha256 or
            value.get('profile_sha256') != profile.profile_sha256 or value.get('mode') != profile.mode or
            value.get('class_id') != class_id or value.get('schema_version') != '1'):
            raise IntakeError('RECORDED_STATE_CONFLICT')
        return value

    # Read-before-write deduplicates deliveries and recovery across process loss.
    status, value = _request(profile.origin, token, path)
    if status == 200:
        return {'state': 'RECONCILED', 'document': validate(value), 'mutations': 0}
    if status != 404:
        raise IntakeError('READ_NOT_AUTHORIZED_OR_UNAVAILABLE')
    try:
        status, value = _request(profile.origin, token, path + '/original', 'PUT', data, {
            'Content-Type': 'application/octet-stream', 'X-Document-Name': class_id + '.' + suffix,
            'X-Document-SHA256': expected_sha256, 'X-Document-Profile-SHA256': profile.profile_sha256,
            'X-Document-Class': class_id, 'X-Document-Schema-Version': '1'})
        if status in (200, 201, 202):
            return {'state': 'QUARANTINED', 'document': validate(value), 'mutations': 1}
        if status in (400, 401, 403, 422):
            raise IntakeError('RECEIVE_REJECTED')
    except (OSError, TimeoutError, json.JSONDecodeError):
        pass
    # Never blindly resend an uncertain write. A GET can establish its outcome.
    status, value = _request(profile.origin, token, path)
    if status == 200:
        return {'state': 'RECONCILED_AFTER_UNCERTAIN', 'document': validate(value), 'mutations': 1}
    raise IntakeError('UNCONFIRMED_USE_SAME_VERSION_AND_READ')


def email_manifest(profile, token, manifest, quarantine, class_by_part):
    root = Path(quarantine).resolve()
    if manifest.get('state') != 'QUARANTINE_EXTRACTED' or manifest.get('automatic_storage_authorized') is not False or manifest.get('tenant_id') != profile.tenant:
        raise IntakeError('QUARANTINE_MANIFEST_REQUIRED')
    results = []
    for attachment in manifest['attachments']:
        name = attachment.get('quarantine_name', '')
        if not re.fullmatch(r'part-[0-9]{4}-[0-9a-f]{64}\.bin', name):
            raise IntakeError('QUARANTINE_PATH')
        path = root / name
        if path.is_symlink() or path.resolve().parent != root or path.stat().st_size > MAX_BYTES:
            raise IntakeError('QUARANTINE_PATH_OR_SIZE')
        part = attachment['part_index']
        results.append(receive(profile, token, 'email', manifest['source'], part,
            class_by_part[part], path.read_bytes(), attachment['sha256']))
    return results


def sftp_object(profile, token, receipt, data, class_id):
    # Provider completion/version evidence is mandatory. A mere filesystem create
    # event does not prove that the writer has finished uploading the object.
    if receipt.get('schema') != 'elite.sftp-retained-object/v1' or receipt.get('tenant_id') != profile.tenant or receipt.get('organization_id') != profile.organization or receipt.get('complete') is not True or receipt.get('retained') is not True:
        raise IntakeError('COMPLETED_RETAINED_OBJECT_REQUIRED')
    return receive(profile, token, 'sftp', receipt['source'], 1, class_id, data, receipt['sha256'])


def azure_sftp_retained_receipt(profile, event, retained, expected_topic, expected_url_prefix):
    """Normalize an authenticated Event Grid event AFTER retained-object binding.

    Does not authenticate a webhook and never fetches an event-provided URL.
    The provider host must validate its delivery identity before calling this.
    """
    profile.validate()
    data = event.get('data', {})
    if (event.get('eventType') != 'Microsoft.Storage.BlobCreated' or event.get('topic') != expected_topic or
        not expected_topic.startswith('/subscriptions/') or data.get('api') != 'SftpCommit'):
        raise IntakeError('SFTP_COMMIT_EVENT_REQUIRED')
    parsed = urllib.parse.urlsplit(expected_url_prefix)
    if parsed.scheme != 'https' or parsed.username or parsed.password or parsed.query or parsed.fragment or not parsed.path.endswith('/'):
        raise IntakeError('SFTP_SCOPE_PREFIX')
    url = data.get('url', '')
    if (not isinstance(url, str) or not url.startswith(expected_url_prefix) or
        retained.get('source_url') != url or not data.get('eTag') or retained.get('source_etag') != data['eTag'] or
        type(data.get('contentLength')) is not int or not 1 <= data['contentLength'] <= MAX_BYTES or
        retained.get('bytes') != data['contentLength'] or retained.get('tenant_id') != profile.tenant or
        retained.get('organization_id') != profile.organization or retained.get('retained') is not True or
        not HEX.fullmatch(retained.get('sha256', ''))):
        raise IntakeError('SFTP_RETAINED_BINDING_REQUIRED')
    return {'schema': 'elite.sftp-retained-object/v1', 'tenant_id': profile.tenant,
            'organization_id': profile.organization, 'complete': True, 'retained': True,
            'source': retained['source'], 'sha256': retained['sha256']}


def main():
    import argparse, os
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--profile', type=Path, required=True)
    parser.add_argument('--channel', choices=('email', 'sftp'), required=True)
    parser.add_argument('--receipt', type=Path, required=True)
    parser.add_argument('--input', type=Path, required=True, help='quarantine directory (email) or retained object (SFTP)')
    parser.add_argument('--classes', type=Path, required=True, help='JSON mapping part indexes to admitted class IDs')
    parser.add_argument('--output', type=Path, required=True, help='new result receipt; never overwrite evidence')
    args = parser.parse_args()
    if args.output.exists():
        raise IntakeError('OUTPUT_EXISTS')
    profile = Profile(**json.loads(args.profile.read_text('utf8')))
    profile.validate()
    token = os.environ.get('DOCUMENT_INTAKE_BEARER', '')
    receipt = json.loads(args.receipt.read_text('utf8'))
    classes = {int(k): v for k, v in json.loads(args.classes.read_text('utf8')).items()}
    if args.channel == 'email':
        results = email_manifest(profile, token, receipt, args.input, classes)
    else:
        if args.input.is_symlink() or args.input.stat().st_size > MAX_BYTES:
            raise IntakeError('OBJECT_SIZE_OR_SYMLINK')
        results = [sftp_object(profile, token, receipt, args.input.read_bytes(), classes[1])]
    # Do not duplicate document contents/fields in the transport evidence.
    safe = [{'state': r['state'], 'document_id': r['document']['document_id'],
             'original_sha256': r['document']['original_sha256'], 'mutations': r['mutations']} for r in results]
    with args.output.open('x', encoding='utf8') as output:
        json.dump({'schema': 'elite.document-intake-result/v1', 'channel': args.channel,
                   'result': 'PASS_QUARANTINE_ONLY', 'documents': safe,
                   'automatic_processing': False, 'production_authorized': False}, output, indent=2)
    print('PASS_QUARANTINE_ONLY', len(safe))


if __name__ == '__main__':
    main()
