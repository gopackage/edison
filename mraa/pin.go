package mraa

// Limits defines the minimum and maximum value of a scaled duty cycle.
type Limits struct {
	Min int
	Max int
}

// Validate verifies that the limits are properly ordered
// with the Min always less than or equal to the Max.
func (l *Limits) Validate() {
	if l.Min > l.Max {
		// Swap
		max := l.Min
		l.Min = l.Max
		l.Max = max
	}
}
