package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateNetwork(t *testing.T) (*Network){
	network := NewNetwork()
	return network
}

func TestCreateNetWork(t *testing.T) {
	network := CreateNetwork(t)
	require.Equal(t, len(*network.Nodes), int(0))
	require.Equal(t, network.IsNetworkDead, false)
}