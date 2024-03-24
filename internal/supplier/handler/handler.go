package handler

import (
	addrModel "backend2/internal/address/dto"
	"backend2/internal/customerror"
	"backend2/internal/handlers"
	"backend2/internal/supplier/db"
	"backend2/internal/supplier/dto"
	"backend2/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	supplierURL  = "/api/supplier"
	suppliersURL = "/api/suppliers"
)

type supplierHandler struct {
	logger *logging.Logger
	repo   *db.SupplierRepository
	ctx    context.Context
}

func NewSupplierHandler(logger *logging.Logger, repo *db.SupplierRepository, ctx context.Context) handlers.Handler {
	return &supplierHandler{
		logger: logger,
		repo:   repo,
		ctx:    ctx,
	}
}

func (s *supplierHandler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodGet, suppliersURL, customerror.Middleware(s.GetAll))
	r.HandlerFunc(http.MethodGet, supplierURL, customerror.Middleware(s.GetOne))
	r.HandlerFunc(http.MethodPost, supplierURL, customerror.Middleware(s.Create))
	r.HandlerFunc(http.MethodPut, supplierURL, customerror.Middleware(s.Update))
	r.HandlerFunc(http.MethodDelete, supplierURL, customerror.Middleware(s.Delete))
}

func (s *supplierHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	s.logger.Tracef("func GetAll")

	all, err := s.repo.FindAll(s.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	suppliers, err := json.Marshal(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(suppliers)
	return nil
}

func (s *supplierHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	s.logger.Tracef("func GetAll")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	one, err := s.repo.FindOne(s.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(one)
	return nil
}

func (s *supplierHandler) Create(w http.ResponseWriter, r *http.Request) error {
	s.logger.Tracef("func Create")

	var newSupplier dto.SupplierDTO
	err := json.NewDecoder(r.Body).Decode(&newSupplier)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	err = s.repo.Create(s.ctx, &newSupplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *supplierHandler) Update(w http.ResponseWriter, r *http.Request) error {
	s.logger.Tracef("func Update")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	var address addrModel.AddressDTO
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	err = s.repo.Update(s.ctx, id, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *supplierHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	s.logger.Tracef("func Delete")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	err := s.repo.Delete(s.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return nil
}
