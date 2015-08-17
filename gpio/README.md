# Intel Edison GPIO using Go

Simplified [Go][] (aka golang) development for the Intel Edison platform.
This library makes it easy to add gpio support to applications running on
Edison. The implementation is tested to work on both Arduino and Intel
mini-breakout boards. Note that if you build your own custom Edison breakout
board, your Edison will probably behave as a mini-breakout and this library
should work and auto-detect your board. If you do encounter problems on a custom
board, please try to force the board type when creating your GPIO pins.

# Installation

We strongly recommend using [gpm][] to develop with this library. If you use
gpm, simply add `github.com/metamech/edison/gpio` to your `Godeps` file and run
`gpm install`.

If you want to develop on the edison library itself,
run `gpm install` to pull the testing frameworks we use,
then `go test ./...` normally.

# TODO

* documentation
* unit tests

[Go]: http://golang.org
[gpm]: https://github.com/pote/gpm
