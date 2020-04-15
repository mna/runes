// Command runes prints Unicode runes to stdout along with UTF8
// encoding information. Specific code points and ranges of code
// points can be provided as arguments.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"golang.org/x/text/unicode/runenames"
)

func main() {
	var (
		flagHelp  = flag.Bool("h", false, "Display this message.")
		flagLHelp = flag.Bool("help", false, "Display this message.")
		flagJSON  = flag.Bool("json", false, "Output JSON data.")
	)
	flag.Usage = usage
	flag.Parse()
	_ = flagJSON

	if *flagHelp || *flagLHelp {
		help()
		return
	}

	var p printer
	if *flagJSON {
		p = &jsonPrinter{}
	} else {
		p = &textPrinter{}
	}

	start, end := rune(12000), rune(13000)
	//args := flag.Args()
	if err := p.printStart(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	printRange(p, start, end)
	if err := p.printEnd(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type printer interface {
	printStart(w io.Writer) error
	printRune(runeInfo) error
	printEnd() error
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
	jp.ris = make([]runeInfo, 0, 1024) // TODO: receive a size hint
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

type runeInfo struct {
	Rune       rune     `json:"rune"`
	Name       string   `json:"name"`
	Valid      bool     `json:"valid"`
	Categories []string `json:"categories"`
	UTF8       []byte   `json:"-"`
	UTF16      []uint16 `json:"utf16"`
	UTF8JSON   []uint16 `json:"utf8"`
}

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

func printRune(p printer, r rune) {
	ri := info(r)
	p.printRune(ri)
}

func printRange(p printer, start, end rune) {
	for i := start; i <= end; i++ {
		// for ranges, ignore invalid utf8 runes
		if !utf8.ValidRune(i) {
			continue
		}
		printRune(p, i)
	}
}

func usage() {
	const usageMsg = `usage: runes [-h] [-json] [CODEPOINT ...] [CPSTART-CPEND ...]`
	fmt.Println(usageMsg)
}

func help() {
	const msg = `
The runes command prints information about Unicode code points. Without argument,
all code points are printed; specific code points can be requested as arguments,
and ranges of code points are supported. Code points starting with '0x' or 'u+'
are considered in hexadecimal (the 'x' and 'u' are case insensitive), otherwise
the number is processed as decimal.

Examples:
    runes
    runes 0x2318 40-60
    runes u+1f970
`
	usage()
	fmt.Println(msg)
}
