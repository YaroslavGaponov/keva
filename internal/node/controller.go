package node

import (
	"context"
	"io"
	"net/http"

	"github.com/YaroslavGaponov/keva/pkg/storage"
	"github.com/gorilla/mux"
)

type Controller struct {
	store storage.Storage
}

func ControllerNew(store storage.Storage) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	key, found := mux.Vars(r)["key"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err := c.store.Get(context.Background(), key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(value)
}

func (c *Controller) GetHash(w http.ResponseWriter, r *http.Request) {
	key, found := mux.Vars(r)["key"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err := c.store.GetHash(context.Background(), key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(value)
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	key, found := mux.Vars(r)["key"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := c.store.Set(context.Background(), key, result); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	key, found := mux.Vars(r)["key"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := c.store.Remove(context.Background(), key); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
