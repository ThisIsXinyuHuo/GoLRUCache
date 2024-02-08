package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

// test Get
func TestGet(t *testing.T) {
	// create a lru
	lru := New(int64(0), nil)
	// add a k-v pair
	lru.Add("key1", String("1234"))
	// see if k-v pair exits behave correctly in our lru
	if val, exist := lru.Get("key1"); !exist || string(val.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	// test for a non-exits pair
	if _, exist := lru.Get("key2"); exist {
		t.Fatalf("cache miss key2 failed")
	}
}

// test DeleteOldest
func TestDeleteOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"

	// create a lru with designed capacity
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	// if lru does not move the least recent one, or store more elements
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("DeleteOldest key1 failed")
	}
}

// test the callback
func TestOnEvicted(t *testing.T) {
	// slice with 0 values
	keys := make([]string, 0)
	// desgin the call back
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("OnEvicted failed, expect keys equals to %s", expect)
	}
}
