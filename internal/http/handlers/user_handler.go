package handlers

import (
	"apiserver/internal/domain"
	"apiserver/internal/http/handlers/dto"
	"apiserver/internal/services"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var ur dto.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		slog.Error("unable decode to json", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if valid := ur.Validate(); valid == false {
		slog.Error("payload wrong", "request", ur)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := services.CreateUserInput{
		Name: ur.Name,
	}

	if err := u.userService.CreateUserFromBusinessRule(r.Context(), input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	userID, err := strconv.Atoi(strID)
	if err != nil {
		slog.Error("method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := u.userService.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Info("user not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		slog.Error("cannot delete user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (u *UserHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)

	if err != nil {
		slog.Error("unable parse URI", "URI", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := u.userService.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			slog.Error("user handler err:", "err", err)
			return
		}
		slog.Error("user handler err:", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (u *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.GetAllUserFromAnyBusinessRule(r.Context())
	if err != nil {
		slog.Error("cannot get user", "err", err)
	}

	if err := writeJSON(w, http.StatusOK, users); err != nil {
		slog.Error("unable encode to json", "err", err)
	}
}
