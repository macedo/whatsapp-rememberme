package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type wrapResponseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *wrapResponseWriter) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *wrapResponseWriter) Status() int {
	return rw.code
}

func RequestLogger(log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestID := ctx.Value(RequestIDCtxKey).(string)

			log.WithFields(logrus.Fields{
				"request_id":     requestID,
				"method":         r.Method,
				"remote-addr":    r.RemoteAddr,
				"http-protocl":   r.Proto,
				"headers":        r.Header,
				"content-length": r.ContentLength,
			}).Infof("HTTP request to %s", r.URL)

			ww := &wrapResponseWriter{w, 200}

			t1 := time.Now()
			defer func() {
				log.WithFields(logrus.Fields{
					"request_id": requestID,
					"status":     ww.Status(),
					"headers":    ww.Header(),
					"duration":   time.Since(t1),
				}).Infof("completed %s", r.URL)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
