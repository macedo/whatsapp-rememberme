package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const headerXRequestID = "X-Request-ID"

type contextkey string

const RequestIDCtxKey = contextkey("requestid")

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(headerXRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDCtxKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
