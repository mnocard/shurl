package main

import (
	"net/http"

	"github.com/go-chi/chi"
	config "github.com/mnocard/shurl/internal/app/config"
	"go.uber.org/zap"
)

var addresses = make(map[string]string)
var addr config.Addr
var sugar zap.SugaredLogger

func createMux() *chi.Mux {
	r := chi.NewRouter()
	r.Use(withLogging)
	r.Post("/", addURL)
	r.Get("/{hash}", getURL)

	return r
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar = *logger.Sugar()

	config.ParseFlags(&addr)

	sugar.Infow(
		"Starting server",
		"addr", addr,
	)

	r := createMux()

	sugar.Info("mux created")

	if err := http.ListenAndServe(addr.FlagRun, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
