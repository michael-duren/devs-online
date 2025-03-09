package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/tui-chat/ui"
)

func main() {
	app := ui.InitialModel()
	app.Logger.Info("Starting Client")
	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		app.Logger.Errorf("oh no looks like we're ngmi. famous last words: %v", err)
		os.Exit(1)
	}
}
