package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/messages"
)

type ChatModel struct {
	Response *messages.DummyResponse
	// TODO: Create msg struct that contains username, other details
	Messages    []string
	Input       textinput.Model
	Credentials *Credentials
}

func NewChatModel() *ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()

	return &ChatModel{
		Response: nil,
		Messages: make([]string, 0),
		Input:    ti,
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return textinput.Blink
}
