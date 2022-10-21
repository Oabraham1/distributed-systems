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
	nodeOne.peers = []*Raft{nodeTwo, nodeThree, nodeFour, nodeFive}
	nodeTwo.peers = []*Raft{nodeOne, nodeThree, nodeFour, nodeFive}
	nodeThree.peers = []*Raft{nodeOne, nodeTwo, nodeFour, nodeFive}
	nodeFour.peers = []*Raft{nodeOne, nodeTwo, nodeThree, nodeFive}
	nodeFive.peers = []*Raft{nodeOne, nodeTwo, nodeThree, nodeFour}

	require.Equal(t, len(nodeOne.peers), 4)
	require.Equal(t, len(nodeTwo.peers), 4)
	require.Equal(t, len(nodeThree.peers),4)
	require.Equal(t, len(nodeFour.peers), 4)
	require.Equal(t, len(nodeFive.peers), 4)

	// No leader should be elected
	require.Nil(t, nodeOne.state.leader)
	require.Nil(t, nodeTwo.state.leader)
	require.Nil(t, nodeThree.state.leader)
	require.Nil(t, nodeFour.state.leader)
	require.Nil(t, nodeFive.state.leader)

	nodeOne.sendRequestVoteToPeer(2)
	nodeOne.sendRequestVoteToPeer(3)
	nodeOne.sendRequestVoteToPeer(4)
	nodeOne.sendRequestVoteToPeer(5)

	// No leader should be elected
	require.Nil(t, nodeOne.state.leader)
}
