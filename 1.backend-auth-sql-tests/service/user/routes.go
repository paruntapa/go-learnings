package user

import (
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

}

func (h *handler) handleRegister(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}
