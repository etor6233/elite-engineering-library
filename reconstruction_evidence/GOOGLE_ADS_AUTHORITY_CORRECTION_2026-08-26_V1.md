# Google Ads authority correction — 2026-08-26 V1

## Contradicción

El lock Elite declaraba `google-ads-python-31.4.0` con commit `b8eb80ae...`. La fuente primaria PyPI vigente al 2026-08-26 identifica `31.3.0` como release más reciente, publicada 2026-08-13. El claim 31.4.0 se congeló antes de crear un adapter.

## Autoridad corregida

- distribución: `google-ads==31.3.0`;
- PyPI wheel: 23.675.722 bytes; SHA-256 `286946249d54dea97ad4c81aa5d466a49060cc726f0846281a4e756d24de9890`;
- PyPI sdist: 12.431.987 bytes; SHA-256 `8f42581e9db633562e2664fd0e583af7c6064d6a2e991a6910c4b5c9cd30b28d`;
- tag GitHub `31.3.0`: commit `f7bf312d26904ea50ad9ef02e395edccfaa1ebb7`;
- source ZIP descargado: 25.342.657 bytes; SHA-256 `bdf483432bf74142c757b07cea54dfc6402b488bf2a007befa526a92c4c0731a`;
- licencia: Apache-2.0; `LICENSE` SHA-256 `c62a1cc62e85e51a18ddd590970e03ffe47017653033cc515781fdb445a7aa5e`.

## Cambios y gates

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` subió a 0.4.3. El ID del lock y el perfil commerce/communications ahora usan `google-ads-python-31.3.0`; planes dependientes fijan 0.4.3. La admisión continúa `INTEGRATION_ONLY`: todavía faltan developer token, OAuth/service account compatible, customer/login customer hierarchy, cuenta test, API access level, costo/cuota, conversion policy y contract probe.

## Fuentes primarias

- PyPI JSON: https://pypi.org/pypi/google-ads/31.3.0/json
- tag oficial: https://github.com/googleads/google-ads-python/tree/f7bf312d26904ea50ad9ef02e395edccfaa1ebb7
- documentación oficial: https://developers.google.com/google-ads/api/docs/client-libs/python
