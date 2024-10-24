package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/Alturino/url-shortener/internal/log"
)

func Logging(next http.Handler) http.Handler {
	startTime := time.Now()
	logger := log.InitLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashcode := uuid.NewString()

		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(log.KeyHashcode, hashcode).
				Any(log.KeyRequestHeader, r.Header).
				Str(log.KeyRequestHost, r.Host).
				Str(log.KeyRequestIp, r.RemoteAddr).
				Str(log.KeyRequestMethod, r.Method).
				Str(log.KeyRequestURI, r.RequestURI).
				Time(log.KeyRequestReceivedAt, startTime)
		})
		c := logger.WithContext(r.Context())
		c = context.WithValue(c, log.KeyHashcode, hashcode)
		c = context.WithValue(c, log.KeyRequestReceivedAt, startTime)
		r = r.WithContext(c)

		next.ServeHTTP(w, r)
	})
}
