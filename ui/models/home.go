package models

type HomeModel struct {
	Name string
}

func NewHomeModel() *HomeModel {
	return &HomeModel{
		Name: "",
	}
}
