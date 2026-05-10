package routes

import (
	"apiserver/internal/http/handlers"
	"net/http"
)

type Handlers struct {
	User    *handlers.UserHandler
	Product *handlers.ProductHandler
}

func Register(mux *http.ServeMux, h Handlers) {
	registerUserRoutes(mux, h.User)
	registerProductRoutes(mux, h.Product)
}
