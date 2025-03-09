package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

const (
	Gray        = "#7E8294"
	Gray_alt    = "#4C566A"
	Background  = "#000"
	Red         = "#EC7279"
	Yellow      = "#ECBE7B"
	Orange      = "#DA8548"
	Green       = "#A0C980"
	Cyan        = "#4DB5BD"
	Dark_cyan   = "#5699AF"
	Blue        = "#6CB6EB"
	Violet      = "#A9A1E1"
	Purple      = "#D38AEA"
	Light_blue  = "#ADD8E6"
	Light_pink  = "#D8BFD8"
	Disabled    = "#676E95"
	Diff_red    = "#FB4934"
	Diff_green  = "#8EC07C"
	Diff_blue   = "#458588"
	Diff_yellow = "#FABD2F"
	White       = "#FFF"
)

func Layout(appModel *models.AppModel, body string) string {
	background := lipgloss.Color(Background)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Background(background).
		Foreground(lipgloss.Color(Violet)).
		Padding(1).
		Width(appModel.Width)
	header := headerStyle.Render("TUI CHAT")

	bodyStyle := lipgloss.NewStyle().
		Background(background).
		Foreground(lipgloss.Color(Gray)).
		Padding(1).
		Width(appModel.Width).
		Height(appModel.Height - 6)

	footerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(Cyan)).
		Background(background).
		Padding(1).
		Width(appModel.Width)
	footer := footerStyle.Render("TUI CHAT BY MICHAEL DUREN")

	appModel.BodyDimensions.Height = appModel.Height - 6
	appModel.BodyDimensions.Width = appModel.Width - 2

	return lipgloss.Place(
		appModel.Width,
		appModel.Height,
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
