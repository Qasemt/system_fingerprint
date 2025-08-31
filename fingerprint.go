// fingerprint.go
package main

import (
	"crypto/sha256"
	"strings"

	"fmt"
)

func getSystemFingerprint() string {
	components, err := getPlatformComponents()
	if err != nil || len(components) == 0 {
		return ""
	}

	// 🔥 ادغام با ||
	joined := strings.Join(components, "||")

	// SHA-256
	hash := sha256.Sum256([]byte(joined))
	return fmt.Sprintf("%x", hash)
}
