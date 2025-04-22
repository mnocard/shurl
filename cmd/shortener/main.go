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

const BaseURI = "http://localhost:8080"
const LinkURI = "/link/"

func getHash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha[0:8]
}

func AddURL(res http.ResponseWriter, req *http.Request) {
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
	res.Write([]byte(BaseURI + LinkURI + hash))
}

func GetURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := req.URL.Path[len(LinkURI):]
	hash2 := chi.URLParam(req, "hash")
	log.Printf("hash 1: %s", hash)
	log.Printf("hash 2: %s", hash2)

	address, exists := addresses[hash2]

	if !exists {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Add("Access-Control-Expose-Headers", "*")
	res.Header().Add("content-type", "text/plain")
	res.Header().Add("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", AddURL)
		r.Route("/link", func(r chi.Router) {
			r.Get("/{hash}", GetURL)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
