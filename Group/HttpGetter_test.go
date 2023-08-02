package group

import (
	"bytes"
	"testing"
)

func TestGetter(t *testing.T) {
	// request: start a HTTP server at localhost:9999 (db: kvs)
	addr := "localhost:9999"
	kvs := map[string]string{
		"Tom":       "cat",
		"Jerry":     "mouse",
		"Tom&Jerry": "friend",
	}

	// test httpGetter
	// define a httpGetter
	httpgetter := NewHttpGetter("http://" + addr + "/gcache/")
	// get values from httpGetter
	for k, v := range kvs {
		if val, err := httpgetter.Get("Tom&Jerry", k); err != nil {
			t.Errorf("[Error]: %v", err)
		} else if bytes.Equal(val, []byte(v)) {
			t.Errorf("[Error]: %v != %v", val, v)
		}
	}
}
