package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/easy-comerce/backend/pkg/timeutil"
)

var (
	Logger *slog.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
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
	})
}

func getLogFileName() string {
	today := timeutil.NowUTC().Format("2006-01-02")
	logsDir := getProjectRoot() + "/logs"
	return filepath.Join(logsDir, fmt.Sprintf("app-%s.log", today))
}

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// Go up from pkg/logger/logger.go to project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
