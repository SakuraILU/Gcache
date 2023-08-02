package singleflight

import (
	"log"
	"sync"
)

type TaskFn func() (interface{}, error)

type SinlgeFlight struct {
	calls map[string]*call
	clk   sync.RWMutex
}

func NewSingleFlight() (s *SinlgeFlight) {
	return &SinlgeFlight{
		calls: make(map[string]*call),
		clk:   sync.RWMutex{},
	}
}

func (s *SinlgeFlight) Do(key string, fn TaskFn) (val interface{}, err error) {
	{
		s.clk.RLock()

		if c, ok := s.calls[key]; ok {
			log.Println("Wait for the result of the same key")
			s.clk.RUnlock()
			val, err = c.wait()
			return
		}
		s.clk.RUnlock()
	}

	c := newCall()
	{
		s.clk.Lock()
		s.calls[key] = c
		s.clk.Unlock()
	}

	val, err = fn()
	c.val, c.err = val, err
	c.done()

	{
		s.clk.Lock()
		delete(s.calls, key)
		s.clk.Unlock()
	}

	return
}
