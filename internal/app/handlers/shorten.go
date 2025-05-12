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

func (h *H) APIShorten(res http.ResponseWriter, req *http.Request) {
	sugar := log.GetLogger()
	sugar.Info("APIShorten. Start")
	if req.Method != http.MethodPost {
		sugar.Errorw("APIShorten. Method error", "Method", req.Method)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var request Request
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		sugar.Error("APIShorten. Read body error")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		sugar.Errorw("APIShorten. Unmarshal error", "error", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := h.storage.Add(string(request.URL))
	if err != nil {
		sugar.Errorw("APIShorten. Get hash error", "error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	addr := config.GetAddresses()
	response := Response{
		Result: addr.FlagBase + "/" + hash,
	}

	data, err := json.Marshal(response)
	if err != nil {
		sugar.Errorw("APIShorten. Marshal error", "error", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(data)
	sugar.Info("APIShorten. End")
}
