package cache

import (
	"sync"
	"time"
)

var (
	DefaultExpiration = time.Hour * 24 * 7
)

type cache[KT comparable, VT any] struct {
	data sync.Map
}

type item[VT any] struct {
	Value    VT
	ExpireAt int64
}

func New[KT comparable, VT any]() *cache[KT, VT] {
	return &cache[KT, VT]{}
}

func (c *cache[KT, VT]) Get(k KT) (VT, bool) {
	x, ok := c.data.Load(k)

	if ok {
		it := x.(item[VT])
		if it.ExpireAt >= time.Now().Unix() {
			return it.Value, true
		}
	}

	return *new(VT), false
}

func (c *cache[KT, VT]) Set(k KT, v VT) {
	c.SetEx(k, v, DefaultExpiration)
}

func (c *cache[KT, VT]) SetEx(k KT, v VT, d time.Duration) {
	c.data.Store(k, item[VT]{
		Value:    v,
		ExpireAt: time.Now().Add(d).Unix(),
	})
}

func (c *cache[KT, VT]) Has(k KT) bool {
	x, ok := c.data.Load(k)
	if !ok {
		return false
	}
	it := x.(item[VT])
	return it.ExpireAt >= time.Now().Unix()
}

func (c *cache[KT, VT]) Del(k KT) {
	c.data.Delete(k)
}
