package main

import "github.com/michael-duren/tui-chat/internal/logging"

func main() {
	logger := logging.NewLogger()
	logger.Info("Starting client")
}
