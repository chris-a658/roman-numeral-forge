package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var scanner *bufio.Scanner
	source := "stdin"

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "romanforge:", err)
			os.Exit(1)
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
		source = os.Args[1]
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

		out, err := convertLine(line)
		if err != nil {
			failed = true
			if pe, ok := err.(*ParseError); ok {
				pe.Line = lineNum
				if pe.Input == "" {
					pe.Input = raw
				}
				fmt.Fprint(os.Stderr, pe.Report())
			} else {
				fmt.Fprintf(os.Stderr, "%s:%d: %v\n", source, lineNum, err)
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
	if n, err := strconv.Atoi(line); err == nil {
		return ToRoman(n)
	}
	n, err := FromRoman(strings.ToUpper(line))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(n), nil
}
