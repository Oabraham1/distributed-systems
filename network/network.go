package network

type Network struct {
	nodes 			int64
	IsNetworkDead 	bool
}

func NewNetwork(nodes int64) *Network {
	return &Network{
		nodes: nodes,
		IsNetworkDead: false,
	}
}

func (n *Network) KillNetwork() {
	n.IsNetworkDead = true
}