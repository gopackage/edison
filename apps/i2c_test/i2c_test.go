package main

import (
	"fmt"
	"github.com/gopackage/edison/mraa"
)

func main() {
	_ = mraa.Init()
	var bus uint = 0

	_, err := mraa.I2cInit(bus)
	if err != nil {
		fmt.Printf("Error initing bus: %d\n", bus)
	}
}
