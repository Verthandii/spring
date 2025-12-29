package utils

import (
	"sync"
	"time"
)

type TimeMultiplier struct {
	t          time.Duration
	mutex      sync.RWMutex
	min, max   time.Duration
	multiplier float64
}

func NewTimeMultiplier(min, max time.Duration, multiplier float64) *TimeMultiplier {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if multiplier < 1 {
		multiplier = 1
	}
	return &TimeMultiplier{
		t:          min,
		min:        min,
		max:        max,
		multiplier: multiplier,
	}
}

func (it *TimeMultiplier) Incr() {
	it.mutex.Lock()
	it.t = time.Duration(float64(it.t) * it.multiplier)
	if it.isValidMax() && it.t > it.max {
		it.t = it.max
	}
	it.mutex.Unlock()
}

func (it *TimeMultiplier) Reset() {
	it.mutex.RLock()
	if it.t == it.min {
		it.mutex.RUnlock()
		return
	}
	it.mutex.RUnlock()

	it.mutex.Lock()
	it.t = it.min
	it.mutex.Unlock()
}

func (it *TimeMultiplier) Sleep() {
	if duration := it.Duration(); duration > 0 {
		time.Sleep(duration)
	}
}

func (it *TimeMultiplier) Duration() time.Duration {
	it.mutex.RLock()
	defer it.mutex.RUnlock()
	return it.t
}

func (it *TimeMultiplier) isValidMax() bool {
	return it.max > 0
}
