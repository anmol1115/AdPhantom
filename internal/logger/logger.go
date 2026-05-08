package logger

import (
	"errors"
	"log"
	"os"

	"github.com/anmol1115/AdPhantom/internal/config"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
)

type Logger struct {
	Level    LogLevel
	MaxSize  int
	MaxCount int
	logger   *log.Logger
	file     *os.File
}

func convertLogLevel(level string) (LogLevel, error) {
	switch level {
	case "debug":
		return Debug, nil
	case "error":
		return Error, nil
	case "info":
		return Info, nil
	}
	return Debug, errors.New("Invalid log level")
}

func Init(loggerConfig *config.LoggingConfig, filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	l := log.New(file, "", log.Ldate|log.Ltime|log.Lmsgprefix)
	ll, err := convertLogLevel(loggerConfig.Level)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Level:    ll,
		MaxSize:  loggerConfig.Size,
		MaxCount: loggerConfig.Count,
		logger:   l,
		file:     file,
	}, nil
}

func (l *Logger) Debug(msg ...string) {
	if l.Level > Debug {
		return
	}

	l.logger.Println("[DEBUG]", msg)
}

func (l *Logger) Info(msg ...string) {
	if l.Level > Info {
		return
	}
	l.logger.Println("[WARN]", msg)
}

func (l *Logger) Error(msg ...string) {
	if l.Level > Error {
		return
	}
	l.logger.Println("[ERROR]", msg)
}

func (l *Logger) Close() {
	l.file.Close()
}
