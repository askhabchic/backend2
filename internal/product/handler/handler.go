package handler

import (
	"backend2/internal/customerror"
	"backend2/internal/handlers"
	"backend2/internal/product/db"
	"backend2/internal/product/dto"
	"backend2/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

const (
	productURL  = "/api/product"
	productsURL = "/api/products"
)

type productHandler struct {
	logger *logging.Logger
	repo   *db.ProductRepository
	ctx    context.Context
}

func NewProductRepository(logger *logging.Logger, repo *db.ProductRepository, ctx context.Context) handlers.Handler {
	return &productHandler{
		logger: logger,
		repo:   repo,
		ctx:    ctx,
	}
}

func (p *productHandler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodGet, productsURL, customerror.Middleware(p.GetAll))
	r.HandlerFunc(http.MethodGet, productURL, customerror.Middleware(p.GetOne))
	r.HandlerFunc(http.MethodPost, productURL, customerror.Middleware(p.Create))
	r.HandlerFunc(http.MethodPut, productURL, customerror.Middleware(p.Update))
	r.HandlerFunc(http.MethodDelete, productURL, customerror.Middleware(p.Delete))
}

func (p *productHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	all, err := p.repo.FindAll(p.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	products, err := json.Marshal(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(products)
	return nil
}

func (p *productHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	one, err := p.repo.FindOne(p.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(one)
	return nil
}

func (p *productHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var newProduct dto.ProductDTO
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	err = p.repo.Create(p.ctx, &newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (p *productHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	amount, err := strconv.Atoi(r.URL.Query().Get("amount"))
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	err = p.repo.Update(p.ctx, id, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (p *productHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("%d Invalid request body", http.StatusBadRequest)
	}

	err := p.repo.Delete(p.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
