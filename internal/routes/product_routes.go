package routes

import (
	"apiserver/internal/http/handlers"
	"net/http"
)

func registerProductRoutes(mux *http.ServeMux, h *handlers.ProductHandler) {
	mux.HandleFunc("/product/{id}", h.GetProduct)
}
