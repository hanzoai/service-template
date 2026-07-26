# hanzoai/template — canonical Go backend skeleton

This is the one-way shape for a Hanzo-stack Go backend. Copy the tree, rename
`template` everywhere, and wire your routes. Do **not** deviate from the
layout without a spec change in `~/work/hanzo/ARCHITECTURE.md`.

The skeleton ships:

- `cmd/app/` — operator CLI (migrate, seed, shell).
- `cmd/appd/` — long-running HTTP daemon.
- `config/config.go` — env-only configuration.
- `api/` — HTTP handlers (public + tenant-scoped).
- `models/` — `orm.Typed[T]` data models.
- `migrations/` — base TS migrations.
- `store/` — thin wrapper around `hanzoai/base/store.MultiTenantStore`.
- `web/` — optional embedded SPA.
- `LLM.md` — next-agent handoff.

The key invariants it enforces:

1. Gateway-upstream assertion at boot (`claims.AssertGatewayUpstream`).
2. `claims.Strip → claims.Inject → claims.RequireGateway` on every
   tenant-scoped route.
3. Per-(org, user) SQLite via `store.ForCtx(ctx)`. No direct file opens.
4. Public endpoints (`/healthz`, `/readyz`, `/metrics`) mount BEFORE the
   tenant middleware chain.
5. Single multi-arch image built in CI — never locally.

To bootstrap a new service:

```
cp -R ~/work/hanzo/_template ~/work/hanzo/<service>
cd ~/work/hanzo/<service>
find . -type f -exec sed -i '' 's/template/<service>/g' {} +
go mod init github.com/hanzoai/<service>
go mod tidy
go test ./...
```

Then replace `models/example.go` and `api/example.go` with your domain.
