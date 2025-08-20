package configs

import (
	"log"
	"os"
	"time"
)

type Level int

const (
	INFO Level = iota
	WARN
	ERROR
	DEBUG
)

var LevelColors = map[Level]string {
	INFO:  "\033[32m", // Green
	WARN:  "\033[33m", // Yellow
	ERROR: "\033[31m", // Red
	DEBUG: "\033[36m", // Cyan
}

type CustomLogger struct {
	out *log.Logger
}

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		out: log.New(os.Stdout, "", 0),
	}
}

func (l *CustomLogger) log(level Level, msg string) {
	color := LevelColors[level]
	levelStr := []string{"INFO", "WARN", "ERROR", "DEBUG"}[level]
	timestamp := time.Now().Format(time.RFC3339)

	l.out.Printf("%s[%s] %s %-5s %s%s",
		color,
		timestamp,
		"|",
		levelStr,
		msg,
		"\033[0m", // Reset
	)
}

func (l *CustomLogger) Info(msg string)  { l.log(INFO, msg) }
func (l *CustomLogger) Warn(msg string)  { l.log(WARN, msg) }
func (l *CustomLogger) Error(msg string) { l.log(ERROR, msg) }
func (l *CustomLogger) Debug(msg string) { l.log(DEBUG, msg) }
