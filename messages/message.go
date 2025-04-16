package messages

import "time"

const (
	ChatMessageType  MessageType = "chat"
	InitMessageType  MessageType = "init"
	LeaveMessageType MessageType = "leave"
	JoinMessageType  MessageType = "join"
)

type MessageType string

// Message
// A wrapper for sending different types of
// messages between peers
type Message struct {
	Type    MessageType `json:"type"`
	Content any         `json:"content"`
	Sender  string      `json:"sender,omitempty"`
	Time    time.Time   `json:"time"`
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
	Date     time.Time
	Message  string
	Username string
}

func NewChatMessage(date time.Time, msg, username string) *Message {
	return &Message{
		Type: ChatMessageType,
		Content: ChatMessage{
			Date:     date,
			Message:  msg,
			Username: username,
		},
		Sender: username,
		Time:   time.Now(),
	}
}

// InitMessage - The message the server
// should send to the client when entering the chat
type InitMessage struct {
	ChatHistory  []ChatMessage
	Participants []Participant
}

func NewInitMessage(ch []ChatMessage, p []Participant) *Message {
	return &Message{
		Type: InitMessageType,
		Content: &InitMessage{
			ChatHistory:  ch,
			Participants: p,
		},
		Time: time.Now(),
	}
}

type LeaveMessage struct {
	Username string `json:"username"`
}

func NewLeaveMessage(username string) *Message {
	return &Message{
		Type: LeaveMessageType,
		Content: &LeaveMessage{
			Username: username,
		},
		Time: time.Now(),
	}
}

type JoinMessage struct {
	Username string `json:"username"`
}

func NewJoinMessage(username string) *Message {
	return &Message{
		Type: LeaveMessageType,
		Content: &JoinMessage{
			Username: username,
		},
		Time: time.Now(),
	}
}
