package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mnocard/shurl/internal/app/config"
	log "github.com/mnocard/shurl/internal/app/logger/zap"
	"github.com/mnocard/shurl/internal/app/storage"
)

type H struct {
	storage storage.S
}

func NewHandler(s storage.S) *H {
	return &H{
		storage: s,
	}
}

func (h *H) AddURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := h.storage.Add(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	addr := config.GetAddresses()
	shortURL := addr.FlagBase + "/" + hash

	sugar := log.GetLogger()
	sugar.Infof("addURL. shortURL: %s, c.FlagRunAddr: %s", shortURL, addr.FlagBase)

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func (h *H) GetURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := chi.URLParam(req, "hash")
	address, err := h.storage.Get(hash)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	sugar := log.GetLogger()
	sugar.Infof("hash: %s, address: %s", hash, address)

	res.Header().Add("Access-Control-Expose-Headers", "*")
	res.Header().Add("content-type", "text/plain")
	res.Header().Add("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
