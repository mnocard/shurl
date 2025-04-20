package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
				url:    "/",
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
				url:    "/",
				method: http.MethodPost,
			},
			want: want{
				contentType: "",
				statusCode:  400,
				response:    "",
			},
		},
	}

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(url)))
	w := httptest.NewRecorder()
	h := http.HandlerFunc(AddURL)
	h(w, r)

	result := w.Result()
	shortURL, _ := io.ReadAll(result.Body)

	fmt.Println(string(shortURL))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, string(shortURL), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetURL)
			h(w, request)

			result := w.Result()

			require.Equal(t, tt.want.statusCode, result.StatusCode)
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.statusCode != http.StatusBadRequest {
				require.Equal(t, url, result.Header.Get("Location"))
			}
		})
	}
}
