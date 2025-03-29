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

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		spinnerStyle.Render(loadingModel.Spinner.View()),
		textStyle.Render("Loading..."),
		textStyle.Render("Press 'q' to quit"),
	)

	return style.Render(content)
}
