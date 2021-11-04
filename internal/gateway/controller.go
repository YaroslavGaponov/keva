package gateway

import (
	"io/ioutil"
	"net/http"

	"github.com/YaroslavGaponov/keva/pkg/client"
	"github.com/YaroslavGaponov/keva/pkg/cluster"
	"github.com/gorilla/mux"
)

type Controller struct {
	api *client.Client
}

func ControllerNew(subscriber cluster.Subscriber) *Controller {
	return &Controller{
		api: client.NewClient(client.ClientOptionsFromEnv(), subscriber),
	}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {

	key, found := mux.Vars(r)["key"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := c.api.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if status, _ := c.api.Set(key, value); status {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {

	key, found := mux.Vars(r)["key"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if status, _ := c.api.Del(key); status {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
