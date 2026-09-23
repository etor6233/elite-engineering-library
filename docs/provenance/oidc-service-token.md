# Service token provenance and admission boundary

AUTHORED: profile/hash parsing, exact authority comparisons, bounded transport, external secret reader, lifetime/cache concurrency glue, host contract and all local fixtures. None is attributed to an external company. Existing Principal/OIDCVerifier remains owned by GO-ENTERPRISE-BACKEND-CORE and is not recopied by this pack.

DEPENDENCY_PIN: golang.org/x/oauth2 v0.36.0, clientcredentials.Config.Token executes the grant and builds its Basic/Post authentication. github.com/coreos/go-oidc/v3 v3.20.0, NewProvider and IDTokenVerifier execute discovery, signature/issuer/audience/expiry verification and RemoteKeySet key rotation. github.com/go-jose/go-jose/v4 v4.1.4 signs only ephemeral test-fixture tokens. No module/version is introduced. Exact local module info, source files, licenses, go.mod and go.sum are hashed in oidc-service-source-receipt.json.

Official method authorities: https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/clientcredentials ; https://pkg.go.dev/github.com/coreos/go-oidc/v3@v3.20.0/oidc ; https://www.rfc-editor.org/rfc/rfc6749#section-4.4 . Sources are inspected; SDK integration tests execute the fixed existing modules. No upstream full-suite execution or worldwide provider compatibility is claimed.

G0/G1 source identity/license, G2 security boundary and G3 failure design are recorded before coding in oidc-service-admission.md. Local G4 unit/contract, G5 claim/transport negatives plus finite fuzz, G6 concurrency/renewal and G7 recovery/TLS checks are evidenced by exact logs. G8 exact pack reconstruction is separately recorded. This scoped evidence does not close root composition SCA/SAST/security/release or certify production.

AUTHORED glue is inevitable here: a provider SDK cannot know the library tenant, organization, connector subject or permission boundary, configured exact endpoints, external secret path or host's cancellation/poll semantics. No custom signing, OAuth token protocol or IdP framework is implemented.
