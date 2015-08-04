package edison

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

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

func writeFile(path string, data []byte) (i int, err error) {
	fmt.Println(">>", path, string(data))
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i = 0
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}

type mux struct {
	pin   int
	value int
}

type pin struct {
	pin   int
	rPin  int
	label string
}

type fsPin struct {
	pin          int
	resistor     int
	levelShifter int
	pwmPin       int
	label        string
	mux          []mux
}

var pwmPins = make(map[int]*pwmPin)
var digitalPins = make(map[int]DigitalPin)

// Pin represents a GPIO pin on the Edison.
type Pin interface {
	Direction(dir string) (err error)
	Export() (err error)
	Unexport() (err error)
}

// DigitalPin abstracts a digital (vs analog) pin.
type DigitalPin interface {
	Pin
	DigitalWrite(int) (err error)
	DigitalRead() (val int, err error)
	PwmWrite(val byte) (err error)
}

var pinMap = map[int]fsPin{
	0: fsPin{
		pin:          130,
		resistor:     216,
		levelShifter: 248,
		pwmPin:       -1,
		mux:          []mux{},
	},
	1: fsPin{
		pin:          131,
		resistor:     217,
		levelShifter: 249,
		pwmPin:       -1,
		mux:          []mux{},
	},
	2: fsPin{
		pin:          128,
		resistor:     218,
		levelShifter: 250,
		pwmPin:       -1,
		mux:          []mux{},
	},
	3: fsPin{
		pin:          12,
		resistor:     219,
		levelShifter: 251,
		pwmPin:       0,
		mux:          []mux{},
	},

	4: fsPin{
		pin:          129,
		resistor:     220,
		levelShifter: 252,
		pwmPin:       -1,
		mux:          []mux{},
	},
	5: fsPin{
		pin:          13,
		resistor:     221,
		levelShifter: 253,
		pwmPin:       1,
		mux:          []mux{},
	},
	6: fsPin{
		pin:          182,
		resistor:     222,
		levelShifter: 254,
		pwmPin:       2,
		mux:          []mux{},
	},
	7: fsPin{
		pin:          48,
		resistor:     223,
		levelShifter: 255,
		pwmPin:       -1,
		mux:          []mux{},
	},
	8: fsPin{
		pin:          49,
		resistor:     224,
		levelShifter: 256,
		pwmPin:       -1,
		mux:          []mux{},
	},
	9: fsPin{
		pin:          183,
		resistor:     225,
		levelShifter: 257,
		pwmPin:       3,
		mux:          []mux{},
	},
	10: fsPin{
		pin:          41,
		resistor:     226,
		levelShifter: 258,
		pwmPin:       4,
		mux: []mux{
			mux{263, HIGH},
			mux{240, LOW},
		},
	},
	11: fsPin{
		pin:          43,
		resistor:     227,
		levelShifter: 259,
		pwmPin:       5,
		mux: []mux{
			mux{262, HIGH},
			mux{241, LOW},
		},
	},
	12: fsPin{
		pin:          42,
		resistor:     228,
		levelShifter: 260,
		pwmPin:       -1,
		mux: []mux{
			mux{242, LOW},
		},
	},
	13: fsPin{
		pin:          40,
		resistor:     229,
		levelShifter: 261,
		pwmPin:       -1,
		mux: []mux{
			mux{243, LOW},
		},
	},
}

func newDigitalRootPin(p int, r int, v ...string) DigitalPin {
	d := &pin{pin: p, rPin: r}
	if len(v) > 0 {
		d.label = v[0]
	} else {
		d.label = "gpio" + strconv.Itoa(d.pin)
	}
	return d
}

func newDigitalPin(p int, v ...string) DigitalPin {
	d := &pin{pin: p}
	if len(v) > 0 {
		d.label = v[0]
	} else {
		d.label = "gpio" + strconv.Itoa(d.pin)
	}
	return d
}

