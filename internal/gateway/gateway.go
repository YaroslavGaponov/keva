package gateway

import (
	"net/http"

	"github.com/YaroslavGaponov/keva/pkg/cluster"
	"github.com/YaroslavGaponov/keva/pkg/utils"
	"github.com/gorilla/mux"
)

const (
	PATH_STORAGE = "/storage"
)

type Gateway struct {
	subscriber cluster.Subscriber
	server  *http.Server
}

func New(subscriber cluster.Subscriber) *Gateway {

	controller := ControllerNew(subscriber)

	router := mux.NewRouter()
	router.HandleFunc(PATH_STORAGE+"/{key}", controller.Get).Methods("GET")
	router.HandleFunc(PATH_STORAGE+"/{key}", controller.Set).Methods("POST", "PUT")
	router.HandleFunc(PATH_STORAGE+"/{key}", controller.Del).Methods("DELETE")

	port := utils.GetEnvVariableOrDefult("PORT", "5555")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return &Gateway{
		subscriber: subscriber,
		server:  server,
	}
}

func (m *Gateway) Start() error {
	return m.server.ListenAndServe()
}
