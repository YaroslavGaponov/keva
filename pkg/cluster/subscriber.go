package cluster

type Subscriber interface {
	GetNodes() (map[string]NodeInfo, error)
}
