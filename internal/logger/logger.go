package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func WithLogLevel(level string) logrus.Level {
	switch level {
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "debug":
		return logrus.DebugLevel
	default:
		log.Print("applying default log level ERROR")
		return logrus.ErrorLevel
	}
}

func WithOutput(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "file":
		f, err := os.OpenFile("logs/logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("could not open log file. Path: logs")
		}
		return f
	default:
		log.Print("applying default output - stdout")
		return os.Stdout
	}
}

const (
	FormatterTypeJson = "json"
	FormatterTypeText = "text"
)

func WithFormatter(formatter string) logrus.Formatter {
	switch strings.ToLower(strings.TrimSpace(formatter)) {
	case FormatterTypeJson:
		return &UTCFormatter{&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			}}}
	case FormatterTypeText:
		return &UTCFormatter{&logrus.TextFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			}}}
	default:
		panic(fmt.Errorf("%v is an invalid formatter type", formatter))
	}
}

type UTCFormatter struct {
	logrus.Formatter
}

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

type FiberLogWriter struct{}

func (l *FiberLogWriter) Write(p []byte) (n int, err error) {
	logrus.Info(string(p))
	return len(p), nil
}
