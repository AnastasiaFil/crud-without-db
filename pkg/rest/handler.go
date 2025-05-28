package rest

import (
	"crud-without-db/internal/domain"
	"crud-without-db/pkg/logger"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Users interface {
	Create(user domain.User) error
	GetByID(id int64) (domain.User, error)
	GetAll() ([]domain.User, error)
	Delete(id int64) error
	Update(id int64, inp domain.User) error
}

type Handler struct {
	usersService Users
	logger       zerolog.Logger
}

func NewHandler(users Users) *Handler {
	return &Handler{
		usersService: users,
		logger:       logger.GetLogger("handler"),
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	users := r.PathPrefix("/users").Subrouter()
	{
		users.HandleFunc("", h.createUser).Methods("POST")
		users.HandleFunc("", h.getAllUsers).Methods("GET")
		users.HandleFunc("/{id}", h.getUserByID).Methods("GET")
		users.HandleFunc("/{id}", h.deleteUser).Methods("DELETE")
		users.HandleFunc("/{id}", h.updateUser).Methods("PUT")
	}

	return r
}

// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User
// @Router /users/{id} [get]
func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		h.logger.Error().Err(err).Str("method", "getUserByID").Msg("Invalid user ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debug().Int64("user_id", id).Msg("Getting user by ID")

	user, err := h.usersService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.logger.Warn().Int64("user_id", id).Msg("User not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.logger.Error().Err(err).Int64("user_id", id).Msg("Failed to get user by ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		h.logger.Error().Err(err).Int64("user_id", id).Msg("Failed to marshal user response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int64("user_id", id).Str("user_name", user.Name).Msg("User retrieved successfully")
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "Create user"
// @Success 201 {object} domain.User
// @Router /users [post]
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error().Err(err).Str("method", "createUser").Msg("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user domain.User
	if err = json.Unmarshal(reqBytes, &user); err != nil {
		h.logger.Error().Err(err).Str("method", "createUser").Msg("Failed to unmarshal user data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debug().
		Str("user_name", user.Name).
		Int("user_age", user.Age).
		Str("user_sex", user.Sex).
		Msg("Creating new user")

	err = h.usersService.Create(user)
	if err != nil {
		h.logger.Error().Err(err).
			Str("user_name", user.Name).
			Msg("Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info().Str("user_name", user.Name).Msg("User created successfully")
	w.WriteHeader(http.StatusCreated)
}

// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		h.logger.Error().Err(err).Str("method", "deleteUser").Msg("Invalid user ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debug().Int64("user_id", id).Msg("Deleting user")

	err = h.usersService.Delete(id)
	if err != nil {
		h.logger.Error().Err(err).Int64("user_id", id).Msg("Failed to delete user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int64("user_id", id).Msg("User deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} domain.User
// @Router /users [get]
func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().Msg("Getting all users")

	users, err := h.usersService.GetAll()
	if err != nil {
		h.logger.Error().Err(err).Str("method", "getAllUsers").Msg("Failed to get all users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		h.logger.Error().Err(err).Str("method", "getAllUsers").Msg("Failed to marshal users response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int("users_count", len(users)).Msg("Retrieved all users successfully")
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Update a user
// @Description Update a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "Update user"
// @Success 200 {object} domain.User
// @Router /users/{id} [put]
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		h.logger.Error().Err(err).Str("method", "updateUser").Msg("Invalid user ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error().Err(err).Int64("user_id", id).Msg("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.User
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		h.logger.Error().Err(err).Int64("user_id", id).Msg("Failed to unmarshal user data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debug().
		Int64("user_id", id).
		Str("user_name", inp.Name).
		Int("user_age", inp.Age).
		Str("user_sex", inp.Sex).
		Msg("Updating user")

	err = h.usersService.Update(id, inp)
	if err != nil {
		h.logger.Error().Err(err).Int64("user_id", id).Msg("Failed to update user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int64("user_id", id).Str("user_name", inp.Name).Msg("User updated successfully")
	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
