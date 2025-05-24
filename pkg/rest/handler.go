package rest

import (
	"crud-without-db/internal/domain"
	"encoding/json"
	"errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
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
}

func NewHandler(users Users) *Handler {
	return &Handler{
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (not recommended for production)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	))

	users := r.PathPrefix("/users").Subrouter()
	{
		users.HandleFunc("", h.createUser).Methods(http.MethodPost)
		users.HandleFunc("", h.getAllUsers).Methods(http.MethodGet)
		users.HandleFunc("/{id}", h.getUserByID).Methods(http.MethodGet)
		users.HandleFunc("/{id}", h.deleteUser).Methods(http.MethodDelete)
		users.HandleFunc("/{id}", h.updateUser).Methods(http.MethodPut)
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
		log.Println("getUserByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.usersService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Println("getUserByID() StatusBadRequest error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getUserByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Println("getUserByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		log.Println("createUser() readAll error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user domain.User
	if err = json.Unmarshal(reqBytes, &user); err != nil {
		log.Println("createUser() unmarshal error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usersService.Create(user)
	if err != nil {
		log.Println("createUser() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		log.Println("deleteUser() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usersService.Delete(id)
	if err != nil {
		log.Println("deleteUser() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} domain.User
// @Router /users [get]
func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usersService.GetAll()
	if err != nil {
		log.Println("getAllUsers() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		log.Println("getAllUsers() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		log.Println("error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("updateUser() readAll error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.User
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		log.Println("updateUser() unmarshal error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usersService.Update(id, inp)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
