package models

import "github.com/charmbracelet/log"

type WindowDemnsions struct {
	Width  int
	Height int
}

type CurrentView string

const (
	Home  CurrentView = "/home"
	Login CurrentView = "/login"
)

type AppModel struct {
	Logger *log.Logger
	*WindowDemnsions
	BodyDimensions *WindowDemnsions
	CurrentView    CurrentView
	// Page Models
	Home *HomeModel
}

func NewAppModel(logger *log.Logger) *AppModel {
	homeModel := NewHomeModel()
	return &AppModel{
		Logger: logger,
		WindowDemnsions: &WindowDemnsions{
			Width:  0,
			Height: 0,
		},
		BodyDimensions: &WindowDemnsions{
			Width:  0,
			Height: 0,
		},
		CurrentView: Home,
		Home:        homeModel,
	}
}
