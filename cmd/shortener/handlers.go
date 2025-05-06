package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	hash "github.com/mnocard/shurl/internal/app/hash"
)

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

	hash := hash.GetHash(body)
	addresses[hash] = string(body)
	shortURL := addr.FlagBase + "/" + hash
	sugar.Infof("addURL. shortURL: %s, c.FlagRunAddr: %s", shortURL, addr.FlagBase)

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

	sugar.Infof("hash: %s, address: %s", hash, address)

	res.Header().Add("Access-Control-Expose-Headers", "*")
	res.Header().Add("content-type", "text/plain")
	res.Header().Add("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
