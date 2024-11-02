package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Otlp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, span := tracer.Start(r.Context(), "middleware otlp")
		defer span.End()

		logger := zerolog.Ctx(r.Context())

		logger.Info().Msg("attaching otelhttp")
		otelhttp.WithRouteTag(r.RequestURI, next)
		logger.Info().Msg("attached otelhttp")

		newR := r.WithContext(c)
		next.ServeHTTP(w, newR)
	})
}
