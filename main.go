// Command runes prints Unicode runes to stdout along with UTF8
// encoding information. Specific code points and ranges of code
// points can be provided as arguments.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	var (
		flagHelp  = flag.Bool("h", false, "Display this message.")
		flagLHelp = flag.Bool("help", false, "Display this message.")
		flagJSON  = flag.Bool("json", false, "Output JSON data.")
		flagAll   = flag.Bool("all", false, "Print all Unicode code points.")
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

	if !*flagAll {
		lastIx := -1
		for i, arg := range args {
			if len(arg) == 0 {
				continue
			}
			if arg == "-" {
				lastIx = i
				break
			}

			var nums []int
			parts := strings.Split(arg, "-")
			if len(parts) > 2 {
				fmt.Fprintf(os.Stderr, "invalid rune argument: too many parts in range %s\n", arg)
				os.Exit(1)
			}
			for _, part := range parts {
				if strings.HasPrefix(part, "u+") || strings.HasPrefix(part, "U+") {
					part = "0x" + part[2:]
				}
				n, err := strconv.ParseUint(part, 0, 32)
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid rune argument: %s\n", arg)
					os.Exit(1)
				}
				nums = append(nums, int(n))
			}

			if len(nums) == 1 {
				rs = append(rs, rune(nums[0]))
				continue
			}
			rs = append(rs, runesInRange(nums[0], nums[1])...)
		}

		// if there are remaining arguments, treat them as strings to print the
		// runes of.
		args = args[lastIx+1:]
		for _, arg := range args {
			rs = append(rs, runesSet(arg)...)
		}
	}

	count := len(rs)
	if *flagAll {
		count = int(unicode.MaxRune)
	}

	if err := p.printStart(os.Stdout, count); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *flagAll {
		for i := 0; i <= unicode.MaxRune; i++ {
			if utf8.ValidRune(rune(i)) {
				if err := printRune(p, rune(i)); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}
	} else {
		if err := printRunes(p, rs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if err := p.printEnd(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runesInRange(start, end int) []rune {
	rs := make([]rune, 0, end-start+1)
	for i := start; i <= end; i++ {
		rs = append(rs, rune(i))
	}
	return rs
}

// returns a slice of runes where each distinct rune in arg is returned.
func runesSet(arg string) []rune {
	rs := make([]rune, len(arg))
	for i, r := range arg {
		rs[i] = r
	}
	return rs
}

func usage() {
	const usageMsg = `usage: runes [-h] [-json] [CODEPT ...] [START-END ...] [- STRING ...]`
	fmt.Println(usageMsg)
}

func help() {
	const msg = `
The runes command prints information about Unicode code points. Specific code
points can be requested as arguments, and ranges of code points are supported
(e.g. 0x17-0x60). Code points starting with '0x' or 'u+' are considered in
hexadecimal (the 'x' and 'u' are case insensitive), otherwise the number is
treated as decimal.

A single dash argument '-' can be used so that subsequent arguments are treated
as strings for which each rune will be printed.

The output follows the order of runes as specified on the command-line,
the same rune will be printed multiple times if it is specified or included
in multiple arguments.

Flags:
  -h,-help           Display this message.
  -json              Output JSON data.
  -all               Print all Unicode code points.

Examples:
    runes -all
    runes -json 0x2318 40-60
    runes u+1f970 0X55-0XA0 - "Some string"
`
	usage()
	fmt.Println(msg)
}
