package models

import "github.com/charmbracelet/log"

type WindowDemnsions struct {
	Width  int
	Height int
}

type AppModel struct {
	Logger *log.Logger
	*WindowDemnsions
}

func NewAppModel(logger *log.Logger) *AppModel {
	return &AppModel{
		Logger: logger,
		WindowDemnsions: &WindowDemnsions{
			Width:  0,
			Height: 0,
		},
	}
}
