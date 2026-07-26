// Command appd is the long-running HTTP daemon for the template service.
//
// Boot sequence:
//
//  1. Load env config.
//  2. Assert the service is behind hanzoai/gateway (HANZO_GATEWAY_UPSTREAM=1).
//     If not, boot refuses — the service NEVER verifies JWT itself.
//  3. Build the per-(org,user) store.
//  4. Wire the mux.
//  5. Listen. On SIGTERM: flush every dirty SQLite, then exit.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hanzoai/base/tools/claims"
	"github.com/hanzoai/template/api"
	"github.com/hanzoai/template/config"
	"github.com/hanzoai/template/store"
)

func main() {
	if err := claims.AssertGatewayUpstream(); err != nil {
		log.Fatalf("boot: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("boot: %v", err)
	}

	mts, err := store.New(cfg)
	if err != nil {
		log.Fatalf("boot: %v", err)
	}

	// Wire the real per-(org,user) store. *base.MultiTenantStore
	// satisfies api.AppStore via its ForCtx(context.Context) signature.
	app := &api.App{Store: mts}

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           api.Router(app),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Graceful shutdown — flush stores before exit.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("shutdown: draining HTTP")
	shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutCtx)

	log.Printf("shutdown: flushing tenant stores")
	if err := mts.Close(shutCtx); err != nil {
		log.Printf("shutdown: store close: %v", err)
	}
	log.Printf("shutdown: done")
	os.Exit(0)
}
