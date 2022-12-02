package network

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/oabraham1/distributed-systems/messages"
)

func CreateNetwork(t *testing.T) (*Network){
	network := NewNetwork()
	return network
}

func CreateNewMessage(id, from int64, to int64, delay int64, messageType messages.MessageType) *messages.Message {
	return &messages.Message{
		ID:     id,
		From:   from,
		To:     to,
		Type:   messageType,
		Delay: delay,
		Recieved: 0,
	}
}

func TestCreateNetWork(t *testing.T) {
	network := CreateNetwork(t)
	require.Equal(t, len(network.Nodes), int(0))
	require.Equal(t, network.IsNetworkDead, false)
}

func TestAddNode(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 5; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	require.Equal(t, len(network.Nodes), int(5))
}

func TestRemoveNode(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 5; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	for i := 0; i < 5; i++ {
		network.RemoveNode(int64(i))
	}
	require.Equal(t, len(network.Nodes), int(0))
}

func TestRemoveMessage(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 5; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	require.Equal(t, len(network.Nodes), int(5))

	for i := 0; i < 5; i++ {
		network.AddMessage(int64(i), int64(i+1), int64(i), 4, messages.STARTLEADERELECTION)
	}
	require.Equal(t, len(network.Messages), int(5))

	for i := 0; i < 5; i++ {
		message := network.RemoveMessage(CreateNewMessage(int64(i), int64(i+1), int64(i), 4, messages.STARTLEADERELECTION))
		require.NotNil(t, message)
	}
	require.Equal(t, len(network.Messages), int(0))
	message := network.RemoveMessage(CreateNewMessage(int64(100), int64(150), int64(130), 4, messages.STARTLEADERELECTION))
	require.Nil(t, message)
}

func TestSendAndRecieveMessage(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 5; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	require.Equal(t, len(network.Nodes), int(5))

	for i := 0; i < 5; i++ {
		network.AddMessage(rand.Int63n(11), int64(i), int64(i+1), 4, messages.STARTLEADERELECTION)
	}
	require.Equal(t, len(network.Messages), int(5))

	for i := 0; i < 4; i++ {
		message := network.Messages[i]
		network.SendMessage(message)
		require.Equal(t, message.Recieved, int64(0))
		require.Equal(t, len(network.GetNode(message.From).Sent), int(1))
		require.Nil(t, network.GetNode(100))

		network.ReceiveMessage(message)
		require.Equal(t, message.Recieved, int64(1))
		require.Equal(t, len(network.GetNode(message.To).Received), int(1))
	}
}

func TestStreamMessages(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 4; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	require.Equal(t, len(network.Nodes), int(4))

	network.AddMessage(rand.Int63n(11), int64(0), int64(1), 4, messages.STARTLEADERELECTION)
	network.AddMessage(rand.Int63n(11), int64(1), int64(0), 4, messages.STARTLEADERELECTION)
	require.Equal(t, len(network.Messages), int(2))

	network.StreamMessages()

	require.Equal(t, len(network.Messages), int(0))
	require.Equal(t, len(network.GetNode(int64(0)).Sent), int(1))
	require.Equal(t, len(network.GetNode(int64(1)).Sent), int(1))
	require.Equal(t, len(network.GetNode(int64(0)).Received), int(1))
	require.Equal(t, len(network.GetNode(int64(1)).Received), int(1))
	require.Equal(t, network.GetNode(int64(1)).Sent[0].Delay, int64(-1))
}

func TestTick(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 5; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	require.Equal(t, len(network.Nodes), int(5))

	for i := 0; i < 5; i++ {
		network.AddMessage(rand.Int63n(11), int64(i), int64(i+1), 4, messages.STARTLEADERELECTION)
	}
	require.Equal(t, len(network.Messages), int(5))

	messages := network.Messages
	for _, message := range network.Messages {
		delay := message.Delay
		network.Tick(message)
		require.Equal(t, message.Delay, delay-1)
	}

	for i, message := range network.Messages {
		require.Equal(t, message.Delay, (messages)[i].Delay)
	}
}

func TestKillNetwork(t *testing.T) {
	network := CreateNetwork(t)
	for i := 0; i < 5; i++ {
		network.AddNode(NewNode(int64(i)))
	}
	require.Equal(t, len(network.Nodes), int(5))
	network.KillNetwork()

	for i := 0; i < 5; i++ {
		require.Equal(t, network.GetNode(int64(i)).IsNodeDead, true)
	}
	require.Equal(t, network.IsNetworkDead, true)
}