package pkg

import (
	"math/rand"
	"time"
)

var (
	Rand random
	seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const (
	CharsetNumber    = "0123456789"
	CharsetWordUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetWordLower = "abcdefghijklmnopqrstuvwxyz"
)

type random struct{}

// 随机数，最大不包含n本身
func (random) Int(n int) int {
	return seed.Intn(n)
}

func (r random) String(len int) string {
	return r.StringWithCharset(len, CharsetWordUpper+CharsetNumber)
}

func (random) StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}
