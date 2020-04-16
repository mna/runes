// Command runes prints Unicode runes to stdout along with UTF8
// encoding information. Specific code points and ranges of code
// points can be provided as arguments.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var (
		flagHelp  = flag.Bool("h", false, "Display this message.")
		flagLHelp = flag.Bool("help", false, "Display this message.")
		flagJSON  = flag.Bool("json", false, "Output JSON data.")
	)
	flag.Usage = usage
	flag.Parse()

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

	args := flag.Args()
	var rs []rune
	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}

		switch p0 := arg[0]; p0 {
		case 'u', 'U', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// it is either a number or a range
			parts := []string{arg}
			if rangeIx := strings.Index(arg, "-"); rangeIx >= 0 {
				parts = []string{arg[:rangeIx], arg[rangeIx+1:]}
			}

		default:
			rs = append(rs, runesSet(arg)...)
		}
	}

	if err := p.printStart(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := p.printEnd(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// decode a command-line argument into a list of runes to print,
// and return true as second value if this is a range in the form
// <start>-<end> (inclusive).
func decode(arg string) (runes []rune, isRange bool, err error) {
	if len(arg) == 0 {
		return nil, false, nil
	}

	p0 := arg[0]
	base := 10
	start := 0
	switch p0 {
	case 'u', 'U':
		if len(arg) == 1 || arg[1] != '+' {
			return runesSet(arg), false, nil
		}
		base = 16
		start = 2 // skip u+
	case '0':
		if len(arg) > 1 && arg[1] == 'x' || arg[1] == 'X' {
			base = 16
			start = 2 // skip 0x
			break
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// ok, decimal number
	default:
		return runesSet(arg), false, nil
	}

	num, err := strconv.ParseUint(arg[start:], base, 32)
	if err != nil {
		return nil, false, err
	}
	runes = append(runes, rune(num))
	return runes, false, nil
}

// returns a slice of runes where each distinct rune in arg is returned.
func runesSet(arg string) []rune {
	m := make(map[rune]bool)
	for _, r := range arg {
		m[r] = true
	}
	rs := make([]rune, 0, len(m))
	for k := range m {
		rs = append(rs, k)
	}
	sort.Slice(rs, func(l, r int) bool {
		lr, rr := rs[l], rs[r]
		return lr < rr
	})
	return rs
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
