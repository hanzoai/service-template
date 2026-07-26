// Package store wires the service's per-(org, user) SQLite storage via
// hanzoai/base/store.MultiTenantStore.
//
// Handlers never open SQLite directly. They call app.Store.ForCtx(ctx) and
// receive a *dbx.DB scoped to the caller's tenant.
package store

import (
	"fmt"
	"strings"

	basestore "github.com/hanzoai/base/store"
	"github.com/hanzoai/base/tools/filesystem"
	"github.com/hanzoai/template/config"
)

// New builds a MultiTenantStore from Config. It resolves the object-store
// URL into a filesystem.System (S3, GCS, or local file) and applies the
// pod-level LRU / idle-TTL settings.
func New(cfg config.Config) (*basestore.MultiTenantStore, error) {
	fs, err := resolveObjectStore(cfg.ObjectStoreURL)
	if err != nil {
		return nil, fmt.Errorf("store: object store %q: %w", cfg.ObjectStoreURL, err)
	}
	return basestore.New(basestore.Options{
		ObjectStore: fs,
		CacheRoot:   cfg.CacheRoot,
		LRUSize:     cfg.LRUSize,
		IdleTTL:     cfg.IdleTTL,
	})
}

// resolveObjectStore supports file:// for dev/test. S3 and GCS wiring goes
// through filesystem.NewS3 with env-driven credentials — intentionally kept
// out of this skeleton so each service explicitly owns its bucket config.
func resolveObjectStore(url string) (*filesystem.System, error) {
	switch {
	case strings.HasPrefix(url, "file://"):
		return filesystem.NewLocal(strings.TrimPrefix(url, "file://"))
	case strings.HasPrefix(url, "s3://"), strings.HasPrefix(url, "gs://"):
		return nil, fmt.Errorf("store: s3/gs backends not wired in skeleton — add NewS3 with org-owned credentials")
	default:
		return nil, fmt.Errorf("store: unsupported scheme in %q", url)
	}
}
