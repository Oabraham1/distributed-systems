package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateNetwork(t *testing.T) (*Network){
	network := NewNetwork(5)
	require.Equal(t, network.nodes, int64(5))
	return network
}

func TestCreateNetWork(t *testing.T) {
	network := CreateNetwork(t)
	network.KillNetwork()
	require.Equal(t, network.IsNetworkDead, true)
}