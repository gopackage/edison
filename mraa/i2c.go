package mraa

import (
	"fmt"
	"os"
)

var (
	I2C_STD  = "std"  // up to 100 Khz
	I2C_FAST = "fast" // up to 400 Khz
	I2C_HIGH = "high" // up to 3.4Mhz
)

type i2c_context struct {
	busnum uint     // the bus number of the /dev/i2c-* device
	fh     *os.File // the file handle to the /dev/i2c-* device
	addr   int      // the address of the i2c slave
	funcs  uint64   // dev/i2c-* device capabilities as per https://www.kernel.org/doc/Documentation/i2c/functionality
}

type I2cContext *i2c_context

type i2c_smbus_ioctl_data_struct struct {
	read_write uint8  // operation direction
	command    uint8  // ioctl command
	size       int    // data size
	data       []byte // data
}

func i2c_init_internal(bus uint) (*i2c_context, error) {
	var err error

	dev := &i2c_context{}
	dev.busnum = bus

	if advance_func.i2c_init_pre != nil {
		if err = advance_func.i2c_init_pre(bus); err != nil {
			goto init_internal_cleanup
		}
	}

	if advance_func.i2c_init_bus_replace != nil {
		if err = advance_func.i2c_init_bus_replace(dev); err != nil {
			goto init_internal_cleanup
		}
	} else {
		fp := fmt.Sprintf("/dev/i2c-%d", bus)
		if dev.fh, err = os.OpenFile(fp, os.O_RDWR, 0664); err != nil {
			err = fmt.Errorf("i2c: Failed to open requested i2c port %s, %s", fp, err)
			goto init_internal_cleanup
		}
	}

init_internal_cleanup:
	if err != nil {
		return nil, err
	} else {
		return dev, nil
	}
}

func I2cInit(bus int) (*i2c_context, error) {
	var board *board_t = plat
	if board == nil {
		return nil, fmt.Errorf("i2c: Platform Not Initialized")
	}

	return i2c_init_internal(uint(board.i2c_bus[bus].bus_id))
}
