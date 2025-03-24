package models

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

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
				Key("address").
				Title("IP Address").
				Value(&m.Address),

			huh.NewInput().
				Key("username").
				Title("Username").
				Value(&m.Username),

			huh.NewInput().
				Key("secret").
				Title("Room Secret").
				Value(&m.Secret),
			huh.NewConfirm().
				Key("login").
				Title("Join Chat").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).WithWidth(45).
		WithShowHelp(true).
		WithShowErrors(true)

	return m
}
