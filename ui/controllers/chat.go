package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Chat(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	return m, nil
}
