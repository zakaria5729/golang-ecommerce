package logger

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var Logger *slog.Logger

func init() {
	initLogger()
}

func initLogger() {
	logsDir := getProjectRoot() + "/logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
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
	logsDir := getProjectRoot() + "/logs"
	return filepath.Join(logsDir, fmt.Sprintf("app-%s.log", today))
}

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// Go up from pkg/logger/logger.go to project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
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
