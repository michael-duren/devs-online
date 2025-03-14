package models

import "github.com/charmbracelet/huh"

type LoginModel struct {
	Address  string
	Username string
	Secret   string
	Form     *huh.Form
}

func NewLoginModel() *LoginModel {
	m := &LoginModel{
		Address:  "",
		Username: "",
		Secret:   "",
	}
	m.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("IP Address").
				Value(&m.Address),

			huh.NewInput().
				Title("Username").
				Value(&m.Username),

			huh.NewInput().
				Title("Room Secret").
				Value(&m.Secret),
		),
	)
	return m
}
