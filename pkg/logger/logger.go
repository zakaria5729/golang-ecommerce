package logger

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var Logger *slog.Logger

func init() {
	initLogger()
}

func initLogger() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	logFile := getLogFileName()
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	Logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func Init() {
	initLogger()
}

func IsInitialized() bool {
	return Logger != nil
}

func getLogFileName() string {
	today := time.Now().Format("2006-01-02")
	return filepath.Join("logs", fmt.Sprintf("app-%s.log", today))
}

func SetLevel(level slog.Level) {
	logFile := getLogFileName()
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to create log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	Logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	}))
}
