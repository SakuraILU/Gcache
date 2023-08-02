package group

import (
	"fmt"
	lru "gcache/Lru"
	singleflight "gcache/SingleFlight"
	"time"
)

type Group struct {
	name       string
	gcache     *cache
	peerpicker PeerPicker
	getter     Getter
	sf         *singleflight.SinlgeFlight
}

func NewGroup(name string, maxbytes int, getter Getter, peerpicker PeerPicker) (g *Group) {
	g = &Group{
		name:       name,
		gcache:     newCache(maxbytes),
		peerpicker: peerpicker,
		getter:     getter,
		sf:         singleflight.NewSingleFlight(),
	}
	return
}

func (g *Group) Get(key string) (val lru.Value, err error) {
	if key == "" {
		err = fmt.Errorf("[Error/Group] empty key")
		return
	}

	v, err := g.sf.Do(key, func() (v interface{}, err error) {
		// check wehter gcache has this value?
		if v, err = g.gcache.get(key); err == nil {
			return
		}

		if v, err = g.loadRemote(key); err == nil {
			return
		}

		if v, err = g.loadLocal(key); err == nil {
			g.gcache.add(key, v.(lru.Value))
			return v, err
		}
		return
	})

	val = v.(lru.Value)

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
	time.Sleep(10 * time.Second)
	bytes, err := g.getter.Get(key)
	if err != nil {
		return
	}
	// get this value from local by getter
	v = lru.NewValue(bytes)
	return
}
