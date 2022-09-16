package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var (
	jwtKey     = []byte("HelloWorld")
	jwtExpires = time.Second * 7200
)

func init() {
	var conf struct {
		Key     string
		Expires int
	}
	viper.UnmarshalKey("token", &conf)
	jwtKey = []byte(conf.Key)
	jwtExpires = time.Second * time.Duration(conf.Expires)
}

// 签名体
type Claims struct {
	Uid int `json:"uid"`
	// ...
	jwt.StandardClaims
}

func (c *Claims) IsExpired() bool {
	return c != nil && time.Now().Unix() > c.ExpiresAt
}
