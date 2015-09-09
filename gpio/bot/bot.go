package bot

import "github.com/gopackage/sysfs"

// Bot automates configuring and sending multi-step pin commands when the pin
// itself isn't needed. This is useful for sending commands to setup/mux pins
// that aren't going to be used by the application itself. The Bot API is
// designed to make it easy to understand what steps are happening in code (as
// opposed to being efficient or fast to execute).
type Bot struct {
	file  sysfs.DeviceFile
	steps []*botstep
}

// NewBot creates a new pin bot for automating pin one-off commands.
func NewBot(file sysfs.DeviceFile) *Bot {
	return &Bot{file: file}
}

// Add queues a command to execute with optional arguments.
func (p *Bot) Add(pin int, commanders ...CommandFunc) {
	for _, cmd := range commanders {
		p.steps = append(p.steps, &botstep{cmd, pin})
	}
}

// Run starts the Bot executing the commands it was configured for. It stops
// when the first command returns an error.
func (p *Bot) Run() error {
	for _, step := range p.steps {
		path, value := step.commander(step.pin)
		err := p.file.Write(path, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// CommandFunc is used to generate the path and value for each bot command.
type CommandFunc func(pin int) (path, value string)

type botstep struct {
	commander func(pin int) (path, value string)
	pin       int
}
