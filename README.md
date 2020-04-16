# runes

Command runes prints Unicode runes to stdout along with the rune's
category and sub-category, its UTF8 and UTF16 encoding information
and its official name. Specific code points and ranges of code
points can be provided as arguments.

* Canonical repository: https://git.sr.ht/~mna/runes
* Issues: https://todo.sr.ht/~mna/runes

## Usage

```
usage: runes [-h] [-json] [CODEPT ...] [START-END ...] [- STRING ...]

The runes command prints information about Unicode code points. Without
argument, all code points are printed; specific code points can be requested as
arguments, and ranges of code points are supported (e.g. 0x17-0x60). Code
points starting with '0x' or 'u+' are considered in hexadecimal (the 'x' and
'u' are case insensitive), otherwise the number is processed as decimal.

A single dash '-' can be used so that subsequent arguments are treated as
strings for which each rune will be printed.

The output follows the order of runes as specified on the command-line,
the same rune will be printed multiple times if it is specified or included
in multiple arguments.

Examples:
    runes
    runes 0x2318 40-60
    runes u+1f970 0X55-0XA0 - "Some string"
```

## License

The [BSD 3-Clause license][bsd].

[bsd]: http://opensource.org/licenses/BSD-3-Clause
