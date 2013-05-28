// Copyright 2013 Daniel Jo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package console defines an API for manipulating a terminal. Currently it is
// built around VT100 escape sequences, though it is conceivable that the API
// may be implemented to work with other terminals.
//
// TODO: Provide a means to accept Unicode input.
package console

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

const (
	_ESC  = "\033"
	_CSI  = _ESC + "["
	clear = _CSI + "2J"
)

// Console is an interface to a terminal, defined by an input stream and an
// output stream.
type Console struct {
	in  io.Reader
	out io.Writer
}

// New returns a Console that receives input from Reader in and outputs to
// Writer out.
func New(in io.Reader, out io.Writer) *Console {
	return &Console{in: in, out: out}
}

// Clear writes the clear escape sequence.
func (c *Console) Clear() {
	c.out.Write([]byte(clear))
}

// MoveUp moves the cursor up by i spaces.
func (c *Console) MoveUp(i int) {
	c.out.Write([]byte(FormatMoveUp(i)))
}

// MoveUp moves the cursor down by i spaces.
func (c *Console) MoveDown(i int) {
	c.out.Write([]byte(FormatMoveDown(i)))
}

// MoveUp moves the cursor right by i spaces.
func (c *Console) MoveRight(i int) {
	c.out.Write([]byte(FormatMoveRight(i)))
}

// MoveUp moves the cursor left by i spaces.
func (c *Console) MoveLeft(i int) {
	c.out.Write([]byte(FormatMoveLeft(i)))
}

// MoveTo moves the cursor to the specified line and column.
func (c *Console) MoveTo(line, column int) {
	c.out.Write([]byte(FormatMoveTo(line, column)))
}

// SetColor sets the current printing colour to col.
func (c *Console) SetColor(col Color) {
	c.out.Write([]byte(col.String()))
}

// PutRune writes the Unicode rune r to the specified line and column.
func (c *Console) PutRune(line, column int, r rune) {
	c.MoveTo(line, column)
	c.WriteRune(r)
}

// WriteRune writes the Unicode rune r to the current cursor location.
func (c *Console) WriteRune(r rune) {
	var (
		bytes [4]byte
		l     int
	)
	l = utf8.EncodeRune(bytes[:], r)
	c.out.Write(bytes[:l])
}

// PutString writes the string str to the specified line and column.
func (c *Console) PutString(line, column int, str string) {
	c.out.Write([]byte(FormatMoveTo(line, column) + str))
}

// PutStringf calls fmt.Sprintf to format the string s with arguments args and
// writes the result to the specified line and column.
func (c *Console) PutStringf(line, column int, s string, args ...interface{}) {
	c.out.Write([]byte(FormatMoveTo(line, column) + fmt.Sprintf(s, args...)))
}

// WriteString writes the string str to the current cursor location.
func (c *Console) WriteString(str string) {
	c.out.Write([]byte(str))
}

// WriteStringf calls fmt.Sprintf to format the string s with arguments args and
// writes the result to the current cursor location.
func (c *Console) WriteStringf(s string, args ...interface{}) {
	c.out.Write([]byte(fmt.Sprintf(s, args...)))
}

// HideCursor prevents the terminal from rendering the cursor.
func (c *Console) HideCursor() {
	c.out.Write([]byte(_CSI + "?25l"))
}

// ShowCursor permits the terminal to render the cursor
func (c *Console) ShowCursor() {
	c.out.Write([]byte(_CSI + "?25h"))
}

// keybuf is a buffer for reading input.
var keybuf [16]byte

// GetKey reads a keystroke from the Console's input stream and returns its key
// code. There is no current support for reading Unicode runes.
func (c *Console) GetKey() int32 {
	n, _ := c.in.Read(keybuf[:])

	return parseKey(keybuf[:n])
}

// AltBuffer switches to the alternate terminal buffer.
func (c *Console) AltBuffer() {
	c.out.Write([]byte(_CSI + "?47h"))
}

// MainBuffer switches to the main terminal buffer.
func (c *Console) MainBuffer() {
	c.out.Write([]byte(_CSI + "?47l"))
}

// FormatClear returns the escape sequence that clears the terminal
func FormatClear() string {
	return clear
}

// FormateMoveUp returns the escape sequence that moves the cursor up i spaces.
func FormatMoveUp(i int) string {
	var istring = strconv.FormatInt(int64(i), 10)
	return _CSI + istring + "A"
}

// FormateMoveDown returns the escape sequence that moves the cursor down i
// spaces.
func FormatMoveDown(i int) string {
	var istring = strconv.FormatInt(int64(i), 10)
	return _CSI + istring + "B"
}

// FormateMoveRight returns the escape sequence that moves the cursor right i
// spaces.
func FormatMoveRight(i int) string {
	var istring = strconv.FormatInt(int64(i), 10)
	return _CSI + istring + "C"
}

// FormateMoveLeft returns the escape sequence that moves the cursor left i
// spaces.
func FormatMoveLeft(i int) string {
	var istring = strconv.FormatInt(int64(i), 10)
	return _CSI + istring + "D"
}

// FormatMoveTo returns the escape sequence that moves the cursor to the
// specified line and column.
func FormatMoveTo(line, column int) string {
	var (
		lstring = strconv.FormatInt(int64(line), 10)
		cstring = strconv.FormatInt(int64(column), 10)
	)
	return _CSI + lstring + ";" + cstring + "H"
}

// Init initializes the terminal to a suitable mode.
func Init() error {
	var (
		term *termios
		err error
	)
	if term, err = getTermios(); err != nil { return err }
	term.rawMode()
	return term.set()
}

// Exit returns the terminal to its default settings.
func Exit() error {
	return defaultTermios.set()
}
