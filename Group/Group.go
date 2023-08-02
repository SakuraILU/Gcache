package group

import (
	"fmt"
	lru "gcache/Lru"
)

type Group struct {
	name       string
	gcache     *cache
	getter     Getter
	peerpicker PeerPicker
}

// var (
// 	groups map[string]*Group = make(map[string]*Group)
// 	glk    sync.RWMutex
// )

func NewGroup(name string, maxbytes int, getter Getter, peerpicker PeerPicker) (g *Group) {
	g = &Group{
		name:       name,
		gcache:     newCache(maxbytes),
		getter:     getter,
		peerpicker: peerpicker,
	}
	return
}

// func GetGroup(name string) (g *Group, err error) {
// 	glk.RLock()
// 	defer glk.RUnlock()

// 	g, ok := groups[name]
// 	if !ok {
// 		err = fmt.Errorf("[ERROR] group %s is not exist", name)
// 	}
// 	return
// }

func (g *Group) Get(key string) (v lru.Value, err error) {
	if key == "" {
		err = fmt.Errorf("[Error/Group] empty key")
		return
	}

	// check wehter gcache has this value?
	if v, err = g.gcache.get(key); err == nil {
		return
	}

	if v, err = g.loadRemote(key); err == nil {
		return
	}

	if v, err = g.loadLocal(key); err == nil {
		g.gcache.add(key, v)
		return
	}

	return
}

func (g *Group) loadRemote(key string) (v lru.Value, err error) {
	// get peer
	peer, err := g.peerpicker.PickPeer(key)
	if err != nil {
		return
	}
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return
	}
	v = lru.NewValue(bytes)
	return
}

func (g *Group) loadLocal(key string) (v lru.Value, err error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return
	}
	// get this value from local by getter
	v = lru.NewValue(bytes)
	return
}
