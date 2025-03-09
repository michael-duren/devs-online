package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/internal/logging"
	"github.com/michael-duren/tui-chat/ui/controllers"
	"github.com/michael-duren/tui-chat/ui/models"
	"github.com/michael-duren/tui-chat/ui/views"
)

type CurrentView string

const (
	Home CurrentView = "/home"
)

type Model struct {
	CurrentView CurrentView
	*models.AppModel

	// Page Models
	Home *models.HomeModel
}

func InitialModel() Model {
	logger := logging.NewLogger("client")
	appModel := models.NewAppModel(logger)
	homeModel := models.NewHomeModel()
	return Model{
		AppModel:    appModel,
		CurrentView: Home,
		Home:        homeModel,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := controllers.Base(*m.AppModel, msg)
	if cmd != nil {
		return m, cmd
	}

	switch m.CurrentView {
	case Home:
		return controllers.Home(m, msg)
	}

	return m, nil
}

func (m Model) View() string {
	// header
	header := "TUI CHAT\n"

	var body string
	switch m.CurrentView {
	case Home:
		body = views.Home(m.Home, m.Logger)
	}

	footer := "\nTUI CHAT BY MICHAEL DUREN"
	return views.Layout(header, body, footer)
}
