package main

import (
	"fmt"

	"github.com/gopackage/cli"
	"github.com/gopackage/edison/gpio"
	"github.com/gopackage/sysfs"
)

func enable(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	mode := command.Args[0].Value
	fmt.Printf("enable pin[%d] to mode %s\n", pin, mode)
}

func disable(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	fmt.Printf("disable pin[%d]\n", pin)
}

func read(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	fmt.Printf("read pin[%d]", pin)
}

func set(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	state := command.Args[1].Value
	fmt.Printf("set pin[%d] to %s\n", pin, state)
}

func main() {
	program := cli.New()
	program.SetVersion("0.1")

	fmt.Printf("Running gpio init\n")
	pins := gpio.NewPins(&sysfs.HardwareFile{})
	err := pins.Init()
	if err != nil {
		fmt.Printf("Init failed: %s\n", err)
	} else {
		fmt.Printf("Init success\n")
	}

	program.Command("enable <pin> <mode>", "enable pwm control on <pin>").SetAction(enable)
	program.Command("disable <pin>", "disable pwm control on <pin>").SetAction(disable)
	program.Command("read <pin>", "read pin state on <pin>").SetAction(read)
	program.Command("set <pin> <state>", "set pin state on <pin> to <state>").SetAction(set)

	program.Parse()
}
