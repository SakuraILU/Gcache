package group

import (
	"log"
	"net/http"
	"strings"
)

type HTTPPool struct {
	self      string
	base_path string
}

func NewHTTPPool(self, base_path string) *HTTPPool {
	log.Printf("HTTPPool is at %v/%v\n", self, base_path)

	return &HTTPPool{
		self:      self,
		base_path: base_path,
	}
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err_handle := func(err error) {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
	}

	path := r.URL.Path
	if path[:len(h.base_path)] != h.base_path {
		log.Fatal("unsupported base path")
	}
	relpath := path[len(h.base_path):]
	strs := strings.SplitN(relpath, "/", 2)
	gname := strs[0]
	key := strs[1]

	g, err := GetGroup(gname)
	if err != nil {
		err_handle(err)
		return
	}

	v, err := g.Get(key)
	if err != nil {
		err_handle(err)
		return
	}

	_, err = w.Write(v.ByteSlice())
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
