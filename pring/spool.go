package pring

import "sync"

type SPool struct {
	c    chan []byte
	pool *sync.Pool
}

func NewSPool(size int) *SPool {
	new := func() interface{} {
		return make([]byte, 1<<16)
	}
	sp := &SPool{
		c:    make(chan []byte, size),
		pool: &sync.Pool{New: new}}
	return sp
}

func (sp *SPool) Get() []byte {
	return sp.pool.Get().([]byte)
}

func (sp *SPool) Put(b []byte) {
	sp.pool.Put(b[:cap(b)])
}

func (sp *SPool) Write(b []byte) {
	sp.c <- b
}

func (sp *SPool) Read() []byte {
	return <-sp.c
}
