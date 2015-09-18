package main

import (
	"fmt"
	"time"

	"github.com/gopackage/cli"
	"github.com/gopackage/edison/mraa"
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

func cycle(program *cli.Program, command *cli.Command, unknownArgs []string) {
	pin := command.Args[0].IntValue(0)
	period := command.Args[1].IntValue(0)
	fmt.Printf("cycle pin[%d]\n", pin)

	pwm, err := mraa.PwmInit(pin)
	if err != nil {
		fmt.Printf("Error opening pwm pin %d: %s\n", pin, err)
		return
	}

	if err = pwm.Period(period); err != nil {
		fmt.Printf("Error setting pwm period to 200 for pin %d: %s\n", pin, err)
	}

	if err = pwm.Enable(1); err != nil {
		fmt.Printf("Error enabling pin %d: %s\n", pin, err)
	}

	var value int

	for {
		value++
		pwm.Scale(value)
		time.Sleep(50000 * time.Microsecond)
		if value > 255 {
			value = 0
		}
		duty, err := pwm.Read()
		if err != nil {
			fmt.Printf("Error from PwmRead: %s\n", err)
		} else {
			fmt.Printf("New duty for pin[%d]: %f\n", pin, duty)
		}
	}
}

func main() {
	program := cli.New()
	program.SetVersion("0.1")

	_ = mraa.Init()

	program.Command("enable <pin>", "enable pwm control on <pin>").SetAction(enable)
	program.Command("disable <pin>", "disable pwm control on <pin>").SetAction(disable)
	program.Command("pulse <pin> <dutycycle> <period>", "control pwm pulse on <pin> for <dutycycle> over <period>").SetAction(pulse)
	program.Command("cycle <pin> <period>", "cycles a pwm through a range of color values on <pin> over <period>").SetAction(cycle)

	program.Parse()
}
