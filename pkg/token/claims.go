package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	JwtConfig struct {
		Key     string
		Expires int
	}
)

// 签名体
type Claims struct {
	Uid int `json:"uid"`
	jwt.StandardClaims
}

func (c *Claims) IsExpired() bool {
	return c != nil && time.Now().Unix() > c.ExpiresAt
}
