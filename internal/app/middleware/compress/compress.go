package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	log "github.com/mnocard/shurl/internal/app/middleware/logger/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar := log.GetLogger()

		acceptEncodingHeaders := r.Header.Values("Accept-Encoding")
		isAcceptGzip := false
		for _, v := range acceptEncodingHeaders {
			if strings.Contains(v, "gzip") {
				isAcceptGzip = true
				break
			}
		}

		if !isAcceptGzip {
			sugar.Info("CompressHandle. AcceptGzip false")
			next.ServeHTTP(w, r)
			return
		}

		// Создаем буфер для записи данных
		buffer := &bytes.Buffer{}
		bufferWriter := gzipWriter{
			ResponseWriter: w,
			Writer:         buffer,
		}

		sugar.Info("CompressHandle. before next.ServeHTTP")
		next.ServeHTTP(bufferWriter, r)
		sugar.Info("CompressHandle. after next.ServeHTTP")

		// Проверяем Content-Type
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "text/html") {
			// Если Content-Type не подходит, отправляем данные как есть
			sugar.Info("CompressHandle. ContentType false")
			w.Write(buffer.Bytes())
			return
		}

		// Устанавливаем заголовок Content-Encoding
		w.Header().Set("Content-Encoding", "gzip")

		// Сжимаем данные из буфера
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			sugar.Infow("CompressHandle. NewWriterLevel", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		// Записываем сжатые данные в ResponseWriter
		_, err = gz.Write(buffer.Bytes())
		if err != nil {
			sugar.Infow("CompressHandle. gz.Write", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func DecompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar := log.GetLogger()

		contentEncodingHeaders := r.Header.Values("Content-Encoding")
		isContainsGzip := slices.Contains(contentEncodingHeaders, "gzip")
		if isContainsGzip {
			sugar.Info("DecompressHandle. isContainsGzip true")
			var reader io.ReadCloser

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			reader = gz
			r.Body = reader
		}

		sugar.Info("DecompressHandle. before next.ServeHTTP")
		next.ServeHTTP(w, r)
		sugar.Info("DecompressHandle. after next.ServeHTTP")
	})
}
