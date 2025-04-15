package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/internal/logging"
	"github.com/michael-duren/tui-chat/internal/server"
	"github.com/michael-duren/tui-chat/ui"
	"github.com/spf13/cobra"
)

func gracefulShutdown(server *server.Server, done chan bool) {
	// listen from interrupt signal from os
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// wait for interrupt
	<-ctx.Done()

	server.ShutdownSockets()

	done <- true
}

func runServe(external bool, addr string, port int, secret string) {
	logger := logging.NewLogger("server")
	logger.Info("sever address: ", addr)
	if external {
		// TODO: Figure out
		logger.Info("serving on external port")
	}

	server := server.NewServer(port, logger, secret)
	done := make(chan bool, 1)
	go gracefulShutdown(server, done)
	err := server.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %v:", err))
	}

	<-done
	log.Print("Graceful shutdown complete.")
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
	var port int
	var addr string
	var secret string
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the t-chat server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("starting server, serving external: %t", external)
			runServe(external, addr, port, secret)
		},
	}
	serveCmd.Flags().BoolVar(&external, "external", false, "Enable external access")
	serveCmd.Flags().StringVar(&addr, "address", "", "Specifies the host address")
	serveCmd.Flags().IntVar(&port, "port", 8080, "Specifies the port")
	serveCmd.Flags().StringVar(&secret, "chat secret", "", "Specifies the secret to get into the chat")

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
