package util

import (
	"math/rand"
	"time"
)

var (
	seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const (
	CharsetNumber    = "0123456789"
	CharsetWordUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetWordLower = "abcdefghijklmnopqrstuvwxyz"
)

// 随机数，最大不包含n本身
func RandInt(n int) int {
	return seed.Intn(n)
}

func RandString(len int) string {
	return RandStringWithCharset(len, CharsetWordUpper+CharsetNumber)
}

func RandStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}
