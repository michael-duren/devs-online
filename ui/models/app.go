package models

import "github.com/charmbracelet/log"

type WindowDemnsions struct {
	Width  int
	Height int
}

type CurrentView string

const (
	HomePath    CurrentView = "/home"
	LoginPath   CurrentView = "/login"
	LoadingPath CurrentView = "/loading"
	ChatPath    CurrentView = "/chat"
)

type AppModel struct {
	Logger *log.Logger
	*WindowDemnsions
	BodyDimensions *WindowDemnsions
	CurrentView    CurrentView
	// Page Models
	Home    *HomeModel
	Login   *LoginModel
	Loading *LoadingModel
	Chat    *ChatModel
}

func NewAppModel(logger *log.Logger) *AppModel {
	homeModel := NewHomeModel()
	loginModel := NewLoginModel()
	loadingModel := NewLoadingModel()
	chatModel := NewChatModel()
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
		CurrentView: HomePath,
		Home:        homeModel,
		Login:       loginModel,
		Loading:     loadingModel,
		Chat:        chatModel,
	}
}
