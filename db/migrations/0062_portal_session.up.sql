-- AUTHORED opaque session storage/CAS glue. No plaintext OAuth credentials.
create table platform.portal_session (
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 session_id_sha256 bytea not null check(octet_length(session_id_sha256)=32),
 version bigint not null default 1 check(version>0),
 state text not null check(state in ('active','refreshing','reauth_required','revoked')),
 operation_id text not null default '',
 ciphertext text not null check(length(ciphertext)<=32768),
 ciphertext_sha256 bytea not null check(octet_length(ciphertext_sha256)=32),
 access_expires_at timestamptz not null,
 absolute_expires_at timestamptz not null,
 revocation_pending boolean not null default false,
 revocation_lease_until timestamptz not null default '-infinity',
 revocation_operation_id text not null default '',
 created_at timestamptz not null default clock_timestamp(),
 updated_at timestamptz not null default clock_timestamp(),
 primary key(profile_sha256,session_id_sha256),
 check(access_expires_at<=absolute_expires_at),
 check((state='active' and not revocation_pending) or state<>'active')
);
create index portal_session_expiry on platform.portal_session(absolute_expires_at);
create table platform.portal_session_retention_summary (
 profile_sha256 text primary key check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 purged_count bigint not null default 0,
 unconfirmed_provider_revocation_count bigint not null default 0,
 updated_at timestamptz not null default clock_timestamp()
);
