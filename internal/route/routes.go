package route

import (
	"net/http"
	"github.com/easy-comerce/backend/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user", handler.UserHandler)
	mux.HandleFunc("/order", handler.OrderHandler)
}