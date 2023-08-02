package main

import (
	"fmt"
	gcache "gcache/Gcache"
	group "gcache/Group"
	"log"
	"net/http"
)

func startServer(name string, addr string, kvs1 map[string]string) {
	htt_urls := []string{
		"http://localhost:9999/gcache/",
		"http://localhost:10000/gcache/",
		"http://localhost:10001/gcache/",
	}

	peers := gcache.NewHTTPPool(addr, "/gcache/")
	var getter1 group.GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs1[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key not found")
	}
	peers.AddGroup(name, 1000, getter1)
	peers.AddRemotePeers(htt_urls...)

	http.ListenAndServe(addr, peers)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var addrs []string = []string{
		"localhost:9999",
		"localhost:10000",
		"localhost:10001",
	}

	kvs := map[string]string{
		"Tom":       "cat",
		"Jerry":     "mouse",
		"Tom&Jerry": "friend",
		"bag":       "thing",
		"ship":      "vehicle",
		"car":       "vehicle",
		"apple":     "fruit",
		"banana":    "fruit",
		"orange":    "fruit",
		"127.0.0.1": "ip\nefgeda\n\neadaw",
	}

	go startServer("thing", addrs[0], kvs)
	go startServer("thing", addrs[1], kvs)
	go startServer("thing", addrs[2], kvs)

	select {}
}
