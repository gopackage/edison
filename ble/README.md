# Intel Edison Bluetooth using Go

Simplified [Go][] (aka golang) development for the Intel Edison platform.
This library makes it easy to add bluetooth support to applications running on
Edison. The implementation uses DBus to communicate with the Bluez bluetooth
stack and assumes a fairly standard Yocto build.

# Installation

We strongly recommend using [gpm][] to develop with this library. If you use
gpm, simply add `github.com/metamech/edison/ble` to your `Godeps` file and run
`gpm install`.

If you want to develop on the edison library itself,
run `gpm install` to pull the testing frameworks we use,
then `go test ./...` normally.

# TODO

* documentation
* unit tests

[Go]: http://golang.org
[gpm]: https://github.com/pote/gpm
