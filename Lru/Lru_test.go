package lru

import (
	"fmt"
	"testing"
)

// easy test
func TestLru1(t *testing.T) {
	// define several values to be cached
	num := 10
	keys := make([]string, num)
	values := make([]Value, num)
	for i := 0; i < num; i++ {
		keys[i] = string("k" + fmt.Sprint(i))
		values[i] = Value{bytes: []byte(string("v" + fmt.Sprint(i)))}
	}

	// define a lru cache
	lru := NewLruCache(1000)
	// add values to lru cache
	for i := 0; i < num; i++ {
		lru.Add(keys[i], values[i])
	}

	// get values from lru cache
	for i := 0; i < num; i++ {
		v, err := lru.Get(keys[i])
		if err != nil {
			t.Errorf("key %v is not exsit", keys[i])
		}
		if v.ToString() != values[i].ToString() {
			t.Errorf("value %v is not equal to %v", v.ToString(), values[i].ToString())
		}
	}
}

// exceed maxbytes test
func TestLru2(t *testing.T) {
	// define several values to be cached
	num := 100
	keys := make([]string, num)
	values := make([]Value, num)
	for i := 0; i < num; i++ {
		// key : 01,02...11,12,99
		keys[i] = fmt.Sprintf("%02d", i)
		values[i] = Value{bytes: []byte(fmt.Sprintf("%02d", i))}
	}
	// every entry {xx, xx} takes 4 bytes
	// so maxbytes is 400
	// but here we set maxbytes to 100
	// so the first 75 entries will be evicted
	// and the last 25 entries will be cached
	maxbytes := 100

	// define a lru cache
	lru := NewLruCache(maxbytes)
	// add values to lru cache
	for i := 0; i < num; i++ {
		lru.Add(keys[i], values[i])
	}

	// get values from lru cache
	for i := 0; i < num; i++ {
		v, err := lru.Get(keys[i])
		if i < maxbytes*3/4 {
			// the first 75 entries will be evicted
			if err == nil {
				t.Errorf("key %v should be evicted", keys[i])
			}
		} else {
			// the last 25 entries will be cached
			if err != nil {
				t.Errorf("key %v is not exsit", keys[i])
			}
			if v.ToString() != values[i].ToString() {
				t.Errorf("value %v is not equal to %v", v.ToString(), values[i].ToString())
			}
		}
	}
}
