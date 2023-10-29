package cache

import (
	"time"
)

type Cleaner struct {
	Interval time.Duration
}

type item struct {
	value    string
	creation int64
	rank     int
}

type LRUCache struct {
	cacheMap map[int]item
	ttl      time.Duration
	cleaner  *Cleaner
	capacity int
	rank     int
}

func New(capacity int, ttl time.Duration, cleanUpInterval time.Duration) *LRUCache {
	newLRUCache := LRUCache{
		cacheMap: map[int]item{},
		ttl:      ttl,
		capacity: capacity,
		rank:     0,
	}

	if cleanUpInterval > 0 {
		clean(cleanUpInterval, &newLRUCache)
	}

	return &newLRUCache
}

func (c *LRUCache) Get(key int) string {
	rankValue := c.rank
	item := c.cacheMap[key]
	item.rank = rankValue
	c.rank++
	return item.value
}

func (c *LRUCache) Put(key int, value string) {
	if len(c.cacheMap) == c.capacity {
		var minRankKey int = -1
		for k, _ := range c.cacheMap {
			if minRankKey == -1 || c.cacheMap[k].rank < c.cacheMap[minRankKey].rank {
				minRankKey = k
			}
		}
		delete(c.cacheMap, minRankKey)
	}
	c.cacheMap[key] = item{value: value, creation: time.Now().UnixNano(), rank: c.rank}
	c.rank++
}

func clean(cleanUpInterval time.Duration, cache *LRUCache) {
	cleaner := &Cleaner{
		Interval: cleanUpInterval,
	}

	cache.cleaner = cleaner
	go cleaner.Cleaning(cache)
}

func (c *Cleaner) Cleaning(cache *LRUCache) {
	ticker := time.NewTicker(c.Interval)

	for {
		select {
		case <-ticker.C:
			cache.purge()
		}
	}
}

func (c *LRUCache) purge() {
	now := time.Now().UnixNano()
	ttl := c.ttl.Nanoseconds()
	for key, data := range c.cacheMap {
		if data.creation+ttl < now {
			delete(c.cacheMap, key)
		}
	}
}
