package network

type Network struct {
	Nodes 			*[]Node
	IsNetworkDead 	bool
}

type Node struct {
	ID 			int64
	IsNodeDead 	bool
}

func NewNetwork() *Network {
	return &Network{
		Nodes: &[]Node{},
		IsNetworkDead: false,
	}
}

func NewNode(id int64) *Node {
	return &Node{
		ID: id,
		IsNodeDead: false,
	}
}

func (network *Network) KillNetwork() {
	for _, node := range *network.Nodes {
		node.IsNodeDead = true
	}
	network.IsNetworkDead = true
}