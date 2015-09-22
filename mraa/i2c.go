package mraa

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	I2C_STD  = "std"  // up to 100 Khz
	I2C_FAST = "fast" // up to 400 Khz
	I2C_HIGH = "high" // up to 3.4Mhz

	I2C_TENBIT         = 0x0704
	I2C_FUNCS  uintptr = 0x0705
	I2C_RDWR           = 0x0707
	I2C_PEC            = 0x0708

	I2C_SMBUS uintptr = 0x0720

	I2C_NOCMD             = 0
	I2C_SMBUS_READ        = 1
	I2C_SMBUS_WRITE uint8 = 0

	I2C_SMBUS_QUICK            = 0
	I2C_SMBUS_BYTE             = 1
	I2C_SMBUS_BYTE_DATA        = 2
	I2C_SMBUS_WORD_DATA        = 3
	I2C_SMBUS_PROC_CALL        = 4
	I2C_SMBUS_BLOCK_DATA       = 5
	I2C_SMBUS_I2C_BLOCK_BROKEN = 6
	I2C_SMBUS_BLOCK_PROC_CALL  = 7
	I2C_SMBUS_I2C_BLOCK_DATA   = 8

	I2C_FUNC_I2C                    = 0x00000001
	I2C_FUNC_10BIT_ADDR             = 0x00000002
	I2C_FUNC_PROTOCOL_MANGLING      = 0x00000004
	I2C_FUNC_SMBUS_PEC              = 0x00000008
	I2C_FUNC_SMBUS_BLOCK_PROC_CALL  = 0x00008000
	I2C_FUNC_SMBUS_QUICK            = 0x00010000
	I2C_FUNC_SMBUS_READ_BYTE        = 0x00020000
	I2C_FUNC_SMBUS_WRITE_BYTE       = 0x00040000
	I2C_FUNC_SMBUS_READ_BYTE_DATA   = 0x00080000
	I2C_FUNC_SMBUS_WRITE_BYTE_DATA  = 0x00100000
	I2C_FUNC_SMBUS_READ_WORD_DATA   = 0x00200000
	I2C_FUNC_SMBUS_WRITE_WORD_DATA  = 0x00400000
	I2C_FUNC_SMBUS_PROC_CALL        = 0x00800000
	I2C_FUNC_SMBUS_READ_BLOCK_DATA  = 0x01000000
	I2C_FUNC_SMBUS_WRITE_BLOCK_DATA = 0x02000000
	I2C_FUNC_SMBUS_READ_I2C_BLOCK   = 0x04000000
	I2C_FUNC_SMBUS_WRITE_I2C_BLOCK  = 0x08000000

	I2C_M_TEN          = 0x10
	I2C_M_RD           = 0x01
	I2C_M_NOSTART      = 0x4000
	I2C_M_REV_DIR_ADDR = 0x2000
	I2C_M_IGNORE_NAK   = 0x1000
	I2C_M_NO_RD_ACK    = 0x0800

	I2C_SMBUS_I2C_BLOCK_MAX = 32
)

type i2c_context struct {
	busnum uint     // the bus number of the /dev/i2c-* device
	fh     *os.File // the file handle to the /dev/i2c-* device
	addr   int      // the address of the i2c slave
	funcs  uintptr  // dev/i2c-* device capabilities as per https://www.kernel.org/doc/Documentation/i2c/functionality
}

type I2cContext *i2c_context

type i2c_smbus_ioctl_data_struct struct {
	read_write uint8  // operation direction
	command    uint8  // ioctl command
	size       int    // data size
	data       []byte // data
}

type i2c_msg struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   []byte
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

		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, dev.fh.Fd(), I2C_FUNCS, dev.funcs)
		if errno != 0 {
			fmt.Printf("Failed to get I2C_FUNC map from device: %d\n", errno)
			dev.funcs = 0
		}
	}

	if advance_func.i2c_init_post != nil {
		if err = advance_func.i2c_init_post(dev); err != nil {
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

func I2cInit(bus uint) (*i2c_context, error) {
	var board *board_t = plat
	if board == nil {
		return nil, fmt.Errorf("i2c: Platform Not Initialized")
	}

	if board.i2c_bus_count == 0 {
		return nil, fmt.Errorf("No i2c buses defined in platform")
	}

	if bus >= board.i2c_bus_count {
		return nil, fmt.Errorf("Above i2c bus count")
	}

	if board.i2c_bus[bus].bus_id == -1 {
		fmt.Printf("Invalid i2c bus, moving to default i2c bus\n")
		bus = board.def_i2c_bus
	}

	pos := board.i2c_bus[bus].sda
	if board.pins[pos].i2c.mux_total > 0 {
		if err := setup_mux_mapped(board.pins[pos].i2c); err != nil {
			return nil, fmt.Errorf("i2c: Failed to set-up i2c scl multiplexer")
		}
	}

	return i2c_init_internal(uint(board.i2c_bus[bus].bus_id))
}

func I2cInitRaw(bus uint) (*i2c_context, error) {
	return i2c_init_internal(bus)
}

func (dev *i2c_context) Read() ([]byte, error) {
	data := make([]byte, 0, 256)
	_, err := dev.fh.Read(data)
	return data, err
}

func (dev *i2c_context) ReadBytes(addr uint, length int) ([]byte, error) {

}

func (dev *i2c_context) Frequency(mode string) error {
	if advance_func.i2c_set_frequency_replace != nil {
		return advance_func.i2c_set_frequency_replace(dev, mode)
	}
	return fmt.Errorf("i2c: Frequency set not supported")
}

func i2c_smbus_access(fd uintptr, read_write, command uint8, size int, data []byte) int {
	var args i2c_smbus_ioctl_data_struct

	args.read_write = read_write
	args.command = command
	args.size = size
	args.data = data

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, I2C_SMBUS, uintptr(unsafe.Pointer(&args)))
	return int(errno)
}

func (dev *i2c_context) Write(command byte, data []byte) int {
	length := len(data) - 1
	if length > I2C_SMBUS_I2C_BLOCK_MAX {
		length = I2C_SMBUS_I2C_BLOCK_MAX
	}

	return i2c_smbus_access(dev.fh.Fd(), I2C_SMBUS_WRITE, command, I2C_SMBUS_I2C_BLOCK_DATA, data)
}
