package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
	Register(router *httprouter.Router)
	GetAll(w http.ResponseWriter, r *http.Request) error
	GetByID(w http.ResponseWriter, r *http.Request) error
	Create(w http.ResponseWriter, r *http.Request) error
	Update(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
}

func NewHandler() {

}
