package main

import (
	"fmt"

	"github.com/gopackage/cli"
)

func scan(program *cli.Program, command *cli.Command, unknownArgs []string) {
	addr := command.Args[0].IntValue(0)
	fmt.Printf("scan addr[%d]\n", addr)
}

func read(program *cli.Program, command *cli.Command, unknownArgs []string) {
	addr := command.Args[0].IntValue(0)
	fmt.Printf("read addr[%d]\n", addr)
}

func write(program *cli.Program, command *cli.Command, unknownArgs []string) {
	addr := command.Args[0].IntValue(0)
	fmt.Printf("write to addr[%d]\n", addr)
}

func main() {
	program := cli.New()
	program.SetVersion("0.1")

	program.Command("scan [start] [end]", "scan the i2c bus for devices from the start to end addresses").SetAction(scan)
	program.Command("read <addr> [length]", "read a certain amount of bytes from the given address").SetAction(read)
	program.Command("write <addr> [bytes]", "write bytes to a given address").SetAction(write)

	program.Parse()
}
