package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Loading(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case LoginResult:
		if msg.conn == nil || msg.err != nil || msg.username == nil {
			log.Warn("error logging into chat: ", "error: ", msg.err)
			m.CurrentView = models.LoginPath
			return m, nil
		}

		m.Chat.Conn = msg.conn
		m.Chat.Username = *msg.username
		m.CurrentView = models.ChatPath
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		default:
			m.Logger.Infof("In default : %v\n", msg)
		}
	}
	var cmd tea.Cmd
	m.Loading.Spinner, cmd = m.Loading.Spinner.Update(msg)
	return m, cmd
}
