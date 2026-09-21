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

package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/module"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/minop/pkg/theme"
	"github.com/cqroot/prompt"
	promptconstants "github.com/cqroot/prompt/constants"
	"github.com/cqroot/prompt/input"
)

// defaultHostsFile is the path used when no hosts file is explicitly
// specified via WithHostsFile.
const (
	defaultHostsFile = "./" + constants.DefaultHostsFile
)

// Cli provides an interactive command-line interface for running
// shell commands across remote hosts.
type Cli struct {
	hostsFile   string
	optMaxProcs int
}

// New creates a new Cli instance with the given options.
func New(opts ...Option) *Cli {
	c := Cli{
		optMaxProcs: constants.DefaultMaxProcs,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// minopTheme customizes the prompt appearance for the interactive REPL.
func minopTheme(msg string, state prompt.State, model string) string {
	s := strings.Builder{}

	s.WriteString(promptconstants.DefaultNormalPromptSuffixStyle.Render(msg))
	s.WriteString(theme.PromptArrow().Render(" › "))
	s.WriteString(model)
	if state != prompt.StateNormal {
		s.WriteString("\n")
	}

	return s.String()
}

// showHelp prints the available CLI commands and is invoked from Run
// when the user types "help" or "h" at the prompt.
func showHelp() {
	fmt.Print(theme.Title().Render("\nMINOP CLI COMMANDS\n"))

	helpEntries := gtypes.NewOrderedMap[string, string]()
	helpEntries.Put("exit", "Quit minop")
	helpEntries.Put("quit", "Quit minop")
	helpEntries.Put("help", "Show help output")
	_ = helpEntries.ForEach(func(k, v string) error {
		fmt.Printf("    %s    %s\n", theme.HelpKey().Render(k), theme.Dim().Render(v))
		return nil
	})

	fmt.Println()
}

// Run starts the interactive CLI loop, reading commands from stdin and
// executing them on the configured remote hosts. It returns when the user
// quits or encounters an error.
func (c Cli) Run() error {
	hostsFile := c.hostsFile
	if hostsFile == "" {
		hostsFile = defaultHostsFile
	}

	e := executor.New(executor.WithMaxProcs(c.optMaxProcs))
	hostGroup, err := executor.LoadHostsFile(hostsFile)
	if err != nil {
		return err
	}
	pool := remote.NewHostPool()

	for {
		val, err := prompt.New(prompt.WithTheme(minopTheme)).Ask("MINOP").
			Input("", input.WithWidth(0), input.WithCharLimit(0))
		if err != nil {
			if errors.Is(err, prompt.ErrUserQuit) {
				return nil
			} else {
				return err
			}
		}

		trimmed := strings.Trim(val, " ")

		if trimmed == "" {
			continue
		}

		if trimmed == "exit" || trimmed == "quit" || trimmed == "q" {
			return nil
		}

		if trimmed == "help" || trimmed == "h" {
			showHelp()
			continue
		}

		m, err := module.NewShell(module.Task{
			Shell: val,
		})
		if err != nil {
			return err
		}
		m.SetGroup(constants.GroupAll)

		err = e.ExecuteOnHosts("", hostGroup, pool, m)
		if err != nil {
			return err
		}

		fmt.Println("")
	}
}
