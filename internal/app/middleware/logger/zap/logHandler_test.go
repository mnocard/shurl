package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetLogger(t *testing.T) {
	// Проверяем, что GetLogger возвращает не nil
	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger() returned nil")
	}

	// Проверяем, что повторный вызов возвращает тот же экземпляр
	logger2 := GetLogger()
	if logger != logger2 {
		t.Fatal("GetLogger() did not return the same instance")
	}
}

func TestWithLogging(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Инициализируем логгер
	GetLogger()
	// Оборачиваем обработчик в middleware WithLogging
	loggedHandler := WithLogging(handler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Выполняем запрос
	loggedHandler.ServeHTTP(rec, req)

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

func TestLoggingResponseWriter(t *testing.T) {
	// Создаем буфер для записи ответа
	writer := httptest.NewRecorder()
	lrw := loggingResponseWriter{
		ResponseWriter: writer,
		responseData:   &responseData{},
	}

	// Пишем данные в loggingResponseWriter
	data := []byte("test data")
	n, err := lrw.Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что количество записанных байтов корректное
	if n != len(data) {
		t.Errorf("expected %d bytes written, got %d", len(data), n)
	}

	// Проверяем, что размер записанных данных обновился
	if lrw.responseData.size != len(data) {
		t.Errorf("expected response size %d, got %d", len(data), lrw.responseData.size)
	}

	// Проверяем WriteHeader
	statusCode := http.StatusCreated
	lrw.WriteHeader(statusCode)
	if lrw.responseData.status != statusCode {
		t.Errorf("expected status %d, got %d", statusCode, lrw.responseData.status)
	}
}

func TestWithLoggingLogs(t *testing.T) {
	// Создаем буфер для записи логов
	var logBuffer bytes.Buffer

	// Настраиваем тестовый логгер
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&logBuffer),
		zapcore.InfoLevel,
	)
	testLogger := zap.New(core).Sugar()

	// Заменяем глобальный логгер на тестовый
	sugar = testLogger

	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Оборачиваем обработчик в middleware WithLogging
	loggedHandler := WithLogging(handler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Выполняем запрос
	loggedHandler.ServeHTTP(rec, req)

	// Проверяем, что логи записаны
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, `uri /test`) {
		t.Errorf("expected log to contain URI, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `method GET`) {
		t.Errorf("expected log to contain method, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `status 200`) {
		t.Errorf("expected log to contain status, got: %s", logOutput)
	}
}
