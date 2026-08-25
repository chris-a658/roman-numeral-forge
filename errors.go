package main

import (
	"fmt"
	"strings"
)

// ParseError describes exactly where and why a conversion failed. Column is a 1-based
// rune index into Input. Line is filled in by the caller once it knows which line of
// the source the error came from; a zero Line means "not applicable" (e.g. a caller
// converting a single string with no surrounding file).
type ParseError struct {
	Line   int
	Column int
	Input  string
	Reason string
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Reason)
	}
	return fmt.Sprintf("column %d: %s", e.Column, e.Reason)
}

// Report renders the error the way a compiler would: the message, the offending line,
// and a caret under the exact column.
func (e *ParseError) Report() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", e.Error())
	if e.Input != "" {
		fmt.Fprintf(&b, "    %s\n", e.Input)
		fmt.Fprintf(&b, "    %s^\n", strings.Repeat(" ", e.Column-1))
	}
	return b.String()
}
