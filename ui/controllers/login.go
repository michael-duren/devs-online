package controllers

import (
	"fmt"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/michael-duren/tui-chat/messages"
	"github.com/michael-duren/tui-chat/ui/models"
)

type LoginResult struct {
	conn     *websocket.Conn
	username *string
	err      error
}

func newLoginResult(conn *websocket.Conn, username *string, err error) *LoginResult {
	return &LoginResult{
		conn:     conn,
		username: username,
		err:      err,
	}
}

func connectToChat(creds *messages.Credentials) tea.Cmd {
	return func() tea.Msg {
		path := fmt.Sprintf("/ws?username=%s&secret=%s", creds.Username, creds.Secret)
		u := url.URL{Scheme: "ws", Host: creds.Address, Path: path}

		c, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Info("unable to connect with websocket. res: %v", res)
			// TODO: check response for reason of not being successful
			return newLoginResult(
				nil,
				nil,
				err,
			)
		}

		return newLoginResult(c, &creds.Username, nil)
	}
}

func Login(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := m.Login.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Login.Form = f
		cmds = append(cmds, cmd)
	}

	if m.Login.Form.State == huh.StateCompleted {
		m.Logger.Infof("Form values - Address: %s, Username: %s, Secret: %s",
			m.Login.Address,
			m.Login.Username,
			m.Login.Secret)
		m.CurrentView = models.LoadingPath
		creds := messages.NewCredentials(
			m.Login.Address,
			m.Login.Username,
			m.Login.Secret,
		)
		m.Chat.Username = creds.Username

		return m, tea.Batch(m.Loading.Init(), connectToChat(creds))
	}

	return m, tea.Batch(cmds...)
}
