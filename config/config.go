// Package config loads environment-driven configuration for the template
// service. It deliberately uses only os.Getenv — no flag soup, no YAML.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config is the full runtime configuration for the service. Every field has
// a default that is safe for dev; production values are populated via env.
type Config struct {
	// ListenAddr is the HTTP bind address. Default ":8080".
	ListenAddr string

	// MetricsAddr is a SEPARATE bind for /metrics. Default ":9100".
	// Never expose /metrics on the public API mux.
	MetricsAddr string

	// ObjectStoreURL is the durable tenant-data backing. Scheme-prefixed:
	//   s3://bucket/prefix
	//   gs://bucket/prefix
	//   file:///abs/path
	ObjectStoreURL string

	// CacheRoot is the in-pod hot SQLite directory. Default "/data/cache".
	CacheRoot string

	// LRUSize caps the number of open SQLite handles per pod.
	LRUSize int

	// IdleTTL evicts handles after this long with no requests.
	IdleTTL time.Duration
}

// Load reads Config from env. Returns an error if a required value is
// missing or malformed.
func Load() (Config, error) {
	c := Config{
		ListenAddr:     envStr("LISTEN_ADDR", ":8080"),
		MetricsAddr:    envStr("METRICS_ADDR", ":9100"),
		ObjectStoreURL: os.Getenv("HANZO_TENANT_OBJECT_STORE"),
		CacheRoot:      envStr("HANZO_TENANT_LOCAL_CACHE", "/data/cache"),
		LRUSize:        envInt("HANZO_TENANT_LRU_SIZE", 1000),
		IdleTTL:        envDur("HANZO_TENANT_IDLE_TTL", 5*time.Minute),
	}
	if c.ObjectStoreURL == "" {
		return c, fmt.Errorf("config: HANZO_TENANT_OBJECT_STORE is required")
	}
	return c, nil
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func envDur(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
