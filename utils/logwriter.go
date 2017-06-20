package utils

import (
	"io"
	"log"
	"strings"
)

type loggerWriter struct {
	logger *log.Logger
}

func (lw loggerWriter) Write(p []byte) (int, error) {
	str := strings.Trim(string(p), "\x00")
	l := len(str)
	lw.logger.Printf("%s", strings.Trim(string(p), "\x00"))
	return l, nil
}

// LogWriter turns a *log.Logger into an io.Writer
func LogWriter(logger *log.Logger) io.Writer {
	return loggerWriter{
		logger: logger,
	}
}
