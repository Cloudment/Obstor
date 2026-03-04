package cmd

// MovingAverage is the interface for an exponentially weighted moving average.
type MovingAverage interface {
	Add(float64)
	Value() float64
	Set(float64)
}

// SimpleEWMA is a simple exponentially weighted moving average
// implementation using the default decay factor (2 / (N + 1)) where N = 30.
type SimpleEWMA struct {
	avg  float64
	init bool
}

const ewmaDecay = 2.0 / (30.0 + 1.0)

// Add adds a new value to the moving average.
func (e *SimpleEWMA) Add(value float64) {
	if !e.init {
		e.avg = value
		e.init = true
		return
	}
	e.avg = ewmaDecay*value + (1-ewmaDecay)*e.avg
}

// Value returns the current value of the moving average.
func (e *SimpleEWMA) Value() float64 {
	return e.avg
}

// Set sets the moving average value directly.
func (e *SimpleEWMA) Set(value float64) {
	e.avg = value
	e.init = true
}
