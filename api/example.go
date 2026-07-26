package api

import "net/http"

// ExampleGet is a placeholder tenant-scoped handler. Replace with your
// domain logic. The middleware chain guarantees:
//
//	claims.FromContext(ctx).OrgID  != ""
//	claims.FromContext(ctx).UserID != ""
//
// So handlers never need to nil-check identity — the chain already 503'd
// the caller if headers were missing.
func (*App) ExampleGet(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ExampleCreate is the write counterpart.
func (*App) ExampleCreate(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
