package middleware

import (
	"net/http"
	"time"
	"log/slog"
)

type responseData struct {
	status	int
	size	int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	*responseData
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.responseData.status = statusCode

	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	if err == nil {
		lrw.responseData.size += size
	}

	if lrw.responseData.status == 0 {
		lrw.responseData.status = http.StatusOK
	}

	return size, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		resData := &responseData{
			status: 0,
			size: 0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData: resData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		slog.Info("recieved request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", resData.status),
		slog.Int("bytes", resData.size),
		slog.String("ip", r.RemoteAddr),
		slog.Duration("duration", duration),
	)
	})
}