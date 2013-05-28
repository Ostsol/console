// Copyright 2013 Daniel Jo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"strconv"
)

// Colour attributes
const (
	RESET = iota
	BRIGHT
	DIM
	_SKIPPED3
	UNDERSCORE
	BLINK
	_SKIPPED6
	REVERSE
	HIDDEN
)

// Colours
const (
	BLACK = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

// Color defines a VT100-compatible colour.
type Color struct {
	Attr, Fore, Back uint8
}

// String returns an escape sequence representing the Color.
func (c Color) String() string {
	return FormatColor(c.Attr, c.Fore, c.Back)
}

// FormatColor returns an escape sequence representing the style defined by the
// attribute attr, foreground colour fore, and background colour back.
func FormatColor(attr, fore, back uint8) string {
	if attr < 0 || attr > 8 {
		panic("colorString: invalid attribute")
	}
	if fore < 0 || fore > 7 {
		panic("colorString: invalid foreground colour")
	}
	if back < 0 || back > 7 {
		panic("colorString: invalid background colour")
	}
	var (
		astring = strconv.FormatInt(int64(attr), 10)
		fstring = strconv.FormatInt(int64(fore+30), 10)
		bstring = strconv.FormatInt(int64(back+40), 10)
	)
	return _CSI + astring + ";" + fstring + ";" + bstring + "m"
}

