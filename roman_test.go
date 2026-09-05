package main

import (
	"strings"
	"testing"
)

func TestToRoman(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{1, "I"},
		{4, "IV"},
		{9, "IX"},
		{14, "XIV"},
		{40, "XL"},
		{90, "XC"},
		{400, "CD"},
		{900, "CM"},
		{1994, "MCMXCIV"},
		{2026, "MMXXVI"},
		{3999, "MMMCMXCIX"},
	}
	for _, c := range cases {
		got, err := ToRoman(c.n)
		if err != nil {
			t.Errorf("ToRoman(%d) returned error: %v", c.n, err)
			continue
		}
		if got != c.want {
			t.Errorf("ToRoman(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestToRomanOutOfRange(t *testing.T) {
	for _, n := range []int{0, -1, 4000, 10000} {
		_, err := ToRoman(n)
		if err == nil {
			t.Errorf("ToRoman(%d) should have errored, got nil", n)
			continue
		}
		pe, ok := err.(*ParseError)
		if !ok {
			t.Errorf("ToRoman(%d) error is not a *ParseError: %T", n, err)
			continue
		}
		if pe.Column != 1 {
			t.Errorf("ToRoman(%d) error column = %d, want 1", n, pe.Column)
		}
	}
}

func TestFromRoman(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"I", 1},
		{"IV", 4},
		{"IX", 9},
		{"XIV", 14},
		{"XL", 40},
		{"XC", 90},
		{"CD", 400},
		{"CM", 900},
		{"MCMXCIV", 1994},
		{"MMXXVI", 2026},
		{"MMMCMXCIX", 3999},
	}
	for _, c := range cases {
		got, err := FromRoman(c.s)
		if err != nil {
			t.Errorf("FromRoman(%q) returned error: %v", c.s, err)
			continue
		}
		if got != c.want {
			t.Errorf("FromRoman(%q) = %d, want %d", c.s, got, c.want)
		}
	}
}

func TestFromRomanRoundTrip(t *testing.T) {
	for n := 1; n <= 3999; n++ {
		s, err := ToRoman(n)
		if err != nil {
			t.Fatalf("ToRoman(%d) returned error: %v", n, err)
		}
		got, err := FromRoman(s)
		if err != nil {
			t.Fatalf("FromRoman(%q) returned error: %v", s, err)
		}
		if got != n {
			t.Fatalf("round trip broke: ToRoman(%d) = %q, FromRoman(%q) = %d", n, s, s, got)
		}
	}
}

func TestFromRomanEmpty(t *testing.T) {
	_, err := FromRoman("")
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("FromRoman(\"\") error is not a *ParseError: %v", err)
	}
	if pe.Column != 1 {
		t.Errorf("FromRoman(\"\") error column = %d, want 1", pe.Column)
	}
}

func TestFromRomanInvalidCharacter(t *testing.T) {
	// mirrors the README example: the bad character is the 5th rune.
	_, err := FromRoman("MCMXA")
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("FromRoman(\"MCMXA\") error is not a *ParseError: %v", err)
	}
	if pe.Column != 5 {
		t.Errorf("FromRoman(\"MCMXA\") error column = %d, want 5", pe.Column)
	}
	if !strings.Contains(pe.Reason, `'A'`) {
		t.Errorf("FromRoman(\"MCMXA\") reason = %q, want it to mention 'A'", pe.Reason)
	}
}

func TestFromRomanNonCanonical(t *testing.T) {
	cases := []struct {
		s          string
		wantColumn int
		wantHint   string
	}{
		{"IIII", 2, "IV"},   // too many repeats: first divergence is the 2nd char
		{"VX", 2, "V"},      // invalid subtractive pair, evaluates to 5
		{"IL", 1, "XLIX"},   // invalid subtractive pair, evaluates to 49
		{"XCIX", 0, ""},     // sanity check this one is actually valid, see below
	}
	for _, c := range cases {
		if c.s == "XCIX" {
			if _, err := FromRoman(c.s); err != nil {
				t.Errorf("FromRoman(%q) should be valid, got error: %v", c.s, err)
			}
			continue
		}
		_, err := FromRoman(c.s)
		pe, ok := err.(*ParseError)
		if !ok {
			t.Errorf("FromRoman(%q) error is not a *ParseError: %v", c.s, err)
			continue
		}
		if pe.Column != c.wantColumn {
			t.Errorf("FromRoman(%q) error column = %d, want %d", c.s, pe.Column, c.wantColumn)
		}
		if !strings.Contains(pe.Reason, c.wantHint) {
			t.Errorf("FromRoman(%q) reason = %q, want it to mention %q", c.s, pe.Reason, c.wantHint)
		}
	}
}

func TestFromRomanOutOfRangeTotal(t *testing.T) {
	_, err := FromRoman("MMMM")
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("FromRoman(\"MMMM\") error is not a *ParseError: %v", err)
	}
	if !strings.Contains(pe.Reason, "4000") {
		t.Errorf("FromRoman(\"MMMM\") reason = %q, want it to mention 4000", pe.Reason)
	}
}

func TestConvertLine(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"1994", "MCMXCIV"},
		{"MCMXCIV", "1994"},
		{"mcmxciv", "1994"}, // lowercase numerals are accepted
	}
	for _, c := range cases {
		got, err := convertLine(c.in)
		if err != nil {
			t.Errorf("convertLine(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("convertLine(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestConvertLineErrors(t *testing.T) {
	for _, in := range []string{"0", "4000", "IIII", "", "VX"} {
		if _, err := convertLine(in); err == nil {
			t.Errorf("convertLine(%q) should have errored", in)
		}
	}
}

func TestConvertToRoman(t *testing.T) {
	got, err := convertToRoman("1994")
	if err != nil {
		t.Fatalf("convertToRoman(\"1994\") returned error: %v", err)
	}
	if got != "MCMXCIV" {
		t.Errorf("convertToRoman(\"1994\") = %q, want %q", got, "MCMXCIV")
	}

	if _, err := convertToRoman("MCMXCIV"); err == nil {
		t.Error("convertToRoman(\"MCMXCIV\") should have errored: it is not an integer")
	}
}

func TestConvertToArabic(t *testing.T) {
	got, err := convertToArabic("mcmxciv")
	if err != nil {
		t.Fatalf("convertToArabic(\"mcmxciv\") returned error: %v", err)
	}
	if got != "1994" {
		t.Errorf("convertToArabic(\"mcmxciv\") = %q, want %q", got, "1994")
	}

	if _, err := convertToArabic("1994"); err == nil {
		t.Error("convertToArabic(\"1994\") should have errored: it is not a roman numeral")
	}
}

func TestParseErrorReport(t *testing.T) {
	pe := &ParseError{Line: 4, Column: 2, Input: "IIII", Reason: `"IIII" is not a valid roman numeral (did you mean "IV"?)`}
	report := pe.Report()
	wantLines := []string{
		`line 4, column 2: "IIII" is not a valid roman numeral (did you mean "IV"?)`,
		"    IIII",
		"     ^",
	}
	for _, want := range wantLines {
		if !strings.Contains(report, want) {
			t.Errorf("Report() = %q, want it to contain %q", report, want)
		}
	}
}

func TestParseErrorWithoutLine(t *testing.T) {
	pe := &ParseError{Column: 3, Reason: "boom"}
	got := pe.Error()
	want := "column 3: boom"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
