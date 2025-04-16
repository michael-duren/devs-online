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

func newLoginResult(conn *websocket.Conn, username *string, err error) LoginResult {
	return LoginResult{
		conn:     conn,
		username: username,
		err:      err,
	}
}

func connectToChat(creds *messages.Credentials, logger *log.Logger) tea.Cmd {
	logger.Infof("in connect to chat with :%v", creds)
	return func() tea.Msg {
		u := url.URL{Scheme: "ws", Host: creds.Address, Path: "/ws"}
		query := u.Query()
		query.Set("username", creds.Username)
		query.Set("secret", creds.Secret)
		u.RawQuery = query.Encode()

		logger.Infof("executing url: %v", u.String())

		c, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
		logger.Infof("websocket dialed: got conn: %v, res: %v, err: %v", c, res, err)
		if res != nil && res.StatusCode < 101 {
			return newLoginResult(
				nil,
				nil,
				fmt.Errorf("bad handshake recieved %d status code", res.StatusCode),
			)
		}
		if err != nil {
			log.Info("unable to connect with websocket. res: %v", res)
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

		return m, tea.Batch(m.Loading.Init(), connectToChat(creds, m.Logger))
	}

	return m, tea.Batch(cmds...)
}
