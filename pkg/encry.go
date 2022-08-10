package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var (
	Encry encry
)

type encry struct{}

func (encry) Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (encry) HmacSha256(s string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
