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
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/minop/pkg/cli"
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
	c := cli.New(logger, cli.WithMaxProcs(flagMaxProcs), cli.WithVerboseLeve(flagVerboseLevel))
	err := c.Run()
	CheckErr(err)
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
