package network

import (
	"time"
)

type Network struct {
	Nodes 			[]*Node
	Messages 		[]*Message
	IsNetworkDead 	bool
}

type Message struct {
	From 		int64
	To 			int64
	Content 	byte
	Delay 		int64
	Recieved 	bool
}

type Node struct {
	ID 			int64
	IsNodeDead 	bool
	Sent 		[]*Message
	Received 	[]*Message
}


func NewNetwork() *Network {
	return &Network{
		IsNetworkDead: false,
	}
}

func NewMessage(from, to, delay int64, content byte) *Message {
	return &Message{
		From: from,
		To: to,
		Content: content,
		Delay: delay,
		Recieved: false,
	}
}

func NewNode(id int64) *Node {
	return &Node{
		ID: id,
		IsNodeDead: false,
	}
}

func (network *Network) AddMessage(from, to, delay int64, content byte) {
	message := NewMessage(from, to, delay, content)
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

func (network *Network) RemoveMessage(message *Message) *Message {
	for i, m := range network.Messages {
		if m.From == message.From && m.To == message.To && m.Content == message.Content {
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

func (network *Network) SendMessage(message *Message) {
	for _, node := range network.Nodes {
		if node.ID == message.From && !node.IsNodeDead && node.ID != message.To {
			node.Sent = append(node.Sent, message)
			return
		}
	}
}

func (network *Network) ReceiveMessage(message *Message) {
	for _, node := range network.Nodes {
		if node.ID == message.To && !node.IsNodeDead {
			node.Received = append(node.Received, message)
			message.Recieved = true
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

func (network *Network) Tick(message *Message) {
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