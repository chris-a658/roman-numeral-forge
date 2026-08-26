# roman-numeral-forge

A command-line converter between Arabic numbers and Roman numerals. Feed it a file or
stdin with one value per line - numbers get converted to numerals, numerals get
converted back to numbers, and it figures out which is which per line.

## Why

Most roman numeral converters either silently accept garbage like `IIII` or `VX`, or
fail with a generic "invalid input" that doesn't say where the problem is. This one
validates strictly - a numeral is only accepted if it's the exact canonical form for
some integer - and when something is wrong it reports the line and column, with a
caret pointing at the offending character:

```
$ echo "MCMXA" | romanforge
line 1, column 5: unexpected character 'A' (valid roman numeral digits are I, V, X, L, C, D, M)
    MCMXA
        ^
```

## Usage

Build:

```
go build -o romanforge .
```

Convert a mix of numbers and numerals, one per line:

```
$ printf "1994\nXL\n2026\nIIII\n" | ./romanforge
MCMXCIV
40
MMXXVI
line 4, column 2: "IIII" is not a valid roman numeral (did you mean "IV"?)
    IIII
     ^
```

Read from a file instead of stdin:

```
./romanforge input.txt
```

Lines starting with `#` are treated as comments and skipped; blank lines are skipped
too. A line that fails to convert is reported to stderr and the rest of the file still
runs, so one bad line doesn't hide problems further down.

## Rules

- Arabic input must be a plain integer from 1 to 3999 (roman numerals have no zero and
  no standard notation past 3999).
- Roman input must be in canonical form: at most three repeated I/X/C/M, no repeated
  V/L/D, and only the six standard subtractive pairs (IV, IX, XL, XC, CD, CM). Anything
  that would round-trip to a different string is rejected, with the canonical form
  suggested in the error.

## Status

Early skeleton. Conversion and error reporting work end to end, with unit tests
covering the round trip for every value in range and the error columns for the
non-canonical cases; there are no explicit direction flags yet.

## License

MIT, see LICENSE.
