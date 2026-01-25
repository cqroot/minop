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
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/operation"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/prompt"
	promptconstants "github.com/cqroot/prompt/constants"
)

type Cli struct {
	optVerboseLevel int
	optMaxProcs     int
}

func New(opts ...Option) *Cli {
	e := Cli{
		optVerboseLevel: 0,
		optMaxProcs:     1,
	}

	for _, opt := range opts {
		opt(&e)
	}

	return &e
}

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

func (c Cli) Run() error {
	hostGroup, err := remote.HostsFromFile(filepath.Join(".", constants.HostFileName))
	if err != nil {
		return err
	}
	pool := remote.NewHostPool()
	e := executor.New(executor.WithMaxProcs(c.optMaxProcs))

	for true {
		val, err := prompt.New(prompt.WithTheme(MinopTheme)).Ask("MINOP").Input("Remote command")
		if err != nil {
			if errors.Is(err, prompt.ErrUserQuit) {
				return nil
			} else {
				return err
			}
		}

		if val == "exit" || val == "quit" {
			return nil
		}

		op, err := operation.NewOpShell(operation.Input{
			Shell: val,
		})
		if err != nil {
			return err
		}
		op.SetRole("all")

		err = e.ExecuteOperation(hostGroup, pool, op)
		if err != nil {
			return err
		}

		fmt.Println("")
	}
	return nil
}
