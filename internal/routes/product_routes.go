package routes

import (
	"net/http"

	"apiserver/internal/http/handlers"
)

func registerProductRoutes(mux *http.ServeMux, h *handlers.ProductHandler) {
	mux.HandleFunc("/product/{id}", h.GetProduct)
}
