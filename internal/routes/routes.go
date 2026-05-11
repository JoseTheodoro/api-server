package routes

import (
	"net/http"

	"apiserver/internal/http/handlers"
)

type Handlers struct {
	User    *handlers.UserHandler
	Product *handlers.ProductHandler
}

func Register(mux *http.ServeMux, h Handlers) {
	registerUserRoutes(mux, h.User)
	registerProductRoutes(mux, h.Product)
}
