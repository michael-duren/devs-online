package views

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Login(m *models.AppModel) string {
	loginModel := m.Login

	switch loginModel.Form.State {
	case huh.StateCompleted:
		m.Logger.Info("You logged in fwiend")
		return "Logged in"
	default:
		v := strings.TrimSuffix(m.Login.Form.View(), "\n\n")
		form := lipgloss.NewStyle().Margin(1, 0).Render(v)

		return lipgloss.Place(
			m.BodyDimensions.Width,
			m.BodyDimensions.Height,
			lipgloss.Center,
			lipgloss.Center,
			form,
		)

	}
}
