package restserver

import (
	"car-dealership/customer/internal"
	"net/http"
)

type Handler struct {
	application *internal.CustomerService
}

func NewHandler(application *internal.CustomerService) *Handler {
	return &Handler{
		application: application,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/customer", h.UpsertCustomer)
	mux.HandleFunc("PUT /v1/customer/{id}", h.UpsertCustomer)
	mux.HandleFunc("GET /v1/customers", h.ListCustomers)
}

func (h *Handler) UpsertCustomer(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request)  {}
