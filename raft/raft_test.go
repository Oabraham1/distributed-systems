package raft

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRaftServer() (*RaftServer){
	raft := NewRaftServer(5)
	return raft
}

func TestCreateRaftServer(t *testing.T) {
	raft := CreateRaftServer()
	require.Equal(t, len(raft.Nodes), int(5))
	require.Equal(t, raft.LeaderID, int64(-1))
}

func TestGetNode(t *testing.T) {
	raft := CreateRaftServer()
	node := raft.GetNode(0)
	require.Equal(t, node.Node.ID, int64(0))
}

func TestElection(t *testing.T) {
	raft := CreateRaftServer()
	raft.StartElection(raft.GetNode(0))
	require.Equal(t, raft.LeaderID, int64(0))
}