// Copyright 2013 Daniel Jo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package frame implements a 2-D frame composed of cells printable to a
// terminal. This package is dependant upon the console package, but is not
// tied to any particular implementation of it.

package frame

import (
	"console"
	"image"
	"unicode/utf8"
)

// Cell is a single character in the Frame. It consists of a unicode rune and
// a console.Color
type Cell struct {
	// R is a unicode rune
	R rune
	// C is a color as defined in the Console package
	C console.Color
}

// Frame is a rectangular grid of Cells, modeled off of the API of the Image package in the standard library.
type Frame struct {
	Data   []Cell
	Bounds image.Rectangle
	Stride int
}

// New creates a new Frame bounded by a rectangle r
func New(r image.Rectangle) *Frame {
	var w, h = r.Dx(), r.Dy()
	return &Frame{
		Data:   make([]Cell, w*h),
		Bounds: r,
		Stride: r.Dx(),
	}
}

// FillRect sets each Cell in the Frame within the rectangle r to Cell c.
func (f *Frame) FillRect(r image.Rectangle, c Cell) {
	r = f.Bounds.Intersect(r)
	if r.Empty() {
		return
	}

	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			f.Set(x, y, c)
		}
	}
}

// Fill sets all Cells in the Frame to c
func (f *Frame) Fill(c Cell) {
	f.FillRect(f.Bounds, c)
}

// ClearRect sets all Cells bounded by r to an empty state. An empty Cell is
// contains a space character with its foreground and background colour set to
// black.
func (f *Frame) ClearRect(r image.Rectangle) {
	f.FillRect(r, Cell{' ', console.Color{}})
}

// Clear calls ClearRect with r set to the frame's bounds
func (f *Frame) Clear() {
	f.ClearRect(f.Bounds)
}

// CellOffset returns the index of a Cell within the Data field of Frame.
func (f *Frame) CellOffset(x, y int) int {
	return (y-f.Bounds.Min.Y)*f.Stride + (x - f.Bounds.Min.X)
}

// At returns the Cell located at coordinate x,y. 0,0 is the top-left corner. If
// the coordinate are not within the bounds of the Frame, a zeroed Cell is
// returned.
func (f *Frame) At(x, y int) Cell {
	if !image.Pt(x, y).In(f.Bounds) {
		return Cell{}
	}

	return f.Data[f.CellOffset(x, y)]
}

// Set replaces the Cell at coordinate x,y to c if the coordinate is within the
// bounds of the Frame.
func (f *Frame) Set(x, y int, c Cell) {
	var pt = image.Pt(x, y)
	if !pt.In(f.Bounds) {
		return
	}

	f.Data[f.CellOffset(x, y)] = c
}

// PutText writes strings into the Frame at position x, y. The input strings
// should not contain any escape sequences. Colors are added by interspersing
// console.Color structs between strings. Text that extends beyond the width of
// the Frame is truncated.
// TODO: Consider parsing strings with escape sequences to extract colours.
// This would unfortunately break potential compatibility with a Windows
// version of the Console package unless an agnostic set of escape sequences is
// developed.
func (f *Frame) PutText(x, y int, vals ...interface{}) {
	if y < f.Bounds.Min.Y || y >= f.Bounds.Max.Y {
		return
	}

	var col console.Color

	for _, v := range vals {
		if x >= f.Bounds.Max.X {
			break
		}

		switch o := v.(type) {
		case string:
			for _, r := range o {
				f.Set(x, y, Cell{r, col})
				x++
			}
		case console.Color:
			col = o
		default:
			panic("Frame.PutText accepts data of only type string and console.Color.")
		}
	}
}

// PutTextRel calls PutText with coordinates relative to the minimum of the
// Frame's bounds.
func (f *Frame) PutTextRel(x, y int, vals ...interface{}) {
	f.PutText(f.Bounds.Min.X+x, f.Bounds.Min.Y+y, vals...)
}

// A reusable buffer meant to minimize future allocations.
var printBuf []byte

// PrintRect outputs all cells within the rectangle r to the Console c.
func (f *Frame) PrintRect(c *console.Console, r image.Rectangle) {
	var (
		col console.Color
	)

	printBuf = printBuf[:0]

	for y := r.Min.Y; y < r.Max.Y; y++ {
		printBuf = append(printBuf, console.FormatMoveTo(y+1, r.Min.X)...)

		for x := r.Min.X; x < r.Max.X; x++ {
			var cell = f.At(x, y)

			if cell.R == 0 {
				cell.R = ' '
			}

			if col != cell.C {
				col = cell.C
				printBuf = append(printBuf, col.String()...)
			}

			var (
				buf [4]byte
				n = utf8.EncodeRune(buf[:], cell.R)
			)
			printBuf = append(printBuf, buf[:n]...)
		}
	}
	c.WriteString(string(printBuf))
}

// Print outputs the entire Frame to Console c.
func (f *Frame) Print(c *console.Console) {
	f.PrintRect(c, f.Bounds)
}

// SubFrame returns a Frame bounded by the rectangle r, sharing the data with
// the receiver.
func (f *Frame) SubFrame(r image.Rectangle) *Frame {
	r = r.Intersect(f.Bounds)

	if r.Empty() {
		return new(Frame)
	}

	return &Frame{
		Data: f.Data[f.CellOffset(r.Min.X, r.Min.Y):],
		Bounds: r,
		Stride: f.Stride,
	}
}

