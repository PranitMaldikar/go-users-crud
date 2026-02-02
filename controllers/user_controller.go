package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-users-crud/models"
	"go-users-crud/services"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(svc *services.UserService) *UserController {
	return &UserController{Service: svc}
}

func (c *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", c.usersCollection)
	mux.HandleFunc("/users/", c.usersItem)
}

func (c *UserController) usersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req models.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		u, err := c.Service.CreateUser(r.Context(), req)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, u)

	case http.MethodGet:
		users, err := c.Service.ListUsers(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "db query failed")
			return
		}
		writeJSON(w, http.StatusOK, users)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (c *UserController) usersItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDFromPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, err := c.Service.GetUser(r.Context(), id)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, u)

	case http.MethodPut:
		var req models.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		u, err := c.Service.UpdateUser(r.Context(), id, req)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, u)

	case http.MethodDelete:
		if err := c.Service.DeleteUser(r.Context(), id); err != nil {
			handleServiceError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseIDFromPath(path string) (int, bool) {
	// expects /users/{id}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] != "users" {
		return 0, false
	}
	id, err := strconv.Atoi(parts[1])
	return id, err == nil
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case services.ErrBadRequest:
		writeError(w, http.StatusBadRequest, "name and email are required")
	case services.ErrNotFound:
		writeError(w, http.StatusNotFound, "user not found")
	case services.ErrConflict:
		writeError(w, http.StatusConflict, "email already exists")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
