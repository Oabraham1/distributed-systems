package messages

import "fmt"

// Message type contains message ID
// Send message to start leader election
// Convert message type into byte array and convert byte array into message type when recieved

type Message struct {
	ID     int64
	From   int64
	To     int64
	Type   MessageType
	Delay 		int64
	Recieved 	int64
}

type MessageType int

const (
	STARTLEADERELECTION MessageType = iota
	APPROVELEADERELCTION
	STOPLEADERELECTION
	REQUESTVOTE
	LEADERCONFIRMATIONREQUEST
	CONFIRMLEADER
	REJECTLEADER
	VOTE
	LEADERELECTION
)

func NewMessage(id int64, from int64, to int64, delay int64, messageType MessageType) *Message {
	return &Message{
		ID:     id,
		From:   from,
		To:     to,
		Type:   messageType,
		Delay: delay,
		Recieved: 0,
	}
}

func (message *Message) String() string {
	return fmt.Sprintf("Message: ID: %d, From: %d, To: %d, Type: %d", message.ID, message.From, message.To, message.Type)
}

func MessageToBytes(message Message) []byte {
	return []byte{byte(message.ID), byte(message.From), byte(message.To), byte(message.Type), byte(message.Delay), byte(message.Recieved)}
}

func BytesToMessage(bytes []byte) *Message {
	return &Message{
		ID:     int64(bytes[0]),
		From:   int64(bytes[1]),
		To:     int64(bytes[2]),
		Type:   MessageType(bytes[3]),
		Delay: int64(bytes[4]),
		Recieved: int64(bytes[5]),
	}
}