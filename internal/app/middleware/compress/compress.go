package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
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
		acceptEncodingHeaders := r.Header.Values("Accept-Encoding")
		contentHeaders := w.Header().Values("Content-Type")
		isAcceptGzip := slices.Contains(acceptEncodingHeaders, "gzip")
		isContent := slices.Contains(contentHeaders, "application/json") || slices.Contains(contentHeaders, "text/html")

		if !isAcceptGzip || !isContent {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func DecompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncodingHeaders := r.Header.Values("Content-Encoding")
		isContainsGzip := slices.Contains(contentEncodingHeaders, "gzip")
		if isContainsGzip {
			var reader io.ReadCloser

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
			r.Body = reader
		}

		next.ServeHTTP(w, r)
	})
}
