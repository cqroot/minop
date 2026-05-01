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
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/minop/pkg/executor"
	"github.com/spf13/cobra"
)

// Output styling for host tree display
var treeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("212"))

var hostStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

var groupStyle = lipgloss.NewStyle().
	Bold(true).Faint(false).Foreground(lipgloss.Color("12"))

// RunHostCmd displays all configured hosts in a tree format.
func RunHostCmd(cmd *cobra.Command, args []string) {
	e := executor.New(
		executor.WithVerboseLevel(flagVerboseLevel),
		executor.WithMaxProcs(flagMaxProcs))

	hostGroup, _, err := e.LoadConfig(flagConfigFile)
	CheckErr(err)

	groups := make([]string, 0, len(hostGroup))
	for group := range hostGroup {
		groups = append(groups, group)
	}

	_, _ = fmt.Fprintln(os.Stdout)
	for idx, group := range groups {
		hosts := hostGroup[group]
		_, _ = fmt.Fprintf(os.Stdout, "  %s\n", groupStyle.Render("• "+group))
		for hostIdx, host := range hosts {
			var branch string
			if hostIdx == len(hosts)-1 {
				branch = treeStyle.Render("  └──")
			} else {
				branch = treeStyle.Render("  ├──")
			}
			_, _ = fmt.Fprintf(os.Stdout, "%s %s\n",
				branch,
				hostStyle.Render(fmt.Sprintf("%s@%s:%d", host.User, host.Address, host.Port)))
		}
		if idx < len(groups)-1 {
			_, _ = fmt.Fprintln(os.Stdout)
		}
	}
}

// NewHostCmd creates the host command that lists all configured hosts.
func NewHostCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "host",
		Short: "List all hosts",
		Long:  "List all hosts.",
		Run:   RunHostCmd,
	}

	return &c
}
