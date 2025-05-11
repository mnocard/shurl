package main

import (
	"net/http"

	"github.com/mnocard/shurl/internal/app/config"
	"github.com/mnocard/shurl/internal/app/handlers"
	log "github.com/mnocard/shurl/internal/app/logger/zap"
	memStorage "github.com/mnocard/shurl/internal/app/storage/memory"
)

func main() {
	addr := config.GetAddresses()

	sugar := log.GetLogger()
	sugar.Infow(
		"Starting server",
		"addr", addr,
	)

	s := memStorage.NewMemoryStorage()
	h := handlers.NewHandler(s)
	r, err := handlers.CreateMux(h)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "create mux")
	}

	sugar.Info("mux created")

	if err := http.ListenAndServe(addr.FlagRun, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
