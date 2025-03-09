package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui/models"
)

func Base(m models.AppModel, msg tea.Msg) (models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}
