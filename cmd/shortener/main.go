package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
)

var addresses = make(map[string]string)

func getHash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha[0:8]
}

func addURL(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("content-type") != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("wrong content-type"))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("bad request"))
		return
	}

	hash := getHash(body)
	addresses[hash] = string(body)

	res.Write([]byte(hash))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
}

func getURL(res http.ResponseWriter, req *http.Request) {
	hash := req.RequestURI[1:]
	address := addresses[hash]

	res.Header().Set("content-type", "text/plain")
	res.Header().Set("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func shortenerPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		addURL(res, req)
	} else if req.Method == http.MethodGet {
		getURL(res, req)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", shortenerPage)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
