package main

import (
	"fmt"
	group "gcache/Group"
	"net/http"
)

func main() {
	// add three groups
	kvs1 := map[string]string{
		"Tom":       "cat",
		"Jerry":     "mouse",
		"Tom&Jerry": "friend",
	}
	var getter1 group.GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs1[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}
	group.NewGroup("Tom&Jerry", 1000, getter1)

	kvs2 := map[string]string{
		"apple":  "fruit",
		"banana": "fruit",
		"orange": "fruit",
	}
	var getter2 group.GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs2[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}
	group.NewGroup("fruit", 1000, getter2)

	kvs3 := map[string]string{
		"bag":  "thing",
		"ship": "vehicle",
		"car":  "vehicle",
	}
	var getter3 group.GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs3[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("[ERROR] key %s is not exist", key)
	}
	group.NewGroup("thing", 1000, getter3)

	addr := "localhost:9999"
	// add three groups to HTTPPool
	peers := group.NewHTTPPool(addr, "/gcache/")
	// serve
	http.ListenAndServe(addr, peers)

}
