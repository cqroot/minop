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
	"os"

	"github.com/cqroot/minop/pkg/hosts"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func RunHostCmd(cmd *cobra.Command, args []string) {
	hostmgr, err := hosts.New("./host.list")
	cobra.CheckErr(err)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"Group", "Address", "User", "Port"})

	isFirst := true
	for group, hosts := range hostmgr.Hosts {
		for _, host := range hosts {
			if isFirst {
				t.AppendRow(table.Row{group, host.Address, host.User, host.Port})
				isFirst = false
			} else {
				t.AppendRow(table.Row{"", host.Address, host.User, host.Port})
			}
		}
		t.AppendSeparator()
		isFirst = true
	}
	t.Render()
}

func NewHostCmd() *cobra.Command {
	hostCmd := cobra.Command{
		Use:   "host",
		Short: "List all hosts.",
		Long:  "List all hosts.",
		Run:   RunHostCmd,
	}

	return &hostCmd
}
