# OpenSSH public regression fixture notice

The four Base64 files under `fixtures/` losslessly reconstruct bytes copied verbatim from the OpenSSH portable `V_10_2_P1` commit `d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3`, tree `c91beec59da56218afced35b14087d362b5bba8e`:

- `regress/unittests/sshsig/testdata/ed25519.pub`, blob `b078e4516fbe44d2cf03e32a98e624bcb2ea6875`;
- `regress/unittests/sshsig/testdata/ed25519.sig`, blob `8e8ff2a8ac197df506ee00eef1d398a8090609c5`;
- `regress/unittests/sshsig/testdata/namespace`, blob `1570cd548baa9998c08cd500e88daa254dbbe66c`;
- `regress/unittests/sshsig/testdata/signed-data`, blob `7df4bedd135c385ed6e5cf99c901e3ad6e44a043`.

OpenSSH portable publishes its full license and copyright notices at `https://github.com/openssh/openssh-portable/blob/d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3/LICENCE` (blob `aeb3017e76c6ca6a53337dfb4af2f434e1628568`, 21,140 bytes, SHA-256 `5bb5b160726ef5756e4f32fe95b35249c294962419650f48d05134b486d27ccb`). Preserve this notice and the linked upstream license when redistributing the fixtures. No OpenSSH executable, source implementation, or private key is embedded by this pack.
