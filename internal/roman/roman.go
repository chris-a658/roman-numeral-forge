// Package roman converts between integers and canonical roman numerals.
package roman

import (
	"fmt"
	"strings"
)

var romanTable = []struct {
	Value  int
	Symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

var romanValues = map[rune]int{
	'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000,
}

// ToRoman converts an integer in [1, 3999] to its canonical roman numeral form.
func ToRoman(n int) (string, error) {
	if n < 1 || n > 3999 {
		return "", &ParseError{
			Column: 1,
			Reason: fmt.Sprintf("%d is out of range: roman numerals represent integers from 1 to 3999", n),
		}
	}
	var b strings.Builder
	remaining := n
	for _, digit := range romanTable {
		for remaining >= digit.Value {
			b.WriteString(digit.Symbol)
			remaining -= digit.Value
		}
	}
	return b.String(), nil
}

// FromRoman parses a roman numeral and returns its integer value. It only accepts
// canonical form - the same string ToRoman would produce for that value - so things
// like "IIII" or "VX" are rejected instead of guessed at. This also means the
// implementation doesn't need a separate hand-rolled grammar for repetition and
// subtractive-pair rules: a string is valid exactly when it round-trips.
func FromRoman(s string) (int, error) {
	if s == "" {
		return 0, &ParseError{Column: 1, Reason: "empty input"}
	}

	runes := []rune(s)
	for i, r := range runes {
		if _, ok := romanValues[r]; !ok {
			return 0, &ParseError{
				Column: i + 1,
				Input:  s,
				Reason: fmt.Sprintf("unexpected character %q (valid roman numeral digits are I, V, X, L, C, D, M)", r),
			}
		}
	}

	total := 0
	for i, r := range runes {
		v := romanValues[r]
		if i+1 < len(runes) && v < romanValues[runes[i+1]] {
			total -= v
		} else {
			total += v
		}
	}

	if total < 1 || total > 3999 {
		return 0, &ParseError{
			Column: 1,
			Input:  s,
			Reason: fmt.Sprintf("%q evaluates to %d, which is out of range: roman numerals represent 1 to 3999", s, total),
		}
	}

	canonical, err := ToRoman(total)
	if err != nil {
		return 0, err
	}
	if canonical != s {
		return 0, &ParseError{
			Column: firstDiff(s, canonical),
			Input:  s,
			Reason: fmt.Sprintf("%q is not a valid roman numeral (did you mean %q?)", s, canonical),
		}
	}

	return total, nil
}

// firstDiff returns the 1-based rune index of the first position where a and b differ.
func firstDiff(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	n := len(ra)
	if len(rb) < n {
		n = len(rb)
	}
	for i := 0; i < n; i++ {
		if ra[i] != rb[i] {
			return i + 1
		}
	}
	return n + 1
}
