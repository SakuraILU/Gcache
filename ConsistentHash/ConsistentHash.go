package consistenthashmap

import (
	"fmt"
	"log"
	"sort"
)

type HashFun func(interface{}) int

type ConsistentHash struct {
	hashs    []int
	hashvals map[int]interface{}
	hashfun  HashFun
	replicas int
}

func NewConsistentHash(hashfun HashFun, replicas int) (c *ConsistentHash) {
	return &ConsistentHash{
		hashs:    make([]int, 0),
		hashvals: make(map[int]interface{}),
		hashfun:  hashfun,
		replicas: replicas,
	}
}

func (c *ConsistentHash) Add(vals ...string) {
	for _, val := range vals {
		for i := 0; i < c.replicas; i++ {
			hash := c.hashfun(fmt.Sprint(i) + val)
			c.hashs = append(c.hashs, hash)
			c.hashvals[hash] = val
			log.Printf("Add key %s into consisitentHash\n", val)
		}
	}

	sort.Ints(c.hashs)
}

func (c *ConsistentHash) Get(val interface{}) (elem interface{}, err error) {
	if val == nil {
		err = fmt.Errorf("[ERROR] empty value to get")
		return
	}

	hash := c.hashfun(val)
	idx := sort.Search(len(c.hashs), func(i int) bool {
		fmt.Printf("hash: %d, c.hashs[i]: %d\n", hash, c.hashs[i])
		return c.hashs[i] >= hash
	})

	// val may be very large...exceed all and == len(c.hashs)
	// %len(c.hashs) to moveTo the front
	// idx --> hash --> val
	elem = c.hashvals[c.hashs[idx%len(c.hashs)]]

	return
}
