package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mnocard/shurl/internal/app/config"
	"github.com/mnocard/shurl/internal/app/hash"
	memStorage "github.com/mnocard/shurl/internal/app/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIShorten(t *testing.T) {
	url := "http://ya.ru"
	hash := hash.GetHash([]byte(url))
	addr := config.GetAddresses()
	result := addr.FlagBase + "/" + hash

	type want struct {
		contentType string
		statusCode  int
		response    Response
	}

	type request struct {
		url    string
		method string
		body   Request
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "apiShorten correct",
			request: request{
				url:    "/api/shorten",
				method: http.MethodPost,
				body:   Request{URL: url},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				response:    Response{Result: result},
			},
		},
		{
			name: "apiShorten wrong method",
			request: request{
				url:    "/api/shorten",
				method: http.MethodGet,
				body:   Request{},
			},
			want: want{
				contentType: "",
				statusCode:  400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request.body)
			request := httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler := NewHandler(memStorage.NewMemoryStorage())
			h := http.HandlerFunc(handler.APIShorten)
			h(w, request)

			result := w.Result()

			require.Equal(t, tt.want.statusCode, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.statusCode != http.StatusBadRequest {
				shortURL, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.NotEmpty(t, shortURL)
			}
		})
	}
}
