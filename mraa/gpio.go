package mraa

import (
	"fmt"
	"os"
)

const (
	// HIGH pin state
	HIGH = 1
	// LOW pin state
	LOW = 0

	// IN input pin dir
	IN = "in"
	// OUT output pin dir
	OUT = "out"
	// OUT Init high
	OUT_INIT_HIGH = "high"
	// OUT Init low
	OUT_INIT_LOW = "low"

	// modes
	MODE_GPIO_STRONG   = "strong"
	MODE_GPIO_PULLUP   = "pullup"
	MODE_GPIO_PULLDOWN = "pulldown"
	MODE_GPIO_HIZ      = "hiz"

	SYSFS_CLASS_GPIO   = "/sys/class/gpio"
	SYSFS_PINMODE_PATH = "/sys/kernel/debug/gpio_debug/gpio"
)

var (
	miniboard    bool     = false
	plat         *board_t = nil
	advance_func adv_func_t
)

type gpio_mux struct {
	pin   uint
	value uint
}

type gpio_context struct {
	pin      int      // the pin number as know to the OS
	phy_pin  int      // pin passed to clean init.  -1 none and raw
	value_fp *os.File // value file pointer, literally
	owner    bool     // if this context originally exported the pin
}

type GpioContext *gpio_context

type pin_t struct {
	pinmap        uint        // sysfs pin
	parent_id     uint        // parent chip id
	mux_total     uint        // number of muxes needed
	mux           [6]gpio_mux // mux array
	output_enable bool        // output enable gpio, for level shifting
	pullup_enable bool        // pull-up enable GPIO, inputs
	// TODO(stephen) complex cap
}

type pin_capabilities_t struct {
	valid     bool
	gpio      bool
	pwm       bool
	fast_gpio bool
	spi       bool
	i2c       bool
	aio       bool
	uart      bool
}

// Structure representing a physical pin
type pininfo_t struct {
	name         string // Pins real world name
	capabilities pin_capabilities_t
	gpio         pin_t
	pwm          pin_t
	aio          pin_t
	// TODO(stephen) mmap
	i2c  pin_t
	spi  pin_t
	uart pin_t
}

type i2c_bus_t struct {
	bus_id int  // ID as exposed in the system
	scl    uint // i2c SCL
	sda    uint // i2c SDA
}

type spi_bus_t struct {
	bus_id     uint // The Bus ID as exposed in the system
	slave_s    uint // slave select
	three_wire bool // Is the bus only a three wire system
	sclk       uint // serial clock
	mosi       uint // master out, slave in
	miso       uint // master in, slave out
	cs         uint // chip select, when used the board is a spi slave
}

type uart_dev_t struct {
	index       uint   // ID as exposted in the system
	rx          int    // uart rx
	tx          int    // uart tx
	device_path string // to store '/dev/ttySI' for example
}

type adv_func_t struct {
	gpio_init_pre  func(pin int) error
	gpio_init_post func(dev *gpio_context) error

	gpio_close_pre func(dev *gpio_context) error

	gpio_mode_replace func(dev *gpio_context, mode string) error
	gpio_mode_pre     func(dev *gpio_context, mode string) error
	gpio_mode_post    func(dev *gpio_context, mode string) error

	gpio_dir_replace func(dev *gpio_context, dir string) error
	gpio_dir_pre     func(dev *gpio_context, dir string) error
	gpio_dir_post    func(dev *gpio_context, dir string) error

	gpio_write_pre  func(dev *gpio_context, value int) error
	gpio_write_post func(dev *gpio_context, value int) error
	//type gpio_mmap_setup (mraa_gpio_context dev, mraa_boolean_t en) (error)

	//type i2c_init_pre (unsigned int bus) (error)
	//type i2c_init_post (mraa_i2c_context dev) (error)
	//type i2c_set_frequency_replace (mraa_i2c_context dev, mraa_i2c_mode_t mode) (error)

	//type aio_get_valid_fp (mraa_aio_context dev) (error)
	//type aio_init_pre (unsigned int aio) (error)
	//type aio_init_post (mraa_aio_context dev) (error)

	pwm_init_replace   func(pin int) (*pwm_context, error)
	pwm_init_pre       func(pin int) error
	pwm_init_post      func(pwm *pwm_context) error
	pwm_period_replace func(pwm *pwm_context, period int) error

	//type spi_init_pre (int bus) (error)
	//type spi_init_post (mraa_spi_context spi) (error)
	//type spi_lsbmode_replace (mraa_spi_context dev, mraa_boolean_t lsb) (error)

	//type uart_init_pre (int index) (error)
	//type uart_init_post (mraa_uart_context uart) (error)
}

