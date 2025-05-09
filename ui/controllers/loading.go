package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/michael-duren/tui-chat/messages"
	"github.com/michael-duren/tui-chat/ui/models"
)

func ListenForWebSocketMessages(conn *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return messages.WebSocketError{
				Err:     err,
				Address: conn.LocalAddr().String(),
			}
		}
		return messages.WebSocketMessage{
			Data: data,
		}
	}
}

func Loading(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	m.Logger.Info("in loading ctlr")
	switch msg := msg.(type) {
	case LoginResult:
		m.Logger.Info("in login resutl")
		if msg.conn == nil || msg.err != nil || msg.username == nil {
			log.Warn("error logging into chat: ", "error: ", msg.err)
			m.CurrentView = models.LoginPath

			return m, nil
		}
		m.Chat.Conn = msg.conn
		m.Chat.Username = *msg.username
		m.CurrentView = models.ChatPath
		return m, ListenForWebSocketMessages(msg.conn)
	case tea.KeyMsg:
		switch msg.String() {
		default:
			m.Logger.Infof("In default : %v\n", msg)
		}
	}
	// m.Logger.Infof("msg: %v", msg)
	var cmd tea.Cmd
	m.Loading.Spinner, cmd = m.Loading.Spinner.Update(msg)
	return m, cmd
}
