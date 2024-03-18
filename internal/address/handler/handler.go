package handler

import (
	"backend2/internal/address/dao"
	"backend2/internal/address/model"
	"backend2/internal/customerror"
	"backend2/internal/handlers"
	"backend2/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

const (
	addressURL = "/api/address/:uuid"
	addrURL    = "/api/address"
)

type addressHandler struct {
	logger *logging.Logger
	dao    *dao.DAO
	ctx    context.Context
}

func NewAddressHandler(logger *logging.Logger, dao *dao.DAO, ctx context.Context) handlers.Handler {
	return &addressHandler{
		logger: logger,
		dao:    dao,
		ctx:    ctx,
	}
}

func (a addressHandler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodGet, addrURL, customerror.Middleware(a.GetAll))
	r.HandlerFunc(http.MethodGet, addressURL, customerror.Middleware(a.GetOne))
	r.HandlerFunc(http.MethodPost, addrURL, customerror.Middleware(a.Create))
	r.HandlerFunc(http.MethodPut, addressURL, customerror.Middleware(a.Update))
	r.HandlerFunc(http.MethodDelete, addressURL, customerror.Middleware(a.Delete))
}

func (a addressHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	a.logger.Infof("func GetAll")
	w.Header().Set("Content-Type", "application/json")

	all, err := a.dao.FindAll(a.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	marshal, err := json.Marshal(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func (a addressHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	a.logger.Infof("func GetByID")
	id := r.URL.Query().Get("id")
	w.Header().Set("Content-Type", "application/json")

	one, err := a.dao.FindOne(a.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	marshal, err := json.Marshal(one)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func (a addressHandler) Create(w http.ResponseWriter, r *http.Request) error {
	a.logger.Infof("func Create")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This method not allowed"))
		return fmt.Errorf("%d %s method not allowed", http.StatusMethodNotAllowed, r.Method)
	}

	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := io.ReadAll(r.Body)
	var addr model.Address
	err := json.Unmarshal(reqBody, &addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	add, err := a.dao.Create(a.ctx, &addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(add)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = w.Write(marshal)
	return nil
}

func (a addressHandler) Update(w http.ResponseWriter, r *http.Request) error {
	a.logger.Infof("func Update")
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	reqBody, _ := io.ReadAll(r.Body)
	var addr *model.Address
	err := json.Unmarshal(reqBody, &addr)

	_, err = a.dao.Update(a.ctx, id, addr)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (a addressHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	a.logger.Infof("func Delete")
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	err := a.dao.Delete(a.ctx, id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	if err != nil {
		return err
	}
	return nil
}
