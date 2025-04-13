package models

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/messages"
)

type Participant struct {
	Username string
	Online   bool
}

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
	Response     *messages.DummyResponse
	Messages     []*ChatMessage
	Input        textinput.Model
	Credentials  *Credentials
	Participants []Participant
}

func NewChatModel() *ChatModel {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()

	return &ChatModel{
		Response:     nil,
		Messages:     make([]*ChatMessage, 0),
		Input:        ti,
		Participants: []Participant{},
	}
}

func (m *ChatModel) AddParticipant(username string, online bool) {
	m.Participants = append(m.Participants, Participant{
		Username: username,
		Online:   online,
	})
}

func (m *ChatModel) Init() tea.Cmd {
	return textinput.Blink
}
