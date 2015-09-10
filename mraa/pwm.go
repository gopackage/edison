package mraa

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	SYSFS_PWM = "/sys/class/pwm"
)

type pwm_context struct {
	pin     int      // the pin number, as known to the os
	chipid  int      // the chip id, which the pwm resides
	duty_fp *os.File // file pointer to the duty file
	period  int      // cache the period to speed up setting duty
	owner   bool     // Owner of the pwm context

	// Automatic scaling of the duty cycle
	limits Limits  // The limits of scaled values
	span   float64 // The distance between the scale limits
	min    float64 // The minimum value of the limits as a float
}

type PwmContext *pwm_context

// SetLimits sets the limits that scaled PWM outputs are based on.
// By default, the scale is 0 - 255 to support unsigned byte values
// common for PWM (e.g. for scaling 24-bit RGB colors).
func (dev *pwm_context) SetLimits(l Limits) {
	l.Validate()
	dev.limits = l
	dev.min = float64(l.Min)
	dev.span = float64(l.Max) - dev.min
}

func (dev *pwm_context) SetupDutyFile() error {
	buf := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/duty_cycle", dev.chipid, dev.pin)
	file, err := os.OpenFile(buf, os.O_RDWR, 0664)
	if err != nil {
		return fmt.Errorf("pwm: error opening duty cycle file: %s", buf)
	}
	dev.duty_fp = file
	return nil
}

func (dev *pwm_context) WritePeriod(period int) error {
	if advance_func.pwm_period_replace != nil {
		err := advance_func.pwm_period_replace(dev, period)
		if err == nil {
			dev.period = period
		}
		return err
	}

	buf := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/period", dev.chipid, dev.pin)
	period_f, err := os.OpenFile(buf, os.O_RDWR, 0664)
	if err != nil {
		return fmt.Errorf("pwm: Failed to open period for writing")
	}
	defer period_f.Close()

	out := fmt.Sprintf("%d", period)
	if _, err := period_f.Write([]byte(out)); err != nil {
		return fmt.Errorf("pwm: Failed to write period to file: %s, %s", out, err)
	}
	fmt.Printf("pwm: Wrote period[%d] to pin[%d]\n", period, dev.pin)

	dev.period = period
	return nil
}

func (dev *pwm_context) Period(us int) error {
	if us < plat.pwm_min_period || us > plat.pwm_max_period {
		return fmt.Errorf("pwm: period vlaue outside platform range")
	}
	return dev.WritePeriod(us * 1000)
}

func (dev *pwm_context) Enable(enable int) error {
	var status int = 0
	if enable != 0 {
		status = 1
	} else {
		status = 0
	}

	buf := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/enable", dev.chipid, dev.pin)

	enable_f, err := os.OpenFile(buf, os.O_RDWR, 664)
	if err != nil {
		return fmt.Errorf("pwm: Failed to open enable for writing for pin %d: %s", dev.pin, err)
	}
	defer enable_f.Close()

	out := fmt.Sprintf("%d", status)
	if _, err = enable_f.Write([]byte(out)); err != nil {
		return fmt.Errorf("pwm: Failed to write to enable for pin %d: %s", dev.pin, err)
	}

	return nil
}

func (dev *pwm_context) ReadPeriod() error {
	buf := fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/period", dev.chipid, dev.pin)

	output, err := ioutil.ReadFile(buf)
	if err != nil {
		return fmt.Errorf("pwm: Error reading period for pin[%d]: %s", dev.pin, err)
	}

	i, err := parseInt(string(output), 32)
	if err != nil {
		return fmt.Errorf("pwm: Can't convert period output to int: %s, %s", string(output), err)
	}

	dev.period = int(i)
	return nil
}

func (dev *pwm_context) ReadDuty() (int, error) {
	if dev.duty_fp == nil {
		if err := dev.SetupDutyFile(); err != nil {
			return -1, fmt.Errorf("pwm: Error setting up duty fp in pwm_read_duty for pin[%d]: %s", dev.pin, err)
		}
	} else {
		if _, err := dev.duty_fp.Seek(0, os.SEEK_SET); err != nil {
			return -1, fmt.Errorf("pwm: error seeking on duty file for pin[%d]:", dev.pin, err)
		}
	}

	// Should never be larger than this, in fact 4096 is probably too large
	var output []byte = make([]byte, 0, 4096)
	if _, err := dev.duty_fp.Read(output); err != nil {
		return -1, fmt.Errorf("pwm: Error in reading duty for pin[%d]: %s", dev.pin, err)
	}

	i, err := parseInt(string(output), 32)
	if err != nil {
		return -1, fmt.Errorf("pwm: Can't convert duty output to int: %s, %s", string(output), err)
	}
	fmt.Printf("pwm: Got duty[%d] for pin[%d]\n", i, dev.pin)

	return int(i), nil
}

func (dev *pwm_context) Read() (float32, error) {
	if err := dev.ReadPeriod(); err != nil {
		return -1.0, fmt.Errorf("pwm: Error reading period in pwm_read for pin[%d]: %s", dev.pin, err)
	}

	if dev.period > 0 {
		duty, err := dev.ReadDuty()
		fmt.Printf("pwm_read: period: %d, duty: %d\n", dev.period, duty)
		if err != nil {
			return -1.0, fmt.Errorf("pwm: Error reading duty in pwm_read for pin[%d]: %s", dev.pin, err)
		}
		return float32(duty) / float32(dev.period), nil
	}

	return 0.0, nil
}

