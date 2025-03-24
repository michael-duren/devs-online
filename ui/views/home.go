package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/tui-chat/ui/models"
)

const welcomeMessage = `
/^^        /^^           /^^                                        
/^^        /^^           /^^                                        
/^^   /^   /^^   /^^     /^^   /^^^   /^^    /^^^ /^^ /^^    /^^    
/^^  /^^   /^^ /^   /^^  /^^ /^^    /^^  /^^  /^^  /^  /^^ /^   /^^ 
/^^ /^ /^^ /^^/^^^^^ /^^ /^^/^^    /^^    /^^ /^^  /^  /^^/^^^^^ /^^
/^ /^    /^^^^/^         /^^ /^^    /^^  /^^  /^^  /^  /^^/^        
/^^        /^^  /^^^^   /^^^   /^^^   /^^    /^^^  /^  /^^  /^^^^   
                                                                    
/^^^ /^^^^^^               /^^^^^           /^^^^        /^^        
     /^^                   /^^   /^^      /^^    /^^     /^^        
     /^^       /^^         /^^    /^^   /^^        /^^   /^^        
     /^^     /^^  /^^      /^^    /^^   /^^        /^^   /^^        
     /^^    /^^    /^^     /^^    /^^   /^^        /^^   /^^        
     /^^     /^^  /^^      /^^   /^^      /^^     /^^    /^^        
     /^^       /^^         /^^^^^    /^^    /^^^^     /^^/^^^^^^^^  
    `

const getStarted = "Press 's' or 'Enter' to get started"

func Home(m *models.AppModel) string {
	homeModel := m.Home
	if homeModel.Name != "" {
		return fmt.Sprintf("Welcome %s to TUI CHAT", homeModel.Name)
	}

	m.Logger.Info("height: ", m.BodyDimensions.Height)
	welcomeStyle := lipgloss.NewStyle().
		// Background(Background).
		Foreground(Cyan)
	getStartedStyle := lipgloss.NewStyle().
		// Background(Background).
		Foreground(Violet).
		Bold(true).
		BorderStyle(lipgloss.DoubleBorder())

	return lipgloss.Place(
		m.BodyDimensions.Width,
		m.BodyDimensions.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			welcomeStyle.Render(welcomeMessage),
			getStartedStyle.Render(getStarted),
		),
	)
}
