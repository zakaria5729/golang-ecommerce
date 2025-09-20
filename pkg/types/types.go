package types

import "net/http"

type MiddlewareHandler func(http.Handler) http.Handler

type HandlerFunc func(http.ResponseWriter, *http.Request)
