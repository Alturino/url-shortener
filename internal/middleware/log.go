package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/Alturino/url-shortener/internal/log"
)

func Logging(next http.Handler) http.Handler {
	logger := log.InitLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashcode := uuid.NewString()

		logger.Info().Msg("attaching request value to logger")
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(log.KeyHashcode, hashcode).
				Any(log.KeyRequestHeader, r.Header).
				Str(log.KeyRequestHost, r.Host).
				Str(log.KeyRequestIp, r.RemoteAddr).
				Str(log.KeyRequestMethod, r.Method).
				Str(log.KeyRequestURI, r.RequestURI).
				Str(log.KeyRequestURL, r.URL.String())
		})
		logger.Info().Msg("attached request value to logger")

		logger.Info().Msg("attaching request value to context")
		c := log.AttachHashcodeToContext(r.Context(), hashcode)
		c = logger.WithContext(c)
		newR := r.WithContext(c)
		logger.Info().Msg("attached request value to context")

		logger.Info().Msg("next handler")
		next.ServeHTTP(w, newR)
	})
}
