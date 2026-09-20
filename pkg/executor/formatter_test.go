/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package executor_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/stretchr/testify/require"
)

// newTestPrinter returns a Printer wired to an in-memory buffer so
// tests can assert on the rendered bytes without touching stdout.
// termWidth is forced via setting the unexported field through a
// subsequent PrintTaskHeader call, so we drive it indirectly: tests
// that need width control build their own Printer with
// executor.NewPrinter and then exercise PrintTaskHeader's behaviour
// without depending on the tty-detection path.
func newTestPrinter(t *testing.T, verboseLevel int) (*executor.Printer, *bytes.Buffer) {
	t.Helper()
	buf := &bytes.Buffer{}
	return executor.NewPrinter(buf, verboseLevel), buf
}

func TestPrinter_PrintKeyValue_EmptyValueIsSkipped(t *testing.T) {
	p, buf := newTestPrinter(t, 0)
	p.PrintKeyValue("", "Stdout", "")
	require.Empty(t, buf.String(), "empty values must not produce any output")
}

func TestPrinter_PrintKeyValue_SingleLineAtVerboseZero(t *testing.T) {
	p, buf := newTestPrinter(t, 0)
	p.PrintKeyValue("    ", "ExitStatus", "0")

	out := buf.String()
	require.Contains(t, out, "ExitStatus:")
	require.Contains(t, out, "0\n")
	require.NotContains(t, out, "\n    0\n",
		"single-line values at verbose 0 must render inline, not as a multi-line block")
}

func TestPrinter_PrintKeyValue_SingleLineWithTrailingNewline(t *testing.T) {
	// A trailing newline is still a single visual line, so the
	// inline path should still apply.
	p, buf := newTestPrinter(t, 0)
	p.PrintKeyValue("", "Stdout", "hello\n")

	out := buf.String()
	require.Contains(t, out, "Stdout:")
	require.Contains(t, out, "hello")
	// The trailing newline of the value must have been stripped so the
	// output ends with exactly one newline from the surrounding fmt.
	require.False(t, strings.HasSuffix(out, "\n\n"),
		"trailing newline in value should not produce a blank line in output")
}

func TestPrinter_PrintKeyValue_MultiLineAtVerboseZero(t *testing.T) {
	p, buf := newTestPrinter(t, 0)
	p.PrintKeyValue("    ", "Stdout", "line1\nline2\nline3\n")

	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	require.Len(t, lines, 4, "expected 1 header + 3 body lines, got %q", out)
	require.Contains(t, lines[0], "Stdout:")
	require.Contains(t, lines[1], "line1")
	require.Contains(t, lines[2], "line2")
	require.Contains(t, lines[3], "line3")
}

func TestPrinter_PrintKeyValue_MultiLineAtHigherVerbose(t *testing.T) {
	// At verbose > 0, even a single-line value should render as a
	// multi-line block to make structured fields easier to scan.
	p, buf := newTestPrinter(t, 2)
	p.PrintKeyValue("", "ExitStatus", "0")

	out := buf.String()
	require.Contains(t, out, "ExitStatus:\n")
	require.Contains(t, out, "    0\n")
}

func TestPrinter_PrintKeyValue_RespectsPrefix(t *testing.T) {
	p, buf := newTestPrinter(t, 0)
	p.PrintKeyValue(">>", "Stdout", "hello")

	out := buf.String()
	require.True(t, strings.HasPrefix(out, ">>    "),
		"prefix should be followed by the 4-space indent, got %q", out)
}

func TestPrinter_PrintHostResult_NilResStillPrintsHeader(t *testing.T) {
	h, err := remote.ParseHostLine("ops:pass@127.0.0.1:9001")
	require.NoError(t, err)

	p, buf := newTestPrinter(t, 0)
	p.PrintHostResult("", h, nil)

	out := buf.String()
	require.Contains(t, out, "ops@127.0.0.1:9001")
	require.Contains(t, out, "[")
	// No key/value rendering happens when res is nil.
	require.NotContains(t, out, "ExitStatus")
}

func TestPrinter_PrintHostResult_RendersKeyValues(t *testing.T) {
	h, err := remote.ParseHostLine("ops:pass@127.0.0.1:9001")
	require.NoError(t, err)

	res := gtypes.NewOrderedMap[string, string]()
	res.Put("ExitStatus", "0")
	res.Put("Stdout", "hello\n")

	p, buf := newTestPrinter(t, 0)
	p.PrintHostResult("    ", h, res)

	out := buf.String()
	require.Contains(t, out, "ops@127.0.0.1:9001")
	require.Contains(t, out, "ExitStatus:")
	require.Contains(t, out, "0")
	require.Contains(t, out, "Stdout:")
	require.Contains(t, out, "hello")
}

func TestPrinter_PrintTaskSeparator(t *testing.T) {
	p, buf := newTestPrinter(t, 0)
	p.PrintTaskSeparator()
	require.Equal(t, "\n", buf.String())
}

func TestPrinter_PrintTaskHeader_WritesHeaderAndSeparator(t *testing.T) {
	// PrintTaskHeader uses stdout's terminal width when not set; in
	// tests we have no tty, so it falls back to a fixed width. We
	// only assert on the name and timestamp appearing, plus that the
	// line ends with a newline.
	p, buf := newTestPrinter(t, 0)
	p.PrintTaskHeader("my-module")

	out := buf.String()
	require.Contains(t, out, "my-module")
	require.Contains(t, out, "•") // delimiter character
	require.True(t, strings.HasSuffix(out, "\n"))
}
