# golgo - Conways Game of Life in Go

This expands upon the simple example of [Conways Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) 
written in go [provided by Google](https://golang.org/doc/play/life.go)
by adding a richer interface (thanks to termbox) and more features

![alt text][logo]
[logo]: https://raw.githubusercontent.com/eazynow/golgo/master/screenshot.png "Screenshot"


## Requirements

golgo should work with go versions 1.2 or greater.


## Installation

To install golgo, use `go get`:
```
go get github.com/eazynow/golgo
```

## Usage

To play the simulation, just run

```
$ $GOPATH/bin/golgo
```


## Todo

This is still work in progress. Amongst other things:

+ Refactor to a proper solution (split out the code into separate files etc)
+ Add tests
+ Add import starting position
+ Make board size configurable (flags)
+ Make color scheme configurable?


## Credits

+ [googles example](https://golang.org/doc/play/life.go) implementation of game of life is used as a basic engine
+ [termbox](https://github.com/nsf/termbox-go) is used as the display engine