type board_t struct {
	phy_pin_count      uint // total IO pins
	gpio_count         uint // gpio count
	aio_count          uint // analog side count
	i2c_bus_count      uint // < usable i2c count
	i2c_bus            [12]i2c_bus_t
	def_i2c_bus        uint // position in array of default i2c bus
	spi_bus_count      uint // usable spi count
	spi_bus            [12]spi_bus_t
	def_spi_bus        uint          // position in array of default spi bus
	adc_raw            uint          // ADC raw bit value
	adc_supported      uint          // ADC supported bit value
	def_uart_dev       uint          // Position in array of default uart
	uart_dev_count     uint          // Usable uart count
	uart_dev           [6]uart_dev_t // Array of UARTs
	pwm_default_period int           // the default PWM period in US
	pwm_max_period     int           // maximum period in us
	pwm_min_period     int           // minimum period in us
	platform_name      string        // platform name string
	pins               []pininfo_t   // pin array
}

func gpio_get_value_file_handle(dev *gpio_context) error {
	du := fmt.Sprintf(SYSFS_CLASS_GPIO+"/gpio%d/value", dev.pin)
	file, err := os.OpenFile(du, os.O_RDWR, 0664)
	if err != nil {
		return fmt.Errorf("Error opening value file %s for pin: %d", du, dev.pin)
	}
	dev.value_fp = file

	return nil
}

func gpio_unexport_force(dev *gpio_context) error {
	unexport, err := os.OpenFile(SYSFS_CLASS_GPIO+"/unexport", os.O_WRONLY, 0664)
	if err == nil {
		return fmt.Errorf("gpio: Failed to open unexport for writing")
	}

	bu := fmt.Sprintf("%d", dev.pin)
	if _, err = unexport.Write([]byte(bu)); err != nil {
		unexport.Close()
		return fmt.Errorf("gpio: Failed to write to unexport")
	}

	unexport.Close()
	// TODO(stephen) mraa_gpio_isr_exit(dev)
	return nil
}

func gpio_unexport(dev *gpio_context) error {
	if dev.owner {
		return gpio_unexport_force(dev)
	}
	return fmt.Errorf("gpio: not owner of pin context")
}

func gpio_close(dev *gpio_context) error {
	var err error

	if advance_func.gpio_close_pre != nil {
		err = advance_func.gpio_close_pre(dev)
	}

	if dev.value_fp != nil {
		dev.value_fp.Close()
	}
	gpio_unexport(dev)
	return err
}

func gpio_write(dev *gpio_context, value int) error {
	var err error

	if dev == nil {
		return fmt.Errorf("gpio: passed in gpio_context is nil")
	}

	// TODO(stephen) handle mmap writes

	if advance_func.gpio_write_pre != nil {
		if err = advance_func.gpio_write_pre(dev, value); err != nil {
			return fmt.Errorf("gpio: error in write pre: %s", err)
		}
	}

	if dev.value_fp == nil {
		if err = gpio_get_value_file_handle(dev); err != nil {
			return fmt.Errorf("gpio: error getting value fp: %s", err)
		}
	}

	if _, err = dev.value_fp.Seek(0, os.SEEK_SET); err != nil {
		return fmt.Errorf("gpio: error seeking on value file for pin: %d", dev.pin)
	}

	buf := fmt.Sprintf("%d", value)
	if _, err = dev.value_fp.Write([]byte(buf)); err != nil {
		fmt.Printf("gpio: error writing to value fp (%d): %s", dev.pin, err)
	}

	if advance_func.gpio_write_post != nil {
		return advance_func.gpio_write_post(dev, value)
	}

	return nil
}

