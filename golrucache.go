package golrucache

import (
	"fmt"
	"log"
	"sync"
)

// user interaction, retrieve the data

// each group is a logical namespace for cached data
// it contains associated data that is spread over various cache entries

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// when the data does not exist in cache, we need to get the data from some source
type Getter interface {
	Get(key string) ([]byte, error)
}

// create a class which implements the Get method
// In some cases, we can cast the user provided getter function to GetterFunc type
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mutex.Lock()
	defer mutex.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	// add to groups map
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mutex.RLock()
	g := groups[name]
	mutex.RUnlock()
	return g
}

// get value from the cache from some data source
func (g *Group) getFromSource(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)

	//cannot find the data from the provided source
	if err != nil {
		return ByteView{}, err
	}

	//return a copy
	value := ByteView{b: cloneBytes(bytes)}
	// add the retuned value to the cache
	g.populateCache(key, value)
	return value, nil
}

// get the data
// late will implement this methond in the distributed setting
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getFromSource(key)
}

// get a value from cache
func (g *Group) Get(key string) (ByteView, error) {
	// the key isn't provided
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// check if the key is in the cache
	if value, exist := g.mainCache.get(key); exist {
		log.Println("[GoLRUCache] hit")
		return value, nil
	}

	// if the key does not exist in the cache, query the data source
	return g.load(key)
}

// add the value to the cache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
