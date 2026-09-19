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
	"os"
	"strings"
	"time"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/theme"
	"golang.org/x/term"
)

const (
	timestampWidth = 19
	// fallbackTerminalWidth is used when stdout is not a tty (e.g.
	// output is piped to a file or another process). 500 is wide
	// enough to keep task headers on one line for most terminal
	// widths while still being sensible for non-tty output.
	fallbackTerminalWidth = 500
)

func getTerminalWidth() int {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			return w
		}
	}
	return fallbackTerminalWidth
}

func printTaskHeader(opName string, termWidth int) {
	delim := ""
	delimLen := termWidth - len(opName) - 2 - timestampWidth
	if delimLen > 0 {
		delim = strings.Repeat("•", delimLen)
	}
	fmt.Printf(
		"%s %s %s\n",
		theme.Task().Render(opName),
		theme.Dim().Render(delim),
		theme.Dim().Render(time.Now().Format("2006-01-02 15:04:05")),
	)
}

func printHostResult(prefix string, h remote.Host, res *gtypes.OrderedMap[string, string], verboseLevel int) {
	hostStr := fmt.Sprintf("%s%s@%s:%d", prefix, h.Username, h.Address, h.Port)
	fmt.Printf("%s  %s\n", theme.HostLine().Render(hostStr),
		theme.Timestamp().Render(time.Now().Format("[2006-01-02 15:04:05]")))

	if res != nil {
		_ = res.ForEach(func(key, val string) error {
			printKeyValue(prefix, key, val, verboseLevel)
			return nil
		})
	}
}

func printKeyValue(prefix string, key string, val string, verboseLevel int) {
	if val == "" {
		return
	}

	indent := fmt.Sprintf("%s    ", prefix)
	if verboseLevel == 0 && (strings.IndexByte(val, '\n') == -1 || strings.IndexByte(val, '\n') == len(val)-1) {
		fmt.Printf("%s%s %s\n", indent, theme.ResultLabel().Render(fmt.Sprintf("%s:", key)), strings.ReplaceAll(val, "\n", ""))
	} else {
		fmt.Printf("%s%s:\n", indent, theme.ResultLabel().Render(key))
		scanner := bufio.NewScanner(strings.NewReader(val))
		for scanner.Scan() {
			fmt.Printf("%s    %s\n", indent, scanner.Text())
		}
	}
}