var GpioWrite func(*gpio_context, int) error = gpio_write

func gpio_dir(dev *gpio_context, dir string) error {
	if advance_func.gpio_dir_replace != nil {
		return advance_func.gpio_dir_replace(dev, dir)
	}
	if advance_func.gpio_dir_pre != nil {
		if err := advance_func.gpio_dir_pre(dev, dir); err != nil {
			return err
		}
	}

	filepath := fmt.Sprintf(SYSFS_CLASS_GPIO+"/gpio%d/direction", dev.pin)

	direction, err := os.OpenFile(filepath, os.O_RDWR, 664)
	defer direction.Close()
	if err != nil {
		// Direction failed to Open.  If HIGH or LOW was passed will try and set
		// If not, fail as usual
		switch dir {
		case OUT_INIT_HIGH:
			return gpio_write(dev, 1)
		case OUT_INIT_LOW:
			return gpio_write(dev, 0)
		default:
			return fmt.Errorf("gpio: error opening gpio file: %s", err)
		}
	}

	if dir != IN && dir != OUT && dir != OUT_INIT_LOW && dir != OUT_INIT_LOW {
		fmt.Errorf("gpio: Direction passed in is incorrect: %s", dir)
	}

	if _, err := direction.Write([]byte(dir)); err != nil {
		return fmt.Errorf("gpio: error writing dir to file: %s", err)
	}

	if advance_func.gpio_dir_post != nil {
		return advance_func.gpio_dir_post(dev, dir)
	}

	return nil
}

var GpioDir func(*gpio_context, string) error = gpio_dir

func gpio_mode(dev *gpio_context, mode string) error {
	if advance_func.gpio_mode_replace != nil {
		return advance_func.gpio_mode_replace(dev, mode)
	}

	if advance_func.gpio_mode_pre != nil {
		if err := advance_func.gpio_mode_pre(dev, mode); err != nil {
			return fmt.Errorf("gpio: error running mode pre: %s", err)
		}
	}

	if dev.value_fp != nil {
		dev.value_fp.Close()
		dev.value_fp = nil
	}

	filepath := fmt.Sprintf(SYSFS_CLASS_GPIO+"/gpio%d/drive", dev.pin)
	drive, err := os.OpenFile(filepath, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("gpio: Failed to open drive for writing: %s", err)
	}
	defer drive.Close()

	if mode != MODE_GPIO_HIZ && mode != MODE_GPIO_PULLUP && mode != MODE_GPIO_STRONG && mode != MODE_GPIO_PULLDOWN {
		return fmt.Errorf("gpio: Entered mode not supported: %s", mode)
	}

	if _, err := drive.Write([]byte(mode)); err != nil {
		fmt.Errorf("gpio: Failed to write to drive mode: %s", err)
	}

	if advance_func.gpio_mode_post != nil {
		return advance_func.gpio_mode_post(dev, mode)
	}

	return nil
}

func gpio_owner(dev *gpio_context, own bool) error {
	if dev == nil {
		fmt.Errorf("gpio: context passed into gpio_owner is nil")
	}
	dev.owner = own
	return nil
}

func setup_mux_mapped(meta pin_t) error {
	var mi uint

	for mi = 0; mi < meta.mux_total; mi++ {
		var mux_i *gpio_context
		mux_i, err := GpioInitRaw(int(meta.mux[mi].pin))
		if err != nil {
			return fmt.Errorf("gpio: error setting up mux pins", err)
		}

		gpio_dir(mux_i, OUT)
		gpio_owner(mux_i, false)

		if err := gpio_write(mux_i, int(meta.mux[mi].value)); err != nil {
			gpio_close(mux_i)
			return fmt.Errorf("Error writing value to mux pin (index: %d)", mi)
		}
		gpio_close(mux_i)
	}

	return nil
}

