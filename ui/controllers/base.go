package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Base(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}
