package helpers

import (
	"encoding/base32"
)

// EncodeTOTP encode TOTP secret in base32
func EncodeTOTP(secret string) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(secret))
}
