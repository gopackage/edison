package gpio

const (
	// HIGH pin state
	HIGH = 1
	// LOW pin state
	LOW = 0
	// IN input pin mode
	IN = "in"
	// OUT output pin mode
	OUT = "out"
	// GPIOPATH is the path to the sysfs gpio
	GPIOPATH = "/sys/class/gpio"
)

// Pin represents a GPIO pin on the Edison.
type Pin interface {
	Export() (err error)
	Unexport() (err error)
}

// DigitalOutputPin is a digital gpio pin set to output mode.
type DigitalOutputPin interface {
	Pin
	// Sets the pin state to high.
	SetHigh() (err error)
	// Sets the pin state to low.
	SetLow() (err error)
}

// DigitalInputPin represents a digital (vs analog) pin in input mode.
type DigitalInputPin interface {
	Pin
	// High returns true if the pin is high.
	High() (isHigh bool, err error)
	// Low returns true if the pin is low.
	Low() (isLow bool, err error)
}

// PWMPin represents an pulse-width modulation (analog) output pin.
type PWMPin interface {
	Pin
	// Enable begins the PWM sending pulses
	Enable() error
	// Disable stops the PWM from sending pulses
	Disable() error
	// Period reads and returns the current PWM pin period value (pulses per second).
	Period() (period int, err error)
	// SetDuty changes the PWM duty cycle of the pin.
	SetDuty(duty int) error
}

// Limits defines the minimum and maximum value of a scaled duty cycle.
type Limits struct {
	Min int
	Max int
}

// Validate verifies that the limits are properly ordered
// with the Min in always less than or equal to the Max.
func (l *Limits) Validate() {
	if l.Min > l.Max {
		// Swap
		max := l.Min
		l.Min = l.Max
		l.Max = max
	}
}
