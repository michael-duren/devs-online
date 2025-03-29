package views

import (
	"fmt"

	"github.com/michael-duren/tui-chat/ui/models"
)

func Loading(m *models.AppModel) string {
	loadingModel := m.Loading
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", loadingModel.Spinner.View())

	return fmt.Sprint(str, "\n")
}
