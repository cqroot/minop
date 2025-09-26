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

package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/minop/pkg/action"
	"github.com/cqroot/minop/pkg/action/command"
	"github.com/cqroot/minop/pkg/constants"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/remote"
	"github.com/cqroot/prompt"
	promptconstants "github.com/cqroot/prompt/constants"
	"github.com/spf13/cobra"
)

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

func RunCliCmd(cmd *cobra.Command, args []string) {
	hostGroup, err := host.Read(filepath.Join(".", constants.HostFileName))
	CheckErr(err)
	rgs := make(map[host.Host]*remote.Remote)

	for true {
		val, err := prompt.New(prompt.WithTheme(MinopTheme)).Ask("MINOP").Input("Remote command")
		if err != nil {
			if errors.Is(err, prompt.ErrUserQuit) {
				return
			} else {
				CheckErr(err)
			}
		}

		if val == "exit" || val == "quit" {
			return
		}

		act, err := command.New(map[string]string{
			"command": val,
		})
		CheckErr(err)

		actWrapper := *action.New("command", "all", act)
		e := executor.New(logger, executor.WithMaxProcs(flagMaxProcs))

		err = e.PrintActionResult(hostGroup, &rgs, actWrapper, "")
		CheckErr(err)

		fmt.Println("")
	}
}

func NewCliCmd() *cobra.Command {
	cliCmd := cobra.Command{
		Use:   "cli",
		Short: "",
		Long:  "",
		Run:   RunCliCmd,
	}

	return &cliCmd
}
