package main

import (
	"fmt"
	"strings"
	"unicode"
)

// table encodes the allowed Unicode code points on 17_408 uint64s (for 1_114_112 bits).
type table [17408]uint64

// set the runes.
func (t *table) set(rs ...rune) {
	for _, r := range rs {
		if r > unicode.MaxRune {
			panic(fmt.Sprintf("%#U is outside the Unicode range", r))
		}
		t[r/64] |= 1 << uint64(r%64)
	}
}

// unset the runes.
func (t *table) unset(rs ...rune) {
	for _, r := range rs {
		if r > unicode.MaxRune {
			panic(fmt.Sprintf("%#U is outside the Unicode range", r))
		}
		t[r/64] &^= 1 << uint64(r%64)
	}
}

// setRange sets all runes in [from, to] inclusively.
func (t *table) setRange(from, to rune) {
	t.setUnsetRange(from, to, true)
}

// unsetRange unsets all runes in [from, to] inclusively.
func (t *table) unsetRange(from, to rune) {
	t.setUnsetRange(from, to, false)
}

func (t *table) setUnsetRange(from, to rune, set bool) {
	if from > to {
		panic(fmt.Sprintf("from rune %#U is greater than to rune %#U", from, to))
	}
	if to > unicode.MaxRune {
		panic(fmt.Sprintf("%#U is outside the Unicode range", to))
	}
	rng := make([]rune, to-from+1)
	for i := from; i <= to; i++ {
		rng[i-from] = i
	}
	if set {
		t.set(rng...)
	} else {
		t.unset(rng...)
	}
}

// is returns true if r is set.
func (t *table) is(r rune) bool {
	if r > unicode.MaxRune {
		return false
	}
	return t[r/64]&(1<<uint64(r%64)) != 0
}

// String returns the string representation of the Unicode table.
func (t *table) String() string {
	var buf strings.Builder
	buf.WriteByte('[')

	var last rune = -1
	writeFromLastTo := func(end rune) {
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}
		if last == end {
			fmt.Fprintf(&buf, "%#U", last)
		} else {
			fmt.Fprintf(&buf, "%#U-%#U", last, end)
		}
	}

	for i := rune(0); i <= unicode.MaxRune; i++ {
		if t.is(i) {
			if last == -1 {
				last = i
			}
			continue
		}
		if last == -1 {
			continue
		}
		writeFromLastTo(i - 1)
		last = -1
	}
	if last != -1 {
		writeFromLastTo(unicode.MaxRune)
	}

	buf.WriteByte(']')
	return buf.String()
}
