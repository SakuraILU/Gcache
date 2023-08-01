package consistenthashmap

import (
	"strconv"
	"testing"
)

func TestConsistentHash1(t *testing.T) {
	hashfun := func(val interface{}) (hash int) {
		v, ok := val.(string)
		if !ok {
			t.Fatalf("[ERROR] invalid value type %T", val)
		}
		hash, err := strconv.Atoi(v)
		if err != nil {
			t.Fatalf("[ERROR] invalid value %v", val)
		}
		return
	}

	c := NewConsistentHash(hashfun, 3)
	// add 10, 6, 1--> 09/19/29, 06/16/26, 02/12/22
	// 02--06--09--12--16--19--22--26--29
	c.Add("9", "6", "2")

	kvs := map[string]string{
		"5":  "6",
		"19": "9",
		"31": "2",
		"24": "6",
		"1":  "2",
		"21": "2",
	}
	for k, v := range kvs {
		elem, err := c.Get(k)
		if err != nil {
			t.Fatalf("[ERROR] %v", err)
		}

		if elem != v {
			t.Fatalf("[ERROR] %v != %v", elem, v)
		}
	}
}
