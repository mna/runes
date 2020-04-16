package main

import (
	"testing"
	"unicode"
)

func TestTable(t *testing.T) {
	var tbl table

	tbl.set([]rune("abcd")...)
	got := tbl.String()
	want := "[U+0061 'a'-U+0064 'd']"
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}

	tbl.setRange('A', 'Z')
	tbl.unsetRange('M', 'Q')
	got = tbl.String()
	want = "[U+0041 'A'-U+004C 'L',U+0052 'R'-U+005A 'Z',U+0061 'a'-U+0064 'd']"
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}

	tbl.setRange('╒', '╟')
	tbl.unsetRange('A', 'z')
	got = tbl.String()
	want = "[U+2552 '╒'-U+255F '╟']"
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}

	tbl.set(unicode.MaxRune)
	got = tbl.String()
	want = "[U+2552 '╒'-U+255F '╟',U+10FFFF]"
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}
