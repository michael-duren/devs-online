package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Sidebar(m *models.AppModel) string {
	sidebarStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Dark_cyan).
		Padding(1).
		Width(36).
		Height(m.BodyDimensions.Height - 2)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(Cyan).
		MarginBottom(1)

	participantStyle := lipgloss.NewStyle().
		Foreground(Gray)

	var participants []string
	participants = append(participants, headerStyle.Render("Participants"))
	participants = append(participants, strings.Repeat("─", 18))
	for _, participant := range m.Chat.Participants {
		// Add an online status indicator
		status := "○"
		if participant.Online {
			status = "●"
		}

		participantLine := fmt.Sprintf("%s %s", status, participant.Username)
		if participant.Username == m.Chat.Credentials.Username {
			participantLine += " (you)"
		}

		participants = append(participants, participantStyle.Render(participantLine))
	}

	content := strings.Join(participants, "\n")
	return sidebarStyle.Render(content)
}
