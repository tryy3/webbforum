package utils

import (
	"encoding/base64"

	"github.com/gorilla/securecookie"
)

// GenerateRandomKey is a wrapper over securecookie.GenerateRandomKey to generate a string
// rather then a byte slice
func GenerateRandomKey(len int) string {
	// Generate random key from securecookie
	data := securecookie.GenerateRandomKey(len)
	if data == nil {
		return ""
	}

	// Encode the random key
	return base64.StdEncoding.EncodeToString(data)
}
