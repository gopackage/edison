package gpio

import "math"

// FromScale returns a converted input from min, max to 0.0...1.0.
func FromScale(input, min, max int) float64 {
	fmin := float64(min)
	fmax := float64(max)
	finput := float64(input)
	return (finput - math.Min(fmin, fmax)) / (math.Max(fmin, fmax) - math.Min(fmin, fmax))
}
