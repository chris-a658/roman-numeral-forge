package main

import (
	"testing"
)

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
