# Mercado Libre marketplace adapter

This authored Go adapter follows the current official Mercado Libre HTTP documentation. It does **not** use or revive Mercado Libre's archived SDK, which the provider marks not functional and unmaintained.

Implemented claims:

- authenticated `GET /items/{id}`;
- authenticated `GET /orders/{id}`;
- authenticated `GET /questions/{id}?api_version=4`;
- seller reconciliation with `GET /questions/search?seller_id=...&status=UNANSWERED&api_version=4`;
- approved `POST /items/validate`, preserving 204 acceptance or 400 validation details;
- notification envelope allowlisting for `items`, `orders_v2`, `shipments` and `questions`, durable fetch jobs and raw-body deduplication.

The template is blocked. Before use, prove the legal application owner, seller authorization, secure OAuth token storage, scopes, site/seller/application IDs, item/category rules, the lack of a preproduction environment, origin controls for notifications, topics, quotas/cost and reconciliation. The bearer token exists only in the named environment variable.

```powershell
go mod verify
go test ./...
go vet ./...
go run ./cmd/meli-marketplace -profile .\provider-profile.json -request .\request.json -output .\evidence\operation-001
```

Question responses can be passed to the `GO-OMNICHANNEL-LEAD-INGRESS` Mercado Libre v4 decoder. This pack deliberately does not publish, modify or delete listings; change price/stock; refresh OAuth tokens; answer questions; acknowledge claims; update shipments; or infer a provider response as a successful business write. Those capabilities require separate official-contract packs, an outbound fence and real-account gates.
