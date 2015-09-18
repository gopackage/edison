package gpio

const (
	UART_DEV_PATH = "/dev/ttyMDF1"
)

func InitMiniboard(b *board_t) {
	miniboard = true
	b.phy_pin_count = 56
	b.aio_count = 0
	b.pwm_default_period = 5000
	b.pwm_max_period = 218453
	b.pwm_min_period = 1

	b.pins = make([]pininfo_t, 56)

	advance_func.gpio_init_post = intel_edison_gpio_init_post
	advance_func.pwm_init_pre = intel_edison_pwm_init_pre
	advance_func.gpio_mode_replace = intel_edison_gpio_mode_replace

	pos := 0
	b.pins[pos].name = "J17-1"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 1, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 182
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].pwm.pinmap = 2
	b.pins[pos].pwm.parent_id = 0
	b.pins[pos].pwm.mux_total = 0
	pos++

	b.pins[pos].name = "J17-2"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J17-3"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J17-4"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J17-5"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 135
	b.pins[pos].gpio.mux_total = 0
	pos++

	b.pins[pos].name = "J17-6"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J17-7"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 1, 0, 0)
	b.pins[pos].gpio.pinmap = 27
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].i2c.pinmap = 1
	b.pins[pos].i2c.mux_total = 0
	pos++

	b.pins[pos].name = "J17-8"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 1, 0, 0)
	b.pins[pos].gpio.pinmap = 20
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].i2c.pinmap = 1
	b.pins[pos].i2c.mux_total = 0
	pos++

	b.pins[pos].name = "J17-9"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 1, 0, 0)
	b.pins[pos].gpio.pinmap = 28
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].i2c.pinmap = 1
	b.pins[pos].i2c.mux_total = 0
	pos++

	b.pins[pos].name = "J17-10"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 1, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 111
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].spi.pinmap = 5
	b.pins[pos].spi.mux_total = 0
	pos++

	b.pins[pos].name = "J17-11"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 1, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 109
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].spi.pinmap = 5
	b.pins[pos].spi.mux_total = 0
	pos++

	b.pins[pos].name = "J17-12"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 1, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 115
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].spi.pinmap = 5
	b.pins[pos].spi.mux_total = 0
	pos++
	b.pins[pos].name = "J17-13"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J17-14"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 128
	b.pins[pos].gpio.parent_id = 0
	b.pins[pos].gpio.mux_total = 0
	pos++

	b.pins[pos].name = "J18-1"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 1, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 13
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].pwm.pinmap = 1
	b.pins[pos].pwm.parent_id = 0
	b.pins[pos].pwm.mux_total = 0
	pos++

	b.pins[pos].name = "J18-2"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 165
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J18-3"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J18-4"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J18-5"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J18-6"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 1, 0, 0)
	b.pins[pos].gpio.pinmap = 19
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].i2c.pinmap = 1
	b.pins[pos].i2c.mux_total = 0
	pos++

	b.pins[pos].name = "J18-7"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 1, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 12
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].pwm.pinmap = 0
	b.pins[pos].pwm.parent_id = 0
	b.pins[pos].pwm.mux_total = 0
	pos++

	b.pins[pos].name = "J18-8"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 1, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 183
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].pwm.pinmap = 3
	b.pins[pos].pwm.parent_id = 0
	b.pins[pos].pwm.mux_total = 0
	pos++
	b.pins[pos].name = "J18-9"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J18-10"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 1, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 110
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].spi.pinmap = 5
	b.pins[pos].spi.mux_total = 0
	pos++
	b.pins[pos].name = "J18-11"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 1, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 114
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].spi.pinmap = 5
	b.pins[pos].spi.mux_total = 0
	pos++

	b.pins[pos].name = "J18-12"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 129
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J18-13"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 1)
	b.pins[pos].gpio.pinmap = 130
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].uart.pinmap = 0
	b.pins[pos].uart.parent_id = 0
	b.pins[pos].uart.mux_total = 0

	pos++
	b.pins[pos].name = "J18-14"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J19-1"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J19-2"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J19-3"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J19-4"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 44
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J19-5"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 46
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J19-6"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 48
	b.pins[pos].gpio.mux_total = 0
	pos++

	b.pins[pos].name = "J19-7"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++

	b.pins[pos].name = "J19-8"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 1)
	b.pins[pos].gpio.pinmap = 131
	b.pins[pos].gpio.mux_total = 0
	b.pins[pos].uart.pinmap = 0
	b.pins[pos].uart.parent_id = 0
	b.pins[pos].uart.mux_total = 0
	pos++

	b.pins[pos].name = "J19-9"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 14
	b.pins[pos].gpio.mux_total = 0
	pos++

	b.pins[pos].name = "J19-10"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 40
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J19-11"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 43
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J19-12"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 77
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J19-13"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 82
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J19-14"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 83
	b.pins[pos].gpio.mux_total = 0
	pos++

	b.pins[pos].name = "J20-1"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J20-2"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J20-3"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 0, 0, 0, 0, 0, 0, 0)
	pos++
	b.pins[pos].name = "J20-4"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 45
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-5"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 47
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-6"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 49
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-7"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 15
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-8"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 84
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-9"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 42
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-10"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 41
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-11"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 78
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-12"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 79
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-13"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 80
	b.pins[pos].gpio.mux_total = 0
	pos++
	b.pins[pos].name = "J20-14"
	b.pins[pos].capabilities = mraa_pincapabilities(1, 1, 0, 0, 0, 0, 0, 0)
	b.pins[pos].gpio.pinmap = 81
	b.pins[pos].gpio.mux_total = 0
	pos++

	// BUS DEFINITIONS
	b.i2c_bus_count = 9
	b.def_i2c_bus = 6
	var ici int
	for ici = 0; ici < 9; ici++ {
		b.i2c_bus[ici].bus_id = -1
	}
	b.i2c_bus[1].bus_id = 1
	b.i2c_bus[1].sda = 7
	b.i2c_bus[1].scl = 19

	b.i2c_bus[6].bus_id = 6
	b.i2c_bus[6].sda = 8
	b.i2c_bus[6].scl = 6

	b.spi_bus_count = 1
	b.def_spi_bus = 0
	b.spi_bus[0].bus_id = 5
	b.spi_bus[0].slave_s = 1
	b.spi_bus[0].cs = 23
	b.spi_bus[0].mosi = 11
	b.spi_bus[0].miso = 24
	b.spi_bus[0].sclk = 10

	b.uart_dev_count = 1
	b.def_uart_dev = 0
	b.uart_dev[0].rx = 26
	b.uart_dev[0].tx = 35
	b.uart_dev[0].device_path = UART_DEV_PATH

	plat = b
}

func InitEdison() error {
	board := board_t{}
	/*tristate*/ _, err := GpioInitRaw(214)
	if err == nil {
		InitMiniboard(&board)
	} else {
		//InitArduino
	}

	return nil
}
