package golrucache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

// test if customized callback function (getter) works
func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")

	if value, _ := f.Get("key"); !reflect.DeepEqual(value, expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Alice":  "160",
	"Nicole": "170",
	"Kate":   "165",
}

func TestGet(t *testing.T) {
	// track the number of times visiting db
	loadCounts := make(map[string]int, len(db))

	grp := NewGroup("height", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if value, exist := db[key]; exist {
				if _, exist := loadCounts[key]; !exist {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(value), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	for key, value := range db {
		// test load value from the callback function
		if view, err := grp.Get(key); err != nil || view.String() != value {
			t.Fatalf("failed to get value of %s", key)
		}

		// test if cache works
		if _, err := grp.Get(key); err != nil || loadCounts[key] > 1 {
			t.Fatalf("cache %s miss", key)
		}
	}

	// test unkonwn keys
	if view, err := grp.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown keys should be empty, but got %s", view)
	}
}
