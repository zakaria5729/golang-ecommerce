package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

var (
	logger          zerolog.Logger
	currentLogFile  *os.File
	asyncFileWriter io.WriteCloser
	loggerMutex     sync.Mutex
)

func init() {
	if err := os.MkdirAll(utils.GetLogFolderPath(), 0755); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	createNewAsyncWritersLogger()
}

func requireCurrentLogFile() {
	if currentLogFile == nil {
		panic("Failed to create log file")
	}
}

func createFileIfNotExists() {
	if currentLogFile != nil && filepath.Base(currentLogFile.Name()) == getLogFileName() {
		return
	}

	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	CloseLogFile(false)
	createNewAsyncWritersLogger()
}

func createNewAsyncWritersLogger() {
	createNewLogFile()
	requireCurrentLogFile()

	var logWriters []io.Writer
	var activeProfile = config.GetActiveProfile()

	asyncFileWriter = diode.NewWriter(currentLogFile, 10000, 10*time.Millisecond, func(missed int) {
		fmt.Fprintf(os.Stderr, "File log dropped %d messages\n", missed)
	})
	logWriters = append(logWriters, asyncFileWriter)

	if activeProfile == c.EnvLocal || activeProfile == c.EnvDev {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 T 15:04:05 Z07:00",
			NoColor:    false,
		}
		logWriters = append(logWriters, consoleWriter)
	}

	output := zerolog.MultiLevelWriter(logWriters...)

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "time_utc"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"

	logger = zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}

func CloseLogFile(writeLog bool) {
	if currentLogFile != nil {
		if writeLog {
			Info("Log file closed successfully")
		}

		if asyncFileWriter != nil {
			asyncFileWriter.Close()
			asyncFileWriter = nil
		}

		time.Sleep(200 * time.Millisecond)
		currentLogFile.Close()
		currentLogFile = nil
	}
}

func createNewLogFile() {
	logFilePath := filepath.Join(utils.GetLogFolderPath(), getLogFileName())
	if logFilePath == "" {
		panic("Log file path can't be empty")
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}

	currentLogFile = logFile
}

func getLogFileName() string {
	today := timeutil.NowUTC().Format(c.LogFileFormat)
	return fmt.Sprintf("app-%s."+c.LogFileExt, today)
}

func applyLog(ev *zerolog.Event, fields []any) *zerolog.Event {
	if len(fields) == 0 {
		return ev
	}

	if len(fields)%2 == 1 {
		last := fields[len(fields)-1]
		ev = ev.Interface(fmt.Sprintf("%v", last), "!BADKEY")
		fields = fields[:len(fields)-1]
	}

	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		val := fields[i+1]
		ev = ev.Interface(key, val)
	}

	return ev
}

func Info(msg string, fields ...any) {
	createFileIfNotExists()
	applyLog(logger.Info(), fields).Msg(msg)
}

func Error(msg string, fields ...any) {
	createFileIfNotExists()
	applyLog(logger.Error(), fields).Msg(msg)
}

func Warn(msg string, fields ...any) {
	createFileIfNotExists()
	applyLog(logger.Warn(), fields).Msg(msg)
}

func Fatal(msg string, fields ...any) {
	createFileIfNotExists()
	applyLog(logger.Fatal(), fields).Msg(msg)
}
