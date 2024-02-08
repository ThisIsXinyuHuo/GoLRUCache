package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	// if maxByte == 0, unlimited capacity
	maxBytes int64
	nbytes   int64
	// Front: most recent
	// End: Lease recent
	ll    *list.List
	cache map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// return the number of elements in the list
func (c *Cache) Len() int {
	return c.ll.Len()
}

// constrcuter for Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get: Search for a key and return the value
func (c *Cache) Get(key string) (value Value, ok bool) {
	// if the key exits:
	if val, exist := c.cache[key]; exist {
		// The most recent node should be moved to front
		c.ll.MoveToFront(val)
		// kv => list's element => element's value => convert it to *entry
		kv := val.Value.(*entry)
		//return the value
		return kv.value, true
	}
	//not exists, return the default value
	return
}

func (c *Cache) DeleteOldest() {
	//last element is the least recent
	val := c.ll.Back()
	if val != nil {
		// delete the last element
		c.ll.Remove(val)
		// kv => list's element => element's value => convert it to *entry
		kv := val.Value.(*entry)
		// delte element in the cache map
		delete(c.cache, kv.key)
		// update the capacity
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// callback
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add: adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	// if the key already exists:
	if val, exist := c.cache[key]; exist {
		// The most recent node should be moved to front
		c.ll.MoveToFront(val)
		// kv => list's element => element's value => convert it to *entry
		kv := val.Value.(*entry)
		// update the used capacity
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		// update the value
		kv.value = value
	} else {
		// if the key does not exist
		// create a key and move it to front
		val := c.ll.PushFront(&entry{key, value})
		c.cache[key] = val
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.DeleteOldest()
	}
}
