package gpio

import (
	"strconv"

	"github.com/gopackage/edison/gpio/bot"
	"github.com/gopackage/sysfs"
)

// Pins manages the hardware information about the Edison pin outs.
// By default, there are two main boards supported by the Edison,
// the Arduino board and the mini-breakout board. Most boards will
// follow the mini-breakout board pinouts (the default Edison pin layout).
// The Arduino board includes some extra i2c hardware that moves some
// of the standard gpio to different pin numbers and adds other pin
// types (like A/D). We can detect the Arduino board by looking for the
// gpio expansion devices on the i2c bus and assuming the arduino pinout.
// We allow users to override the auto-detection in case they want to force
// a particular pinout if auto-detection is not accurate.
//
// TODO autodetect arduino board
// TODO add mini-breakout pin outs
type Pins struct {
	Standard bool             // True if the pinout is a standard (mini-breakout)
	File     sysfs.DeviceFile // The device file implementation to use for generating pins
}

// NewPins creates a new pins factory object.
func NewPins(file sysfs.DeviceFile) *Pins {
	return &Pins{File: file}
}

// Autodetect the board type and set up the corresponding pin maps.
func (p *Pins) Autodetect() error {
	// TODO auto-detect - right now everything is a standard pinout
	p.Standard = true

	return nil
}

// Init sets up the pins according to whether we are running on a standard
// (mini breakout) or non-standard (arduino) board. You shuld run Autodetect()
// or set the board pin out type manually before calling Init().
func (p *Pins) Init() error {

	/*
		tristate, err := GpioInitRaw(214)
		if err != nil {
			return fmt.Errorf("Error opening tristate 214")
			miniboard = true
		}
	*/

	return nil
}

// DigitalOut creates a new digital output pin.
func (p *Pins) DigitalOut(pin int) (DigitalOutputPin, error) {
	info := pinMap[pin]
	b := bot.NewBot(nil)

	b.Add(info.pin, bot.Export(), bot.DirOut())
	b.Add(info.resistor, bot.Export(), bot.DirIn())
	b.Add(info.levelShifter, bot.Export(), bot.DirOut(), bot.High())

	err := p.setupMux(b, &info)
	if err != nil {
		return nil, err
	}
	err = b.Run()
	if err != nil {
		return nil, err
	}

	return &BasicOutputPin{pin: info.pin, rPin: pin, label: "gpio" + strconv.Itoa(pin), file: p.File}, nil
}

// DigitalIn creates a new digital input pin.
func (p *Pins) DigitalIn(pin int) (DigitalInputPin, error) {
	info := pinMap[pin]
	b := bot.NewBot(nil)

	b.Add(info.pin, bot.Export(), bot.DirIn())
	b.Add(info.resistor, bot.Export(), bot.DirOut(), bot.Low())
	b.Add(info.levelShifter, bot.Export(), bot.DirOut(), bot.Low())

	err := p.setupMux(b, &info)
	if err != nil {
		return nil, err
	}
	err = b.Run()
	if err != nil {
		return nil, err
	}

	return &BasicInputPin{pin: info.pin, rPin: pin, label: "gpio" + strconv.Itoa(pin), file: p.File}, nil
}

// PWM creates a new PWM output pin.
func (p *Pins) PWM(pin int) (PWMPin, error) {
	/*
		out, err := p.DigitalOut(pin)
		if err != nil {
			return nil, err
		}
		out.SetHigh()
		fpin := pinMap[out.rPin]
		if fpin.pwm != -1 {
			if err = p.file.Write("/sys/kernel/debug/gpio_debug/gpio"+strconv.Itoa(p.pin)+"/current_pinmux", "mode"+mode); err != nil {
				return
			}
			pwms[fpin.pwm] = NewBasicPWMPin(fpin.pwm, p.Pile)
			if err = pwms[fpin.pwm].Export(); err != nil {
				return
			}
			if err = pwms[fpin.pwm].Enable("1"); err != nil {
				return
			}
		}
	*/
	pwm, err := NewBasicPWMPin(pin, p.File)
	return pwm, err
	//return nil, errors.New("Not a PWM pin")
}

// setupMux configures multiplexer pins that are required by a pin being
// constructed by adding the appropriate actions to the provided pin.Bot.
func (p *Pins) setupMux(b *bot.Bot, info *PinInfo) error {
	if len(info.mux) > 0 {
		for _, mux := range info.mux {
			if mux.value == HIGH {
				b.Add(mux.pin, bot.Export(), bot.DirOut(), bot.High())
			} else {
				b.Add(mux.pin, bot.Export(), bot.DirOut(), bot.Low())
			}
		}
	}
	return nil
}

// PinInfo captures internal information about the pins on the Edison.
// Pin implementations use the information to properly configure and
// command the pin hardware through the sysfs driver.
type PinInfo struct {
	pin          int
	resistor     int
	levelShifter int
	pwm          int    // PWM pin or -1 for pins without PWM
	label        string // gpio file name
	mux          []mux  // List of mux pins that must be set to swizzle the pin
}

// mux tracks the pins and mux values that must be set on pin setup for each pin.
type mux struct {
	pin   int
	value int
}

/*

// Init initializes the Edison gpio pins. This global initialization will
// be removed and individual pins will initialize as they are created.
//
// TODO[scoward] remove this global init and move to individual pin inits (see #5)
func Init() error {
	var err error
	tristate := newDigitalPin(214)
	if err = tristate.Export(); err != nil {
		return err
	}
	if err = tristate.Direction(OUT); err != nil {
		return err
	}
	if err = tristate.DigitalWrite(LOW); err != nil {
		return err
	}

	for _, i := range []int{263, 262} {
		io := newDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(OUT); err != nil {
			return err
		}
		if err = io.DigitalWrite(HIGH); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		io := newDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(OUT); err != nil {
			return err
		}
		if err = io.DigitalWrite(LOW); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}

	}

	for _, i := range []int{111, 115, 114, 109} {
		if err = changePinMode(i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []int{131, 129, 40} {
		if err = changePinMode(i, "0"); err != nil {
			return err
		}
	}

	err = tristate.DigitalWrite(HIGH)
	return nil
}
*/
