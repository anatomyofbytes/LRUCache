package lru

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	capacity int
	ll       *list.List
	cache    map[string]*list.Element
	mu       sync.RWMutex
}

type entry struct {
	key   string
	value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		ll:       list.New(),
		// попробовать добавить capacicty
		cache: make(map[string]*list.Element),
	}
}

func (c *LRUCache) Add(key string, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.ll = list.New()
	}
	ee, ok := c.cache[key]
	if ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return ok
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.capacity != 0 && c.ll.Len() > c.capacity {
		c.removeOldest()
	}
	return ok
}

func (c *LRUCache) Get(key string) (value interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cache == nil {
		return nil, false
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LRUCache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return true
	}
	ele, ok := c.cache[key]
	if ok {
		c.removeElement(ele)
	}
	return ok
}

func (c *LRUCache) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}
