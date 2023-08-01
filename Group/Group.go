package group

import (
	"fmt"
	lru "gcache/Lru"
	"sync"
)

type Group struct {
	gcache *cache
	getter Getter
}

var (
	groups map[string]*Group = make(map[string]*Group)
	glk    sync.RWMutex
)

func NewGroup(name string, maxbytes int, getter Getter) (g *Group) {
	glk.Lock()
	defer glk.Unlock()

	g = &Group{
		gcache: newCache(maxbytes),
		getter: getter,
	}
	groups[name] = g
	return
}

func GetGroup(name string) (g *Group, err error) {
	glk.RLock()
	defer glk.RUnlock()

	g, ok := groups[name]
	if !ok {
		err = fmt.Errorf("[ERROR] group %s is not exist", name)
	}
	return
}

func (g *Group) Get(key string) (v lru.Value, err error) {
	if key == "" {
		err = fmt.Errorf("[Error/Group] empty key")
		return
	}

	// check wehter gcache has this value?
	if v, err = g.gcache.get(key); err != nil {
		v, err = g.loadByGetter(key)
		if err != nil {
			return
		}

		g.gcache.add(key, v)
	}

	return
}

func (g *Group) loadByGetter(key string) (v lru.Value, err error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return
	}
	// get this value from local by getter
	v = lru.NewValue(bytes)
	return
}
