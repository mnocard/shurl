package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
)

var addresses = make(map[string]string)

const baseUri = "http://localhost:8080"
const linkUri = "/link/"

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
	res.Write([]byte(baseUri + linkUri + hash))
}

func GetURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := req.URL.Path[len(linkUri):]
	address, exists := addresses[hash]

	if !exists {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Add("Access-Control-Expose-Headers", "Authorization")
	res.Header().Set("content-type", "text/plain")
	res.Header().Add("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", AddURL)
	mux.HandleFunc(linkUri, GetURL)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
