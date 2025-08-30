package main

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/route"
	"github.com/easy-comerce/backend/pkg/config"
)

func main() {
	cfg := config.Load()
	mux := http.NewServeMux()
	route.RegisterRoutes(mux)

	http.ListenAndServe(":"+cfg.Port, mux)
}
