package views

import (
	"fmt"

	"github.com/michael-duren/tui-chat/ui/models"
)

func Chat(m *models.AppModel) string {
	chatModel := m.Chat

	return fmt.Sprint("chat page", chatModel, "\n\n", *m.Chat.Response.Message)
}
