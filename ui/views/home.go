package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

const welcomeTwo = `
   _____     ____    _        
  |  __ \   / __ \  | |       
  | |  | | | |  | | | |       
  | |  | | | |  | | | |       
  | |__| | | |__| | | |____   
  |_____/   \____/  |______|

  Devs     Online   😉
    `

const getStarted = "\nPress 's' or 'Enter' to get started"

func Home(m *models.AppModel) string {
	homeModel := m.Home
	if homeModel.Name != "" {
		return fmt.Sprintf("Welcome %s to TUI CHAT", homeModel.Name)
	}

	m.Logger.Info("height: ", m.BodyDimensions.Height)
	welcomeStyle := lipgloss.NewStyle().
		Foreground(Purple).
		Bold(true)

	getStartedStyle := lipgloss.NewStyle().
		Foreground(Violet).
		Bold(true)

	return lipgloss.Place(
		m.BodyDimensions.Width,
		m.BodyDimensions.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			welcomeStyle.Render(welcomeTwo),
			getStartedStyle.Render(getStarted),
		),
	)
}
