package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"romanforge/internal/roman"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run holds the entire CLI: flag parsing, the convert loop, and the summary line.
// It takes its args and I/O as parameters, rather than reaching for os.Args and the
// real stdin/stdout, so tests can drive it without spawning a subprocess.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("romanforge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	toRoman := fs.Bool("to-roman", false, "treat every line as an integer and convert it to a roman numeral")
	toArabic := fs.Bool("to-arabic", false, "treat every line as a roman numeral and convert it to an integer")
	strict := fs.Bool("strict", false, "exit immediately on the first conversion error instead of reporting it and continuing")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *toRoman && *toArabic {
		fmt.Fprintln(stderr, "romanforge: -to-roman and -to-arabic are mutually exclusive")
		return 1
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

	if fileArgs := fs.Args(); len(fileArgs) > 0 {
		f, err := os.Open(fileArgs[0])
		if err != nil {
			fmt.Fprintln(stderr, "romanforge:", err)
			return 1
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
		source = fileArgs[0]
	} else {
		scanner = bufio.NewScanner(stdin)
	}

	lineNum := 0
	converted := 0
	failedCount := 0
	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		out, err := convert(line)
		if err != nil {
			failedCount++
			if pe, ok := err.(*roman.ParseError); ok {
				pe.Line = lineNum
				if pe.Input == "" {
					pe.Input = raw
				}
				fmt.Fprint(stderr, pe.Report())
			} else {
				fmt.Fprintf(stderr, "%s:%d: %v\n", source, lineNum, err)
			}
			if *strict {
				return 1
			}
			continue
		}
		converted++
		fmt.Fprintln(stdout, out)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(stderr, "romanforge:", err)
		return 1
	}

	if converted+failedCount > 0 {
		fmt.Fprintf(stderr, "romanforge: %d converted, %d failed\n", converted, failedCount)
	}
	if failedCount > 0 {
		return 1
	}
	return 0
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
