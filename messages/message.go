package messages

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

const (
	ChatMessageType  MessageType = "chat"
	InitMessageType  MessageType = "init"
	LeaveMessageType MessageType = "leave"
	JoinMessageType  MessageType = "join"

	// TODO: Implement
	ShutdownMessageType MessageType = "shutdown"
)

type WebSocketMessage struct {
	Data []byte
}

type WebSocketError struct {
	Address string
	Err     error
}

func (e *WebSocketError) Error() string {
	return fmt.Sprintf("error connecting to %s: error: %v", e.Address, e.Err)
}

type MessageType string

// Message
// A wrapper for sending different types of
// messages between peers
type Message struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	Sender  string      `json:"sender,omitempty"`
	Time    time.Time   `json:"time"`
}

type Participant struct {
	Username string
	Online   bool
}

// ChatMessage
// The chat message sent by the server
// or client
type ChatMessage struct {
	Message  string
	Username string
}

func NewChatMessage(msg, username string) *Message {
	chatMsg := ChatMessage{
		Message:  msg,
		Username: username,
	}
	encoded, err := json.Marshal(chatMsg)
	if err != nil {
		log.Errorf("unable to encode chat msg: %v", err)
		return nil
	}

	return &Message{
		Type:    ChatMessageType,
		Content: string(encoded),
		Sender:  username,
		Time:    time.Now(),
	}
}

// InitMessage - The message the server
// should send to the client when entering the chat
type InitMessage struct {
	ChatHistory  []Message
	Participants []Participant
}

func NewInitMessage(ch []Message, p []Participant) *Message {
	initMsg := &InitMessage{
		ChatHistory:  ch,
		Participants: p,
	}
	encoded, err := json.Marshal(initMsg)
	if err != nil {
		log.Errorf("unable to encode chat msg: %v", err)
		return nil
	}

	return &Message{
		Type:    InitMessageType,
		Content: string(encoded),
		Time:    time.Now(),
	}
}

type LeaveMessage struct {
	Username string `json:"username"`
}

func NewLeaveMessage(username string) *Message {
	leaveMsg := LeaveMessage{
		Username: username,
	}

	encoded, err := json.Marshal(leaveMsg)
	if err != nil {
		log.Errorf("unable to encode leave msg: %v", err)
		return nil
	}

	return &Message{
		Type:    LeaveMessageType,
		Content: string(encoded),
		Time:    time.Now(),
	}
}

type JoinMessage struct {
	Username string `json:"username"`
}

func NewJoinMessage(username string) *Message {
	joinMsg := &JoinMessage{
		Username: username,
	}

	encoded, err := json.Marshal(joinMsg)
	if err != nil {
		log.Errorf("unable to encode join msg: %v", err)
		return nil
	}

	return &Message{
		Type:    LeaveMessageType,
		Content: string(encoded),
		Time:    time.Now(),
	}
}
