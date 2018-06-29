package main

type Cache struct {
	cache map[int64]struct{}
}

func NewCache() *Cache {
	c := new(Cache)
	c.cache = make(map[int64]struct{})
	return c
}

func (c *Cache) Exists(k int64) bool {
	if _, ok := c.cache[k]; ok {
		return true
	}
	c.cache[k] = struct{}{}
	return false
}

type StringCache struct {
	cache map[string]int64
}

func NewStringCache() *StringCache {
	c := new(StringCache)
	c.cache = make(map[string]int64, 1000)
	return c
}

func (c *StringCache) Exists(k string) (int64, bool) {
	if value, ok := c.cache[k]; ok {
		return value, true
	}
	return 0, false
}

func (c *StringCache) Add(k string, v int64) {
	c.cache[k] = v
}
