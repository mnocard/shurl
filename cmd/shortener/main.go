package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mnocard/shurl/internal/app/config"
	"github.com/mnocard/shurl/internal/app/handlers"
	log "github.com/mnocard/shurl/internal/app/logger/zap"
	memStorage "github.com/mnocard/shurl/internal/app/storage/memoryStorage"
)

func createMux(h *handlers.H) *chi.Mux {
	r := chi.NewRouter()
	r.Use(log.WithLogging)

	r.Post("/", h.AddURL)
	r.Get("/{hash}", h.GetURL)

	return r
}

func main() {
	addr := config.GetAddresses()

	sugar := log.GetLogger()
	sugar.Infow(
		"Starting server",
		"addr", addr,
	)

	h := handlers.NewHandler(memStorage.NewMemoryStorage())
	r := createMux(h)

	sugar.Info("mux created")

	if err := http.ListenAndServe(addr.FlagRun, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
