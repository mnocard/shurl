package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/mnocard/shurl/internal/app"
)

var addresses = make(map[string]string)
var addr config.Addr

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
	shortURL := addr.FlagBase + "/" + hash
	log.Printf("addURL. shortURL: %s, c.FlagRunAddr: %s", shortURL, addr.FlagBase)

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func getURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := chi.URLParam(req, "hash")
	address, exists := addresses[hash]
	if !exists {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("hash: %s, address: %s", hash, address)

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
	log.Print("main start")
	config.ParseFlags(&addr)
	log.Printf("main parseFlags. config.flagRunAddr: %s, config.flagBaseAddr: %s", addr.FlagRun, addr.FlagBase)
	r := createMux()
	log.Print("main createMux")
	log.Fatal(http.ListenAndServe(addr.FlagRun, r))
}
