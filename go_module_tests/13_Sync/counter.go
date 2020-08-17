package main

import "sync"

type Counter struct {
	count int
	mx    sync.Mutex
}

//NewCounter is a Counter constructor that returns a pointer
func NewCounter() *Counter {
	return &Counter{}
}

//Value returns the counter value
func (c *Counter) Value() int {
	return c.count
}

//Increment increments the counter.count with one
func (c *Counter) Increment() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.count++
}
