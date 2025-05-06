package main

import (
	"net/http"

	"github.com/go-chi/chi"
	config "github.com/mnocard/shurl/internal/app/config"
)

var addresses = make(map[string]string)
var addr config.Addr
var sugar = getLogger()

func createMux() *chi.Mux {
	r := chi.NewRouter()
	r.Use(withLogging)
	r.Post("/", addURL)
	r.Get("/{hash}", getURL)

	return r
}

func main() {
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
