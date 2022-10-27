package raft

import (
	"testing"
	
	"github.com/stretchr/testify/require"
)

func TestElection(t *testing.T) {
	// Create 5 raft nodes
	nodeOne := NewRaft(1, []*Raft{})
	nodeTwo := NewRaft(2, []*Raft{})
	nodeThree := NewRaft(3, []*Raft{})
	nodeFour := NewRaft(4, []*Raft{})
	nodeFive := NewRaft(5, []*Raft{})

	require.Equal(t, nodeOne.me, int64(1))
	require.Equal(t, nodeTwo.me, int64(2))
	require.Equal(t, nodeThree.me, int64(3))
	require.Equal(t, nodeFour.me, int64(4))
	require.Equal(t, nodeFive.me, int64(5))

	// Add the nodes to each other's peers
	// TODO: Add a method to add a peer to a node
	nodeOne.peers = []*Raft{nodeTwo, nodeThree, nodeFour, nodeFive}
	nodeTwo.peers = []*Raft{nodeOne, nodeThree, nodeFour, nodeFive}
	nodeThree.peers = []*Raft{nodeOne, nodeTwo, nodeFour, nodeFive}
	nodeFour.peers = []*Raft{nodeOne, nodeTwo, nodeThree, nodeFive}
	nodeFive.peers = []*Raft{nodeOne, nodeTwo, nodeThree, nodeFour}

	// TODO: Refractor into loop
	require.Equal(t, len(nodeOne.peers), 4)
	require.Equal(t, len(nodeTwo.peers), 4)
	require.Equal(t, len(nodeThree.peers),4)
	require.Equal(t, len(nodeFour.peers), 4)
	require.Equal(t, len(nodeFive.peers), 4)

	// TODO: Refractor into loop
	// No leader should be elected
	require.Nil(t, nodeOne.state.leader)
	require.Nil(t, nodeTwo.state.leader)
	require.Nil(t, nodeThree.state.leader)
	require.Nil(t, nodeFour.state.leader)
	require.Nil(t, nodeFive.state.leader)

	// Kill the cluster
	// TODO: Refractor into loop
	nodeOne.KillRaftServer()
	nodeTwo.KillRaftServer()
	nodeThree.KillRaftServer()
	nodeFour.KillRaftServer()
	nodeFive.KillRaftServer()

	require.Equal(t, nodeOne.IsRaftServerDead(), true)
	require.Equal(t, nodeTwo.IsRaftServerDead(), true)
	require.Equal(t, nodeThree.IsRaftServerDead(), true)
	require.Equal(t, nodeFour.IsRaftServerDead(), true)
	require.Equal(t, nodeFive.IsRaftServerDead(), true)

	// Run the cluster
	nodeOne.run()
}


// Notes
// Create Network Class that represents a bidirectional network
// Network should have a number of nodes in its constructor
// Each connection between two nodes is a queue
// Each node has a dictionary of connections to other nodes which is a dictionary of queues
// Send a message of type raft Message to a node

// Create a priority queue in the constructor of the network that represents the message queue
// Set delay in constructor
// Create a method that sends a message to a node that has a delay then adds the message to the queue with parameters of message and random delay
// Create a method to recieve messages from the queue
   // If delay is 0 then return the message from the queue

// Create a tick method that decrements the delay of each message in the queue and sleeps for 1 second