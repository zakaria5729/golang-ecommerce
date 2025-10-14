package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
)

var (
	Logger *slog.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		logsDir := filepath.Join(utils.GetProjectRootPath(), "logs")
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
	today := timeutil.NowUTC().Format(c.LogFileFormat)
	logsDir := filepath.Join(utils.GetProjectRootPath(), "logs")
	return filepath.Join(logsDir, fmt.Sprintf("app-%s.log", today))
}
