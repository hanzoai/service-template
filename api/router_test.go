package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hanzoai/dbx"
)

// fakeStore lets us exercise Router without standing up a real
// MultiTenantStore + object bucket.
type fakeStore struct{}

func (fakeStore) ForCtx(_ context.Context) (*dbx.DB, error) { return nil, nil }

// TestRouter_HealthzPublic — /healthz does not require gateway identity
// and returns 200 even with no headers.
func TestRouter_HealthzPublic(t *testing.T) {
	app := &App{Store: fakeStore{}}
	h := Router(app)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("/healthz = %d, want 200", rec.Code)
	}
}

// TestRouter_ReadyzWhenStoreNil — /readyz returns 503 when the store is
// not wired. This is the one load-bearing check Readyz can make without
// depending on the service's downstream specifics.
func TestRouter_ReadyzWhenStoreNil(t *testing.T) {
	app := &App{Store: nil}
	h := Router(app)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("/readyz (nil store) = %d, want 503", rec.Code)
	}
}

// TestRouter_ReadyzWhenStoreWired — /readyz returns 200 once the store
// is attached.
func TestRouter_ReadyzWhenStoreWired(t *testing.T) {
	app := &App{Store: fakeStore{}}
	h := Router(app)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("/readyz (wired) = %d, want 200", rec.Code)
	}
}

// TestRouter_TenantRouteRejectsNoIdentity — /v1/* goes through
// claims.Chain, so a request without gateway headers is 503'd.
func TestRouter_TenantRouteRejectsNoIdentity(t *testing.T) {
	app := &App{Store: fakeStore{}}
	h := Router(app)

	req := httptest.NewRequest(http.MethodGet, "/v1/examples/42", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("/v1/examples (no identity) = %d, want 503", rec.Code)
	}
}

// TestRouter_TenantRouteStripsForgedIdentity — a request with forged
// X-User-Id / X-Org-Id is stripped by claims.Chain and then 503'd by
// RequireGateway because the stripped headers leave no identity.
func TestRouter_TenantRouteStripsForgedIdentity(t *testing.T) {
	app := &App{Store: fakeStore{}}
	h := Router(app)

	req := httptest.NewRequest(http.MethodGet, "/v1/examples/42", nil)
	req.Header.Set("X-User-Id", "attacker")
	req.Header.Set("X-Org-Id", "victim")
	req.Header.Set("X-Roles", "admin")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("/v1/examples (forged) = %d, want 503 after strip", rec.Code)
	}
}
