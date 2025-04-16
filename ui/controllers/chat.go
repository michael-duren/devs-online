package controllers

import (
	"encoding/json"
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/messages"
	"github.com/michael-duren/tui-chat/ui/models"
)

func handleWebSocketMessage(m *models.AppModel, websocketMsg messages.WebSocketMessage) (*models.AppModel, tea.Cmd) {
	var msg messages.Message
	conn := m.Chat.Conn
	if err := json.Unmarshal(websocketMsg.Data, &msg); err != nil {
		m.Logger.Errorf("unable to unmarshal websocket msg error: %v", err)
		return m, ListenForWebSocketMessages(conn)
	}

	if msg.Type != messages.InitMessageType {
		// we want to add to history for everything except an init msg
		m.Chat.Messages = append(m.Chat.Messages, msg)
	}

	switch msg.Type {
	// note: for chat msgs all we want to do is append so we don't have to do anything else here
	case messages.JoinMessageType:
		var joinMsg messages.JoinMessage
		if err := json.Unmarshal([]byte(msg.Content), &joinMsg); err != nil {
			log.Errorf("error decoding the join msg: %v", err)
		} else {
			m.Chat.Participants = append(m.Chat.Participants, messages.Participant{
				Username: joinMsg.Username,
				Online:   true,
			})
		}
	case messages.LeaveMessageType:
		var leaveMsg messages.LeaveMessage
		if err := json.Unmarshal([]byte(msg.Content), &leaveMsg); err != nil {
			log.Errorf("error decoding the leave msg:  %v", err)
		} else {
			for i, p := range m.Chat.Participants {
				if p.Username == leaveMsg.Username {
					m.Chat.Participants = slices.Delete(m.Chat.Participants, i, i+1)
					break
				}
			}
		}
	case messages.InitMessageType:
		var initMsg messages.InitMessage

		if err := json.Unmarshal([]byte(msg.Content), &initMsg); err != nil {
			log.Errorf("error decoding the init msg:  %v", err)
		} else {
			m.Chat.Messages = initMsg.ChatHistory
			m.Chat.Participants = initMsg.Participants
		}
	case messages.ShutdownMessageType:
		m.Logger.Infof("server is shutting down, shutting down app")
		_ = m.Chat.Conn.Close()
		return m, tea.Quit
	}

	return m, ListenForWebSocketMessages(conn)
}

func Chat(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.Chat.Input.Value()
			if val == "" {
				return m, nil
			}

			chatMsg := messages.NewChatMessage(val, m.Chat.Username)
			if err := m.Chat.Conn.WriteJSON(chatMsg); err != nil {
				m.Logger.Errorf("failed to send chat msg from client: %v", err)
			}

			m.Chat.Input.Reset()
			m.Chat.Messages = append(m.Chat.Messages, *chatMsg)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case messages.WebSocketMessage:
		return handleWebSocketMessage(m, msg)
	case messages.WebSocketError:
		m.Logger.Errorf("websocket err: %v", msg.Err)
		m.CurrentView = models.LoginPath
		m.Login.NetworkError = fmt.Sprintf("there was an issue recieving messages in the chat: %v", msg.Err)
		return m, nil
	}

	var cmd tea.Cmd
	m.Chat.Input, cmd = m.Chat.Input.Update(msg)
	return m, tea.Batch(cmd, ListenForWebSocketMessages(m.Chat.Conn))
}
