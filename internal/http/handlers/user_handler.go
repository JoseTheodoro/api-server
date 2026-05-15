package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"apiserver/internal/domain"
	"apiserver/internal/http/handlers/dto"
	"apiserver/internal/services"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	rCtx := r.Context()
	rCtx, span := otel.Tracer("app/user-handler").Start(rCtx, "user.handler_create")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "create"),
	)

	var ur dto.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "decode json: payload wrong")

		slog.Error("unable decode to json", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if valid := ur.Validate(); valid == false {
		span.RecordError(errors.New("invalid user"))
		span.SetStatus(codes.Error, "invalid user payload wrong")
		slog.Error("payload wrong", "request", ur)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := services.UserCreateInput{
		FirstName: ur.FirstName,
		LastName:  ur.LastName,
		Genre:     ur.Genre,
		Email:     ur.Email,
		DateBirth: ur.DateBirth,
		Password:  ur.Password,
	}

	user, err := u.userService.CreateUserFromBusinessRule(rCtx, input)
	if err != nil {
		slog.Error("handler error on create user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	span.SetStatus(codes.Ok, "handler user success")
	writeJSON(w, http.StatusOK, user)
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

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	var userUpdateRequest dto.UserUpdateRequest
	userID, err := strconv.Atoi(strID)
	if err != nil {
		slog.Error("unable to convert path", "URI", r.URL.Path, "err", err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// create a request and validate
	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		slog.Error("unable decode payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userUpdateRequest.ID = userID

	if valid := userUpdateRequest.Validate(); valid == false {
		slog.Error("request invalid", "rquest", userUpdateRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create a update user input
	input := services.UserUpdateInput{
		ID:        userUpdateRequest.ID,
		FirstName: userUpdateRequest.FirstName,
		LastName:  userUpdateRequest.LastName,
		Email:     userUpdateRequest.Email,
		Genre:     userUpdateRequest.Genre,
		DateBirth: userUpdateRequest.DateBirth,
		Password:  userUpdateRequest.Password,
	}

	fmt.Printf("HandlernewUser=%v\n", input)

	if err := u.userService.UpdateUser(r.Context(), input); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Info("handler user not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		slog.Error("handler error update user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
