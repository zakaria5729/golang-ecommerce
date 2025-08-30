package handler

import (
	"net/http"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order endpoint"))
}
