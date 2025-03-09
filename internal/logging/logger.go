package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

type ProjectName string

const (
	Server ProjectName = "server"
	Client ProjectName = "client"
)

func NewLogger(project ProjectName) *log.Logger {
	logFileName := fmt.Sprintf("%d_%02d_%02d_log.txt", time.Now().Year(), time.Now().Month(), time.Now().Day())

	logDir := filepath.Join("logs", string(project))
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("could not create log dir %v", err)
	}

	logPath := filepath.Join(logDir, logFileName)
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFile, err = os.Create(logPath)
		if err != nil {
			fmt.Println(err)
			panic("could not create log file")
		}
	}

	logger := log.NewWithOptions(logFile, log.Options{
		Level:           log.DebugLevel,
		ReportTimestamp: true,
		ReportCaller:    true,
	})
	return logger
}
