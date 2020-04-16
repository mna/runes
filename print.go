package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"golang.org/x/text/unicode/runenames"
)

type printer interface {
	printStart(w io.Writer) error
	printRune(runeInfo) error
	printEnd() error
}

type runeInfo struct {
	Rune       rune     `json:"rune"`
	Name       string   `json:"name"`
	Valid      bool     `json:"valid"`
	Categories []string `json:"categories"`
	UTF8       []byte   `json:"-"`
	UTF16      []uint16 `json:"utf16"`
	UTF8JSON   []uint16 `json:"utf8"`
}

// return a filled runeInfo struct for that rune r.
func info(r rune) runeInfo {
	var buf [utf8.UTFMax]byte

	ri := runeInfo{Rune: r}
	ri.Valid = utf8.ValidRune(r)
	if !ri.Valid {
		return ri
	}

	ri.Name = runenames.Name(r)

	n := utf8.EncodeRune(buf[:], r)
	ri.UTF8 = buf[:n]
	r1, r2 := utf16.EncodeRune(r)
	if r1 == utf8.RuneError && r2 == utf8.RuneError {
		ri.UTF16 = []uint16{uint16(r)}
	} else {
		ri.UTF16 = []uint16{uint16(r1), uint16(r2)}
	}

	var cats []string
	for nm, rt := range unicode.Categories {
		if unicode.Is(rt, r) {
			cats = append(cats, nm)
		}
	}
	sort.Strings(cats)
	ri.Categories = cats
	return ri
}

// print a single rune.
func printRune(p printer, r rune) error {
	ri := info(r)
	return p.printRune(ri)
}

// print an explicit list of runes.
func printRunes(p printer, rs []rune) error {
	for _, r := range rs {
		if err := printRune(p, r); err != nil {
			return err
		}
	}
	return nil
}

type textPrinter struct {
	bw *bufio.Writer
}

func (tp *textPrinter) printStart(w io.Writer) error {
	tp.bw = bufio.NewWriter(w)
	return nil
}

func (tp *textPrinter) printRune(ri runeInfo) error {
	var catgs string
	if ri.Valid {
		catgs = fmt.Sprintf("%v", ri.Categories)
	} else {
		catgs = "[!]"
	}
	fmt.Fprintf(tp.bw, "%-7s", catgs)

	wd := runewidth.RuneWidth(ri.Rune)
	rn := fmt.Sprintf("%#U", ri.Rune)
	if n := len(rn) + wd - 1; n < 15 {
		rn += strings.Repeat(" ", 15-n)
	}
	fmt.Fprintf(tp.bw, "%s", rn)

	u8 := fmt.Sprintf("[% X]", ri.UTF8)
	fmt.Fprintf(tp.bw, "%-12s", u8)

	var u16 string
	if len(ri.UTF16) == 2 {
		u16 = fmt.Sprintf("[%X %X]", ri.UTF16[0], ri.UTF16[1])
	} else {
		u16 = fmt.Sprintf("[%X]", ri.UTF16[0])
	}
	fmt.Fprintf(tp.bw, "%-12s", u16)

	fmt.Fprintln(tp.bw, ri.Name)
	return nil
}

func (tp *textPrinter) printEnd() error {
	return tp.bw.Flush()
}

type jsonPrinter struct {
	ris []runeInfo
	bw  *bufio.Writer
}

func (jp *jsonPrinter) printStart(w io.Writer) error {
	jp.bw = bufio.NewWriter(w)
	jp.ris = make([]runeInfo, 0, 1024) // TODO: receive a size hint?
	return nil
}

func (jp *jsonPrinter) printRune(ri runeInfo) error {
	ri.UTF8JSON = make([]uint16, len(ri.UTF8))
	for i, b := range ri.UTF8 {
		ri.UTF8JSON[i] = uint16(b)
	}
	jp.ris = append(jp.ris, ri)
	return nil
}

func (jp *jsonPrinter) printEnd() error {
	enc := json.NewEncoder(jp.bw)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jp.ris); err != nil {
		return err
	}
	return jp.bw.Flush()
}