func (dev *pwm_context) WriteDuty(duty int) error {
	if dev.duty_fp == nil {
		if err := dev.SetupDutyFile(); err != nil {
			return fmt.Errorf("pwm: Error writing duty: %s", err)
		}
	}

	buf := fmt.Sprintf("%d", duty)
	if _, err := dev.duty_fp.Write([]byte(buf)); err != nil {
		return fmt.Errorf("pwm: Error writing duty %d to pin[%d]: %s", duty, dev.pin, err)
	}

	return nil
}

func (dev *pwm_context) WritePerc(percentage float32) error {
	if dev.period == -1 {
		if err := dev.ReadPeriod(); err != nil {
			return fmt.Errorf("pwm: pwm_write: %s", err)
		}
	}

	fmt.Printf("pwm: In pwm_write, period is set to: %d\n", dev.period)
	fmt.Printf("pwm: In pwm_write, percentage is set to %f\n", percentage)

	var duty int = 0
	if percentage > 1.0 {
		fmt.Printf("pwm: number greater than 1 entered, defaulting to 100 percent\n")
		duty = dev.period
	} else {
		duty = int(percentage * float32(dev.period))
	}
	fmt.Printf("pwm: pin[%d] writing new duty[%d]\n", dev.pin, duty)
	return dev.WriteDuty(duty)
}

// Scale sets the duty cycle based on the currently set Limits for the PWM.
func (dev *pwm_context) Scale(value int) error {
	if dev.period == -1 {
		if err := dev.ReadPeriod(); err != nil {
			return fmt.Errorf("pwm: error running Scale: %s", err)
		}
	}

	duty := (float64(value) - dev.min) / dev.span
	fmt.Printf("pwm: Scaling pin[%d] from value: %d to duty: %f\n", dev.pin, value, duty)
	return dev.WriteDuty(int(float64(dev.period) * duty))
}

func PwmInitRaw(chipin, pin int) (*pwm_context, error) {
	dev := &pwm_context{}
	dev.duty_fp = nil
	dev.chipid = chipin
	dev.pin = pin
	dev.period = -1
	dev.SetLimits(Limits{0, 255})

	directory := fmt.Sprintf(SYSFS_PWM+"/pwmchip%d/pwm%d", dev.chipid, dev.pin)
	s, err := os.Stat(directory)
	if err == nil && s.IsDir() {
		fmt.Printf("pwm: Pin already exported, continuing\n")
		dev.owner = false
	} else {
		buffer := fmt.Sprintf("/sys/class/pwm/pwmchip%d/export", dev.chipid)
		export_f, err := os.OpenFile(buffer, os.O_WRONLY, 0664)
		if err != nil {
			return nil, fmt.Errorf("pwm: failed to open export for writing: %s", err)
		}
		defer export_f.Close()

		out := fmt.Sprintf("%d", dev.pin)
		if _, err := export_f.Write([]byte(out)); err != nil {
			return nil, fmt.Errorf("pwm: Failed to write to export!  Potentially already enabled: %s", err)
		}
		dev.owner = true
		dev.Period(plat.pwm_default_period)
	}
	if err := dev.SetupDutyFile(); err != nil {
		return nil, fmt.Errorf("pwm: Error setting up duty file: %s", err)
	}
	return dev, nil
}

func PwmInit(pin int) (*pwm_context, error) {
	if advance_func.pwm_init_replace != nil {
		return advance_func.pwm_init_replace(pin)
	}

	if advance_func.pwm_init_pre != nil {
		if err := advance_func.pwm_init_pre(pin); err != nil {
			return nil, fmt.Errorf("pwm: error running init pre: %s", err)
		}
	}

	if plat == nil {
		return nil, fmt.Errorf("pwm: Platform not initialized")
	}

	if plat.pins[pin].capabilities.pwm != true {
		return nil, fmt.Errorf("pwm: pin not capable of pwm")
	}

	if plat.pins[pin].capabilities.gpio == true {
		// This deserves more investigation
		// TODO(stephen) figure out what the last comment means
		mux_i, err := GpioInitRaw(int(plat.pins[pin].gpio.pinmap))
		if err != nil {
			return nil, fmt.Errorf("pwm: error in gpio->pwm transition (init raw): %s", err)
		}
		if err = gpio_dir(mux_i, OUT); err != nil {
			return nil, fmt.Errorf("pwm: error in gpio->pwm transition (dir): %s", err)
		}
		if err = gpio_write(mux_i, 1); err != nil {
			return nil, fmt.Errorf("pwm: error in gpio->pwm transition (write): %s", err)
		}
		if err = gpio_close(mux_i); err != nil {
			return nil, fmt.Errorf("pwm: error in gpio->pwm transition (close): %s", err)
		}
	}

	if plat.pins[pin].pwm.mux_total > 0 {
		if err := setup_mux_mapped(plat.pins[pin].pwm); err != nil {
			return nil, fmt.Errorf("pwm: Failed to setup multiplexer")
		}
	}

	chip := int(plat.pins[pin].pwm.parent_id)
	pinn := int(plat.pins[pin].pwm.pinmap)

	if advance_func.pwm_init_post != nil {
		pret, err := PwmInitRaw(chip, pinn)
		if err != nil {
			return nil, fmt.Errorf("pwm: error creating pwm pin: %s", err)
		}
		if err := advance_func.pwm_init_post(pret); err != nil {
			return nil, fmt.Errorf("pwm: error creating pret: %s", err)
		}
		return pret, nil
	}
	return PwmInitRaw(chip, pinn)
}
