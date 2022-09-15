package kademlia

import "math"

type Counter struct {
	i    int64
	lock chan struct{}
}

// Make a new thread safe counter
func MakeCounter() *Counter {
	c := &Counter{
		i:    0,
		lock: make(chan struct{}, 1),
	}
	c.lock <- struct{}{}
	return c
}

// Get next integer from counter. Will always be greater than 0
func (c *Counter) GetNext() int64 {
	<-c.lock

	if c.i == math.MaxInt64 {
		c.i = 0
	}
	c.i++
	x := c.i

	c.lock <- struct{}{}
	return x
}
