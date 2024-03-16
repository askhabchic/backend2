package handler

import (
	"backend2/internal/client/dao"
	"backend2/internal/client/model"
	"backend2/internal/customerror"
	"backend2/internal/handlers"
	"backend2/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

const (
	clientURL  = "/api/client/:uuid"
	clientsURL = "/api/clients"
)

type clientHandler struct {
	logger *logging.Logger
	dao    *dao.DAO
	ctx    context.Context
}

func NewClientHandler(logger *logging.Logger, dao *dao.DAO, ctx context.Context) handlers.Handler {
	return &clientHandler{
		logger: logger,
		dao:    dao,
		ctx:    ctx,
	}
}

func (c *clientHandler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodGet, clientsURL, customerror.Middleware(c.GetAll))
	r.HandlerFunc(http.MethodGet, clientURL, customerror.Middleware(c.GetByID))
	r.HandlerFunc(http.MethodPost, clientsURL, customerror.Middleware(c.Create))
	r.HandlerFunc(http.MethodPut, clientsURL, customerror.Middleware(c.Update))
	r.HandlerFunc(http.MethodDelete, clientURL, customerror.Middleware(c.Delete))
}

func (c *clientHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	c.logger.Infof("func GetAll")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	limit, _ := strconv.Atoi(vars["limit"])
	offset, _ := strconv.Atoi(vars["offset"])

	all, err := c.dao.FindAll(c.ctx, limit, offset)
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

func (c *clientHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	c.logger.Infof("func GetByID")
	vars := mux.Vars(r)
	name := vars["name"]
	surname := vars["surname"]

	w.Header().Set("Content-Type", "application/json")

	one, err := c.dao.FindOne(c.ctx, name, surname)
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

func (c *clientHandler) Create(w http.ResponseWriter, r *http.Request) error {
	c.logger.Infof("func Create")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This method not allowed"))
		return fmt.Errorf("%d %s method not allowed", http.StatusMethodNotAllowed, r.Method)
	}

	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := io.ReadAll(r.Body)
	var cl model.Client
	err := json.Unmarshal(reqBody, &cl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	cli, err := c.dao.Create(c.ctx, &cl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(cli)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = w.Write(marshal)
	return nil
}

func (c *clientHandler) Update(w http.ResponseWriter, r *http.Request) error {
	c.logger.Infof("func Update")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	addr := vars["address_id"]

	err := c.dao.Update(c.ctx, id, addr)
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

func (c *clientHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	c.logger.Infof("func Delete")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.dao.Delete(c.ctx, id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	if err != nil {
		return err
	}
	return nil
}
