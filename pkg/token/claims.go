package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtKey     = []byte("HelloWorld")
	jwtExpires = time.Hour * 2
)

type JWTConfig struct {
	Key     string
	Expires int
}

func SetConf(conf JWTConfig) {
	jwtKey = []byte(conf.Key)
	jwtExpires = time.Duration(conf.Expires) * time.Second
}

// 签名体
type Claims struct {
	Uid int `json:"uid"`
	jwt.StandardClaims
}

func (c *Claims) IsExpired() bool {
	return c != nil && time.Now().Unix() > c.ExpiresAt
}
