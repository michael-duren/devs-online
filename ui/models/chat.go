package models

import "github.com/michael-duren/tui-chat/ui/messages"

type ChatModel struct {
	Response *messages.DummyResponse
}

func NewChatModel() *ChatModel {
	return &ChatModel{
		Response: nil,
	}
}
