package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

type Middleware func(http.Handler) http.Handler

const name = "github.com/Alturino/url-shortener"

var tracer = otel.Tracer(name)

func CreateStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}
		return next
	}
}
