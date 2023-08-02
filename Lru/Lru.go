package lru

import (
	"container/list"
	"fmt"
	"log"
)

type LruCache struct {
	kvs map[string]*list.Element
	lst *list.List

	curbytes int
	maxbytes int
}

func NewLruCache(maxbytes int) *LruCache {
	return &LruCache{
		kvs:      make(map[string]*list.Element),
		lst:      list.New(),
		curbytes: 0,
		maxbytes: maxbytes,
	}
}

func (l *LruCache) Add(key string, val Value) {
	if e, ok := l.kvs[key]; ok {
		// key remains same, val changes
		oldval := e.Value.(entry).value
		e.Value = entry{key: key, value: val}
		l.lst.MoveToFront(e)
		l.curbytes += val.Len() - oldval.Len()
	} else {
		// new key and val
		e = l.lst.PushFront(entry{key: key, value: val})
		l.kvs[key] = e
		l.curbytes += len(key) + val.Len()
	}

	if l.curbytes > l.maxbytes {
		// evict the oldest cached values if exceed max volume
		l.evict()
	}
	// l.Traverse()
}

func (l *LruCache) Get(key string) (v Value, err error) {
	e, ok := l.kvs[key]
	if !ok {
		err = fmt.Errorf("key %v is not exsit", key)
		return
	}

	v = e.Value.(entry).value
	return
}

func (l *LruCache) Traverse() {
	for e := l.lst.Front(); e != nil; e = e.Next() {
		fmt.Printf("key %v, val (%s, %v)\n", e.Value.(entry).key, e.Value.(entry).key, e.Value.(entry).value)
	}
	// for e := l.lst.Back(); e != nil; e = e.Prev() {
	// 	fmt.Printf("%v\n", e.Value)
	// }
	fmt.Println("ok")
}

func (l *LruCache) evict() {
	fmt.Printf("evict %d--> %d\n", l.curbytes, l.maxbytes)
	for l.curbytes > l.maxbytes {
		fmt.Printf("size %d\n", l.lst.Len())
		// get the back one
		e := l.lst.Back()
		l.Traverse()
		log.Printf("[EVICT]: %v\n", e.Value)
		kv := e.Value.(entry)
		l.curbytes -= kv.value.Len() + len(kv.key)
		if l.curbytes < 0 {
			log.Fatal("curbytes must >= 0")
		}

		l.lst.Remove(e)
		delete(l.kvs, kv.key)
	}
}
