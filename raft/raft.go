package raft

import (
	"math/rand"

	"github.com/oabraham1/distributed-systems/messages"
	"github.com/oabraham1/distributed-systems/network"
)

type RaftServer struct {
	Network 	*network.Network
	Nodes 		[]*RaftNode
	LeaderID 	int64
}

type RaftNode struct {
	Node			*network.Node
	Confirmations 	int64
	GoAhead			int64
	Votes 			int64
	VotedFor 		int64
}

func NewRaftNode(id int64) *RaftNode {
	return &RaftNode{
		Node: network.NewNode(id),
		Confirmations: 0,
		Votes: 0,
		VotedFor: -1,
	}
}

func NewRaftServer(nodes int64) *RaftServer {
	raft := &RaftServer{
		Network: network.NewNetwork(),
		LeaderID: -1,
	}
	for i := 0; i < int(nodes); i++ {
		raftNode := NewRaftNode(int64(i))
		raft.Nodes = append(raft.Nodes, raftNode)
	}
	return raft
}

func (raft *RaftServer) GetNode(id int64) *RaftNode {
	for _, node := range raft.Nodes {
		if node.Node.ID == id && !node.Node.IsNodeDead {
			return node
		}
	}
	return nil
}

func (raft *RaftServer) GetLeader() *RaftNode {
	if raft.LeaderID == -1 {
		return nil
	}
	return raft.GetNode(raft.LeaderID)
}

func (raft *RaftServer) ReceiveMessage(message *messages.Message) {
	reciever := raft.GetNode(message.To)
	if reciever == nil {
		return
	}
	sender := raft.GetNode(message.From)
	if sender == nil {
		return
	}

	reciever.Node.Received = append(reciever.Node.Received, message)
	
	// Handle message for different types of messages
	if message.Type == messages.STARTLEADERELECTION {
		if reciever.GoAhead == 0 {
			reciever.GoAhead = 1
			newMessage := messages.NewMessage(rand.Int63n(30), reciever.Node.ID, sender.Node.ID, rand.Int63n(12), messages.APPROVELEADERELCTION)
			raft.Network.Messages = append(raft.Network.Messages, newMessage)
			reciever.Node.Sent = append(reciever.Node.Sent, message)
		}

	} 	else if message.Type == messages.APPROVELEADERELCTION {
		reciever.GoAhead += 1

	}	else if message.Type == messages.REQUESTVOTE {
		if reciever.VotedFor == -1 {
			reciever.VotedFor = sender.Node.ID
			newMessage := messages.NewMessage(rand.Int63n(30), reciever.Node.ID, sender.Node.ID, rand.Int63n(12), messages.VOTE)
			raft.Network.Messages = append(raft.Network.Messages, newMessage)
			reciever.Node.Sent = append(reciever.Node.Sent, message)
		}

	}	else if message.Type == messages.VOTE {
		reciever.Votes += 1

	}	else if message.Type == messages.LEADERCONFIRMATIONREQUEST {
		if sender.Votes > int64(len(raft.Nodes)/2) {
			newMessage := messages.NewMessage(rand.Int63n(30), reciever.Node.ID, sender.Node.ID, rand.Int63n(12), messages.CONFIRMLEADER)
			raft.Network.Messages = append(raft.Network.Messages, newMessage)
			reciever.Node.Sent = append(reciever.Node.Sent, message)
		}

	}	else if message.Type == messages.CONFIRMLEADER {
		reciever.Confirmations += 1
	}
}

func (raft *RaftServer) StreamMessages() {
	for len(raft.Network.Messages) > int(0) {
		for _, message := range raft.Network.Messages {
			if message.Delay == 0 {
				m := raft.Network.RemoveMessage(message)
				if m != nil {
					raft.Network.SendMessage(m)
					raft.ReceiveMessage(m)
				}
			}
			raft.Network.Tick(message)
		}
	}
}

func (raft *RaftServer) StartElection(candidate *RaftNode) {
	if raft.GetLeader() != nil {
		return
	}

	// Send Message to all nodes to check if election can be started
	for _, node := range raft.Nodes {
		if node.Node.ID != candidate.Node.ID {
			raft.Network.AddMessage(1, candidate.Node.ID, node.Node.ID, rand.Int63n(12), messages.STARTLEADERELECTION)
		}
	}

	raft.StreamMessages()

	// If majority of nodes confirm election can be started then send request vote to all nodes
	if candidate.GoAhead > int64(len(raft.Nodes)/2) {
		// Candidate votes for itself
		candidate.VotedFor = candidate.Node.ID
		candidate.Votes += 1

		// Send Message to all nodes to vote for candidate
		for _, node := range raft.Nodes {
			if node.Node.ID != candidate.Node.ID {
				raft.Network.AddMessage(1, candidate.Node.ID, node.Node.ID, rand.Int63n(12), messages.REQUESTVOTE)
			}
		}
	}

	raft.StreamMessages()

	// Check if candidate has majority votes
	if candidate.Votes > int64(len(raft.Nodes)/2) {
		raft.LeaderID = candidate.Node.ID
	}

	// raft.StreamMessages()
		

	// // Check if node got all confirmations
	// if candidate.Confirmations > int64(len(raft.Nodes)/2) {
	// 	raft.LeaderID = candidate.Node.ID
	// }

	// Reset votes
	for _, node := range raft.Nodes {
		node.Votes = 0
		node.VotedFor = -1
	}


	// Reset confirmations
	for _, node := range raft.Nodes {
		node.Confirmations = 0
	}
}
