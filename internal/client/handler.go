package client

import (
	"backend2/internal/customerror"
	"backend2/internal/handlers"
	"backend2/pkg/logging"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

const (
	clientURL  = "/client/:uuid"
	clientsURL = "/clients"
)

type clientHandler struct {
	logger *logging.Logger
	srvc   *Service
}

func (c *clientHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO
}

func NewClientHandler(logger *logging.Logger, s *Service) handlers.Handler {
	return &clientHandler{
		logger: logger,
		srvc:   s,
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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	limit, _ := strconv.Atoi(vars["limit"])
	offset, _ := strconv.Atoi(vars["offset"])

	all, err := c.srvc.FindAll(context.TODO(), limit, offset)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(all)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (c *clientHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	name := vars["name"]
	surname := vars["surname"]

	w.Header().Set("Content-Type", "application/json")

	one, err := c.srvc.FindOne(context.TODO(), name, surname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	marshal, err := json.Marshal(one)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func (c *clientHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := io.ReadAll(r.Body)
	var cl Client
	err := json.Unmarshal(reqBody, &cl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	cli, err := c.srvc.Create(context.TODO(), &cl)
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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	addr := vars["address_id"]

	err := c.srvc.Update(context.TODO(), id, addr)
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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.srvc.Delete(context.TODO(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	if err != nil {
		return err
	}
	return nil
}
