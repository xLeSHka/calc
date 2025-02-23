package cache

import (
	"github.com/xLeSHka/calc/internal/pkg/config"
	"sync"
	"time"
)

type Cache struct {
	addictionTime      time.Duration
	subtractionTime    time.Duration
	divisionTime       time.Duration
	multiplicationTime time.Duration
	exponentiationTime time.Duration
	unaryMinusTime     time.Duration
	logarithmTime      time.Duration
	squareRootTime     time.Duration
	mu                 *sync.RWMutex
}

func New(config config.Config) *Cache {
	return &Cache{
		addictionTime:      config.TimeAddiction,
		subtractionTime:    config.TimeSubtraction,
		divisionTime:       config.TimeDivision,
		multiplicationTime: config.TimeMultiplication,
		exponentiationTime: config.TimeExponentiation,
		unaryMinusTime:     config.TimeUnaryMinus,
		logarithmTime:      config.TimeLogarithm,
		squareRootTime:     config.TimeSquareRoot,
		mu:                 &sync.RWMutex{},
	}
}
func (t *Cache) AddictionTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.addictionTime
}
func (t *Cache) SubtractionTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.subtractionTime
}
func (t *Cache) DivisionTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.divisionTime
}
func (t *Cache) MultiplicationTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.multiplicationTime
}
func (t *Cache) ExponentiationTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.exponentiationTime
}
func (t *Cache) UnaryMinusTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.unaryMinusTime
}
func (t *Cache) LogarithmTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.logarithmTime
}
func (t *Cache) SquareRootTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.squareRootTime
}
