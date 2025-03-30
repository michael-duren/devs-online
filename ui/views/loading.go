package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Loading(m *models.AppModel) string {
	loadingModel := m.Loading

	style := lipgloss.NewStyle().
		Foreground(Cyan).
		Padding(1).
		Align(lipgloss.Center)

	spinnerStyle := lipgloss.NewStyle().
		Foreground(Violet).
		Bold(true)

	textStyle := lipgloss.NewStyle().
		Foreground(Gray)

	content := lipgloss.Place(
		m.BodyDimensions.Width,
		m.BodyDimensions.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			spinnerStyle.Render(loadingModel.Spinner.View()),
			textStyle.Render("\n"),
			textStyle.Render("Loading..."),
		),
	)

	return style.Render(content)
}
