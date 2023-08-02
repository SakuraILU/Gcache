package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

func newCall() (c *call) {
	c = &call{
		wg: sync.WaitGroup{},
	}
	c.wg.Add(1)
	return
}

func (c *call) wait() (v interface{}, err error) {
	c.wg.Wait()
	return c.val, c.err
}

func (c *call) done() {
	c.wg.Done()
}
