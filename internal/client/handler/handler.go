package handler

import (
	addrModel "backend2/internal/address/dto"
	"backend2/internal/client/db"
	"backend2/internal/client/dto"
	"backend2/internal/customerror"
	"backend2/internal/handlers"
	"backend2/pkg/logging"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	clientURL  = "/api/client"
	clientsURL = "/api/clients"
)

type clientHandler struct {
	logger *logging.Logger
	repo   *db.ClientRepository
	ctx    context.Context
}

func NewClientHandler(logger *logging.Logger, repo *db.ClientRepository, ctx context.Context) handlers.Handler {
	return &clientHandler{
		logger: logger,
		repo:   repo,
		ctx:    ctx,
	}
}

func (h *clientHandler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodGet, clientsURL, customerror.Middleware(h.GetAll))
	r.HandlerFunc(http.MethodGet, clientURL, customerror.Middleware(h.GetOne))
	r.HandlerFunc(http.MethodPost, clientsURL, customerror.Middleware(h.Create))
	r.HandlerFunc(http.MethodPut, clientURL, customerror.Middleware(h.Update))
	r.HandlerFunc(http.MethodDelete, clientURL, customerror.Middleware(h.Delete))
}

//TODO check all if/return/error and duplication of code

func (h *clientHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	h.logger.Trace("func GetAll")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	all, err := h.repo.FindAll(h.ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	clients, err := json.Marshal(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(clients)
	return nil
}

func (h *clientHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	h.logger.Trace("func GetByID")
	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")

	w.Header().Set("Content-Type", "application/json")

	client, err := h.repo.FindOne(h.ctx, name, surname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(client)
	return nil
}

func (h *clientHandler) Create(w http.ResponseWriter, r *http.Request) error {
	h.logger.Trace("func Create")

	var newClient dto.ClientDTO
	err := json.NewDecoder(r.Body).Decode(&newClient)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	err = h.repo.Create(h.ctx, &newClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *clientHandler) Update(w http.ResponseWriter, r *http.Request) error {
	h.logger.Trace("func Update")

	id := r.URL.Query().Get("id")
	var address addrModel.AddressDTO
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}
	err = h.repo.Update(h.ctx, id, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *clientHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	h.logger.Trace("func Delete")
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")

	err := h.repo.Delete(h.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
