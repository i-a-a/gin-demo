package token

import (
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	charsetLen            = len(charset)
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Token struct {
	// 有效期较长。后端需要存储，自己限制其过期时间。
	RefreshToken string `json:"refresh_token"`
	// 有效期短。是个JWT， 后端不需要存储
	AccessToken string `json:"access_token"`
	// JWT的过期时间，前端存储，判断是否需要刷新
	ExpireAt int64 `json:"expire_at"`
}

func GenerateTokens(uid int) Token {
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[seededRand.Intn(charsetLen)]
	}
	var token Token
	token.RefreshToken = string(b)
	token.AccessToken, token.ExpireAt = generateJWT(uid)
	return token
}

func generateJWT(uid int) (string, int64) {
	if uid == 0 {
		logrus.Error("生成JWT时，uid为0")
	}
	claims := Claims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpires).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return token, claims.ExpiresAt
}

// 解析JWT，注意过期不是错误，自己判断 IsExpired()
func ParseJWT(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	var claims *Claims
	if tokenClaims != nil {
		claims = tokenClaims.Claims.(*Claims)
	}

	return claims, err
}
