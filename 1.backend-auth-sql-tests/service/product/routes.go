package product

import (
	"backend-apis/types"
	"backend-apis/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string][]*types.Product{"Product Details: ": ps})

}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handleGetProduct).Methods(http.MethodGet)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {

	var payload types.Product

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product details %v", errors))
		return
	}

	_, err := h.store.CheckDuplicateProducts(&payload)

	if err != nil {
		utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("Product already exists"))
		return
	}

	err = h.store.CreateProducts(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
		CreatedAt:   payload.CreatedAt,
	})

	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	response := map[string]string{"Status": "Product Created"}

	utils.WriteJson(w, http.StatusOK, response)
}
