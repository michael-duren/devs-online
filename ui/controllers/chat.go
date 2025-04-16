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

func Chat(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	m.Logger.Infof("in chat ctlr: participants %v", m.Chat.Participants)
	m.Logger.Infof("in chat ctlr: msgs %v", m.Chat.Messages)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.Chat.Input.Value()
			if val == "" {
				return m, nil
			}

			chatMsg := messages.NewChatMessage(val, m.Chat.Username)
			m.Logger.Infof("in chat ctlr writing: %v", chatMsg)
			if err := m.Chat.Conn.WriteJSON(chatMsg); err != nil {
				m.Logger.Errorf("failed to send chat msg from client: %v", err)
			}

			m.Chat.Input.Reset()
			m.Chat.Messages = append(m.Chat.Messages, *chatMsg)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case messages.WebSocketMessage:
		m.Logger.Infof("in chat ctlr recieving websocketmsg: %v", msg)
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

func handleWebSocketMessage(m *models.AppModel, websocketMsg messages.WebSocketMessage) (*models.AppModel, tea.Cmd) {
	var msg messages.Message
	m.Logger.Infof("in handleWebSocketMessage: %v", m)
	m.Logger.Infof("Raw WebSocket Message Data: %v", string(websocketMsg.Data))
	conn := m.Chat.Conn
	if err := json.Unmarshal(websocketMsg.Data, &msg); err != nil {
		m.Logger.Errorf("unable to unmarshal websocket msg error: %v", err)
		return m, ListenForWebSocketMessages(conn)
	}
	m.Logger.Infof("Decoded Message: Type=%v, Content=%v, Sender=%v",
		msg.Type, msg.Content, msg.Sender)

	switch msg.Type {
	case messages.ChatMessageType:
		var chatMsg messages.ChatMessage
		if err := json.Unmarshal([]byte(msg.Content), &chatMsg); err != nil {
			m.Logger.Errorf("error decoding chat msg: %v", err)
			m.Logger.Errorf("Problematic content: %v", msg.Content)
		} else {
			m.Logger.Infof("Decoded Chat Message: %+v", chatMsg)
			m.Chat.Messages = append(m.Chat.Messages, msg)
		}
	case messages.JoinMessageType:
		var joinMsg messages.JoinMessage
		if err := json.Unmarshal([]byte(msg.Content), &joinMsg); err != nil {
			log.Errorf("error decoding the join msg: %v", err)
		} else {
			m.Chat.Participants = append(m.Chat.Participants, messages.Participant{
				Username: joinMsg.Username,
				Online:   true,
			})
			m.Chat.Messages = append(m.Chat.Messages, msg)
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

			m.Chat.Messages = append(m.Chat.Messages, msg)
		}
	case messages.InitMessageType:
		var initMsg messages.InitMessage
		m.Logger.Infof("in init msg: %v", msg)

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
	default:
		m.Logger.Infof("unexpected msg: %v", msg)
	}

	return m, ListenForWebSocketMessages(conn)
}
