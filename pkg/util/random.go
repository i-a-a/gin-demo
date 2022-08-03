package util

import (
	"math/rand"
	"time"
)

type charset string

const (
	CharsetNumber    charset = "0123456789"
	CharsetWordUpper charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetWordLower charset = "abcdefghijklmnopqrstuvwxyz"
)

var (
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 随机数，最大不包含n本身
func RandomInt(n int) int {
	return seededRand.Intn(n)
}

func RandomString(len int) string {
	return RandomStringWithCharset(len, CharsetWordUpper+CharsetNumber)
}

// 可以打组合拳(😂)
func RandomStringWithCharset(length int, charset charset) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
