package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Layout(m *models.AppModel) string {
	headerStyle := lipgloss.NewStyle().
		// Background(Background).
		Foreground(Violet).
		Padding(1).
		Width(m.Width)
	header := headerStyle.Render("TUI CHAT")

	bodyStyle := lipgloss.NewStyle().
		Bold(true).
		// Background(Background).
		Foreground(Gray).
		Padding(1).
		Width(m.Width).
		Height(m.Height - 6)

	footerStyle := lipgloss.NewStyle().
		Foreground(Cyan).
		// Background(Background).
		Padding(1).
		Width(m.Width)
	footer := footerStyle.Render("TUI CHAT BY MICHAEL DUREN")

	m.BodyDimensions.Height = m.Height - 9
	m.BodyDimensions.Width = m.Width - 2

	var body string
	switch m.CurrentView {
	case models.Home:
		body = Home(m)
	case models.Login:
		body = Login(m)
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			bodyStyle.Render(body),
			footer,
		),
	)
}
