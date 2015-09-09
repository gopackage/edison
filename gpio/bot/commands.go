package bot

import (
	"fmt"
	"strconv"
)

// Static generates a static text gpio export Commander.
func Static(path, value string) CommandFunc {
	return func(pin int) (string, string) { return path, value }
}

// Export generates a gpio export Commander.
func Export() CommandFunc {
	return func(pin int) (string, string) {
		return "/sys/class/gpio/export", strconv.Itoa(pin)
	}
}

// Unexport generates a gpio export Commander.
func Unexport() CommandFunc {
	return func(pin int) (string, string) {
		return "/sys/class/gpio/unexport", strconv.Itoa(pin)
	}
}

// DirOut generates a gpio direction output Commander.
func DirOut() CommandFunc {
	return func(pin int) (string, string) {
		return fmt.Sprintf("/sys/class/gpio/%d/direction", pin), "out"
	}
}

// DirIn generates a gpio direction output Commander.
func DirIn() CommandFunc {
	return func(pin int) (string, string) {
		return fmt.Sprintf("/sys/class/gpio/%d/direction", pin), "in"
	}
}

// Low generates a gpio pin "low" output command.
func Low() CommandFunc {
	return func(pin int) (string, string) {
		return fmt.Sprintf("/sys/class/gpio/%d/value", pin), "0"
	}
}

// High generates a gpio pin "high" output command.
func High() CommandFunc {
	return func(pin int) (string, string) {
		return fmt.Sprintf("/sys/class/gpio/%d/value", pin), "1"
	}
}

// PinMode generates a PWM pin mux mode output command.
func PinMode(mode int) CommandFunc {
	return func(pin int) (string, string) {
		return fmt.Sprintf("/sys/kernel/debug/gpio_debug/gpio%d/current_pinmux", pin),
			fmt.Sprintf("mode%d", mode)
	}
}
