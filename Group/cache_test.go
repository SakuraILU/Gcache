package group

import (
	"fmt"
	lru "gcache/Lru"
	"sync"
	"testing"
)

// single goroutine
func TestCache(t *testing.T) {
	kvs := make(map[string]lru.Value)
	num := 100
	for i := 0; i < num; i++ {
		kvs[fmt.Sprintf("%02d", i)] = lru.NewValue([]byte(fmt.Sprintf("%02d", i)))
	}

	gcache := newCache(1000)
	for k, v := range kvs {
		gcache.add(k, v)
	}

	for k, v := range kvs {
		if val, err := gcache.get(k); err != nil {
			t.Errorf("[Error]: %v", err)
		} else if val.ToString() != v.ToString() {
			t.Errorf("[Error]: %v != %v", val, v)
		}
	}
}

// multi goroutine
func TestCache2(t *testing.T) {
	kvs := make(map[string]lru.Value)
	num := 100
	for i := 0; i < num; i++ {
		kvs[fmt.Sprintf("%02d", i)] = lru.NewValue([]byte(fmt.Sprintf("%02d", i)))
	}

	gcache := newCache(1000)
	// multi add
	gnum := 10
	wg := sync.WaitGroup{}
	wg.Add(gnum)
	for i := 0; i < gnum; i++ {
		go func(i int) {
			for k, v := range kvs {
				gcache.add(k, v)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// multi get
	for k, v := range kvs {
		go func(k string, v lru.Value) {
			if val, err := gcache.get(k); err != nil {
				t.Errorf("[Error]: %v", err)
			} else if val.ToString() != v.ToString() {
				t.Errorf("[Error]: %v != %v", val, v)
			}
		}(k, v)
	}
}

// multi goroutine with exceed maxbytes
func TestCache3(t *testing.T) {
	kvs := make(map[string]lru.Value)
	num := 1000
	for i := 0; i < num; i++ {
		kvs[fmt.Sprintf("%02d", i)] = lru.NewValue([]byte(fmt.Sprintf("%02d", i)))
	}

	gcache := newCache(1000)
	// multi add
	gnum := 40
	wg := sync.WaitGroup{}
	wg.Add(gnum)
	for i := 0; i < gnum; i++ {
		go func() {
			for k, v := range kvs {
				gcache.add(k, v)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// multi get
	for i := 0; i < gnum; i++ {
		go func() {
			// the first 3/4 entries will be evicted
			for i := 0; i < num; i++ {
				if i < num*3/4 {
					if _, err := gcache.get(fmt.Sprintf("%02d", i)); err == nil {
						t.Errorf("[Error]: key %v should be evicted", i)
					}
				} else {
					if val, err := gcache.get(fmt.Sprintf("%02d", i)); err != nil {
						t.Errorf("[Error]: %v", err)
					} else {
						v := kvs[fmt.Sprintf("%02d", i)]
						if val.ToString() != v.ToString() {
							t.Errorf("[Error]: %v != %v", val, kvs[fmt.Sprintf("%02d", i)])
						}
					}
				}
			}
		}()
	}
}
