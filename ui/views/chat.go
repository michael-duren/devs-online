package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Chat(m *models.AppModel) string {
	chatModel := m.Chat

	sidebar := Sidebar(m)
	sidebarWidth := 38

	chatWidth := m.BodyDimensions.Width - sidebarWidth - 2
	messageStyle := lipgloss.NewStyle().
		Foreground(Gray).
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Violet)

	inputStyle := lipgloss.NewStyle().
		Foreground(Cyan).
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Purple)

	messageHeight := m.BodyDimensions.Height - 7
	inputHeight := 3

	var messages []string
	for _, msg := range chatModel.Messages {
		s := fmt.Sprintf(
			"%s (%s): %s",
			msg.Username,
			msg.Date.Format(time.Kitchen),
			msg.Message,
		)
		messages = append(messages, s)
	}
	messageContent := strings.Join(messages, "\n")
	if messageContent == "" {
		messageContent = "No messages yet..."
	}

	messageArea := messageStyle.
		Height(messageHeight).
		Width(chatWidth).
		Render(messageContent)

	inputArea := inputStyle.
		Height(inputHeight).
		Width(chatWidth).
		Render(chatModel.Input.View())

	chatArea := lipgloss.JoinVertical(
		lipgloss.Top,
		messageArea,
		inputArea,
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		sidebar,
		chatArea,
	)
}
