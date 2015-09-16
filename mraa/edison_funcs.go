package mraa

import (
	"fmt"
	"os"
)

var (
	outputen   = [...]int{248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 232, 233, 234, 235, 236, 237}
	pullup_map = [...]int{216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 208, 209, 210, 211, 212, 213}
)

func intel_edison_pinmode_change(sysfs int, mode int) error {
	if mode < 0 {
		return nil
	}

	buffer := fmt.Sprintf(SYSFS_PINMODE_PATH+"%d/current_pinmux", sysfs)
	modef, err := os.OpenFile(buffer, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("edison: Failed to open SoC pinmode for opening: %s", err)
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

func intel_edison_gpio_mode_replace(dev *gpio_context, mode string) error {
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
	case MODE_GPIO_STRONG:
		break
	case MODE_GPIO_PULLUP:
		value = 1
		break
	case MODE_GPIO_PULLDOWN:
		value = 0
		break
	case MODE_GPIO_HIZ:
		return nil
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

func intel_edison_i2c_init_pre(bus uint) error {
	if miniboard == false {
		// TODO(stephen): handle non-miniboard init pre
	} else {
		if bus != 6 && bus != 1 {
			fmt.Printf("edison: You can't use that bus, switching to bus 6\n")
			bus = 6
		}
		scl := plat.pins[plat.i2c_bus[bus].scl].gpio.pinmap
		sda := plat.pins[plat.i2c_bus[bus].sda].gpio.pinmap
		intel_edison_pinmode_change(int(sda), 1)
		intel_edison_pinmode_change(int(scl), 1)
	}

	return nil
}

func intel_edison_i2c_freq(dev *i2c_context, mode string) error {
	var sysnode *os.File = nil
	var err error

	switch dev.busnum {
	case 1:
		sysnode, err = os.OpenFile("/sys/devices/pci0000:00/0000:00:08.0/i2c_dw_sysnode/mode", os.O_RDWR, 0664)
		break
	case 6:
		sysnode, err = os.OpenFile("/sys/devices/pci0000:00/0000:00:09.0/i2c_dw_sysnode/mode", os.O_RDWR, 0664)
		break
	default:
		err = fmt.Errorf("edisoN: i2c bus selected does not support frequency changes")
	}

	if err != nil {
		return err
	}
	defer sysnode.Close()

	buf := ""
	switch mode {
	case I2C_STD:
		buf = "std"
		break
	case I2C_FAST:
		buf = "fast"
		break
	case I2C_HIGH:
		buf = "high"
		break
	default:
		return fmt.Errorf("edison: Invalid i2c mode selected")
	}

	if _, err = sysnode.Write([]byte(buf)); err != nil {
		return fmt.Errorf("edison: Error writing to sysnode %s: %s", sysnode, err)
	}
	return nil
}
