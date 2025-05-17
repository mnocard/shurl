package main

import (
	"net/http"

	"github.com/mnocard/shurl/internal/app/config"
	"github.com/mnocard/shurl/internal/app/handlers"
	log "github.com/mnocard/shurl/internal/app/middleware/logger/zap"
	fStorage "github.com/mnocard/shurl/internal/app/storage/file"
)

func main() {
	config := config.GetConfig()

	sugar := log.GetLogger()
	sugar.Infow(
		"Starting server",
		"config", config,
	)

	s, err := fStorage.NewFileStorage(config.FileStoragePath)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "create file storage")
	}
	defer s.Close()

	h := handlers.NewHandler(s)
	r, err := handlers.CreateMux(h)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "create mux")
	}

	sugar.Info("mux created")

	if err := http.ListenAndServe(config.FlagRun, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
