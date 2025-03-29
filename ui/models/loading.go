package models

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type LoadingModel struct {
	Spinner spinner.Model
}

func NewLoadingModel() *LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &LoadingModel{
		Spinner: s,
	}
}
