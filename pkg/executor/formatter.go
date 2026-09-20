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

package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/theme"
	"golang.org/x/term"
)

const (
	// timestampWidth matches the literal width of the "2006-01-02 15:04:05"
	// format used for task-header timestamps. Keep this in sync with
	// the format string in PrintTaskHeader.
	timestampWidth = 19
	// fallbackTerminalWidth is used when stdout is not a tty (e.g.
	// output is piped to a file or another process). 500 is wide
	// enough to keep task headers on one line for most terminal
	// widths while still being sensible for non-tty output.
	fallbackTerminalWidth = 500
)

// Printer renders task, host, and key/value output for an Executor run.
// It owns the output writer, the detected terminal width, and the
// verbosity level that controls single-line vs multi-line value
// formatting, so callers (Executor, tests, alternative front-ends)
// don't need to wire these concerns themselves.
type Printer struct {
	out          io.Writer
	termWidth    int
	verboseLevel int
}

// NewPrinter builds a Printer that writes to out at the given verbosity
// level. The terminal width is detected lazily on the first call to
// PrintTaskHeader so non-tty writers (tests, file redirection) get the
// fallback width without having to query term.GetSize.
func NewPrinter(out io.Writer, verboseLevel int) *Printer {
	return &Printer{
		out:          out,
		verboseLevel: verboseLevel,
	}
}

// PrintTaskHeader writes the banner for a single module: the module
// name, a "•"-filled separator sized to the terminal width, and the
// current timestamp.
func (p *Printer) PrintTaskHeader(moduleName string) {
	if p.termWidth == 0 {
		p.termWidth = detectTerminalWidth()
	}

	delim := ""
	delimLen := p.termWidth - len(moduleName) - 2 - timestampWidth
	if delimLen > 0 {
		delim = strings.Repeat("•", delimLen)
	}
	fmt.Fprintf(p.out,
		"%s %s %s\n",
		theme.Task().Render(moduleName),
		theme.Dim().Render(delim),
		theme.Dim().Render(time.Now().Format("2006-01-02 15:04:05")),
	)
}

// PrintTaskSeparator writes a blank line between consecutive module
// outputs. Callers invoke this after each module finishes to keep
// successive banners visually separated.
func (p *Printer) PrintTaskSeparator() {
	fmt.Fprintln(p.out)
}

// PrintHostResult writes one host's header line and, when res is
// non-nil, walks its key/value pairs through PrintKeyValue. prefix is
// prepended to the host label so output nests visually under the
// task header.
func (p *Printer) PrintHostResult(prefix string, h remote.Host, res *gtypes.OrderedMap[string, string]) {
	hostStr := remote.HostStr(h, prefix)
	fmt.Fprintf(p.out, "%s  %s\n",
		theme.HostLine().Render(hostStr),
		theme.Timestamp().Render(time.Now().Format("[2006-01-02 15:04:05]")),
	)

	if res != nil {
		_ = res.ForEach(func(key, val string) error {
			p.PrintKeyValue(prefix, key, val)
			return nil
		})
	}
}

// PrintKeyValue renders a single key/value entry. Empty values are
// skipped. At verbose level 0, single-line values are rendered inline
// as "key: value"; multi-line values and any value at higher verbosity
// are rendered as a header line followed by indented body lines.
func (p *Printer) PrintKeyValue(prefix string, key string, val string) {
	if val == "" {
		return
	}

	indent := prefix + "    "
	if p.verboseLevel == 0 && isSingleLineValue(val) {
		keyLabel := theme.ResultLabel().Render(fmt.Sprintf("%s:", key))
		fmt.Fprintf(p.out, "%s%s %s\n", indent, keyLabel, strings.ReplaceAll(val, "\n", ""))
		return
	}

	fmt.Fprintf(p.out, "%s%s:\n", indent, theme.ResultLabel().Render(key))
	scanner := bufio.NewScanner(strings.NewReader(val))
	for scanner.Scan() {
		fmt.Fprintf(p.out, "%s    %s\n", indent, scanner.Text())
	}
}

// isSingleLineValue reports whether val contains at most one trailing
// newline, i.e. it fits on a single visual line. A value with no
// newlines or with only a single trailing newline qualifies.
func isSingleLineValue(val string) bool {
	nl := strings.IndexByte(val, '\n')
	return nl == -1 || nl == len(val)-1
}

// detectTerminalWidth returns stdout's width when it is a tty, and
// the package fallback otherwise.
func detectTerminalWidth() int {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			return w
		}
	}
	return fallbackTerminalWidth
}
