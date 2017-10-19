package lib

import (
	"encoding/hex"

	"crypto/sha1"

	"github.com/satori/go.uuid"
)

func RndKey() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}

func Encode(secret string) string {
	var hash = sha1.Sum([]byte(secret))
	return hex.EncodeToString(hash[:])
}