// GetDigitalPin creates a pin object given pin number and output.
// TODO specify what pin number being used here
// TODO move to enum types for direction (rather than using a string)
func GetDigitalPin(p int, dir string) (pin DigitalPin, err error) {
	fpin := pinMap[p]
	if digitalPins[fpin.pin] == nil {
		digitalPins[fpin.pin] = newDigitalRootPin(fpin.pin, p)
		if err = digitalPins[fpin.pin].Export(); err != nil {
			return
		}
		digitalPins[fpin.resistor] = newDigitalPin(fpin.resistor)
		if err = digitalPins[fpin.resistor].Export(); err != nil {
			return
		}
		digitalPins[fpin.levelShifter] = newDigitalPin(fpin.levelShifter)
		if err = digitalPins[fpin.levelShifter].Export(); err != nil {
			return
		}

		if len(fpin.mux) > 0 {
			for _, mux := range fpin.mux {
				digitalPins[mux.pin] = newDigitalPin(mux.pin)
				if err = digitalPins[mux.pin].Export(); err != nil {
					return
				}
				if err = digitalPins[mux.pin].Direction(OUT); err != nil {
					return
				}
				if err = digitalPins[mux.pin].DigitalWrite(mux.value); err != nil {
					return
				}
			}
		}
	}

	if dir == IN {
		if err = digitalPins[fpin.pin].Direction(IN); err != nil {
			return
		}
		if err = digitalPins[fpin.resistor].Direction(OUT); err != nil {
			return
		}
		if err = digitalPins[fpin.resistor].DigitalWrite(LOW); err != nil {
			return
		}
		if err = digitalPins[fpin.levelShifter].Direction(OUT); err != nil {
			return
		}
		if err = digitalPins[fpin.levelShifter].DigitalWrite(LOW); err != nil {
			return
		}
	} else if dir == OUT {
		if err = digitalPins[fpin.pin].Direction(OUT); err != nil {
			return
		}
		if err = digitalPins[fpin.resistor].Direction(IN); err != nil {
			return
		}
		if err = digitalPins[fpin.levelShifter].Direction(OUT); err != nil {
			return
		}
		if err = digitalPins[fpin.levelShifter].DigitalWrite(HIGH); err != nil {
			return
		}
	}

	pin = digitalPins[fpin.pin]
	return pin, err
}

func (p *pin) DigitalWrite(b int) (err error) {
	str := fmt.Sprintf("%v/%s/value", GPIOPATH, p.label)
	fmt.Printf("Writing value (%d) to pin (%d)\n", b, p.pin)
	_, err = writeFile(str, []byte(strconv.Itoa(b)))
	return err
}

func (p *pin) Export() (err error) {
	fmt.Printf("Exporting pin: %d\n", p.pin)
	if _, err := writeFile(GPIOPATH+"/export", []byte(strconv.Itoa(p.pin))); err != nil {
		// If EBUSY then the pin has already been exported
		if err.(*os.PathError).Err != syscall.EBUSY {
			return err
		}
	}
	return nil
}

func (p *pin) Unexport() (err error) {
	if _, err := writeFile(GPIOPATH+"/unexport", []byte(strconv.Itoa(p.pin))); err != nil {
		// If EINVAL then the pin is reserved in the system and can't be unexported
		if err.(*os.PathError).Err != syscall.EINVAL {
			return err
		}
	}
	return nil
}

func (p *pin) Direction(dir string) (err error) {
	str := fmt.Sprintf("%v/%s/direction", GPIOPATH, p.label)
	fmt.Printf("Setting direction (%s) for pin (%d)\n", dir, p.pin)
	_, err = writeFile(str, []byte(dir))
	return err
}

func writeValue(pin int, val byte) (err error) {
	p, err := GetDigitalPin(pin, "out")
	if err != nil {
		return
	}
	return p.DigitalWrite(int(val))
}

// changePinMode writes pin mode to current_pinmux file
func changePinMode(pin int, mode string) (err error) {
	_, err = writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+strconv.Itoa(pin)+"/current_pinmux",
		[]byte("mode"+mode),
	)
	return
}

func (p *pin) PwmWrite(val byte) (err error) {
	fpin := pinMap[p.rPin]
	if fpin.pwmPin != -1 {
		if pwmPins[fpin.pwmPin] == nil {
			if err = writeValue(p.rPin, 1); err != nil {
				return
			}
			if err = changePinMode(fpin.pin, "1"); err != nil {
				return
			}
			pwmPins[fpin.pwmPin] = newPwmPin(fpin.pwmPin)
			if err = pwmPins[fpin.pwmPin].export(); err != nil {
				return
			}
			if err = pwmPins[fpin.pwmPin].enable("1"); err != nil {
				return
			}
		}

		periodStr, err := pwmPins[fpin.pwmPin].period()
		if err != nil {
			return err
		}
		period, err := strconv.Atoi(periodStr)
		if err != nil {
			return err
		}
		duty := FromScale(float64(val), 0, 255.0)
		return pwmPins[fpin.pwmPin].writeDuty(strconv.Itoa(int(float64(period) * duty)))
	}
	return errors.New("Not a PWM pin")
}

func (p *pin) DigitalRead() (val int, err error) {
	buf, err := readFile(fmt.Sprintf("%v/%d/value", GPIOPATH, p.pin))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(buf[0]))
}
