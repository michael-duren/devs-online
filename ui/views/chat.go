package views

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/messages"
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

	var msgs []string
	for _, msg := range chatModel.Messages {
		var s string
		switch msg.Type {
		case messages.ChatMessageType:
			var chatMsg messages.ChatMessage
			if err := json.Unmarshal([]byte(msg.Content), &chatMsg); err != nil {
				log.Errorf("unable to unmarshal json: %v", err)
			}
			s = fmt.Sprintf(
				"%s (%s): %s",
				chatMsg.Username,
				chatMsg.Date.Format(time.Kitchen),
				chatMsg.Message,
			)
		}
		// TODO: Add other cases

		msgs = append(msgs, s)
	}
	messageContent := strings.Join(msgs, "\n")
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
