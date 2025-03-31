package models

import (
	"errors"
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
				Value(&m.Address).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("please enter the ip address you wish to connect to")
					}
					return nil
				},
				),

			huh.NewInput().
				Key("username").
				Title("Username").
				Value(&m.Username).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("please enter a username for the chat")
					}
					return nil
				},
				),

			huh.NewInput().
				Key("secret").
				Title("Room Secret").
				Value(&m.Secret).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("please enter the chat secret")
					}
					return nil
				},
				),

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
