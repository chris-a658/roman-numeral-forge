package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunSummaryAndExitCodes(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		stdin      string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "all lines convert",
			stdin:      "1994\nMCMXCIV\n",
			wantCode:   0,
			wantStdout: "MCMXCIV\n1994\n",
			wantStderr: "romanforge: 2 converted, 0 failed\n",
		},
		{
			name:     "one bad line keeps going and exits nonzero",
			stdin:    "1994\nIIII\n",
			wantCode: 1,
			// MCMXCIV converts fine, IIII fails but doesn't stop the run.
			wantStdout: "MCMXCIV\n",
		},
		{
			name:     "strict mode stops at the first error",
			stdin:    "IIII\n1994\n",
			args:     []string{"-strict"},
			wantCode: 1,
			// 1994 is never reached.
			wantStdout: "",
		},
		{
			name:       "blank lines and comments produce no summary",
			stdin:      "\n# a comment\n   \n",
			wantCode:   0,
			wantStdout: "",
			wantStderr: "",
		},
		{
			name:     "conflicting direction flags",
			args:     []string{"-to-roman", "-to-arabic"},
			stdin:    "",
			wantCode: 1,
			wantStderr: "romanforge: -to-roman and -to-arabic are mutually exclusive\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(c.args, strings.NewReader(c.stdin), &stdout, &stderr)
			if code != c.wantCode {
				t.Errorf("run() exit code = %d, want %d (stderr: %s)", code, c.wantCode, stderr.String())
			}
			if stdout.String() != c.wantStdout {
				t.Errorf("run() stdout = %q, want %q", stdout.String(), c.wantStdout)
			}
			if c.wantStderr != "" && stderr.String() != c.wantStderr {
				t.Errorf("run() stderr = %q, want %q", stderr.String(), c.wantStderr)
			}
		})
	}
}

func TestRunToRomanFlagPinsDirection(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-to-roman"}, strings.NewReader("9\n"), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() exit code = %d, want 0 (stderr: %s)", code, stderr.String())
	}
	if got, want := stdout.String(), "IX\n"; got != want {
		t.Errorf("run() stdout = %q, want %q", got, want)
	}
}

func TestRunUnknownFlagReturnsUsageExitCode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-not-a-real-flag"}, strings.NewReader(""), &stdout, &stderr)
	if code != 2 {
		t.Errorf("run() exit code = %d, want 2", code)
	}
}
