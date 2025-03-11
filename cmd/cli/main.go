package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/internal/logging"
	"github.com/michael-duren/tui-chat/ui"
	"github.com/spf13/cobra"
)

func runServe(external bool) {
	logger := logging.NewLogger("server")
	logger.Info("Starting server")
	if external {
		logger.Info("serving on external port")
	}
}

func runChat() {
	app := ui.InitialModel()
	app.Logger.Info("Starting Client")
	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		app.Logger.Errorf("oh no looks like we're ngmi. famous last words: %v", err)
		os.Exit(1)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "t-chat",
		Short: "T-Chat CLI Application",
	}

	var external bool
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the t-chat server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("starting server, serving external: %t", external)
			runServe(external)
		},
	}
	serveCmd.Flags().BoolVar(&external, "external", false, "Enable external access")

	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Start a chat session",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("starting chat session")
			runChat()
		},
	}

	rootCmd.AddCommand(serveCmd, chatCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
