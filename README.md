# Intel Edison using Go

Simplified [Go][] (aka golang) development for the Intel Edison platform.
This library provides basic wrappers around common Edison Go applications
built to run on Edison modules.

# Installation

We strongly recommend using [gpm][] to develop with this library. If you use
gpm, simply add `github.com/metamech/edison` to your `Godeps` file and run
`gpm install`.

If you want to develop on the edison library itself,
run `gpm install` to pull the testing frameworks we use,
then `go test ./...` normally.

# TODO

* documentation
* unit tests

[Go]: http://golang.org
[gpm]: https://github.com/pote/gpm
