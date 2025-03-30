package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Chat(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Chat.Input.Value() != "" {
				m.Chat.Messages = append(m.Chat.Messages, m.Chat.Input.Value())
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
