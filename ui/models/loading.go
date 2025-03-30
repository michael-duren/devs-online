package models

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoadingModel struct {
	Spinner spinner.Model
}

func NewLoadingModel() *LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &LoadingModel{
		Spinner: s,
	}
}

func (m *LoadingModel) Init() tea.Cmd {
	return m.Spinner.Tick
}
