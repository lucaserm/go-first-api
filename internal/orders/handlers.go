package orders

import (
	"log"
	"net/http"

	"github.com/lucaserm/ecom/internal/json"
	"github.com/lucaserm/ecom/internal/products"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var payload createOrderParams

	if err := json.Read(r, &payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), payload)

	if err != nil {
		log.Println(err)

		if err == products.ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
