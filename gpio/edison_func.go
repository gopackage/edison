package gpio

import (
	"fmt"
	"os"
)

func intel_edison_pinmode_change(sysfs int, mode int) error {
	if mode < 0 {
		return nil
	}

	buffer := fmt.Sprintf(SYSFS_PINMODE_PATH+"%i/current_pinmux", sysfs)
	modef, err := os.OpenFile(buffer, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("edison: Failed to open SoC pinmode for opening")
	}

	_, err = modef.Write([]byte(fmt.Sprintf("mode%d", mode)))
	modef.Close()

	return err
}

func intel_edison_pwm_init_pre(pin int) error {
	if miniboard == true {
		return intel_edison_pinmode_change(int(plat.pins[pin].gpio.pinmap), 1)
	}

	// TODO(stephen) handle non miniboard case

	return nil
}

func intel_edison_gpio_init_post(dev *gpio_context) error {
	if dev == nil {
		return fmt.Errorf("edison: passed in gpio_context is null")
	}

	var sysfs, mode int
	if miniboard == true {
		sysfs = dev.pin
		mode = 0
	} else {
		// TODO(stephen) Implement non miniboard option
	}

	return intel_edison_pinmode_change(sysfs, mode)
}

func intel_edison_gpio_mode_replace(dev *gpio_context, mode int) error {
	if dev.value_fp != nil {
		if err := dev.value_fp.Close(); err != nil {
			return fmt.Errorf("edison: error closing value_fp for pin: %d", dev.pin)
		}
		dev.value_fp = nil
	}

	pullup, err := GpioInitRaw(pullup_map[dev.phy_pin])
	if err != nil {
		return fmt.Errorf("edison: error creating pullup pin: %d", dev.phy_pin)
	}

	if err := gpio_dir(pullup, IN); err != nil {
		gpio_close(pullup)
		return fmt.Errorf("edison: Failed to setup gpio mode-pullup: %s", err)
	}

	value := -1
	switch mode {
	case HIGH:
		value = 1
		break
	case LOW:
		value = 0
		break
	default:
		return fmt.Errorf("edison: invalid mode sent to mode_replace: %d", mode)
	}

	if value != -1 {
		if err := gpio_dir(pullup, OUT); err != nil {
			gpio_close(pullup)
			return fmt.Errorf("edison: Error setting pullup")
		}
		if err := gpio_write(pullup, value); err != nil {
			gpio_close(pullup)
			return fmt.Errorf("edison: Error setting pullup")
		}
	}

	return gpio_close(pullup)
}
