package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/michael-duren/tui-chat/messages"
)

type ChatModel struct {
	Messages     []messages.Message
	Input        textinput.Model
	Username     string
	Participants []messages.Participant
	Conn         *websocket.Conn
}

func NewChatModel() *ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()

	return &ChatModel{
		Messages:     make([]messages.Message, 0),
		Input:        ti,
		Participants: []messages.Participant{},
	}
}

func (m *ChatModel) AddParticipant(username string, online bool) {
	m.Participants = append(m.Participants, messages.Participant{
		Username: username,
		Online:   online,
	})
}

func (m *ChatModel) Init() tea.Cmd {
	return textinput.Blink
}
