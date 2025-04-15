package messages

import "time"

const (
	ChatMessageType MessageType = "chat"
	InitType        MessageType = "init"
	UserAlertType   MessageType = "user_alert"
)

type MessageType string

type Message interface {
	Type() MessageType
}

// Participant
type Participant struct {
	Username string
	Online   bool
}

// ChatMessage
// The chat message sent by the server
// or client
type ChatMessage struct {
	Date        time.Time
	Message     string
	Username    string
	messageType MessageType
}

func (c *ChatMessage) Type() MessageType {
	return c.messageType
}

func NewChatMessage(date time.Time, msg, username string) *ChatMessage {
	return &ChatMessage{
		Date:        date,
		Message:     msg,
		Username:    username,
		messageType: ChatMessageType,
	}
}

// InitMessage - The message the server
// should send to the client when entering the chat
type InitMessage struct {
	ChatHistory  []ChatMessage
	Participants []Participant
}

// UserAction - Used if a user leaves the chat
type UserAction struct {
	Participant Participant
}

//
// type Message[T any] struct {
// 	MessageType MessageType
// 	Payload     T
// }
