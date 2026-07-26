// Package api exposes the service's HTTP surface.
//
// Middleware order — fixed, not configurable — is provided by the
// canonical [claims.Chain] helper in hanzoai/base/tools/claims:
//
//	claims.Strip        → drop any forged identity headers
//	claims.Inject       → parse the canonical 3 headers into Claims
//	claims.RequireGateway → 503 if either X-User-Id or X-Org-Id is empty
//
// Tenant-scoped routes MUST mount through claims.Chain. Public endpoints
// (/healthz, /readyz, /metrics) mount BEFORE the chain and do not require
// identity.
package api

import (
	"context"
	"net/http"

	"github.com/hanzoai/base/tools/claims"
	"github.com/hanzoai/dbx"
)

// Router wires the full HTTP mux. Takes the app handle (which owns the
// MultiTenantStore and config). Returned handler is ready to serve.
func Router(app *App) http.Handler {
	mux := http.NewServeMux()

	// Public, unauthenticated probes. Do NOT run them through claims.*
	// or we'd 503 on every healthcheck when gateway headers are absent.
	mux.HandleFunc("GET /healthz", Healthz)
	mux.HandleFunc("GET /readyz", app.Readyz)

	// Tenant-scoped routes. claims.Chain enforces the canonical
	// Strip→Inject→RequireGateway order at compile time: there's
	// nothing for a service author to mis-order.
	tenant := http.NewServeMux()
	tenant.HandleFunc("GET /v1/examples/{id}", app.ExampleGet)
	tenant.HandleFunc("POST /v1/examples", app.ExampleCreate)
	mux.Handle("/v1/", claims.Chain(tenant))

	return mux
}

// App is the service handle: config + store + hooks. Handlers live here as
// methods so they can close over the store without global state.
type App struct {
	// Store is the per-(org,user) SQLite store. Handlers call
	// Store.ForCtx(r.Context()) to obtain a *dbx.DB scoped to the caller.
	Store AppStore
}

// AppStore is the minimum surface the service's handlers require from the
// tenant store. It matches *base.MultiTenantStore.ForCtx; declared here
// so handlers depend on the interface (testable with a fake) rather than
// the concrete import. A nil AppStore is a boot error — main.go must
// wire the real store before calling Router.
type AppStore interface {
	ForCtx(ctx context.Context) (*dbx.DB, error)
}
