package cluster

type Publisher interface {
	RegisterNode(info NodeInfo) error
}
