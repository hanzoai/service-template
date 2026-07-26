package api

import "net/http"

// Healthz is a liveness probe: 200 as long as the process is up. Do not
// add dependency checks here — kubelet's liveness probe MUST NOT fail
// during a transient downstream outage. That's what Readyz is for.
func Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true}`))
}

// Readyz is a readiness probe: 200 only when the tenant store is wired
// and the process is ready to serve traffic. We verify the single
// invariant this template actually owns — app.Store is not nil — and
// leave deeper store health (object-storage ping, IAM reach) to
// callers who override Readyz. The probe is intentionally cheap:
// kubelet calls it every second.
func (a *App) Readyz(w http.ResponseWriter, _ *http.Request) {
	if a == nil || a.Store == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"ready":false,"reason":"store not wired"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ready":true}`))
}
