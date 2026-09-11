package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"romanforge/internal/roman"
)

func main() {
	toRoman := flag.Bool("to-roman", false, "treat every line as an integer and convert it to a roman numeral")
	toArabic := flag.Bool("to-arabic", false, "treat every line as a roman numeral and convert it to an integer")
	strict := flag.Bool("strict", false, "exit immediately on the first conversion error instead of reporting it and continuing")
	flag.Parse()

	if *toRoman && *toArabic {
		fmt.Fprintln(os.Stderr, "romanforge: -to-roman and -to-arabic are mutually exclusive")
		os.Exit(1)
	}
	convert := convertLine
	switch {
	case *toRoman:
		convert = convertToRoman
	case *toArabic:
		convert = convertToArabic
	}

	var scanner *bufio.Scanner
	source := "stdin"

	if args := flag.Args(); len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "romanforge:", err)
			os.Exit(1)
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
		source = args[0]
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	failed := false
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		out, err := convert(line)
		if err != nil {
			failed = true
			if pe, ok := err.(*roman.ParseError); ok {
				pe.Line = lineNum
				if pe.Input == "" {
					pe.Input = raw
				}
				fmt.Fprint(os.Stderr, pe.Report())
			} else {
				fmt.Fprintf(os.Stderr, "%s:%d: %v\n", source, lineNum, err)
			}
			if *strict {
				os.Exit(1)
			}
			continue
		}
		fmt.Println(out)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "romanforge:", err)
		os.Exit(1)
	}
	if failed {
		os.Exit(1)
	}
}

// convertLine converts a single line in either direction: a plain integer becomes a
// roman numeral, anything else is parsed as a roman numeral and turned back into a
// number.
func convertLine(line string) (string, error) {
	if _, err := strconv.Atoi(line); err == nil {
		return convertToRoman(line)
	}
	return convertToArabic(line)
}

// convertToRoman parses line as a plain integer and converts it to a roman numeral.
// Used when -to-roman pins the direction instead of letting convertLine guess it.
func convertToRoman(line string) (string, error) {
	n, err := strconv.Atoi(line)
	if err != nil {
		return "", &roman.ParseError{Column: 1, Input: line, Reason: fmt.Sprintf("%q is not a valid integer", line)}
	}
	return roman.ToRoman(n)
}

// convertToArabic parses line as a roman numeral and converts it to an integer.
// Used when -to-arabic pins the direction instead of letting convertLine guess it.
func convertToArabic(line string) (string, error) {
	n, err := roman.FromRoman(strings.ToUpper(line))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(n), nil
}
