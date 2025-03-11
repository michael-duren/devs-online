package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Home(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s", "enter":
			m.Logger.Debug("In case")
			m.CurrentView = models.Login
			return m, nil
		default:
			m.Logger.Infof("In default : %v\n", msg)
		}
	}

	return m, nil
}
