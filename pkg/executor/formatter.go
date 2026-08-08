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

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/remote"
	"golang.org/x/term"
)

const (
	timestampWidth = 19
)

var (
	labelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	taskStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dimStyle       = lipgloss.NewStyle().Faint(true)
	hostStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	timestampStyle = lipgloss.NewStyle().Faint(true)
)

func getTerminalWidth() int {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			return w
		}
	}
	return 500
}

func printTaskHeader(opName string, termWidth int) {
	delim := ""
	delimLen := termWidth - len(opName) - 2 - timestampWidth
	if delimLen > 0 {
		delim = strings.Repeat("•", delimLen)
	}
	fmt.Printf("%s %s %s\n",
		taskStyle.Render(opName),
		dimStyle.Render(delim),
		dimStyle.Render(time.Now().Format("2006-01-02 15:04:05")),
	)
}

func printHostResult(prefix string, h remote.Host, res *gtypes.OrderedMap[string, string], verboseLevel int) {
	hostStr := fmt.Sprintf("%s%s@%s:%d", prefix, h.User, h.Address, h.Port)
	fmt.Printf("%s  %s\n", hostStyle.Render(hostStr),
		timestampStyle.Render(time.Now().Format("[2006-01-02 15:04:05]")))

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
		fmt.Printf("%s%s %s\n", indent, labelStyle.Render(fmt.Sprintf("%s:", key)), strings.ReplaceAll(val, "\n", ""))
	} else {
		fmt.Printf("%s%s:\n", indent, labelStyle.Render(key))
		scanner := bufio.NewScanner(strings.NewReader(val))
		for scanner.Scan() {
			fmt.Printf("%s    %s\n", indent, scanner.Text())
		}
	}
}
