package counter

import "sync"

type Counter struct {
	expressionCounter int64
	taskCounter       int64
	mu                *sync.Mutex
}

func New() *Counter {
	return &Counter{
		mu:                &sync.Mutex{},
		expressionCounter: 0,
		taskCounter:       0,
	}
}
func (c *Counter) ExprInc() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	cur := c.expressionCounter
	c.expressionCounter++
	return cur
}
func (c *Counter) TaskInc() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	cur := c.taskCounter
	c.taskCounter++
	return cur
}
