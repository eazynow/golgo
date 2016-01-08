// Googles implementation of Conway's Game of Life, now with added features
// For the original version see https://golang.org/doc/play/life.go
package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const (
	gameBorderCol   = 2
	gameBorderRow   = 1
	textColor       = termbox.ColorWhite
	backgroundColor = termbox.ColorBlue
	cellColor       = termbox.ColorWhite
	boardColor      = termbox.ColorBlack
)

var titles = []string{
	"Conways Game of Life",
	"--------------------",
	"eazynow 2016",
}

var instructions = []string{
	"Controls:",
	"",
	"p    pause",
	"s    step",
	"r    randomize",
	"",
	"q    quit",
	"",
}

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {

	f := Field{w: w, h: h}
	f.Clear()

	return &f
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

func (f *Field) Clear() {
	f.s = make([][]bool, f.h)
	for i := range f.s {
		f.s[i] = make([]bool, f.w)
	}
}

func (f *Field) Randomize() {
	f.Clear()
	for i := 0; i < (f.w * f.h / 4); i++ {
		f.Set(rand.Intn(f.w), rand.Intn(f.h), true)
	}
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b       *Field
	w, h       int
	generation int
	paused     bool
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) (*Life, error) {

	if err := termbox.Init(); err != nil {
		return nil, err
	}

	a := NewField(w, h)

	a.Randomize()

	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
		generation: 1,
		paused:     false,
	}, nil
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
	l.generation++
}

func (l *Life) Close() {
	termbox.Close()
}

func (l *Life) Pause() {
	l.paused = !l.paused
}

func (l *Life) IsPaused() bool {
	return l.paused
}

// Render renders the board in termbox
func (l *Life) Render() {
	termbox.Clear(backgroundColor, backgroundColor)

	titleX := l.w + (gameBorderCol * 2)
	titleY := gameBorderRow
	for y, t := range titles {
		tbprint(titleX, titleY+y, textColor, backgroundColor, t)
	}

	instrucX := titleX
	instrucY := titleY + len(titles) + 2
	for y, i := range instructions {
		tbprint(instrucX, instrucY+y, textColor, backgroundColor, i)
	}

	pauseX := titleX
	pauseY := gameBorderRow + l.h - 2

	pauseMsg := "RUNNING"
	if l.paused {
		pauseMsg = "PAUSED"
	}

	tbprint(pauseX, pauseY, textColor, backgroundColor, pauseMsg)

	genX := titleX
	genY := pauseY + 1

	tbprint(genX, genY, textColor, backgroundColor, fmt.Sprintf("Generation: %d", l.generation))

	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			color := boardColor
			if l.a.Alive(x, y) {
				color = cellColor
			}
			termbox.SetCell(x+gameBorderCol, y+gameBorderRow, ' ', color, color)
		}
	}

	termbox.Flush()
}

func (l *Life) Randomize() {
	l.a.Randomize()
	l.generation = 1
}

func (l *Life) Run() {
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Ch == 'p':
					l.Pause()
				case ev.Ch == 's':
					if !l.IsPaused() {
						l.Pause()
					}
					l.Step()
				case ev.Ch == 'r':
					l.Randomize()
				case ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD:
					return
				}
			}
		default:
			if !l.paused {
				l.Step()
			}
			l.Render()
			time.Sleep(time.Second / 30)
		}
	}
}

func main() {
	l, err := NewLife(80, 30)
	if err != nil {
		fmt.Printf("Could not start the game: %s\n", err.Error())
		return
	}

	defer l.Close()

	l.Run()
}

// Function tbprint draws a string.
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
