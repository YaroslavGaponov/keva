package node

import (
	"net"
	"net/http"

	"github.com/YaroslavGaponov/keva/pkg/cluster"
	"github.com/YaroslavGaponov/keva/pkg/storage"
	"github.com/YaroslavGaponov/keva/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	PATH_STORAGE = "storage"
)

type Node struct {
	info cluster.NodeInfo

	server   *http.Server
	listener net.Listener

	publisher cluster.Publisher
	storage   storage.Storage
}

func New(publisher cluster.Publisher, storage storage.Storage) (*Node, error) {

	controller := ControllerNew(storage)

	listener, err := utils.GetTcpListener(utils.SelfIP().String())
	if err != nil {
		return nil, err
	}

	info := cluster.NodeInfo{
		NodeId: uuid.New().String(),
		Schema: "http",
		Host:   listener.Addr().(*net.TCPAddr).IP.String(),
		Port:   listener.Addr().(*net.TCPAddr).Port,
		Path:   PATH_STORAGE,
	}

	router := mux.NewRouter()
	router.HandleFunc("/"+PATH_STORAGE+"/{key}/hash", controller.GetHash).Methods("GET")
	router.HandleFunc("/"+PATH_STORAGE+"/{key}", controller.Get).Methods("GET")
	router.HandleFunc("/"+PATH_STORAGE+"/{key}", controller.Set).Methods("POST", "PUT")
	router.HandleFunc("/"+PATH_STORAGE+"/{key}", controller.Del).Methods("DELETE")

	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: router,
	}

	n := Node{
		info:      info,
		publisher: publisher,
		storage:   storage,
		server:    server,
		listener:  listener,
	}

	return &n, nil
}

func (n *Node) Start() error {
	if err := n.publisher.RegisterNode(n.info); err != nil {
		return err
	}

	return n.server.Serve(n.listener)
}

func (n *Node) Stop() error {
	if err := n.storage.Close(); err != nil {
		return err
	}
	return n.server.Close()
}
