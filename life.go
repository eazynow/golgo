// Googles implementation of Conway's Game of Life, updated to display in termbox.
// For the original version see https://golang.org/doc/play/life.go
package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const backgroundColor = termbox.ColorBlue
const cellColor = termbox.ColorWhite

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
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
	a, b   *Field
	w, h   int
	paused bool
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) (*Life, error) {

	if err := termbox.Init(); err != nil {
		return nil, err
	}

	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
		paused: false,
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
}

func (l *Life) Close() {
	termbox.Close()
}

func (l *Life) Pause() {
	l.paused = !l.paused
}

// Render renders the board in termbox
func (l *Life) Render() {
	termbox.Clear(backgroundColor, backgroundColor)

	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			if l.a.Alive(x, y) {
				termbox.SetCell(x, y, ' ', cellColor, cellColor)
			}
		}
	}

	termbox.Flush()
}

func main() {
	l, err := NewLife(80, 30)
	if err != nil {
		fmt.Printf("Could not start the game: %s\n", err.Error())
		return
	}

	defer l.Close()

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
				case ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD:
					return
				}
			}
		default:
			if !l.paused {
				l.Step()
				l.Render()
			}
			time.Sleep(time.Second / 30)
		}
	}
}
