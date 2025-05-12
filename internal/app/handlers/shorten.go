package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mnocard/shurl/internal/app/config"
	log "github.com/mnocard/shurl/internal/app/logger/zap"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func (h *H) ApiShorten(res http.ResponseWriter, req *http.Request) {
	sugar := log.GetLogger()
	sugar.Info("ApiShorten. Start")
	if req.Method != http.MethodPost {
		sugar.Errorw("ApiShorten. Method error", "Method", req.Method)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var request Request
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		sugar.Error("ApiShorten. Read body error")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		sugar.Errorw("ApiShorten. Unmarshal error", "error", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := h.storage.Add(string(request.URL))
	if err != nil {
		sugar.Errorw("ApiShorten. Get hash error", "error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	addr := config.GetAddresses()
	response := Response{
		Result: addr.FlagBase + "/" + hash,
	}

	data, err := json.Marshal(response)
	if err != nil {
		sugar.Errorw("ApiShorten. Marshal error", "error", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write(data)
	sugar.Info("ApiShorten. End")
}
