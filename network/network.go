package network

import (
	"time"
	"github.com/oabraham1/distributed-systems/messages"
)

type Network struct {
	Nodes 			[]*Node
	Messages 		[]*messages.Message
	IsNetworkDead 	bool
}

type Node struct {
	ID 			int64
	IsNodeDead 	bool
	Sent 		[]*messages.Message
	Received 	[]*messages.Message
}

func NewNetwork() *Network {
	return &Network{
		IsNetworkDead: false,
	}
}

func NewNode(id int64) *Node {
	return &Node{
		ID: id,
		IsNodeDead: false,
	}
}

func (network *Network) AddMessage(id, from, to, delay int64, messageType messages.MessageType) {
	message := messages.NewMessage(id, from, to, delay, messageType)
	network.Messages = append(network.Messages, message)
}

func (network *Network) AddNode(node *Node) {
	network.Nodes = append(network.Nodes, node)
}

func (network *Network) GetNode(id int64) *Node {
	for _, node := range network.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

func (network *Network) RemoveMessage(message *messages.Message) *messages.Message {
	for i, m := range network.Messages {
		if m.String() == message.String() {
			network.Messages = append((network.Messages)[:i], (network.Messages)[i+1:]...)
			return m
		}
	}
	
	return nil
}

func (network *Network) RemoveNode(id int64) {
	for i, n := range network.Nodes {
		if n.ID == id {
			network.Nodes = append((network.Nodes)[:i], (network.Nodes)[i+1:]...)
			n.KillNode()
		}
	}
}

func (network *Network) SendMessage(message *messages.Message) {
	for _, node := range network.Nodes {
		if node.ID == message.From && !node.IsNodeDead && node.ID != message.To {
			node.Sent = append(node.Sent, message)
			return
		}
	}
}

func (network *Network) ReceiveMessage(message *messages.Message) {
	for _, node := range network.Nodes {
		if node.ID == message.To && !node.IsNodeDead {
			node.Received = append(node.Received, message)
			message.Recieved = 1
			return
		}
	}
}

func (network *Network) StreamMessages() {
	for len(network.Messages) > int(0) {
		for _, message := range network.Messages {
			if message.Delay == 0 {
				m := network.RemoveMessage(message)
				if m != nil {
					network.SendMessage(m)
					network.ReceiveMessage(m)
				}
			}
			network.Tick(message)
		}
	}
}

func (network *Network) Tick(message *messages.Message) {
	message.Delay--
	time.Sleep(time.Second)
}

func (node *Node) KillNode() {
	node.IsNodeDead = true
}

func (network *Network) KillNetwork() {
	for _, node := range network.Nodes {
		node.KillNode()
	}
	network.IsNetworkDead = true
}