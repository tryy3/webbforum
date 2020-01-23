package log

import (
	"io"
	"strings"

	"github.com/apex/log"
)

// NewAuthbossLogger creates a new AuthBossLogger
func NewAuthbossLogger() io.Writer {
	return &AuthbossLogger{}
}

// AuthBossLogger takes care of logging all of the authboss related stuff
type AuthbossLogger struct{}

// Write sends a log message to the logger
func (a *AuthbossLogger) Write(p []byte) (int, error) {
	log.Info(strings.TrimRight(string(p), "\n"))
	return 0, nil
}
