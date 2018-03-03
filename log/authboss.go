package log

import (
	"io"

	"github.com/labstack/gommon/log"
)

func NewAuthbossLogger() io.Writer {
	return &AuthbossLogger{}
}

type AuthbossLogger struct{}

func (a *AuthbossLogger) Write(p []byte) (int, error) {
	log.Info(p)
	return 0, nil
}
