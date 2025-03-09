package views

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Home(m *models.HomeModel, logger *log.Logger) string {
	if m.Name != "" {
		return fmt.Sprintf("Welcome %s to TUI CHAT", m.Name)
	}

	return "Welcome to TUI CHAT"
}
