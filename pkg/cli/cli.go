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

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/gtypes"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/operation"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/prompt"
	promptconstants "github.com/cqroot/prompt/constants"
	"github.com/cqroot/prompt/input"
)

// defaultConfigFile is the path used when no config file is explicitly specified.
const defaultConfigFile = "./" + constants.DefaultConfigFile

// Cli provides an interactive command-line interface for remote operations.
type Cli struct {
	configFile      string
	optVerboseLevel int
	optMaxProcs     int
}

// New creates a new Cli instance with the given options.
func New(opts ...Option) *Cli {
	c := Cli{
		optVerboseLevel: 0,
		optMaxProcs:     1,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// MinopTheme customizes the prompt appearance.
func MinopTheme(msg string, state prompt.State, model string) string {
	s := strings.Builder{}

	s.WriteString(promptconstants.DefaultNormalPromptSuffixStyle.Render(msg))
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(" › "))
	s.WriteString(model)
	if state != prompt.StateNormal {
		s.WriteString("\n")
	}

	return s.String()
}

// Output styling for help display
var (
	helpTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	helpKeyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	helpValStyle   = lipgloss.NewStyle().Faint(true)
)

// ShowHelp displays the available CLI commands.
func ShowHelp() {
	fmt.Print(helpTitleStyle.Render("\nMINOP CLI COMMANDS\n"))

	helpEntries := gtypes.NewOrderedMap[string, string]()
	helpEntries.Put("exit", "Quit minop")
	helpEntries.Put("quit", "Quit minop")
	helpEntries.Put("help", "Show help output")
	_ = helpEntries.ForEach(func(k, v string) error {
		fmt.Printf("    %s    %s\n", helpKeyStyle.Render(k), helpValStyle.Render(v))
		return nil
	})

	fmt.Println()
}

// Run starts the interactive CLI loop, reading commands from stdin and
// executing them on the configured remote hosts. It returns when the user
// quits or encounters an error.
func (c Cli) Run() error {
	configFile := c.configFile
	if configFile == "" {
		configFile = defaultConfigFile
	}

	e := executor.New(executor.WithMaxProcs(c.optMaxProcs))
	hostGroup, _, err := e.LoadConfig(configFile)
	if err != nil {
		return err
	}
	pool := remote.NewHostPool()

	for {
		val, err := prompt.New(prompt.WithTheme(MinopTheme)).Ask("MINOP").
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
			ShowHelp()
			continue
		}

		op, err := operation.NewOpShell(operation.Input{
			Shell: val,
		})
		if err != nil {
			return err
		}
		op.SetRole(constants.RoleAll)

		err = e.ExecuteOperation(hostGroup, pool, op)
		if err != nil {
			return err
		}

		fmt.Println("")
	}
}
