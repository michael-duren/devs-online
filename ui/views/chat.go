package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Chat(m *models.AppModel) string {
	chatModel := m.Chat

	messageStyle := lipgloss.NewStyle().
		Foreground(Gray).
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Violet)

	inputStyle := lipgloss.NewStyle().
		Foreground(Cyan).
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Violet)

	messageHeight := m.BodyDimensions.Height - 8 // Leave space for input
	inputHeight := 3

	var messages []string
	messages = append(messages, chatModel.Messages...)
	messageContent := strings.Join(messages, "\n")
	if messageContent == "" {
		messageContent = "No messages yet..."
	}

	messageArea := messageStyle.
		Height(messageHeight).
		Width(m.BodyDimensions.Width - 2).
		Render(messageContent)

	inputArea := inputStyle.
		Height(inputHeight).
		Width(m.BodyDimensions.Width - 2).
		Render(chatModel.Input.View())

	return lipgloss.JoinVertical(
		lipgloss.Top,
		messageArea,
		inputArea,
	)
}
