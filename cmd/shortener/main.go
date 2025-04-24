package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var addresses = make(map[string]string)

func getHash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha[0:8]
}

func addURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := getHash(body)
	addresses[hash] = string(body)

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.flagBaseAddr + "/" + hash))
}

func getURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := chi.URLParam(req, "hash")
	address, exists := addresses[hash]
	log.Printf("hash: %s", hash)
	log.Printf("address: %s", address)

	if !exists {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Add("Access-Control-Expose-Headers", "*")
	res.Header().Add("content-type", "text/plain")
	res.Header().Add("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func createMux() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", addURL)
	r.Get("/{hash}", getURL)

	return r
}

func main() {
	parseFlags()
	r := createMux()
	log.Fatal(http.ListenAndServe(config.flagRunAddr, r))
}
