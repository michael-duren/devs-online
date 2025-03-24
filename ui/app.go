package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/internal/logging"
	"github.com/michael-duren/tui-chat/ui/controllers"
	"github.com/michael-duren/tui-chat/ui/models"
	"github.com/michael-duren/tui-chat/ui/views"
)

type Model struct {
	*models.AppModel
}

func InitialModel() Model {
	logger := logging.NewLogger("client")
	appModel := models.NewAppModel(logger)
	return Model{
		AppModel: appModel,
	}
}

func (m Model) Init() tea.Cmd {
	return m.Login.Form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := controllers.Base(m.AppModel, msg)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	return views.Layout(m.AppModel)
}
