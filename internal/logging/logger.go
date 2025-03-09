package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func NewLogger() *log.Logger {
	logFileName := fmt.Sprintf("%d_%02d_%02d_log.txt", time.Now().Year(), time.Now().Month(), time.Now().Day())
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFile, err = os.Create(logFileName)
		if err != nil {
			panic("could not create log file")
		}
	}

	logger := log.New(logFile)
	return logger
}
