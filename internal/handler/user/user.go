package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	model "github.com/Avito-courses/l11-examples/internal/model/user"
)

type Controller struct {
	repo userRepo
}

func NewUserController(service userRepo) *Controller {
	return &Controller{repo: service}
}

// Get обработчик получения конкретного юзера
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := c.repo.GetByID(r.Context(), id)
	if err != nil {
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, http.StatusOK, ModelToResponse(*user))
}

func (c *Controller) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (c *Controller) writeError(w http.ResponseWriter, err error) {
	message, status := c.mapError(err)
	c.writeJSON(w, status, map[string]string{"error": message})
}

func (c *Controller) mapError(err error) (string, int) {
	switch {
	case errors.Is(err, model.ErrUserNotFound):
		return "User not found", http.StatusNotFound
	case errors.Is(err, model.ErrPhoneExists):
		return "User with this phone already exists", http.StatusConflict
	default:
		return "Internal server error", http.StatusInternalServerError
	}
}
