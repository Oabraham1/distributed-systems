package messages

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateNewMessage(from int64, to int64, id int64, messageType MessageType, delay int64) *Message {
	return &Message{
		ID:     id,
		From:   from,
		To:     to,
		Type:   messageType,
		Delay: delay,
		Recieved: 0,
	}
}

func TestMessageToString(t *testing.T) {
	message := CreateNewMessage(1, 2, 3, STARTLEADERELECTION, rand.Int63n(10))
	require.Equal(t, message.String(), "Message: ID: 3, From: 1, To: 2, Type: 0")
}

func TestMessageToBytes(t *testing.T) {
	message := CreateNewMessage(1, 2, 3, STARTLEADERELECTION, rand.Int63n(10))
	bytes := MessageToBytes(*message)
	require.Equal(t, bytes, []byte{3, 1, 2, 0, byte(message.Delay), byte(message.Recieved)})
}

func TestBytesToMessage(t *testing.T) {
	message := CreateNewMessage(1, 2, 3, STARTLEADERELECTION, rand.Int63n(10))
	bytes := MessageToBytes(*message)
	message2 := BytesToMessage(bytes)
	require.Equal(t, message, message2)
}