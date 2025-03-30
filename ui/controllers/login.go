package controllers

import (
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/ui/messages"
	"github.com/michael-duren/tui-chat/ui/models"
)

func makeDummyRequest(url string, logger *log.Logger) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 5 * time.Second}
		res, err := c.Get(url)
		if err != nil {
			logger.Errorf("error response: %v", err)
			return messages.DummyResponse{Err: err}
		}

		defer func() { _ = res.Body.Close() }()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("error reading response body: ", err)
			return messages.DummyResponse{Err: err}
		}

		payload := string(bodyBytes)
		return messages.DummyResponse{
			StatusCode: res.StatusCode,
			Message:    &payload,
			Err:        nil,
		}
	}
}

func Login(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s", "enter":
			m.Logger.Debug("In case")
			return m, nil
		default:
			m.Logger.Infof("In default : %v\n", msg)
		}
	}

	var cmds []tea.Cmd

	form, cmd := m.Login.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Login.Form = f
		cmds = append(cmds, cmd)
	}

	if m.Login.Form.State == huh.StateCompleted {
		// TODO: Update with actual logic
		m.Logger.Info("in state completed")
		url := "https://swapi.dev/api/people/11"
		m.Logger.Infof("Form values - Address: %s, Username: %s, Secret: %s",
			m.Login.Address,
			m.Login.Username,
			m.Login.Secret)
		m.CurrentView = models.Loading
		return m, tea.Batch(m.Loading.Init(), makeDummyRequest(url, m.Logger))
	}

	return m, tea.Batch(cmds...)
}
