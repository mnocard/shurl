package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	config "github.com/mnocard/shurl/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddURLHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}

	type request struct {
		url    string
		method string
		body   []byte
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "addURL correct",
			request: request{
				url:    "/",
				method: http.MethodPost,
				body:   []byte("http://ya.ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				response:    "",
			},
		},
		{
			name: "addURL wrong method",
			request: request{
				url:    "/",
				method: http.MethodGet,
				body:   []byte("http://ya.ru"),
			},
			want: want{
				contentType: "",
				statusCode:  400,
				response:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.url, bytes.NewReader(tt.request.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(AddURL)
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

func TestGetURLHandler(t *testing.T) {
	url := "http://ya.ru"

	type want struct {
		contentType string
		statusCode  int
		response    string
	}

	type request struct {
		url    string
		method string
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "getURL correct",
			request: request{
				method: http.MethodGet,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  307,
				response:    url,
			},
		},
		{
			name: "getURL wrong method",
			request: request{
				method: http.MethodPost,
			},
			want: want{
				contentType: "",
				statusCode:  405,
				response:    "",
			},
		},
	}

	log.Print("NewRouter")
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", AddURL)
		r.Route("/link", func(r chi.Router) {
			r.Get("/{hash}", GetURL)
		})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	config.ParseFlags(&addr)
	log.Print("AddURL")
	req, _ := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader([]byte(url)))
	resp, _ := http.DefaultClient.Do(req)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	shortURL := string(data)

	log.Print("ts.URL: " + ts.URL)
	log.Print("shortURL: " + shortURL)

	if addr.FlagBase != "" {
		shortURL = strings.Replace(shortURL, addr.FlagBase, ts.URL, 1)
	}

  log.Print("shortURL: " + shortURL)

	for _, tt := range tests {
		log.Print("GetURL")
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.request.method, shortURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			log.Print(string(respBody))

			require.Equal(t, tt.want.statusCode, resp.StatusCode)
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			if tt.want.statusCode != resp.StatusCode {
				require.Equal(t, url, resp.Header.Get("Location"))
			}
		})
	}
}
