package golrucache

import (
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// when the data does not exist in cache, we need to get the data from some source
type Getter interface {
	Get(key string) ([]byte, error)
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)
