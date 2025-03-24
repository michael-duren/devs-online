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
		// TODO: Render form
		// var addr string
		// if m.Login.Form.GetString("address") != "" {
		// 	addr = "Address: " + m.Login.Form.GetString("address")
		v := strings.TrimSuffix(m.Login.Form.View(), "\n\n")
		form := lipgloss.NewStyle().Margin(1, 0).Render(v)
		m.Logger.Info("value of form", v)
		// errors := m.Login.Form.Errors()
		//
		// lipgloss.NewStyle().Padding(1, 4, 0, 1).Render(form)

		return lipgloss.Place(
			m.BodyDimensions.Width,
			m.BodyDimensions.Height,
			lipgloss.Center,
			lipgloss.Center,
			form,
		)

	}
}
