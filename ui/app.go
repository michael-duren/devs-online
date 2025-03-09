package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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
	Logger      *log.Logger
	CurrentView CurrentView

	// Page Models
	Home *models.HomeModel
}

func InitialModel() Model {
	logger := logging.NewLogger("client")
	homeModel := models.NewHomeModel()
	return Model{
		// TODO: Add page models here
		Logger:      logger,
		CurrentView: Home,
		Home:        homeModel,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := controllers.Base(m, msg)
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
	s := "TUI CHAT\n"

	switch m.CurrentView {
	case Home:
		s += views.Home(m.Home, m.Logger)
	}

	s += "\nTUI CHAT BY MICHAEL DUREN"
	return s
}
