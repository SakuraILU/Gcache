package gcache

import (
	"fmt"
	consistenthash "gcache/ConsistentHash"
	group "gcache/Group"
	"gcache/pb"
	"hash/crc32"
	"log"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

type HTTPPool struct {
	self      string
	base_path string

	// groups
	groups map[string]*group.Group
	glk    sync.RWMutex

	// peers
	peers           map[string]*group.HttpGetter
	consistent_hash *consistenthash.ConsistentHash
	plk             sync.RWMutex
}

func NewHTTPPool(self, base_path string) *HTTPPool {
	log.Printf("HTTPPool is at %v/%v\n", self, base_path)

	return &HTTPPool{
		self:      self,
		base_path: base_path,

		groups: make(map[string]*group.Group),
		glk:    sync.RWMutex{},

		peers:           make(map[string]*group.HttpGetter),
		consistent_hash: consistenthash.NewConsistentHash(hashfun, 3),
		plk:             sync.RWMutex{},
	}
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err_handle := func(err error) {
		log.Println(err.Error())
		data := &pb.Response{Value: []byte(err.Error())}
		body, err := proto.Marshal(data)
		w.Write(body)
	}

	path := r.URL.Path
	if path[:len(h.base_path)] != h.base_path {
		log.Fatal("unsupported base path")
	}
	relpath := path[len(h.base_path):]
	strs := strings.SplitN(relpath, "/", 2)
	gname := strs[0]
	key := strs[1]

	g, ok := h.groups[gname]
	if !ok {
		err_handle(fmt.Errorf("[ERROR] group %s is not exist", gname))
		return
	}

	v, err := g.Get(key)
	if err != nil {
		err_handle(err)
		return
	}

	// marshal data based on protobuf
	data := &pb.Response{Value: v.ByteSlice()}
	body, err := proto.Marshal(data)
	_, err = w.Write(body)
	if err != nil {
		err_handle(err)
		return
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		err_handle(err)
		return
	}
}

func (h *HTTPPool) AddGroup(name string, maxbytes int, getter group.Getter) {
	h.glk.Lock()
	defer h.glk.Unlock()

	if _, ok := h.groups[name]; ok {
		log.Fatalf("[ERROR] group %s is already exist", name)
	}

	h.groups[name] = group.NewGroup(name, maxbytes, getter, h)
}

func (h *HTTPPool) GetGroup(name string) (g *group.Group, err error) {
	h.glk.RLock()
	defer h.glk.RUnlock()

	g, ok := h.groups[name]
	if !ok {
		err = fmt.Errorf("[ERROR] group %s is not exist", name)
	}
	return
}

func (h *HTTPPool) AddRemotePeers(urls ...string) {
	h.plk.Lock()
	defer h.plk.Unlock()

	for _, url := range urls {
		h.peers[url] = group.NewHttpGetter(url)
		h.consistent_hash.Add(url)
	}
}

func (h *HTTPPool) PickPeer(key string) (httpgetter *group.HttpGetter, err error) {
	h.plk.RLock()
	defer h.plk.RUnlock()

	iurl, err := h.consistent_hash.Get(key)
	if err != nil {
		return
	}
	url, ok := iurl.(string)
	if !ok {
		err = fmt.Errorf("invalid url")
		return
	}
	// log.Printf("url %s, self %s", url, h.self)
	if url == "http://"+h.self+h.base_path {
		err = fmt.Errorf("self url...not remote")
		return
	}

	httpgetter = h.peers[url]
	return
}

func hashfun(i interface{}) int {
	str, ok := i.(string)
	if !ok {
		log.Fatal("invalid value to hash, not a string")
	}

	code := crc32.ChecksumIEEE([]byte(str))

	// log.Println("hash", str, code)

	return int(code)
}
