package counter

import "sync"

type Counter struct {
	counter int64
	mu      *sync.Mutex
}

func New() *Counter {
	return &Counter{
		mu:      &sync.Mutex{},
		counter: 1,
	}
}

func (c *Counter) Int() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	cur := c.counter
	c.counter++
	return cur
}
func (c *Counter) Restart() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter = 1
}
