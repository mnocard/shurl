package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecompressHandler(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Оборачиваем обработчик в middleware WithLogging
	decompressHandler := DecompressHandle(handler)

	// Создаем сжатое тело запроса
	var compressedBody strings.Builder
	gz := gzip.NewWriter(&compressedBody)
	_, err := gz.Write([]byte(strings.Repeat("Hello, world ", 20)))
	if err != nil {
		t.Fatalf("failed to write gzip data: %v", err)
	}
	gz.Close()

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(compressedBody.String()))
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/html")

	rec := httptest.NewRecorder()

	// Выполняем запрос
	decompressHandler.ServeHTTP(rec, req)

	// Проверяем, что статус ответа корректный
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Проверяем, что тело ответа корректное
	expectedBody := "Hello, World!"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestCompressHandler(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Оборачиваем обработчик в middleware WithLogging
	compressHandler := CompressHandle(handler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()

	// Выполняем запрос
	compressHandler.ServeHTTP(rec, req)

	// Проверяем, что статус ответа корректный
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Проверяем, что тело ответа корректное
	expectedBody := "Hello, World!"

	gz, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
		return
	}
	defer gz.Close()

	// Читаем распакованное тело ответа
	decompressedBody, err := io.ReadAll(gz)
	if err != nil {
		t.Fatalf("failed to read decompressed body: %v", err)
	}

	if string(decompressedBody) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
