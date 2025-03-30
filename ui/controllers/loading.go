package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/messages"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Loading(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.DummyResponse:
		m.CurrentView = models.Chat
		m.Chat.Response = &msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			m.Logger.Infof("In default : %v\n", msg)
		}
	}
	var cmd tea.Cmd
	m.Loading.Spinner, cmd = m.Loading.Spinner.Update(msg)
	return m, cmd
}