func gpio_read(dev *gpio_context) (int, error) {
	if dev == nil {
		return -1, nil
	}

	// TODO(stephen) Do mmap stuff here

	if dev.value_fp == nil {
		if err := gpio_get_value_file_handle(dev); err != nil {
			return -1, fmt.Errorf("gpio: Failed to get value file pointer: %s", err)
		}
	} else {
		if _, err := dev.value_fp.Seek(0, os.SEEK_SET); err != nil {
			return -1, fmt.Errorf("gpio: error seeking on value file for pin[%d]:", dev.pin, err)
		}
	}

	// Should never be larger than this, in fact 4096 is probably too large
	var buf []byte = make([]byte, 0, 4096)
	if _, err := dev.value_fp.Read(buf); err != nil {
		return -1, fmt.Errorf("gpio: Error reading from pin[%d]: %s", dev.pin, err)
	}

	// Reset the seek for some reason
	_, _ = dev.value_fp.Seek(0, os.SEEK_SET)

	i, err := parseInt(string(buf), 32)
	if err != nil {
		return -1, fmt.Errorf("gpio: Can't convert buf to int: %s, %s", string(buf), err)
	}

	return int(i), nil
}

func GpioInit(pin int) (*gpio_context, error) {
	var err error

	if plat == nil {
		return nil, fmt.Errorf("gpio: platform is not initialized")
	}

	if pin < 0 || uint(pin) > plat.phy_pin_count {
		return nil, fmt.Errorf("gpio: pin %d beyond platform definition", pin)
	}

	// TODO(stephen) check capabilities

	if plat.pins[pin].gpio.mux_total > 0 {
		if err = setup_mux_mapped(plat.pins[pin].gpio); err != nil {
			return nil, fmt.Errorf("gpio: unable to setup muxes")
		}
	}

	r, err := GpioInitRaw(int(plat.pins[pin].gpio.pinmap))
	if err != nil {
		return nil, fmt.Errorf("gpio: gpio_init_raw(%d) return error: %s", pin, err)
	}
	r.phy_pin = pin

	if advance_func.gpio_init_post != nil {
		if err := advance_func.gpio_init_post(r); err != nil {
			return r, fmt.Errorf("gpio: init_post returned an error: %s", err)
		}
	}

	return r, nil
}

func GpioInitRaw(pin int) (*gpio_context, error) {
	fmt.Printf("gpio: Running GpioInitRwaw for pin %d\n", pin)
	if pin < 0 {
		return nil, fmt.Errorf("gpio: pin number < 0")
	}

	if advance_func.gpio_init_pre != nil {
		if err := advance_func.gpio_init_pre(pin); err != nil {
			return nil, fmt.Errorf("gpio: error in pre init: %s\n", err)
		}
	}

	dev := &gpio_context{}
	dev.pin = pin
	dev.phy_pin = -1
	dev.value_fp = nil

	// make sure the pin is exported
	directory := SYSFS_CLASS_GPIO + fmt.Sprintf("/gpio%d/", dev.pin)
	stat, err := os.Stat(directory)
	// if file already exists it's already owned
	// if file doesn't exist attempt to export the pin
	if err == nil && stat.IsDir() {
		dev.owner = false
	} else {
		toWrite := fmt.Sprintf("%d", dev.pin)
		_, err := writeFile(SYSFS_CLASS_GPIO+"/export", []byte(toWrite))
		if err != nil {
			return nil, fmt.Errorf("gpio: failed to write %d to export: %s\n", dev.pin, err)
		}
		dev.owner = true
	}

	return dev, nil
}
