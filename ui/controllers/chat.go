package controllers

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/messages"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Chat(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Chat.Input.Value() != "" {
				m.Chat.Messages = append(
					m.Chat.Messages,
					messages.NewChatMessage(
						time.Now(),
						m.Chat.Input.Value(),
						m.Chat.Username,
					),
				)
				m.Chat.Input.Reset()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Chat.Input, cmd = m.Chat.Input.Update(msg)
	return m, cmd
}
