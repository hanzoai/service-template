# LLM.md — next-agent handoff for hanzoai/template

This is the canonical skeleton every new Hanzo Go backend starts from. It is
kept minimal on purpose: if a feature is not in every service, it does not
belong here.

## What the skeleton enforces

- `config/config.go` pulls env; nothing else. No flag soup. No YAML.
- `cmd/appd/main.go` calls `claims.AssertGatewayUpstream()` at boot. The
  daemon refuses to start without `HANZO_GATEWAY_UPSTREAM=1`.
- `api/router.go` is the one place middleware is wired. Top-to-bottom order:
  `claims.Strip`, `claims.Inject`, `claims.RequireGateway` for tenant routes.
  Public routes (`/healthz`, `/readyz`) mount separately, before the chain.
- `store/store.go` is the only caller of `hanzoai/base/store`. Handlers call
  `app.Store.ForCtx(ctx)` and get a `*dbx.DB` scoped to the caller's
  (org, user).
- `models/example.go` shows the `orm.Typed[Example]` pattern. Replace with
  your domain types — do not grow a domain-specific layer on top.

## What NOT to add

- JWT middleware of any kind. The gateway is the sole verifier.
- A second cache library. `luxfi/cache` is it.
- `nginx`, `caddy`, or `envoy` configs. Ingress is `hanzoai/ingress`.
- Random `SUMMARY.md` / `STATUS.md` / `NOTES.md`. Only `LLM.md`.
- A Helm chart. Kustomize lives in `hanzoai/platform`.

## Dependencies

Only three direct imports in `go.mod`:

- `github.com/hanzoai/base` — runtime + claims + store.
- `github.com/hanzoai/orm` — data layer.
- `github.com/luxfi/cache` — cache primitives.

Anything else is a review-blocking addition. Third-party libs go through
`hanzoai/base/tools/*` or get their own module.
