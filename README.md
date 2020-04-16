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
```

## Example output

```
$ ./runes u+1f970 0x12 u+2134-u+213a - Martin
[S So] U+1F970 '🥰'   [F0 9F A5 B0] [D83E DD70] SMILING FACE WITH SMILING EYES AND THREE HEARTS
[C Cc] U+0012         [12]          [12]        <control>
[L Ll] U+2134 'ℴ'     [E2 84 B4]    [2134]      SCRIPT SMALL O
[L Lo] U+2135 'ℵ'     [E2 84 B5]    [2135]      ALEF SYMBOL
[L Lo] U+2136 'ℶ'     [E2 84 B6]    [2136]      BET SYMBOL
[L Lo] U+2137 'ℷ'     [E2 84 B7]    [2137]      GIMEL SYMBOL
[L Lo] U+2138 'ℸ'     [E2 84 B8]    [2138]      DALET SYMBOL
[L Ll] U+2139 'ℹ'     [E2 84 B9]    [2139]      INFORMATION SOURCE
[S So] U+213A '℺'     [E2 84 BA]    [213A]      ROTATED CAPITAL Q
[L Lu] U+004D 'M'     [4D]          [4D]        LATIN CAPITAL LETTER M
[L Ll] U+0061 'a'     [61]          [61]        LATIN SMALL LETTER A
[L Ll] U+0072 'r'     [72]          [72]        LATIN SMALL LETTER R
[L Ll] U+0074 't'     [74]          [74]        LATIN SMALL LETTER T
[L Ll] U+0069 'i'     [69]          [69]        LATIN SMALL LETTER I
[L Ll] U+006E 'n'     [6E]          [6E]        LATIN SMALL LETTER N
```

## License

The [BSD 3-Clause license][bsd].

[bsd]: http://opensource.org/licenses/BSD-3-Clause
