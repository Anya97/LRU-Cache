package cache

import (
	"time"
)

type Cleaner struct {
	Interval time.Duration
}

type Node struct {
	previousElemLink *Node
	nextElemLink     *Node
	key              *string
}

type item struct {
	value    interface{}
	creation time.Time
	node     *Node
}

type LRUCache struct {
	cacheMap      map[string]*item
	ttl           time.Duration
	cleaner       *Cleaner
	capacity      int
	firstElemLink *Node
	lastElemLink  *Node
}

func New(capacity int, ttl time.Duration, cleanUpInterval time.Duration) *LRUCache {
	newLRUCache := LRUCache{
		cacheMap: map[string]*item{},
		ttl:      ttl,
		capacity: capacity,
	}

	if cleanUpInterval > 0 {
		clean(cleanUpInterval, &newLRUCache)
	}

	return &newLRUCache
}

func (c *LRUCache) Get(key string) interface{} {
	if element := c.get(key); element != nil {
		return element.value
	}
	return nil
}

func (c *LRUCache) get(key string) *item {
	itemInCache, ok := c.cacheMap[key]
	if !ok {
		return nil
	}
	c.moveToBack(itemInCache)
	return itemInCache
}

func (c *LRUCache) Put(key string, value interface{}) {
	if element := c.get(key); element != nil {
		c.cacheMap[key] = &item{value: value, creation: time.Now(), node: element.node}
		return
	}
	if len(c.cacheMap) == c.capacity {
		c.removeFirst()
	}
	valueItem := item{value: value, creation: time.Now(), node: &Node{key: &key}}
	c.cacheMap[key] = &valueItem
	c.pushBack(valueItem.node)
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
	for key, data := range c.cacheMap {
		if data.creation.Add(c.ttl).After(time.Now()) {
			c.unlinkNode(data.node)
			delete(c.cacheMap, key)
		}
	}
}

func (c *LRUCache) moveToBack(element *item) {
	if c.lastElemLink == element.node {
		return
	}
	if c.firstElemLink == element.node {
		c.firstElemLink = c.firstElemLink.nextElemLink
		c.firstElemLink.previousElemLink = nil
	} else {
		element.node.previousElemLink.nextElemLink = element.node.nextElemLink
		element.node.nextElemLink.previousElemLink = element.node.previousElemLink
	}
	c.lastElemLink.nextElemLink = element.node
	element.node.previousElemLink = c.lastElemLink
	element.node.nextElemLink = nil
	c.lastElemLink = element.node
}

func (c *LRUCache) pushBack(node *Node) {
	if c.firstElemLink == nil {
		c.firstElemLink = node
		c.lastElemLink = node
	} else {
		c.lastElemLink.nextElemLink = node
		node.previousElemLink = c.lastElemLink
		c.lastElemLink = node
	}
}

func (c *LRUCache) removeFirst() {
	delete(c.cacheMap, *c.firstElemLink.key)
	c.firstElemLink = c.firstElemLink.nextElemLink
	c.firstElemLink.previousElemLink = nil
}

func (c *LRUCache) unlinkNode(node *Node) {
	if node.previousElemLink != nil {
		node.previousElemLink.nextElemLink = node.nextElemLink
	} else {
		c.firstElemLink = node.nextElemLink
	}
	if node.nextElemLink != nil {
		node.nextElemLink.previousElemLink = node.previousElemLink
	} else {
		c.lastElemLink = node.previousElemLink
	}
}
