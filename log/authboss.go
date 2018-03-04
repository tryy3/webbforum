package log

import (
	"io"
	"strings"

	"github.com/apex/log"
)

func NewAuthbossLogger() io.Writer {
	return &AuthbossLogger{}
}

type AuthbossLogger struct{}

func (a *AuthbossLogger) Write(p []byte) (int, error) {
	log.Info(strings.TrimRight(string(p), "\n"))
	return 0, nil
}
