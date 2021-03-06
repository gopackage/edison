# Intel Edison using Go

Simplified [Go][] (aka golang) development for the Intel Edison platform.
This library makes it fast and easy to develop on Edison modules. We provide
high level APIs to read/write/control gpio, i2c, and bluetooth (ble).

# Installation

We strongly recommend using [gpm][] to develop with this library. If you use
gpm, simply add `github.com/metamech/edison` to your `Godeps` file and run
`gpm install`.

If you want to develop on the edison library itself,
run `gpm install` to pull the testing frameworks we use,
then `go test ./...` normally.

# TODO

* i2c support
* documentation
* unit tests

[Go]: http://golang.org
[gpm]: https://github.com/pote/gpm
