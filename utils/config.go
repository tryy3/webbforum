package utils

import "time"

// Config is the general config for everything
type Config struct {
	// HTTP
	HTTPIP   string
	HTTPPort int

	// XSRF
	XSRFName string

	// Cookie Settings
	CookieStoreKey []byte
	CookieExpiry   time.Duration

	// Session Settings
	SessionStoreKey []byte
	SessionName     string

	// Email
	SMTPHost     string
	SMTPUsername string
	SMTPPassword string
	SMTPIdentity string
	SMTPEmail    string
	SMTPName     string
}
