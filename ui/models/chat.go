package models

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/messages"
)

type ChatMessage struct {
	Date     time.Time
	Message  string
	Username string
}

func NewChatMessage(date time.Time, msg, username string) *ChatMessage {
	return &ChatMessage{
		Date:     date,
		Message:  msg,
		Username: username,
	}
}

type ChatModel struct {
	Response    *messages.DummyResponse
	Messages    []*ChatMessage
	Input       textinput.Model
	Credentials *Credentials
}

func NewChatModel() *ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()

	return &ChatModel{
		Response: nil,
		Messages: make([]*ChatMessage, 0),
		Input:    ti,
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return textinput.Blink
}
