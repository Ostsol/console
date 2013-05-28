// Copyright 2013 Daniel Jo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
)

// Key codes. Unprintable codes are given values within the bounds of the
// "private use area" of the Unicode standard [U+E000..U+F8FF].
const (
	K_TAB       = 9
	K_ENTER     = 13
	K_ESCAPE    = 27
	K_BACKSPACE = 127

	K_UP = 0xE000 + iota
	K_DOWN
	K_RIGHT
	K_LEFT
	K_HOME
	K_END
	K_PAGEUP
	K_PAGEDOWN
	K_INSERT
	K_DELETE
	K_KPHOME
	K_KPEND
	K_KPPAGEUP
	K_KPPAGEDOWN
)

func parseSS3(buf []byte) int32 {
	if len(buf) < 1 {
		return 0
	}

	switch buf[0] {
	case 'H':
		return K_HOME
	case 'F':
		return K_END
	default:
		return 0
	}

	panic("unreachable")
}

func parseCSI(buf []byte) int32 {
	if len(buf) < 1 {
		return 0
	}

	switch string(buf) {
	case "A":
		return K_UP
	case "B":
		return K_DOWN
	case "C":
		return K_RIGHT
	case "D":
		return K_LEFT
	case "2~":
		return K_INSERT
	case "3~":
		return K_DELETE
	case "5~":
		return K_PAGEUP
	case "6~":
		return K_PAGEDOWN
	default:
		return 0
	}

	panic("unreachable")
}

func parseESC(buf []byte) int32 {
	if len(buf) < 1 {
		return K_ESCAPE
	}

	switch buf[0] {
	case '[':
		return parseCSI(buf[1:])
	case 'O':
		return parseSS3(buf[1:])
	default:
		return 0
	}

	panic("unreachable")
}

func parseKey(buf []byte) int32 {
	if len(buf) < 1 {
		return 0
	}

	switch buf[0] {
	case '\033':
		return parseESC(buf[1:])
	default:
		return int32(buf[0])
	}

	panic("unreachable")
}

