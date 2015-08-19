package main

import (
	"fmt"

	"github.com/gopackage/cli"
)

func enable(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	fmt.Printf("enable pin[%d]\n", pin)
}

func disable(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	fmt.Printf("disable pin[%d]\n", pin)
}

func pulse(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	duty := command.Args[1].IntValue(0)
	period := command.Args[2].IntValue(0)
	fmt.Printf("pulse on pin[%d] duty cycle: %d period: %d\n", pin, duty, period)
}

func main() {
	program := cli.New()
	program.SetVersion("0.1")

	program.Command("enable <pin>", "enable pwm control on <pin>").SetAction(enable)
	program.Command("disable <pin>", "disable pwm control on <pin>").SetAction(disable)
	program.Command("pulse <pin> <dutycycle> <period>", "control pwm pulse on <pin> for <dutycycle> over <period>").SetAction(pulse)

	program.Parse()
}
