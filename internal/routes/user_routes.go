package routes

import (
	"net/http"

	"apiserver/internal/http/handlers"
)

func registerUserRoutes(mux *http.ServeMux, h *handlers.UserHandler) {
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("DELETE /users/{id}", h.DeleteUser)
	mux.HandleFunc("GET /users", h.GetAll)
	mux.HandleFunc("GET /users/{id}", h.FindByID)
	mux.HandleFunc("PUT /users/{id}", h.UpdateUser)
}
