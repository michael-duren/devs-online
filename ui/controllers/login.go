package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
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

func connectToChat(creds *models.Credentials) tea.Cmd {
	return func() tea.Msg {
		u := url.URL{Scheme: "ws", Host: creds.Address, Path: "/ws"}
		credStr, err := json.Marshal(*models.NewCredentialDto(creds.Username, creds.Secret))
		authHeader := http.Header{}
		authHeader.Set("Authorization", string(credStr))
		if err != nil {
			return newLoginResult(
				nil,
				nil,
				err,
			)
		}

		c, res, err := websocket.DefaultDialer.Dial(u.String(), authHeader)
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
		m.CurrentView = models.Loading
		creds := models.NewCredentials(
			m.Login.Address,
			m.Login.Username,
			m.Login.Secret,
		)
		m.Chat.Username = creds.Username

		return m, tea.Batch(m.Loading.Init(), connectToChat(creds))
	}

	return m, tea.Batch(cmds...)
}
